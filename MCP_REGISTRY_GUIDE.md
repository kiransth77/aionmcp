# Publishing AionMCP to the MCP Registry

## Overview

To publish AionMCP to the official Model Context Protocol (MCP) registry on GitHub, follow these steps.

## Prerequisites

1. **GitHub Account**: With push access to your repository
2. **Go 1.21+**: Installed and working
3. **Git**: For version control
4. **Authenticated with GitHub**: For pushing releases

## Step 1: Prepare Your Repository

### 1.1 Ensure MCP Server Implementation
Your server should expose tools via the MCP protocol. Currently, AionMCP exposes tools via:
- HTTP endpoints (`/api/v1/tools`)
- gRPC APIs (AgentService)

For proper MCP registry listing, ensure your server can be called by Claude and other MCP clients.

### 1.2 Update README with MCP Information
Add to your README.md:

```markdown
## MCP Server Setup

AionMCP is a Model Context Protocol server that exposes API specifications as tools.

### Installation

```bash
go install github.com/aionmcp/aionmcp/cmd/server@latest
```

### Configuration

Add to Claude Desktop's `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "aionmcp": {
      "command": "aionmcp",
      "args": ["--port", "8080"],
      "env": {
        "AIONMCP_LEARNING_ENABLED": "true"
      }
    }
  }
}
```

### Available Tools

- **Echo**: Echo back input for testing
- **Status**: Get server status information
- **Dynamic Tools**: Any tools imported from OpenAPI/GraphQL/AsyncAPI specs

See `/api/v1/tools` endpoint for full list of available tools.
```

## Step 2: Create a Proper Release

### 2.1 Create a git tag

```bash
cd /Users/kiran/Documents/GitHub/aionmcp
git tag -a v0.1.0 -m "Initial release: Multi-spec API importer with self-learning"
git push origin v0.1.0
```

### 2.2 Build binaries for multiple platforms

```bash
# macOS
GOOS=darwin GOARCH=amd64 go build -o dist/aionmcp-darwin-amd64 ./cmd/server
GOOS=darwin GOARCH=arm64 go build -o dist/aionmcp-darwin-arm64 ./cmd/server

# Linux
GOOS=linux GOARCH=amd64 go build -o dist/aionmcp-linux-amd64 ./cmd/server
GOOS=linux GOARCH=arm64 go build -o dist/aionmcp-linux-arm64 ./cmd/server

# Windows
GOOS=windows GOARCH=amd64 go build -o dist/aionmcp-windows-amd64.exe ./cmd/server
```

### 2.3 Create GitHub Release

1. Go to https://github.com/aionmcp/aionmcp/releases
2. Click "Create a new release"
3. Tag version: `v0.1.0`
4. Release title: "AionMCP v0.1.0 - Initial Release"
5. Description:
```markdown
## Features

- **Multi-Protocol Support**: OpenAPI, GraphQL, and AsyncAPI
- **Autonomous Learning**: Self-improving error recovery
- **Dynamic Tools**: Hot-reloadable without restart
- **Auto-Documentation**: Self-generating docs and insights

## Breaking Changes

None - Initial release.

## Installation

```bash
go install github.com/aionmcp/aionmcp/cmd/server@v0.1.0
```

See MCP_REGISTRY_GUIDE.md for setup instructions.
```

6. Upload binaries from `dist/` folder
7. Publish release

## Step 3: Submit to MCP Registry

### 3.1 Fork the MCP Registry

1. Go to https://github.com/modelcontextprotocol/servers
2. Click "Fork"
3. Clone your fork:
```bash
git clone https://github.com/YOUR_USERNAME/servers.git
cd servers
```

### 3.2 Create Server Entry

Create `servers/src/aionmcp/index.json`:

```json
{
  "name": "AionMCP",
  "description": "Autonomous Go MCP server that dynamically imports OpenAPI, GraphQL, and AsyncAPI specifications",
  "author": "Your Name",
  "homepage": "https://github.com/aionmcp/aionmcp",
  "repository": {
    "type": "git",
    "url": "https://github.com/aionmcp/aionmcp"
  },
  "language": "go",
  "mcpVersion": "1.0",
  "tools": {
    "enabled": true,
    "count": "3+"
  },
  "requirements": {
    "go": "1.21+"
  },
  "installation": {
    "type": "github",
    "owner": "aionmcp",
    "repo": "aionmcp",
    "bin": "aionmcp"
  },
  "documentation": "https://github.com/aionmcp/aionmcp/blob/main/README.md",
  "features": [
    "Multi-spec API import (OpenAPI, GraphQL, AsyncAPI)",
    "Autonomous learning and error recovery",
    "Dynamic tool registration",
    "Self-generating documentation",
    "Real-time metrics and monitoring"
  ]
}
```

