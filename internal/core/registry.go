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

// ToolRegistryEvent represents events in the tool registry
type ToolRegistryEvent struct {
	Type      ToolEventType `json:"type"`
	ToolName  string        `json:"tool_name"`
	Metadata  ToolMetadata  `json:"metadata"`
	Timestamp time.Time     `json:"timestamp"`
}

// ToolEventType represents the type of tool registry event
type ToolEventType string

const (
	ToolEventAdded   ToolEventType = "tool_added"
	ToolEventRemoved ToolEventType = "tool_removed"
	ToolEventUpdated ToolEventType = "tool_updated"
)

// ToolRegistryEventHandler handles tool registry events
type ToolRegistryEventHandler func(event ToolRegistryEvent)

// eventHandlerEntry wraps a handler with its unique ID
type eventHandlerEntry struct {
	id      int
	handler ToolRegistryEventHandler
}

// ToolRegistry manages the collection of available tools with dynamic registration
// It implements the types.ToolRegistry interface
type ToolRegistry struct {
	mu               sync.RWMutex
	tools            map[string]Tool
	versions         map[string]string // tool name -> version
	sources          map[string]string // tool name -> source identifier
	eventHandlers    []eventHandlerEntry
	nextHandlerID    int
	logger           *zap.Logger
	handlerSemaphore chan struct{} // Limits concurrent event handler executions
}

// NewToolRegistry creates a new tool registry with dynamic capabilities
func NewToolRegistry(logger *zap.Logger) *ToolRegistry {
	registry := &ToolRegistry{
		tools:            make(map[string]Tool),
		versions:         make(map[string]string),
		sources:          make(map[string]string),
		eventHandlers:    make([]eventHandlerEntry, 0),
		nextHandlerID:    1,
		logger:           logger,
		handlerSemaphore: make(chan struct{}, 50), // Limit to 50 concurrent handlers
	}

	// Register built-in tools for iteration 0
	registry.registerBuiltinTools()
	
	return registry
}

// Register adds a tool to the registry with version and source tracking
func (r *ToolRegistry) Register(tool Tool) error {
	return r.RegisterWithSource(tool, "unknown", "")
}

// RegisterWithSource adds a tool to the registry with source information
func (r *ToolRegistry) RegisterWithSource(tool Tool, sourceID, version string) error {
	r.mu.Lock()

	name := tool.Name()
	if name == "" {
		r.mu.Unlock()
		return fmt.Errorf("tool name cannot be empty")
	}

	eventType := ToolEventAdded
	if _, exists := r.tools[name]; exists {
		eventType = ToolEventUpdated
		r.logger.Warn("Tool already exists, updating", 
			zap.String("tool", name),
			zap.String("old_version", r.versions[name]),
			zap.String("new_version", version))
	}

	r.tools[name] = tool
	r.versions[name] = version
	r.sources[name] = sourceID

	r.logger.Info("Tool registered", 
		zap.String("tool", name),
		zap.String("description", tool.Description()),
		zap.String("version", version),
		zap.String("source", sourceID))

	// Prepare event while still holding lock
	event := ToolRegistryEvent{
		Type:      eventType,
		ToolName:  name,
		Metadata:  tool.Metadata(),
		Timestamp: time.Now(),
	}
	r.mu.Unlock()
	
	// Emit event after releasing lock to avoid deadlock
	r.emitEvent(event)

	return nil
}

// RegisterBatch adds multiple tools atomically
func (r *ToolRegistry) RegisterBatch(tools []Tool, sourceID string) error {
	r.mu.Lock()

	// Validate all tools first
	for _, tool := range tools {
		if tool.Name() == "" {
			r.mu.Unlock()
			return fmt.Errorf("tool name cannot be empty")
		}
	}

	// Register all tools
	events := make([]ToolRegistryEvent, 0, len(tools))
	for _, tool := range tools {
		name := tool.Name()
		metadata := tool.Metadata()
		
		eventType := ToolEventAdded
		if _, exists := r.tools[name]; exists {
			eventType = ToolEventUpdated
		}

		r.tools[name] = tool
		r.versions[name] = metadata.Version
		r.sources[name] = sourceID

		events = append(events, ToolRegistryEvent{
			Type:      eventType,
			ToolName:  name,
			Metadata:  metadata,
			Timestamp: time.Now(),
		})
	}

	r.logger.Info("Batch tool registration completed",
		zap.Int("count", len(tools)),
		zap.String("source", sourceID))

	r.mu.Unlock()

	// Emit all events after releasing lock to avoid deadlock
	for _, event := range events {
		r.emitEvent(event)
	}

	return nil
}

// UnregisterBySource removes all tools from a specific source
func (r *ToolRegistry) UnregisterBySource(sourceID string) error {
	r.mu.Lock()

	var removedTools []string
	for name, source := range r.sources {
		if source == sourceID {
			removedTools = append(removedTools, name)
		}
	}

	var events []ToolRegistryEvent
	for _, name := range removedTools {
		tool := r.tools[name]
		delete(r.tools, name)
		delete(r.versions, name)
		delete(r.sources, name)

		r.logger.Info("Tool unregistered by source", 
			zap.String("tool", name),
			zap.String("source", sourceID))

		// Prepare event
		events = append(events, ToolRegistryEvent{
			Type:      ToolEventRemoved,
			ToolName:  name,
			Metadata:  tool.Metadata(),
			Timestamp: time.Now(),
		})
	}

	r.logger.Info("Batch tool removal by source completed",
		zap.Int("count", len(removedTools)),
		zap.String("source", sourceID))

	r.mu.Unlock()

	// Emit events after releasing lock to avoid deadlock
	for _, event := range events {
		r.emitEvent(event)
	}

	return nil
}

