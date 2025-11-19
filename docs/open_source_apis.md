# Open Source APIs for AionMCP Testing

This document lists recommended open-source APIs that can be used to test and demonstrate AionMCP's capabilities.

## Free Public APIs (No Authentication Required)

### 1. **JSONPlaceholder** (Fake Online REST API)
- **URL**: https://jsonplaceholder.typicode.com
- **OpenAPI Spec**: https://github.com/jsmini/json-placeholder-swagger
- **Description**: Perfect for testing - provides fake data for posts, comments, users, etc.
- **Features**:
  - GET, POST, PUT, DELETE operations
  - Simple CRUD endpoints
  - No authentication needed
  - No rate limits
- **Example Endpoints**:
  - `GET /posts` - List all posts
  - `GET /posts/{id}` - Get a specific post
  - `GET /comments?postId={id}` - Get comments for a post
  - `GET /users` - List all users

### 2. **OpenWeather API** (with free tier)
- **URL**: https://api.openweathermap.org
- **OpenAPI Spec**: Available at https://openweathermap.org/api
- **Description**: Weather data API with free tier
- **Features**:
  - Current weather data
  - 5-day forecast
  - Geolocation data
  - Free tier available (requires API key)
- **Example Use Case**: Create tools for weather queries across your MCP

### 3. **REST Countries API**
- **URL**: https://restcountries.com
- **Description**: Provides information about countries
- **Features**:
  - No authentication required
  - Comprehensive country data
  - Multiple endpoints
  - CORS enabled
- **Example Endpoints**:
  - `GET /v3.1/all` - Get all countries
  - `GET /v3.1/name/{name}` - Search by country name
  - `GET /v3.1/region/{region}` - Get by region
  - `GET /v3.1/currency/{code}` - Get by currency

### 4. **PokéAPI**
- **URL**: https://pokeapi.co
- **Description**: RESTful Pokémon API with complete Pokémon data
- **Features**:
  - No authentication required
  - Fully documented
  - Well-structured REST API
  - Excellent test data
- **Example Endpoints**:
  - `GET /pokemon` - List Pokémon
  - `GET /pokemon/{id or name}` - Get specific Pokémon details
  - `GET /type` - List types
  - `GET /ability` - List abilities

### 5. **SpaceX API** (via r-spacex/SpaceX-API)
- **URL**: https://api.spacexdata.com/v4
- **Description**: Open-source SpaceX data API
- **Features**:
  - No authentication required
  - Launches, rockets, missions data
  - Well-documented REST API
- **Example Endpoints**:
  - `GET /launches` - List all launches
  - `GET /rockets` - List all rockets
  - `GET /missions` - List all missions

### 6. **Public APIs Directory**
- **URL**: https://github.com/public-apis/public-apis
- **Description**: Curated list of free public APIs
- **Features**:
  - Hundreds of free APIs
  - Categorized by type
  - Authentication requirements documented

## Self-Hosted / Local APIs for Testing

### 1. **Mock Server with JSON Server**
```bash
# Install globally
npm install -g json-server

# Create a db.json file with test data
echo '{
  "posts": [
    { "id": 1, "title": "First Post", "content": "Hello World" },
    { "id": 2, "title": "Second Post", "content": "More content" }
  ],
  "comments": [
    { "id": 1, "text": "Great post!", "postId": 1 }
  ]
}' > db.json

# Start the server
json-server --watch db.json --port 3000
```

### 2. **Local Node.js Express Server**
Create a simple test API in Go:
```bash
# Use the existing petstore example as a template
go run examples/mock-api/server.go
```

## Recommended Testing Strategy

### Phase 1: Existing Example APIs
1. Start with the included `petstore.yaml` OpenAPI specification
2. Test the blog GraphQL API with `blog.graphql`
3. Test async events with `user-events.yaml`

### Phase 2: Public APIs (No Auth)
1. **Test REST API Imports**: Use JSONPlaceholder or REST Countries API
   ```bash
   # Download OpenAPI spec
   curl -o countries-api.yaml https://raw.githubusercontent.com/fawazahmed0/open-weather-map-api/main/openapi.yaml
   
   # Import into AionMCP
   curl -X POST http://localhost:8080/api/v1/specs/import \
     -H "Content-Type: application/json" \
     -d '{"name": "countries", "url": "https://restcountries.com"}'
   ```

2. **Test Tool Execution**: Execute imported tools
   ```bash
   curl -X POST http://localhost:8080/api/v1/tools/list-countries/execute \
     -H "Content-Type: application/json" \
     -d '{}'
   ```

### Phase 3: Agent Integration
1. Test agent connecting to the AionMCP server
2. Verify tool discovery and execution
3. Test learning/reflection on agent interactions

## OpenAPI Spec Downloads

### Easy-to-use specs for testing:
1. **Stripe Public API** (Complex but well-documented)
   - https://raw.githubusercontent.com/stripe/openapi/master/openapi/spec3.yaml

2. **GitHub API** (Excellent documentation)
   - https://raw.githubusercontent.com/github/rest-api-description/main/descriptions/api.github.com.json

3. **Slack API** (Large but realistic)
   - https://raw.githubusercontent.com/slackapi/slack-api-schemas/main/openapi-3.yaml

## Testing with cURL Examples

### Import a JSONPlaceholder API Spec
```bash
# Create an OpenAPI spec for JSONPlaceholder
cat > /tmp/jsonplaceholder-openapi.yaml << 'EOF'
openapi: 3.0.0
info:
  title: JSONPlaceholder API
  version: 1.0.0
  description: Fake REST API for testing
servers:
  - url: https://jsonplaceholder.typicode.com
paths:
  /posts:
    get:
      summary: List all posts
      operationId: listPosts
      responses:
        '200':
          description: List of posts
  /posts/{id}:
    get:
      summary: Get a specific post
      operationId: getPost
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Post details
  /users:
    get:
      summary: List all users
      operationId: listUsers
      responses:
        '200':
          description: List of users
EOF

# Import the spec into AionMCP
curl -X POST http://localhost:8080/api/v1/specs/import \
  -H "Content-Type: application/yaml" \
  -d @/tmp/jsonplaceholder-openapi.yaml
```

### List Available Tools
```bash
curl -s http://localhost:8080/api/v1/tools | jq '.'
```

### Execute a Tool
```bash
curl -X POST http://localhost:8080/api/v1/tools/listPosts/execute \
  -H "Content-Type: application/json" \
  -d '{}'
```

## Recommended Learning Path

1. **Start Simple**: JSONPlaceholder (no auth, simple CRUD)
2. **Add Complexity**: REST Countries (filtering, search)
3. **Real-world Data**: PokéAPI (larger dataset, relationships)
4. **Advanced**: GitHub API (authentication, complex operations)
5. **Multiple Specs**: Mix OpenAPI, GraphQL, and AsyncAPI

## Monitoring & Learning

Once the server is running with imported APIs:
1. Check the learning engine output in `/tmp/aionmcp-logs`
2. View discovered tools via `/api/v1/tools`
3. Monitor execution patterns via learning dashboard
4. Verify reflection generation in `docs/reflections/`

## References

- **Public APIs Repository**: https://github.com/public-apis/public-apis
- **JSONPlaceholder**: https://jsonplaceholder.typicode.com
- **OpenAPI Initiative**: https://www.openapis.org
- **REST Countries**: https://restcountries.com
- **PokéAPI**: https://pokeapi.co
- **SpaceX API**: https://api.spacexdata.com