### 3.3 Add Configuration Examples

Create `servers/src/aionmcp/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "aionmcp": {
      "command": "go",
      "args": ["run", "github.com/aionmcp/aionmcp/cmd/server@latest"],
      "env": {
        "AIONMCP_PORT": "8080",
        "AIONMCP_LEARNING_ENABLED": "true"
      }
    }
  }
}
```

Or for installed binary:

```json
{
  "mcpServers": {
    "aionmcp": {
      "command": "aionmcp",
      "env": {
        "AIONMCP_PORT": "8080",
        "AIONMCP_LEARNING_ENABLED": "true"
      }
    }
  }
}
```

### 3.4 Create Pull Request

```bash
git add servers/src/aionmcp/
git commit -m "Add AionMCP to registry"
git push origin main
```

Then:
1. Go to https://github.com/modelcontextprotocol/servers
2. Click "Compare & pull request"
3. Add description:
```markdown
## AionMCP Server Listing

Adds AionMCP - an autonomous Go MCP server for multi-spec API imports.

### Description
AionMCP dynamically imports OpenAPI, GraphQL, and AsyncAPI specifications and exposes them as tools. It features autonomous learning, hot-reloadable tools, and self-generating documentation.

### Key Features
- Multi-protocol support (OpenAPI, GraphQL, AsyncAPI)
- Autonomous learning engine
- Dynamic tool registration
- Auto-generated documentation

### Resources
- **GitHub**: https://github.com/aionmcp/aionmcp
- **Documentation**: [Link to docs]
- **Status**: Production-ready MVP

Addresses #[issue-number]
```

4. Submit PR and address any feedback from maintainers

## Step 4: Optimize for MCP Clients

### 4.1 Implement MCP Tool Discovery

Ensure `/api/v1/tools` endpoint returns proper MCP-compatible tool descriptions:

```json
{
  "tools": [
    {
      "name": "echo",
      "description": "Echo back the input message for testing purposes",
      "inputSchema": {
        "type": "object",
        "properties": {
          "message": {
            "type": "string",
            "description": "Message to echo"
          }
        },
        "required": ["message"]
      }
    }
  ]
}
```

### 4.2 Support Installation Methods

- Binary installation via `go install`
- Docker image (optional but recommended)
- Direct source compilation

## Step 5: Maintenance

### Version Bumping

When making updates:
1. Update version in `cmd/server/main.go`
2. Update `go.mod` if needed
3. Create git tag: `git tag -a v0.2.0 -m "Release message"`
4. Create GitHub release with binaries
5. Update MCP registry entry if needed (submit PR)

### Monitoring

- Watch GitHub issues for user feedback
- Monitor MCP registry for updates/requirements
- Keep dependencies updated

## Troubleshooting

### Server not appearing in Claude Desktop

1. Verify the registry entry format is correct
2. Ensure binaries are working: `aionmcp --version`
3. Check configuration in `claude_desktop_config.json`
4. Restart Claude Desktop
5. Check logs: `~/.local/state/Claude/debug.log` (Linux) or similar

### Tools not loading

1. Verify `/api/v1/tools` returns valid JSON
2. Check tool schema matches MCP spec
3. Verify server is running on configured port
4. Check server logs for errors

## Additional Resources

- **MCP Specification**: https://spec.modelcontextprotocol.io/
- **MCP Registry**: https://github.com/modelcontextprotocol/servers
- **Claude Desktop Config**: https://claude.ai/mcp
- **Go Project Template**: https://golang.org/doc/code

## Next Steps

1. [ ] Create release v0.1.0 with git tag
2. [ ] Build multi-platform binaries
3. [ ] Create GitHub Release
4. [ ] Fork MCP servers registry
5. [ ] Create AionMCP entry in registry
6. [ ] Submit PR to registry
7. [ ] Monitor PR for feedback
8. [ ] Merge and celebrate! 🎉
