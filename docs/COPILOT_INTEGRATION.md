# GitHub Copilot Integration with AionMCP

This guide explains how to integrate GitHub Copilot with AionMCP to make your API tools available to Copilot.

## Understanding Agents vs Tools

- **Tools**: API operations imported from OpenAPI, GraphQL, or AsyncAPI specs
- **Agents**: AI systems (like Copilot) that can execute tools and access APIs

## How It Works

AionMCP acts as an MCP (Model Context Protocol) server that provides:
1. **Tools**: Functions and APIs imported from your specifications
2. **Agent Registration**: A system for agents (like Copilot) to connect and register sessions

## Setting Up Copilot Integration

### Step 1: Start the AionMCP Server

1. Open the AionMCP extension in VS Code
2. Click **"Start AionMCP Server"** in the Server Status view
3. Verify the server is running (check status bar shows ✓)

### Step 2: Import API Specifications

To make APIs available to Copilot:

1. In the **Tools** view, click **"Import API Specification"**
2. Select the specification type:
   - `openapi` - OpenAPI/Swagger specs
   - `graphql` - GraphQL schemas
   - `asyncapi` - AsyncAPI event specs
3. Provide the path to the spec file (local or URL)
4. Tools from the spec will be imported and available

### Step 3: Verify Tools Are Available

1. Check the **Tools** view to see imported tools
2. Check the **Server Status** view for "Tool Count"
3. Verify via API:
   ```bash
   curl http://localhost:8080/api/v1/tools
   ```

### Step 4: Configure Copilot to Use AionMCP

To make tools available in Claude Desktop or GitHub Copilot, the MCP server needs to be configured in your client.

#### For Claude Desktop:

1. **Edit config file**: `~/.config/Claude/claude_desktop_config.json`
   
2. **Add AionMCP configuration**:
   ```json
   {
     "mcpServers": {
       "aionmcp": {
         "command": "/path/to/aionmcp",
         "args": [],
         "env": {
           "MCP_PORT": "9090"
         }
       }
     }
   }
   ```

3. **Restart Claude Desktop**

4. **Verify**: Go to **Settings → Configure Tools → Model Packages** and look for AionMCP tools

## Legacy: Register Copilot Agent (Experimental)

The extension includes an experimental "Register GitHub Copilot Agent" button that programmatically registers a Copilot agent session:

1. Click **"Register GitHub Copilot Agent"** in the Connected Agents view
2. You should see confirmation in the logs
3. The agent will appear in the Connected Agents list with status ACTIVE

**Note**: This creates a session-based registration that expires after 1 hour. For production use, configure Claude Desktop as described above.

## Troubleshooting

### Copilot Agent Not Showing

**Problem**: Registered agent doesn't appear in Connected Agents view

**Solutions**:
1. Verify server is running (check status bar)
2. Check server logs for registration errors
3. Try re-registering the agent
4. Check network connectivity on port 8080

### Tools Not Available to Copilot

**Problem**: Imported tools don't execute through Copilot

**Solutions**:
1. Verify tools are in the Tools view
2. Check agent session is active (not expired)
3. Verify Copilot has necessary permissions
4. Check server logs for execution errors

### Session Timeout Errors

**Problem**: "Agent session expired" error

**Solutions**:
1. Re-register the agent
2. Increase timeout in server configuration
3. Check network connectivity

## Architecture

```
┌─────────────────┐
│ GitHub Copilot  │
└────────┬────────┘
         │
    gRPC/HTTP
         │
┌────────▼────────┐
│  AionMCP Server │
│  ┌────────────┐ │
│  │ Agent Mgmt │ │
│  ├────────────┤ │
│  │ Tool Reg   │ │
│  ├────────────┤ │
│  │ Execution  │ │
│  └────────────┘ │
└─────────────────┘
         │
    ┌────▼─────┐
    │API Specs  │
    │(OpenAPI,  │
    │GraphQL,   │
    │AsyncAPI)  │
    └───────────┘
```

## Next Steps

1. **Import multiple specifications** to build a comprehensive API catalog
2. **Configure Claude Desktop** to access AionMCP tools
3. **Test tool execution** through Copilot
4. **Monitor usage** via server stats and logs
5. **Customize capabilities** in agent metadata

## API Reference

### Register Mock Agent

```
POST /api/v1/agents/register-mock
```

**Response**:
```json
{
  "session_id": "string (UUID)",
  "agent_id": "github-copilot",
  "agent_name": "GitHub Copilot",
  "expires_at": "number (Unix timestamp)",
  "message": "string"
}
```

### List Agents

```
GET /api/v1/agents
```

**Response**:
```json
{
  "agents": [
    {
      "session_id": "string",
      "agent_id": "string",
      "agent_name": "string",
      "version": "string",
      "status": "AGENT_STATUS_ACTIVE",
      "created_at": "string (RFC3339)"
    }
  ]
}
```

### Import Specification

```
POST /api/v1/import-spec
```

**Request**:
```json
{
  "spec_type": "openapi|graphql|asyncapi",
  "path": "string (file path or URL)",
  "name": "string (optional)"
}
```

**Response**:
```json
{
  "source": {
    "path": "string",
    "type": "string"
  },
  "tools": [
    {
      "name": "string",
      "description": "string",
      "version": "string"
    }
  ],
  "warnings": ["string"],
  "errors": ["string"],
  "duration": "number (ms)"
}
```

## Support

For issues or questions:
1. Check the troubleshooting section above
2. Review server logs in the extension
3. Open an issue on GitHub
