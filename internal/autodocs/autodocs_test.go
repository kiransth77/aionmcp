package autodocs

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"
)

// TestDocumentationSystem tests the core documentation functionality
func TestDocumentationSystem(t *testing.T) {
	// Setup
	projectRoot := "../../" // Relative to internal/autodocs
	dataSource := NewLearningDataSource(projectRoot, "")
	engine := NewEngine(projectRoot, dataSource)

	t.Run("Changelog Generation", func(t *testing.T) {
		request := GenerationRequest{
			Type:       DocumentTypeChangelog,
			OutputPath: filepath.Join("test_output", "changelog.md"),
			DateRange: &DateRange{
				StartDate: time.Now().AddDate(0, 0, -7),
				EndDate:   time.Now(),
			},
			IncludeData: true,
			Format:      "markdown",
		}

		result, err := engine.Generate(request)
		if err != nil {
			t.Fatalf("Changelog generation failed: %v", err)
		}

		if !result.Success {
			t.Fatalf("Changelog generation was not successful: %s", result.Error)
		}

		if result.ContentLength == 0 {
			t.Error("Generated changelog has no content")
		}

		t.Logf("✅ Changelog generated: %d bytes", result.ContentLength)
	})

	t.Run("Reflection Generation", func(t *testing.T) {
		request := GenerationRequest{
			Type:       DocumentTypeReflection,
			OutputPath: filepath.Join("test_output", "reflection.md"),
			DateRange: &DateRange{
				StartDate: time.Now().Truncate(24 * time.Hour),
				EndDate:   time.Now(),
			},
			IncludeData: true,
			Format:      "markdown",
		}

		result, err := engine.Generate(request)
		if err != nil {
			t.Fatalf("Reflection generation failed: %v", err)
		}

		if !result.Success {
			t.Fatalf("Reflection generation was not successful: %s", result.Error)
		}

		if result.ContentLength == 0 {
			t.Error("Generated reflection has no content")
		}

		t.Logf("✅ Reflection generated: %d bytes", result.ContentLength)
	})

	t.Run("README Generation", func(t *testing.T) {
		request := GenerationRequest{
			Type:        DocumentTypeReadme,
			OutputPath:  filepath.Join("test_output", "README.md"),
			IncludeData: true,
			Format:      "markdown",
		}

		result, err := engine.Generate(request)
		if err != nil {
			t.Fatalf("README generation failed: %v", err)
		}

		if !result.Success {
			t.Fatalf("README generation was not successful: %s", result.Error)
		}

		if result.ContentLength == 0 {
			t.Error("Generated README has no content")
		}

		t.Logf("✅ README generated: %d bytes", result.ContentLength)
	})

	t.Run("Generate All Documents", func(t *testing.T) {
		results, err := engine.GenerateAll()
		if err != nil {
			t.Fatalf("Generate all failed: %v", err)
		}

		if len(results) == 0 {
			t.Fatal("No documents were generated")
		}

		successCount := 0
		for _, result := range results {
			if result.Success {
				successCount++
			}
		}

		t.Logf("✅ Generated %d documents, %d successful", len(results), successCount)
	})

	t.Run("Engine Statistics", func(t *testing.T) {
		stats := engine.GetStats()

		if stats["registered_generators"] == nil {
			t.Error("No generators registered")
		}

		t.Logf("✅ Engine stats: %v", stats)
	})

	t.Run("Generation History", func(t *testing.T) {
		history, err := engine.GetGenerationHistory(10)
		if err != nil {
			t.Fatalf("Failed to get generation history: %v", err)
		}

		if len(history) == 0 {
			t.Error("No generation history found")
		}

		t.Logf("✅ Generation history: %d entries", len(history))
	})
}

// TestDataSources tests data source functionality
func TestDataSources(t *testing.T) {
	projectRoot := "../../"
	
	t.Run("Git Data Source", func(t *testing.T) {
		gitDS := NewGitDataSource(projectRoot)
		
		// Test project info
		info, err := gitDS.GetProjectInfo()
		if err != nil {
			t.Fatalf("Failed to get project info: %v", err)
		}
		
		if info["current_branch"] == nil {
			t.Error("Current branch not found in project info")
		}
		
		t.Logf("✅ Project info: %v", info["current_branch"])
		
		// Test commits
		dateRange := DateRange{
			StartDate: time.Now().AddDate(0, 0, -7),
			EndDate:   time.Now(),
		}
		
		commits, err := gitDS.GetCommits(dateRange)
		if err != nil {
			t.Fatalf("Failed to get commits: %v", err)
		}
		
		t.Logf("✅ Found %d commits in last 7 days", len(commits))
	})

	t.Run("Learning Data Source", func(t *testing.T) {
		learningDS := NewLearningDataSource(projectRoot, "")
		
		// Test learning snapshot (should return mock data)
		snapshot, err := learningDS.GetLearningSnapshot()
		if err != nil {
			t.Fatalf("Failed to get learning snapshot: %v", err)
		}
		
		if snapshot.TotalExecutions == 0 {
			t.Error("Learning snapshot has no executions")
		}
		
		t.Logf("✅ Learning snapshot: %d executions, %.1f%% success rate", 
			snapshot.TotalExecutions, snapshot.SuccessRate*100)
	})
}

// BenchmarkDocumentGeneration benchmarks document generation performance
func BenchmarkDocumentGeneration(b *testing.B) {
	projectRoot := "../../"
	dataSource := NewLearningDataSource(projectRoot, "")
	engine := NewEngine(projectRoot, dataSource)

	request := GenerationRequest{
		Type:       DocumentTypeChangelog,
		OutputPath: "/tmp/benchmark_changelog.md",
		DateRange: &DateRange{
			StartDate: time.Now().AddDate(0, 0, -7),
			EndDate:   time.Now(),
		},
		IncludeData: true,
		Format:      "markdown",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Generate(request)
		if err != nil {
			b.Fatalf("Generation failed: %v", err)
		}
	}
}

// Example function for documentation
func ExampleEngine_Generate() {
	// Create a documentation engine
	projectRoot := "."
	dataSource := NewLearningDataSource(projectRoot, "")
	engine := NewEngine(projectRoot, dataSource)

	// Generate a changelog
	request := GenerationRequest{
		Type:       DocumentTypeChangelog,
		OutputPath: "docs/changelog.md",
		DateRange: &DateRange{
			StartDate: time.Now().AddDate(0, 0, -30), // Last 30 days
			EndDate:   time.Now(),
		},
		IncludeData: true,
		Format:      "markdown",
	}

	result, err := engine.Generate(request)
	if err != nil {
		fmt.Printf("Generation failed: %v\n", err)
		return
	}

	fmt.Printf("Changelog generated: %s (%d bytes)\n", 
		result.OutputPath, result.ContentLength)
}