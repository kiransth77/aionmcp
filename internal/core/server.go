package core

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/aionmcp/aionmcp/internal/selflearn"
	"github.com/aionmcp/aionmcp/pkg/importer"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Server represents the main AionMCP server
type Server struct {
	logger          *zap.Logger
	httpServer      *http.Server
	grpcServer      *grpc.Server
	toolRegistry    *ToolRegistry
	importerManager *importer.ImporterManager
	fileWatcher     *importer.FileWatcher
	learningEngine  *selflearn.Engine
	shutdown        chan struct{}
	wg              sync.WaitGroup
}

// NewServer creates a new AionMCP server instance
func NewServer(logger *zap.Logger) (*Server, error) {
	// Initialize tool registry
	registry := NewToolRegistry(logger)

	// Initialize importer manager
	importerManager := importer.NewImporterManager(registry)
	
	// Register importers
	importerManager.RegisterImporter(importer.NewOpenAPIImporter())
	importerManager.RegisterImporter(importer.NewGraphQLImporter())
	importerManager.RegisterImporter(importer.NewAsyncAPIImporter())

	// Initialize file watcher
	fileWatcher, err := importer.NewFileWatcher(importerManager, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}

	// Initialize self-learning engine
	learningConfig := selflearn.DefaultCollectionConfig()
	learningConfig.Enabled = viper.GetBool("learning.enabled")
	if learningConfig.Enabled {
		if sampleRate := viper.GetFloat64("learning.sample_rate"); sampleRate > 0 {
			learningConfig.SampleRate = sampleRate
		}
		if retentionDays := viper.GetInt("learning.retention_days"); retentionDays > 0 {
			learningConfig.RetentionPeriod = time.Duration(retentionDays) * 24 * time.Hour
		}
	}

	// Create learning storage
	storagePath := viper.GetString("storage.path")
	if storagePath == "" {
		storagePath = "./data/aionmcp.db"
	}
	learningStorage, err := selflearn.NewBoltStorage(storagePath, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create learning storage: %w", err)
	}

	// Create learning engine
	learningEngine := selflearn.NewEngine(learningConfig, learningStorage, logger)

	// Create HTTP server with Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	
	// Add request logging middleware
	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		
		logger.Info("HTTP request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", time.Since(start)),
		)
	})

	// Setup HTTP routes
	setupHTTPRoutes(router, registry, importerManager, fileWatcher, learningEngine, logger)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("server.port")),
		Handler: router,
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()
	// TODO: Register gRPC services in iteration 4

	return &Server{
		logger:          logger,
		httpServer:      httpServer,
		grpcServer:      grpcServer,
		toolRegistry:    registry,
		importerManager: importerManager,
		fileWatcher:     fileWatcher,
		learningEngine:  learningEngine,
		shutdown:        make(chan struct{}),
	}, nil
}

// Run starts the server and blocks until context is cancelled
func (s *Server) Run(ctx context.Context) error {
	s.logger.Info("Starting AionMCP server",
		zap.String("http_port", s.httpServer.Addr),
		zap.Int("grpc_port", viper.GetInt("server.grpc_port")))

	// Start HTTP server
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("HTTP server failed", zap.Error(err))
		}
	}()

	// Start gRPC server
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", viper.GetInt("server.grpc_port")))
		if err != nil {
			s.logger.Error("Failed to listen on gRPC port", zap.Error(err))
			return
		}

		if err := s.grpcServer.Serve(lis); err != nil {
			s.logger.Error("gRPC server failed", zap.Error(err))
		}
	}()

	s.logger.Info("AionMCP server started successfully")

	// Wait for shutdown signal
	<-ctx.Done()
	s.logger.Info("Shutting down AionMCP server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		s.logger.Error("Failed to shutdown HTTP server", zap.Error(err))
	}

	// Shutdown gRPC server
	s.grpcServer.GracefulStop()

	// Stop file watcher
	s.fileWatcher.Stop()

	// Wait for all goroutines to finish
	s.wg.Wait()

	return nil
}

