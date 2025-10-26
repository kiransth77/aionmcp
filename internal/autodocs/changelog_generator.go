package autodocs

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ChangelogGenerator generates changelog documents from git history
type ChangelogGenerator struct {
	dataSource DataSource
}

// NewChangelogGenerator creates a new changelog generator
func NewChangelogGenerator(dataSource DataSource) *ChangelogGenerator {
	return &ChangelogGenerator{
		dataSource: dataSource,
	}
}

// Generate creates a changelog document
func (c *ChangelogGenerator) Generate(request GenerationRequest) (*GenerationResult, error) {
	if err := c.Validate(request); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	
	// Determine date range
	dateRange := DateRange{
		StartDate: time.Now().AddDate(0, -1, 0), // Default to last month
		EndDate:   time.Now(),
	}
	
	if request.DateRange != nil {
		dateRange = *request.DateRange
	}
	
	// Get commits and project info
	commits, err := c.dataSource.GetCommits(dateRange)
	if err != nil {
		return &GenerationResult{
			Type:    request.Type,
			Success: false,
			Error:   fmt.Sprintf("failed to get commits: %v", err),
		}, nil
	}
	
	projectInfo, err := c.dataSource.GetProjectInfo()
	if err != nil {
		return &GenerationResult{
			Type:    request.Type,
			Success: false,
			Error:   fmt.Sprintf("failed to get project info: %v", err),
		}, nil
	}
	
	// Generate changelog content
	content, metadata, err := c.generateChangelog(commits, projectInfo, dateRange)
	if err != nil {
		return &GenerationResult{
			Type:    request.Type,
			Success: false,
			Error:   fmt.Sprintf("failed to generate changelog: %v", err),
		}, nil
	}
	
	// Write to file
	if err := c.writeToFile(request.OutputPath, content); err != nil {
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
func (c *ChangelogGenerator) GetSupportedTypes() []DocumentType {
	return []DocumentType{DocumentTypeChangelog}
}

// Validate checks if the generation request is valid
func (c *ChangelogGenerator) Validate(request GenerationRequest) error {
	if request.Type != DocumentTypeChangelog {
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

// generateChangelog creates the changelog content
func (c *ChangelogGenerator) generateChangelog(commits []GitCommit, projectInfo map[string]interface{}, dateRange DateRange) (string, *DocumentMetadata, error) {
	var content strings.Builder
	
	// Header
	content.WriteString("# Changelog\n\n")
	content.WriteString("All notable changes to this project will be documented in this file.\n\n")
	content.WriteString("The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),\n")
	content.WriteString("and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).\n\n")
	
	// Auto-generation notice
	content.WriteString(fmt.Sprintf("*This changelog was automatically generated on %s*\n\n", time.Now().Format("2006-01-02 15:04:05")))
	
	if len(commits) == 0 {
		content.WriteString("## No changes in the specified date range\n\n")
		content.WriteString(fmt.Sprintf("Date range: %s to %s\n\n", 
			dateRange.StartDate.Format("2006-01-02"), 
			dateRange.EndDate.Format("2006-01-02")))
	} else {
		// Group commits by date (daily entries)
		dailyCommits := c.groupCommitsByDate(commits)
		
		// Sort dates in descending order
		var dates []string
		for date := range dailyCommits {
			dates = append(dates, date)
		}
		sort.Slice(dates, func(i, j int) bool {
			return dates[i] > dates[j] // Descending order
		})
		
		// Generate entries for each date
		for _, date := range dates {
			dayCommits := dailyCommits[date]
			c.generateDayEntry(&content, date, dayCommits)
		}
		
		// Summary section
		c.generateSummary(&content, commits, dateRange)
	}
	
	// Metadata
	metadata := &DocumentMetadata{
		Version:     "1.0",
		GeneratedAt: time.Now(),
		DataSources: []string{"git"},
		CommitRange: &CommitRange{
			StartDate:   dateRange.StartDate,
			EndDate:     dateRange.EndDate,
			CommitCount: len(commits),
		},
		Tags: map[string]string{
			"format": "keepachangelog",
		},
	}
	
	if len(commits) > 0 {
		metadata.CommitRange.StartCommit = commits[len(commits)-1].ShortHash
		metadata.CommitRange.EndCommit = commits[0].ShortHash
	}
	
	return content.String(), metadata, nil
}

// groupCommitsByDate groups commits by their date
func (c *ChangelogGenerator) groupCommitsByDate(commits []GitCommit) map[string][]GitCommit {
	dailyCommits := make(map[string][]GitCommit)
	
	for _, commit := range commits {
		date := commit.Date.Format("2006-01-02")
		dailyCommits[date] = append(dailyCommits[date], commit)
	}
	
	return dailyCommits
}

// generateDayEntry generates a changelog entry for a specific day
func (c *ChangelogGenerator) generateDayEntry(content *strings.Builder, date string, commits []GitCommit) {
	// Parse date for better formatting
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		parsedDate = time.Now()
	}
	
	content.WriteString(fmt.Sprintf("## %s\n\n", parsedDate.Format("2006-01-02 (Monday)")))
	
	// Categorize commits
	categories := c.categorizeCommits(commits)
	
	// Define category order and display names
	categoryOrder := []string{"breaking", "feature", "fix", "perf", "docs", "refactor", "test", "chore", "style", "ci", "other"}
	categoryNames := map[string]string{
		"breaking": "💥 Breaking Changes",
		"feature":  "✨ Features",
		"fix":      "🐛 Bug Fixes",
		"perf":     "⚡ Performance",
		"docs":     "📚 Documentation",
		"refactor": "♻️ Code Refactoring",
		"test":     "✅ Tests",
		"chore":    "🔧 Chores",
		"style":    "🎨 Styles",
		"ci":       "👷 CI/CD",
		"other":    "📦 Other",
	}
	
	// Write each category
	for _, category := range categoryOrder {
		categoryCommits := categories[category]
		if len(categoryCommits) == 0 {
			continue
		}
		
		content.WriteString(fmt.Sprintf("### %s\n\n", categoryNames[category]))
		
		for _, commit := range categoryCommits {
			c.writeCommitEntry(content, commit)
		}
		content.WriteString("\n")
	}
}

// categorizeCommits categorizes commits by type
func (c *ChangelogGenerator) categorizeCommits(commits []GitCommit) map[string][]GitCommit {
	categories := make(map[string][]GitCommit)
	
	// Initialize categories
	categoryNames := []string{"breaking", "feature", "fix", "perf", "docs", "refactor", "test", "chore", "style", "ci", "other"}
	for _, name := range categoryNames {
		categories[name] = []GitCommit{}
	}
	
	// Categorize each commit
	for _, commit := range commits {
		category := c.categorizeCommit(commit)
		categories[category] = append(categories[category], commit)
	}
	
	return categories
}

// categorizeCommit determines the category of a commit
func (c *ChangelogGenerator) categorizeCommit(commit GitCommit) string {
	subject := strings.ToLower(commit.Subject)
	
	// Check for breaking changes first
	if strings.Contains(subject, "breaking") || strings.Contains(subject, "!:") || strings.Contains(commit.Body, "BREAKING CHANGE") {
		return "breaking"
	}
	
	// Define patterns for different categories
	patterns := map[string][]string{
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
	
	for category, keywords := range patterns {
		for _, keyword := range keywords {
			if strings.Contains(subject, keyword) {
				return category
			}
		}
	}
	
	return "other"
}

// writeCommitEntry writes a single commit entry
func (c *ChangelogGenerator) writeCommitEntry(content *strings.Builder, commit GitCommit) {
	// Format: - subject (shortHash) by author
	content.WriteString(fmt.Sprintf("- %s ([`%s`](../../commit/%s))", 
		commit.Subject, commit.ShortHash, commit.Hash))
	
	// Add author if different from previous commit
	content.WriteString(fmt.Sprintf(" by %s", commit.Author))
	
	// Add file change stats if significant
	if commit.ChangedFiles > 0 {
		content.WriteString(fmt.Sprintf(" (%d files", commit.ChangedFiles))
		if commit.Insertions > 0 || commit.Deletions > 0 {
			content.WriteString(fmt.Sprintf(", +%d/-%d lines", commit.Insertions, commit.Deletions))
		}
		content.WriteString(")")
	}
	
	content.WriteString("\n")
	
	// Add body if it contains important information and is not too long
	if len(commit.Body) > 0 && len(commit.Body) < 200 && !strings.Contains(strings.ToLower(commit.Body), "signed-off-by") {
		// Format body as indented text
		bodyLines := strings.Split(strings.TrimSpace(commit.Body), "\n")
		for _, line := range bodyLines {
			if strings.TrimSpace(line) != "" {
				content.WriteString(fmt.Sprintf("  %s\n", strings.TrimSpace(line)))
			}
		}
	}
}

// generateSummary generates a summary section
func (c *ChangelogGenerator) generateSummary(content *strings.Builder, commits []GitCommit, dateRange DateRange) {
	content.WriteString("## Summary\n\n")
	
	// Basic statistics
	content.WriteString(fmt.Sprintf("**Period:** %s to %s\n\n", 
		dateRange.StartDate.Format("2006-01-02"), 
		dateRange.EndDate.Format("2006-01-02")))
	
	content.WriteString(fmt.Sprintf("**Total commits:** %d\n\n", len(commits)))
	
	// Category breakdown
	categories := c.categorizeCommits(commits)
	content.WriteString("**Changes by type:**\n\n")
	
	categoryNames := map[string]string{
		"feature":  "Features",
		"fix":      "Bug Fixes", 
		"perf":     "Performance",
		"docs":     "Documentation",
		"refactor": "Refactoring",
		"test":     "Tests",
		"chore":    "Chores",
		"style":    "Styles",
		"ci":       "CI/CD",
		"other":    "Other",
		"breaking": "Breaking Changes",
	}
	
	for category, commits := range categories {
		if len(commits) > 0 {
			content.WriteString(fmt.Sprintf("- %s: %d\n", categoryNames[category], len(commits)))
		}
	}
	
	// Author statistics
	authors := make(map[string]int)
	totalInsertions := 0
	totalDeletions := 0
	totalFiles := 0
	
	for _, commit := range commits {
		authors[commit.Author]++
		totalInsertions += commit.Insertions
		totalDeletions += commit.Deletions
		totalFiles += commit.ChangedFiles
	}
	
	content.WriteString(fmt.Sprintf("\n**Contributors:** %d\n\n", len(authors)))
	
	// Top contributors
	type authorStat struct {
		name   string
		commits int
	}
	
	var authorStats []authorStat
	for name, count := range authors {
		authorStats = append(authorStats, authorStat{name, count})
	}
	
	sort.Slice(authorStats, func(i, j int) bool {
		return authorStats[i].commits > authorStats[j].commits
	})
	
	for i, stat := range authorStats {
		if i >= 5 { // Top 5 contributors
			break
		}
		content.WriteString(fmt.Sprintf("- %s: %d commits\n", stat.name, stat.commits))
	}
	
	// Code statistics
	content.WriteString(fmt.Sprintf("\n**Code changes:**\n"))
	content.WriteString(fmt.Sprintf("- Files changed: %d\n", totalFiles))
	content.WriteString(fmt.Sprintf("- Lines added: +%d\n", totalInsertions))
	content.WriteString(fmt.Sprintf("- Lines removed: -%d\n", totalDeletions))
	content.WriteString(fmt.Sprintf("- Net change: %+d lines\n\n", totalInsertions-totalDeletions))
}

// writeToFile writes content to the specified file path
func (c *ChangelogGenerator) writeToFile(outputPath, content string) error {
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