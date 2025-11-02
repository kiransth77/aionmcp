package autodocs

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// GetHealthStatus returns a health status string based on the score
func GetHealthStatus(score int) string {
	if score >= 90 {
		return "Excellent"
	} else if score >= 80 {
		return "Good"
	} else if score >= 70 {
		return "Fair"
	} else if score >= 50 {
		return "Needs Attention"
	} else {
		return "Critical"
	}
}

// WriteToFile writes content to the specified file path
func WriteToFile(outputPath, content string) error {
	// Ensure directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Write file
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	return nil
}

// CalculateHealthScore calculates system health score based on learning snapshot
// This is a shared utility used across multiple generators to ensure consistent scoring
func CalculateHealthScore(learning *LearningSnapshot) int {
	score := 100
	
	// Deduct for low success rate
	if learning.SuccessRate < 1.0 {
		score -= int((1.0 - learning.SuccessRate) * 50) // maxSuccessRateDeduction = 50
	}
	
	// Deduct for high latency
	if learning.AvgLatency > 0 {
		latencyMs := float64(learning.AvgLatency) / float64(time.Millisecond)
		if latencyMs > 1000 {
			score -= 20 // highLatencyDeduction = 20
		} else if latencyMs > 500 {
			score -= 10 // mediumLatencyDeduction = 10
		}
	}
	
	// Deduct for critical insights
	for _, insight := range learning.ActiveInsights {
		if insight.Priority == "critical" {
			score -= 15 // criticalIssueDeduction = 15
		} else if insight.Priority == "high" {
			score -= 5 // highPriorityDeduction = 5
		}
	}
	
	// Ensure minimum score
	if score < 0 {
		score = 0
	}
	
	return score
}

// CommitCategorizationPatterns defines the keywords used to categorize commits
// This is a shared constant to ensure consistency across git analysis and changelog generation
var CommitCategorizationPatterns = map[string][]string{
	"feature":  {"feat:", "feature:", "add:", "implement", "new"},
	"fix":      {"fix:", "bug:", "bugfix:", "hotfix:", "patch:"},
	"perf":     {"perf:", "performance:", "optimize", "speed", "improve performance"},
	"docs":     {"docs:", "doc:", "documentation", "readme", "changelog"},
	"refactor": {"refactor:", "cleanup:", "clean:", "reorganize"},
	"test":     {"test:", "tests:", "testing:", "spec:"},
	"chore":    {"chore:", "bump:", "update:", "upgrade:", "version:", "deps:"},
	"style":    {"style:", "format:", "lint:", "prettier:"},
	"ci":       {"ci:", "build:", "deploy:", "pipeline:", "github:", "actions:"},
}
