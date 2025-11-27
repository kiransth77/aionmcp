package main

import (
	"fmt"
	"log"

	"github.com/kiransth77/aionmcp/pkg/toon"
)

// Example demonstrating TOON (Token-Oriented Object Notation) context for LLMs
func main() {
	fmt.Println("=== TOON Context Builder Example ===\n")

	// Example 1: Basic Context
	basicExample()

	// Example 2: Conversation with Tools
	conversationExample()

	// Example 3: Token Optimization
	tokenOptimizationExample()

	// Example 4: Different Renderers
	rendererExample()
}

func basicExample() {
	fmt.Println("1. Basic Context Creation")
	fmt.Println("------------------------")

	// Create a TOON context for an AI agent
	builder := toon.NewBuilder("aionmcp", "github-copilot").
		WithModel("gpt-4").
		WithSessionID("demo-session-001").
		WithUser("developer@example.com").
		WithEnvironment("development")

	// Add a tool (API endpoint)
	builder.AddTool(toon.ToolContext{
		ID:          "petstore-list-pets",
		Name:        "List Pets",
		Description: "Retrieve a list of all available pets from the store",
		Type:        "openapi",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"status": map[string]interface{}{
					"type":        "string",
					"description": "Filter by status (available, pending, sold)",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Maximum number of results",
				},
			},
		},
		Tags:       []string{"api", "petstore", "read"},
		Confidence: 0.95,
	})

	// Add state information
	builder.AddState("api_calls_made", 0)
	builder.AddState("cache_size", "512MB")

	// Add conversation messages
	builder.AddMessage("user", "Show me all available pets")
	builder.AddMessage("assistant", "I'll retrieve the available pets for you")

	// Build and display
	context := builder.Build()
	fmt.Println(context.ToString())

	// Get token count
	tokenCount := context.TokenCount()
	fmt.Printf("Estimated token count: %d\n\n", tokenCount)
}

func conversationExample() {
	fmt.Println("2. Conversation Context")
	fmt.Println("----------------------")

	// Create a conversation context with system prompt
	systemPrompt := `You are a helpful API assistant. You have access to various APIs 
and can help users interact with them. Always use tools to fetch data when needed.`

	cc := toon.NewConversationContext(systemPrompt)

	// Add tools
	cc.AddTool(toon.ToolContext{
		ID:          "openapi-petstore",
		Name:        "Petstore API",
		Description: "Access to pet store operations",
		Type:        "openapi",
		Tags:        []string{"petstore"},
	})

	cc.AddTool(toon.ToolContext{
		ID:          "graphql-users",
		Name:        "Users GraphQL",
		Description: "Query user data via GraphQL",
		Type:        "graphql",
		Tags:        []string{"users"},
	})

	// Add conversation history
	cc.AddMessage("user", "List all available pets")
	cc.AddMessage("assistant", "I'll fetch the pets from the Petstore API")
	cc.AddMessage("user", "Also show me user information for user123")
	cc.AddMessage("assistant", "I can get both. Let me query the APIs...")

	// Display system context
	fmt.Println("System Context (compact):")
	fmt.Println(cc.GetSystemContext(true))

	fmt.Println("\nConversation History (last 3 messages):")
	fmt.Println(cc.GetConversationHistory(3, false))

	// Token estimation
	fmt.Printf("Estimated tokens: %d\n\n", cc.EstimateTokens())
}

func tokenOptimizationExample() {
	fmt.Println("3. Token Optimization")
	fmt.Println("--------------------")

	// Create context with verbose format
	builder := toon.NewBuilder("aionmcp", "copilot").
		WithModel("gpt-4")

	// Add multiple tools
	for i := 1; i <= 5; i++ {
		builder.AddTool(toon.ToolContext{
			ID:          fmt.Sprintf("tool-%d", i),
			Name:        fmt.Sprintf("Tool %d", i),
			Description: fmt.Sprintf("This is tool number %d with a longer description", i),
			Type:        "openapi",
		})
	}

	// Add conversation history
	for i := 1; i <= 5; i++ {
		builder.AddMessage("user", fmt.Sprintf("Question %d: Tell me about tool %d", i, i))
		builder.AddMessage("assistant", fmt.Sprintf("Answer %d: Tool %d is useful for...", i, i))
	}

	context := builder.Build()

	// Compare formats
	verbose, _ := builder.ToJSON()
	compact, _ := builder.ToCompactJSON()

	fmt.Printf("Verbose JSON: %d characters\n", len(verbose))
	fmt.Printf("Compact JSON: %d characters\n", len(compact))
	fmt.Printf("Tokens saved: ~%d tokens (%.1f%% reduction)\n\n",
		(len(verbose)-len(compact))/4, float64(len(verbose)-len(compact))/float64(len(verbose))*100)

	// Token summary
	fmt.Printf("Total estimated tokens: %d\n\n", context.TokenCount())
}

func rendererExample() {
	fmt.Println("4. Context Rendering")
	fmt.Println("-------------------")

	cc := toon.NewConversationContext("You are helpful")

	cc.AddTool(toon.ToolContext{
		ID:          "tool-1",
		Name:        "Sample Tool",
		Description: "A sample tool for demonstration",
		Type:        "openapi",
	})

	cc.AddMessage("user", "Hello")
	cc.AddMessage("assistant", "Hi there!")

	// Render in different formats
	formats := []string{"compact", "text", "markdown", "json"}

	for _, format := range formats {
		fmt.Printf("Format: %s\n", format)
		fmt.Println("----------")

		renderer := toon.NewContextRenderer(format)
		output := renderer.Render(cc)

		// Limit output for readability
		if len(output) > 200 {
			fmt.Println(output[:200] + "...\n")
		} else {
			fmt.Println(output + "\n")
		}
	}
}

func init() {
	// Handle any initialization if needed
	log.SetFlags(log.Lshortfile)
}
