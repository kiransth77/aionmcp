# TOON (Token-Oriented Object Notation) - LLM Context Format

**Optional Feature**: TOON context formatting for optimized LLM interaction.

---

## 📝 Overview

TOON is a token-efficient format for preparing context and tool information for Large Language Models. It helps:

- **Reduce token usage** by compacting information
- **Improve clarity** with structured formatting
- **Optimize costs** when working with paid LLM APIs
- **Better context management** for multi-turn conversations

---

## 🚀 Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/kiransth77/aionmcp/pkg/toon"
)

func main() {
    // Create a TOON context builder
    ctx := toon.NewBuilder("aionmcp", "github-copilot").
        WithModel("gpt-4").
        WithSessionID("session-123").
        WithUser("user@example.com")
    
    // Add tools
    ctx.AddTool(toon.ToolContext{
        ID:          "petstore-list-pets",
        Name:        "List Pets",
        Description: "Get list of available pets",
        Type:        "openapi",
    })
    
    // Add conversation history
    ctx.AddMessage("user", "Show me available pets")
    ctx.AddMessage("assistant", "I'll list the pets for you...")
    
    // Convert to JSON
    json, _ := ctx.ToJSON()
    fmt.Println(json)
}
```

### Token-Efficient Format

```go
// Compact JSON (minimal tokens)
compactJSON, _ := ctx.ToCompactJSON()
fmt.Println(compactJSON) // {"version":"1.0","timestamp":"...","metadata":{...},...}

// Token count estimate
tokenCount := ctx.Build().TokenCount()
fmt.Printf("Estimated tokens: %d\n", tokenCount)
```

---

## 🎯 Use Cases

### 1. Optimized Model Context

Send only necessary context to reduce API costs:

```go
context := toon.NewConversationContext(
    "You are a helpful API assistant.",
).
WithTool(myTool).
GetSystemContext(compact=true) // Returns compact version
```

### 2. Multi-Turn Conversations

Manage conversation history efficiently:

```go
conversation := toon.NewConversationContext("System prompt")

conversation.AddMessage("user", "Call the petstore API")
conversation.AddMessage("assistant", "I'll help with that")
conversation.AddMessage("user", "List all cats")

// Prune if too many tokens
conversation.PruneHistory(maxTokens=2000)
```

### 3. Tool Discovery for Models

Format available tools for model consumption:

```go
formatter := toon.NewToolFormatter(compact=true)
description := formatter.FormatToolDescription(
    "petstore-list",
    "List Pets",
    "Returns available pets",
    inputSchema,
    []string{"api", "petstore"},
)
```

### 4. Context Rendering

Render in different formats for different consumers:

```go
renderer := toon.NewContextRenderer("compact") // or "json", "text", "markdown"
output := renderer.Render(conversationContext)
```

---

## 📊 TOON Structure

### Context JSON

```json
{
  "version": "1.0",
  "timestamp": "2025-11-27T22:30:00Z",
  "metadata": {
    "agent": "aionmcp",
    "model": "gpt-4",
    "session_id": "session-123",
    "user": "user@example.com",
    "source": "github-copilot",
    "environment": "production"
  },
  "tools": [
    {
      "id": "petstore-list-pets",
      "name": "List Pets",
      "description": "Get list of available pets",
      "type": "openapi",
      "input_schema": { "type": "object", "properties": {} },
      "tags": ["api", "petstore"],
      "confidence": 0.95
    }
  ],
  "state": {
    "conversation_count": 5,
    "api_calls": 3
  },
  "messages": [
    {
      "role": "user",
      "content": "List all available pets",
      "timestamp": "2025-11-27T22:30:00Z"
    },
    {
      "role": "assistant",
      "content": "I'll retrieve the pet list...",
      "timestamp": "2025-11-27T22:30:01Z"
    }
  ]
}
```

### Compact Notation

For token efficiency, compact versions use abbreviated keys:

```json
{"v":"1.0","ag":"aionmcp","src":"copilot","t":3,"m":5,"s":2}
```

---

## 🔧 API Reference

### Builder Pattern

```go
// Create builder
builder := toon.NewBuilder(agent, source)

// Configure
builder.WithModel(model)
builder.WithSessionID(sessionID)
builder.WithUser(user)
builder.WithEnvironment(env)

// Add content
builder.AddTool(toolContext)
builder.AddState(key, value)
builder.AddMessage(role, content)
builder.AddToolMessage(role, content, toolID)

// Build and export
context := builder.Build()
json, err := builder.ToJSON()
compact, err := builder.ToCompactJSON()
```

### Conversation Context

```go
// Create
cc := toon.NewConversationContext("System prompt")

// Add content
cc.AddTool(tool)
cc.AddMessage(role, content)

// Query
systemContext := cc.GetSystemContext(compact)
history := cc.GetConversationHistory(maxMessages, compact)
tokens := cc.EstimateTokens()

// Manage
cc.PruneHistory(maxTokens) // Remove old messages if over limit
```

### Rendering

```go
renderer := toon.NewContextRenderer("compact")
// Formats: "json", "text", "compact", "markdown"

