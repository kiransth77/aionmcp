package importer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aionmcp/aionmcp/pkg/types"
	"github.com/getkin/kin-openapi/openapi3"
)

// OpenAPIImporter handles OpenAPI 3.x specifications
type OpenAPIImporter struct {
	loader *openapi3.Loader
}

// NewOpenAPIImporter creates a new OpenAPI importer
func NewOpenAPIImporter() *OpenAPIImporter {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	return &OpenAPIImporter{
		loader: loader,
	}
}

// GetType returns the specification type
func (i *OpenAPIImporter) GetType() SpecType {
	return SpecTypeOpenAPI
}

// Supports checks if this importer can handle the given source
func (i *OpenAPIImporter) Supports(source SpecSource) bool {
	return source.Type == SpecTypeOpenAPI
}

// Validate checks if the specification is valid
func (i *OpenAPIImporter) Validate(ctx context.Context, source SpecSource) error {
	_, err := i.loadSpec(ctx, source.Path)
	return err
}

// Import parses the OpenAPI specification and generates tools
func (i *OpenAPIImporter) Import(ctx context.Context, source SpecSource) (*ImportResult, error) {
	start := time.Now()
	
	result := &ImportResult{
		Source:    source,
		Tools:     []types.Tool{},
		Errors:    []error{},
		Warnings:  []string{},
		Timestamp: start,
	}

	// Load the specification
	doc, err := i.loadSpec(ctx, source.Path)
	if err != nil {
		result.Errors = append(result.Errors, err)
		result.Duration = time.Since(start)
		return result, err
	}

	// Validate the loaded specification
	if err := doc.Validate(ctx); err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Specification validation warning: %v", err))
	}

	// Generate tools from paths
	for path, pathItem := range doc.Paths.Map() {
		// Generate tools for each HTTP method
		methods := map[string]*openapi3.Operation{
			"GET":     pathItem.Get,
			"POST":    pathItem.Post,
			"PUT":     pathItem.Put,
			"PATCH":   pathItem.Patch,
			"DELETE":  pathItem.Delete,
			"HEAD":    pathItem.Head,
			"OPTIONS": pathItem.Options,
		}

		for method, operation := range methods {
			if operation == nil {
				continue
			}

			tool, err := i.createToolFromOperation(source, doc, path, method, operation)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("failed to create tool for %s %s: %w", method, path, err))
				continue
			}

			result.Tools = append(result.Tools, tool)
		}
	}

	result.Duration = time.Since(start)
	return result, nil
}

// loadSpec loads an OpenAPI specification from file or URL
func (i *OpenAPIImporter) loadSpec(ctx context.Context, path string) (*openapi3.T, error) {
	// Check if it's a URL
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		parsedURL, err := url.Parse(path)
		if err != nil {
			return nil, fmt.Errorf("invalid URL: %w", err)
		}
		return i.loader.LoadFromURI(parsedURL)
	}

	// Load from file
	return i.loader.LoadFromFile(path)
}

// createToolFromOperation creates an MCP tool from an OpenAPI operation
func (i *OpenAPIImporter) createToolFromOperation(source SpecSource, doc *openapi3.T, path, method string, operation *openapi3.Operation) (types.Tool, error) {
	tool := &OpenAPITool{
		source:    source,
		doc:       doc,
		path:      path,
		method:    method,
		operation: operation,
	}

	return tool, nil
}

// OpenAPITool represents a tool generated from an OpenAPI operation
type OpenAPITool struct {
	source    SpecSource
	doc       *openapi3.T
	path      string
	method    string
	operation *openapi3.Operation
}

// Name returns the tool name
func (t *OpenAPITool) Name() string {
	// Use operationId if available, otherwise generate from path and method
	if t.operation.OperationID != "" {
		return fmt.Sprintf("openapi.%s.%s", t.source.ID, t.operation.OperationID)
	}

	// Generate name from path and method
	cleanPath := strings.ReplaceAll(strings.Trim(t.path, "/"), "/", "_")
	cleanPath = strings.ReplaceAll(cleanPath, "{", "")
	cleanPath = strings.ReplaceAll(cleanPath, "}", "")
	return fmt.Sprintf("openapi.%s.%s_%s", t.source.ID, strings.ToLower(t.method), cleanPath)
}

// Description returns the tool description
func (t *OpenAPITool) Description() string {
	if t.operation.Summary != "" {
		return t.operation.Summary
	}
	if t.operation.Description != "" {
		return t.operation.Description
	}
	return fmt.Sprintf("%s %s operation from %s", t.method, t.path, t.source.Name)
}

