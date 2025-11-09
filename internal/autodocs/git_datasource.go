package autodocs

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// GitDataSource implements DataSource for git repository information
type GitDataSource struct {
	repoPath string
}

// NewGitDataSource creates a new git data source
func NewGitDataSource(repoPath string) *GitDataSource {
	return &GitDataSource{
		repoPath: repoPath,
	}
}

// GetCommits retrieves git commits within a date range
func (g *GitDataSource) GetCommits(dateRange DateRange) ([]GitCommit, error) {
	// Format git log command with date range
	sinceDate := dateRange.StartDate.Format("2006-01-02")
	untilDate := dateRange.EndDate.Format("2006-01-02")
	
	// Use pretty format to get structured commit information
	cmd := exec.Command("git", "log",
		"--since="+sinceDate,
		"--until="+untilDate,
		"--pretty=format:%H|%h|%an|%ae|%ai|%s|%b",
		"--numstat",
		"--no-merges",
	)
	cmd.Dir = g.repoPath
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git log: %w", err)
	}
	
	return g.parseGitLog(string(output))
}

// parseGitLog parses the output of git log command
func (g *GitDataSource) parseGitLog(logOutput string) ([]GitCommit, error) {
	lines := strings.Split(logOutput, "\n")
	var commits []GitCommit
	var currentCommit *GitCommit
	
	commitPattern := regexp.MustCompile(`^([a-f0-9]{40})\|([a-f0-9]+)\|([^|]+)\|([^|]+)\|([^|]+)\|([^|]*)\|(.*)`)
	statPattern := regexp.MustCompile(`^(\d+)\s+(\d+)\s+(.+)`)
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Check if this is a commit header line
		if matches := commitPattern.FindStringSubmatch(line); matches != nil {
			// Save previous commit if exists
			if currentCommit != nil {
				commits = append(commits, *currentCommit)
			}
			
			// Parse commit date
			commitDate, err := time.Parse("2006-01-02 15:04:05 -0700", matches[5])
			if err != nil {
				commitDate = time.Now() // Fallback
			}
			
			// Create new commit
			currentCommit = &GitCommit{
				Hash:      matches[1],
				ShortHash: matches[2],
				Author:    matches[3],
				Email:     matches[4],
				Date:      commitDate,
				Subject:   matches[6],
				Body:      matches[7],
			}
		} else if statMatches := statPattern.FindStringSubmatch(line); statMatches != nil && currentCommit != nil {
			// Parse file statistics
			insertions, _ := strconv.Atoi(statMatches[1])
			deletions, _ := strconv.Atoi(statMatches[2])
			
			currentCommit.Insertions += insertions
			currentCommit.Deletions += deletions
			currentCommit.ChangedFiles++
		}
	}
	
	// Add the last commit
	if currentCommit != nil {
		commits = append(commits, *currentCommit)
	}
	
	return commits, nil
}

// GetLearningSnapshot retrieves current learning system data
func (g *GitDataSource) GetLearningSnapshot() (*LearningSnapshot, error) {
	// This would typically call the learning system API
	// For now, return a placeholder implementation
	return &LearningSnapshot{
		TotalExecutions: 0,
		SuccessRate:     1.0,
		AvgLatency:      0,
		TopTools:        []ToolUsageInfo{},
		ErrorBreakdown:  map[string]int{},
		RecentPatterns:  []PatternSummary{},
		ActiveInsights:  []InsightSummary{},
		SnapshotTime:    time.Now(),
	}, nil
}

// GetProjectInfo retrieves general project information
func (g *GitDataSource) GetProjectInfo() (map[string]interface{}, error) {
	info := make(map[string]interface{})
	
	// Get current branch
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = g.repoPath
	if output, err := cmd.Output(); err == nil {
		info["current_branch"] = strings.TrimSpace(string(output))
	}
	
	// Get latest commit hash
	cmd = exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = g.repoPath
	if output, err := cmd.Output(); err == nil {
		info["latest_commit"] = strings.TrimSpace(string(output))
	}
	
	// Get total commit count
	cmd = exec.Command("git", "rev-list", "--count", "HEAD")
	cmd.Dir = g.repoPath
	if output, err := cmd.Output(); err == nil {
		if count, err := strconv.Atoi(strings.TrimSpace(string(output))); err == nil {
			info["total_commits"] = count
		}
	}
	
	// Get repository status
	cmd = exec.Command("git", "status", "--porcelain")
	cmd.Dir = g.repoPath
	if output, err := cmd.Output(); err == nil {
		dirty := len(strings.TrimSpace(string(output))) > 0
		info["dirty"] = dirty
	}
	
	// Get remote URL
	cmd = exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir = g.repoPath
	if output, err := cmd.Output(); err == nil {
		info["remote_url"] = strings.TrimSpace(string(output))
	}
	
	// Get creation date (first commit)
	cmd = exec.Command("git", "log", "--reverse", "--format=%ai", "--max-count=1")
	cmd.Dir = g.repoPath
	if output, err := cmd.Output(); err == nil {
		if firstCommitDate, err := time.Parse("2006-01-02 15:04:05 -0700", strings.TrimSpace(string(output))); err == nil {
			info["created_at"] = firstCommitDate
		}
	}
	
	info["snapshot_time"] = time.Now()
	
	return info, nil
}

