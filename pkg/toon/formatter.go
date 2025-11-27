package toon

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ToolFormatter converts tools to TOON format for model consumption
type ToolFormatter struct {
	compact bool // If true, minimize tokens
}

// NewToolFormatter creates a new tool formatter
func NewToolFormatter(compact bool) *ToolFormatter {
	return &ToolFormatter{
		compact: compact,
	}
}

// FormatToolDescription creates a TOON-formatted tool description
// Optimized for token efficiency while maintaining clarity
func (tf *ToolFormatter) FormatToolDescription(
	toolID string,
	name string,
	description string,
	inputSchema map[string]interface{},
	tags []string,
) string {
	if tf.compact {
		return tf.formatCompact(toolID, name, description, tags)
	}
	return tf.formatVerbose(toolID, name, description, inputSchema, tags)
}

// formatCompact returns minimal token representation
func (tf *ToolFormatter) formatCompact(
	toolID string,
	name string,
	description string,
	tags []string,
) string {
	return fmt.Sprintf(
		"[%s] %s: %s %s",
		toolID, name, description, formatTags(tags),
	)
}

// formatVerbose returns detailed token representation
func (tf *ToolFormatter) formatVerbose(
	toolID string,
	name string,
	description string,
	inputSchema map[string]interface{},
	tags []string,
) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Tool: %s\n", name))
	sb.WriteString(fmt.Sprintf("ID: %s\n", toolID))
	sb.WriteString(fmt.Sprintf("Description: %s\n", description))

	if len(tags) > 0 {
		sb.WriteString(fmt.Sprintf("Tags: %s\n", strings.Join(tags, ", ")))
	}

	if len(inputSchema) > 0 {
		sb.WriteString("Parameters:\n")
		sb.WriteString(formatSchema(inputSchema, "  "))
	}

	return sb.String()
}

// formatTags returns a formatted tag string
func formatTags(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	return fmt.Sprintf("[%s]", strings.Join(tags, ","))
}

// formatSchema returns a formatted schema string
func formatSchema(schema map[string]interface{}, indent string) string {
	var sb strings.Builder

	if properties, ok := schema["properties"].(map[string]interface{}); ok {
		for key, prop := range properties {
			propMap, _ := prop.(map[string]interface{})
			propType := getSchemaType(propMap)
			required := "optional"
			if reqs, ok := schema["required"].([]interface{}); ok {
				for _, req := range reqs {
					if req == key {
						required = "required"
						break
					}
				}
			}
			sb.WriteString(fmt.Sprintf("%s%s (%s, %s)\n", indent, key, propType, required))
		}
	}

	return sb.String()
}

// getSchemaType extracts type from schema map
func getSchemaType(schema map[string]interface{}) string {
	if t, ok := schema["type"]; ok {
		return fmt.Sprintf("%v", t)
	}
	return "any"
}

// ConversationContext represents a model conversation with TOON context
type ConversationContext struct {
	SystemPrompt string
	Tools        []ToolContext
	History      []Message
	State        map[string]interface{}
}

// NewConversationContext creates a new conversation context
func NewConversationContext(systemPrompt string) *ConversationContext {
	return &ConversationContext{
		SystemPrompt: systemPrompt,
		Tools:        make([]ToolContext, 0),
		History:      make([]Message, 0),
		State:        make(map[string]interface{}),
	}
}

// AddTool adds a tool to the conversation
func (cc *ConversationContext) AddTool(tool ToolContext) {
	cc.Tools = append(cc.Tools, tool)
}

// AddMessage adds a message to the conversation history
func (cc *ConversationContext) AddMessage(role string, content string) {
	cc.History = append(cc.History, Message{
		Role:      role,
		Content:   content,
		Timestamp: getCurrentTimestamp(),
	})
}

// GetSystemContext returns the system context formatted for the model
func (cc *ConversationContext) GetSystemContext(compact bool) string {
	var sb strings.Builder

	sb.WriteString(cc.SystemPrompt)
	sb.WriteString("\n\n")

	if compact {
		sb.WriteString(fmt.Sprintf("[TOOLS:%d] ", len(cc.Tools)))
		var names []string
		for _, tool := range cc.Tools {
			names = append(names, tool.Name)
		}
		sb.WriteString(strings.Join(names, ", "))
	} else {
		sb.WriteString("Available Tools:\n")
		for i, tool := range cc.Tools {
			sb.WriteString(fmt.Sprintf("%d. %s - %s\n", i+1, tool.Name, tool.Description))
		}
	}

	return sb.String()
}

