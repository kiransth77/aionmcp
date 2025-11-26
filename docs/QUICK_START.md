# AionMCP Quick Start Guide

## 60-Second Setup

### Step 1: Install the Extension (30 seconds)
1. Open VS Code
2. Go to Extensions (Cmd+Shift+X / Ctrl+Shift+X)
3. Search for "AionMCP"
4. Click "Install"

### Step 2: Build the Server (30 seconds)
```bash
cd /Users/kiran/Documents/GitHub/aionmcp
go build -o bin/aionmcp-server ./cmd/server/main.go
```

### Step 3: One-Click Setup (5 minutes total)
1. In VS Code, click the **AionMCP** icon in the sidebar (left activity bar)
2. Click the **play button** (▶) to start the server
3. Wait for "Server Status: Healthy" in the tree view
4. Click the **⚙️ settings gear icon** in the Server Status section
5. A dialog will appear saying "AionMCP configured successfully"
6. Close Claude Desktop completely
7. Reopen Claude Desktop
8. **Done!** AionMCP now appears in Claude's tools list

## What Just Happened?

The extension automatically:
- ✅ Found your AionMCP server binary
- ✅ Created `~/.config/Claude/claude_desktop_config.json`
- ✅ Added AionMCP as an MCP server
- ✅ Showed you the exact configuration file

## Using AionMCP with Claude

### In Claude Desktop:
1. Start a conversation with Claude
2. Look for **"AionMCP Tools"** in the tools selector
3. You'll see tools from any specs you've imported
4. Use them naturally in your conversation

### Examples:
```
You: "Using the Petstore API tools, list all available pets"
Claude: [Uses AionMCP tools to query the API]

You: "Can you query the weather API for London?"
Claude: [Uses imported weather spec tools via AionMCP]
```

## Importing Your First API Spec

### Option 1: Via Extension (Easiest)
1. In VS Code, go to AionMCP sidebar
2. Click the **Files icon** (in the Tools section)
3. Select spec type (OpenAPI, GraphQL, AsyncAPI)
4. Provide path/URL to your spec
5. Click OK
6. Tools appear in the tree view automatically

### Option 2: Via API
```bash
curl -X POST http://localhost:8080/api/v1/import-spec \
  -H "Content-Type: application/json" \
  -d '{
    "type": "openapi",
    "path": "./examples/specs/petstore.yaml",
    "name": "Petstore API"
  }'
```

### Example Specs Included
- **Petstore API** (OpenAPI): `./examples/specs/petstore.yaml`
- **Blog Schema** (GraphQL): `./examples/specs/blog.graphql`
- **User Events** (AsyncAPI): `./examples/specs/user-events.yaml`

Try these to see AionMCP in action:
```bash
# Using the extension UI is easier!
# Or via curl:
curl -X POST http://localhost:8080/api/v1/import-spec \
  -H "Content-Type: application/json" \
  -d '{
    "type": "openapi",
    "path": "./examples/specs/petstore.yaml",
    "name": "Petstore API"
  }'
```

## Architecture Overview

```
Claude Desktop (or any MCP client)
         │
         ▼
  ~/.config/Claude/claude_desktop_config.json
         │
         ▼
  AionMCP Server (localhost:8080/9090)
         │
    ┌────┴──────┐
    │            │
┌───▼──┐    ┌───▼──┐
│ Spec │    │Tool  │
│Parser│    │Call  │
│      │    │Handler
│ • OpenAPI │
│ • GraphQL │
│ • AsyncAPI
│           
└───────────────────────┘
         │
         ▼
    Your APIs
```

## Common Tasks

### Start the Server
```bash
# Via VS Code: Click the ▶ button in AionMCP sidebar
# Or manually:
cd /Users/kiran/Documents/GitHub/aionmcp
./bin/aionmcp-server
```

### Stop the Server
```bash
# Via VS Code: Click the ⏹ button in AionMCP sidebar
# Or in terminal: Ctrl+C
```

### Check Server Status
```bash
curl http://localhost:8080/api/v1/health | jq .
```

### View Server Logs
```bash
# Via VS Code: Click "View Server Logs" in AionMCP sidebar
# Or manually:
cd /Users/kiran/Documents/GitHub/aionmcp
tail -f bin/aionmcp-server.log
```

### Import an API Spec
```bash
# Easiest: Use VS Code extension UI
# Click the file icon in Tools section → follow prompts

# Or via curl:
curl -X POST http://localhost:8080/api/v1/import-spec \
  -H "Content-Type: application/json" \
  -d '{
    "type": "openapi",
    "path": "./path/to/your/api.yaml",
    "name": "My API"
  }'
```

### List All Imported Tools
```bash
curl http://localhost:8080/api/v1/tools | jq .
```

### Execute a Tool Directly
```bash
curl -X POST http://localhost:8080/api/v1/tools/{tool-id}/invoke \
  -H "Content-Type: application/json" \
  -d '{"param": "value"}'
```

## Troubleshooting

### Q: "Server is not running"
**A:** Click the ▶ button in the AionMCP sidebar to start it

### Q: "Configure Claude Desktop option not appearing"
**A:** 
1. Make sure the server is running (see above)
2. Restart VS Code
3. Try the configure button again

### Q: "Claude Desktop doesn't recognize AionMCP"
**A:**
1. Verify config was created: `cat ~/.config/Claude/claude_desktop_config.json`
2. Close Claude Desktop completely (Cmd+Q)
3. Wait 2 seconds
4. Reopen Claude Desktop
5. Check tools list

### Q: "Binary not found error"
**A:**
1. Rebuild: `go build -o bin/aionmcp-server ./cmd/server/main.go`
2. Or copy to system location: `sudo cp bin/aionmcp-server /usr/local/bin/`
3. Try configure again

### Q: "Can't import API spec"
**A:**
1. Make sure the file path is correct
2. For URLs, ensure they're publicly accessible
3. Verify the spec format (OpenAPI 3.0, GraphQL, or AsyncAPI 2.0+)
4. Check server logs for details

### Q: "Tools not appearing in Claude"
**A:**
1. Import the spec via VS Code extension
2. Wait a few seconds
3. Restart Claude Desktop completely
4. Refresh tools list in Claude

## Next Steps

1. **Import your first API**: See [Importing Your First API Spec](#importing-your-first-api-spec)
2. **Use with Claude**: Open Claude Desktop and look for AionMCP tools
3. **Learn more**: See [docs/README.md](./README.md) for full documentation
4. **Explore examples**: Check `./examples/specs/` for test specifications

## Resources

- **Full Documentation**: [docs/README.md](./README.md)
- **Architecture Details**: [docs/architecture.md](./architecture.md)
- **Claude Desktop Setup**: [docs/CLAUDE_DESKTOP_AUTO_CONFIG.md](./CLAUDE_DESKTOP_AUTO_CONFIG.md)
- **Copilot Integration**: [docs/COPILOT_INTEGRATION.md](./COPILOT_INTEGRATION.md)
- **API Reference**: [docs/README.md#api-endpoints](./README.md#api-endpoints)

## Support

- **Issues**: Report on GitHub
- **Questions**: Check the troubleshooting section
- **Contributing**: Fork and submit PRs

---

**Ready to start?** 👉 [Install the Extension](#60-second-setup)

**Want more details?** 👉 [Full Documentation](./README.md)

**Questions?** 👉 [Troubleshooting](#troubleshooting)
