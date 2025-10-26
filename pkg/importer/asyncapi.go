package importer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aionmcp/aionmcp/pkg/types"
)

// AsyncAPIImporter handles AsyncAPI specifications
type AsyncAPIImporter struct{}

// NewAsyncAPIImporter creates a new AsyncAPI importer
func NewAsyncAPIImporter() *AsyncAPIImporter {
	return &AsyncAPIImporter{}
}

// GetType returns the specification type
func (i *AsyncAPIImporter) GetType() SpecType {
	return SpecTypeAsyncAPI
}

// Supports checks if this importer can handle the given source
func (i *AsyncAPIImporter) Supports(source SpecSource) bool {
	return source.Type == SpecTypeAsyncAPI
}

// Validate checks if the AsyncAPI specification is valid
func (i *AsyncAPIImporter) Validate(ctx context.Context, source SpecSource) error {
	content, err := i.loadSpec(source.Path)
	if err != nil {
		return err
	}

	// Simple validation - check if it's valid JSON/YAML
	var spec map[string]interface{}
	if err := json.Unmarshal(content, &spec); err != nil {
		// Try YAML parsing as fallback
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	// Check for required AsyncAPI fields
	if _, exists := spec["asyncapi"]; !exists {
		return fmt.Errorf("missing required 'asyncapi' field")
	}

	return nil
}

// Import parses the AsyncAPI specification and generates tools
func (i *AsyncAPIImporter) Import(ctx context.Context, source SpecSource) (*ImportResult, error) {
	start := time.Now()
	
	result := &ImportResult{
		Source:    source,
		Tools:     []types.Tool{},
		Errors:    []error{},
		Warnings:  []string{},
		Timestamp: start,
	}

	// Load the specification
	content, err := i.loadSpec(source.Path)
	if err != nil {
		result.Errors = append(result.Errors, err)
		result.Duration = time.Since(start)
		return result, err
	}

	// Parse the AsyncAPI document as JSON
	var spec map[string]interface{}
	if err := json.Unmarshal(content, &spec); err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("failed to parse AsyncAPI spec: %w", err))
		result.Duration = time.Since(start)
		return result, err
	}

	// Extract channels
	channels, ok := spec["channels"].(map[string]interface{})
	if !ok {
		result.Warnings = append(result.Warnings, "No channels found in AsyncAPI specification")
		result.Duration = time.Since(start)
		return result, nil
	}

	// Generate tools from channels
	for channelName, channelData := range channels {
		channel, ok := channelData.(map[string]interface{})
		if !ok {
			continue
		}

		// Create publish tools
		if publish, exists := channel["publish"]; exists {
			tool := i.createPublishTool(source, spec, channelName, channel, publish)
			result.Tools = append(result.Tools, tool)
		}

		// Create subscribe tools
		if subscribe, exists := channel["subscribe"]; exists {
			tool := i.createSubscribeTool(source, spec, channelName, channel, subscribe)
			result.Tools = append(result.Tools, tool)
		}
	}

	// Add warning if no servers are defined
	if servers, exists := spec["servers"]; !exists || len(servers.(map[string]interface{})) == 0 {
		result.Warnings = append(result.Warnings, "No servers defined in AsyncAPI spec, tools may need manual configuration")
	}

	result.Duration = time.Since(start)
	return result, nil
}

// loadSpec loads an AsyncAPI specification from file
func (i *AsyncAPIImporter) loadSpec(path string) ([]byte, error) {
	// For now, only support file loading
	// TODO: Add URL support for AsyncAPI specs
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return nil, fmt.Errorf("URL loading not yet supported for AsyncAPI specs")
	}

	return os.ReadFile(path)
}

// createPublishTool creates a tool for publishing messages to a channel
func (i *AsyncAPIImporter) createPublishTool(source SpecSource, spec map[string]interface{}, channelName string, channel map[string]interface{}, publish interface{}) types.Tool {
	return &AsyncAPITool{
		source:      source,
		spec:        spec,
		channelName: channelName,
		channel:     channel,
		operation:   "publish",
	}
}

// createSubscribeTool creates a tool for subscribing to messages from a channel
func (i *AsyncAPIImporter) createSubscribeTool(source SpecSource, spec map[string]interface{}, channelName string, channel map[string]interface{}, subscribe interface{}) types.Tool {
	return &AsyncAPITool{
		source:      source,
		spec:        spec,
		channelName: channelName,
		channel:     channel,
		operation:   "subscribe",
	}
}

// AsyncAPITool represents a tool generated from an AsyncAPI operation
type AsyncAPITool struct {
	source      SpecSource
	spec        map[string]interface{}
	channelName string
	channel     map[string]interface{}
	operation   string // "publish" or "subscribe"
}

// Name returns the tool name
func (t *AsyncAPITool) Name() string {
	// Clean channel name for use in tool name
	cleanChannel := strings.ReplaceAll(t.channelName, "/", "_")
	cleanChannel = strings.ReplaceAll(cleanChannel, "{", "")
	cleanChannel = strings.ReplaceAll(cleanChannel, "}", "")
	return fmt.Sprintf("asyncapi.%s.%s_%s", t.source.ID, t.operation, cleanChannel)
}

// Description returns the tool description
func (t *AsyncAPITool) Description() string {
	switch t.operation {
	case "publish":
		return fmt.Sprintf("Publish message to %s channel", t.channelName)
	case "subscribe":
		return fmt.Sprintf("Subscribe to messages from %s channel", t.channelName)
	}

	return fmt.Sprintf("AsyncAPI %s operation on channel %s", t.operation, t.channelName)
}

