package feedback

import (
	"time"
)

// Models for the public feedback API

// ExecutionSummary represents a summarized view of execution data
type ExecutionSummary struct {
	ID           string                 `json:"id"`
	ToolName     string                 `json:"tool_name"`
	Timestamp    time.Time              `json:"timestamp"`
	Duration     time.Duration          `json:"duration"`
	Success      bool                   `json:"success"`
	ErrorType    string                 `json:"error_type,omitempty"`
	SourceType   string                 `json:"source_type"`
	Context      map[string]interface{} `json:"context,omitempty"`
}

// PatternSummary represents a summarized view of detected patterns
type PatternSummary struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Frequency   int       `json:"frequency"`
	Confidence  float64   `json:"confidence"`
	FirstSeen   time.Time `json:"first_seen"`
	LastSeen    time.Time `json:"last_seen"`
}

// InsightSummary represents a summarized view of generated insights
type InsightSummary struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Priority    string    `json:"priority"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Suggestion  string    `json:"suggestion"`
	CreatedAt   time.Time `json:"created_at"`
}

// LearningOverview provides a high-level view of learning data
type LearningOverview struct {
	TotalExecutions    int64                    `json:"total_executions"`
	SuccessRate        float64                  `json:"success_rate"`
	AverageLatency     time.Duration            `json:"average_latency"`
	ErrorBreakdown     map[string]int           `json:"error_breakdown"`
	TopTools           []ToolStatSummary        `json:"top_tools"`
	RecentPatterns     []PatternSummary         `json:"recent_patterns"`
	ActiveInsights     []InsightSummary         `json:"active_insights"`
	LastUpdated        time.Time                `json:"last_updated"`
}

// ToolStatSummary represents statistics for a specific tool
type ToolStatSummary struct {
	Name           string        `json:"name"`
	ExecutionCount int64         `json:"execution_count"`
	SuccessRate    float64       `json:"success_rate"`
	AverageLatency time.Duration `json:"average_latency"`
	LastUsed       time.Time     `json:"last_used"`
}

// ToolInsights represents insights and patterns for a specific tool
type ToolInsights struct {
	ToolName      string           `json:"tool_name"`
	Stats         ToolStatSummary  `json:"stats"`
	Patterns      []PatternSummary `json:"patterns"`
	Insights      []InsightSummary `json:"insights"`
	Suggestions   []string         `json:"suggestions"`
}

// ErrorAnalysis provides detailed error analysis
type ErrorAnalysis struct {
	ErrorType      string               `json:"error_type"`
	TotalCount     int                  `json:"total_count"`
	Percentage     float64              `json:"percentage"`
	TopTools       []ToolErrorSummary   `json:"top_tools"`
	RecentTrend    string               `json:"recent_trend"` // "increasing", "decreasing", "stable"
	Suggestions    []string             `json:"suggestions"`
}

// ToolErrorSummary represents error statistics for a tool
type ToolErrorSummary struct {
	ToolName   string  `json:"tool_name"`
	ErrorCount int     `json:"error_count"`
	Percentage float64 `json:"percentage"`
}

// PerformanceAnalysis provides performance insights
type PerformanceAnalysis struct {
	OverallLatency    time.Duration              `json:"overall_latency"`
	SlowestTools      []ToolPerformanceSummary   `json:"slowest_tools"`
	FastestTools      []ToolPerformanceSummary   `json:"fastest_tools"`
	LatencyTrend      string                     `json:"latency_trend"` // "improving", "degrading", "stable"
	Recommendations   []string                   `json:"recommendations"`
}

// ToolPerformanceSummary represents performance statistics for a tool
type ToolPerformanceSummary struct {
	ToolName        string        `json:"tool_name"`
	AverageLatency  time.Duration `json:"average_latency"`
	MedianLatency   time.Duration `json:"median_latency,omitempty"`
	P95Latency      time.Duration `json:"p95_latency,omitempty"`
	ExecutionCount  int64         `json:"execution_count"`
}

// UsageAnalysis provides usage pattern insights
type UsageAnalysis struct {
	TotalExecutions   int64                    `json:"total_executions"`
	UniqueSessions    int64                    `json:"unique_sessions,omitempty"`
	TopTools          []ToolUsageSummary       `json:"top_tools"`
	UsageDistribution map[string]float64       `json:"usage_distribution"` // tool -> percentage
	PeakUsageHours    []int                    `json:"peak_usage_hours"`   // hours of day (0-23)
	Insights          []string                 `json:"insights"`
}

// ToolUsageSummary represents usage statistics for a tool
type ToolUsageSummary struct {
	ToolName       string    `json:"tool_name"`
	ExecutionCount int64     `json:"execution_count"`
	Percentage     float64   `json:"percentage"`
	LastUsed       time.Time `json:"last_used"`
	FirstUsed      time.Time `json:"first_used"`
}

// LearningConfig represents the current learning configuration
type LearningConfig struct {
	Enabled              bool          `json:"enabled"`
	SampleRate           float64       `json:"sample_rate"`
	RetentionPeriod      time.Duration `json:"retention_period"`
	AsyncProcessing      bool          `json:"async_processing"`
	IncludeSuccessful    bool          `json:"include_successful"`
	IncludeInputOutput   bool          `json:"include_input_output"`
	PIIFilterEnabled     bool          `json:"pii_filter_enabled"`
	MaxInputSize         int           `json:"max_input_size"`
	MaxOutputSize        int           `json:"max_output_size"`
}