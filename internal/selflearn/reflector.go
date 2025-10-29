package selflearn

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// Reflector generates insights and suggestions based on patterns and data
type Reflector struct {
	storage  Storage
	analyzer *Analyzer
	logger   *zap.Logger
}

// NewReflector creates a new insight reflector
func NewReflector(storage Storage, analyzer *Analyzer, logger *zap.Logger) *Reflector {
	return &Reflector{
		storage:  storage,
		analyzer: analyzer,
		logger:   logger,
	}
}

// GenerateInsights generates actionable insights based on current patterns and data
func (r *Reflector) GenerateInsights(ctx context.Context) ([]Insight, error) {
	r.logger.Info("Starting insight generation")

	var insights []Insight

	// Generate insights from error patterns
	errorInsights, err := r.generateErrorInsights(ctx)
	if err != nil {
		r.logger.Error("Failed to generate error insights", zap.Error(err))
	} else {
		insights = append(insights, errorInsights...)
	}

	// Generate insights from performance patterns
	perfInsights, err := r.generatePerformanceInsights(ctx)
	if err != nil {
		r.logger.Error("Failed to generate performance insights", zap.Error(err))
	} else {
		insights = append(insights, perfInsights...)
	}

	// Generate insights from usage patterns
	usageInsights, err := r.generateUsageInsights(ctx)
	if err != nil {
		r.logger.Error("Failed to generate usage insights", zap.Error(err))
	} else {
		insights = append(insights, usageInsights...)
	}

	// Generate configuration insights
	configInsights, err := r.generateConfigurationInsights(ctx)
	if err != nil {
		r.logger.Error("Failed to generate configuration insights", zap.Error(err))
	} else {
		insights = append(insights, configInsights...)
	}

	// Store generated insights
	for _, insight := range insights {
		if err := r.storage.StoreInsight(ctx, insight); err != nil {
			r.logger.Error("Failed to store insight",
				zap.String("insight_id", insight.ID),
				zap.Error(err))
		}
	}

	r.logger.Info("Insight generation completed", zap.Int("insights_generated", len(insights)))
	return insights, nil
}

// generateErrorInsights creates insights based on error patterns
func (r *Reflector) generateErrorInsights(ctx context.Context) ([]Insight, error) {
	patterns, err := r.storage.GetPatterns(ctx, PatternTypeError, 50)
	if err != nil {
		return nil, fmt.Errorf("failed to get error patterns: %w", err)
	}

	var insights []Insight

	for _, pattern := range patterns {
		if pattern.Confidence < 0.6 {
			continue // Skip low-confidence patterns
		}

		var priority Priority
		var suggestion string

		switch {
		case pattern.Frequency >= 50:
			priority = PriorityCritical
			suggestion = fmt.Sprintf("Immediate attention required: %s errors occur very frequently (%d times). Consider reviewing the tool configuration, endpoint availability, or implementing retry logic.", 
				pattern.Metadata["error_type"], pattern.Frequency)
		case pattern.Frequency >= 20:
			priority = PriorityHigh
			suggestion = fmt.Sprintf("High priority: %s errors in %s tool need investigation. Review error logs and consider implementing error handling improvements.",
				pattern.Metadata["error_type"], pattern.Metadata["tool_name"])
		case pattern.Frequency >= 10:
			priority = PriorityMedium
			suggestion = fmt.Sprintf("Monitor %s tool for recurring %s errors. Consider implementing better error messaging or user guidance.",
				pattern.Metadata["tool_name"], pattern.Metadata["error_type"])
		default:
			priority = PriorityLow
			suggestion = fmt.Sprintf("Track %s errors in %s tool. May indicate edge cases or rare scenarios.",
				pattern.Metadata["error_type"], pattern.Metadata["tool_name"])
		}

		insight := Insight{
			ID:          r.generateInsightID(),
			Type:        InsightTypeReliability,
			Priority:    priority,
			Title:       fmt.Sprintf("Recurring %s Errors in %s", pattern.Metadata["error_type"], pattern.Metadata["tool_name"]),
			Description: fmt.Sprintf("Pattern detected: %s (Confidence: %.1f%%)", pattern.Description, pattern.Confidence*100),
			Suggestion:  suggestion,
			Evidence: []string{
				fmt.Sprintf("Error frequency: %d occurrences", pattern.Frequency),
				fmt.Sprintf("Pattern confidence: %.1f%%", pattern.Confidence*100),
				fmt.Sprintf("Time range: %s to %s", pattern.FirstSeen.Format("2006-01-02"), pattern.LastSeen.Format("2006-01-02")),
			},
			CreatedAt: time.Now().UTC(),
			Metadata: map[string]string{
				"tool_name":   pattern.Metadata["tool_name"],
				"error_type":  pattern.Metadata["error_type"],
				"pattern_id":  pattern.ID,
				"source_type": "error_pattern",
			},
		}

		insights = append(insights, insight)
	}

	return insights, nil
}

