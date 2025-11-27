# Docker Deployment Guide

AionMCP can be deployed as a Docker container for easy deployment across different environments.

---

## 🐳 Using Pre-built Image from GitHub Container Registry

### Pull and Run

```bash
# Pull the latest image
docker pull ghcr.io/kiransth77/aionmcp:latest

# Run the container
docker run -d \
  --name aionmcp \
  -p 8080:8080 \
  -p 9090:9090 \
  -v $(pwd)/examples/specs:/root/examples/specs:ro \
  ghcr.io/kiransth77/aionmcp:latest
```

### Verify It's Running

```bash
# Check container status
docker ps | grep aionmcp

# Test the health endpoint
curl http://localhost:8080/api/v1/health

# List available tools
curl http://localhost:8080/api/v1/tools | jq
```

---

## 🔧 Build Locally

### Build from Source

```bash
# Clone repository
git clone https://github.com/kiransth77/aionmcp.git
cd aionmcp

# Build image
docker build -t aionmcp:latest .

# Tag for registry (optional)
docker tag aionmcp:latest ghcr.io/kiransth77/aionmcp:latest
```

### Run Locally Built Image

```bash
docker run -d \
  --name aionmcp-local \
  -p 8080:8080 \
  -p 9090:9090 \
  -v $(pwd)/examples/specs:/root/examples/specs:ro \
  aionmcp:latest
```

---

## 🐋 Using Docker Compose

### Quick Start

```bash
# Start services
docker-compose up -d

# View logs
docker-compose logs -f aionmcp

# Stop services
docker-compose down
```

### Configuration

Edit `docker-compose.yml` to customize:
- Port mappings
- Environment variables
- Volume mounts
- Resource limits

---

## 🔗 Available Tags

When pulling from GitHub Container Registry, use these tags:

| Tag | Description |
|-----|-------------|
| `latest` | Latest release (main branch) |
| `v0.1.0` | Specific version |
| `v0.1` | Latest patch for v0.1 |
| `main` | Latest main branch |
| `main-abc123` | Specific commit |

### Examples

```bash
# Latest version
docker pull ghcr.io/kiransth77/aionmcp:latest

# Specific version
docker pull ghcr.io/kiransth77/aionmcp:v0.1.0

# Latest major.minor
docker pull ghcr.io/kiransth77/aionmcp:v0.1

# Main branch
docker pull ghcr.io/kiransth77/aionmcp:main
```

---

## 📝 Environment Variables

Configure the container using environment variables:

```bash
docker run -d \
  --name aionmcp \
  -e LOG_LEVEL=debug \
  -e API_HOST=0.0.0.0 \
  -e API_PORT=8080 \
  -e GRPC_HOST=0.0.0.0 \
  -e GRPC_PORT=9090 \
  -p 8080:8080 \
  -p 9090:9090 \
  ghcr.io/kiransth77/aionmcp:latest
```

### Available Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_LEVEL` | `info` | Logging level (debug, info, warn, error) |
| `API_HOST` | `0.0.0.0` | HTTP API host |
| `API_PORT` | `8080` | HTTP API port |
| `GRPC_HOST` | `0.0.0.0` | gRPC host |
| `GRPC_PORT` | `9090` | gRPC port |

---

## 💾 Volume Mounts

### Specs Directory

```bash
# Mount your specs directory
docker run -d \
  --name aionmcp \
  -v /path/to/your/specs:/root/examples/specs:ro \
  -p 8080:8080 \
  ghcr.io/kiransth77/aionmcp:latest
```

### Data Directory

```bash
# Mount data directory for persistence
docker run -d \
  --name aionmcp \
  -v aionmcp-data:/root/data \
  -p 8080:8080 \
  ghcr.io/kiransth77/aionmcp:latest
```

### Create Named Volume

```bash
# Create a named volume
docker volume create aionmcp-data

# Use it
docker run -d \
  --name aionmcp \
  -v aionmcp-data:/root/data \
  ghcr.io/kiransth77/aionmcp:latest
```

---

## 🌐 Port Mappings

