package importer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aionmcp/aionmcp/pkg/types"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
)

// GraphQLImporter handles GraphQL schemas
type GraphQLImporter struct {
	endpoint string // Default GraphQL endpoint
}

// NewGraphQLImporter creates a new GraphQL importer
func NewGraphQLImporter() *GraphQLImporter {
	return &GraphQLImporter{}
}

// GetType returns the specification type
func (i *GraphQLImporter) GetType() SpecType {
	return SpecTypeGraphQL
}

// Supports checks if this importer can handle the given source
func (i *GraphQLImporter) Supports(source SpecSource) bool {
	return source.Type == SpecTypeGraphQL
}

// Validate checks if the GraphQL schema is valid
func (i *GraphQLImporter) Validate(ctx context.Context, source SpecSource) error {
	schemaString, err := i.loadSchema(source.Path)
	if err != nil {
		return err
	}

	// Parse the schema
	_, err = parser.Parse(parser.ParseParams{
		Source: schemaString,
	})
	return err
}

// Import parses the GraphQL schema and generates tools
func (i *GraphQLImporter) Import(ctx context.Context, source SpecSource) (*ImportResult, error) {
	start := time.Now()
	
	result := &ImportResult{
		Source:    source,
		Tools:     []types.Tool{},
		Errors:    []error{},
		Warnings:  []string{},
		Timestamp: start,
	}

	// Load the schema
	schemaString, err := i.loadSchema(source.Path)
	if err != nil {
		result.Errors = append(result.Errors, err)
		result.Duration = time.Since(start)
		return result, err
	}

	// Parse the schema
	doc, err := parser.Parse(parser.ParseParams{
		Source: schemaString,
	})
	if err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("failed to parse GraphQL schema: %w", err))
		result.Duration = time.Since(start)
		return result, err
	}

	// Extract endpoint from metadata or use default
	endpoint := source.Metadata["endpoint"]
	if endpoint == "" {
		endpoint = "http://localhost:4000/graphql" // Default GraphQL endpoint
		result.Warnings = append(result.Warnings, "No GraphQL endpoint specified in metadata, using default: "+endpoint)
	}

	// Generate tools from queries and mutations
	for _, def := range doc.Definitions {
		if typeDef, ok := def.(*ast.ObjectDefinition); ok {
			switch typeDef.Name.Value {
			case "Query":
				for _, field := range typeDef.Fields {
					tool := i.createQueryTool(source, endpoint, field, schemaString)
					result.Tools = append(result.Tools, tool)
				}
			case "Mutation":
				for _, field := range typeDef.Fields {
					tool := i.createMutationTool(source, endpoint, field, schemaString)
					result.Tools = append(result.Tools, tool)
				}
			}
		}
	}

	result.Duration = time.Since(start)
	return result, nil
}

// loadSchema loads a GraphQL schema from file or URL
func (i *GraphQLImporter) loadSchema(path string) (string, error) {
	// Check if it's a URL
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		resp, err := http.Get(path)
		if err != nil {
			return "", fmt.Errorf("failed to fetch schema from URL: %w", err)
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read schema response: %w", err)
		}

		return string(bodyBytes), nil
	}

	// Load from file
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read schema file: %w", err)
	}

	return string(content), nil
}

// createQueryTool creates a tool for a GraphQL query
func (i *GraphQLImporter) createQueryTool(source SpecSource, endpoint string, field *ast.FieldDefinition, schema string) types.Tool {
	return &GraphQLTool{
		source:    source,
		endpoint:  endpoint,
		field:     field,
		schema:    schema,
		operation: "query",
	}
}

// createMutationTool creates a tool for a GraphQL mutation
func (i *GraphQLImporter) createMutationTool(source SpecSource, endpoint string, field *ast.FieldDefinition, schema string) types.Tool {
	return &GraphQLTool{
		source:    source,
		endpoint:  endpoint,
		field:     field,
		schema:    schema,
		operation: "mutation",
	}
}

// GraphQLTool represents a tool generated from a GraphQL operation
type GraphQLTool struct {
	source    SpecSource
	endpoint  string
	field     *ast.FieldDefinition
	schema    string
	operation string // "query" or "mutation"
}

// Name returns the tool name
func (t *GraphQLTool) Name() string {
	return fmt.Sprintf("graphql.%s.%s_%s", t.source.ID, t.operation, t.field.Name.Value)
}

// Description returns the tool description
func (t *GraphQLTool) Description() string {
	// Try to extract description from directives or comments
	description := fmt.Sprintf("GraphQL %s: %s", t.operation, t.field.Name.Value)
	
	// Add argument information
	if len(t.field.Arguments) > 0 {
		var args []string
		for _, arg := range t.field.Arguments {
			args = append(args, arg.Name.Value)
		}
		description += fmt.Sprintf(" (args: %s)", strings.Join(args, ", "))
	}

	return description
}

