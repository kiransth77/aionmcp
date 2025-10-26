package core

import (
	"fmt"
	"sync"
	"time"

	"github.com/aionmcp/aionmcp/pkg/types"
	"go.uber.org/zap"
)

// Tool type alias for compatibility
type Tool = types.Tool

// ToolMetadata type alias for compatibility
type ToolMetadata = types.ToolMetadata

// ToolRegistry manages the collection of available tools
type ToolRegistry struct {
	mu     sync.RWMutex
	tools  map[string]Tool
	logger *zap.Logger
}

// NewToolRegistry creates a new tool registry
func NewToolRegistry(logger *zap.Logger) *ToolRegistry {
	registry := &ToolRegistry{
		tools:  make(map[string]Tool),
		logger: logger,
	}

	// Register built-in tools for iteration 0
	registry.registerBuiltinTools()
	
	return registry
}

// Register adds a tool to the registry
func (r *ToolRegistry) Register(tool Tool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := tool.Name()
	if name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}

	if _, exists := r.tools[name]; exists {
		r.logger.Warn("Tool already exists, replacing", zap.String("tool", name))
	}

	r.tools[name] = tool
	r.logger.Info("Tool registered", 
		zap.String("tool", name),
		zap.String("description", tool.Description()))

	return nil
}

// Unregister removes a tool from the registry
func (r *ToolRegistry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tools[name]; !exists {
		return fmt.Errorf("tool '%s' not found", name)
	}

	delete(r.tools, name)
	r.logger.Info("Tool unregistered", zap.String("tool", name))

	return nil
}

// Get retrieves a tool by name
func (r *ToolRegistry) Get(name string) (Tool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tool, exists := r.tools[name]
	if !exists {
		return nil, fmt.Errorf("tool '%s' not found", name)
	}

	return tool, nil
}

// ListTools returns metadata for all registered tools
func (r *ToolRegistry) ListTools() []ToolMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make([]ToolMetadata, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool.Metadata())
	}

	return tools
}

// Count returns the number of registered tools
func (r *ToolRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.tools)
}

// registerBuiltinTools adds some basic tools for iteration 0
func (r *ToolRegistry) registerBuiltinTools() {
	// Echo tool for testing
	echoTool := &EchoTool{}
	r.Register(echoTool)

	// Status tool
	statusTool := &StatusTool{registry: r}
	r.Register(statusTool)
}

// EchoTool - simple tool for testing MCP functionality
type EchoTool struct{}

func (t *EchoTool) Name() string {
	return "echo"
}

func (t *EchoTool) Description() string {
	return "Echoes back the input message for testing purposes"
}

func (t *EchoTool) Execute(input any) (any, error) {
	return map[string]any{
		"echo":      input,
		"timestamp": time.Now().Unix(),
		"tool":      t.Name(),
	}, nil
}

func (t *EchoTool) Metadata() ToolMetadata {
	return ToolMetadata{
		Name:        t.Name(),
		Description: t.Description(),
		Version:     "1.0.0",
		Source:      "builtin",
		Tags:        []string{"test", "utility"},
		Schema: map[string]any{
			"input": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"message": map[string]any{
						"type":        "string",
						"description": "Message to echo back",
					},
				},
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// StatusTool - provides information about the registry
type StatusTool struct {
	registry *ToolRegistry
}

func (t *StatusTool) Name() string {
	return "status"
}

func (t *StatusTool) Description() string {
	return "Returns information about the tool registry and server status"
}

func (t *StatusTool) Execute(input any) (any, error) {
	return map[string]any{
		"tool_count": t.registry.Count(),
		"timestamp":  time.Now().Unix(),
		"iteration":  "0",
		"status":     "active",
	}, nil
}

func (t *StatusTool) Metadata() ToolMetadata {
	return ToolMetadata{
		Name:        t.Name(),
		Description: t.Description(),
		Version:     "1.0.0",
		Source:      "builtin",
		Tags:        []string{"system", "status"},
		Schema: map[string]any{
			"input":  map[string]any{"type": "object"},
			"output": map[string]any{"type": "object"},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}