# Step-by-Step Testing with Public APIs

This guide walks through testing AionMCP with real open-source APIs using the VS Code extension and command line.

## Part 1: Verify Server is Running (5 minutes)

### Step 1: Open VS Code
```bash
cd /Users/kiran/Documents/GitHub/aionmcp
code .
```

### Step 2: Start AionMCP Server
1. Press `Cmd+Shift+P` (macOS) or `Ctrl+Shift+P` (Linux/Windows)
2. Search for "AionMCP: Start Server"
3. Click the command

You should see:
- Output panel opens showing server logs
- `[SUCCESS] AionMCP server is running and connected.`
- Status bar shows: `$(server-process) AionMCP $(check)`

### Step 3: Verify Server Health
Open terminal and run:
```bash
curl -s http://localhost:8080/ | jq '.'
```

Expected response:
```json
{
  "status": "AionMCP server is running"
}
```

## Part 2: Test Built-in Tools (5 minutes)

### Step 1: List Available Tools
```bash
curl -s http://localhost:8080/api/v1/tools | jq '.tools[] | {name: .name, description: .description}'
```

Expected output:
```json
[
  {
    "name": "echo",
    "description": "Echoes back the input message for testing purposes"
  },
  {
    "name": "status",
    "description": "Returns information about the tool registry and server status"
  }
]
```

### Step 2: Execute Echo Tool
```bash
curl -X POST http://localhost:8080/api/v1/tools/echo/execute \
  -H "Content-Type: application/json" \
  -d '{"args": {"message": "Hello from AionMCP"}}'
```

Expected response:
```json
{
  "result": "Echo: Hello from AionMCP"
}
```

### Step 3: Check Server Status
```bash
curl -X POST http://localhost:8080/api/v1/tools/status/execute \
  -H "Content-Type: application/json" \
  -d '{}'
```

## Part 3: Test with JSONPlaceholder API (15 minutes)

JSONPlaceholder provides a free fake REST API for testing. No authentication needed.

### Step 1: Create OpenAPI Specification

Save this as `jsonplaceholder-spec.yaml`:

```yaml
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
      parameters:
        - name: _limit
          in: query
          description: Limit results
          schema:
            type: integer
      responses:
        '200':
          description: Array of posts

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

  /comments:
    get:
      operationId: listComments
      summary: Get comments
      tags:
        - Comments
      parameters:
        - name: postId
          in: query
          description: Filter by post ID
          schema:
            type: integer
      responses:
        '200':
          description: Array of comments
```

### Step 2: Import the Spec

The import feature is being developed. For now, test the API directly:

```bash
# Direct test: List all posts (limited to 5)
curl -s 'https://jsonplaceholder.typicode.com/posts?_limit=5' | jq '.'
```

Output:
```json
[
  {
    "userId": 1,
    "id": 1,
    "title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
    "body": "quia et suscipit..."
  },
  ...
]
```

### Step 3: Test Different Endpoints

```bash
# Get specific post
curl -s 'https://jsonplaceholder.typicode.com/posts/1' | jq '.'

# Get users
curl -s 'https://jsonplaceholder.typicode.com/users' | jq '.[0]'

# Get comments for post 1
curl -s 'https://jsonplaceholder.typicode.com/comments?postId=1' | jq '.'
```

## Part 4: Test with REST Countries API (15 minutes)

### Step 1: Get All Countries

```bash
# Get all countries (returns large result)
curl -s 'https://restcountries.com/v3.1/all' | jq '.[0] | {name: .name.common, capital: .capital, region: .region}'
```

Output example:
```json
{
  "name": "Afghanistan",
  "capital": [
    "Kabul"
  ],
  "region": "Asia"
}
```

### Step 2: Search by Name

```bash
# Find countries by name
curl -s 'https://restcountries.com/v3.1/name/japan' | jq '.[] | {name: .name.common, population: .population, capital: .capital}'
```

Output:
```json
{
  "name": "Japan",
  "population": 125124989,
  "capital": [
    "Tokyo"
  ]
}
```

### Step 3: Filter by Region

```bash
# Get countries in a region
curl -s 'https://restcountries.com/v3.1/region/Asia' | jq '.[].name.common' | head -10
```

### Step 4: Get by Currency

```bash
# Find countries by currency code
curl -s 'https://restcountries.com/v3.1/currency/usd' | jq '.[].name.common' | head -5
```

## Part 5: Test with PokéAPI (10 minutes)

### Step 1: List Pokémon

```bash
# Get first 5 Pokémon
curl -s 'https://pokeapi.co/api/v2/pokemon?limit=5' | jq '.results'
```

Output:
```json
[
  {
    "name": "bulbasaur",
    "url": "https://pokeapi.co/api/v2/pokemon/1/"
  },
  ...
]
```

### Step 2: Get Pokémon Details

```bash
# Get detailed info on a Pokémon
curl -s 'https://pokeapi.co/api/v2/pokemon/pikachu' | jq '{
  name: .name,
  height: .height,
  weight: .weight,
  types: [.types[].type.name],
  abilities: [.abilities[].ability.name]
}'
```

Output:
```json
{
  "name": "pikachu",
  "height": 4,
  "weight": 60,
  "types": ["electric"],
  "abilities": ["static", "lightning-rod"]
}
```

### Step 3: Get Pokémon Type Info

