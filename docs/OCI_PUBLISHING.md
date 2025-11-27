# Publishing OCI Image to GitHub Container Registry

This guide explains how to publish AionMCP as an OCI image to GitHub Container Registry (ghcr.io) and update the MCP Registry.

---

## 🐳 Automated Publishing via GitHub Actions

### Trigger Automatic Build

The GitHub Actions workflow (`publish-oci.yml`) automatically builds and publishes the image when:

1. **Push to main branch** - Auto-builds on any push
2. **Push a tag** - Auto-builds on semantic version tags (e.g., `v0.1.0`, `v0.2.0`)
3. **Manual trigger** - Via `workflow_dispatch` from Actions tab

### Available Tags

When you push changes, the following tags are automatically created:

| Tag | When Created |
|-----|--------------|
| `latest` | Every push to main |
| `main` | Every push to main |
| `main-{commit-sha}` | Every push to main |
| `v0.1.0` | When you push tag `v0.1.0` |
| `v0.1` | When you push tag `v0.1.0` (major.minor) |
| `{branch-name}` | When you push to any branch |

### View Builds

1. Go to your GitHub repository
2. Click "Actions" tab
3. Select "Publish OCI Image to GitHub Container Registry"
4. View build logs and status

---

## 🚀 How to Publish (Step by Step)

### Option 1: Automatic via Push to Main (Recommended)

```bash
# Make changes
git add .
git commit -m "feat: your feature"

# Push to main - this automatically triggers the build
git push origin main

# Check Actions tab to monitor the build
```

### Option 2: Manual via Tag Release

```bash
# Create an annotated tag
git tag -a v0.1.0 -m "Release version 0.1.0"

# Push the tag - this triggers the build
git push origin v0.1.0

# Or push all tags
git push origin --tags
```

### Option 3: Manual via Actions Workflow Dispatch

1. Go to GitHub repository → Actions
2. Select "Publish OCI Image to GitHub Container Registry"
3. Click "Run workflow"
4. Select branch (default: main)
5. Click "Run workflow"

---

## 📝 Image Location

Once published, the image is available at:

```
ghcr.io/kiransth77/aionmcp:latest
ghcr.io/kiransth77/aionmcp:v0.1.0
ghcr.io/kiransth77/aionmcp:v0.1
ghcr.io/kiransth77/aionmcp:main
```

### Pull the Image

```bash
# Pull latest
docker pull ghcr.io/kiransth77/aionmcp:latest

# Pull specific version
docker pull ghcr.io/kiransth77/aionmcp:v0.1.0

# Pull main branch
docker pull ghcr.io/kiransth77/aionmcp:main
```

---

## 🔐 Authentication for Private Registry

The GitHub Actions workflow uses `GITHUB_TOKEN` which is automatically available.

To pull the image manually:

### Authenticate with GitHub Token

```bash
# Option 1: Use docker login
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin

# Option 2: Create a personal access token (PAT)
# 1. Go to Settings → Developer settings → Personal access tokens
# 2. Create token with 'packages:read' scope
# 3. Use token as password

echo $PAT_TOKEN | docker login ghcr.io -u USERNAME --password-stdin
```

### View Available Images

```bash
# List all packages
curl https://api.github.com/users/kiransth77/packages \
  -H "Authorization: token $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json"

# Or visit
# https://github.com/kiransth77?tab=packages&repo_name=aionmcp
```

---

## 🔍 Verify Published Image

### Check Image on GitHub

1. Go to repository → "Packages" section (right sidebar)
2. Look for "aionmcp" package
3. View available versions and tags

### Test Image Locally

```bash
# Pull the image
docker pull ghcr.io/kiransth77/aionmcp:latest

# Run it
docker run -d \
  --name test-aionmcp \
  -p 8080:8080 \
  ghcr.io/kiransth77/aionmcp:latest

# Test it works
curl http://localhost:8080/api/v1/health

# Cleanup
docker stop test-aionmcp
docker rm test-aionmcp
```

### Get Image Digest

```bash
# Get SHA256 digest
docker inspect --format='{{.RepoDigests}}' ghcr.io/kiransth77/aionmcp:latest

# Or query registry API
curl https://ghcr.io/v2/kiransth77/aionmcp/manifests/latest \
  -H "Accept: application/vnd.docker.distribution.manifest.v2+json" \
  -H "Authorization: Bearer $GITHUB_TOKEN"
```

---

## 📊 Update MCP Registry

Once the Docker image is published to GitHub Container Registry, you can update the MCP Registry entry to include the OCI package.

### Update server.json

Add the OCI package reference:

