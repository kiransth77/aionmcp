# Using AionMCP with GitHub Copilot (Model-Independent)

## Overview

This guide shows how to integrate AionMCP tools with **GitHub Copilot** in a model-independent way using simple HTTP REST API calls.

## Architecture

```
GitHub Copilot (or any agent)
         │
         ▼ HTTP REST API
    AionMCP Server (localhost:8080)
         │
    ┌────┴──────┐
    │            │
┌───▼──┐    ┌───▼──┐
│Tool  │    │Tool  │
│Registry    │Executor
│            │
│ • Tools    │ • Execution
│ • Schemas  │ • Results
└────────────────┴──────────────┘
         │
         ▼
    Your APIs (OpenAPI, GraphQL, AsyncAPI)
```

## Quick Start: 3 Steps

### Step 1: Start AionMCP Server
```bash
cd /Users/kiran/Documents/GitHub/aionmcp
./bin/aionmcp-server &
```

Or use the VS Code extension:
- Open AionMCP in the sidebar
- Click the ▶️ play button

### Step 2: Verify Server is Running
```bash
curl http://localhost:8080/api/v1/health

# Should return:
# {"status":"healthy","version":"0.1.0",...}
```

### Step 3: Integrate with Your Agent

Use the REST API endpoints to:
- Discover tools: `GET /api/v1/tools`
- Get tool details: `GET /api/v1/tools/{id}`
- Execute tools: `POST /api/v1/tools/{id}/invoke`

## API Reference

### List All Tools

```bash
GET /api/v1/tools
```

**Response:**
```json
{
  "protocol": "1.0",
  "tools": [
    {
      "name": "echo",
      "description": "Echoes back the input message",
      "schema": {
        "input": {
          "properties": {
            "message": { "type": "string" }
          }
        }
      }
    },
    {
      "name": "petstore.listPets",
      "description": "List all available pets",
      "schema": {
        "input": {
          "properties": {
            "limit": { "type": "integer" }
          }
        }
      }
    }
  ]
}
```

### Get Tool Details

```bash
GET /api/v1/tools/{tool-name}
```

Example:
```bash
curl http://localhost:8080/api/v1/tools/echo
```

**Response:**
```json
{
  "name": "echo",
  "description": "Echoes back the input message for testing purposes",
  "version": "1.0.0",
  "source": "builtin",
  "tags": ["test", "utility"],
  "schema": {
    "input": {
      "type": "object",
      "properties": {
        "message": {
          "type": "string",
          "description": "Message to echo back"
        }
      }
    }
  }
}
```

### Execute a Tool

```bash
POST /api/v1/tools/{tool-name}/invoke
Content-Type: application/json

{
  "message": "Hello from GitHub Copilot!"
}
```

Example:
```bash
curl -X POST http://localhost:8080/api/v1/tools/echo/invoke \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello from Copilot!"}'
```

**Response:**
```json
{
  "success": true,
  "result": {
    "message": "Hello from Copilot!"
  },
  "execution_time_ms": 2
}
```

## Integration Examples

### JavaScript/TypeScript (GitHub Copilot Extension)

```typescript
// Service class for AionMCP integration
class AionMCPClient {
  private baseUrl = 'http://localhost:8080/api/v1';

  async listTools() {
    const response = await fetch(`${this.baseUrl}/tools`);
    const data = await response.json();
    return data.tools;
  }

  async getToolSchema(toolName: string) {
    const response = await fetch(`${this.baseUrl}/tools/${toolName}`);
    return response.json();
  }

  async executeTool(toolName: string, params: any) {
    const response = await fetch(`${this.baseUrl}/tools/${toolName}/invoke`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(params)
    });
    return response.json();
  }
}

// Usage in GitHub Copilot context
const client = new AionMCPClient();

// 1. Discover tools available to Copilot
const availableTools = await client.listTools();
console.log('Available tools:', availableTools.map(t => t.name));

// 2. Get a specific tool's schema
const echoToolSchema = await client.getToolSchema('echo');
console.log('Echo tool schema:', echoToolSchema.schema);

// 3. Execute a tool
const result = await client.executeTool('echo', {
  message: 'Hello from GitHub Copilot!'
});
console.log('Result:', result);
```

