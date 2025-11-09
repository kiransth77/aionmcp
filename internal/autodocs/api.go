package autodocs

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// APIHandler handles HTTP requests for documentation operations
type APIHandler struct {
	engine DocumentEngine
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(engine DocumentEngine) *APIHandler {
	return &APIHandler{
		engine: engine,
	}
}

// RegisterRoutes registers documentation API routes
func (h *APIHandler) RegisterRoutes(router *gin.Engine) {
	docs := router.Group("/api/v1/docs")
	{
		// Documentation generation
		docs.POST("/generate", h.GenerateDocument)
		docs.POST("/generate/all", h.GenerateAllDocuments)
		docs.POST("/generate/daily", h.GenerateDaily)
		docs.POST("/generate/weekly", h.GenerateWeekly)

		// Generation history and status
		docs.GET("/history", h.GetGenerationHistory)
		docs.GET("/stats", h.GetStats)

		// Scheduled generation
		docs.POST("/schedule", h.ScheduleGeneration)
		docs.GET("/schedule", h.GetScheduledJobs)
		docs.DELETE("/schedule/:jobId", h.CancelScheduledJob)
		docs.POST("/schedule/process", h.ProcessScheduledJobs)

		// Health and status
		docs.GET("/health", h.GetDocumentationHealth)
		docs.GET("/types", h.GetSupportedTypes)
	}
}

// GenerateDocument generates a specific document
func (h *APIHandler) GenerateDocument(c *gin.Context) {
	var request GenerationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate the request
	engine, ok := h.engine.(*Engine)
	if ok {
		if err := engine.ValidateRequest(request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid generation request",
				"details": err.Error(),
			})
			return
		}
	}

	result, err := h.engine.Generate(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Generation failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

// GenerateAllDocuments generates all supported document types
func (h *APIHandler) GenerateAllDocuments(c *gin.Context) {
	results, err := h.engine.GenerateAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Generation failed",
			"details": err.Error(),
		})
		return
	}

	// Count successes and failures
	successCount := 0
	failureCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
		} else {
			failureCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"results":      results,
		"total":        len(results),
		"successful":   successCount,
		"failed":       failureCount,
		"generated_at": time.Now(),
	})
}

// GenerateDaily generates daily documentation (reflection + README update)
func (h *APIHandler) GenerateDaily(c *gin.Context) {
	// Cast to concrete type to access GenerateDaily method
	engine, ok := h.engine.(*Engine)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Engine does not support daily generation",
		})
		return
	}

	results, err := engine.GenerateDaily()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Daily generation failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results":      results,
		"type":         "daily",
		"generated_at": time.Now(),
	})
}

// GenerateWeekly generates weekly documentation (changelog)
func (h *APIHandler) GenerateWeekly(c *gin.Context) {
	// Cast to concrete type to access GenerateWeekly method
	engine, ok := h.engine.(*Engine)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Engine does not support weekly generation",
		})
		return
	}

	results, err := engine.GenerateWeekly()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Weekly generation failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results":      results,
		"type":         "weekly",
		"generated_at": time.Now(),
	})
}

// GetGenerationHistory returns recent generation history
func (h *APIHandler) GetGenerationHistory(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	history, err := h.engine.GetGenerationHistory(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get generation history",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"history": history,
		"limit":   limit,
		"count":   len(history),
	})
}

// GetStats returns documentation engine statistics
func (h *APIHandler) GetStats(c *gin.Context) {
	// Cast to concrete type to access GetStats method
	engine, ok := h.engine.(*Engine)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Engine does not support stats",
		})
		return
	}

	stats := engine.GetStats()

	c.JSON(http.StatusOK, gin.H{
		"stats":      stats,
		"updated_at": time.Now(),
	})
}

