# AionMCP - Autonomous Go MCP Server

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Status](https://img.shields.io/badge/Status-Alpha-orange.svg)](https://github.com/yourusername/aionmcp)

AionMCP is an autonomous Go-based Model Context Protocol (MCP) server that dynamically imports OpenAPI, GraphQL, and AsyncAPI specifications and exposes them as tools to agents. Features self-learning, context-awareness, and auto-documentation using Clean/Hexagonal architecture.

## 🚀 Features

- **Multi-Spec Support**: OpenAPI 3.x, GraphQL schemas, and AsyncAPI specifications
- **Dynamic Tool Generation**: Automatically creates MCP tools from API specifications
- **Hot-Reload**: File system watching for automatic spec reloading
- **Dual Protocol**: HTTP REST and gRPC APIs for agent integration
- **Self-Learning Engine**: Execution feedback and autonomous reflection (planned)
- **Clean Architecture**: Hexagonal/Clean architecture with clear separation of concerns

## 📋 Current Status (Iteration 1 Complete)

✅ **Core Infrastructure**
- MCP server with HTTP and gRPC endpoints
- Tool registry with dynamic registration
- Configuration management with Viper
- Structured logging with Zap

✅ **Spec Importers**
- OpenAPI 3.x importer with full HTTP operation mapping
- GraphQL schema parser with query/mutation tool generation  
- AsyncAPI importer for event-based tools
- File watching for hot-reload capability

✅ **API Endpoints**
- `GET /api/v1/health` - Server health check
- `GET /api/v1/mcp/tools` - List all available tools
- `POST /api/v1/specs` - Load new API specifications
- `GET /api/v1/specs` - List loaded specifications

## 🏗️ Architecture

```
cmd/                     # Application entry points
internal/core/           # Core business logic
├── server.go           # Main HTTP/gRPC server
└── registry.go         # Tool registration system
internal/adapters/       # External integrations (planned)
internal/registry/       # Tool registry implementation  
internal/selflearn/      # Learning and reflection engine (planned)
pkg/importer/           # Specification importers
├── openapi.go          # OpenAPI 3.x support
├── graphql.go          # GraphQL schema support
├── asyncapi.go         # AsyncAPI support
└── watcher.go          # File system watching
pkg/agent/              # Agent integration APIs (planned)
pkg/types/              # Shared type definitions
docs/                   # Auto-generated documentation
examples/specs/         # Example API specifications
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- Git

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/aionmcp.git
   cd aionmcp
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Build the server**:
   ```bash
   # Windows
   go build -o bin/aionmcp.exe cmd/server/main.go
   
   # Linux/macOS
   go build -o bin/aionmcp cmd/server/main.go
   ```

4. **Run the server**:
   ```bash
   # Windows
   ./bin/aionmcp.exe
   
   # Linux/macOS
   ./bin/aionmcp
   ```

### Testing the Server

1. **Check server health**:
   ```bash
   curl http://localhost:8080/api/v1/health
   ```

2. **Load an OpenAPI specification**:
   ```bash
   curl -X POST http://localhost:8080/api/v1/specs \
     -H "Content-Type: application/json" \
     -d '{
       "id": "petstore",
       "type": "openapi", 
       "path": "./examples/specs/petstore.yaml"
     }'
   ```

3. **List generated tools**:
   ```bash
   curl http://localhost:8080/api/v1/mcp/tools
   ```

## 📝 Configuration

AionMCP uses Viper for configuration management. You can configure the server using:

1. **Environment variables** (prefixed with `AIONMCP_`):
   ```bash
   export AIONMCP_SERVER_PORT=8080
   export AIONMCP_LOG_LEVEL=debug
   ```

2. **Configuration file** (`config.yaml`):
   ```yaml
   server:
     port: 8080
     grpc_port: 9090
   log:
     level: info
     format: json
   storage:
     type: boltdb
     path: ./data/aionmcp.db
   ```

## 🔧 Development

### Project Structure
The project follows Clean/Hexagonal architecture principles:

- **Core Domain**: Business logic in `internal/core/`
- **Adapters**: External integrations in `internal/adapters/`
- **Ports**: Interfaces in `pkg/` packages
- **Infrastructure**: Database, file system, network in respective adapters

### Adding New Importers

1. Implement the `Importer` interface in `pkg/importer/`
2. Register the importer in `internal/core/server.go`
3. Add tests in the corresponding `_test.go` file

### Running Tests
```bash
go test ./...
```

## 📚 Documentation

- [Architecture Documentation](docs/architecture.md)
- [Development Iterations](docs/iterations/)
- [AI Coding Guidelines](.github/copilot-instructions.md)

## 🗺️ Roadmap

### Iteration 1: Spec Importer Layer ✅
- Multi-format API spec parsing
- Dynamic tool generation
- Hot-reload capability

### Iteration 2: Self-Learning Engine (Next)
- Execution feedback and reflection
- Error pattern analysis
- Autonomous improvement suggestions

### Iteration 3: Autonomous Documentation
- Auto-generated docs and changelogs
- Reflection reports

### Iteration 4: Agent Integration APIs
- Enhanced REST and gRPC interfaces
- Hot-plug tool capabilities

### Iteration 5: Task Orchestration
- Multi-step adaptive execution
- Task chaining and workflow management

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙋‍♂️ Support

- Create an [Issue](https://github.com/yourusername/aionmcp/issues) for bug reports or feature requests
- Start a [Discussion](https://github.com/yourusername/aionmcp/discussions) for questions or ideas

## 🌟 Acknowledgments

- [Model Context Protocol](https://modelcontextprotocol.io/) for the MCP specification
- [OpenAPI Initiative](https://www.openapis.org/) for OpenAPI specifications
- [GraphQL Foundation](https://graphql.org/) for GraphQL schemas
- [AsyncAPI Initiative](https://www.asyncapi.com/) for AsyncAPI specifications