package selflearn

import (
	"context"
	"time"
)

// Storage defines the interface for storing and retrieving learning data
type Storage interface {
	// Execution records
	StoreExecution(ctx context.Context, record ExecutionRecord) error
	GetExecution(ctx context.Context, id string) (ExecutionRecord, error)
	GetExecutionsByTool(ctx context.Context, toolName string, limit int) ([]ExecutionRecord, error)
	GetExecutionsByTimeRange(ctx context.Context, start, end time.Time, limit int) ([]ExecutionRecord, error)
	GetExecutionStats(ctx context.Context) (LearningStats, error)

	// Patterns
	StorePattern(ctx context.Context, pattern Pattern) error
	GetPattern(ctx context.Context, id string) (Pattern, error)
	GetPatterns(ctx context.Context, patternType PatternType, limit int) ([]Pattern, error)
	UpdatePattern(ctx context.Context, pattern Pattern) error
	DeletePattern(ctx context.Context, id string) error

	// Insights
	StoreInsight(ctx context.Context, insight Insight) error
	GetInsight(ctx context.Context, id string) (Insight, error)
	GetInsights(ctx context.Context, insightType InsightType, limit int) ([]Insight, error)
	GetInsightsByPriority(ctx context.Context, priority Priority, limit int) ([]Insight, error)
	UpdateInsight(ctx context.Context, insight Insight) error
	DeleteInsight(ctx context.Context, id string) error

	// Maintenance
	Cleanup(ctx context.Context, retentionPeriod time.Duration) error
	Close() error
}