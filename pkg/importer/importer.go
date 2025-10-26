package importer

import (
	"context"
	"fmt"
	"time"

	"github.com/aionmcp/aionmcp/pkg/types"
)

// SpecType represents the type of API specification
type SpecType string

const (
	SpecTypeOpenAPI  SpecType = "openapi"
	SpecTypeGraphQL  SpecType = "graphql"
	SpecTypeAsyncAPI SpecType = "asyncapi"
)

// SpecSource represents a specification source
type SpecSource struct {
	ID          string            `json:"id"`
	Type        SpecType          `json:"type"`
	Path        string            `json:"path"`        // File path or URL
	Name        string            `json:"name"`        // Human-readable name
	Description string            `json:"description"` // Description of the API
	Metadata    map[string]string `json:"metadata"`    // Additional metadata
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// ImportResult contains the result of importing a specification
type ImportResult struct {
	Source    SpecSource     `json:"source"`
	Tools     []types.Tool   `json:"tools"`
	Errors    []error        `json:"errors"`
	Warnings  []string       `json:"warnings"`
	Duration  time.Duration  `json:"duration"`
	Timestamp time.Time      `json:"timestamp"`
}

// SpecImporter is the interface for importing API specifications
type SpecImporter interface {
	// GetType returns the specification type this importer handles
	GetType() SpecType
	
	// Validate checks if the specification is valid for this importer
	Validate(ctx context.Context, source SpecSource) error
	
	// Import parses the specification and generates tools
	Import(ctx context.Context, source SpecSource) (*ImportResult, error)
	
	// Supports checks if this importer can handle the given source
	Supports(source SpecSource) bool
}

// ToolRegistry interface to avoid circular imports
type ToolRegistry interface {
	Register(tool types.Tool) error
	Unregister(name string) error
}

// ImporterManager manages all specification importers
type ImporterManager struct {
	importers map[SpecType]SpecImporter
	registry  ToolRegistry
	sources   map[string]SpecSource // source ID -> source
}

// NewImporterManager creates a new importer manager
func NewImporterManager(registry ToolRegistry) *ImporterManager {
	return &ImporterManager{
		importers: make(map[SpecType]SpecImporter),
		registry:  registry,
		sources:   make(map[string]SpecSource),
	}
}

// RegisterImporter registers a new specification importer
func (m *ImporterManager) RegisterImporter(importer SpecImporter) {
	m.importers[importer.GetType()] = importer
}

// ImportSpec imports a specification and registers the generated tools
func (m *ImporterManager) ImportSpec(ctx context.Context, source SpecSource) (*ImportResult, error) {
	// Find appropriate importer
	importer, exists := m.importers[source.Type]
	if !exists {
		return nil, fmt.Errorf("no importer found for spec type: %s", source.Type)
	}

	// Validate specification
	if err := importer.Validate(ctx, source); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Import and generate tools
	result, err := importer.Import(ctx, source)
	if err != nil {
		return nil, fmt.Errorf("import failed: %w", err)
	}

	// Register tools with the registry
	for _, tool := range result.Tools {
		if err := m.registry.Register(tool); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("failed to register tool %s: %w", tool.Name(), err))
		}
	}

	// Store source information
	m.sources[source.ID] = source

	return result, nil
}

// RemoveSpec removes a specification and unregisters its tools
func (m *ImporterManager) RemoveSpec(ctx context.Context, sourceID string) error {
	source, exists := m.sources[sourceID]
	if !exists {
		return fmt.Errorf("specification source not found: %s", sourceID)
	}

	// Find importer
	importer, exists := m.importers[source.Type]
	if !exists {
		return fmt.Errorf("no importer found for spec type: %s", source.Type)
	}

	// Re-import to get tool names (we could cache this for efficiency)
	result, err := importer.Import(ctx, source)
	if err != nil {
		return fmt.Errorf("failed to re-import for removal: %w", err)
	}

	// Unregister tools
	for _, tool := range result.Tools {
		if err := m.registry.Unregister(tool.Name()); err != nil {
			// Log warning but continue
			continue
		}
	}

	// Remove source
	delete(m.sources, sourceID)

	return nil
}

// ReloadSpec reloads a specification (useful for file watching)
func (m *ImporterManager) ReloadSpec(ctx context.Context, sourceID string) (*ImportResult, error) {
	source, exists := m.sources[sourceID]
	if !exists {
		return nil, fmt.Errorf("specification source not found: %s", sourceID)
	}

	// Remove existing tools
	if err := m.RemoveSpec(ctx, sourceID); err != nil {
		return nil, fmt.Errorf("failed to remove existing spec: %w", err)
	}

	// Re-import
	source.UpdatedAt = time.Now()
	return m.ImportSpec(ctx, source)
}

// ListSources returns all registered specification sources
func (m *ImporterManager) ListSources() []SpecSource {
	sources := make([]SpecSource, 0, len(m.sources))
	for _, source := range m.sources {
		sources = append(sources, source)
	}
	return sources
}

// GetSource returns a specific specification source
func (m *ImporterManager) GetSource(sourceID string) (SpecSource, bool) {
	source, exists := m.sources[sourceID]
	return source, exists
}

// GetSupportedTypes returns all supported specification types
func (m *ImporterManager) GetSupportedTypes() []SpecType {
	types := make([]SpecType, 0, len(m.importers))
	for specType := range m.importers {
		types = append(types, specType)
	}
	return types
}