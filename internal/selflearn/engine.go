package selflearn

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// Engine is the main self-learning engine that coordinates feedback collection,
// analysis, and insight generation
type Engine struct {
	collector *Collector
	storage   Storage
	analyzer  *Analyzer
	reflector *Reflector
	config    CollectionConfig
	logger    *zap.Logger
}

// NewEngine creates a new self-learning engine
func NewEngine(config CollectionConfig, storage Storage, logger *zap.Logger) *Engine {
	collector := NewCollector(config, storage, logger)
	analyzer := NewAnalyzer(storage, logger)
	reflector := NewReflector(storage, analyzer, logger)

	return &Engine{
		collector: collector,
		storage:   storage,
		analyzer:  analyzer,
		reflector: reflector,
		config:    config,
		logger:    logger,
	}
}

// RecordExecution records the execution of a tool for learning purposes
func (e *Engine) RecordExecution(ctx context.Context, toolName, sourceType string, input interface{}, output interface{}, err error, duration time.Duration) error {
	execCtx := ExecutionContext{
		ToolName:   toolName,
		SourceType: sourceType,
		Metadata:   make(map[string]interface{}),
	}

	// Extract additional context from the context if available
	if sessionID := ctx.Value("session_id"); sessionID != nil {
		if sid, ok := sessionID.(string); ok {
			execCtx.SessionID = sid
		}
	}
	if requestID := ctx.Value("request_id"); requestID != nil {
		if rid, ok := requestID.(string); ok {
			execCtx.RequestID = rid
		}
	}
	if userAgent := ctx.Value("user_agent"); userAgent != nil {
		if ua, ok := userAgent.(string); ok {
			execCtx.UserAgent = ua
		}
	}

	return e.collector.CollectExecution(ctx, execCtx, input, output, err, duration)
}

// AnalyzePatterns triggers pattern analysis on existing execution data
func (e *Engine) AnalyzePatterns(ctx context.Context) ([]Pattern, error) {
	return e.analyzer.AnalyzePatterns(ctx)
}

// GenerateInsights triggers insight generation based on current patterns and data
func (e *Engine) GenerateInsights(ctx context.Context) ([]Insight, error) {
	return e.reflector.GenerateInsights(ctx)
}

// GetStats returns overall learning statistics
func (e *Engine) GetStats(ctx context.Context) (LearningStats, error) {
	stats, err := e.storage.GetExecutionStats(ctx)
	if err != nil {
		return stats, err
	}

	// Enhance stats with recent patterns and insights
	patterns, err := e.storage.GetPatterns(ctx, "", 5)
	if err != nil {
		e.logger.Warn("Failed to get recent patterns", zap.Error(err))
	} else {
		stats.RecentPatterns = patterns
	}

	insights, err := e.storage.GetInsights(ctx, "", 10)
	if err != nil {
		e.logger.Warn("Failed to get active insights", zap.Error(err))
	} else {
		stats.ActiveInsights = insights
	}

	return stats, nil
}

// GetToolInsights returns insights specific to a tool
func (e *Engine) GetToolInsights(ctx context.Context, toolName string) ([]Insight, error) {
	// Get all insights and filter by tool
	insights, err := e.storage.GetInsights(ctx, "", 100)
	if err != nil {
		return nil, err
	}

	var toolInsights []Insight
	for _, insight := range insights {
		if toolName, exists := insight.Metadata["tool_name"]; exists && toolName == toolName {
			toolInsights = append(toolInsights, insight)
		}
	}

	return toolInsights, nil
}

// GetErrorPatterns returns error patterns, optionally filtered by tool
func (e *Engine) GetErrorPatterns(ctx context.Context, toolName string) ([]Pattern, error) {
	patterns, err := e.storage.GetPatterns(ctx, PatternTypeError, 50)
	if err != nil {
		return nil, err
	}

	if toolName == "" {
		return patterns, nil
	}

	// Filter by tool name
	var toolPatterns []Pattern
	for _, pattern := range patterns {
		if tool, exists := pattern.Metadata["tool_name"]; exists && tool == toolName {
			toolPatterns = append(toolPatterns, pattern)
		}
	}

	return toolPatterns, nil
}

// RunMaintenance performs maintenance tasks like cleanup and analysis
func (e *Engine) RunMaintenance(ctx context.Context) error {
	e.logger.Info("Starting self-learning maintenance")

	// Cleanup old data
	if err := e.storage.Cleanup(ctx, e.config.RetentionPeriod); err != nil {
		e.logger.Error("Failed to cleanup old data", zap.Error(err))
	}

	// Run pattern analysis
	patterns, err := e.analyzer.AnalyzePatterns(ctx)
	if err != nil {
		e.logger.Error("Failed to analyze patterns", zap.Error(err))
	} else {
		e.logger.Info("Pattern analysis completed", zap.Int("patterns_found", len(patterns)))
	}

	// Generate insights
	insights, err := e.reflector.GenerateInsights(ctx)
	if err != nil {
		e.logger.Error("Failed to generate insights", zap.Error(err))
	} else {
		e.logger.Info("Insight generation completed", zap.Int("insights_generated", len(insights)))
	}

	e.logger.Info("Self-learning maintenance completed")
	return nil
}

// UpdateConfig updates the engine configuration
func (e *Engine) UpdateConfig(config CollectionConfig) {
	e.config = config
	e.collector.UpdateConfig(config)
	e.logger.Info("Engine configuration updated")
}

// GetConfig returns the current engine configuration
func (e *Engine) GetConfig() CollectionConfig {
	return e.config
}

// Close shuts down the learning engine
func (e *Engine) Close() error {
	e.logger.Info("Shutting down self-learning engine")
	return e.storage.Close()
}