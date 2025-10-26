# AionMCP - Autonomous Go MCP Server

## Overview
AionMCP is an autonomous Go-based Model Context Protocol (MCP) server that dynamically imports API specifications (OpenAPI, GraphQL, AsyncAPI) and transforms them into MCP tools. It features self-learning capabilities, context-awareness, and autonomous documentation generation.

## Architecture
The project follows Clean/Hexagonal architecture principles:

```
cmd/                     # Application entry points
├── server/             # Main MCP server application

internal/               # Private application code
├── core/               # Core business logic
├── adapters/           # External integrations
├── registry/           # Tool registration system
└── selflearn/          # Learning and reflection engine

pkg/                    # Public library code
└── agent/              # Agent integration APIs

docs/                   # Documentation
data/                   # Data storage
└── feedback/           # Learning data storage
```

## Current Status
**Iteration 0**: ✅ **COMPLETED**
- Project scaffolding and structure
- Base MCP server with HTTP and gRPC endpoints
- Initial tool registry with built-in tools
- Configuration management with Viper
- Structured logging with Zap
- Graceful shutdown handling

## Quick Start

### Prerequisites
- Go 1.21 or later
- Make sure `$GOPATH/bin` is in your `$PATH`

### Running the Server
```bash
# Build and run
go run cmd/server/main.go

# Or build first
go build -o bin/aionmcp cmd/server/main.go
./bin/aionmcp
```

The server will start on:
- HTTP: `http://localhost:8080`
- gRPC: `localhost:9090`

### API Endpoints

#### Health Check
```bash
curl http://localhost:8080/api/v1/health
```

#### List Available Tools
```bash
curl http://localhost:8080/api/v1/mcp/tools
```

#### Tool Invocation (Echo Tool)
```bash
curl -X POST http://localhost:8080/api/v1/mcp/tools/echo/invoke \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello, AionMCP!"}'
```

## Configuration
Configuration can be provided via:
1. `config.yaml` file in the current directory or `./config/` subdirectory
2. Environment variables with `AIONMCP_` prefix

### Default Configuration
```yaml
server:
  port: 8080
  grpc_port: 9090

mcp:
  protocol_version: "1.0"

storage:
  type: "boltdb"
  path: "./data/aionmcp.db"

log:
  level: "info"
  format: "json"
```

### Environment Variables
```bash
export AIONMCP_SERVER_PORT=8080
export AIONMCP_LOG_LEVEL=debug
```

## Built-in Tools
The server includes several built-in tools for testing and system information:

1. **echo** - Echoes back input messages for testing
2. **status** - Returns registry and server status information

## Next Iterations
- **Iteration 1**: Spec Importer Layer (OpenAPI, GraphQL, AsyncAPI)
- **Iteration 2**: Self-Learning & Reflection Engine
- **Iteration 3**: Autonomous Documentation
- **Iteration 4**: Agent Integration Layer
- **Iteration 5**: Competitive Analysis
- **Iteration 6**: Task Orchestration & Adaptive Behavior

## Development

### Adding New Tools
Tools must implement the `Tool` interface:

```go
type Tool interface {
    Name() string
    Description() string
    Execute(input any) (any, error)
    Metadata() ToolMetadata
}
```

### Testing
```bash
# Run tests
go test ./...

# Run with coverage
go test -v -cover ./...
```

## Contributing
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run tests and ensure they pass
6. Submit a pull request

## License
[License information to be added]