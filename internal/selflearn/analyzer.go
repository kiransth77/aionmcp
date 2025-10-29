package selflearn

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// Analyzer performs pattern analysis on execution data
type Analyzer struct {
	storage Storage
	logger  *zap.Logger
}

// NewAnalyzer creates a new pattern analyzer
func NewAnalyzer(storage Storage, logger *zap.Logger) *Analyzer {
	return &Analyzer{
		storage: storage,
		logger:  logger,
	}
}

// AnalyzePatterns analyzes execution data to identify patterns
func (a *Analyzer) AnalyzePatterns(ctx context.Context) ([]Pattern, error) {
	a.logger.Info("Starting pattern analysis")

	var patterns []Pattern

	// Analyze error patterns
	errorPatterns, err := a.analyzeErrorPatterns(ctx)
	if err != nil {
		a.logger.Error("Failed to analyze error patterns", zap.Error(err))
	} else {
		patterns = append(patterns, errorPatterns...)
	}

	// Analyze performance patterns
	perfPatterns, err := a.analyzePerformancePatterns(ctx)
	if err != nil {
		a.logger.Error("Failed to analyze performance patterns", zap.Error(err))
	} else {
		patterns = append(patterns, perfPatterns...)
	}

	// Analyze usage patterns
	usagePatterns, err := a.analyzeUsagePatterns(ctx)
	if err != nil {
		a.logger.Error("Failed to analyze usage patterns", zap.Error(err))
	} else {
		patterns = append(patterns, usagePatterns...)
	}

	// Store discovered patterns
	for _, pattern := range patterns {
		if err := a.storage.StorePattern(ctx, pattern); err != nil {
			a.logger.Error("Failed to store pattern", 
				zap.String("pattern_id", pattern.ID),
				zap.Error(err))
		}
	}

	a.logger.Info("Pattern analysis completed", zap.Int("patterns_found", len(patterns)))
	return patterns, nil
}

// analyzeErrorPatterns identifies common error patterns
func (a *Analyzer) analyzeErrorPatterns(ctx context.Context) ([]Pattern, error) {
	// Get recent executions with errors
	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour) // Last 24 hours
	
	executions, err := a.storage.GetExecutionsByTimeRange(ctx, startTime, endTime, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to get executions: %w", err)
	}

	// Group errors by type and tool
	errorGroups := make(map[string]*errorGroup)
	
	for _, exec := range executions {
		if exec.Success {
			continue
		}

		key := fmt.Sprintf("%s_%s", exec.ToolName, exec.ErrorType)
		if group, exists := errorGroups[key]; exists {
			group.count++
			group.lastSeen = exec.Timestamp
			if exec.Timestamp.Before(group.firstSeen) {
				group.firstSeen = exec.Timestamp
			}
			// Track unique error messages
			group.errorMessages[exec.Error] = true
		} else {
			errorGroups[key] = &errorGroup{
				toolName:      exec.ToolName,
				errorType:     exec.ErrorType,
				count:         1,
				firstSeen:     exec.Timestamp,
				lastSeen:      exec.Timestamp,
				errorMessages: map[string]bool{exec.Error: true},
			}
		}
	}

	var patterns []Pattern
	
	// Convert significant error groups to patterns
	for _, group := range errorGroups {
		if group.count >= 3 { // Threshold for pattern recognition
			pattern := Pattern{
				ID:          a.generatePatternID(),
				Type:        PatternTypeError,
				Description: fmt.Sprintf("Recurring %s errors in %s tool", group.errorType, group.toolName),
				Frequency:   group.count,
				Confidence:  a.calculateConfidence(group.count, len(executions)),
				FirstSeen:   group.firstSeen,
				LastSeen:    group.lastSeen,
				Metadata: map[string]string{
					"tool_name":       group.toolName,
					"error_type":      string(group.errorType),
					"unique_messages": fmt.Sprintf("%d", len(group.errorMessages)),
				},
			}
			patterns = append(patterns, pattern)
		}
	}

	return patterns, nil
}