```json
{
  "$schema": "https://static.modelcontextprotocol.io/schemas/2025-10-17/server.schema.json",
  "name": "io.github.kiransth77/aionmcp",
  "title": "AionMCP",
  "description": "Dynamic API tool generator for OpenAPI, GraphQL, and AsyncAPI specifications",
  "repository": {
    "url": "https://github.com/kiransth77/aionmcp",
    "source": "github"
  },
  "version": "0.1.0",
  "author": {"name": "Kiran Shrestha"},
  "license": "MIT",
  "homepage": "https://github.com/kiransth77/aionmcp"
}
```

Note: OCI package support in server.json is pending official registry support.

### Verify Registry Entry

```bash
# Query the registry
curl "https://registry.modelcontextprotocol.io/v0.1/servers?search=aionmcp"

# Should show server is active
```

---

## 🔄 Workflow Details

### Build Process

The GitHub Actions workflow does the following:

1. **Checkout**: Clone repository code
2. **Setup Buildx**: Configure Docker BuildKit for multi-platform builds
3. **Authenticate**: Log in to GitHub Container Registry using `GITHUB_TOKEN`
4. **Extract Metadata**: Generate semantic version tags
5. **Build & Push**: 
   - Build Docker image from Dockerfile
   - Push to `ghcr.io/kiransth77/aionmcp`
   - Tag with version, branch, and commit info
   - Cache layers for faster rebuilds
6. **Summarize**: Create GitHub step summary with pull command

### Build Cache

The workflow uses GitHub Actions cache to speed up builds:

```yaml
cache-from: type=gha
cache-to: type=gha,mode=max
```

This means:
- First build: ~2-3 minutes
- Subsequent builds: ~30-60 seconds (with cache hits)

---

## 🐛 Troubleshooting

### Build Fails

1. Check the Actions tab for error messages
2. Common issues:
   - Go version incompatibility
   - Missing dependencies in Dockerfile
   - Binary build errors

### Image Push Fails

```
Error response from daemon: Error response from registry
```

Possible causes:
- Token expired → Workflow recreates token automatically
- Package size too large → Check image size with `docker history`
- Registry unavailable → Check GitHub status

### Can't Pull Image

```bash
# Error: no access to this repository
```

Solutions:
- Check image visibility (Settings → Packages)
- Ensure you're authenticated:
  ```bash
  docker login ghcr.io
  ```
- Use correct username (lowercase)

### Multi-platform Build Issues

To build for multiple architectures:

```bash
# Verify buildx supports your architecture
docker buildx ls

# Or use docker buildx create --use
docker buildx create --name multiarch --use
```

---

## 📚 Dockerfile Details

The Dockerfile uses a multi-stage build:

```dockerfile
# Stage 1: Builder (compile Go code)
FROM golang:1.21-alpine AS builder
# ... build the binary

# Stage 2: Runtime (lightweight Alpine Linux)
FROM alpine:latest
# ... copy binary and run
```

Benefits:
- Final image size: ~50-60MB (vs 300+MB with single stage)
- Faster deployments
- Smaller attack surface
- No build tools in runtime image

---

## 🚀 Next Steps

After publishing to ghcr.io:

1. **Update Deployment Docs**
   ```bash
   # Add OCI section to DOCKER_DEPLOYMENT.md
   ```

2. **Create Release Notes**
   - Include Docker pull command
   - Document new OCI support

3. **Announce on Channels**
   - GitHub Discussions
   - Community channels
   - Social media

4. **Update Registry**
   - Once registry supports OCI packages
   - Add OCI reference to server.json
   - Republish to MCP Registry

---

## 📖 Useful Commands

### Build Locally

```bash
# Build image
docker build -t aionmcp:local .

# Build for specific architecture
docker buildx build --platform linux/amd64 -t aionmcp:local .

# Build for multiple architectures (requires buildx)
docker buildx build --platform linux/amd64,linux/arm64 -t aionmcp:local .
```

### Inspect Image

```bash
# Show layers and size
docker history ghcr.io/kiransth77/aionmcp:latest

# Inspect image details
docker inspect ghcr.io/kiransth77/aionmcp:latest

# Show image stats
docker images ghcr.io/kiransth77/aionmcp
```

### Manage Tags

```bash
# Tag existing image
docker tag ghcr.io/kiransth77/aionmcp:v0.1.0 ghcr.io/kiransth77/aionmcp:latest

# Push tag
docker push ghcr.io/kiransth77/aionmcp:latest

# Remove local image
docker rmi ghcr.io/kiransth77/aionmcp:latest
```

---

## 📞 Support

For issues or questions:
- Check GitHub Actions logs for build errors
- Review [Docker documentation](https://docs.docker.com/)
- Check [GitHub Container Registry docs](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- Open GitHub issue for problems

---

**Status**: ✅ Automated OCI publishing configured  
**Frequency**: Automatic on every main push + tag creation  
**Registry**: ghcr.io/kiransth77/aionmcp  
**Image Size**: ~50-60MB  
**Supported Architectures**: linux/amd64, linux/arm64 (see Dockerfile)
