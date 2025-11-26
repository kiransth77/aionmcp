# Setting Up GitHub Copilot to Use AionMCP

## Overview

This guide explains how to configure GitHub Copilot (in VS Code) to use AionMCP as a tool server for accessing APIs and other services.

## Architecture

```
GitHub Copilot (VS Code)
         │
         │ Uses tool APIs
         ▼
    AionMCP Server (HTTP REST)
         │
         │ Manages
         ▼
    Tool Registry
         │
    ┌────┴──────────┐
    │               │
    ▼               ▼
  OpenAPI      GraphQL/AsyncAPI
    APIs           Specs
```

## Prerequisites

1. **VS Code** with GitHub Copilot extension installed
2. **AionMCP Server** running locally on port 8080
3. **API Specifications** imported into AionMCP

## Step 1: Start AionMCP Server

### Option A: Via VS Code Extension (Easiest)

1. Open VS Code
2. Look for "AionMCP" in the left sidebar
3. Click the **▶️ Play** button in the Server Status section
4. Wait for "Server Status: Healthy" message

### Option B: Manual Terminal

```bash
cd /Users/kiran/Documents/GitHub/aionmcp
./bin/aionmcp-server &
```

### Step 2: Verify Server is Running

Open a terminal and run:
```bash
curl http://localhost:8080/api/v1/health

# Should return:
# {"status":"healthy","version":"0.1.0",...}
```

## Step 2: Import API Specifications

You need to make tools available in AionMCP before Copilot can use them.

### Option A: Via VS Code Extension

1. In VS Code AionMCP sidebar, go to **Tools** section
2. Click the **file icon** (📄) to import
3. Select specification type: `openapi`, `graphql`, or `asyncapi`
4. Enter the file path or URL to your API spec
5. Optional: Give it a friendly name
6. Click OK

The tools now appear in the Tools tree!

### Option B: Via REST API

```bash
curl -X POST http://localhost:8080/api/v1/import-spec \
  -H "Content-Type: application/json" \
  -d '{
    "type": "openapi",
    "path": "./examples/specs/petstore.yaml",
    "name": "Petstore API"
  }'
```

### Example: Import the Petstore API

```bash
curl -X POST http://localhost:8080/api/v1/import-spec \
  -H "Content-Type: application/json" \
  -d '{
    "type": "openapi",
    "path": "/Users/kiran/Documents/GitHub/aionmcp/examples/specs/petstore.yaml",
    "name": "Petstore API"
  }'
```

## Step 3: Verify Tools are Available

Check that tools are now accessible:

```bash
# List all tools
curl http://localhost:8080/api/v1/tools | jq '.tools[].name'

# Should show something like:
# echo
# status
# openapi.petstore.listPets
# openapi.petstore.getPet
# ... etc
```

## Step 4: Using Tools with GitHub Copilot

### In a Copilot Chat

GitHub Copilot can now use these tools! In a Copilot chat window, you can:

**Example 1: Ask Copilot to use a tool**
```
You: "Can you list all available pets using the petstore API?"

Copilot:
1. Discovers available tools from AionMCP
2. Identifies the "listPets" tool
3. Calls GET /api/v1/tools/openapi.petstore.listPets
4. Executes the tool with appropriate parameters
5. Returns the results to you

Response: "Here are the available pets: [...]"
```

**Example 2: Use multiple tools in conversation**
```
You: "Show me pet #1 and then tell me about a similar dog breed"

Copilot:
1. Executes getPet tool with id=1
2. Gets the result (e.g., a "Dachshund")
3. Answers your question about dog breeds
```

### In Code Files

You can also reference tools in code:

```typescript
// In a .ts file, Copilot can help you integrate tools
const listPets = async () => {
  const response = await fetch('http://localhost:8080/api/v1/tools/openapi.petstore.listPets/invoke', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ limit: 10 })
  });
  return response.json();
};

// Copilot can help you understand and use the API
// Just ask questions about the tools!
```

## Configuring Copilot Settings (Optional)

You can enhance Copilot's tool discovery by creating a configuration file:

### Create `.copilot/aionmcp-tools.json`

```json
{
  "toolServer": {
    "baseUrl": "http://localhost:8080",
    "apiVersion": "v1",
    "endpoints": {
      "listTools": "/api/v1/tools",
      "getTool": "/api/v1/tools/{id}",
      "executeTool": "/api/v1/tools/{id}/invoke"
    },
    "discoveryInterval": 5000
  }
}
```

This helps Copilot:
- ✅ Discover tools periodically
- ✅ Know where to find the API
- ✅ Understand the endpoint patterns
- ✅ Execute tools efficiently

## Common Use Cases

### 1. REST API Queries

```
You: "Using the Petstore API, find all pets with the tag 'available'"

Copilot:
- Discovers: openapi.petstore.listPets
- Executes with filter parameters
- Returns filtered results
```

