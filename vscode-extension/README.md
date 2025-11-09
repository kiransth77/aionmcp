# AionMCP VS Code Extension

The AionMCP extension brings the power of the Autonomous Model Context Protocol server directly into your VS Code environment.

## Features

- **Server Management**: Start, stop, and restart the AionMCP server with one click
- **Tool Explorer**: Browse and discover available API tools in a tree view
- **Tool Execution**: Execute tools with interactive parameter forms
- **Agent Monitoring**: View connected agents and their status
- **Real-time Logs**: Monitor server activity and execution results
- **Dashboard**: Comprehensive overview of server metrics and statistics
- **API Spec Import**: Drag-and-drop support for OpenAPI, GraphQL, and AsyncAPI files

## Quick Start

1. Install the extension
2. Open the AionMCP activity bar panel
3. Click "Start Server" in the Server Status view
4. Explore available tools in the Tools view
5. Execute tools by right-clicking and selecting "Execute Tool"

## Configuration

- `aionmcp.serverPath`: Custom path to AionMCP binary (leave empty to use bundled version)
- `aionmcp.serverPort`: HTTP server port (default: 8080)
- `aionmcp.grpcPort`: gRPC server port (default: 50051)
- `aionmcp.autoStart`: Automatically start server when VS Code opens
- `aionmcp.logLevel`: Server log level (debug, info, warn, error)
- `aionmcp.specDirectories`: Directories to watch for API specifications

## Requirements

- VS Code 1.85.0 or higher
- No additional dependencies required (binary is bundled)

## Extension Commands

- `AionMCP: Start Server` - Start the AionMCP server
- `AionMCP: Stop Server` - Stop the AionMCP server
- `AionMCP: Restart Server` - Restart the AionMCP server
- `AionMCP: Refresh Tools` - Refresh the tools list
- `AionMCP: Import API Specification` - Import API spec files
- `AionMCP: Show Dashboard` - Open the dashboard webview
- `AionMCP: View Logs` - Show server logs

## Known Issues

- Server may take a few seconds to start on first launch
- Tool parameter validation is basic (full JSON schema validation coming soon)

## Release Notes

### 0.1.0

Initial release with core functionality:
- Server lifecycle management
- Tool discovery and execution
- Agent monitoring
- Real-time logging
- Dashboard interface