package types

// ToolRegistry defines the interface for tool registry operations
type ToolRegistry interface {
	// Basic registry operations
	Get(name string) (Tool, error)
	ListTools() []ToolMetadata
	Count() int

	// Dynamic registration
	Register(tool Tool) error
	RegisterWithSource(tool Tool, sourceID, version string) error
	RegisterBatch(tools []Tool, sourceID string) error
	Unregister(name string) error
	UnregisterBySource(sourceID string) error

	// Version and source tracking
	GetVersion(name string) (string, error)
	GetSource(name string) (string, error)
	ListToolsBySource(sourceID string) []ToolMetadata
	GetToolSources() []string

	// Statistics
	GetRegistryStats() map[string]interface{}
}
