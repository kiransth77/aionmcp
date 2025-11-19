# Testing AionMCP with Open Source APIs

A step-by-step guide to test AionMCP with real open-source APIs and agents.

## Quick Start: JSONPlaceholder (Recommended First Test)

### Step 1: Start the AionMCP Server

```bash
cd /Users/kiran/Documents/GitHub/aionmcp
./server
```

Expected output:
```json
{"level":"info","ts":1763548667.182084,"msg":"Starting AionMCP server"}
{"level":"info","ts":1763548667.182190,"msg":"Tool registered","tool":"echo"}
{"level":"info","ts":1763548667.182202,"msg":"Tool registered","tool":"status"}
{"level":"info","ts":1763548667.183100,"msg":"AionMCP server started successfully"}
```

Server should be listening on:
- HTTP: `http://localhost:8080`
- gRPC: `localhost:9090`

### Step 2: Test Health Check

```bash
curl -s http://localhost:8080/ | jq '.'
```

Expected response:
```json
{
  "status": "AionMCP server is running"
}
```

### Step 3: List Available Tools

```bash
curl -s http://localhost:8080/api/v1/tools | jq '.'
```

You should see the built-in tools:
- `echo` - Test echo tool
- `status` - Server status tool

### Step 4: Import JSONPlaceholder API Spec

Create the OpenAPI specification:

```bash
cat > /tmp/jsonplaceholder-spec.yaml << 'EOF'
openapi: 3.0.0
info:
  title: JSONPlaceholder API
  version: 1.0.0
  description: Free fake REST API for testing
servers:
  - url: https://jsonplaceholder.typicode.com

paths:
  /posts:
    get:
      operationId: listPosts
      summary: Get all posts
      tags:
        - Posts
      responses:
        '200':
          description: Array of posts
          content:
            application/json:
              schema:
                type: array

  /posts/{postId}:
    get:
      operationId: getPost
      summary: Get a specific post
      tags:
        - Posts
      parameters:
        - name: postId
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Post object

  /users:
    get:
      operationId: listUsers
      summary: Get all users
      tags:
        - Users
      responses:
        '200':
          description: Array of users

  /users/{userId}:
    get:
      operationId: getUser
      summary: Get a specific user
      tags:
        - Users
      parameters:
        - name: userId
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: User object

  /comments:
    get:
      operationId: listComments
      summary: Get all comments
      tags:
        - Comments
      parameters:
        - name: postId
          in: query
          schema:
            type: integer
          description: Filter by post ID
      responses:
        '200':
          description: Array of comments
EOF
```

### Step 5: Execute Tools from the Spec

```bash
# List all posts
curl -X POST http://localhost:8080/api/v1/tools/listPosts/execute \
  -H "Content-Type: application/json" \
  -d '{}' | jq '.' | head -30

# Get a specific post
curl -X POST http://localhost:8080/api/v1/tools/getPost/execute \
  -H "Content-Type: application/json" \
  -d '{"postId": 1}' | jq '.'

# List all users
curl -X POST http://localhost:8080/api/v1/tools/listUsers/execute \
  -H "Content-Type: application/json" \
  -d '{}' | jq '.' | head -20

# Get comments for a post
curl -X POST http://localhost:8080/api/v1/tools/listComments/execute \
  -H "Content-Type: application/json" \
  -d '{"postId": 1}' | jq '.'
```

## Testing with REST Countries API

### Step 1: Create OpenAPI Spec for REST Countries

```bash
cat > /tmp/countries-spec.yaml << 'EOF'
openapi: 3.0.0
info:
  title: REST Countries API
  version: 3.1.0
  description: Get information about countries
servers:
  - url: https://restcountries.com/v3.1

paths:
  /all:
    get:
      operationId: getAllCountries
      summary: Get all countries
      responses:
        '200':
          description: Array of all countries

  /name/{name}:
    get:
      operationId: getCountriesByName
      summary: Search countries by name
      parameters:
        - name: name
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Array of matching countries

  /region/{region}:
    get:
      operationId: getCountriesByRegion
      summary: Get countries by region
      parameters:
        - name: region
          in: path
          required: true
          schema:
            type: string
            enum: [Africa, Americas, Asia, Europe, Oceania]
      responses:
        '200':
          description: Array of countries in region

  /currency/{code}:
    get:
      operationId: getCountriesByCurrency
      summary: Get countries using a currency
      parameters:
        - name: code
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Array of countries using currency
EOF
```

### Step 2: Test REST Countries Tools

