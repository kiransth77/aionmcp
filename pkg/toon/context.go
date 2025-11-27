package toon

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Context represents a TOON (Token-Oriented Object Notation) context for LLM interactions
type Context struct {
	Version   string                 `json:"version"`
	Timestamp string                 `json:"timestamp"`
	Metadata  Metadata               `json:"metadata"`
	Tools     []ToolContext          `json:"tools,omitempty"`
	State     map[string]interface{} `json:"state,omitempty"`
	Messages  []Message              `json:"messages,omitempty"`
}

// Metadata contains contextual information about the interaction
type Metadata struct {
	Agent       string `json:"agent"`
	Source      string `json:"source"`
	SessionID   string `json:"session_id,omitempty"`
	User        string `json:"user,omitempty"`
	Model       string `json:"model,omitempty"`
	Environment string `json:"environment,omitempty"`
}

// ToolContext represents a tool in TOON format for LLM consumption
type ToolContext struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Type         string                 `json:"type"` // openapi, graphql, asyncapi, custom
	InputSchema  map[string]interface{} `json:"input_schema"`
	OutputSchema map[string]interface{} `json:"output_schema,omitempty"`
	Examples     []Example              `json:"examples,omitempty"`
	Tags         []string               `json:"tags,omitempty"`
	Confidence   float64                `json:"confidence,omitempty"` // 0.0-1.0
}

// Example shows sample usage of a tool
type Example struct {
	Input       map[string]interface{} `json:"input"`
	Output      map[string]interface{} `json:"output,omitempty"`
	Description string                 `json:"description,omitempty"`
}

// Message represents a conversation message with context
type Message struct {
	Role      string                 `json:"role"` // user, assistant, system, tool
	Content   string                 `json:"content"`
	Timestamp string                 `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	ToolID    string                 `json:"tool_id,omitempty"`
}

// Builder helps construct TOON contexts
type Builder struct {
	context Context
}

// NewBuilder creates a new TOON context builder
func NewBuilder(agent string, source string) *Builder {
	return &Builder{
		context: Context{
			Version:   "1.0",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Metadata: Metadata{
				Agent:  agent,
				Source: source,
			},
			Tools:    make([]ToolContext, 0),
			State:    make(map[string]interface{}),
			Messages: make([]Message, 0),
		},
	}
}

// WithModel sets the LLM model
func (b *Builder) WithModel(model string) *Builder {
	b.context.Metadata.Model = model
	return b
}

// WithSessionID sets the session identifier
func (b *Builder) WithSessionID(sessionID string) *Builder {
	b.context.Metadata.SessionID = sessionID
	return b
}

// WithUser sets the user identifier
func (b *Builder) WithUser(user string) *Builder {
	b.context.Metadata.User = user
	return b
}

// WithEnvironment sets the environment (dev, staging, prod)
func (b *Builder) WithEnvironment(env string) *Builder {
	b.context.Metadata.Environment = env
	return b
}

// AddTool adds a tool to the context
func (b *Builder) AddTool(tool ToolContext) *Builder {
	b.context.Tools = append(b.context.Tools, tool)
	return b
}

// AddState adds a state entry
func (b *Builder) AddState(key string, value interface{}) *Builder {
	b.context.State[key] = value
	return b
}

// AddMessage adds a conversation message
func (b *Builder) AddMessage(role string, content string) *Builder {
	b.context.Messages = append(b.context.Messages, Message{
		Role:      role,
		Content:   content,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
	return b
}

// AddToolMessage adds a tool-related message
func (b *Builder) AddToolMessage(role string, content string, toolID string) *Builder {
	b.context.Messages = append(b.context.Messages, Message{
		Role:      role,
		Content:   content,
		ToolID:    toolID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
	return b
}

// Build returns the constructed TOON context
func (b *Builder) Build() Context {
	return b.context
}

// ToJSON converts context to JSON string
func (b *Builder) ToJSON() (string, error) {
	data, err := json.MarshalIndent(b.context, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal TOON context: %w", err)
	}
	return string(data), nil
}

// ToCompactJSON converts context to compact JSON (minimal tokens)
func (b *Builder) ToCompactJSON() (string, error) {
	data, err := json.Marshal(b.context)
	if err != nil {
		return "", fmt.Errorf("failed to marshal TOON context: %w", err)
	}
	return string(data), nil
}

// ToString returns a human-readable string representation
func (c *Context) ToString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("TOON Context v%s\n", c.Version))
	sb.WriteString(fmt.Sprintf("Agent: %s | Source: %s\n", c.Metadata.Agent, c.Metadata.Source))
	if c.Metadata.Model != "" {
		sb.WriteString(fmt.Sprintf("Model: %s\n", c.Metadata.Model))
	}
	sb.WriteString(fmt.Sprintf("Tools: %d | Messages: %d | State: %d\n",
		len(c.Tools), len(c.Messages), len(c.State)))

	if len(c.Tools) > 0 {
		sb.WriteString("\nAvailable Tools:\n")
		for _, tool := range c.Tools {
			sb.WriteString(fmt.Sprintf("  - %s (%s): %s\n", tool.Name, tool.Type, tool.Description))
		}
	}
	return sb.String()
}

// Summary returns a token-efficient summary of the context
func (c *Context) Summary() string {
	return fmt.Sprintf(
		"{\"v\":%q,\"ag\":%q,\"src\":%q,\"t\":%d,\"m\":%d,\"s\":%d}",
		c.Version, c.Metadata.Agent, c.Metadata.Source,
		len(c.Tools), len(c.Messages), len(c.State),
	)
}

// TokenCount estimates the token count of the context
func (c *Context) TokenCount() int {
	// Rough estimation: 1 token ~= 4 characters
	json, err := json.Marshal(c)
	if err != nil {
		return 0
	}
	return len(json) / 4
}

// FormatForModel returns a formatted string optimized for model consumption
// Includes context summary and available tools
func (c *Context) FormatForModel() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[CONTEXT] Agent: %s | Session: %s\n", c.Metadata.Agent, c.Metadata.SessionID))

	sb.WriteString(fmt.Sprintf("[TOOLS] Available: %d\n", len(c.Tools)))
	for i, tool := range c.Tools {
		if i < 5 { // Show max 5 tools to save tokens
			sb.WriteString(fmt.Sprintf("  %d. %s - %s\n", i+1, tool.Name, tool.Description))
		}
	}
	if len(c.Tools) > 5 {
		sb.WriteString(fmt.Sprintf("  ... and %d more tools\n", len(c.Tools)-5))
	}

	if len(c.Messages) > 0 {
		sb.WriteString(fmt.Sprintf("[HISTORY] Recent messages: %d\n", len(c.Messages)))
		// Show last 3 messages
		start := len(c.Messages) - 3
		if start < 0 {
			start = 0
		}
		for _, msg := range c.Messages[start:] {
			content := msg.Content
			if len(content) > 100 {
				content = content[:100] + "..."
			}
			sb.WriteString(fmt.Sprintf("  [%s] %s\n", msg.Role, content))
		}
	}

	return sb.String()
}
