# Mobile Deployment Guide

This guide covers deploying and configuring AionMCP server for mobile application access.

## Server Deployment

### Prerequisites

- Go 1.21+ installed
- Domain name (for production)
- SSL/TLS certificate (for HTTPS)
- Reverse proxy (recommended: Nginx or Caddy)

### Basic Deployment

#### 1. Build the Server

```bash
# Clone the repository
git clone https://github.com/kiransth77/aionmcp.git
cd aionmcp

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o bin/aionmcp cmd/server/main.go

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o bin/aionmcp cmd/server/main.go
```

#### 2. Create Configuration

Create `config.yaml`:

```yaml
server:
  port: 8080
  grpc_port: 50051
  host: "0.0.0.0"

learning:
  enabled: true
  sample_rate: 1.0
  retention_days: 30

storage:
  path: "./data/aionmcp.db"

security:
  api_keys:
    - "your-secure-api-key-here"
  cors:
    enabled: true
    allowed_origins:
      - "https://your-mobile-app.com"
      - "http://localhost:3000"  # For development
    allowed_methods:
      - GET
      - POST
      - PUT
      - DELETE
    allowed_headers:
      - Content-Type
      - Authorization
      - X-API-Key

logging:
  level: "info"
  format: "json"
```

#### 3. Run the Server

```bash
./bin/aionmcp --config config.yaml
```

### Docker Deployment

#### Dockerfile

Create `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o aionmcp cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/aionmcp .
COPY --from=builder /app/config.yaml .

EXPOSE 8080 50051

CMD ["./aionmcp", "--config", "config.yaml"]
```

#### Build and Run

```bash
# Build image
docker build -t aionmcp:latest .

# Run container
docker run -d \
  -p 8080:8080 \
  -p 50051:50051 \
  -v $(pwd)/data:/root/data \
  -v $(pwd)/config.yaml:/root/config.yaml \
  --name aionmcp \
  aionmcp:latest
```

#### Docker Compose

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  aionmcp:
    build: .
    ports:
      - "8080:8080"
      - "50051:50051"
    volumes:
      - ./data:/root/data
      - ./config.yaml:/root/config.yaml
    environment:
      - AIONMCP_LOG_LEVEL=info
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8080/api/v1/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

### Cloud Deployment

#### AWS (EC2)

1. **Launch EC2 Instance**
   - Choose Ubuntu 22.04 LTS
   - Instance type: t3.small or larger
   - Configure security groups (ports 80, 443, 8080, 50051)

2. **Install Dependencies**
   ```bash
   sudo apt update
   sudo apt install -y golang-1.21 nginx certbot python3-certbot-nginx
   ```

3. **Deploy Application**
   ```bash
   # Upload binary and config
   scp -i key.pem bin/aionmcp config.yaml ubuntu@your-server.com:~
   
   # SSH to server
   ssh -i key.pem ubuntu@your-server.com
   
   # Run server
   ./aionmcp --config config.yaml
   ```

4. **Setup as Service**

Create `/etc/systemd/system/aionmcp.service`:

```ini
[Unit]
Description=AionMCP Server
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu
ExecStart=/home/ubuntu/aionmcp --config /home/ubuntu/config.yaml
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl daemon-reload
sudo systemctl enable aionmcp
sudo systemctl start aionmcp
```

#### Google Cloud (GKE)

Create `kubernetes/deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: aionmcp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: aionmcp
  template:
    metadata:
      labels:
        app: aionmcp
    spec:
      containers:
      - name: aionmcp
        image: gcr.io/your-project/aionmcp:latest
        ports:
        - containerPort: 8080
        - containerPort: 50051
        env:
        - name: AIONMCP_LOG_LEVEL
          value: "info"
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: aionmcp-service
spec:
  type: LoadBalancer
  ports:
  - port: 80
    targetPort: 8080
    name: http
  - port: 50051
    targetPort: 50051
    name: grpc
  selector:
    app: aionmcp
```

Deploy:
```bash
kubectl apply -f kubernetes/deployment.yaml
```

## Reverse Proxy Configuration

### Nginx

Create `/etc/nginx/sites-available/aionmcp`:

```nginx
# HTTP to HTTPS redirect
server {
    listen 80;
    server_name api.yourdomain.com;
    return 301 https://$server_name$request_uri;
}

# HTTPS server
server {
    listen 443 ssl http2;
    server_name api.yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/api.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.yourdomain.com/privkey.pem;
    
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # REST API
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # CORS headers
        add_header Access-Control-Allow-Origin * always;
        add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS" always;
        add_header Access-Control-Allow-Headers "Content-Type, Authorization, X-API-Key" always;
        
        if ($request_method = 'OPTIONS') {
            return 204;
        }
    }

    # gRPC
    location /grpc/ {
        grpc_pass grpc://localhost:50051;
        grpc_set_header X-Real-IP $remote_addr;
        grpc_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;
    limit_req zone=api_limit burst=20 nodelay;
}
```

