# AionMCP - Model-Independent Tool Server for AI Agents

**A unified, model-independent tool server that works with GitHub Copilot, Claude, and any AI agent via standard HTTP REST API.**

## ✨ What is AionMCP?

AionMCP is a **universal tool server** that makes any API (OpenAPI, GraphQL, AsyncAPI) available to **any AI model or agent** through a simple HTTP REST API. No model-specific configuration needed.

```
GitHub Copilot / Claude / Any Agent
         │
         │ HTTP REST API
         ▼
    AionMCP Server
         │
         ▼
    Your APIs (OpenAPI, GraphQL, AsyncAPI)
```

## 🚀 Quick Start (5 Minutes)

### 1. Start the Server
```bash
cd /Users/kiran/Documents/GitHub/aionmcp
./bin/aionmcp-server
```

Or use VS Code: Click ▶️ in AionMCP sidebar

### 2. Import an API Spec
```bash
curl -X POST http://localhost:8080/api/v1/import-spec \
  -H "Content-Type: application/json" \
  -d '{
    "type": "openapi",
    "path": "./examples/specs/petstore.yaml",
    "name": "Petstore API"
  }'
```

### 3. Use with GitHub Copilot
Open GitHub Copilot chat and ask:
```
"List all available pets using the petstore API"
```

Copilot will automatically discover and execute the tool!

## 🎯 Why Model-Independent?

### ✅ Works with Any Agent
- GitHub Copilot
- Claude Desktop
- Custom agents
- Future LLMs
- Internal tools

### ✅ No Configuration
- No manual JSON file editing
- No model-specific setup
- Just HTTP API calls

### ✅ Future-Proof
- Not tied to MCP protocol
- Adapts to new models
- Single integration point

### ✅ Developer-Friendly
- Standard REST API
- Works from any language (Python, TypeScript, Go, Java, etc.)
- Simple JSON format

## 📚 Documentation

### Getting Started
- **[SETUP_GITHUB_COPILOT.md](./docs/SETUP_GITHUB_COPILOT.md)** - Step-by-step setup guide
- **[QUICK_START.md](./docs/QUICK_START.md)** - Quick reference

### Integration & API
- **[GITHUB_COPILOT_INTEGRATION.md](./docs/GITHUB_COPILOT_INTEGRATION.md)** - Full integration guide with code examples
- **[MODEL_INDEPENDENT_TOOLS.md](./docs/MODEL_INDEPENDENT_TOOLS.md)** - Architecture & design philosophy

### Project Overview
- **[IMPLEMENTATION_SUMMARY.md](./docs/IMPLEMENTATION_SUMMARY.md)** - Complete feature overview
- **[v0.1.0_RELEASE_NOTES.md](./docs/v0.1.0_RELEASE_NOTES.md)** - What's new in this version

## 🔌 REST API Endpoints

```
GET  /api/v1/health              # Server health check
GET  /api/v1/tools               # List all tools
GET  /api/v1/tools/{id}          # Get tool details
POST /api/v1/tools/{id}/invoke   # Execute tool
POST /api/v1/import-spec         # Import API specification
```

### Examples

```bash
# List all tools
curl http://localhost:8080/api/v1/tools | jq .

# Execute a tool
curl -X POST http://localhost:8080/api/v1/tools/echo/invoke \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello from AionMCP!"}'
```

## 💻 Integration Examples

### TypeScript
```typescript
const tools = await fetch('http://localhost:8080/api/v1/tools');
const result = await fetch('http://localhost:8080/api/v1/tools/echo/invoke', {
  method: 'POST',
  body: JSON.stringify({ message: 'Hello!' })
});
```

### Python
```python
import requests
tools = requests.get('http://localhost:8080/api/v1/tools').json()
result = requests.post('http://localhost:8080/api/v1/tools/echo/invoke',
                       json={'message': 'Hello!'})
```

### Go
```go
resp, _ := http.Get("http://localhost:8080/api/v1/tools")
resp, _ := http.Post("http://localhost:8080/api/v1/tools/echo/invoke", ...)
```

See [GITHUB_COPILOT_INTEGRATION.md](./docs/GITHUB_COPILOT_INTEGRATION.md) for more examples.

## ✨ Key Features

- ✅ **Dynamic Tool Import** - Import OpenAPI, GraphQL, or AsyncAPI specs
- ✅ **Multi-Protocol Support** - HTTP REST and gRPC access
- ✅ **Agent-Agnostic** - Works with any AI model or custom agent
- ✅ **Self-Learning** - Tracks execution and learns from patterns
- ✅ **Hot-Reload** - File watcher auto-imports spec changes
- ✅ **Auto-Documentation** - Generates docs from API specs
- ✅ **Agent Registration** - Track connected agents
- ✅ **VS Code Extension** - Full server management UI