// GetConversationHistory returns formatted conversation history
func (cc *ConversationContext) GetConversationHistory(maxMessages int, compact bool) string {
	var sb strings.Builder

	start := len(cc.History) - maxMessages
	if start < 0 {
		start = 0
	}

	for _, msg := range cc.History[start:] {
		if compact {
			sb.WriteString(fmt.Sprintf("[%s] %s\n", shortRole(msg.Role), msg.Content))
		} else {
			sb.WriteString(fmt.Sprintf("%s: %s\n", msg.Role, msg.Content))
		}
	}

	return sb.String()
}

// shortRole returns short role notation for token efficiency
func shortRole(role string) string {
	switch role {
	case "user":
		return "U"
	case "assistant":
		return "A"
	case "system":
		return "S"
	case "tool":
		return "T"
	default:
		return role[:1]
	}
}

// EstimateTokens estimates token count for the context
// Used to stay within model token limits
func (cc *ConversationContext) EstimateTokens() int {
	// Rough estimation: ~1 token per 4 characters
	count := len(cc.SystemPrompt) / 4

	for _, tool := range cc.Tools {
		count += len(tool.Name) / 4
		count += len(tool.Description) / 4
	}

	for _, msg := range cc.History {
		count += len(msg.Content) / 4
	}

	return count
}

// PruneHistory removes old messages if token count exceeds limit
func (cc *ConversationContext) PruneHistory(maxTokens int) {
	if cc.EstimateTokens() <= maxTokens {
		return
	}

	// Keep only recent messages
	if len(cc.History) > 1 {
		cc.History = cc.History[1:]
		cc.PruneHistory(maxTokens) // Recursively prune if needed
	}
}

// getCurrentTimestamp returns current timestamp in RFC3339 format
func getCurrentTimestamp() string {
	return formattedNow()
}

// formattedNow returns formatted current time (defined separately to allow mocking)
var formattedNow = func() string {
	return ""
}

// ContextRenderer renders context in different formats
type ContextRenderer struct {
	format string // "json", "text", "compact", "markdown"
}

// NewContextRenderer creates a new context renderer
func NewContextRenderer(format string) *ContextRenderer {
	return &ContextRenderer{
		format: format,
	}
}

// Render renders the conversation context in the specified format
func (cr *ContextRenderer) Render(cc *ConversationContext) string {
	switch cr.format {
	case "json":
		return cr.renderJSON(cc)
	case "text":
		return cr.renderText(cc)
	case "compact":
		return cr.renderCompact(cc)
	case "markdown":
		return cr.renderMarkdown(cc)
	default:
		return cr.renderText(cc)
	}
}

// renderJSON renders as JSON
func (cr *ContextRenderer) renderJSON(cc *ConversationContext) string {
	ctx := Context{
		Version: "1.0",
		Metadata: Metadata{
			Agent: "aionmcp",
		},
		Tools:    cc.Tools,
		Messages: cc.History,
		State:    cc.State,
	}
	data, _ := json.Marshal(ctx)
	return string(data)
}

// renderText renders as readable text
func (cr *ContextRenderer) renderText(cc *ConversationContext) string {
	var sb strings.Builder

	sb.WriteString(cc.SystemPrompt)
	sb.WriteString("\n\n")

	sb.WriteString("=== Available Tools ===\n")
	for i, tool := range cc.Tools {
		sb.WriteString(fmt.Sprintf("%d. %s\n   %s\n\n", i+1, tool.Name, tool.Description))
	}

	sb.WriteString("=== Conversation History ===\n")
	for _, msg := range cc.History {
		sb.WriteString(fmt.Sprintf("[%s]\n%s\n\n", strings.ToUpper(msg.Role), msg.Content))
	}

	return sb.String()
}

// renderCompact renders in compact format for token efficiency
func (cr *ContextRenderer) renderCompact(cc *ConversationContext) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("P:%q\n", cc.SystemPrompt[:min(50, len(cc.SystemPrompt))]))
	sb.WriteString(fmt.Sprintf("T:%d\n", len(cc.Tools)))
	sb.WriteString(fmt.Sprintf("H:%d\n", len(cc.History)))

	return sb.String()
}

// renderMarkdown renders as markdown
func (cr *ContextRenderer) renderMarkdown(cc *ConversationContext) string {
	var sb strings.Builder

	sb.WriteString("# Conversation Context\n\n")
	sb.WriteString("## System Prompt\n\n")
	sb.WriteString(cc.SystemPrompt)
	sb.WriteString("\n\n## Available Tools\n\n")

	for _, tool := range cc.Tools {
		sb.WriteString(fmt.Sprintf("### %s\n\n", tool.Name))
		sb.WriteString(fmt.Sprintf("%s\n\n", tool.Description))
	}

	sb.WriteString("## Conversation History\n\n")
	for _, msg := range cc.History {
		sb.WriteString(fmt.Sprintf("**%s:**\n\n%s\n\n", msg.Role, msg.Content))
	}

	return sb.String()
}

// min returns minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