Enable and reload:
```bash
sudo ln -s /etc/nginx/sites-available/aionmcp /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### Caddy

Create `Caddyfile`:

```
api.yourdomain.com {
    # Automatic HTTPS
    
    # REST API
    reverse_proxy /api/* localhost:8080
    
    # gRPC
    reverse_proxy /grpc/* h2c://localhost:50051
    
    # CORS
    @cors_preflight {
        method OPTIONS
    }
    handle @cors_preflight {
        header Access-Control-Allow-Origin *
        header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS"
        header Access-Control-Allow-Headers "Content-Type, Authorization, X-API-Key"
        respond 204
    }
    
    header Access-Control-Allow-Origin *
    
    # Rate limiting
    rate_limit {
        zone static_site {
            key {remote_host}
            events 100
            window 1m
        }
    }
}
```

Run Caddy:
```bash
caddy run --config Caddyfile
```

## SSL/TLS Configuration

### Let's Encrypt (Certbot)

```bash
# Install certbot
sudo apt install certbot python3-certbot-nginx

# Obtain certificate
sudo certbot --nginx -d api.yourdomain.com

# Auto-renewal (already configured by certbot)
sudo certbot renew --dry-run
```

### Custom Certificate

```bash
# Generate self-signed certificate (development only)
openssl req -x509 -newkey rsa:4096 \
  -keyout key.pem -out cert.pem \
  -days 365 -nodes \
  -subj "/CN=api.yourdomain.com"

# Install certificate
sudo cp cert.pem /etc/ssl/certs/aionmcp.crt
sudo cp key.pem /etc/ssl/private/aionmcp.key
```

## Security Best Practices

### API Key Management

1. **Generate Secure Keys**
   ```bash
   openssl rand -base64 32
   ```

2. **Environment Variables**
   ```bash
   export AIONMCP_API_KEY="your-secure-key"
   ```

3. **Key Rotation**
   - Rotate keys every 90 days
   - Maintain key versioning
   - Support multiple active keys

### Network Security

1. **Firewall Rules**
   ```bash
   # UFW (Ubuntu)
   sudo ufw allow 80/tcp
   sudo ufw allow 443/tcp
   sudo ufw allow 22/tcp  # SSH only from trusted IPs
   sudo ufw enable
   ```

2. **Fail2Ban**
   ```bash
   sudo apt install fail2ban
   sudo systemctl enable fail2ban
   ```

### Rate Limiting

Configure in `config.yaml`:

```yaml
security:
  rate_limiting:
    enabled: true
    requests_per_minute: 60
    burst: 100
```

## Monitoring

### Health Checks

```bash
# HTTP endpoint
curl https://api.yourdomain.com/api/v1/health

# Response
{
  "status": "healthy",
  "timestamp": 1699564800,
  "version": "0.1.0"
}
```

### Metrics

Access Prometheus metrics:
```
https://api.yourdomain.com/metrics
```

### Logging

Configure structured logging:

```yaml
logging:
  level: "info"
  format: "json"
  output: "/var/log/aionmcp/server.log"
```

## Mobile-Specific Considerations

### CORS Configuration

Enable CORS for mobile web views:

```yaml
security:
  cors:
    enabled: true
    allowed_origins:
      - "*"  # Development only
      - "https://your-app.com"  # Production
    allowed_methods:
      - GET
      - POST
      - PUT
      - DELETE
      - OPTIONS
    allowed_headers:
      - Content-Type
      - Authorization
      - X-API-Key
    max_age: 3600
```

### Mobile Network Optimization

1. **Connection Pooling**
   - Enable HTTP/2
   - Configure keep-alive

2. **Compression**
   - Enable gzip compression
   - Optimize JSON responses

3. **CDN Integration**
   - Use CDN for static content
   - Cache frequently accessed data

### Bandwidth Considerations

Configure response size limits:

```yaml
server:
  max_response_size: 5242880  # 5MB
  compression:
    enabled: true
    level: 6
```

## Troubleshooting

### Common Issues

1. **CORS Errors**
   - Verify CORS configuration
   - Check allowed origins
   - Test OPTIONS preflight

2. **SSL Certificate Issues**
   - Verify certificate validity
   - Check certificate chain
   - Test with `openssl s_client`

3. **Connection Timeouts**
   - Increase timeout values
   - Check firewall rules
   - Verify network connectivity

### Debug Mode

Enable debug logging:

```yaml
logging:
  level: "debug"
```

## Performance Optimization

### Caching

Implement Redis caching:

```yaml
cache:
  enabled: true
  type: "redis"
  redis:
    host: "localhost"
    port: 6379
    db: 0
```

### Load Balancing

Configure multiple instances behind load balancer:

```nginx
upstream aionmcp_backend {
    least_conn;
    server 10.0.1.10:8080;
    server 10.0.1.11:8080;
    server 10.0.1.12:8080;
}

server {
    location /api/ {
        proxy_pass http://aionmcp_backend;
    }
}
```

## Backup and Recovery

### Database Backup

```bash
# Backup BoltDB
cp data/aionmcp.db backups/aionmcp-$(date +%Y%m%d).db

# Automated backup script
#!/bin/bash
BACKUP_DIR="/backups"
DATE=$(date +%Y%m%d-%H%M%S)
cp data/aionmcp.db "$BACKUP_DIR/aionmcp-$DATE.db"
find "$BACKUP_DIR" -name "aionmcp-*.db" -mtime +30 -delete
```

### Configuration Backup

```bash
# Backup configuration
tar czf config-backup.tar.gz config.yaml data/
```

## Next Steps

1. Configure monitoring and alerting
2. Set up automated backups
3. Implement CI/CD pipeline
4. Configure auto-scaling
5. Test disaster recovery procedures

## Resources

- [AionMCP Documentation](../docs/)
- [Mobile Integration Guide](./mobile_integration.md)
- [Nginx Documentation](https://nginx.org/en/docs/)
- [Caddy Documentation](https://caddyserver.com/docs/)
- [Let's Encrypt Documentation](https://letsencrypt.org/docs/)