## 📦 What's Included

```
aionmcp/
├── bin/aionmcp-server              # Main executable
├── cmd/server/                     # Server entry point
├── internal/
│   ├── core/                       # HTTP/gRPC servers & tool registry
│   ├── selflearn/                  # Learning engine & BoltDB storage
│   └── autodocs/                   # Documentation generation
├── pkg/
│   ├── importer/                   # OpenAPI/GraphQL/AsyncAPI parsers
│   ├── agent/                      # Agent registration API
│   └── feedback/                   # Feedback collection
├── vscode-extension/               # VS Code extension
├── docs/                           # Comprehensive documentation
├── examples/specs/                 # Example API specifications
│   ├── petstore.yaml               # OpenAPI 3.0 example
│   ├── blog.graphql                # GraphQL example
│   └── user-events.yaml            # AsyncAPI example
└── go.mod                          # Go dependencies
```

## 🚀 Deployment Options

### Local Development
```bash
./bin/aionmcp-server
# Access: http://localhost:8080
```

### Docker
```bash
docker run -p 8080:8080 aionmcp:latest
```

### Cloud (AWS/GCP/Azure)
Deploy the binary to your cloud provider, access via public URL.

### Embedded
```go
import "github.com/aionmcp/aionmcp/internal/core"
server := core.NewServer(logger)
server.Start()
```

## 🔐 Security

### Development (Current)
- Running on localhost only
- No authentication required

### Production (Recommended)
- Add API key authentication
- Use HTTPS/TLS
- Add rate limiting
- Implement request validation

## 🤖 Supported Models/Agents

- ✅ **GitHub Copilot** - VS Code, GitHub.com
- ✅ **Claude Desktop** - Via REST API
- ✅ **Custom Agents** - Any programming language
- ✅ **Internal Tools** - Direct API access
- ✅ **Future LLMs** - Any HTTP-capable client

## 📋 System Requirements

- **Go 1.21+** (if building from source)
- **8 GB RAM** (recommended)
- **Port 8080** (HTTP), 9090 (gRPC) available
- **Disk space** for API specifications

## 🛠️ Building from Source

```bash
# Clone the repository
git clone https://github.com/kiransth77/aionmcp.git
cd aionmcp

# Build the server
go build -o bin/aionmcp-server ./cmd/server/main.go

# Run
./bin/aionmcp-server
```

## 📖 Documentation Structure

- **Setup & Getting Started**: See [SETUP_GITHUB_COPILOT.md](./docs/SETUP_GITHUB_COPILOT.md)
- **API Integration**: See [GITHUB_COPILOT_INTEGRATION.md](./docs/GITHUB_COPILOT_INTEGRATION.md)
- **Architecture**: See [MODEL_INDEPENDENT_TOOLS.md](./docs/MODEL_INDEPENDENT_TOOLS.md)
- **Quick Reference**: See [QUICK_START.md](./docs/QUICK_START.md)
- **Project Overview**: See [IMPLEMENTATION_SUMMARY.md](./docs/IMPLEMENTATION_SUMMARY.md)

## 🎯 Version & Status

- **Version**: 0.1.0
- **Architecture**: Model-Independent HTTP REST API
- **Status**: Production Ready ✅

## �� Next Steps

1. **Start the server**: `./bin/aionmcp-server`
2. **Import your APIs**: Use VS Code extension or REST API
3. **Use with agents**: GitHub Copilot, Claude, or custom agents
4. **Deploy**: Cloud, Docker, or embedded

## 📞 Support

- **Setup Issues**: See [SETUP_GITHUB_COPILOT.md](./docs/SETUP_GITHUB_COPILOT.md)
- **API Questions**: See [GITHUB_COPILOT_INTEGRATION.md](./docs/GITHUB_COPILOT_INTEGRATION.md)
- **Architecture**: See [MODEL_INDEPENDENT_TOOLS.md](./docs/MODEL_INDEPENDENT_TOOLS.md)
- **Issues**: Report on GitHub
- **Contributing**: Fork and submit PRs

## 📄 License

MIT License - See LICENSE file

---

## 🎉 Get Started Now

```bash
# Start the server
./bin/aionmcp-server

# In another terminal, import Petstore API
curl -X POST http://localhost:8080/api/v1/import-spec \
  -H "Content-Type: application/json" \
  -d '{"type":"openapi","path":"./examples/specs/petstore.yaml"}'

# Open GitHub Copilot and ask:
# "List all available pets using the petstore API"
```

**Questions?** Check the comprehensive documentation in `/docs/`

---

**AionMCP**: One server, any model, any agent. 🚀
