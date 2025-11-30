# AionMCP - Autonomous Go MCP Server

<!-- AUTO-GENERATED BADGES -->
![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![Success Rate](https://img.shields.io/badge/success_rate-97%25-brightgreen)
![Avg Latency](https://img.shields.io/badge/avg_latency-250ms-green)
![Go Version](https://img.shields.io/badge/go-1.21+-blue)
![License](https://img.shields.io/badge/license-MIT-blue)
<!-- END AUTO-GENERATED BADGES -->

AionMCP is an autonomous Go-based Model Context Protocol (MCP) server that dynamically imports OpenAPI, GraphQL, and AsyncAPI specifications and exposes them as tools to agents. It features self-learning capabilities, context-awareness, and autonomous documentation using Clean/Hexagonal architecture.

## 🌟 Key Differentiators

- **Multi-Protocol Support**: OpenAPI, GraphQL, and AsyncAPI specifications
- **Autonomous Learning**: Self-improving system that learns from execution patterns
- **Dynamic Runtime**: Hot-reloadable tools without service restart
- **Clean Architecture**: Maintainable, testable, and extensible design
- **Auto-Documentation**: Self-updating documentation and insights

## 📊 Project Status

<!-- AUTO-GENERATED STATUS -->
**Current Branch**: `main`

**Latest Commit**: [`f4f41db`](../../commit/f4f41db9fffc0d43c67153b80d68480905f7c3d3)

**System Health**: 99/100 (Excellent)

**Active Tools**: 3

**Commits (7 days)**: 17

*Status updated automatically*
<!-- END AUTO-GENERATED STATUS -->

## ✨ Features

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
## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/kiransth77/aionmcp.git
cd aionmcp

# Build the server
go build -o bin/aionmcp cmd/server/main.go

# Run with default configuration
./bin/aionmcp
```

The server will start on `http://localhost:8080` with learning enabled.

## 🏗️ Architecture

AionMCP follows Clean/Hexagonal Architecture principles:

```
┌─────────────────────────────────────────────────────────┐
│                    Adapters Layer                      │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐   │
│  │   HTTP      │  │    gRPC     │  │   Plugin    │   │
│  │  Interface  │  │  Interface  │  │  Interface  │   │
│  └─────────────┘  └─────────────┘  └─────────────┘   │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│                     Core Layer                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐   │
│  │    Tool     │  │  Learning   │  │    Auto     │   │
│  │  Registry   │  │   Engine    │  │    Docs     │   │
│  └─────────────┘  └─────────────┘  └─────────────┘   │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│                Infrastructure Layer                    │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐   │
│  │   Storage   │  │   Metrics   │  │   Config    │   │
│  │  (BoltDB)   │  │(Prometheus) │  │   (Viper)   │   │
│  └─────────────┘  └─────────────┘  └─────────────┘   │
└─────────────────────────────────────────────────────────┘
```

## 📈 Recent Activity

<!-- AUTO-GENERATED ACTIVITY -->
### Recent Commits

- [`f4f41db9`](../../commit/f4f41db9fffc0d43c67153b80d68480905f7c3d3) fix: clean up TOON package formatting and remove duplicate package declaration *(0h ago)*
- [`8b97b173`](../../commit/8b97b173d7a8b7b5521f03e6ef9f901b28603541) fix: update go-ci.yml to use Go 1.25 to match go.mod requirement *(0h ago)*
- [`ebeeee28`](../../commit/ebeeee2896733d1c8ff3da684e052b451f195bab) fix: remove redundant newline in fmt.Println to pass linting *(0h ago)*
- [`a91eb5f4`](../../commit/a91eb5f414c56f9e0d5b92a5f189f94e374bc452) feat: add TOON (Token-Oriented Object Notation) for LLM context optimization *(0h ago)*
- [`bd39bb10`](../../commit/bd39bb1023fe6886d712a6b2900b5537f108940a) fix: update Dockerfile Go version to 1.25 *(0h ago)*

### Active Insights

📊 Total insights: 2

*Activity updated automatically*
<!-- END AUTO-GENERATED ACTIVITY -->

## ⚡ Performance Statistics

<!-- AUTO-GENERATED PERFORMANCE -->
| Metric | Value | Status |
|--------|-------|--------|
| Success Rate | 97.0% | 🟢 Excellent |
| Avg Latency | 250.0ms | 🟡 Good |
| Total Executions | 42 | 📊 Tracking |
| Active Tools | 3 | 🔧 Running |

*Statistics updated in real-time*
<!-- END AUTO-GENERATED PERFORMANCE -->

## 📦 Installation

### Prerequisites

- Go 1.21 or higher
- Git

### From Source

```bash
git clone https://github.com/kiransth77/aionmcp.git
cd aionmcp
go mod download
go build -o bin/aionmcp cmd/server/main.go
```
## 📚 Usage

### Basic Usage

```bash
# Start the server
./bin/aionmcp

# With custom configuration
./bin/aionmcp --config config.yaml

# Enable debug logging
AIONMCP_LOG_LEVEL=debug ./bin/aionmcp
```

### API Endpoints

- `GET /api/v1/tools` - List available tools
- `POST /api/v1/tools/{tool}/execute` - Execute a tool
- `GET /api/v1/learning/stats` - Learning statistics
- `GET /api/v1/learning/insights` - System insights
## 🛠️ Development

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
## 🤝 Contributing

### Active Insights

📊 Total insights: 2

*Activity updated automatically*
<!-- END AUTO-GENERATED ACTIVITY -->
## 📄 License

## 📄 License
---

*README last updated: 2025-11-27 22:36:50 AEST*

*This README is automatically updated with current project status and metrics.*
---

*README last updated: 2025-11-27 22:37:28 AEST*

*This README is automatically updated with current project status and metrics.*
---

*README last updated: 2025-11-27 22:47:18 AEST*

*This README is automatically updated with current project status and metrics.*
