package types

import "time"

// Tool represents an MCP tool interface
type Tool interface {
	Name() string
	Description() string
	Execute(input any) (any, error)
	Metadata() ToolMetadata
}

// ToolMetadata contains metadata about a tool
type ToolMetadata struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Version     string         `json:"version"`
	Source      string         `json:"source"` // openapi, graphql, asyncapi
	Tags        []string       `json:"tags"`
	Schema      map[string]any `json:"schema"` // Input/output schema
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