// generatePerformanceInsights creates insights based on performance patterns
func (r *Reflector) generatePerformanceInsights(ctx context.Context) ([]Insight, error) {
	patterns, err := r.storage.GetPatterns(ctx, PatternTypePerformance, 20)
	if err != nil {
		return nil, fmt.Errorf("failed to get performance patterns: %w", err)
	}

	var insights []Insight

	for _, pattern := range patterns {
		priority := PriorityMedium
		if pattern.Frequency >= 100 {
			priority = PriorityHigh
		}

		suggestion := fmt.Sprintf("Performance optimization needed for %s tool. Consider implementing caching, optimizing API calls, or adding timeout configurations. Average latency: %s",
			pattern.Metadata["tool_name"], pattern.Metadata["average_latency"])

		insight := Insight{
			ID:          r.generateInsightID(),
			Type:        InsightTypePerformance,
			Priority:    priority,
			Title:       fmt.Sprintf("Performance Issues in %s Tool", pattern.Metadata["tool_name"]),
			Description: fmt.Sprintf("Tool shows consistently slow performance: %s", pattern.Description),
			Suggestion:  suggestion,
			Evidence: []string{
				fmt.Sprintf("Average latency: %s", pattern.Metadata["average_latency"]),
				fmt.Sprintf("Execution count: %s", pattern.Metadata["execution_count"]),
				fmt.Sprintf("Success rate: %s%%", pattern.Metadata["success_rate"]),
			},
			CreatedAt: time.Now().UTC(),
			Metadata: map[string]string{
				"tool_name":       pattern.Metadata["tool_name"],
				"average_latency": pattern.Metadata["average_latency"],
				"pattern_id":      pattern.ID,
				"source_type":     "performance_pattern",
			},
		}

		insights = append(insights, insight)
	}

	return insights, nil
}

// generateUsageInsights creates insights based on usage patterns
func (r *Reflector) generateUsageInsights(ctx context.Context) ([]Insight, error) {
	patterns, err := r.storage.GetPatterns(ctx, PatternTypeUsage, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage patterns: %w", err)
	}

	var insights []Insight

	for _, pattern := range patterns {
		suggestion := fmt.Sprintf("Consider optimizing the %s tool for better performance since it represents a significant portion of usage. Also consider load balancing or caching strategies.",
			pattern.Metadata["tool_name"])

		insight := Insight{
			ID:          r.generateInsightID(),
			Type:        InsightTypeUsage,
			Priority:    PriorityMedium,
			Title:       fmt.Sprintf("High Usage Pattern: %s Tool", pattern.Metadata["tool_name"]),
			Description: pattern.Description,
			Suggestion:  suggestion,
			Evidence: []string{
				fmt.Sprintf("Usage percentage: %s%%", pattern.Metadata["usage_percentage"]),
				fmt.Sprintf("Total executions: %s", pattern.Metadata["execution_count"]),
			},
			CreatedAt: time.Now().UTC(),
			Metadata: map[string]string{
				"tool_name":        pattern.Metadata["tool_name"],
				"usage_percentage": pattern.Metadata["usage_percentage"],
				"pattern_id":       pattern.ID,
				"source_type":      "usage_pattern",
			},
		}

		insights = append(insights, insight)
	}

	return insights, nil
}

// generateConfigurationInsights creates insights based on system configuration
func (r *Reflector) generateConfigurationInsights(ctx context.Context) ([]Insight, error) {
	stats, err := r.storage.GetExecutionStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get execution stats: %w", err)
	}

	var insights []Insight

	// Check overall success rate
	if stats.SuccessRate < 0.8 && stats.TotalExecutions > 100 {
		priority := PriorityHigh
		if stats.SuccessRate < 0.5 {
			priority = PriorityCritical
		}

		insight := Insight{
			ID:          r.generateInsightID(),
			Type:        InsightTypeConfiguration,
			Priority:    priority,
			Title:       "Low Overall Success Rate",
			Description: fmt.Sprintf("System-wide success rate is %.1f%%, indicating potential configuration issues", stats.SuccessRate*100),
			Suggestion:  "Review system configuration, endpoint URLs, authentication settings, and network connectivity. Consider implementing health checks and monitoring.",
			Evidence: []string{
				fmt.Sprintf("Success rate: %.1f%%", stats.SuccessRate*100),
				fmt.Sprintf("Total executions: %d", stats.TotalExecutions),
				fmt.Sprintf("Error breakdown available for detailed analysis"),
			},
			CreatedAt: time.Now().UTC(),
			Metadata: map[string]string{
				"success_rate":      fmt.Sprintf("%.2f", stats.SuccessRate),
				"total_executions":  fmt.Sprintf("%d", stats.TotalExecutions),
				"source_type":       "system_stats",
			},
		}

		insights = append(insights, insight)
	}

	// Check for network-heavy error patterns
	if networkErrors, exists := stats.ErrorBreakdown[ErrorTypeNetwork]; exists && networkErrors > int(stats.TotalExecutions)/10 {
		insight := Insight{
			ID:          r.generateInsightID(),
			Type:        InsightTypeConfiguration,
			Priority:    PriorityHigh,
			Title:       "High Network Error Rate",
			Description: fmt.Sprintf("Network errors account for %d of %d total executions, suggesting connectivity issues", networkErrors, stats.TotalExecutions),
			Suggestion:  "Review network configuration, implement retry logic with exponential backoff, and consider circuit breaker patterns for external API calls.",
			Evidence: []string{
				fmt.Sprintf("Network errors: %d", networkErrors),
				fmt.Sprintf("Error percentage: %.1f%%", float64(networkErrors)/float64(stats.TotalExecutions)*100),
			},
			CreatedAt: time.Now().UTC(),
			Metadata: map[string]string{
				"network_errors":   fmt.Sprintf("%d", networkErrors),
				"total_executions": fmt.Sprintf("%d", stats.TotalExecutions),
				"source_type":      "error_analysis",
			},
		}

		insights = append(insights, insight)
	}

	return insights, nil
}

// generateInsightID generates a unique ID for insights
func (r *Reflector) generateInsightID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		zap.L().Error("failed to generate random bytes for insight ID", zap.Error(err))
		return fmt.Sprintf("insight_fallback_%d", time.Now().UnixNano())
	}
	return "insight_" + hex.EncodeToString(bytes)
}