// GetCommitsSince retrieves commits since a specific commit hash
func (g *GitDataSource) GetCommitsSince(sinceCommit string) ([]GitCommit, error) {
	cmd := exec.Command("git", "log",
		sinceCommit+"..HEAD",
		"--pretty=format:%H|%h|%an|%ae|%ai|%s|%b",
		"--numstat",
		"--no-merges",
	)
	cmd.Dir = g.repoPath
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git log since %s: %w", sinceCommit, err)
	}
	
	return g.parseGitLog(string(output))
}

// GetTags retrieves git tags with their information
func (g *GitDataSource) GetTags() ([]map[string]interface{}, error) {
	cmd := exec.Command("git", "tag", "-l", "--sort=-version:refname", "--format=%(refname:short)|%(objectname)|%(creatordate:iso8601)")
	cmd.Dir = g.repoPath
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git tags: %w", err)
	}
	
	var tags []map[string]interface{}
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		
		parts := strings.Split(line, "|")
		if len(parts) >= 3 {
			tag := map[string]interface{}{
				"name": parts[0],
				"hash": parts[1],
			}
			
			if tagDate, err := time.Parse("2006-01-02 15:04:05 -0700", parts[2]); err == nil {
				tag["date"] = tagDate
			}
			
			tags = append(tags, tag)
		}
	}
	
	return tags, nil
}

// GetCurrentVersion attempts to determine the current version from git tags
func (g *GitDataSource) GetCurrentVersion() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	cmd.Dir = g.repoPath
	
	output, err := cmd.Output()
	if err != nil {
		// If no tags exist, return a default version
		return "v0.1.0", nil
	}
	
	return strings.TrimSpace(string(output)), nil
}

// CategorizeCommit categorizes a commit based on its message
func (g *GitDataSource) CategorizeCommit(commit GitCommit) string {
	subject := strings.ToLower(commit.Subject)
	
	// Use shared categorization patterns from utils
	for category, keywords := range CommitCategorizationPatterns {
		for _, keyword := range keywords {
			if strings.Contains(subject, keyword) {
				return category
			}
		}
	}
	
	// Check for breaking changes
	if strings.Contains(subject, "breaking") || strings.Contains(subject, "!:") {
		return "breaking"
	}
	
	return "other"
}

// GetCommitStats returns statistics about commits in a date range
func (g *GitDataSource) GetCommitStats(dateRange DateRange) (map[string]interface{}, error) {
	commits, err := g.GetCommits(dateRange)
	if err != nil {
		return nil, err
	}
	
	stats := map[string]interface{}{
		"total_commits":    len(commits),
		"total_insertions": 0,
		"total_deletions":  0,
		"total_files":      0,
		"categories":       make(map[string]int),
		"authors":          make(map[string]int),
		"daily_activity":   make(map[string]int),
	}
	
	categories := stats["categories"].(map[string]int)
	authors := stats["authors"].(map[string]int)
	daily := stats["daily_activity"].(map[string]int)
	
	totalInsertions := 0
	totalDeletions := 0
	totalFiles := 0
	
	for _, commit := range commits {
		// Update totals
		totalInsertions += commit.Insertions
		totalDeletions += commit.Deletions
		totalFiles += commit.ChangedFiles
		
		// Categorize commit
		category := g.CategorizeCommit(commit)
		categories[category]++
		
		// Track authors
		authors[commit.Author]++
		
		// Track daily activity
		day := commit.Date.Format("2006-01-02")
		daily[day]++
	}
	
	stats["total_insertions"] = totalInsertions
	stats["total_deletions"] = totalDeletions
	stats["total_files"] = totalFiles
	
	// Calculate averages
	if len(commits) > 0 {
		stats["avg_insertions_per_commit"] = totalInsertions / len(commits)
		stats["avg_deletions_per_commit"] = totalDeletions / len(commits)
		stats["avg_files_per_commit"] = totalFiles / len(commits)
	}
	
	return stats, nil
}