```bash
# List all Pokémon types
curl -s 'https://pokeapi.co/api/v2/type' | jq '.results[] | .name' | head -10
```

## Part 6: Create a Test Agent (20 minutes)

### Step 1: Create `test_agent.py`

```python
#!/usr/bin/env python3
"""
Test agent that uses AionMCP with public APIs
"""
import requests
import json

class TestAgent:
    def __init__(self, server_url="http://localhost:8080"):
        self.server_url = server_url
        self.api = f"{server_url}/api/v1"
    
    def test_builtin_tools(self):
        """Test built-in AionMCP tools"""
        print("=== Testing Built-in Tools ===\n")
        
        # Test echo
        resp = requests.post(
            f"{self.api}/tools/echo/execute",
            json={"args": {"message": "Test message"}}
        )
        print(f"Echo: {resp.json()}")
        
        # Test status
        resp = requests.post(
            f"{self.api}/tools/status/execute",
            json={}
        )
        print(f"Status: {resp.json()}\n")
    
    def test_jsonplaceholder(self):
        """Test with JSONPlaceholder API"""
        print("=== Testing JSONPlaceholder ===\n")
        
        # Get posts
        resp = requests.get("https://jsonplaceholder.typicode.com/posts?_limit=3")
        posts = resp.json()
        print(f"Found {len(posts)} posts (limited to 3)")
        for post in posts:
            print(f"  - [{post['id']}] {post['title']}")
        print()
    
    def test_rest_countries(self):
        """Test with REST Countries API"""
        print("=== Testing REST Countries ===\n")
        
        # Get country
        resp = requests.get("https://restcountries.com/v3.1/name/japan")
        country = resp.json()[0]
        print(f"Country: {country['name']['common']}")
        print(f"Population: {country['population']:,}")
        print(f"Region: {country['region']}")
        print(f"Subregion: {country['subregion']}\n")
    
    def test_pokemon_api(self):
        """Test with PokéAPI"""
        print("=== Testing PokéAPI ===\n")
        
        # Get Pokémon
        resp = requests.get("https://pokeapi.co/api/v2/pokemon/charizard")
        pokemon = resp.json()
        print(f"Name: {pokemon['name'].title()}")
        print(f"Height: {pokemon['height'] / 10} m")
        print(f"Weight: {pokemon['weight'] / 10} kg")
        print(f"Types: {', '.join([t['type']['name'] for t in pokemon['types']])}")
        print()

if __name__ == "__main__":
    agent = TestAgent()
    
    try:
        agent.test_builtin_tools()
        agent.test_jsonplaceholder()
        agent.test_rest_countries()
        agent.test_pokemon_api()
        print("✓ All tests completed successfully!")
    except Exception as e:
        print(f"✗ Error: {e}")
```

### Step 2: Run the Agent

```bash
# Install requests if needed
pip3 install requests

# Run the test agent
python3 test_agent.py
```

Expected output:
```
=== Testing Built-in Tools ===

Echo: {'result': 'Echo: Test message'}
Status: {'tool_count': 2, 'uptime_seconds': 123}

=== Testing JSONPlaceholder ===

Found 3 posts (limited to 3)
  - [1] sunt aut facere repellat provident occaecati excepturi optio reprehenderit
  - [2] qui est esse
  - [3] ea molestias quasi exercitationem repudiandae et enim quasi cinnam

=== Testing REST Countries ===

Country: Japan
Population: 125,124,989
Region: Asia
Subregion: Eastern Asia

=== Testing PokéAPI ===

Name: Charizard
Height: 1.7 m
Weight: 90 kg
Types: flying, fire

✓ All tests completed successfully!
```

## Part 7: Monitor Server Activity (5 minutes)

### Check Server Logs
```bash
# In VS Code: Select "AionMCP Server" output tab
# Or in terminal:
tail -f /tmp/aionmcp.log
```

### View Tool Metrics
```bash
# Check which tools have been used
curl -s http://localhost:8080/api/v1/tools | jq '.tools | length'
```

### Stop Server
```bash
# From VS Code: Cmd+Shift+P > "AionMCP: Stop Server"
# Or: Cmd+C in terminal
# Or: killall server
```

## Verification Checklist

- [ ] Server starts and logs show "listening on :8080"
- [ ] `curl http://localhost:8080/` returns healthy status
- [ ] Built-in echo and status tools work
- [ ] JSONPlaceholder API data retrieves successfully
- [ ] REST Countries API searches work
- [ ] PokéAPI Pokémon lookups succeed
- [ ] Test agent runs without errors
- [ ] VS Code extension controls server successfully
- [ ] Server logs show tool execution
- [ ] Learning engine records metrics

## Troubleshooting

### Connection refused
```bash
# Check if server is running
curl -v http://localhost:8080/

# Start server if not running
cd /Users/kiran/Documents/GitHub/aionmcp
./server
```

### Port already in use
```bash
# Find what's using port 8080
lsof -i :8080

# Kill if it's an old server
killall server

# Use different port
HTTP_PORT=8081 ./server
```

### Python agent errors
```bash
# Check dependencies
pip3 install requests

# Run with error details
python3 test_agent.py 2>&1
```

## Next Steps

1. Import custom OpenAPI specs
2. Set up multi-agent orchestration
3. Monitor learning insights
4. Analyze execution patterns
5. Generate reflection documents

---

**Total Testing Time**: ~60 minutes
**Date**: November 19, 2025
