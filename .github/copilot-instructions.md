# AionMCP - Autonomous Go MCP Server

## Project Overview
AionMCP is an autonomous Go-based Model Context Protocol (MCP) server that dynamically imports OpenAPI, GraphQL, and AsyncAPI specifications and exposes them as tools to agents. Features self-learning, context-awareness, and auto-documentation using Clean/Hexagonal architecture.

## Project Structure
```
cmd/                     # Application entry points
internal/core/           # Core business logic
internal/adapters/       # External integrations
internal/registry/       # Tool registration system
internal/selflearn/      # Learning and reflection engine
pkg/agent/              # Agent integration APIs
docs/                   # Auto-generated documentation
data/feedback/          # Learning data storage
```

## Core Libraries & Dependencies
- **API Framework**: `github.com/gin-gonic/gin`
- **gRPC**: `google.golang.org/grpc`
- **Plugin System**: `github.com/hashicorp/go-plugin`
- **Configuration**: `github.com/spf13/viper`
- **Logging**: `go.uber.org/zap`
- **Storage**: `github.com/boltdb/bolt`, `github.com/dgraph-io/badger`
- **Task Queue**: `github.com/hibiken/asynq`
- **Metrics**: `github.com/prometheus/client_golang`

## MCP Protocol Implementation
Core interfaces follow this pattern:
```go
type Tool interface {
    Name() string
    Description() string
    Execute(input any) (any, error)
}

type ToolRegistry struct {
    mu    sync.RWMutex
    tools map[string]Tool
}
```

MCP envelope format:
```json
{
  "protocol": "MCP/1.0",
  "type": "tool_call",
  "tool": "openapi.petstore.listPets",
  "args": { "limit": 10 },
  "context": { "session": "abc123" }
}
```

## Development Iterations

### Iteration 0: Project Scaffolding
**Goal**: Base MCP server structure
**Outputs**: `cmd/server/main.go`, `internal/core/server.go`, registry stub, `docs/README.md`, `docs/architecture.md`

### Iteration 1: Spec Importer Layer  
**Goal**: Multi-format API spec parsing
**Outputs**: `pkg/importer/{openapi,graphql,asyncapi}.go`, `docs/spec_importer.md`
**Constraints**: Go interfaces, hot-reloadable, cached tools
```go
func ImportOpenAPI(path string) ([]Tool, error) {
    spec, err := openapi3.NewLoader().LoadFromFile(path)
    if err != nil { return nil, err }
    var tools []Tool
    for name, op := range spec.Paths {
        tools = append(tools, NewOpenAPITool(name, op))
    }
    return tools, nil
}
```

### Iteration 2: Self-Learning Engine
**Goal**: Execution feedback and reflection
**Outputs**: `pkg/feedback/engine.go`, `docs/self_learning.md`
**Features**: Log results, record errors/retries, generate `reflection.md`
```go
func RecordReflection(tool, errMsg, suggestion string) {
    ref := Reflection{
        Timestamp: time.Now(),
        ToolName: tool,
        Error: errMsg,
        Suggestion: suggestion,
    }
    SaveToBoltDB(ref)
}
```

### Iteration 3: Autonomous Documentation
**Goal**: Auto-generated docs and changelogs
**Outputs**: `docs/changelog.md`, `docs/reflections/YYYY-MM-DD.md`, auto-updated `README.md`

### Iteration 4: Agent Integration APIs
**Goal**: REST + gRPC interfaces for agents
**Outputs**: `pkg/agent/{api,server}.go`, `docs/agent_integration.md`
**Requirements**: List/invoke tools dynamically, hot-plug new tools

### Iteration 5: Competitive Analysis
**Goal**: Research and differentiation documentation
**Outputs**: `docs/comparison.md`
**Targets**: openapi-mcp-server, api-to-mcp, openapi-mcp-generator
**Differentiators**: Multi-spec support, Go-native, self-learning, task orchestration

### Iteration 6: Task Orchestration
**Goal**: Multi-step adaptive execution
**Outputs**: `internal/core/orchestrator.go`, `docs/iterations.md`
**Features**: Task chaining, adaptive behavior, registry confidence scores

## Coding Conventions

### Error Handling
- Use structured errors with context: `fmt.Errorf("failed to parse OpenAPI spec %s: %w", specPath, err)`
- Implement retry logic with exponential backoff for external API calls
- Log all errors with sufficient context for learning system analysis

### Documentation Standards
- Document all architecture decisions in `docs/adr/` (Architecture Decision Records)
- Maintain iteration logs in `docs/iterations/` with what/why/learnings format
- Use Go doc comments for all public APIs
- Include usage examples in package documentation

### Testing Strategy
- Unit tests for all parsers and core logic
- Integration tests with sample API specifications
- End-to-end tests with real MCP client interactions
- Benchmark tests for performance-critical paths (spec parsing, tool execution)

## Learning System Design
- Store execution metadata in structured format (JSON/SQLite)
- Analyze failure patterns: API response codes, timeout frequencies, parameter validation errors
- Implement suggestion engine for common failure remediation
- Track tool usage patterns for optimization opportunities

## Competitive Analysis Targets
Research and differentiate from:
- `openapi-mcp-server`: Direct OpenAPI to MCP bridge
- `api-to-mcp`: Generic API conversion utility  
- `openapi-mcp-generator`: Static code generation approach

Document unique value propositions:
- Dynamic runtime adaptation vs static generation
- Multi-protocol support (OpenAPI + GraphQL + AsyncAPI)
- Autonomous learning and error recovery
- Embeddable architecture for various agent frameworks

## Iteration Strategy
1. **MVP**: Basic OpenAPI spec import with simple MCP tool generation
2. **Enhancement**: Add GraphQL and AsyncAPI support, hot-reload capability
3. **Intelligence**: Implement learning layer and failure recovery
4. **Modularity**: Create embeddable adapters and standalone deployment options
5. **Optimization**: Performance tuning, competitive analysis, feature gaps

Always explain implementation decisions, document learnings, and maintain clear separation between MVP and advanced features.