// setupHTTPRoutes configures HTTP API routes
func setupHTTPRoutes(router *gin.Engine, registry *ToolRegistry, importerManager *importer.ImporterManager, fileWatcher *importer.FileWatcher, learningEngine *selflearn.Engine, logger *zap.Logger) {
	api := router.Group("/api/v1")
	
	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"version":   "0.1.0",
			"iteration": "1",
		})
	})

	// MCP endpoints
	mcp := api.Group("/mcp")
	
	// List available tools
	mcp.GET("/tools", func(c *gin.Context) {
		tools := registry.ListTools()
		c.JSON(http.StatusOK, gin.H{
			"protocol": viper.GetString("mcp.protocol_version"),
			"tools":    tools,
		})
	})

	// Tool invocation endpoint
	mcp.POST("/tools/:name/invoke", func(c *gin.Context) {
		toolName := c.Param("name")
		startTime := time.Now()
		
		var request map[string]interface{}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		// Get tool from registry
		tool, err := registry.Get(toolName)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("tool not found: %s", toolName)})
			return
		}

		// Execute tool and measure duration
		result, err := tool.Execute(request)
		duration := time.Since(startTime)

		// Record execution for learning (async, non-blocking)
		// Capture err value and metadata before goroutine to avoid race condition
		execErr := err
		metadata := tool.Metadata()
		sourceType := "builtin"
		if metadata.Source != "" {
			sourceType = metadata.Source
		}
		
		go func(st string) {
			// Record the execution
			if recordErr := learningEngine.RecordExecution(
				context.Background(),
				toolName,
				st,
				request,
				result,
				execErr,
				duration,
			); recordErr != nil {
				logger.Warn("Failed to record execution for learning",
					zap.String("tool", toolName),
					zap.Error(recordErr))
			}
		}(sourceType)

		if err != nil {
			logger.Error("Tool execution failed",
				zap.String("tool", toolName),
				zap.Duration("duration", duration),
				zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		logger.Info("Tool executed successfully",
			zap.String("tool", toolName),
			zap.Duration("duration", duration))

		c.JSON(http.StatusOK, gin.H{
			"tool":   toolName,
			"result": result,
		})
	})

	// Importer management endpoints
	specs := api.Group("/specs")
	
	// List specification sources
	specs.GET("/", func(c *gin.Context) {
		sources := importerManager.ListSources()
		c.JSON(http.StatusOK, gin.H{
			"sources": sources,
		})
	})

	// Import a new specification
	specs.POST("/", func(c *gin.Context) {
		var req struct {
			ID          string            `json:"id" binding:"required"`
			Type        string            `json:"type" binding:"required"`
			Path        string            `json:"path" binding:"required"`
			Name        string            `json:"name"`
			Description string            `json:"description"`
			Metadata    map[string]string `json:"metadata"`
			EnableWatch bool              `json:"enable_watch"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create spec source
		source := importer.SpecSource{
			ID:          req.ID,
			Type:        importer.SpecType(req.Type),
			Path:        req.Path,
			Name:        req.Name,
			Description: req.Description,
			Metadata:    req.Metadata,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Import the specification
		result, err := importerManager.ImportSpec(c.Request.Context(), source)
		if err != nil {
			logger.Error("Failed to import specification",
				zap.String("source_id", req.ID),
				zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Enable file watching if requested
		if req.EnableWatch {
			if err := fileWatcher.WatchSpec(source); err != nil {
				logger.Warn("Failed to enable file watching",
					zap.String("source_id", req.ID),
					zap.Error(err))
				result.Warnings = append(result.Warnings, fmt.Sprintf("File watching could not be enabled: %v", err))
			}
		}

		logger.Info("Specification imported successfully",
			zap.String("source_id", req.ID),
			zap.String("type", req.Type),
			zap.Int("tools_count", len(result.Tools)))

		c.JSON(http.StatusCreated, gin.H{
			"result": result,
		})
	})

	// Get specification details
	specs.GET("/:id", func(c *gin.Context) {
		sourceID := c.Param("id")
		source, exists := importerManager.GetSource(sourceID)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "specification not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"source":      source,
			"is_watching": fileWatcher.IsWatching(sourceID),
		})
	})

	// Reload a specification
	specs.POST("/:id/reload", func(c *gin.Context) {
		sourceID := c.Param("id")
		
		result, err := importerManager.ReloadSpec(c.Request.Context(), sourceID)
		if err != nil {
			logger.Error("Failed to reload specification",
				zap.String("source_id", sourceID),
				zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		logger.Info("Specification reloaded successfully",
			zap.String("source_id", sourceID),
			zap.Int("tools_count", len(result.Tools)))

		c.JSON(http.StatusOK, gin.H{
			"result": result,
		})
	})

	// Remove a specification
	specs.DELETE("/:id", func(c *gin.Context) {
		sourceID := c.Param("id")
		
		// Stop watching if enabled
		if fileWatcher.IsWatching(sourceID) {
			if err := fileWatcher.UnwatchSpec(sourceID); err != nil {
				logger.Warn("Failed to stop watching specification",
					zap.String("source_id", sourceID),
					zap.Error(err))
			}
		}

		// Remove the specification
		if err := importerManager.RemoveSpec(c.Request.Context(), sourceID); err != nil {
			logger.Error("Failed to remove specification",
				zap.String("source_id", sourceID),
				zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		logger.Info("Specification removed successfully",
			zap.String("source_id", sourceID))

		c.JSON(http.StatusNoContent, nil)
	})

	// List supported specification types
	specs.GET("/types", func(c *gin.Context) {
		types := importerManager.GetSupportedTypes()
		c.JSON(http.StatusOK, gin.H{
			"supported_types": types,
		})
	})

	// Learning endpoints
	learning := api.Group("/learning")

	// Get overall learning statistics
	learning.GET("/stats", func(c *gin.Context) {
		stats, err := learningEngine.GetStats(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get learning stats"})
			return
		}
		c.JSON(http.StatusOK, stats)
	})

	// Get insights
	learning.GET("/insights", func(c *gin.Context) {
		insightType := c.Query("type")
		priority := c.Query("priority")
		
		var insights []selflearn.Insight
		var err error

		if priority != "" {
			insights, err = learningEngine.GetInsightsByPriority(c.Request.Context(), selflearn.Priority(priority), 50)
		} else {
			insights, err = learningEngine.GetInsights(c.Request.Context(), selflearn.InsightType(insightType), 50)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get insights"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"insights": insights})
	})

	// Get patterns
	learning.GET("/patterns", func(c *gin.Context) {
		patternType := c.Query("type")
		toolName := c.Query("tool")

		if toolName != "" {
			patterns, err := learningEngine.GetErrorPatterns(c.Request.Context(), toolName)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tool patterns"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"patterns": patterns})
			return
		}

		patterns, err := learningEngine.GetPatterns(c.Request.Context(), selflearn.PatternType(patternType), 50)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get patterns"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"patterns": patterns})
	})

	// Get tool-specific insights
	learning.GET("/tools/:name/insights", func(c *gin.Context) {
		toolName := c.Param("name")
		insights, err := learningEngine.GetToolInsights(c.Request.Context(), toolName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tool insights"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"tool_name": toolName, "insights": insights})
	})

	// Trigger manual analysis
	learning.POST("/analyze", func(c *gin.Context) {
		patterns, err := learningEngine.AnalyzePatterns(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze patterns"})
			return
		}

		insights, err := learningEngine.GenerateInsights(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate insights"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"patterns_found": len(patterns),
			"insights_generated": len(insights),
		})
	})

	// Get/update learning configuration
	learning.GET("/config", func(c *gin.Context) {
		config := learningEngine.GetConfig()
		c.JSON(http.StatusOK, config)
	})
}