// Unregister removes a tool from the registry
func (r *ToolRegistry) Unregister(name string) error {
	r.mu.Lock()

	tool, exists := r.tools[name]
	if !exists {
		r.mu.Unlock()
		return fmt.Errorf("tool '%s' not found", name)
	}

	delete(r.tools, name)
	delete(r.versions, name)
	delete(r.sources, name)
	
	r.logger.Info("Tool unregistered", zap.String("tool", name))

	// Prepare event
	event := ToolRegistryEvent{
		Type:      ToolEventRemoved,
		ToolName:  name,
		Metadata:  tool.Metadata(),
		Timestamp: time.Now(),
	}
	r.mu.Unlock()

	// Emit event after releasing lock to avoid deadlock
	r.emitEvent(event)

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

// GetVersion returns the version of a specific tool
func (r *ToolRegistry) GetVersion(name string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	version, exists := r.versions[name]
	if !exists {
		return "", fmt.Errorf("tool '%s' not found", name)
	}
	return version, nil
}

// GetSource returns the source of a specific tool
func (r *ToolRegistry) GetSource(name string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	source, exists := r.sources[name]
	if !exists {
		return "", fmt.Errorf("tool '%s' not found", name)
	}
	return source, nil
}

// ListToolsBySource returns tools from a specific source
func (r *ToolRegistry) ListToolsBySource(sourceID string) []ToolMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var tools []ToolMetadata
	for name, source := range r.sources {
		if source == sourceID {
			if tool, exists := r.tools[name]; exists {
				tools = append(tools, tool.Metadata())
			}
		}
	}
	return tools
}

// GetToolSources returns all unique source identifiers
func (r *ToolRegistry) GetToolSources() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sourceSet := make(map[string]bool)
	for _, source := range r.sources {
		sourceSet[source] = true
	}

	sources := make([]string, 0, len(sourceSet))
	for source := range sourceSet {
		sources = append(sources, source)
	}
	return sources
}

// AddEventHandler adds an event handler for tool registry changes and returns a handler ID
func (r *ToolRegistry) AddEventHandler(handler ToolRegistryEventHandler) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	handlerID := r.nextHandlerID
	r.nextHandlerID++
	
	r.eventHandlers = append(r.eventHandlers, eventHandlerEntry{
		id:      handlerID,
		handler: handler,
	})
	
	return handlerID
}

// RemoveEventHandler removes an event handler by its ID
func (r *ToolRegistry) RemoveEventHandler(handlerID int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	for i, entry := range r.eventHandlers {
		if entry.id == handlerID {
			// Remove by replacing with last element and truncating
			r.eventHandlers[i] = r.eventHandlers[len(r.eventHandlers)-1]
			r.eventHandlers = r.eventHandlers[:len(r.eventHandlers)-1]
			return true
		}
	}
	
	r.logger.Warn("Attempted to remove non-existent event handler", zap.Int("handler_id", handlerID))
	return false
}

// emitEvent sends an event to all registered handlers with bounded concurrency
func (r *ToolRegistry) emitEvent(event ToolRegistryEvent) {
	// Don't hold the lock while calling handlers to avoid deadlocks
	r.mu.RLock()
	handlers := make([]eventHandlerEntry, len(r.eventHandlers))
	copy(handlers, r.eventHandlers)
	r.mu.RUnlock()

	for _, entry := range handlers {
		go func(h ToolRegistryEventHandler, registry *ToolRegistry) {
			// Acquire semaphore slot (blocks if at capacity)
			registry.handlerSemaphore <- struct{}{}
			defer func() {
				// Release semaphore slot
				<-registry.handlerSemaphore
				
				if recovered := recover(); recovered != nil {
					registry.logger.Error("Tool registry event handler panic", 
						zap.String("event_type", string(event.Type)),
						zap.String("tool_name", event.ToolName),
						zap.Any("panic", recovered))
				}
			}()
			h(event)
		}(entry.handler, r)
	}
}

// GetRegistryStats returns statistics about the registry
func (r *ToolRegistry) GetRegistryStats() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sourceStats := make(map[string]int)
	for _, source := range r.sources {
		sourceStats[source]++
	}

	// Inline sources calculation to avoid nested lock acquisition
	sourceSet := make(map[string]bool)
	for _, source := range r.sources {
		sourceSet[source] = true
	}

	sources := make([]string, 0, len(sourceSet))
	for source := range sourceSet {
		sources = append(sources, source)
	}

	return map[string]interface{}{
		"total_tools":     len(r.tools),
		"sources":         sources,
		"tools_by_source": sourceStats,
		"event_handlers":  len(r.eventHandlers),
	}
}

// registerBuiltinTools adds some basic tools for iteration 0
func (r *ToolRegistry) registerBuiltinTools() {
	// Echo tool for testing
	echoTool := &EchoTool{}
	r.RegisterWithSource(echoTool, "builtin", "1.0.0")

	// Status tool
	statusTool := &StatusTool{registry: r}
	r.RegisterWithSource(statusTool, "builtin", "1.0.0")
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