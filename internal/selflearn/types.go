package selflearn

import (
	"time"
)

// ExecutionRecord represents a single tool execution with metadata
type ExecutionRecord struct {
	ID           string                 `json:"id"`
	ToolName     string                 `json:"tool_name"`
	Timestamp    time.Time              `json:"timestamp"`
	Duration     time.Duration          `json:"duration"`
	Success      bool                   `json:"success"`
	Input        interface{}            `json:"input,omitempty"`
	Output       interface{}            `json:"output,omitempty"`
	Error        string                 `json:"error,omitempty"`
	ErrorType    ErrorType              `json:"error_type,omitempty"`
	Context      map[string]interface{} `json:"context,omitempty"`
	RetryCount   int                    `json:"retry_count"`
	SourceType   string                 `json:"source_type"` // openapi, graphql, asyncapi, builtin
}

// ErrorType represents the classification of errors
type ErrorType string

const (
	ErrorTypeNetwork       ErrorType = "network"
	ErrorTypeValidation    ErrorType = "validation"
	ErrorTypeConfiguration ErrorType = "configuration"
	ErrorTypePerformance   ErrorType = "performance"
	ErrorTypeLogic         ErrorType = "logic"
	ErrorTypeUnknown       ErrorType = "unknown"
)

// Pattern represents a detected pattern in execution data
type Pattern struct {
	ID          string            `json:"id"`
	Type        PatternType       `json:"type"`
	Description string            `json:"description"`
	Frequency   int               `json:"frequency"`
	Confidence  float64           `json:"confidence"`
	FirstSeen   time.Time         `json:"first_seen"`
	LastSeen    time.Time         `json:"last_seen"`
	Metadata    map[string]string `json:"metadata"`
}

// PatternType represents the type of pattern detected
type PatternType string

const (
	PatternTypeError       PatternType = "error"
	PatternTypePerformance PatternType = "performance"
	PatternTypeUsage       PatternType = "usage"
	PatternTypeSuccess     PatternType = "success"
)

// Insight represents a learning insight or suggestion
type Insight struct {
	ID          string            `json:"id"`
	Type        InsightType       `json:"type"`
	Priority    Priority          `json:"priority"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Suggestion  string            `json:"suggestion"`
	Evidence    []string          `json:"evidence"`
	CreatedAt   time.Time         `json:"created_at"`
	Metadata    map[string]string `json:"metadata"`
}

// InsightType represents the type of insight
type InsightType string

const (
	InsightTypeOptimization    InsightType = "optimization"
	InsightTypeConfiguration   InsightType = "configuration"
	InsightTypeReliability     InsightType = "reliability"
	InsightTypePerformance     InsightType = "performance"
	InsightTypeUsage          InsightType = "usage"
)

// Priority represents the priority level of an insight
type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityMedium   Priority = "medium"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

// LearningStats represents overall learning statistics
type LearningStats struct {
	TotalExecutions   int64             `json:"total_executions"`
	SuccessRate       float64           `json:"success_rate"`
	AverageLatency    time.Duration     `json:"average_latency"`
	ErrorBreakdown    map[ErrorType]int `json:"error_breakdown"`
	TopTools          []ToolStat        `json:"top_tools"`
	RecentPatterns    []Pattern         `json:"recent_patterns"`
	ActiveInsights    []Insight         `json:"active_insights"`
	LastUpdated       time.Time         `json:"last_updated"`
}

// ToolStat represents statistics for a specific tool
type ToolStat struct {
	Name           string        `json:"name"`
	ExecutionCount int64         `json:"execution_count"`
	SuccessRate    float64       `json:"success_rate"`
	AverageLatency time.Duration `json:"average_latency"`
	LastUsed       time.Time     `json:"last_used"`
}

// CollectionConfig represents configuration for feedback collection
type CollectionConfig struct {
	Enabled              bool          `json:"enabled"`
	SampleRate           float64       `json:"sample_rate"`           // 0.0 to 1.0
	MaxInputSize         int           `json:"max_input_size"`        // bytes
	MaxOutputSize        int           `json:"max_output_size"`       // bytes
	RetentionPeriod      time.Duration `json:"retention_period"`     // how long to keep records
	PIIFilterEnabled     bool          `json:"pii_filter_enabled"`   // filter out PII data
	AsyncProcessing      bool          `json:"async_processing"`     // process feedback asynchronously
	IncludeSuccessful    bool          `json:"include_successful"`   // collect data for successful executions
	IncludeInputOutput   bool          `json:"include_input_output"` // include actual input/output data
}

// DefaultCollectionConfig returns a sensible default configuration
func DefaultCollectionConfig() CollectionConfig {
	return CollectionConfig{
		Enabled:              true,
		SampleRate:           1.0, // collect all executions by default
		MaxInputSize:         1024,
		MaxOutputSize:        4096,
		RetentionPeriod:      30 * 24 * time.Hour, // 30 days
		PIIFilterEnabled:     true,
		AsyncProcessing:      true,
		IncludeSuccessful:    true,
		IncludeInputOutput:   true,
	}
}