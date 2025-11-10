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
**Current Branch**: `copilot/implement-tool-for-mobile`

**Latest Commit**: [`835449a`](../../commit/835449aceeb0a53910b0605a45d4268d5ebd1d75)

**System Health**: 99/100 (Excellent)

**Active Tools**: 3

**Commits (7 days)**: 7

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

- [`835449a`](../../commit/835449aceeb0a53910b0605a45d4268d5ebd1d75) Final: Add complete mobile platform section content to README *(0h ago)*
- [`f1ed5da`](../../commit/f1ed5da50ed49e42ba634fd973dc859f1566a27b) Update README generator to preserve mobile platform section *(0h ago)*
- [`f61a487`](../../commit/f61a4879d0fa2ca5857f51e2e68e77a9e035c3c3) Add mobile platform section to README with documentation links *(0h ago)*
- [`9ef32db`](../../commit/9ef32db05a211c58164728d045d91bc2098ebe40) Add comprehensive mobile platform support documentation and examples *(0h ago)*
- [`6955631`](../../commit/6955631cbf55d022fde4f957ca72029a208eb4be) Initial plan for mobile platform support *(0h ago)*

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
## 📱 Mobile Platform Support

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
## 🛠️ Development

### Key Findings

- **Self-Learning and Failure Recovery**: Analysis of reflection records in the learning engine (`docs/reflections/`) indicates recurring parameter validation errors in OpenAPI tool executions. Adaptive retry mechanisms and enhanced validation feedback have been implemented in the importer module to address these issues.
  
- **Hot-Reload Stability**: Hot-reload functionality demonstrates reliability for OpenAPI and GraphQL specifications. However, AsyncAPI event streams use a simple 500ms debounce mechanism in the watcher implementation (`pkg/importer/watcher.go`) to mitigate excessive reloads in high-frequency scenarios.

- **Documentation Automation**: The autodocs generators (`internal/autodocs/`) now correlate changelog and reflection outputs with tool confidence scores, facilitating efficient identification and resolution of unreliable tools.

- **Example Specifications**: Provided sample specifications (`examples/specs/petstore.yaml`, `examples/specs/blog.graphql`, `examples/specs/user-events.yaml`) support integration testing and developer onboarding.

### Future Enhancements

- Incorporate additional sample specifications to validate authentication workflows and large-scale schemas.
- Implement a health-check endpoint for the watcher subsystem to monitor reload backoff status.
- Integrate release automation for cross-platform binary artifacts upon version tagging (see `.github/workflows/release.yml`).
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

*README last updated: 2025-11-10 10:04:05 UTC*

*This README is automatically updated with current project status and metrics.*
