package autodocs

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"
)

const (
	// maxHistoryEntries is the maximum number of generation results to keep in history
	maxHistoryEntries = 100
)

// Engine implements the DocumentEngine interface
type Engine struct {
	generators  map[DocumentType]Generator
	dataSource  DataSource
	projectRoot string
	history     []GenerationResult
	historyMu   sync.RWMutex
	scheduledJobs map[string]*ScheduledJob
	mu          sync.RWMutex
}

// ScheduledJob represents a scheduled documentation generation job
type ScheduledJob struct {
	ID       string
	DocType  DocumentType
	Schedule string
	NextRun  time.Time
	Active   bool
}

// NewEngine creates a new documentation engine
func NewEngine(projectRoot string, dataSource DataSource) *Engine {
	engine := &Engine{
		generators:    make(map[DocumentType]Generator),
		dataSource:    dataSource,
		projectRoot:   projectRoot,
		history:       make([]GenerationResult, 0),
		scheduledJobs: make(map[string]*ScheduledJob),
	}
	
	// Register default generators
	engine.RegisterGenerator(NewChangelogGenerator(dataSource))
	engine.RegisterGenerator(NewReflectionGenerator(dataSource))
	engine.RegisterGenerator(NewReadmeGenerator(dataSource, projectRoot))
	
	return engine
}

// RegisterGenerator adds a new document generator
func (e *Engine) RegisterGenerator(generator Generator) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	supportedTypes := generator.GetSupportedTypes()
	for _, docType := range supportedTypes {
		e.generators[docType] = generator
	}
	
	return nil
}

// Generate creates a document using the appropriate generator
func (e *Engine) Generate(request GenerationRequest) (*GenerationResult, error) {
	e.mu.RLock()
	generator, exists := e.generators[request.Type]
	e.mu.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("no generator found for document type: %s", request.Type)
	}
	
	// Set default output path if not provided
	if request.OutputPath == "" {
		request.OutputPath = e.getDefaultOutputPath(request.Type)
	}
	
	// Set default format if not provided
	if request.Format == "" {
		request.Format = "markdown"
	}
	
	// Generate the document
	result, err := generator.Generate(request)
	if err != nil {
		return nil, fmt.Errorf("generation failed: %w", err)
	}
	
	// Add to history
	e.addToHistory(*result)
	
	return result, nil
}

// GenerateAll creates all supported document types
func (e *Engine) GenerateAll() ([]GenerationResult, error) {
	e.mu.RLock()
	docTypes := make([]DocumentType, 0, len(e.generators))
	for docType := range e.generators {
		docTypes = append(docTypes, docType)
	}
	e.mu.RUnlock()
	
	var results []GenerationResult
	var errors []error
	
	// Generate each document type
	for _, docType := range docTypes {
		request := GenerationRequest{
			Type:        docType,
			OutputPath:  e.getDefaultOutputPath(docType),
			Format:      "markdown",
			IncludeData: true,
		}
		
		// Set appropriate date range for different document types
		switch docType {
		case DocumentTypeChangelog:
			// Last 30 days for changelog
			request.DateRange = &DateRange{
				StartDate: time.Now().AddDate(0, 0, -30),
				EndDate:   time.Now(),
			}
		case DocumentTypeReflection:
			// Today for reflection
			today := time.Now()
			startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
			request.DateRange = &DateRange{
				StartDate: startOfDay,
				EndDate:   startOfDay.Add(24 * time.Hour),
			}
		}
		
		result, err := e.Generate(request)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to generate %s: %w", docType, err))
			// Create a failed result for tracking
			result = &GenerationResult{
				Type:        docType,
				Success:     false,
				Error:       err.Error(),
				GeneratedAt: time.Now(),
			}
		}
		
		results = append(results, *result)
	}
	
	// Return errors if any occurred
	if len(errors) > 0 {
		return results, fmt.Errorf("generation errors occurred: %v", errors)
	}
	
	return results, nil
}

// ScheduleGeneration sets up automatic document generation
func (e *Engine) ScheduleGeneration(docType DocumentType, schedule string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	// Validate that we have a generator for this document type
	if _, exists := e.generators[docType]; !exists {
		return fmt.Errorf("no generator found for document type: %s", docType)
	}
	
	// Parse schedule and calculate next run time
	nextRun, err := e.parseSchedule(schedule)
	if err != nil {
		return fmt.Errorf("invalid schedule format: %w", err)
	}
	
	// Create job ID
	jobID := fmt.Sprintf("%s_%s", docType, schedule)
	
	// Create scheduled job
	job := &ScheduledJob{
		ID:       jobID,
		DocType:  docType,
		Schedule: schedule,
		NextRun:  nextRun,
		Active:   true,
	}
	
	e.scheduledJobs[jobID] = job
	
	return nil
}

