package core

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/aionmcp/aionmcp/pkg/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestTool implements the types.Tool interface for testing
type TestTool struct {
	name        string
	description string
	version     string
	source      string
}

func (t *TestTool) Name() string {
	return t.name
}

func (t *TestTool) Description() string {
	return t.description
}

func (t *TestTool) Execute(input any) (any, error) {
	return map[string]any{
		"tool":   t.name,
		"input":  input,
		"result": "success",
	}, nil
}

func (t *TestTool) Metadata() types.ToolMetadata {
	return types.ToolMetadata{
		Name:        t.name,
		Description: t.description,
		Version:     t.version,
		Source:      t.source,
		Tags:        []string{"test"},
		Schema:      map[string]any{"type": "object"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func TestToolRegistry_Register(t *testing.T) {
	logger := zap.NewNop()
	registry := NewToolRegistry(logger)

	tool := &TestTool{
		name:        "test-tool",
		description: "A test tool",
		version:     "1.0.0",
		source:      "test",
	}

	err := registry.Register(tool)
	assert.NoError(t, err)

	// Verify tool was registered
	retrievedTool, err := registry.Get("test-tool")
	assert.NoError(t, err)
	assert.Equal(t, "test-tool", retrievedTool.Name())
	assert.Equal(t, "A test tool", retrievedTool.Description())

	// Check count
	assert.Equal(t, 3, registry.Count()) // 2 builtin + 1 new
}

func TestToolRegistry_RegisterWithSource(t *testing.T) {
	logger := zap.NewNop()
	registry := NewToolRegistry(logger)

	tool := &TestTool{
		name:        "source-tool",
		description: "A tool with source",
		version:     "2.0.0",
		source:      "test-source",
	}

	err := registry.RegisterWithSource(tool, "custom-source", "2.1.0")
	assert.NoError(t, err)

	// Verify version and source tracking
	version, err := registry.GetVersion("source-tool")
	assert.NoError(t, err)
	assert.Equal(t, "2.1.0", version)

	source, err := registry.GetSource("source-tool")
	assert.NoError(t, err)
	assert.Equal(t, "custom-source", source)
}

func TestToolRegistry_RegisterBatch(t *testing.T) {
	logger := zap.NewNop()
	registry := NewToolRegistry(logger)

	tools := []types.Tool{
		&TestTool{name: "batch-tool-1", description: "Batch tool 1", version: "1.0.0"},
		&TestTool{name: "batch-tool-2", description: "Batch tool 2", version: "1.0.0"},
		&TestTool{name: "batch-tool-3", description: "Batch tool 3", version: "1.0.0"},
	}

	err := registry.RegisterBatch(tools, "batch-source")
	assert.NoError(t, err)

	// Verify all tools were registered
	for i, tool := range tools {
		retrievedTool, err := registry.Get(tool.Name())
		assert.NoError(t, err)
		assert.Equal(t, tools[i].Name(), retrievedTool.Name())

		source, err := registry.GetSource(tool.Name())
		assert.NoError(t, err)
		assert.Equal(t, "batch-source", source)
	}

	assert.Equal(t, 5, registry.Count()) // 2 builtin + 3 batch
}

func TestToolRegistry_UnregisterBySource(t *testing.T) {
	logger := zap.NewNop()
	registry := NewToolRegistry(logger)

	// Register tools from different sources
	tool1 := &TestTool{name: "source1-tool", description: "Tool from source 1"}
	tool2 := &TestTool{name: "source2-tool", description: "Tool from source 2"}
	tool3 := &TestTool{name: "source1-tool2", description: "Another tool from source 1"}

	registry.RegisterWithSource(tool1, "source1", "1.0.0")
	registry.RegisterWithSource(tool2, "source2", "1.0.0")
	registry.RegisterWithSource(tool3, "source1", "1.0.0")

	initialCount := registry.Count()
	assert.Equal(t, 5, initialCount) // 2 builtin + 3 new

	// Remove all tools from source1
	err := registry.UnregisterBySource("source1")
	assert.NoError(t, err)

	// Verify source1 tools are gone
	_, err = registry.Get("source1-tool")
	assert.Error(t, err)

	_, err = registry.Get("source1-tool2")
	assert.Error(t, err)

	// Verify source2 tool still exists
	_, err = registry.Get("source2-tool")
	assert.NoError(t, err)

	assert.Equal(t, 3, registry.Count()) // 2 builtin + 1 remaining
}

func TestToolRegistry_ListToolsBySource(t *testing.T) {
	logger := zap.NewNop()
	registry := NewToolRegistry(logger)

	// Register tools from different sources
	tool1 := &TestTool{name: "api-tool1", description: "API Tool 1"}
	tool2 := &TestTool{name: "api-tool2", description: "API Tool 2"}
	tool3 := &TestTool{name: "gql-tool1", description: "GraphQL Tool 1"}

	registry.RegisterWithSource(tool1, "openapi", "1.0.0")
	registry.RegisterWithSource(tool2, "openapi", "1.0.0")
	registry.RegisterWithSource(tool3, "graphql", "1.0.0")

	// List tools by source
	openAPITools := registry.ListToolsBySource("openapi")
	assert.Len(t, openAPITools, 2)

	toolNames := make([]string, len(openAPITools))
	for i, tool := range openAPITools {
		toolNames[i] = tool.Name
	}
	assert.Contains(t, toolNames, "api-tool1")
	assert.Contains(t, toolNames, "api-tool2")

	graphqlTools := registry.ListToolsBySource("graphql")
	assert.Len(t, graphqlTools, 1)
	assert.Equal(t, "gql-tool1", graphqlTools[0].Name)
}

func TestToolRegistry_GetToolSources(t *testing.T) {
	logger := zap.NewNop()
	registry := NewToolRegistry(logger)

	// Register tools from different sources
	tool1 := &TestTool{name: "tool1", description: "Tool 1"}
	tool2 := &TestTool{name: "tool2", description: "Tool 2"}
	tool3 := &TestTool{name: "tool3", description: "Tool 3"}

	registry.RegisterWithSource(tool1, "source-a", "1.0.0")
	registry.RegisterWithSource(tool2, "source-b", "1.0.0")
	registry.RegisterWithSource(tool3, "source-a", "1.0.0")

	sources := registry.GetToolSources()
	assert.Contains(t, sources, "builtin")  // From builtin tools
	assert.Contains(t, sources, "source-a")
	assert.Contains(t, sources, "source-b")
	assert.Len(t, sources, 3)
}

func TestToolRegistry_EventHandlers(t *testing.T) {
	logger := zap.NewNop()
	registry := NewToolRegistry(logger)

	var mu sync.Mutex
	var receivedEvents []ToolRegistryEvent
	handler := func(event ToolRegistryEvent) {
		mu.Lock()
		defer mu.Unlock()
		receivedEvents = append(receivedEvents, event)
	}

	handlerID := registry.AddEventHandler(handler)
	assert.Greater(t, handlerID, 0)

	// Register a tool
	tool := &TestTool{
		name:        "event-tool",
		description: "Tool for event testing",
	}

	registry.Register(tool)

	// Wait a bit for goroutine to execute
	time.Sleep(10 * time.Millisecond)

	// Verify event was received
	mu.Lock()
	assert.Len(t, receivedEvents, 1)
	assert.Equal(t, ToolEventAdded, receivedEvents[0].Type)
	assert.Equal(t, "event-tool", receivedEvents[0].ToolName)
	mu.Unlock()

	// Unregister the tool
	registry.Unregister("event-tool")

	// Wait a bit for goroutine to execute
	time.Sleep(10 * time.Millisecond)

	// Verify remove event was received
	mu.Lock()
	assert.Len(t, receivedEvents, 2)
	assert.Equal(t, ToolEventRemoved, receivedEvents[1].Type)
	assert.Equal(t, "event-tool", receivedEvents[1].ToolName)
	mu.Unlock()
	
	// Test handler removal
	removed := registry.RemoveEventHandler(handlerID)
	assert.True(t, removed)
	
	// Register another tool - handler should not receive event
	tool2 := &TestTool{
		name:        "event-tool-2",
		description: "Tool for testing removed handler",
	}
	registry.Register(tool2)
	
	// Wait a bit
	time.Sleep(10 * time.Millisecond)
	
	// Verify no new events were received (still 2)
	mu.Lock()
	assert.Len(t, receivedEvents, 2)
	mu.Unlock()
	
	// Test removing non-existent handler
	removed = registry.RemoveEventHandler(999)
	assert.False(t, removed)
}

func TestToolRegistry_GetRegistryStats(t *testing.T) {
	logger := zap.NewNop()
	registry := NewToolRegistry(logger)

	// Register tools from different sources
	tool1 := &TestTool{name: "stats-tool1", description: "Tool 1"}
	tool2 := &TestTool{name: "stats-tool2", description: "Tool 2"}

	registry.RegisterWithSource(tool1, "stats-source", "1.0.0")
	registry.RegisterWithSource(tool2, "stats-source", "1.0.0")

	stats := registry.GetRegistryStats()

	assert.Equal(t, 4, stats["total_tools"]) // 2 builtin + 2 new
	assert.Contains(t, stats["sources"], "builtin")
	assert.Contains(t, stats["sources"], "stats-source")

	toolsBySource := stats["tools_by_source"].(map[string]int)
	assert.Equal(t, 2, toolsBySource["builtin"])
	assert.Equal(t, 2, toolsBySource["stats-source"])
}

func TestToolRegistry_Replace(t *testing.T) {
	logger := zap.NewNop()
	registry := NewToolRegistry(logger)

	// Register initial tool
	tool1 := &TestTool{
		name:        "replaceable-tool",
		description: "Original tool",
		version:     "1.0.0",
	}

	err := registry.RegisterWithSource(tool1, "test-source", "1.0.0")
	assert.NoError(t, err)

	version, err := registry.GetVersion("replaceable-tool")
	assert.NoError(t, err)
	assert.Equal(t, "1.0.0", version)

	// Replace with newer version
	tool2 := &TestTool{
		name:        "replaceable-tool",
		description: "Updated tool",
		version:     "2.0.0",
	}

	err = registry.RegisterWithSource(tool2, "test-source", "2.0.0")
	assert.NoError(t, err)

	// Verify tool was updated
	updatedTool, err := registry.Get("replaceable-tool")
	assert.NoError(t, err)
	assert.Equal(t, "Updated tool", updatedTool.Description())

	version, err = registry.GetVersion("replaceable-tool")
	assert.NoError(t, err)
	assert.Equal(t, "2.0.0", version)

	// Count should remain the same
	assert.Equal(t, 3, registry.Count()) // 2 builtin + 1 updated
}

func TestToolRegistry_GetNonexistent(t *testing.T) {
	logger := zap.NewNop()
	registry := NewToolRegistry(logger)

	_, err := registry.Get("nonexistent-tool")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	_, err = registry.GetVersion("nonexistent-tool")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	_, err = registry.GetSource("nonexistent-tool")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestToolRegistry_EmptyName(t *testing.T) {
	logger := zap.NewNop()
	registry := NewToolRegistry(logger)

	tool := &TestTool{
		name:        "", // Empty name
		description: "Tool with empty name",
	}

	err := registry.Register(tool)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name cannot be empty")

	err = registry.RegisterWithSource(tool, "test", "1.0.0")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name cannot be empty")
}

// Benchmark tests
func BenchmarkToolRegistry_Register(b *testing.B) {
	logger := zap.NewNop()
	registry := NewToolRegistry(logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tool := &TestTool{
			name:        fmt.Sprintf("bench-tool-%d", i),
			description: "Benchmark tool",
		}
		registry.Register(tool)
	}
}

func BenchmarkToolRegistry_Get(b *testing.B) {
	logger := zap.NewNop()
	registry := NewToolRegistry(logger)

	// Pre-register tools
	for i := 0; i < 1000; i++ {
		tool := &TestTool{
			name:        fmt.Sprintf("bench-tool-%d", i),
			description: "Benchmark tool",
		}
		registry.Register(tool)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		registry.Get(fmt.Sprintf("bench-tool-%d", i%1000))
	}
}

func BenchmarkToolRegistry_ListTools(b *testing.B) {
	logger := zap.NewNop()
	registry := NewToolRegistry(logger)

	// Pre-register tools
	for i := 0; i < 1000; i++ {
		tool := &TestTool{
			name:        fmt.Sprintf("bench-tool-%d", i),
			description: "Benchmark tool",
		}
		registry.Register(tool)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		registry.ListTools()
	}
}