### Python (Custom Agent)

```python
import requests
import json

class AionMCPClient:
    def __init__(self, base_url='http://localhost:8080/api/v1'):
        self.base_url = base_url
        self.session = requests.Session()
    
    def list_tools(self):
        """Get all available tools"""
        response = self.session.get(f'{self.base_url}/tools')
        response.raise_for_status()
        return response.json()['tools']
    
    def get_tool_schema(self, tool_name):
        """Get a specific tool's schema"""
        response = self.session.get(f'{self.base_url}/tools/{tool_name}')
        response.raise_for_status()
        return response.json()
    
    def execute_tool(self, tool_name, params):
        """Execute a tool with given parameters"""
        response = self.session.post(
            f'{self.base_url}/tools/{tool_name}/invoke',
            json=params
        )
        response.raise_for_status()
        return response.json()

# Usage example
if __name__ == '__main__':
    client = AionMCPClient()
    
    # List all tools
    tools = client.list_tools()
    print(f"Available tools: {[t['name'] for t in tools]}")
    
    # Get tool schema
    schema = client.get_tool_schema('echo')
    print(f"Echo tool schema: {json.dumps(schema, indent=2)}")
    
    # Execute a tool
    result = client.execute_tool('echo', {'message': 'Hello from Python!'})
    print(f"Result: {result}")
```

### Go (Embed in Copilot)

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type AionMCPClient struct {
    BaseURL string
}

type Tool struct {
    Name        string `json:"name"`
    Description string `json:"description"`
}

type ToolsResponse struct {
    Tools []Tool `json:"tools"`
}

type ExecutionResult struct {
    Success bool        `json:"success"`
    Result  interface{} `json:"result"`
}

func (c *AionMCPClient) ListTools() ([]Tool, error) {
    resp, err := http.Get(fmt.Sprintf("%s/tools", c.BaseURL))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var tr ToolsResponse
    json.NewDecoder(resp.Body).Decode(&tr)
    return tr.Tools, nil
}

func (c *AionMCPClient) ExecuteTool(toolName string, params interface{}) (*ExecutionResult, error) {
    body, _ := json.Marshal(params)
    resp, err := http.Post(
        fmt.Sprintf("%s/tools/%s/invoke", c.BaseURL, toolName),
        "application/json",
        bytes.NewBuffer(body),
    )
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result ExecutionResult
    json.NewDecoder(resp.Body).Decode(&result)
    return &result, nil
}

func main() {
    client := &AionMCPClient{BaseURL: "http://localhost:8080/api/v1"}
    
    // List tools
    tools, _ := client.ListTools()
    fmt.Println("Available tools:", tools)
    
    // Execute tool
    result, _ := client.ExecuteTool("echo", map[string]string{
        "message": "Hello from Go!",
    })
    fmt.Println("Result:", result)
}
```

### cURL (Direct Commands)

```bash
# List all tools
curl http://localhost:8080/api/v1/tools | jq .

# List just tool names
curl http://localhost:8080/api/v1/tools | jq '.tools[].name'

# Get specific tool details
curl http://localhost:8080/api/v1/tools/echo | jq .

# Execute a tool
curl -X POST http://localhost:8080/api/v1/tools/echo/invoke \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Test message"
  }'
```

## Importing API Specifications

Before using tools with your APIs, import the specifications:

### Via VS Code Extension (Easiest)
1. Open the AionMCP sidebar in VS Code
2. Click the **file icon** in the Tools section
3. Select the spec type (OpenAPI, GraphQL, AsyncAPI)
4. Provide the file path or URL
5. The tools automatically appear and are callable

### Via REST API
```bash
curl -X POST http://localhost:8080/api/v1/import-spec \
  -H "Content-Type: application/json" \
  -d '{
    "type": "openapi",
    "path": "./examples/specs/petstore.yaml",
    "name": "Petstore API"
  }'