// GetGenerationHistory returns recent generation results
func (e *Engine) GetGenerationHistory(limit int) ([]GenerationResult, error) {
	e.historyMu.RLock()
	defer e.historyMu.RUnlock()
	
	// Return all if limit is 0 or negative
	if limit <= 0 {
		return e.history, nil
	}
	
	// Return last N results
	start := 0
	if len(e.history) > limit {
		start = len(e.history) - limit
	}
	
	return e.history[start:], nil
}

// GenerateDaily generates daily documentation (reflection + updated README)
func (e *Engine) GenerateDaily() ([]GenerationResult, error) {
	var results []GenerationResult
	
	// Generate daily reflection
	today := time.Now()
	reflectionDate := today.Format("2006-01-02")
	reflectionPath := filepath.Join(e.projectRoot, "docs", "reflections", reflectionDate+".md")
	
	reflectionRequest := GenerationRequest{
		Type:       DocumentTypeReflection,
		OutputPath: reflectionPath,
		DateRange: &DateRange{
			StartDate: time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()),
			EndDate:   today,
		},
		IncludeData: true,
		Format:      "markdown",
	}
	
	reflectionResult, err := e.Generate(reflectionRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to generate daily reflection: %w", err)
	}
	results = append(results, *reflectionResult)
	
	// Update README
	readmeRequest := GenerationRequest{
		Type:        DocumentTypeReadme,
		OutputPath:  filepath.Join(e.projectRoot, "README.md"),
		IncludeData: true,
		Format:      "markdown",
	}
	
	readmeResult, err := e.Generate(readmeRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to update README: %w", err)
	}
	results = append(results, *readmeResult)
	
	return results, nil
}

// GenerateWeekly generates weekly documentation (changelog + summary)
func (e *Engine) GenerateWeekly() ([]GenerationResult, error) {
	var results []GenerationResult
	
	// Generate weekly changelog
	weekAgo := time.Now().AddDate(0, 0, -7)
	now := time.Now()
	
	changelogRequest := GenerationRequest{
		Type:       DocumentTypeChangelog,
		OutputPath: filepath.Join(e.projectRoot, "docs", "changelog.md"),
		DateRange: &DateRange{
			StartDate: weekAgo,
			EndDate:   now,
		},
		IncludeData: true,
		Format:      "markdown",
	}
	
	changelogResult, err := e.Generate(changelogRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to generate weekly changelog: %w", err)
	}
	results = append(results, *changelogResult)
	
	return results, nil
}

// ProcessScheduledJobs runs any scheduled documentation generation jobs
func (e *Engine) ProcessScheduledJobs() error {
	e.mu.RLock()
	jobs := make([]*ScheduledJob, 0, len(e.scheduledJobs))
	for _, job := range e.scheduledJobs {
		if job.Active && time.Now().After(job.NextRun) {
			jobs = append(jobs, job)
		}
	}
	e.mu.RUnlock()
	
	for _, job := range jobs {
		// Generate the document
		request := GenerationRequest{
			Type:        job.DocType,
			OutputPath:  e.getDefaultOutputPath(job.DocType),
			IncludeData: true,
			Format:      "markdown",
		}
		
		// Set appropriate date range based on schedule
		switch job.Schedule {
		case "daily":
			today := time.Now()
			startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
			request.DateRange = &DateRange{
				StartDate: startOfDay,
				EndDate:   today,
			}
		case "weekly":
			request.DateRange = &DateRange{
				StartDate: time.Now().AddDate(0, 0, -7),
				EndDate:   time.Now(),
			}
		case "monthly":
			request.DateRange = &DateRange{
				StartDate: time.Now().AddDate(0, -1, 0),
				EndDate:   time.Now(),
			}
		}
		
		_, err := e.Generate(request)
		if err != nil {
			return fmt.Errorf("failed to process scheduled job %s: %w", job.ID, err)
		}
		
		// Update next run time
		e.mu.Lock()
		nextRun, err := e.parseSchedule(job.Schedule)
		if err == nil {
			job.NextRun = nextRun
		}
		e.mu.Unlock()
	}
	
	return nil
}