### 2. Data Transformation

```
You: "Get the top 5 most expensive pets and format them as a table"

Copilot:
- Uses petstore API to get pets
- Sorts and filters
- Formats as markdown table
```

### 3. Multi-Step Workflows

```
You: "Create a new pet in the store and then confirm it was created"

Copilot:
- Executes createPet tool
- Uses getPet to verify
- Confirms success
```

### 4. Integration with Code

```
You: "Generate JavaScript code to list all pets"

Copilot:
- Understands the tool schema
- Generates proper fetch/axios code
- Includes error handling
```

## Troubleshooting

### "Copilot can't find tools"

**Solution:**
1. Verify server is running: `curl http://localhost:8080/api/v1/health`
2. Check tools are imported: `curl http://localhost:8080/api/v1/tools`
3. Restart VS Code
4. Try importing a spec again

### "Tools are listed but Copilot won't execute them"

**Solution:**
1. Check tool schema: `curl http://localhost:8080/api/v1/tools/{tool-name}`
2. Verify parameters are correct
3. Try executing manually: 
   ```bash
   curl -X POST http://localhost:8080/api/v1/tools/{tool-name}/invoke \
     -H "Content-Type: application/json" \
     -d '{}'
   ```
4. Check AionMCP server logs for errors

### "Tool execution returns errors"

**Solution:**
1. Verify the API is accessible
2. Check tool parameters match API requirements
3. Look at AionMCP logs for details:
   ```bash
   # In VS Code AionMCP sidebar, click "View Server Logs"
   ```

### "Can't import API specification"

**Solution:**
1. Verify the file path is correct (use absolute path)
2. Ensure the spec is valid YAML/JSON
3. Try via curl to see exact error:
   ```bash
   curl -X POST http://localhost:8080/api/v1/import-spec \
     -H "Content-Type: application/json" \
     -d '{"type":"openapi","path":"/path/to/spec.yaml"}'
   ```

## Advanced: Custom Tool Server

You can create custom tools beyond API specs:

### Create a Custom Tool

```go
package main

type CustomTool struct {
    name string
}

func (t *CustomTool) Execute(params map[string]interface{}) (interface{}, error) {
    // Custom logic here
    return result, nil
}
```

### Register with AionMCP

The tool is automatically available to Copilot via the REST API!

## Best Practices

1. **Start Small**: Begin with one API spec
2. **Test Manually**: Verify tools work before using with Copilot
3. **Document Tools**: Add clear descriptions to help Copilot
4. **Monitor Performance**: Check execution times, adjust as needed
5. **Keep Updated**: Refresh API specs as they change
6. **Use Error Handling**: Expect tool calls to fail, handle gracefully

## Security Notes

### Development (Current)
- ✅ Running on localhost (127.0.0.1)
- ✅ No authentication required
- ✅ Safe for local use only

### Production
- 🔒 Deploy behind authentication (API keys)
- 🔒 Use HTTPS/TLS for transport
- 🔒 Add rate limiting
- 🔒 Implement request validation
- 🔒 Add audit logging

## Next Steps

1. **Start the AionMCP server** (Step 1)
2. **Import an API spec** (Step 2)
3. **Open GitHub Copilot chat** in VS Code
4. **Ask Copilot to use a tool**
5. **Watch it work!**

## Example Commands

### Terminal Setup
```bash
# Terminal 1: Start AionMCP
cd ~/Documents/GitHub/aionmcp
./bin/aionmcp-server

# Terminal 2: Import Petstore API
curl -X POST http://localhost:8080/api/v1/import-spec \
  -H "Content-Type: application/json" \
  -d '{
    "type": "openapi",
    "path": "./examples/specs/petstore.yaml",
    "name": "Petstore"
  }'

# Terminal 3: Verify tools
curl http://localhost:8080/api/v1/tools | jq '.tools[].name'
```

### VS Code Setup
1. Open AionMCP sidebar
2. Click ▶️ to start server
3. Click 📄 to import spec
4. Open Copilot chat
5. Use tools!

## Support & Resources

- **GitHub Copilot Docs**: https://github.com/features/copilot
- **VS Code Copilot Extension**: https://marketplace.visualstudio.com/items?itemName=GitHub.copilot
- **AionMCP Docs**: See [MODEL_INDEPENDENT_TOOLS.md](./MODEL_INDEPENDENT_TOOLS.md)
- **API Reference**: See [GITHUB_COPILOT_INTEGRATION.md](./GITHUB_COPILOT_INTEGRATION.md)

---

**Quick Summary**:
1. Start AionMCP server
2. Import API specs
3. Use with GitHub Copilot via REST API
4. Enjoy model-independent tool access!

**Get Started**: `cd /Users/kiran/Documents/GitHub/aionmcp && ./bin/aionmcp-server`
