# Model-Independent Tool Integration

## Architecture

AionMCP provides tools to **any agent or AI model** through a standard HTTP REST API - not through model-specific protocols.

```
GitHub Copilot ─┐
Claude Desktop  ├─→ HTTP/REST API (localhost:8080)
Any LLM Agent   ├─→ AionMCP Server
Custom Apps     ┘
                │
                ▼
        Tool Registry (OpenAPI, GraphQL, AsyncAPI)
                │
                ▼
            Your APIs
```

## Why This Approach?

✅ **Model Independent**: No dependency on specific LLM platforms  
✅ **Agent Agnostic**: Works with GitHub Copilot, Claude, custom agents, etc.  
✅ **Future Proof**: New LLMs can integrate without code changes  
✅ **Simple Integration**: Just HTTP API calls  
✅ **Flexible Deployment**: Can be hosted anywhere  

## REST API Interface

### List All Available Tools
```bash
GET /api/v1/tools

Response:
{
  "tools": [
    {
      "id": "openapi.petstore.listPets",
      "name": "List Pets",
      "description": "Get all available pets",
      "category": "openapi",
      "parameters": {
        "limit": { "type": "integer", "required": false }
      }
    }
  ]
}
```

### Get Tool Details
```bash
GET /api/v1/tools/{tool-id}

Response:
{
  "id": "openapi.petstore.listPets",
  "name": "List Pets",
  "description": "Get all available pets",
  "category": "openapi",
  "parameters": {
    "limit": {
      "type": "integer",
      "required": false,
      "description": "Max number of pets to return"
    }
  },
  "returns": {
    "type": "array",
    "items": {
      "type": "object",
      "properties": {
        "id": { "type": "integer" },
        "name": { "type": "string" }
      }
    }
  }
}
```

### Execute a Tool
```bash
POST /api/v1/tools/{tool-id}/invoke

Request Body:
{
  "limit": 10
}

Response:
{
  "success": true,
  "result": [
    { "id": 1, "name": "Fluffy" },
    { "id": 2, "name": "Buddy" }
  ],
  "execution_time_ms": 124
}
```

## Integration Examples

### GitHub Copilot Integration

GitHub Copilot can be configured to use AionMCP tools by calling the REST API in a custom extension or through direct integration.

### Example: Custom Copilot Extension

```typescript
// This is how a Copilot agent would use AionMCP
async function getAionMCPTools() {
  const response = await fetch('http://localhost:8080/api/v1/tools');
  const { tools } = await response.json();
  return tools;
}

async function executeTool(toolId, params) {
  const response = await fetch(`http://localhost:8080/api/v1/tools/${toolId}/invoke`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(params)
  });
  return response.json();
}
```

### Example: Custom Python Agent

```python
import requests

def get_tools():
    response = requests.get('http://localhost:8080/api/v1/tools')
    return response.json()['tools']

def execute_tool(tool_id, params):
    response = requests.post(
        f'http://localhost:8080/api/v1/tools/{tool_id}/invoke',
        json=params
    )
    return response.json()

# Use the tools
tools = get_tools()
for tool in tools:
    print(f"Tool: {tool['name']}")
    
# Execute a tool
result = execute_tool('openapi.petstore.listPets', {'limit': 5})
print(result['result'])
```

### Example: cURL (Direct API Calls)

```bash
# List tools
curl http://localhost:8080/api/v1/tools | jq .

# Execute a tool
curl -X POST http://localhost:8080/api/v1/tools/openapi.petstore.listPets/invoke \
  -H "Content-Type: application/json" \
  -d '{"limit": 5}'
```

## Server Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/v1/health` | GET | Health check |
| `/api/v1/tools` | GET | List all tools |
| `/api/v1/tools/{id}` | GET | Get tool details |
| `/api/v1/tools/{id}/invoke` | POST | Execute tool |
| `/api/v1/import-spec` | POST | Import API specification |
| `/api/v1/server-stats` | GET | Server statistics |

## Configuration for Different Agents

### GitHub Copilot

GitHub Copilot can be configured via:
1. **Direct API Integration**: In a custom Copilot plugin/extension
2. **Environment Variables**: Point Copilot client to AionMCP server
3. **Configuration File**: Custom Copilot config with tool registry

