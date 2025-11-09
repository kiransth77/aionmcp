package autodocs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ReflectionGenerator generates daily reflection documents using learning insights
type ReflectionGenerator struct {
	dataSource DataSource
}

// NewReflectionGenerator creates a new reflection generator
func NewReflectionGenerator(dataSource DataSource) *ReflectionGenerator {
	return &ReflectionGenerator{
		dataSource: dataSource,
	}
}

// Generate creates a reflection document
func (r *ReflectionGenerator) Generate(request GenerationRequest) (*GenerationResult, error) {
	if err := r.Validate(request); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Determine the reflection date (default to today)
	reflectionDate := time.Now()
	if request.DateRange != nil {
		reflectionDate = request.DateRange.StartDate
	}

	// Get learning snapshot and project info
	learningSnapshot, err := r.dataSource.GetLearningSnapshot()
	if err != nil {
		return &GenerationResult{
			Type:    request.Type,
			Success: false,
			Error:   fmt.Sprintf("failed to get learning snapshot: %v", err),
		}, nil
	}

	projectInfo, err := r.dataSource.GetProjectInfo()
	if err != nil {
		return &GenerationResult{
			Type:    request.Type,
			Success: false,
			Error:   fmt.Sprintf("failed to get project info: %v", err),
		}, nil
	}

	// Get commits for the day
	dayStart := time.Date(reflectionDate.Year(), reflectionDate.Month(), reflectionDate.Day(), 0, 0, 0, 0, reflectionDate.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	commits, err := r.dataSource.GetCommits(DateRange{
		StartDate: dayStart,
		EndDate:   dayEnd,
	})
	if err != nil {
		// Don't fail if we can't get commits, just log it
		commits = []GitCommit{}
	}

	// Generate reflection content
	content, metadata, err := r.generateReflection(reflectionDate, learningSnapshot, projectInfo, commits)
	if err != nil {
		return &GenerationResult{
			Type:    request.Type,
			Success: false,
			Error:   fmt.Sprintf("failed to generate reflection: %v", err),
		}, nil
	}

	// Write to file
	if err := WriteToFile(request.OutputPath, content); err != nil {
		return &GenerationResult{
			Type:    request.Type,
			Success: false,
			Error:   fmt.Sprintf("failed to write file: %v", err),
		}, nil
	}

	return &GenerationResult{
		Type:          request.Type,
		OutputPath:    request.OutputPath,
		Success:       true,
		GeneratedAt:   time.Now(),
		ContentLength: len(content),
		Metadata:      metadata,
	}, nil
}

// GetSupportedTypes returns the document types this generator supports
func (r *ReflectionGenerator) GetSupportedTypes() []DocumentType {
	return []DocumentType{DocumentTypeReflection}
}

// Validate checks if the generation request is valid
func (r *ReflectionGenerator) Validate(request GenerationRequest) error {
	if request.Type != DocumentTypeReflection {
		return fmt.Errorf("unsupported document type: %s", request.Type)
	}

	if request.OutputPath == "" {
		return fmt.Errorf("output path is required")
	}

	if request.Format != "" && request.Format != "markdown" {
		return fmt.Errorf("unsupported format: %s (only markdown supported)", request.Format)
	}

	return nil
}

// generateReflection creates the reflection document content
func (r *ReflectionGenerator) generateReflection(date time.Time, learning *LearningSnapshot, projectInfo map[string]interface{}, commits []GitCommit) (string, *DocumentMetadata, error) {
	var content strings.Builder

	// Header
	content.WriteString(fmt.Sprintf("# Daily Reflection - %s\n\n", date.Format("January 2, 2006")))
	content.WriteString(fmt.Sprintf("*Generated automatically at %s*\n\n", time.Now().Format("15:04:05 MST")))

	// Executive Summary
	r.generateExecutiveSummary(&content, learning, commits)

	// Development Activity
	r.generateDevelopmentActivity(&content, commits, projectInfo)

	// Learning Insights
	r.generateLearningInsights(&content, learning)

	// Performance Analysis
	r.generatePerformanceAnalysis(&content, learning)

	// Error Analysis
	r.generateErrorAnalysis(&content, learning)

	// Tool Usage Patterns
	r.generateToolUsagePatterns(&content, learning)

	// Recommendations
	r.generateRecommendations(&content, learning, commits)

	// Goals and Focus Areas
	r.generateGoalsAndFocus(&content, learning)

	// Metadata
	metadata := &DocumentMetadata{
		Version:       "1.0",
		GeneratedAt:   time.Now(),
		DataSources:   []string{"learning_system", "git"},
		LearningStats: learning,
		Tags: map[string]string{
			"reflection_date": date.Format("2006-01-02"),
			"type":            "daily_reflection",
		},
	}

	if len(commits) > 0 {
		metadata.CommitRange = &CommitRange{
			StartDate:   commits[len(commits)-1].Date,
			EndDate:     commits[0].Date,
			CommitCount: len(commits),
		}
	}

	return content.String(), metadata, nil
}

// generateExecutiveSummary creates an executive summary
func (r *ReflectionGenerator) generateExecutiveSummary(content *strings.Builder, learning *LearningSnapshot, commits []GitCommit) {
	content.WriteString("## 📊 Executive Summary\n\n")

	// Key metrics
	content.WriteString("### Key Metrics\n\n")
	content.WriteString(fmt.Sprintf("- **Total Executions**: %d\n", learning.TotalExecutions))
	content.WriteString(fmt.Sprintf("- **Success Rate**: %.1f%%\n", learning.SuccessRate*100))

	if learning.AvgLatency > 0 {
		latencyMs := float64(learning.AvgLatency) / float64(time.Millisecond)
		content.WriteString(fmt.Sprintf("- **Average Latency**: %.1fms\n", latencyMs))
	}

	content.WriteString(fmt.Sprintf("- **Commits Today**: %d\n", len(commits)))
	content.WriteString(fmt.Sprintf("- **Active Insights**: %d\n", len(learning.ActiveInsights)))
	content.WriteString(fmt.Sprintf("- **Patterns Detected**: %d\n\n", len(learning.RecentPatterns)))

	// Overall health assessment
	healthScore := CalculateHealthScore(learning)
	healthStatus := GetHealthStatus(healthScore)
	
	content.WriteString("### System Health\n\n")
	content.WriteString(fmt.Sprintf("**Overall Health Score**: %d/100 (%s)\n\n", healthScore, healthStatus))

	// Quick wins identified
	if len(learning.ActiveInsights) > 0 {
		highPriorityCount := 0
		for _, insight := range learning.ActiveInsights {
			if insight.Priority == "high" || insight.Priority == "critical" {
				highPriorityCount++
			}
		}

		if highPriorityCount > 0 {
			content.WriteString(fmt.Sprintf("⚠️ **%d high-priority issues** require immediate attention.\n\n", highPriorityCount))
		}
	}
}

// generateDevelopmentActivity creates development activity section
func (r *ReflectionGenerator) generateDevelopmentActivity(content *strings.Builder, commits []GitCommit, projectInfo map[string]interface{}) {
	content.WriteString("## 💻 Development Activity\n\n")

	if len(commits) == 0 {
		content.WriteString("No commits were made today.\n\n")
		return
	}

	// Commit summary
	totalInsertions := 0
	totalDeletions := 0
	totalFiles := 0
	authors := make(map[string]int)

	for _, commit := range commits {
		totalInsertions += commit.Insertions
		totalDeletions += commit.Deletions
		totalFiles += commit.ChangedFiles
		authors[commit.Author]++
	}

	content.WriteString(fmt.Sprintf("### Commit Summary\n\n"))
	content.WriteString(fmt.Sprintf("- **Commits**: %d\n", len(commits)))
	content.WriteString(fmt.Sprintf("- **Files Changed**: %d\n", totalFiles))
	content.WriteString(fmt.Sprintf("- **Lines Added**: +%d\n", totalInsertions))
	content.WriteString(fmt.Sprintf("- **Lines Removed**: -%d\n", totalDeletions))
	content.WriteString(fmt.Sprintf("- **Net Change**: %+d lines\n", totalInsertions-totalDeletions))
	content.WriteString(fmt.Sprintf("- **Active Contributors**: %d\n\n", len(authors)))

	// Recent commits
	content.WriteString("### Recent Commits\n\n")
	for i, commit := range commits {
		if i >= 5 { // Show max 5 recent commits
			break
		}

		content.WriteString(fmt.Sprintf("- **%s** ([`%s`](../../commit/%s))\n",
			commit.Subject, commit.ShortHash, commit.Hash))
		content.WriteString(fmt.Sprintf("  *%s at %s*\n",
			commit.Author, commit.Date.Format("15:04")))

		if commit.ChangedFiles > 0 {
			content.WriteString(fmt.Sprintf("  %d files, +%d -%d lines\n",
				commit.ChangedFiles, commit.Insertions, commit.Deletions))
		}
		content.WriteString("\n")
	}
}

// generateLearningInsights creates learning insights section
func (r *ReflectionGenerator) generateLearningInsights(content *strings.Builder, learning *LearningSnapshot) {
	content.WriteString("## 🧠 Learning Insights\n\n")

	if len(learning.ActiveInsights) == 0 {
		content.WriteString("No active insights at this time. The system is learning from ongoing executions.\n\n")
		return
	}

	// Group insights by priority
	criticalInsights := []InsightSummary{}
	highInsights := []InsightSummary{}
	mediumInsights := []InsightSummary{}
	lowInsights := []InsightSummary{}

	for _, insight := range learning.ActiveInsights {
		switch insight.Priority {
		case "critical":
			criticalInsights = append(criticalInsights, insight)
		case "high":
			highInsights = append(highInsights, insight)
		case "medium":
			mediumInsights = append(mediumInsights, insight)
		case "low":
			lowInsights = append(lowInsights, insight)
		}
	}

	// Critical insights
	if len(criticalInsights) > 0 {
		content.WriteString("### 🚨 Critical Issues\n\n")
		for _, insight := range criticalInsights {
			content.WriteString(fmt.Sprintf("**%s**\n", insight.Title))
			content.WriteString(fmt.Sprintf("%s\n\n", insight.Description))
			content.WriteString(fmt.Sprintf("*Recommendation: %s*\n\n", insight.Suggestion))
		}
	}

	// High priority insights
	if len(highInsights) > 0 {
		content.WriteString("### ⚡ High Priority\n\n")
		for _, insight := range highInsights {
			content.WriteString(fmt.Sprintf("- **%s**: %s\n", insight.Title, insight.Description))
			content.WriteString(fmt.Sprintf("  *%s*\n\n", insight.Suggestion))
		}
	}

	// Medium priority insights (show max 3)
	if len(mediumInsights) > 0 {
		content.WriteString("### 📋 Medium Priority\n\n")
		for i, insight := range mediumInsights {
			if i >= 3 {
				content.WriteString(fmt.Sprintf("*...and %d more medium priority insights*\n\n", len(mediumInsights)-i))
				break
			}
			content.WriteString(fmt.Sprintf("- %s: %s\n", insight.Title, insight.Description))
		}
		content.WriteString("\n")
	}
}

// generatePerformanceAnalysis creates performance analysis section
func (r *ReflectionGenerator) generatePerformanceAnalysis(content *strings.Builder, learning *LearningSnapshot) {
	content.WriteString("## ⚡ Performance Analysis\n\n")

	if learning.AvgLatency == 0 {
		content.WriteString("No performance data available.\n\n")
		return
	}

	latencyMs := float64(learning.AvgLatency) / float64(time.Millisecond)

	content.WriteString(fmt.Sprintf("- **Average Response Time**: %.1fms\n", latencyMs))

	// Performance assessment
	var perfAssessment string
	var perfEmoji string
	if latencyMs < 100 {
		perfAssessment = "Excellent"
		perfEmoji = "🟢"
	} else if latencyMs < 500 {
		perfAssessment = "Good"
		perfEmoji = "🟡"
	} else if latencyMs < 1000 {
		perfAssessment = "Fair"
		perfEmoji = "🟠"
	} else {
		perfAssessment = "Needs Improvement"
		perfEmoji = "🔴"
	}

	content.WriteString(fmt.Sprintf("- **Performance Rating**: %s %s\n\n", perfEmoji, perfAssessment))

	// Top performing tools
	if len(learning.TopTools) > 0 {
		content.WriteString("### Fastest Tools\n\n")
		for i, tool := range learning.TopTools {
			if i >= 3 { // Show top 3
				break
			}

			toolLatencyMs := float64(tool.AvgLatency) / float64(time.Millisecond)
			content.WriteString(fmt.Sprintf("- **%s**: %.1fms avg (%.1f%% success)\n",
				tool.Name, toolLatencyMs, tool.SuccessRate*100))
		}
		content.WriteString("\n")
	}
}

// generateErrorAnalysis creates error analysis section
func (r *ReflectionGenerator) generateErrorAnalysis(content *strings.Builder, learning *LearningSnapshot) {
	content.WriteString("## 🐛 Error Analysis\n\n")

	if len(learning.ErrorBreakdown) == 0 {
		content.WriteString("✅ No errors detected in recent executions.\n\n")
		return
	}

	totalErrors := 0
	for _, count := range learning.ErrorBreakdown {
		totalErrors += count
	}

	content.WriteString(fmt.Sprintf("**Total Errors**: %d\n\n", totalErrors))

	content.WriteString("### Error Breakdown\n\n")
	for errorType, count := range learning.ErrorBreakdown {
		percentage := float64(count) / float64(totalErrors) * 100
		content.WriteString(fmt.Sprintf("- **%s**: %d (%.1f%%)\n", errorType, count, percentage))
	}
	content.WriteString("\n")

	// Error patterns
	errorPatterns := []PatternSummary{}
	for _, pattern := range learning.RecentPatterns {
		if pattern.Type == "error" {
			errorPatterns = append(errorPatterns, pattern)
		}
	}

	if len(errorPatterns) > 0 {
		content.WriteString("### Error Patterns\n\n")
		for _, pattern := range errorPatterns {
			content.WriteString(fmt.Sprintf("- **%s** (seen %d times)\n", pattern.Description, pattern.Frequency))
			content.WriteString(fmt.Sprintf("  *First seen: %s, Last seen: %s*\n\n",
				pattern.FirstSeen.Format("Jan 2 15:04"), pattern.LastSeen.Format("Jan 2 15:04")))
		}
	}
}

// generateToolUsagePatterns creates tool usage patterns section
func (r *ReflectionGenerator) generateToolUsagePatterns(content *strings.Builder, learning *LearningSnapshot) {
	content.WriteString("## 🔧 Tool Usage Patterns\n\n")

	if len(learning.TopTools) == 0 {
		content.WriteString("No tool usage data available.\n\n")
		return
	}

	content.WriteString("### Most Used Tools\n\n")

	totalExecutions := 0
	for _, tool := range learning.TopTools {
		totalExecutions += tool.ExecutionCount
	}

	for i, tool := range learning.TopTools {
		if i >= 5 { // Show top 5
			break
		}

		usagePercentage := float64(tool.ExecutionCount) / float64(totalExecutions) * 100
		content.WriteString(fmt.Sprintf("- **%s**: %d executions (%.1f%%)\n",
			tool.Name, tool.ExecutionCount, usagePercentage))
		content.WriteString(fmt.Sprintf("  Success Rate: %.1f%%, Last Used: %s\n\n",
			tool.SuccessRate*100, tool.LastUsed.Format("Jan 2 15:04")))
	}

	// Usage patterns
	usagePatterns := []PatternSummary{}
	for _, pattern := range learning.RecentPatterns {
		if pattern.Type == "usage" {
			usagePatterns = append(usagePatterns, pattern)
		}
	}

	if len(usagePatterns) > 0 {
		content.WriteString("### Usage Patterns\n\n")
		for _, pattern := range usagePatterns {
			content.WriteString(fmt.Sprintf("- %s\n", pattern.Description))
		}
		content.WriteString("\n")
	}
}

// generateRecommendations creates recommendations section
func (r *ReflectionGenerator) generateRecommendations(content *strings.Builder, learning *LearningSnapshot, commits []GitCommit) {
	content.WriteString("## 💡 Recommendations\n\n")

	recommendations := []string{}

	// Based on insights
	criticalCount := 0
	highCount := 0
	for _, insight := range learning.ActiveInsights {
		if insight.Priority == "critical" {
			criticalCount++
		} else if insight.Priority == "high" {
			highCount++
		}
	}

	if criticalCount > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("🚨 **Immediate Action Required**: Address %d critical issues before continuing development", criticalCount))
	}

	if highCount > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("⚡ **High Priority**: Schedule time to resolve %d high-priority issues", highCount))
	}

	// Based on performance
	if learning.AvgLatency > 0 {
		latencyMs := float64(learning.AvgLatency) / float64(time.Millisecond)
		if latencyMs > 1000 {
			recommendations = append(recommendations,
				"🐌 **Performance**: Consider optimizing slow-performing tools (>1s average)")
		}
	}

	// Based on error rate
	if learning.SuccessRate < 0.95 {
		recommendations = append(recommendations,
			"🐛 **Reliability**: Focus on improving error handling and success rates")
	}

	// Based on commit activity
	if len(commits) == 0 {
		recommendations = append(recommendations,
			"📝 **Development**: No commits today - consider making incremental progress")
	} else if len(commits) > 10 {
		recommendations = append(recommendations,
			"🔄 **Process**: High commit frequency - consider batching related changes")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations,
			"✅ **Keep Going**: System is performing well, continue current practices")
	}

	for _, rec := range recommendations {
		content.WriteString(fmt.Sprintf("- %s\n", rec))
	}
	content.WriteString("\n")
}