```

## Use Case: Chat with GitHub Copilot Using AionMCP Tools

### Scenario
User: "List all available pets from the Petstore API"

### Flow
1. **Copilot recognizes** the request needs tool access
2. **Copilot queries** AionMCP: `GET /api/v1/tools`
3. **Copilot identifies** relevant tools (e.g., `petstore.listPets`)
4. **Copilot gets schema**: `GET /api/v1/tools/petstore.listPets`
5. **Copilot executes**: `POST /api/v1/tools/petstore.listPets/invoke`
6. **Copilot returns** results to user

### Code in Copilot Extension
```typescript
async function handleUserRequest(userQuery: string) {
  const tools = await client.listTools();
  
  // Find relevant tool based on query
  const relevantTool = tools.find(t => 
    t.description.toLowerCase().includes('pets')
  );
  
  if (relevantTool) {
    const schema = await client.getToolSchema(relevantTool.name);
    
    // Execute with appropriate parameters
    const result = await client.executeTool(relevantTool.name, {
      limit: 10
    });
    
    return `Here are the available pets: ${JSON.stringify(result.result)}`;
  }
}
```

## Model Independence Benefits

✅ **Not Coupled to Claude Desktop**: Works with any agent/model  
✅ **Simple HTTP**: Any programming language can integrate  
✅ **No Model-Specific Config**: Standard REST API  
✅ **Easy Switching**: Change models without rewriting integration  
✅ **Future Ready**: New LLM platforms integrate easily  
✅ **Scalable**: Can deploy to cloud with multiple instances  

## Deployment Options

### Local Development (Default)
```bash
./bin/aionmcp-server
# Access: http://localhost:8080
```

### Docker
```bash
docker run -p 8080:8080 aionmcp:latest
```

### Cloud (AWS Lambda, Google Cloud Run, etc.)
```bash
# Deploy anywhere that supports Go
# Access via public URL instead of localhost
```

### Behind Reverse Proxy (Production)
```bash
# Use nginx/Apache to proxy requests
# Add authentication/authorization layer
# Enable HTTPS for security
```

## Security Considerations

### Development (Current)
- ✅ Running on localhost only
- ✅ No authentication required (local)

### Production (Recommended)
- 🔒 Add API key authentication
- 🔒 Use HTTPS/TLS
- 🔒 Add rate limiting
- 🔒 Implement request validation
- 🔒 Add execution timeouts
- 🔒 Enable CORS with specific origins

## Troubleshooting

### Server not running
```bash
# Check if server is listening
lsof -i :8080

# If not, start it
cd /Users/kiran/Documents/GitHub/aionmcp
./bin/aionmcp-server
```

### Can't connect from Copilot
```bash
# Verify server is accessible
curl http://localhost:8080/api/v1/health

# Check firewall rules
sudo lsof -i :8080
```

### Tool execution fails
```bash
# Check tool schema first
curl http://localhost:8080/api/v1/tools/{tool-name}

# Verify parameters match schema
# Check server logs for detailed error info
```

### No tools available
```bash
# Import a spec first
curl http://localhost:8080/api/v1/tools
# Should show at least 'echo' and 'status' built-in tools

# Import an API spec
curl -X POST http://localhost:8080/api/v1/import-spec \
  -H "Content-Type: application/json" \
  -d '{
    "type": "openapi",
    "path": "./examples/specs/petstore.yaml",
    "name": "Petstore API"
  }'
```

## Next Steps

1. **Install & start** the AionMCP server
2. **Import an API spec** (OpenAPI, GraphQL, or AsyncAPI)
3. **Integrate** with your agent using the REST API
4. **Test** with sample queries
5. **Deploy** to production when ready

## Additional Resources

- [Model-Independent Tools Architecture](./MODEL_INDEPENDENT_TOOLS.md)
- [API Documentation](./README.md#api-endpoints)
- [Quick Start Guide](./QUICK_START.md)
- [Examples](../examples/)

## Support

- **Questions**: Check the examples above
- **Issues**: Report on GitHub
- **Contributing**: Fork and submit PRs

---

**Key Principle**: AionMCP exposes tools via standard HTTP REST API. Any agent (GitHub Copilot, Claude, custom agents, etc.) can use them without model-specific configuration.

**Get Started**: `curl http://localhost:8080/api/v1/tools`