// Execute performs the GraphQL operation
func (t *GraphQLTool) Execute(input any) (any, error) {
	// Parse input
	inputMap, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("input must be a JSON object")
	}

	// Extract variables from input
	variables, _ := inputMap["variables"].(map[string]interface{})
	if variables == nil {
		variables = make(map[string]interface{})
	}

	// Copy non-variables fields as variables
	for key, value := range inputMap {
		if key != "variables" {
			variables[key] = value
		}
	}

	// Build GraphQL query/mutation
	query := t.buildQuery(variables)

	// Create GraphQL request
	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	// Execute GraphQL request
	response, err := t.executeGraphQLRequest(requestBody)
	if err != nil {
		return nil, fmt.Errorf("GraphQL request failed: %w", err)
	}

	return response, nil
}

// buildQuery builds the GraphQL query/mutation string
func (t *GraphQLTool) buildQuery(variables map[string]interface{}) string {
	// Build arguments string
	var argsBuilder strings.Builder
	var varsBuilder strings.Builder
	
	if len(t.field.Arguments) > 0 {
		argsBuilder.WriteString("(")
		varsBuilder.WriteString("(")
		
		for i, arg := range t.field.Arguments {
			if i > 0 {
				argsBuilder.WriteString(", ")
				varsBuilder.WriteString(", ")
			}
			
			argName := arg.Name.Value
			argsBuilder.WriteString(fmt.Sprintf("%s: $%s", argName, argName))
			
			// Get type from AST
			typeStr := t.getTypeString(arg.Type)
			varsBuilder.WriteString(fmt.Sprintf("$%s: %s", argName, typeStr))
		}
		
		argsBuilder.WriteString(")")
		varsBuilder.WriteString(")")
	}

	// Build the query
	var queryBuilder strings.Builder
	queryBuilder.WriteString(t.operation)
	if varsBuilder.Len() > 2 { // More than just "()"
		queryBuilder.WriteString(" ")
		queryBuilder.WriteString(varsBuilder.String())
	}
	queryBuilder.WriteString(" { ")
	queryBuilder.WriteString(t.field.Name.Value)
	if argsBuilder.Len() > 2 { // More than just "()"
		queryBuilder.WriteString(argsBuilder.String())
	}
	
	// Add selection set (simplified - in real implementation, this would be more sophisticated)
	queryBuilder.WriteString(" { __typename } ")
	queryBuilder.WriteString(" }")

	return queryBuilder.String()
}

// getTypeString converts AST type to string
func (t *GraphQLTool) getTypeString(typeNode ast.Type) string {
	switch node := typeNode.(type) {
	case *ast.Named:
		return node.Name.Value
	case *ast.NonNull:
		return t.getTypeString(node.Type) + "!"
	case *ast.List:
		return "[" + t.getTypeString(node.Type) + "]"
	default:
		return "String" // Default fallback
	}
}

// executeGraphQLRequest executes the HTTP request to the GraphQL endpoint
func (t *GraphQLTool) executeGraphQLRequest(requestBody map[string]interface{}) (interface{}, error) {
	// Marshal request body
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", t.endpoint, strings.NewReader(string(bodyBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Execute request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode GraphQL response: %w", err)
	}

	// Check for GraphQL errors
	if errors, exists := response["errors"]; exists {
		return map[string]interface{}{
			"errors":      errors,
			"data":        response["data"],
			"status_code": resp.StatusCode,
			"endpoint":    t.endpoint,
		}, nil
	}

	return map[string]interface{}{
		"data":        response["data"],
		"status_code": resp.StatusCode,
		"endpoint":    t.endpoint,
	}, nil
}

// Metadata returns tool metadata
func (t *GraphQLTool) Metadata() types.ToolMetadata {
	// Build input schema from GraphQL field arguments
	inputSchema := map[string]interface{}{
		"type":       "object",
		"properties": make(map[string]interface{}),
		"required":   []string{},
	}

	properties := inputSchema["properties"].(map[string]interface{})
	var required []string

	// Add field arguments to schema
	for _, arg := range t.field.Arguments {
		argSchema := map[string]interface{}{
			"type":        "string", // Simplified type mapping
			"description": fmt.Sprintf("GraphQL argument: %s", arg.Name.Value),
		}

		properties[arg.Name.Value] = argSchema

		// Check if argument is non-null (required)
		if _, isNonNull := arg.Type.(*ast.NonNull); isNonNull {
			required = append(required, arg.Name.Value)
		}
	}

	// Add variables object for complex inputs
	properties["variables"] = map[string]interface{}{
		"type":        "object",
		"description": "GraphQL variables object",
	}

	inputSchema["required"] = required

	return types.ToolMetadata{
		Name:        t.Name(),
		Description: t.Description(),
		Version:     "1.0.0",
		Source:      string(SpecTypeGraphQL),
		Tags:        []string{"graphql", t.operation, "api"},
		Schema: map[string]interface{}{
			"input": inputSchema,
			"output": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"data":        map[string]interface{}{"type": "object"},
					"errors":      map[string]interface{}{"type": "array"},
					"status_code": map[string]interface{}{"type": "integer"},
					"endpoint":    map[string]interface{}{"type": "string"},
				},
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}