// generateGoalsAndFocus creates goals and focus areas section
func (r *ReflectionGenerator) generateGoalsAndFocus(content *strings.Builder, learning *LearningSnapshot) {
	content.WriteString("## 🎯 Goals & Focus Areas\n\n")

	content.WriteString("### Tomorrow's Focus\n\n")

	// Dynamic focus areas based on current state
	focusAreas := []string{}

	// Critical issues first
	criticalCount := 0
	for _, insight := range learning.ActiveInsights {
		if insight.Priority == "critical" {
			criticalCount++
		}
	}

	if criticalCount > 0 {
		focusAreas = append(focusAreas,
			fmt.Sprintf("🚨 Resolve %d critical system issues", criticalCount))
	}

	// Performance improvements
	if learning.SuccessRate < 0.95 {
		focusAreas = append(focusAreas,
			"📈 Improve system reliability and error handling")
	}

	// Documentation and learning
	if len(learning.RecentPatterns) > 3 {
		focusAreas = append(focusAreas,
			"📚 Review and act on detected patterns")
	}

	// Default focus areas
	if len(focusAreas) == 0 {
		focusAreas = append(focusAreas,
			"🔧 Continue feature development")
		focusAreas = append(focusAreas,
			"📊 Monitor system performance")
		focusAreas = append(focusAreas,
			"✅ Maintain code quality")
	}

	for _, area := range focusAreas {
		content.WriteString(fmt.Sprintf("- %s\n", area))
	}

	content.WriteString("\n### Success Metrics\n\n")
	content.WriteString("- Maintain >95% success rate\n")
	content.WriteString("- Keep average latency <500ms\n")
	content.WriteString("- Address all critical insights\n")
	content.WriteString("- Make meaningful progress on features\n\n")

	content.WriteString("---\n\n")
	content.WriteString("*This reflection was generated to help improve system performance and development practices. ")
	content.WriteString("Review regularly and adjust focus areas based on emerging patterns and insights.*\n")
}

// calculateHealthScore calculates an overall health score
func (r *ReflectionGenerator) calculateHealthScore(learning *LearningSnapshot) int {
	score := 100

	// Deduct for low success rate
	if learning.SuccessRate < 1.0 {
		score -= int((1.0 - learning.SuccessRate) * 50) // Up to -50 points
	}

	// Deduct for high latency
	if learning.AvgLatency > 0 {
		latencyMs := float64(learning.AvgLatency) / float64(time.Millisecond)
		if latencyMs > 1000 {
			score -= 20 // -20 for >1s latency
		} else if latencyMs > 500 {
			score -= 10 // -10 for >500ms latency
		}
	}

	// Deduct for critical insights
	for _, insight := range learning.ActiveInsights {
		if insight.Priority == "critical" {
			score -= 15 // -15 per critical issue
		} else if insight.Priority == "high" {
			score -= 5 // -5 per high priority issue
		}
	}

	// Ensure minimum score
	if score < 0 {
		score = 0
	}

	return score
}

// getHealthStatus returns a health status string
func (r *ReflectionGenerator) getHealthStatus(score int) string {
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

// writeToFile writes content to the specified file path
func (r *ReflectionGenerator) writeToFile(outputPath, content string) error {
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