// Execute performs the API call
func (t *OpenAPITool) Execute(input any) (any, error) {
	// Parse input parameters
	params, err := t.parseInput(input)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	// Build the request URL
	baseURL := ""
	if len(t.doc.Servers) > 0 {
		baseURL = t.doc.Servers[0].URL
	}

	// Replace path parameters
	requestPath := t.path
	for paramName, paramValue := range params.Path {
		placeholder := fmt.Sprintf("{%s}", paramName)
		requestPath = strings.ReplaceAll(requestPath, placeholder, fmt.Sprintf("%v", paramValue))
	}

	// Build full URL
	fullURL, err := url.JoinPath(baseURL, requestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	// Add query parameters
	if len(params.Query) > 0 {
		parsedURL, err := url.Parse(fullURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse URL: %w", err)
		}
		
		query := parsedURL.Query()
		for key, value := range params.Query {
			query.Add(key, fmt.Sprintf("%v", value))
		}
		parsedURL.RawQuery = query.Encode()
		fullURL = parsedURL.String()
	}

	// Create HTTP request
	req, err := http.NewRequest(t.method, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for key, value := range params.Headers {
		req.Header.Set(key, fmt.Sprintf("%v", value))
	}

	// Add request body for POST, PUT, PATCH
	if params.Body != nil && (t.method == "POST" || t.method == "PUT" || t.method == "PATCH") {
		bodyBytes, err := json.Marshal(params.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		req.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute the request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var responseBody interface{}
	if resp.Header.Get("Content-Type") == "application/json" {
		if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
			return nil, fmt.Errorf("failed to decode JSON response: %w", err)
		}
	} else {
		// For non-JSON responses, return as string
		bodyBytes := make([]byte, resp.ContentLength)
		if _, err := resp.Body.Read(bodyBytes); err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		responseBody = string(bodyBytes)
	}

	return map[string]interface{}{
		"status_code": resp.StatusCode,
		"headers":     resp.Header,
		"body":        responseBody,
		"request_url": fullURL,
		"method":      t.method,
	}, nil
}

// RequestParams holds parsed request parameters
type RequestParams struct {
	Path    map[string]interface{} `json:"path"`
	Query   map[string]interface{} `json:"query"`
	Headers map[string]interface{} `json:"headers"`
	Body    interface{}            `json:"body"`
}

// parseInput parses the input into request parameters
func (t *OpenAPITool) parseInput(input any) (*RequestParams, error) {
	params := &RequestParams{
		Path:    make(map[string]interface{}),
		Query:   make(map[string]interface{}),
		Headers: make(map[string]interface{}),
	}

	// Convert input to map
	inputMap, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("input must be a JSON object")
	}

	// Extract parameters based on OpenAPI operation definition
	for _, param := range t.operation.Parameters {
		paramName := param.Value.Name
		paramValue, exists := inputMap[paramName]
		if !exists && param.Value.Required {
			return nil, fmt.Errorf("required parameter '%s' is missing", paramName)
		}

		if exists {
			switch param.Value.In {
			case "path":
				params.Path[paramName] = paramValue
			case "query":
				params.Query[paramName] = paramValue
			case "header":
				params.Headers[paramName] = paramValue
			}
		}
	}

	// Extract request body
	if body, exists := inputMap["body"]; exists {
		params.Body = body
	}

	return params, nil
}

// Metadata returns tool metadata
func (t *OpenAPITool) Metadata() types.ToolMetadata {
	// Build input schema from OpenAPI parameters
	inputSchema := map[string]interface{}{
		"type":       "object",
		"properties": make(map[string]interface{}),
		"required":   []string{},
	}

	properties := inputSchema["properties"].(map[string]interface{})
	var required []string

	// Add parameters to schema
	for _, param := range t.operation.Parameters {
		paramSchema := map[string]interface{}{
			"type":        "string", // Simplified for now
			"description": param.Value.Description,
		}

		properties[param.Value.Name] = paramSchema

		if param.Value.Required {
			required = append(required, param.Value.Name)
		}
	}

	// Add request body if present
	if t.operation.RequestBody != nil {
		properties["body"] = map[string]interface{}{
			"type":        "object",
			"description": "Request body",
		}
	}

	inputSchema["required"] = required

	return types.ToolMetadata{
		Name:        t.Name(),
		Description: t.Description(),
		Version:     "1.0.0",
		Source:      string(SpecTypeOpenAPI),
		Tags:        []string{"openapi", "api", strings.ToLower(t.method)},
		Schema: map[string]interface{}{
			"input": inputSchema,
			"output": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"status_code": map[string]interface{}{"type": "integer"},
					"headers":     map[string]interface{}{"type": "object"},
					"body":        map[string]interface{}{"type": "object"},
					"request_url": map[string]interface{}{"type": "string"},
					"method":      map[string]interface{}{"type": "string"},
				},
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}