```bash
# Get all countries (limit output)
curl -s http://localhost:8080/api/v1/tools/getAllCountries/execute \
  -H "Content-Type: application/json" \
  -d '{}' | jq '.[0:2]'

# Search for specific country
curl -s http://localhost:8080/api/v1/tools/getCountriesByName/execute \
  -H "Content-Type: application/json" \
  -d '{"name": "united"}' | jq '.'

# Get countries by region
curl -s http://localhost:8080/api/v1/tools/getCountriesByRegion/execute \
  -H "Content-Type: application/json" \
  -d '{"region": "Asia"}' | jq '.[].name'

# Get countries by currency
curl -s http://localhost:8080/api/v1/tools/getCountriesByCurrency/execute \
  -H "Content-Type: application/json" \
  -d '{"code": "USD"}' | jq '.[].name'
```

## Testing with PokéAPI

```bash
cat > /tmp/pokemon-spec.yaml << 'EOF'
openapi: 3.0.0
info:
  title: PokéAPI
  version: 2.0.0
  description: Pokémon data API
servers:
  - url: https://pokeapi.co/api/v2

paths:
  /pokemon:
    get:
      operationId: listPokemon
      summary: List Pokémon
      parameters:
        - name: limit
          in: query
          schema:
            type: integer
            default: 20
        - name: offset
          in: query
          schema:
            type: integer
            default: 0
      responses:
        '200':
          description: Array of Pokémon

  /pokemon/{id}:
    get:
      operationId: getPokemon
      summary: Get Pokémon details
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Pokémon details

  /type:
    get:
      operationId: listTypes
      summary: List Pokémon types
      responses:
        '200':
          description: Array of types
EOF
```

Test it:
```bash
# List first 5 Pokémon
curl -s http://localhost:8080/api/v1/tools/listPokemon/execute \
  -H "Content-Type: application/json" \
  -d '{"limit": 5}' | jq '.results[] | .name'

# Get specific Pokémon
curl -s http://localhost:8080/api/v1/tools/getPokemon/execute \
  -H "Content-Type: application/json" \
  -d '{"id": "pikachu"}' | jq '.name, .types'
```

## Testing with VS Code Extension

### Step 1: Ensure Server is Running

```bash
cd /Users/kiran/Documents/GitHub/aionmcp && ./server &
```

### Step 2: Open the AionMCP Dashboard

1. Open VS Code
2. Press `Cmd+Shift+P` (or Ctrl+Shift+P on Linux/Windows)
3. Search for "AionMCP: Open Dashboard"
4. View the dashboard with tools and agents

### Step 3: Refresh Tools

1. In the "Tools" sidebar, click the refresh button
2. Select a tool to view its details
3. Double-click a tool to open the executor

### Step 4: Execute Tool from Extension

1. Open the tool executor for any tool
2. Enter parameters if required
3. Click "Execute"
4. View results in the output panel

## Monitoring Server Activity

### View Server Logs

```bash
# Follow logs in real-time
tail -f /tmp/aionmcp.log
```

### Check Imported Specifications

```bash
curl -s http://localhost:8080/api/v1/specs | jq '.'
```

### View Tool Registry

```bash
curl -s http://localhost:8080/api/v1/tools | jq '.tools | length'
```

## Performance Testing

### Test Large Data Retrieval

```bash
# Time the request
time curl -s http://localhost:8080/api/v1/tools/listPosts/execute \
  -H "Content-Type: application/json" \
  -d '{}' > /dev/null
```

### Concurrent Requests

```bash
# Make 10 concurrent requests
for i in {1..10}; do
  curl -s http://localhost:8080/api/v1/tools/listPosts/execute \
    -H "Content-Type: application/json" \
    -d '{}' &
done
wait
echo "All requests completed"
```

## Troubleshooting

### Server not responding
```bash
# Check if server is running
lsof -i :8080
ps aux | grep server

# Restart server
killall server 2>/dev/null
./server
```

### Tools not appearing
```bash
# Check tool registry
curl -s http://localhost:8080/api/v1/tools | jq '.tools[]' | head -20
```

### Tool execution failing
```bash
# Check server logs for errors
tail -50 /tmp/aionmcp.log
```

## Next Steps

1. Create a custom agent that uses these tools
2. Set up the learning engine to analyze patterns
3. Monitor reflections in `docs/reflections/`
4. Build a specialized API aggregator using AionMCP

## References

- **JSONPlaceholder**: https://jsonplaceholder.typicode.com
- **REST Countries**: https://restcountries.com
- **PokéAPI**: https://pokeapi.co
- **AionMCP API**: See `docs/api.md`
