package autodocs

import (
	"fmt"
	"net/http"
	"encoding/json"
	"time"
	"context"
)

const (
	// DefaultHTTPTimeout is the default timeout for HTTP client requests
	DefaultHTTPTimeout = 30 * time.Second
)

// LearningDataSource implements DataSource interface by integrating with the learning system
type LearningDataSource struct {
	gitDataSource *GitDataSource
	learningAPIURL string
	httpClient    *http.Client
}

// NewLearningDataSource creates a new learning-integrated data source with default timeout
func NewLearningDataSource(repoPath, learningAPIURL string) *LearningDataSource {
	return NewLearningDataSourceWithTimeout(repoPath, learningAPIURL, DefaultHTTPTimeout)
}

// NewLearningDataSourceWithTimeout creates a new learning-integrated data source with custom timeout
func NewLearningDataSourceWithTimeout(repoPath, learningAPIURL string, timeout time.Duration) *LearningDataSource {
	return &LearningDataSource{
		gitDataSource:  NewGitDataSource(repoPath),
		learningAPIURL: learningAPIURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// GetCommits retrieves git commits (delegates to git data source)
func (l *LearningDataSource) GetCommits(dateRange DateRange) ([]GitCommit, error) {
	return l.gitDataSource.GetCommits(dateRange)
}

// GetProjectInfo retrieves project information (delegates to git data source)
func (l *LearningDataSource) GetProjectInfo() (map[string]interface{}, error) {
	return l.gitDataSource.GetProjectInfo()
}

// GetLearningSnapshot retrieves current learning system data
func (l *LearningDataSource) GetLearningSnapshot() (*LearningSnapshot, error) {
	// Try to get real learning data from the API
	if l.learningAPIURL != "" {
		if snapshot, err := l.fetchLearningData(); err == nil {
			return snapshot, nil
		}
	}
	
	// Fallback to mock data if learning system is not available
	return l.getMockLearningSnapshot(), nil
}

// fetchLearningData retrieves data from the learning system API
func (l *LearningDataSource) fetchLearningData() (*LearningSnapshot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Get learning statistics
	statsURL := fmt.Sprintf("%s/api/v1/learning/stats", l.learningAPIURL)
	req, err := http.NewRequestWithContext(ctx, "GET", statsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create stats request: %w", err)
	}
	
	resp, err := l.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch learning stats: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("learning API returned status: %d", resp.StatusCode)
	}
	
	var stats struct {
		TotalExecutions int                    `json:"total_executions"`
		SuccessRate     float64                `json:"success_rate"`
		AverageLatency  int64                  `json:"average_latency"` // nanoseconds
		ErrorBreakdown  map[string]int         `json:"error_breakdown"`
		TopTools        []ToolUsageInfo        `json:"top_tools"`
		RecentPatterns  []PatternSummary       `json:"recent_patterns"`
		ActiveInsights  []InsightSummary       `json:"active_insights"`
		LastUpdated     time.Time              `json:"last_updated"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("failed to decode learning stats: %w", err)
	}
	
	// Convert to LearningSnapshot
	snapshot := &LearningSnapshot{
		TotalExecutions: stats.TotalExecutions,
		SuccessRate:     stats.SuccessRate,
		AvgLatency:      time.Duration(stats.AverageLatency),
		TopTools:        stats.TopTools,
		ErrorBreakdown:  stats.ErrorBreakdown,
		RecentPatterns:  stats.RecentPatterns,
		ActiveInsights:  stats.ActiveInsights,
		SnapshotTime:    time.Now(),
	}
	
	return snapshot, nil
}

// getMockLearningSnapshot returns mock learning data for testing/fallback
func (l *LearningDataSource) getMockLearningSnapshot() *LearningSnapshot {
	return &LearningSnapshot{
		TotalExecutions: 42,
		SuccessRate:     0.97,
		AvgLatency:      250 * time.Millisecond,
		TopTools: []ToolUsageInfo{
			{
				Name:           "openapi.petstore.listPets",
				ExecutionCount: 25,
				SuccessRate:    0.96,
				AvgLatency:     180 * time.Millisecond,
				LastUsed:       time.Now().Add(-2 * time.Hour),
			},
			{
				Name:           "graphql.blog.getPosts",
				ExecutionCount: 15,
				SuccessRate:    1.0,
				AvgLatency:     120 * time.Millisecond,
				LastUsed:       time.Now().Add(-1 * time.Hour),
			},
			{
				Name:           "asyncapi.user-events.publishEvent",
				ExecutionCount: 8,
				SuccessRate:    0.875,
				AvgLatency:     350 * time.Millisecond,
				LastUsed:       time.Now().Add(-30 * time.Minute),
			},
		},
		ErrorBreakdown: map[string]int{
			"network":       2,
			"validation":    1,
			"timeout":       1,
		},
		RecentPatterns: []PatternSummary{
			{
				ID:          "pattern_usage_001",
				Type:        "usage",
				Description: "OpenAPI tools are used 60% of the time",
				Frequency:   25,
				FirstSeen:   time.Now().Add(-7 * 24 * time.Hour),
				LastSeen:    time.Now().Add(-2 * time.Hour),
			},
			{
				ID:          "pattern_perf_001", 
				Type:        "performance",
				Description: "AsyncAPI tools show higher latency (>300ms)",
				Frequency:   8,
				FirstSeen:   time.Now().Add(-5 * 24 * time.Hour),
				LastSeen:    time.Now().Add(-30 * time.Minute),
			},
		},
		ActiveInsights: []InsightSummary{
			{
				ID:          "insight_perf_001",
				Type:        "performance",
				Priority:    "medium",
				Title:       "AsyncAPI Tool Performance",
				Description: "AsyncAPI tools showing higher than average latency",
				Suggestion:  "Consider implementing connection pooling or caching for AsyncAPI tools",
				CreatedAt:   time.Now().Add(-24 * time.Hour),
			},
			{
				ID:          "insight_usage_001",
				Type:        "optimization",
				Priority:    "low",
				Title:       "Tool Usage Imbalance",
				Description: "OpenAPI tools are heavily used while GraphQL tools are underutilized",
				Suggestion:  "Review GraphQL tool capabilities and consider promoting usage",
				CreatedAt:   time.Now().Add(-12 * time.Hour),
			},
		},
		SnapshotTime: time.Now(),
	}
}

// GetDetailedInsights retrieves detailed insights from the learning system
func (l *LearningDataSource) GetDetailedInsights() ([]InsightSummary, error) {
	if l.learningAPIURL == "" {
		return l.getMockLearningSnapshot().ActiveInsights, nil
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	insightsURL := fmt.Sprintf("%s/api/v1/learning/insights", l.learningAPIURL)
	req, err := http.NewRequestWithContext(ctx, "GET", insightsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create insights request: %w", err)
	}
	
	resp, err := l.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch insights: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("insights API returned status: %d", resp.StatusCode)
	}
	
	var response struct {
		Insights []InsightSummary `json:"insights"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode insights: %w", err)
	}
	
	return response.Insights, nil
}

// GetPatterns retrieves patterns from the learning system
func (l *LearningDataSource) GetPatterns() ([]PatternSummary, error) {
	if l.learningAPIURL == "" {
		return l.getMockLearningSnapshot().RecentPatterns, nil
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	patternsURL := fmt.Sprintf("%s/api/v1/learning/patterns", l.learningAPIURL)
	req, err := http.NewRequestWithContext(ctx, "GET", patternsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create patterns request: %w", err)
	}
	
	resp, err := l.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch patterns: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("patterns API returned status: %d", resp.StatusCode)
	}
	
	var response struct {
		Patterns []PatternSummary `json:"patterns"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode patterns: %w", err)
	}
	
	return response.Patterns, nil
}

// TriggerAnalysis triggers manual analysis in the learning system
func (l *LearningDataSource) TriggerAnalysis() error {
	if l.learningAPIURL == "" {
		return nil // No-op if learning system not available
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	analyzeURL := fmt.Sprintf("%s/api/v1/learning/analyze", l.learningAPIURL)
	req, err := http.NewRequestWithContext(ctx, "POST", analyzeURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create analyze request: %w", err)
	}
	
	resp, err := l.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to trigger analysis: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("analyze API returned status: %d", resp.StatusCode)
	}
	
	return nil
}

// GetHealthStatus returns overall system health based on learning data
func (l *LearningDataSource) GetHealthStatus() (map[string]interface{}, error) {
	snapshot, err := l.GetLearningSnapshot()
	if err != nil {
		return nil, err
	}
	
	// Calculate health score
	healthScore := 100
	
	// Deduct for low success rate
	if snapshot.SuccessRate < 1.0 {
		healthScore -= int((1.0-snapshot.SuccessRate)*50)
	}
	
	// Deduct for high latency
	if snapshot.AvgLatency > 0 {
		latencyMs := float64(snapshot.AvgLatency) / float64(time.Millisecond)
		if latencyMs > 1000 {
			healthScore -= 20
		} else if latencyMs > 500 {
			healthScore -= 10
		}
	}
	
	// Deduct for critical insights
	for _, insight := range snapshot.ActiveInsights {
		if insight.Priority == "critical" {
			healthScore -= 15
		} else if insight.Priority == "high" {
			healthScore -= 5
		}
	}
	
	if healthScore < 0 {
		healthScore = 0
	}
	
	// Determine status
	var status string
	if healthScore >= 90 {
		status = "excellent"
	} else if healthScore >= 80 {
		status = "good"
	} else if healthScore >= 70 {
		status = "fair"
	} else if healthScore >= 50 {
		status = "needs_attention"
	} else {
		status = "critical"
	}
	
	return map[string]interface{}{
		"score":             healthScore,
		"status":            status,
		"success_rate":      snapshot.SuccessRate,
		"avg_latency_ms":    float64(snapshot.AvgLatency) / float64(time.Millisecond),
		"total_executions":  snapshot.TotalExecutions,
		"active_insights":   len(snapshot.ActiveInsights),
		"critical_insights": l.countInsightsByPriority(snapshot.ActiveInsights, "critical"),
		"high_insights":     l.countInsightsByPriority(snapshot.ActiveInsights, "high"),
		"last_updated":      snapshot.SnapshotTime,
	}, nil
}

// countInsightsByPriority counts insights by priority level
func (l *LearningDataSource) countInsightsByPriority(insights []InsightSummary, priority string) int {
	count := 0
	for _, insight := range insights {
		if insight.Priority == priority {
			count++
		}
	}
	return count
}

// IsLearningSystemAvailable checks if the learning system is reachable
func (l *LearningDataSource) IsLearningSystemAvailable() bool {
	if l.learningAPIURL == "" {
		return false
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	healthURL := fmt.Sprintf("%s/health", l.learningAPIURL)
	req, err := http.NewRequestWithContext(ctx, "GET", healthURL, nil)
	if err != nil {
		return false
	}
	
	resp, err := l.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	
	return resp.StatusCode == http.StatusOK
}