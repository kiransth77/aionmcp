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
**Current Branch**: `copilot/sub-pr-2-again`

**Latest Commit**: [`4211070`](../../commit/42110703db96b228d6fadbaf2ccc633768e1c849)

**System Health**: 99/100 (Excellent)

**Active Tools**: 3

**Commits (7 days)**: 9

*Status updated automatically*
<!-- END AUTO-GENERATED STATUS -->

## ✨ Features

### Core Capabilities

- **Multi-Spec Import**: Automatically imports and converts API specifications
- **Dynamic Tool Registry**: Hot-reload tools without service restart
- **Self-Learning Engine**: Analyzes patterns and generates insights
- **Autonomous Documentation**: Auto-generates changelogs and reflections
- **Performance Monitoring**: Real-time execution metrics and optimization
- **Error Recovery**: Intelligent error handling and pattern detection

### API Support

- **OpenAPI 3.0+**: REST API specifications with full schema support
- **GraphQL**: Query and mutation support with type introspection
- **AsyncAPI**: Event-driven API specifications
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

- [`4211070`](../../commit/42110703db96b228d6fadbaf2ccc633768e1c849) Initial plan *(0h ago)*
- [`7724811`](../../commit/772481134d3b465bc56a0547288e8dfd5cc6d611) refactor: Extract magic numbers to constants and simplify code structure *(4h ago)*
- [`ed72de9`](../../commit/ed72de972116228d3cb6e1c3a258e2d61e087f62) Initial plan *(4h ago)*
- [`cd486cf`](../../commit/cd486cfe4f821fe30f18bfba9d3c5aacf23ace34) refactor: Extract duplicated logic and add missing constants *(5h ago)*
- [`e8a2a4b`](../../commit/e8a2a4bb652a2477bac85a073e25e16821385931) Initial plan *(5h ago)*

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

### Local Development

```bash
# Run tests
go test ./...

# Run with hot reload
go run cmd/server/main.go

# Build for production
go build -ldflags "-s -w" -o bin/aionmcp cmd/server/main.go
```
## 🤝 Contributing

### Development Process

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request
## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

*README last updated: 2025-11-02 11:10:32 UTC*

*This README is automatically updated with current project status and metrics.*