output := renderer.Render(conversationContext)
```

---

## 💡 Optimization Tips

### 1. Use Compact Mode for Large Contexts

```go
// For contexts with many tools/messages
formatter := toon.NewToolFormatter(compact=true)
context := cc.GetSystemContext(compact=true)
```

### 2. Prune Conversation History

```go
// Keep only recent messages within token limit
if context.EstimateTokens() > 3000 {
    context.PruneHistory(2500)
}
```

### 3. Selective Tool Inclusion

```go
// Only include relevant tools
relevantTools := filterToolsByCategory(allTools, "api")
for _, tool := range relevantTools {
    builder.AddTool(tool)
}
```

### 4. Token-Efficient Rendering

```go
// Use compact JSON for transmission
json, _ := builder.ToCompactJSON() // Smaller than formatted
// Use readable format for development
json, _ := builder.ToJSON() // Pretty-printed
```

---

## 📈 Performance Examples

### Token Savings

```
Full verbose context:
  - Tools: 1500 chars
  - Messages: 5000 chars
  - Metadata: 500 chars
  Total: 7000 chars ≈ 1750 tokens

Compact context:
  - Tools (compact): 600 chars
  - Messages (compact): 2000 chars
  - Metadata (compact): 150 chars
  Total: 2750 chars ≈ 688 tokens

Savings: 60% tokens reduction
```

### Cost Impact (GPT-4 pricing as example)

```
Verbose: 1750 tokens input × $0.03/1K = $0.0525
Compact: 688 tokens input × $0.03/1K = $0.0206
Savings per request: $0.0319 (61%)
Savings at 1000 requests: $31.90
```

---

## 🔌 Integration Examples

### With GitHub Copilot

```go
// Format context for Copilot
builder := toon.NewBuilder("aionmcp", "github-copilot")
builder.WithModel("gpt-4")
builder.AddTool(availableTool)

// Send to Copilot
context, _ := builder.ToJSON()
copilotAPI.SetContext(context)
```

### With Claude

```go
builder := toon.NewBuilder("aionmcp", "anthropic-claude")
builder.WithModel("claude-3-opus")

// Claude context format
systemPrompt := builder.Build().FormatForModel()
claudeAPI.CreateMessage(systemPrompt, messages)
```

### With Custom LLM

```go
renderer := toon.NewContextRenderer("markdown")
context := renderer.Render(conversationContext)

// Send to your LLM
response := yourLLM.Complete(context)
```

---

## 🎓 Best Practices

1. **Use Compact Mode for Production**
   - Reduces token usage and API costs
   - Faster response times

2. **Monitor Token Count**
   - Check `EstimateTokens()` regularly
   - Prune history if approaching limits

3. **Include Relevant Context Only**
   - Filter tools by relevance
   - Remove completed tasks from state

4. **Timestamp Messages**
   - Helps models understand sequence
   - Useful for debugging

5. **Use Confidence Scores**
   - Indicate reliability of tools
   - Helps models choose better options

---

## 📚 Examples

### Full Multi-Tool Context

```go
builder := toon.NewBuilder("aionmcp", "copilot").
    WithModel("gpt-4").
    WithSessionID("abc123")

// Add multiple tools
for _, tool := range tools {
    builder.AddTool(tool)
}

// Add conversation
builder.AddMessage("user", "What can I do?")
builder.AddMessage("assistant", "I have access to these APIs...")

// Add state
builder.AddState("conversation_count", 1)
builder.AddState("rate_limited", false)

context, _ := builder.ToJSON()
```

### Token-Conscious Conversation

```go
cc := toon.NewConversationContext("You are helpful")

// Add messages while monitoring tokens
for {
    tokens := cc.EstimateTokens()
    if tokens > 3500 {
        cc.PruneHistory(3000)
    }
    
    cc.AddMessage("user", userInput)
    // ... call model ...
    cc.AddMessage("assistant", response)
}
```

---

## 🐛 Troubleshooting

### Context Too Large

```go
// Solution: Use compact format and prune history
compact, _ := ctx.ToCompactJSON()
cc.PruneHistory(maxTokens)
```

### Too Much Information

```go
// Solution: Filter tools and use compact rendering
relevantTools := filterByTag(allTools, "petstore")
renderer := toon.NewContextRenderer("compact")
```

### Integration Issues

```go
// Verify format is correct
json, err := ctx.ToJSON()
if err != nil {
    log.Printf("Failed to serialize context: %v", err)
}
```

---

## 📞 Support

- API Documentation: See inline Go docs
- Examples: `examples/toon/`
- Issues: GitHub Issues with `[toon]` prefix

---

## 🎯 Future Enhancements

- [ ] TOON schema validation
- [ ] Streaming context updates
- [ ] Context versioning
- [ ] Multi-language support
- [ ] Advanced token optimization
- [ ] Context caching
- [ ] Distributed context sharing

---

**Status**: Optional feature, ready for production use  
**Version**: 1.0  
**Last Updated**: November 27, 2025