// Execute performs the AsyncAPI operation
func (t *AsyncAPITool) Execute(input any) (any, error) {
	// Parse input
	inputMap, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("input must be a JSON object")
	}

	// Get server information
	servers, exists := t.spec["servers"].(map[string]interface{})
	if !exists || len(servers) == 0 {
		return nil, fmt.Errorf("no servers defined in AsyncAPI specification")
	}

	// Use first server (in production, this should be configurable)
	var serverURL string
	var protocol string
	for _, serverData := range servers {
		if server, ok := serverData.(map[string]interface{}); ok {
			if url, exists := server["url"].(string); exists {
				serverURL = url
			}
			if prot, exists := server["protocol"].(string); exists {
				protocol = prot
			}
			break
		}
	}

	switch t.operation {
	case "publish":
		return t.executePublish(inputMap, serverURL, protocol)
	case "subscribe":
		return t.executeSubscribe(inputMap, serverURL, protocol)
	default:
		return nil, fmt.Errorf("unsupported operation: %s", t.operation)
	}
}

// executePublish handles message publishing
func (t *AsyncAPITool) executePublish(input map[string]interface{}, serverURL, protocol string) (interface{}, error) {
	// Extract message payload
	payload, exists := input["payload"]
	if !exists {
		return nil, fmt.Errorf("payload is required for publish operation")
	}

	// For now, return a simulation response
	// TODO: Implement actual message publishing based on protocol (MQTT, AMQP, WebSocket, etc.)
	result := map[string]interface{}{
		"operation":    "publish",
		"channel":      t.channelName,
		"payload":      payload,
		"server_url":   serverURL,
		"protocol":     protocol,
		"timestamp":    time.Now().Unix(),
		"status":       "simulated", // Indicates this is a simulation
		"message":      "Message publishing is simulated. Actual implementation depends on protocol adapter.",
	}

	// Add any headers from input
	if headers, exists := input["headers"]; exists {
		result["headers"] = headers
	}

	return result, nil
}

// executeSubscribe handles message subscription
func (t *AsyncAPITool) executeSubscribe(input map[string]interface{}, serverURL, protocol string) (interface{}, error) {
	// Extract subscription parameters
	timeout := 30 // Default timeout in seconds
	if timeoutVal, exists := input["timeout"]; exists {
		if timeoutFloat, ok := timeoutVal.(float64); ok {
			timeout = int(timeoutFloat)
		}
	}

	// For now, return a simulation response
	// TODO: Implement actual message subscription based on protocol
	result := map[string]interface{}{
		"operation":    "subscribe",
		"channel":      t.channelName,
		"server_url":   serverURL,
		"protocol":     protocol,
		"timeout":      timeout,
		"timestamp":    time.Now().Unix(),
		"status":       "simulated",
		"message":      "Message subscription is simulated. Actual implementation depends on protocol adapter.",
		"simulated_messages": []map[string]interface{}{
			{
				"payload":   map[string]interface{}{"message": "Simulated message 1"},
				"timestamp": time.Now().Unix(),
			},
		},
	}

	return result, nil
}

// Metadata returns tool metadata
func (t *AsyncAPITool) Metadata() types.ToolMetadata {
	// Build input schema based on operation type
	inputSchema := map[string]interface{}{
		"type":       "object",
		"properties": make(map[string]interface{}),
		"required":   []string{},
	}

	properties := inputSchema["properties"].(map[string]interface{})
	var required []string

	switch t.operation {
	case "publish":
		// Publish operations require a payload
		properties["payload"] = map[string]interface{}{
			"type":        "object",
			"description": "Message payload to publish",
		}
		properties["headers"] = map[string]interface{}{
			"type":        "object",
			"description": "Optional message headers",
		}
		required = append(required, "payload")

	case "subscribe":
		// Subscribe operations can have timeout and filter parameters
		properties["timeout"] = map[string]interface{}{
			"type":        "integer",
			"description": "Subscription timeout in seconds",
			"default":     30,
		}
		properties["filter"] = map[string]interface{}{
			"type":        "object",
			"description": "Optional message filter criteria",
		}
	}

	inputSchema["required"] = required

	// Build output schema
	outputSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"operation":    map[string]interface{}{"type": "string"},
			"channel":      map[string]interface{}{"type": "string"},
			"server_url":   map[string]interface{}{"type": "string"},
			"protocol":     map[string]interface{}{"type": "string"},
			"timestamp":    map[string]interface{}{"type": "integer"},
			"status":       map[string]interface{}{"type": "string"},
		},
	}

	// Add operation-specific output properties
	switch t.operation {
	case "publish":
		outputSchema["properties"].(map[string]interface{})["payload"] = map[string]interface{}{"type": "object"}
	case "subscribe":
		outputSchema["properties"].(map[string]interface{})["messages"] = map[string]interface{}{
			"type": "array",
			"items": map[string]interface{}{"type": "object"},
		}
	}

	return types.ToolMetadata{
		Name:        t.Name(),
		Description: t.Description(),
		Version:     "1.0.0",
		Source:      string(SpecTypeAsyncAPI),
		Tags:        []string{"asyncapi", "messaging", t.operation},
		Schema: map[string]interface{}{
			"input":  inputSchema,
			"output": outputSchema,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}