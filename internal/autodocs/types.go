package autodocs

import (
	"time"
)

// DocumentType represents the type of document being generated
type DocumentType string

const (
	DocumentTypeChangelog    DocumentType = "changelog"
	DocumentTypeReflection   DocumentType = "reflection"
	DocumentTypeReadme       DocumentType = "readme"
	DocumentTypeArchitecture DocumentType = "architecture"
)

// GenerationRequest represents a request to generate documentation
type GenerationRequest struct {
	Type        DocumentType `json:"type"`
	OutputPath  string       `json:"output_path"`
	DateRange   *DateRange   `json:"date_range,omitempty"`
	IncludeData bool         `json:"include_data"`
	Format      string       `json:"format"` // markdown, html, json
}

// DateRange specifies a time range for documentation generation
type DateRange struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// GenerationResult contains the result of document generation
type GenerationResult struct {
	Type          DocumentType `json:"type"`
	OutputPath    string       `json:"output_path"`
	Success       bool         `json:"success"`
	Error         string       `json:"error,omitempty"`
	GeneratedAt   time.Time    `json:"generated_at"`
	ContentLength int          `json:"content_length"`
	Metadata      interface{}  `json:"metadata,omitempty"`
}

// DocumentMetadata contains metadata about generated documents
type DocumentMetadata struct {
	Version       string            `json:"version"`
	GeneratedAt   time.Time         `json:"generated_at"`
	DataSources   []string          `json:"data_sources"`
	CommitRange   *CommitRange      `json:"commit_range,omitempty"`
	LearningStats *LearningSnapshot `json:"learning_stats,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
}

// CommitRange represents a range of git commits
type CommitRange struct {
	StartCommit string    `json:"start_commit"`
	EndCommit   string    `json:"end_commit"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	CommitCount int       `json:"commit_count"`
}

// LearningSnapshot contains learning system data for documentation
type LearningSnapshot struct {
	TotalExecutions int              `json:"total_executions"`
	SuccessRate     float64          `json:"success_rate"`
	AvgLatency      time.Duration    `json:"avg_latency"`
	TopTools        []ToolUsageInfo  `json:"top_tools"`
	ErrorBreakdown  map[string]int   `json:"error_breakdown"`
	RecentPatterns  []PatternSummary `json:"recent_patterns"`
	ActiveInsights  []InsightSummary `json:"active_insights"`
	SnapshotTime    time.Time        `json:"snapshot_time"`
}

// ToolUsageInfo contains usage information for a tool
type ToolUsageInfo struct {
	Name           string        `json:"name"`
	ExecutionCount int           `json:"execution_count"`
	SuccessRate    float64       `json:"success_rate"`
	AvgLatency     time.Duration `json:"avg_latency"`
	LastUsed       time.Time     `json:"last_used"`
}

// PatternSummary contains summary information about detected patterns
type PatternSummary struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Frequency   int       `json:"frequency"`
	FirstSeen   time.Time `json:"first_seen"`
	LastSeen    time.Time `json:"last_seen"`
}

// InsightSummary contains summary information about generated insights
type InsightSummary struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Priority    string    `json:"priority"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Suggestion  string    `json:"suggestion"`
	CreatedAt   time.Time `json:"created_at"`
}

// GitCommit represents a git commit for changelog generation
type GitCommit struct {
	Hash         string    `json:"hash"`
	ShortHash    string    `json:"short_hash"`
	Author       string    `json:"author"`
	Email        string    `json:"email"`
	Date         time.Time `json:"date"`
	Subject      string    `json:"subject"`
	Body         string    `json:"body"`
	ChangedFiles int       `json:"changed_files"`
	Insertions   int       `json:"insertions"`
	Deletions    int       `json:"deletions"`
}

// ChangelogEntry represents an entry in the changelog
type ChangelogEntry struct {
	Version         string      `json:"version"`
	Date            time.Time   `json:"date"`
	Features        []GitCommit `json:"features"`
	Fixes           []GitCommit `json:"fixes"`
	Performance     []GitCommit `json:"performance"`
	Docs            []GitCommit `json:"docs"`
	Other           []GitCommit `json:"other"`
	BreakingChanges []GitCommit `json:"breaking_changes"`
}

// Generator interface for different document generators
type Generator interface {
	// Generate creates a document of the specified type
	Generate(request GenerationRequest) (*GenerationResult, error)

	// GetSupportedTypes returns the document types this generator supports
	GetSupportedTypes() []DocumentType

	// Validate checks if the generation request is valid for this generator
	Validate(request GenerationRequest) error
}

// DataSource interface for retrieving data for documentation
type DataSource interface {
	// GetCommits retrieves git commits within a date range
	GetCommits(dateRange DateRange) ([]GitCommit, error)

	// GetLearningSnapshot retrieves current learning system data
	GetLearningSnapshot() (*LearningSnapshot, error)

	// GetProjectInfo retrieves general project information
	GetProjectInfo() (map[string]interface{}, error)
}

// DocumentEngine coordinates the generation of various documents
type DocumentEngine interface {
	// RegisterGenerator adds a new document generator
	RegisterGenerator(generator Generator) error

	// Generate creates a document using the appropriate generator
	Generate(request GenerationRequest) (*GenerationResult, error)

	// GenerateAll creates all supported document types
	GenerateAll() ([]GenerationResult, error)

	// ScheduleGeneration sets up automatic document generation
	ScheduleGeneration(docType DocumentType, schedule string) error

	// GetGenerationHistory returns recent generation results
	GetGenerationHistory(limit int) ([]GenerationResult, error)
}
