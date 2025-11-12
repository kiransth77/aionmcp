# AionMCP - Autonomous Go MCP Server

<!-- AUTO-GENERATED BADGES -->
![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![Success Rate](https://img.shields.io/badge/success_rate-97%25-brightgreen)
![Avg Latency](https://img.shields.io/badge/avg_latency-250ms-green)
![Go Version](https://img.shields.io/badge/go-1.21+-blue)
![License](https://img.shields.io/badge/license-MIT-blue)
[![Sponsor](https://img.shields.io/badge/Sponsor-%E2%9D%A4-red)](https://github.com/sponsors/kiransth77)
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

**Latest Commit**: [`63682e2`](../../commit/63682e2af452c91e8d8b33014180ee8082ae1d72)

**System Health**: 99/100 (Excellent)

**Active Tools**: 3

**Commits (7 days)**: 8

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

- [`63682e2`](../../commit/63682e2af452c91e8d8b33014180ee8082ae1d72) Add downloadable demo applications for Android and iOS *(2d ago)*
- [`835449a`](../../commit/835449aceeb0a53910b0605a45d4268d5ebd1d75) Final: Add complete mobile platform section content to README *(2d ago)*
- [`f1ed5da`](../../commit/f1ed5da50ed49e42ba634fd973dc859f1566a27b) Update README generator to preserve mobile platform section *(2d ago)*
- [`f61a487`](../../commit/f61a4879d0fa2ca5857f51e2e68e77a9e035c3c3) Add mobile platform section to README with documentation links *(2d ago)*
- [`9ef32db`](../../commit/9ef32db05a211c58164728d045d91bc2098ebe40) Add comprehensive mobile platform support documentation and examples *(2d ago)*

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

AionMCP provides full support for Android and iOS mobile applications through REST API and gRPC interfaces.

### 🎉 Demo Apps Available!

Download complete, ready-to-use mobile applications:

- **Android**: [Download APK](https://github.com/kiransth77/aionmcp/releases) | [Source Code](examples/mobile/android-app/)
- **iOS**: TestFlight Beta (Coming Soon) | [Source Code](examples/mobile/ios-app/)

See [Demo Apps Guide](examples/mobile/README.md) for installation and usage.

### Platform Support

- **Android**: Kotlin/Java integration with Retrofit and gRPC
- **iOS**: Swift integration with Alamofire and gRPC-Swift
- **Cross-Platform**: REST API compatible with React Native, Flutter, and other frameworks

### Quick Start

**Try the Demo Apps:**
1. Download for [Android](https://github.com/kiransth77/aionmcp/releases) or iOS (Coming Soon)
2. Install and open the app
3. Configure your AionMCP server URL in Settings
4. Start exploring tools!

**Build Your Own:**

*Android (Kotlin)*:
```kotlin
val client = AionMCPClient("https://your-server.com")
val tools = client.api.listTools()
```

*iOS (Swift)*:
```swift
let client = AionMCPClient(baseURL: "https://your-server.com")
let tools = try await client.listTools()
```

### Documentation

- 📖 [Complete Mobile Integration Guide](docs/mobile_integration.md)
- 📱 [Demo Applications Guide](examples/mobile/README.md)
- 🤖 [Android Code Examples](examples/mobile/android/)
- 🍎 [iOS Code Examples](examples/mobile/ios/)
- 🚀 [Mobile Deployment Guide](docs/mobile_deployment.md)

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

*README last updated: 2025-11-12 11:41:30 UTC*

*This README is automatically updated with current project status and metrics.*