// GetScheduledJobs returns all scheduled jobs
func (e *Engine) GetScheduledJobs() []*ScheduledJob {
	e.mu.RLock()
	defer e.mu.RUnlock()
	
	jobs := make([]*ScheduledJob, 0, len(e.scheduledJobs))
	for _, job := range e.scheduledJobs {
		jobs = append(jobs, job)
	}
	
	return jobs
}

// CancelScheduledJob cancels a scheduled job
func (e *Engine) CancelScheduledJob(jobID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	job, exists := e.scheduledJobs[jobID]
	if !exists {
		return fmt.Errorf("job not found: %s", jobID)
	}
	
	job.Active = false
	return nil
}

// GetStats returns engine statistics
func (e *Engine) GetStats() map[string]interface{} {
	e.historyMu.RLock()
	e.mu.RLock()
	defer e.historyMu.RUnlock()
	defer e.mu.RUnlock()
	
	stats := map[string]interface{}{
		"registered_generators": len(e.generators),
		"total_generations":     len(e.history),
		"scheduled_jobs":        len(e.scheduledJobs),
	}
	
	// Count successful generations
	successCount := 0
	for _, result := range e.history {
		if result.Success {
			successCount++
		}
	}
	
	if len(e.history) > 0 {
		stats["success_rate"] = float64(successCount) / float64(len(e.history))
	} else {
		stats["success_rate"] = 0.0
	}
	
	// Recent generation statistics
	recent := make(map[DocumentType]int)
	cutoff := time.Now().AddDate(0, 0, -7) // Last 7 days
	
	for _, result := range e.history {
		if result.GeneratedAt.After(cutoff) {
			recent[result.Type]++
		}
	}
	
	stats["recent_generations"] = recent
	
	// Active scheduled jobs
	activeJobs := 0
	for _, job := range e.scheduledJobs {
		if job.Active {
			activeJobs++
		}
	}
	stats["active_scheduled_jobs"] = activeJobs
	
	return stats
}

// Helper methods

// getDefaultOutputPath returns the default output path for a document type
func (e *Engine) getDefaultOutputPath(docType DocumentType) string {
	switch docType {
	case DocumentTypeChangelog:
		return filepath.Join(e.projectRoot, "docs", "changelog.md")
	case DocumentTypeReflection:
		date := time.Now().Format("2006-01-02")
		return filepath.Join(e.projectRoot, "docs", "reflections", date+".md")
	case DocumentTypeReadme:
		return filepath.Join(e.projectRoot, "README.md")
	case DocumentTypeArchitecture:
		return filepath.Join(e.projectRoot, "docs", "architecture.md")
	default:
		return filepath.Join(e.projectRoot, "docs", string(docType)+".md")
	}
}

// parseSchedule parses a schedule string and returns the next run time
func (e *Engine) parseSchedule(schedule string) (time.Time, error) {
	now := time.Now()
	
	switch schedule {
	case "daily":
		// Next day at midnight
		tomorrow := now.AddDate(0, 0, 1)
		return time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location()), nil
	case "weekly":
		// Next week at midnight on Monday
		daysUntilMonday := (8 - int(now.Weekday())) % 7
		// If today is Monday, daysUntilMonday will be 0, so schedule for today
		if daysUntilMonday == 0 {
			daysUntilMonday = 7
		}
		nextMonday := now.AddDate(0, 0, daysUntilMonday)
		return time.Date(nextMonday.Year(), nextMonday.Month(), nextMonday.Day(), 0, 0, 0, 0, nextMonday.Location()), nil
	case "monthly":
		// Next month at midnight on the 1st
		nextMonth := now.AddDate(0, 1, 0)
		return time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, nextMonth.Location()), nil
	case "hourly":
		// Next hour
		return now.Add(time.Hour), nil
	default:
		return time.Time{}, fmt.Errorf("unsupported schedule: %s", schedule)
	}
}

// addToHistory adds a generation result to the history
func (e *Engine) addToHistory(result GenerationResult) {
	e.historyMu.Lock()
	defer e.historyMu.Unlock()
	
	e.history = append(e.history, result)
	
	// Keep only last maxHistoryEntries results
	if len(e.history) > maxHistoryEntries {
		e.history = e.history[len(e.history)-maxHistoryEntries:]
	}
}

// ValidateRequest validates a generation request
func (e *Engine) ValidateRequest(request GenerationRequest) error {
	e.mu.RLock()
	generator, exists := e.generators[request.Type]
	e.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("unsupported document type: %s", request.Type)
	}
	
	return generator.Validate(request)
}