### Other Agents

Any agent that supports:
- HTTP/REST API calls ✅
- JSON request/response format ✅
- Dynamic tool discovery ✅

Can integrate with AionMCP.

## Running AionMCP as a Tool Server

### Start the Server
```bash
# Via VS Code extension
# Click the ▶ Play button in AionMCP sidebar

# Or manually
cd /Users/kiran/Documents/GitHub/aionmcp
./bin/aionmcp-server
```

### Verify It's Running
```bash
curl http://localhost:8080/api/v1/health
# Response: {"status":"healthy","version":"0.1.0",...}
```

### Import Tools
```bash
# Via VS Code extension: Click file icon in Tools section
# Or via API:
curl -X POST http://localhost:8080/api/v1/import-spec \
  -H "Content-Type: application/json" \
  -d '{
    "type": "openapi",
    "path": "./examples/specs/petstore.yaml",
    "name": "Petstore API"
  }'
```

### Use with Any Agent

Any agent/application can now:
```
GET /api/v1/tools              → Discover available tools
GET /api/v1/tools/{id}         → Get tool schema
POST /api/v1/tools/{id}/invoke → Execute tool
```

## Security Considerations

### Current (Development)
- ✅ Running on localhost (127.0.0.1)
- ✅ No authentication required (local only)

### Production (Future)
- 🔄 Add API key authentication
- 🔄 Add rate limiting
- 🔄 Add CORS policy configuration
- 🔄 Add request validation and sanitization
- 🔄 Add execution timeout enforcement
- 🔄 Add audit logging for tool executions

## Deployment Options

### Option 1: Local Development
```bash
# Run on localhost
./bin/aionmcp-server
# Access: http://localhost:8080
```

### Option 2: Docker Container
```bash
# Run in container for isolation
docker run -p 8080:8080 aionmcp:latest
```

### Option 3: Cloud Deployment
```bash
# Deploy to cloud with authentication/security
# Examples: AWS Lambda, Google Cloud Run, Azure Functions
```

### Option 4: Embedded
```go
// Embed AionMCP in your Go application
import "github.com/aionmcp/aionmcp/internal/core"

server := core.NewServer(logger)
server.Start()
defer server.Stop()
```

## Example: Using with GitHub Copilot

### Scenario: Chat with Copilot using AionMCP Tools

```
User: "Can you list the available pets in the Petstore API?"

Copilot Agent Flow:
1. Call GET /api/v1/tools
   → Gets: [openapi.petstore.listPets, openapi.petstore.getPet, ...]
2. Recognize "list pets" matches "openapi.petstore.listPets"
3. Get tool schema: GET /api/v1/tools/openapi.petstore.listPets
   → Gets parameter requirements
4. Call tool: POST /api/v1/tools/openapi.petstore.listPets/invoke
   → Body: {"limit": 10}
   → Response: [{"id": 1, "name": "Fluffy"}, ...]
5. Respond to user with results

Copilot: "There are 2 pets available:
- ID 1: Fluffy
- ID 2: Buddy"
```

## Advantages Over MCP

| Aspect | MCP Protocol | HTTP REST (AionMCP) |
|--------|-------------|-------------------|
| Agent Independence | Protocol-specific | Any HTTP client |
| Integration | Requires MCP support | Any language/platform |
| Discoverability | Protocol-defined | Simple GET request |
| Execution | gRPC/Protocol Buffers | Standard JSON |
| Future-Proof | Tied to MCP evolution | Industry-standard HTTP |
| Hosting | Specific process | Any web server |
| Scalability | Limited | Can use reverse proxy, load balancer |

## Next Steps

1. **Ensure server is running** with imported tools
2. **Test the REST API** to verify tools are accessible
3. **Integrate with your agent** (GitHub Copilot, custom agent, etc.)
4. **Deploy** wherever your agent needs it

## Support

For help with:
- **Server**: See [docs/README.md](./README.md)
- **API**: See [docs/QUICK_START.md](./QUICK_START.md)
- **Agents**: Check the integration examples above
- **Issues**: Report on GitHub

---

**Key Principle**: AionMCP is a **Tool Server**, not a model-specific integration. Any agent can use it via standard HTTP REST API.
