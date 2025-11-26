# Automatic Claude Desktop Configuration

## Overview

AionMCP now provides a one-click VS Code extension command to automatically configure Claude Desktop to use AionMCP as an MCP server. This eliminates the need for manual JSON file editing and makes AionMCP truly model-independent.

## Features

✅ **Automatic Configuration**: Single extension command to set up Claude Desktop  
✅ **Binary Detection**: Automatically finds AionMCP server binary in common locations  
✅ **Config File Management**: Creates or updates `~/.config/Claude/claude_desktop_config.json`  
✅ **Error Handling**: Clear feedback on success or failure  
✅ **Model Independent**: Works with Claude Desktop and any MCP-compatible client  

## How It Works

### Extension Command: "Configure Claude Desktop for AionMCP"

1. **Open VS Code Command Palette**: `Cmd+Shift+P` (macOS) / `Ctrl+Shift+P` (Linux/Windows)
2. **Search for command**: Type "Configure Claude Desktop"
3. **Execute**: Press Enter to run the configuration

### What Happens Behind the Scenes

The command:
1. **Finds the binary**: Searches for AionMCP server in:
   - `~/go/bin/aionmcp-server` (Go workspace)
   - `/usr/local/bin/aionmcp-server` (System binary)
   - `/opt/homebrew/bin/aionmcp-server` (Homebrew, macOS)
   - `~/.local/bin/aionmcp-server` (User local binary)

2. **Creates config directory**: `~/.config/Claude/` (if needed)

3. **Updates configuration**: Modifies `claude_desktop_config.json` to add:
   ```json
   {
     "mcpServers": {
       "aionmcp": {
         "command": "/path/to/aionmcp-server",
         "args": [],
         "env": {}
       }
     }
   }
   ```

4. **Shows confirmation**: Displays success message with option to view config file

5. **Provides instructions**: Guides user to restart Claude Desktop

## Usage Flow

```
┌─────────────────────────────────┐
│  Install AionMCP Extension      │
└────────────┬────────────────────┘
             │
             ▼
┌─────────────────────────────────┐
│  Start AionMCP Server           │
│  (via extension button)          │
└────────────┬────────────────────┘
             │
             ▼
┌─────────────────────────────────┐
│  Configure Claude Desktop       │
│  (via extension command)        │
└────────────┬────────────────────┘
             │
             ▼
┌─────────────────────────────────┐
│  Restart Claude Desktop         │
│  (user action)                  │
└────────────┬────────────────────┘
             │
             ▼
┌─────────────────────────────────┐
│  Claude Desktop recognizes      │
│  AionMCP in tools list          │
└─────────────────────────────────┘
```

## Configuration File Location

- **macOS/Linux**: `~/.config/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

## Binary Paths Checked

The extension automatically searches for the AionMCP server binary in these locations (in order):

1. **Go Workspace**: `~/go/bin/aionmcp-server`
2. **System Binary**: `/usr/local/bin/aionmcp-server`
3. **Homebrew (macOS)**: `/opt/homebrew/bin/aionmcp-server`
4. **User Local**: `~/.local/bin/aionmcp-server`

If the binary is not found in any of these locations, the configuration will fail with a helpful error message.

## Building and Installing the Binary

### Option 1: Build from Source

```bash
# Clone the repository
git clone https://github.com/kiransth77/aionmcp.git
cd aionmcp

# Build the server
go build -o bin/aionmcp-server ./cmd/server/main.go

# Install to system path (optional)
sudo cp bin/aionmcp-server /usr/local/bin/
```

### Option 2: Add to Go Workspace

```bash
# Build to Go workspace
go build -o ~/go/bin/aionmcp-server ./cmd/server/main.go
```

### Option 3: Homebrew (Coming Soon)

```bash
brew install aionmcp
```

## Troubleshooting

### Error: "AionMCP binary not found"

**Solution**: Build the server binary and place it in one of the supported locations:

```bash
# Build the binary
cd /Users/kiran/Documents/GitHub/aionmcp
go build -o bin/aionmcp-server ./cmd/server/main.go

# Or copy to system location
sudo cp bin/aionmcp-server /usr/local/bin/
```

### Claude Desktop doesn't recognize AionMCP after configuration

**Solution**: Restart Claude Desktop completely:

1. Close Claude Desktop entirely
2. Wait 2 seconds
3. Reopen Claude Desktop
4. Check the tools list - AionMCP should now appear

### Configuration file created but Claude Desktop still doesn't work

**Steps to debug**:

1. Verify the config file exists:
   ```bash
   cat ~/.config/Claude/claude_desktop_config.json
   ```

2. Verify the binary path in the config is correct:
   ```bash
   which aionmcp-server
   ```

3. Test that the binary is executable:
   ```bash
   /path/to/aionmcp-server --version
   ```

4. Check Claude Desktop logs:
   ```bash
   tail -f ~/.claude-desktop/logs/*.log
   ```

## Model Independence

This configuration approach is model-independent because:

1. **MCP Standard**: Uses Model Context Protocol (MCP) standard, supported by multiple platforms
2. **Any LLM Client**: Works with Claude Desktop, and will work with future MCP-compatible clients
3. **Not Claude-Specific**: Configuration is generic MCP format, not Claude-specific
4. **Future Ready**: As more LLM platforms adopt MCP, this same configuration will work with them

### Supported Clients

- ✅ Claude Desktop (primary)
- 🔄 Other MCP-compatible clients (future support)

## Advanced Configuration

You can manually edit the configuration file to customize behavior:

```json
{
  "mcpServers": {
    "aionmcp": {
      "command": "/path/to/aionmcp-server",
      "args": ["--config", "/custom/config.yaml"],
      "env": {
        "AIONMCP_LOG_LEVEL": "debug",
        "AIONMCP_PORT": "8080"
      }
    }
  }
}
```

### Available Environment Variables

- `AIONMCP_LOG_LEVEL`: Log level (debug, info, warn, error)
- `AIONMCP_PORT`: HTTP server port (default: 8080)
- `AIONMCP_GRPC_PORT`: gRPC server port (default: 9090)
- `AIONMCP_DB_PATH`: BoltDB storage path (default: ./data/aionmcp.db)

## Related Documentation

- [Copilot Integration Guide](./COPILOT_INTEGRATION.md)
- [Architecture Overview](./architecture.md)
- [API Reference](./README.md#api-endpoints)

## Support

For issues or questions about Claude Desktop configuration:

1. Check the [Troubleshooting](#troubleshooting) section above
2. Review [Claude Desktop Documentation](https://claude.ai/docs)
3. Report issues on [GitHub](https://github.com/kiransth77/aionmcp/issues)
