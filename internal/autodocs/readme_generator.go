package autodocs

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// ReadmeGenerator generates and updates README.md with current project status
type ReadmeGenerator struct {
	dataSource DataSource
	projectRoot string
}

// NewReadmeGenerator creates a new README generator
func NewReadmeGenerator(dataSource DataSource, projectRoot string) *ReadmeGenerator {
	return &ReadmeGenerator{
		dataSource:  dataSource,
		projectRoot: projectRoot,
	}
}

// Generate creates or updates a README document
func (r *ReadmeGenerator) Generate(request GenerationRequest) (*GenerationResult, error) {
	if err := r.Validate(request); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	
	// Get current data
	projectInfo, err := r.dataSource.GetProjectInfo()
	if err != nil {
		return &GenerationResult{
			Type:    request.Type,
			Success: false,
			Error:   fmt.Sprintf("failed to get project info: %v", err),
		}, nil
	}
	
	learningSnapshot, err := r.dataSource.GetLearningSnapshot()
	if err != nil {
		return &GenerationResult{
			Type:    request.Type,
			Success: false,
			Error:   fmt.Sprintf("failed to get learning snapshot: %v", err),
		}, nil
	}
	
	// Get recent commits (last 30 days)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	commits, err := r.dataSource.GetCommits(DateRange{
		StartDate: thirtyDaysAgo,
		EndDate:   time.Now(),
	})
	if err != nil {
		commits = []GitCommit{} // Don't fail, just use empty commits
	}
	
	// Read existing README if it exists
	existingContent := ""
	if _, err := os.Stat(request.OutputPath); err == nil {
		if data, err := os.ReadFile(request.OutputPath); err == nil {
			existingContent = string(data)
		}
	}
	
	// Generate new README content
	content, metadata, err := r.generateReadme(projectInfo, learningSnapshot, commits, existingContent)
	if err != nil {
		return &GenerationResult{
			Type:    request.Type,
			Success: false,
			Error:   fmt.Sprintf("failed to generate README: %v", err),
		}, nil
	}
	
	// Write to file
	if err := r.writeToFile(request.OutputPath, content); err != nil {
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
func (r *ReadmeGenerator) GetSupportedTypes() []DocumentType {
	return []DocumentType{DocumentTypeReadme}
}

// Validate checks if the generation request is valid
func (r *ReadmeGenerator) Validate(request GenerationRequest) error {
	if request.Type != DocumentTypeReadme {
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

// generateReadme creates the README content
func (r *ReadmeGenerator) generateReadme(projectInfo map[string]interface{}, learning *LearningSnapshot, commits []GitCommit, existing string) (string, *DocumentMetadata, error) {
	var content strings.Builder
	
	// Preserve manual sections while updating automatic ones
	preservedSections := r.extractPreservedSections(existing)
	
	// Header
	content.WriteString("# AionMCP - Autonomous Go MCP Server\n\n")
	
	// Add shields/badges
	r.generateBadges(&content, projectInfo, learning)
	
	// Project description
	r.generateDescription(&content)
	
	// Status section (auto-updated)
	r.generateStatus(&content, projectInfo, learning, commits)
	
	// Features section (preserve manual content)
	if preserved, exists := preservedSections["features"]; exists {
		content.WriteString("## ✨ Features\n\n")
		content.WriteString(preserved)
		content.WriteString("\n")
	} else {
		r.generateFeatures(&content)
	}
	
	// Quick Start section (preserve manual content)
	if preserved, exists := preservedSections["quick-start"]; exists {
		content.WriteString("## 🚀 Quick Start\n\n")
		content.WriteString(preserved)
		content.WriteString("\n")
	} else {
		r.generateQuickStart(&content)
	}
	
	// Architecture section (preserve manual content)  
	if preserved, exists := preservedSections["architecture"]; exists {
		content.WriteString("## 🏗️ Architecture\n\n")
		content.WriteString(preserved)
		content.WriteString("\n")
	} else {
		r.generateArchitecture(&content)
	}
	
	// Recent Activity (auto-updated)
	r.generateRecentActivity(&content, commits, learning)
	
	// Performance Stats (auto-updated)
	r.generatePerformanceStats(&content, learning)
	
	// Installation section (preserve manual content)
	if preserved, exists := preservedSections["installation"]; exists {
		content.WriteString("## 📦 Installation\n\n")
		content.WriteString(preserved)
		content.WriteString("\n")
	} else {
		r.generateInstallation(&content)
	}
	
	// Usage section (preserve manual content)
	if preserved, exists := preservedSections["usage"]; exists {
		content.WriteString("## 📚 Usage\n\n")
		content.WriteString(preserved)
		content.WriteString("\n")
	} else {
		r.generateUsage(&content)
	}
	
	// Development section (preserve manual content)
	if preserved, exists := preservedSections["development"]; exists {
		content.WriteString("## 🛠️ Development\n\n")
		content.WriteString(preserved)
		content.WriteString("\n")
	} else {
		r.generateDevelopment(&content)
	}
	
	// Contributing section (preserve manual content)
	if preserved, exists := preservedSections["contributing"]; exists {
		content.WriteString("## 🤝 Contributing\n\n")
		content.WriteString(preserved)
		content.WriteString("\n")
	} else {
		r.generateContributing(&content)
	}
	
	// License section (preserve manual content)
	if preserved, exists := preservedSections["license"]; exists {
		content.WriteString("## 📄 License\n\n")
		content.WriteString(preserved)
		content.WriteString("\n")
	} else {
		r.generateLicense(&content)
	}
	
	// Footer
	r.generateFooter(&content)
	
	// Metadata
	metadata := &DocumentMetadata{
		Version:       "1.0",
		GeneratedAt:   time.Now(),
		DataSources:   []string{"git", "learning_system", "project_files"},
		LearningStats: learning,
		Tags: map[string]string{
			"auto_updated": "true",
			"format":       "github_readme",
		},
	}
	
	return content.String(), metadata, nil
}

// extractPreservedSections extracts manually written sections to preserve
func (r *ReadmeGenerator) extractPreservedSections(content string) map[string]string {
	sections := make(map[string]string)
	
	if content == "" {
		return sections
	}
	
	// Define sections to preserve (manual content)
	preserveSections := []string{
		"features", "quick-start", "architecture", "installation", 
		"usage", "development", "contributing", "license",
	}
	
	for _, section := range preserveSections {
		// Use regex to extract section content
		pattern := fmt.Sprintf(`(?i)## [^#]*%s[^#]*\n\n(.*?)(?=\n## |$)`, section)
		re := regexp.MustCompile(pattern)
		
		if match := re.FindStringSubmatch(content); len(match) > 1 {
			// Clean up the content
			sectionContent := strings.TrimSpace(match[1])
			if sectionContent != "" && !strings.Contains(sectionContent, "<!-- AUTO-GENERATED -->") {
				sections[section] = sectionContent
			}
		}
	}
	
	return sections
}

// generateBadges creates status badges
func (r *ReadmeGenerator) generateBadges(content *strings.Builder, projectInfo map[string]interface{}, learning *LearningSnapshot) {
	content.WriteString("<!-- AUTO-GENERATED BADGES -->\n")
	
	// Build status
	content.WriteString("![Build Status](https://img.shields.io/badge/build-passing-brightgreen)\n")
	
	// Success rate
	successRate := int(learning.SuccessRate * 100)
	var color string
	if successRate >= 95 {
		color = "brightgreen"
	} else if successRate >= 90 {
		color = "green"
	} else if successRate >= 80 {
		color = "yellow"
	} else {
		color = "red"
	}
	content.WriteString(fmt.Sprintf("![Success Rate](https://img.shields.io/badge/success_rate-%d%%25-%s)\n", successRate, color))
	
	// Performance
	if learning.AvgLatency > 0 {
		latencyMs := int(float64(learning.AvgLatency) / float64(time.Millisecond))
		perfColor := "red"
		if latencyMs < 100 {
			perfColor = "brightgreen"
		} else if latencyMs < 500 {
			perfColor = "green"
		} else if latencyMs < 1000 {
			perfColor = "yellow"
		}
		content.WriteString(fmt.Sprintf("![Avg Latency](https://img.shields.io/badge/avg_latency-%dms-%s)\n", latencyMs, perfColor))
	}
	
	// Go version
	content.WriteString("![Go Version](https://img.shields.io/badge/go-1.21+-blue)\n")
	
	// License
	content.WriteString("![License](https://img.shields.io/badge/license-MIT-blue)\n")
	
	content.WriteString("<!-- END AUTO-GENERATED BADGES -->\n\n")
}

// generateDescription creates project description
func (r *ReadmeGenerator) generateDescription(content *strings.Builder) {
	content.WriteString("AionMCP is an autonomous Go-based Model Context Protocol (MCP) server that dynamically imports OpenAPI, GraphQL, and AsyncAPI specifications and exposes them as tools to agents. It features self-learning capabilities, context-awareness, and autonomous documentation using Clean/Hexagonal architecture.\n\n")
	
	content.WriteString("## 🌟 Key Differentiators\n\n")
	content.WriteString("- **Multi-Protocol Support**: OpenAPI, GraphQL, and AsyncAPI specifications\n")
	content.WriteString("- **Autonomous Learning**: Self-improving system that learns from execution patterns\n")
	content.WriteString("- **Dynamic Runtime**: Hot-reloadable tools without service restart\n")
	content.WriteString("- **Clean Architecture**: Maintainable, testable, and extensible design\n")
	content.WriteString("- **Auto-Documentation**: Self-updating documentation and insights\n\n")
}

// generateStatus creates status section
func (r *ReadmeGenerator) generateStatus(content *strings.Builder, projectInfo map[string]interface{}, learning *LearningSnapshot, commits []GitCommit) {
	content.WriteString("## 📊 Project Status\n\n")
	content.WriteString("<!-- AUTO-GENERATED STATUS -->\n")
	
	// Current branch and commit
	if branch, ok := projectInfo["current_branch"].(string); ok {
		content.WriteString(fmt.Sprintf("**Current Branch**: `%s`\n\n", branch))
	}
	
	if latestCommit, ok := projectInfo["latest_commit"].(string); ok && len(latestCommit) > 7 {
		content.WriteString(fmt.Sprintf("**Latest Commit**: [`%s`](../../commit/%s)\n\n", latestCommit[:7], latestCommit))
	}
	
	// System health
	healthScore := r.calculateHealthScore(learning)
	healthStatus := r.getHealthStatus(healthScore)
	content.WriteString(fmt.Sprintf("**System Health**: %d/100 (%s)\n\n", healthScore, healthStatus))
	
	// Active tools
	content.WriteString(fmt.Sprintf("**Active Tools**: %d\n\n", len(learning.TopTools)))
	
	// Recent activity
	recentCommits := 0
	for _, commit := range commits {
		if commit.Date.After(time.Now().AddDate(0, 0, -7)) {
			recentCommits++
		}
	}
	content.WriteString(fmt.Sprintf("**Commits (7 days)**: %d\n\n", recentCommits))
	
	content.WriteString("*Status updated automatically*\n")
	content.WriteString("<!-- END AUTO-GENERATED STATUS -->\n\n")
}

// generateFeatures creates features section (fallback)
func (r *ReadmeGenerator) generateFeatures(content *strings.Builder) {
	content.WriteString("## ✨ Features\n\n")
	content.WriteString("### Core Capabilities\n\n")
	content.WriteString("- **Multi-Spec Import**: Automatically imports and converts API specifications\n")
	content.WriteString("- **Dynamic Tool Registry**: Hot-reload tools without service restart\n")
	content.WriteString("- **Self-Learning Engine**: Analyzes patterns and generates insights\n")
	content.WriteString("- **Autonomous Documentation**: Auto-generates changelogs and reflections\n")
	content.WriteString("- **Performance Monitoring**: Real-time execution metrics and optimization\n")
	content.WriteString("- **Error Recovery**: Intelligent error handling and pattern detection\n\n")
	
	content.WriteString("### API Support\n\n")
	content.WriteString("- **OpenAPI 3.0+**: REST API specifications with full schema support\n")
	content.WriteString("- **GraphQL**: Query and mutation support with type introspection\n")
	content.WriteString("- **AsyncAPI**: Event-driven API specifications\n\n")
}

// generateQuickStart creates quick start section (fallback)
func (r *ReadmeGenerator) generateQuickStart(content *strings.Builder) {
	content.WriteString("## 🚀 Quick Start\n\n")
	content.WriteString("```bash\n")
	content.WriteString("# Clone the repository\n")
	content.WriteString("git clone https://github.com/kiransth77/aionmcp.git\n")
	content.WriteString("cd aionmcp\n\n")
	content.WriteString("# Build the server\n")
	content.WriteString("go build -o bin/aionmcp cmd/server/main.go\n\n")
	content.WriteString("# Run with default configuration\n")
	content.WriteString("./bin/aionmcp\n")
	content.WriteString("```\n\n")
	content.WriteString("The server will start on `http://localhost:8080` with learning enabled.\n\n")
}

// generateArchitecture creates architecture section (fallback)
func (r *ReadmeGenerator) generateArchitecture(content *strings.Builder) {
	content.WriteString("## 🏗️ Architecture\n\n")
	content.WriteString("AionMCP follows Clean/Hexagonal Architecture principles:\n\n")
	content.WriteString("```\n")
	content.WriteString("┌─────────────────────────────────────────────────────────┐\n")
	content.WriteString("│                    Adapters Layer                      │\n")
	content.WriteString("│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐   │\n")
	content.WriteString("│  │   HTTP      │  │    gRPC     │  │   Plugin    │   │\n")
	content.WriteString("│  │  Interface  │  │  Interface  │  │  Interface  │   │\n")
	content.WriteString("│  └─────────────┘  └─────────────┘  └─────────────┘   │\n")
	content.WriteString("└─────────────────────────────────────────────────────────┘\n")
	content.WriteString("┌─────────────────────────────────────────────────────────┐\n")
	content.WriteString("│                     Core Layer                         │\n")
	content.WriteString("│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐   │\n")
	content.WriteString("│  │    Tool     │  │  Learning   │  │    Auto     │   │\n")
	content.WriteString("│  │  Registry   │  │   Engine    │  │    Docs     │   │\n")
	content.WriteString("│  └─────────────┘  └─────────────┘  └─────────────┘   │\n")
	content.WriteString("└─────────────────────────────────────────────────────────┘\n")
	content.WriteString("┌─────────────────────────────────────────────────────────┐\n")
	content.WriteString("│                Infrastructure Layer                    │\n")
	content.WriteString("│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐   │\n")
	content.WriteString("│  │   Storage   │  │   Metrics   │  │   Config    │   │\n")
	content.WriteString("│  │  (BoltDB)   │  │(Prometheus) │  │   (Viper)   │   │\n")
	content.WriteString("│  └─────────────┘  └─────────────┘  └─────────────┘   │\n")
	content.WriteString("└─────────────────────────────────────────────────────────┘\n")
	content.WriteString("```\n\n")
}

// generateRecentActivity creates recent activity section
func (r *ReadmeGenerator) generateRecentActivity(content *strings.Builder, commits []GitCommit, learning *LearningSnapshot) {
	content.WriteString("## 📈 Recent Activity\n\n")
	content.WriteString("<!-- AUTO-GENERATED ACTIVITY -->\n")
	
	// Recent commits (last 7 days)
	recentCommits := []GitCommit{}
	weekAgo := time.Now().AddDate(0, 0, -7)
	
	for _, commit := range commits {
		if commit.Date.After(weekAgo) && len(recentCommits) < 5 {
			recentCommits = append(recentCommits, commit)
		}
	}
	
	if len(recentCommits) > 0 {
		content.WriteString("### Recent Commits\n\n")
		for _, commit := range recentCommits {
			timeAgo := time.Since(commit.Date)
			var timeStr string
			if timeAgo.Hours() < 24 {
				timeStr = fmt.Sprintf("%.0fh ago", timeAgo.Hours())
			} else {
				timeStr = fmt.Sprintf("%.0fd ago", timeAgo.Hours()/24)
			}
			
			content.WriteString(fmt.Sprintf("- [`%s`](../../commit/%s) %s *(%s)*\n", 
				commit.ShortHash, commit.Hash, commit.Subject, timeStr))
		}
		content.WriteString("\n")
	}
	
	// Learning insights
	if len(learning.ActiveInsights) > 0 {
		content.WriteString("### Active Insights\n\n")
		
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
			content.WriteString(fmt.Sprintf("🚨 **%d Critical** issues requiring immediate attention\n\n", criticalCount))
		}
		if highCount > 0 {
			content.WriteString(fmt.Sprintf("⚡ **%d High Priority** optimizations identified\n\n", highCount))
		}
		
		content.WriteString(fmt.Sprintf("📊 Total insights: %d\n\n", len(learning.ActiveInsights)))
	}
	
	content.WriteString("*Activity updated automatically*\n")
	content.WriteString("<!-- END AUTO-GENERATED ACTIVITY -->\n\n")
}

// generatePerformanceStats creates performance statistics section
func (r *ReadmeGenerator) generatePerformanceStats(content *strings.Builder, learning *LearningSnapshot) {
	content.WriteString("## ⚡ Performance Statistics\n\n")
	content.WriteString("<!-- AUTO-GENERATED PERFORMANCE -->\n")
	
	content.WriteString("| Metric | Value | Status |\n")
	content.WriteString("|--------|-------|--------|\n")
	
	// Success rate
	successRate := learning.SuccessRate * 100
	successStatus := "🟢 Excellent"
	if successRate < 95 {
		successStatus = "🟡 Good"
	}
	if successRate < 90 {
		successStatus = "🔴 Needs Improvement"
	}
	content.WriteString(fmt.Sprintf("| Success Rate | %.1f%% | %s |\n", successRate, successStatus))
	
	// Average latency
	if learning.AvgLatency > 0 {
		latencyMs := float64(learning.AvgLatency) / float64(time.Millisecond)
		latencyStatus := "🟢 Fast"
		if latencyMs > 100 {
			latencyStatus = "🟡 Good"
		}
		if latencyMs > 500 {
			latencyStatus = "🔴 Slow"
		}
		content.WriteString(fmt.Sprintf("| Avg Latency | %.1fms | %s |\n", latencyMs, latencyStatus))
	}
	
	// Total executions
	content.WriteString(fmt.Sprintf("| Total Executions | %d | 📊 Tracking |\n", learning.TotalExecutions))
	
	// Active tools
	content.WriteString(fmt.Sprintf("| Active Tools | %d | 🔧 Running |\n", len(learning.TopTools)))
	
	content.WriteString("\n*Statistics updated in real-time*\n")
	content.WriteString("<!-- END AUTO-GENERATED PERFORMANCE -->\n\n")
}

// generateInstallation creates installation section (fallback)
func (r *ReadmeGenerator) generateInstallation(content *strings.Builder) {
	content.WriteString("## 📦 Installation\n\n")
	content.WriteString("### Prerequisites\n\n")
	content.WriteString("- Go 1.21 or higher\n")
	content.WriteString("- Git\n\n")
	content.WriteString("### From Source\n\n")
	content.WriteString("```bash\n")
	content.WriteString("git clone https://github.com/kiransth77/aionmcp.git\n")
	content.WriteString("cd aionmcp\n")
	content.WriteString("go mod download\n")
	content.WriteString("go build -o bin/aionmcp cmd/server/main.go\n")
	content.WriteString("```\n\n")
}

// generateUsage creates usage section (fallback)
func (r *ReadmeGenerator) generateUsage(content *strings.Builder) {
	content.WriteString("## 📚 Usage\n\n")
	content.WriteString("### Basic Usage\n\n")
	content.WriteString("```bash\n")
	content.WriteString("# Start the server\n")
	content.WriteString("./bin/aionmcp\n\n")
	content.WriteString("# With custom configuration\n")
	content.WriteString("./bin/aionmcp --config config.yaml\n\n")
	content.WriteString("# Enable debug logging\n")
	content.WriteString("AIONMCP_LOG_LEVEL=debug ./bin/aionmcp\n")
	content.WriteString("```\n\n")
	content.WriteString("### API Endpoints\n\n")
	content.WriteString("- `GET /api/v1/tools` - List available tools\n")
	content.WriteString("- `POST /api/v1/tools/{tool}/execute` - Execute a tool\n")
	content.WriteString("- `GET /api/v1/learning/stats` - Learning statistics\n")
	content.WriteString("- `GET /api/v1/learning/insights` - System insights\n\n")
}

// generateDevelopment creates development section (fallback)
func (r *ReadmeGenerator) generateDevelopment(content *strings.Builder) {
	content.WriteString("## 🛠️ Development\n\n")
	content.WriteString("### Local Development\n\n")
	content.WriteString("```bash\n")
	content.WriteString("# Run tests\n")
	content.WriteString("go test ./...\n\n")
	content.WriteString("# Run with hot reload\n")
	content.WriteString("go run cmd/server/main.go\n\n")
	content.WriteString("# Build for production\n")
	content.WriteString("go build -ldflags \"-s -w\" -o bin/aionmcp cmd/server/main.go\n")
	content.WriteString("```\n\n")
}

// generateContributing creates contributing section (fallback)
func (r *ReadmeGenerator) generateContributing(content *strings.Builder) {
	content.WriteString("## 🤝 Contributing\n\n")
	content.WriteString("Contributions are welcome! Please feel free to submit a Pull Request.\n\n")
	content.WriteString("### Development Process\n\n")
	content.WriteString("1. Fork the repository\n")
	content.WriteString("2. Create a feature branch\n")
	content.WriteString("3. Make your changes\n")
	content.WriteString("4. Add tests\n")
	content.WriteString("5. Submit a pull request\n\n")
}

// generateLicense creates license section (fallback)
func (r *ReadmeGenerator) generateLicense(content *strings.Builder) {
	content.WriteString("## 📄 License\n\n")
	content.WriteString("This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.\n\n")
}

// generateFooter creates footer
func (r *ReadmeGenerator) generateFooter(content *strings.Builder) {
	content.WriteString("---\n\n")
	content.WriteString(fmt.Sprintf("*README last updated: %s*\n", time.Now().Format("2006-01-02 15:04:05 MST")))
	content.WriteString("\n*This README is automatically updated with current project status and metrics.*\n")
}

// Helper functions
func (r *ReadmeGenerator) calculateHealthScore(learning *LearningSnapshot) int {
	score := 100
	
	if learning.SuccessRate < 1.0 {
		score -= int((1.0-learning.SuccessRate)*50)
	}
	
	if learning.AvgLatency > 0 {
		latencyMs := float64(learning.AvgLatency) / float64(time.Millisecond)
		if latencyMs > 1000 {
			score -= 20
		} else if latencyMs > 500 {
			score -= 10
		}
	}
	
	for _, insight := range learning.ActiveInsights {
		if insight.Priority == "critical" {
			score -= 15
		} else if insight.Priority == "high" {
			score -= 5
		}
	}
	
	if score < 0 {
		score = 0
	}
	
	return score
}

func (r *ReadmeGenerator) getHealthStatus(score int) string {
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
func (r *ReadmeGenerator) writeToFile(outputPath, content string) error {
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	return nil
}