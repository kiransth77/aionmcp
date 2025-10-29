package autodocs

import (
	"fmt"
	"os"
	"path/filepath"
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