// ScheduleGeneration schedules automatic document generation
func (h *APIHandler) ScheduleGeneration(c *gin.Context) {
	var request struct {
		DocumentType DocumentType `json:"document_type" binding:"required"`
		Schedule     string       `json:"schedule" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	err := h.engine.ScheduleGeneration(request.DocumentType, request.Schedule)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to schedule generation",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Generation scheduled successfully",
		"document_type": request.DocumentType,
		"schedule":      request.Schedule,
		"scheduled_at":  time.Now(),
	})
}

// GetScheduledJobs returns all scheduled generation jobs
func (h *APIHandler) GetScheduledJobs(c *gin.Context) {
	// Cast to concrete type to access GetScheduledJobs method
	engine, ok := h.engine.(*Engine)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Engine does not support scheduled jobs",
		})
		return
	}

	jobs := engine.GetScheduledJobs()

	c.JSON(http.StatusOK, gin.H{
		"jobs":  jobs,
		"count": len(jobs),
	})
}

// CancelScheduledJob cancels a scheduled generation job
func (h *APIHandler) CancelScheduledJob(c *gin.Context) {
	jobID := c.Param("jobId")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Job ID is required",
		})
		return
	}

	// Cast to concrete type to access CancelScheduledJob method
	engine, ok := h.engine.(*Engine)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Engine does not support job cancellation",
		})
		return
	}

	err := engine.CancelScheduledJob(jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Failed to cancel job",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Job cancelled successfully",
		"job_id":       jobID,
		"cancelled_at": time.Now(),
	})
}

// ProcessScheduledJobs manually triggers processing of scheduled jobs
func (h *APIHandler) ProcessScheduledJobs(c *gin.Context) {
	// Cast to concrete type to access ProcessScheduledJobs method
	engine, ok := h.engine.(*Engine)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Engine does not support scheduled job processing",
		})
		return
	}

	err := engine.ProcessScheduledJobs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process scheduled jobs",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Scheduled jobs processed successfully",
		"processed_at": time.Now(),
	})
}

// GetDocumentationHealth returns documentation system health status
func (h *APIHandler) GetDocumentationHealth(c *gin.Context) {
	// Cast to concrete type to access additional methods
	engine, ok := h.engine.(*Engine)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Engine health check not supported",
		})
		return
	}

	stats := engine.GetStats()

	// Determine health status
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"components": map[string]interface{}{
			"generators": map[string]interface{}{
				"status": "healthy",
				"count":  stats["registered_generators"],
			},
			"history": map[string]interface{}{
				"status":            "healthy",
				"total_generations": stats["total_generations"],
				"success_rate":      stats["success_rate"],
			},
			"scheduler": map[string]interface{}{
				"status":      "healthy",
				"active_jobs": stats["active_scheduled_jobs"],
			},
		},
	}

	// Check if success rate is below threshold
	if successRate, ok := stats["success_rate"].(float64); ok && successRate < 0.9 {
		health["status"] = "degraded"
		health["components"].(map[string]interface{})["history"].(map[string]interface{})["status"] = "degraded"
		health["message"] = "Documentation generation success rate is below 90%"
	}

	c.JSON(http.StatusOK, health)
}

// GetSupportedTypes returns supported document types
func (h *APIHandler) GetSupportedTypes(c *gin.Context) {
	types := []DocumentType{
		DocumentTypeChangelog,
		DocumentTypeReflection,
		DocumentTypeReadme,
		DocumentTypeArchitecture,
	}

	typeInfo := make(map[DocumentType]interface{})
	for _, docType := range types {
		typeInfo[docType] = map[string]interface{}{
			"supported_formats": []string{"markdown"},
			"auto_scheduling":   docType != DocumentTypeArchitecture,
			"description":       h.getTypeDescription(docType),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"supported_types": types,
		"type_info":       typeInfo,
		"total_types":     len(types),
	})
}

// Helper function to get type descriptions
func (h *APIHandler) getTypeDescription(docType DocumentType) string {
	switch docType {
	case DocumentTypeChangelog:
		return "Automatically generated changelog from git commits with categorization"
	case DocumentTypeReflection:
		return "Daily reflection documents with learning insights and system health"
	case DocumentTypeReadme:
		return "Auto-updating README with current project status and metrics"
	case DocumentTypeArchitecture:
		return "Architecture documentation with system overview and components"
	default:
		return "Custom document type"
	}
}

// MiddlewareRequestLogging logs API requests for documentation operations
func (h *APIHandler) MiddlewareRequestLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Log request details
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		statusCode := c.Writer.Status()

		// Only log documentation API requests
		if len(path) > 12 && path[:12] == "/api/v1/docs" {
			fmt.Printf("[DOCS-API] %s %s %d %v %s\n",
				method, path, statusCode, latency, clientIP)
		}
	}
}

// WebhookHandler handles webhook requests for automatic documentation generation
func (h *APIHandler) WebhookHandler(c *gin.Context) {
	var payload struct {
		Event  string                 `json:"event"`
		Data   map[string]interface{} `json:"data"`
		Source string                 `json:"source"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid webhook payload",
		})
		return
	}

	// Handle different webhook events
	switch payload.Event {
	case "git.push":
		// Generate changelog and update README
		h.handleGitPushWebhook(c, payload.Data)
	case "learning.insight":
		// Generate reflection document
		h.handleLearningInsightWebhook(c, payload.Data)
	case "scheduled.generation":
		// Process scheduled jobs
		h.ProcessScheduledJobs(c)
	default:
		c.JSON(http.StatusOK, gin.H{
			"message": "Webhook received but no action taken",
			"event":   payload.Event,
		})
	}
}

// handleGitPushWebhook handles git push webhook events
func (h *APIHandler) handleGitPushWebhook(c *gin.Context, data map[string]interface{}) {
	// Generate updated README and changelog
	requests := []GenerationRequest{
		{
			Type:        DocumentTypeReadme,
			IncludeData: true,
			Format:      "markdown",
		},
		{
			Type: DocumentTypeChangelog,
			DateRange: &DateRange{
				StartDate: time.Now().AddDate(0, 0, -7), // Last week
				EndDate:   time.Now(),
			},
			IncludeData: true,
			Format:      "markdown",
		},
	}

	var results []GenerationResult
	for _, request := range requests {
		result, err := h.engine.Generate(request)
		if err != nil {
			results = append(results, GenerationResult{
				Type:    request.Type,
				Success: false,
				Error:   err.Error(),
			})
		} else {
			results = append(results, *result)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Git push webhook processed",
		"results": results,
	})
}

// handleLearningInsightWebhook handles learning insight webhook events
func (h *APIHandler) handleLearningInsightWebhook(c *gin.Context, data map[string]interface{}) {
	// Generate reflection document
	request := GenerationRequest{
		Type:        DocumentTypeReflection,
		IncludeData: true,
		Format:      "markdown",
	}

	result, err := h.engine.Generate(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate reflection",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Learning insight webhook processed",
		"result":  result,
	})
}