| Port | Protocol | Description |
|------|----------|-------------|
| 8080 | HTTP | REST API |
| 9090 | gRPC | gRPC API |

### Custom Ports

```bash
# Map to different host ports
docker run -d \
  --name aionmcp \
  -p 3000:8080 \
  -p 5000:9090 \
  ghcr.io/kiransth77/aionmcp:latest
```

---

## 🔍 Monitoring

### Check Container Status

```bash
# Show running containers
docker ps

# Show container logs
docker logs aionmcp

# Follow logs
docker logs -f aionmcp

# Show last 50 lines
docker logs --tail 50 aionmcp
```

### Execute Commands in Container

```bash
# Interactive shell
docker exec -it aionmcp sh

# Run command
docker exec aionmcp ls -la /root
```

### Health Check

```bash
# Test health endpoint
curl http://localhost:8080/api/v1/health

# Response should be:
# {"status":"ok"}
```

---

## 🚀 Production Deployment

### Kubernetes Example

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: aionmcp
spec:
  replicas: 2
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
        image: ghcr.io/kiransth77/aionmcp:v0.1.0
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 9090
          name: grpc
        env:
        - name: LOG_LEVEL
          value: "info"
        livenessProbe:
          httpGet:
            path: /api/v1/health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /api/v1/health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
        volumeMounts:
        - name: specs
          mountPath: /root/examples/specs
          readOnly: true
      volumes:
      - name: specs
        configMap:
          name: aionmcp-specs
---
apiVersion: v1
kind: Service
metadata:
  name: aionmcp-service
spec:
  selector:
    app: aionmcp
  ports:
  - port: 8080
    name: http
    targetPort: 8080
  - port: 9090
    name: grpc
    targetPort: 9090
  type: LoadBalancer
```

### Docker Swarm Example

```bash
# Initialize swarm
docker swarm init

# Create service
docker service create \
  --name aionmcp \
  --publish 8080:8080 \
  --publish 9090:9090 \
  --replicas 2 \
  ghcr.io/kiransth77/aionmcp:latest

# Scale service
docker service scale aionmcp=3

# View services
docker service ls

# View service logs
docker service logs aionmcp
```

---

## 🔒 Security Considerations

### Non-Root User (Future)

The current image runs as root. For production, consider:

```dockerfile
RUN addgroup -g 1000 aionmcp && \
    adduser -D -u 1000 -G aionmcp aionmcp

USER aionmcp
```

### Read-Only Volumes

```bash
# Mount specs as read-only
docker run -d \
  -v /path/to/specs:/root/examples/specs:ro \
  ghcr.io/kiransth77/aionmcp:latest
```

### Network Isolation

```bash
# Create custom network
docker network create aionmcp-net

# Run container on network
docker run -d \
  --network aionmcp-net \
  --name aionmcp \
  ghcr.io/kiransth77/aionmcp:latest
```

---

## 🐛 Troubleshooting

### Container Won't Start

```bash
# Check logs
docker logs aionmcp

# Inspect container
docker inspect aionmcp

# Try running interactively
docker run -it ghcr.io/kiransth77/aionmcp:latest sh
```

### Port Already in Use

```bash
# Find what's using port 8080
lsof -i :8080

# Use different port
docker run -d -p 3000:8080 ghcr.io/kiransth77/aionmcp:latest
```

### Permission Denied

```bash
# Check volume permissions
docker run -it \
  -v /path/to/specs:/root/examples/specs:ro \
  ghcr.io/kiransth77/aionmcp:latest sh -c "ls -la /root/examples/specs"
```

### Out of Memory

```bash
# Increase memory limit
docker run -d \
  -m 512m \
  ghcr.io/kiransth77/aionmcp:latest
```

---

## 📚 Resources

- [Docker Documentation](https://docs.docker.com/)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [Kubernetes Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/)
- [Docker Compose](https://docs.docker.com/compose/)

---

## 📞 Support

For issues or questions:
- [GitHub Issues](https://github.com/kiransth77/aionmcp/issues)
- [GitHub Discussions](https://github.com/kiransth77/aionmcp/discussions)