// analyzePerformancePatterns identifies performance-related patterns
func (a *Analyzer) analyzePerformancePatterns(ctx context.Context) ([]Pattern, error) {
	stats, err := a.storage.GetExecutionStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get execution stats: %w", err)
	}

	var patterns []Pattern

	// Check for slow tools (above average + 2 standard deviations)
	avgLatency := stats.AverageLatency
	for _, toolStat := range stats.TopTools {
		if toolStat.AverageLatency > avgLatency*2 {
			pattern := Pattern{
				ID:          a.generatePatternID(),
				Type:        PatternTypePerformance,
				Description: fmt.Sprintf("Tool %s shows consistently slow performance", toolStat.Name),
				Frequency:   int(toolStat.ExecutionCount),
				Confidence:  0.8, // High confidence for performance metrics
				FirstSeen:   time.Now().Add(-7 * 24 * time.Hour), // Approximate
				LastSeen:    toolStat.LastUsed,
				Metadata: map[string]string{
					"tool_name":        toolStat.Name,
					"average_latency":  toolStat.AverageLatency.String(),
					"execution_count":  fmt.Sprintf("%d", toolStat.ExecutionCount),
					"success_rate":     fmt.Sprintf("%.2f", toolStat.SuccessRate),
				},
			}
			patterns = append(patterns, pattern)
		}
	}

	return patterns, nil
}

// analyzeUsagePatterns identifies usage-related patterns
func (a *Analyzer) analyzeUsagePatterns(ctx context.Context) ([]Pattern, error) {
	stats, err := a.storage.GetExecutionStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get execution stats: %w", err)
	}

	var patterns []Pattern

	// Identify heavily used tools
	if len(stats.TopTools) > 0 && stats.TotalExecutions > 0 {
		topTool := stats.TopTools[0]
		usagePercentage := float64(topTool.ExecutionCount) / float64(stats.TotalExecutions) * 100

		if usagePercentage > 50 { // More than 50% of executions
			pattern := Pattern{
				ID:          a.generatePatternID(),
				Type:        PatternTypeUsage,
				Description: fmt.Sprintf("Tool %s dominates usage with %.1f%% of all executions", topTool.Name, usagePercentage),
				Frequency:   int(topTool.ExecutionCount),
				Confidence:  0.9,
				FirstSeen:   time.Now().Add(-30 * 24 * time.Hour), // Approximate
				LastSeen:    topTool.LastUsed,
				Metadata: map[string]string{
					"tool_name":        topTool.Name,
					"usage_percentage": fmt.Sprintf("%.1f", usagePercentage),
					"execution_count":  fmt.Sprintf("%d", topTool.ExecutionCount),
				},
			}
			patterns = append(patterns, pattern)
		}
	}

	return patterns, nil
}

// calculateConfidence calculates confidence score for a pattern
func (a *Analyzer) calculateConfidence(frequency, totalSamples int) float64 {
	if totalSamples == 0 {
		return 0.0
	}

	// Simple confidence calculation based on frequency and sample size
	ratio := float64(frequency) / float64(totalSamples)
	
	// Base confidence on ratio and sample size
	confidence := ratio
	if frequency >= 10 {
		confidence += 0.2 // Boost for higher frequency
	}
	if frequency >= 50 {
		confidence += 0.1 // Additional boost for very high frequency
	}

	// Cap at 1.0
	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}

// generatePatternID generates a unique ID for patterns
func (a *Analyzer) generatePatternID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		a.logger.Error("failed to generate random pattern ID", zap.Error(err))
		// Fallback: use timestamp-based ID
		return fmt.Sprintf("pattern_fallback_%d", time.Now().UnixNano())
	}
	return "pattern_" + hex.EncodeToString(bytes)
}

// errorGroup represents a group of similar errors
type errorGroup struct {
	toolName      string
	errorType     ErrorType
	count         int
	firstSeen     time.Time
	lastSeen      time.Time
	errorMessages map[string]bool
}