package core

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/aionmcp/aionmcp/internal/selflearn"
	"github.com/aionmcp/aionmcp/pkg/agent"
	agentpb "github.com/aionmcp/aionmcp/pkg/agent/proto"
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
	agentServer     *agent.AgentServer
	agentAPI        *agent.AgentAPI
	learningEngine  *selflearn.Engine
	shutdown        chan struct{}
	wg              sync.WaitGroup
	serverCtx       context.Context // Server-scoped context for background operations
	cancelFunc      context.CancelFunc
	startTime       time.Time // Track server startup time for uptime calculation
}

// NewServer creates a new AionMCP server instance
func NewServer(logger *zap.Logger) (*Server, error) {
	logger.Info("NewServer: Starting initialization")

	// Initialize tool registry
	registry := NewToolRegistry(logger)
	logger.Info("NewServer: Tool registry created")

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

	// Initialize agent server and API
	agentServer := agent.NewAgentServer(logger, registry)
	agentAPI := agent.NewAgentAPI(logger, registry, agentServer)
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
	logger.Info("NewServer: Creating learning storage", zap.String("path", storagePath))
	learningStorage, err := selflearn.NewBoltStorage(storagePath, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create learning storage: %w", err)
	}
	logger.Info("NewServer: Learning storage created")

	// Create learning engine (ensure storage cleanup on error)
	learningEngine := selflearn.NewEngine(learningConfig, learningStorage, logger)
	if learningEngine == nil {
		learningStorage.Close()
		return nil, fmt.Errorf("failed to create learning engine")
	}

	// Create HTTP server with Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Add a simple root endpoint for health checks
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "AionMCP server is running"})
	})

	// Add request logging middleware
	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()

		// Skip logging for the root health check to avoid noise
		if c.Request.URL.Path == "/" {
			return
		}

		logger.Info("HTTP request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", time.Since(start)),
		)
	})

	// Create server-scoped context for background operations
	serverCtx, cancelFunc := context.WithCancel(context.Background())

	// Setup HTTP routes
	setupHTTPRoutes(router, registry, importerManager, fileWatcher, agentAPI, agentServer, learningEngine, logger, serverCtx, time.Now())

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("server.port")),
		Handler: router,
	}
	logger.Info("HTTP server configured", zap.String("addr", httpServer.Addr))

	// Create gRPC server and register agent service
	grpcServer := grpc.NewServer()
	agentpb.RegisterAgentServiceServer(grpcServer, agentServer)

	return &Server{
		logger:          logger,
		httpServer:      httpServer,
		grpcServer:      grpcServer,
		toolRegistry:    registry,
		importerManager: importerManager,
		fileWatcher:     fileWatcher,
		agentServer:     agentServer,
		agentAPI:        agentAPI,
		learningEngine:  learningEngine,
		shutdown:        make(chan struct{}),
		serverCtx:       serverCtx,
		cancelFunc:      cancelFunc,
		startTime:       time.Now(),
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
		s.logger.Info("HTTP server listening", zap.String("addr", s.httpServer.Addr))
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

	// Give servers a moment to start
	time.Sleep(100 * time.Millisecond)

	// Wait for shutdown signal
	<-ctx.Done()
	s.logger.Info("Shutting down AionMCP server...")

	// Cancel server-scoped context to stop background operations
	s.cancelFunc()

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
func setupHTTPRoutes(router *gin.Engine, registry *ToolRegistry, importerManager *importer.ImporterManager, fileWatcher *importer.FileWatcher, agentAPI *agent.AgentAPI, agentServer *agent.AgentServer, learningEngine *selflearn.Engine, logger *zap.Logger, serverCtx context.Context, startTime time.Time) {
	api := router.Group("/api/v1")

	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"version":   "0.1.0",
			"iteration": "4",
		})
	})

	// Agent integration routes
	agentAPI.RegisterRoutes(api)

	// Register a mock agent for testing (development endpoint) - at API root level
	api.POST("/register-copilot-agent", func(c *gin.Context) {
		ctx := c.Request.Context()

		// Register GitHub Copilot mock agent
		req := &agentpb.RegisterAgentRequest{
			AgentId:      "github-copilot",
			AgentName:    "GitHub Copilot",
			AgentVersion: "1.0.0",
			Capabilities: &agentpb.AgentCapabilities{
				SupportedProtocols:      []string{"mcp/1.0"},
				SupportedToolTypes:      []string{"openapi", "graphql"},
				SupportsStreaming:       true,
				SupportsAsyncInvocation: true,
				MaxConcurrentTools:      10,
				PreferredFormats:        []string{"json"},
			},
			Metadata: map[string]string{
				"type":     "ai-agent",
				"provider": "github",
				"llm":      "gpt-4",
			},
			SessionTimeoutSeconds: 3600, // 1 hour timeout
		}

		resp, err := agentServer.RegisterAgent(ctx, req)
		if err != nil {
			logger.Error("Failed to register mock Copilot agent", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		logger.Info("Mock Copilot agent registered successfully",
			zap.String("session_id", resp.SessionId))

		c.JSON(http.StatusOK, gin.H{
			"session_id": resp.SessionId,
			"agent_id":   "github-copilot",
			"agent_name": "GitHub Copilot",
			"expires_at": resp.ExpiresAtUnix,
			"message":    "GitHub Copilot agent registered successfully",
		})
	})

	// List available tools
	api.GET("/tools", func(c *gin.Context) {
		tools := registry.ListTools()
		c.JSON(http.StatusOK, gin.H{
			"protocol": viper.GetString("mcp.protocol_version"),
			"tools":    tools,
		})
	})

	// Tool invocation endpoint
	api.POST("/tools/:name/execute", func(c *gin.Context) {
		toolName := c.Param("name")
		startTime := time.Now()

		var request struct {
			Args map[string]interface{} `json:"args"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body, expected 'args' field"})
			return
		}

		// Get tool from registry
		tool, err := registry.Get(toolName)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("tool not found: %s", toolName)})
			return
		}

		// Execute tool and measure duration
		result, err := tool.Execute(request.Args)
		duration := time.Since(startTime)

		// Record execution for learning (async, non-blocking)
		go func(ctx context.Context, engine *selflearn.Engine, log *zap.Logger, tn string, req, res interface{}, execErr error, dur time.Duration) {
			metadata := tool.Metadata()
			sourceType := "builtin"
			if metadata.Source != "" {
				sourceType = metadata.Source
			}
			if recordErr := engine.RecordExecution(ctx, tn, sourceType, req, res, execErr, dur); recordErr != nil {
				log.Warn("Failed to record execution for learning", zap.String("tool", tn), zap.Error(recordErr))
			}
		}(serverCtx, learningEngine, logger, toolName, request.Args, result, err, duration)

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

	// Learning statistics
	api.GET("/stats", func(c *gin.Context) {
		stats, err := learningEngine.GetStats(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get learning stats"})
			return
		}
		c.JSON(http.StatusOK, stats)
	})

	// Server statistics endpoint for VS Code extension
	api.GET("/server-stats", func(c *gin.Context) {
		uptime := time.Since(startTime).Seconds()
		toolCount := registry.Count()

		// Get learning stats for execution count and success rate
		learningStats, err := learningEngine.GetStats(c.Request.Context())
		executionCount := int64(0)
		successRate := 0.0
		if err == nil {
			executionCount = learningStats.TotalExecutions
			successRate = learningStats.SuccessRate
		}

		// Get agent count from agent server
		agentCount := agentServer.GetAgentCount()

		c.JSON(http.StatusOK, gin.H{
			"uptime":         uptime,
			"toolCount":      toolCount,
			"agentCount":     agentCount,
			"executionCount": executionCount,
			"successRate":    successRate,
		})
	})

	// Import API specification endpoint
	api.POST("/import-spec", func(c *gin.Context) {
		var request struct {
			SpecType string `json:"spec_type" binding:"required"` // "openapi", "graphql", "asyncapi"
			Path     string `json:"path" binding:"required"`      // File path or URL
			Name     string `json:"name"`                         // Optional name for the spec
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
			return
		}

		// Validate spec type
		var specType importer.SpecType
		switch request.SpecType {
		case "openapi":
			specType = importer.SpecTypeOpenAPI
		case "graphql":
			specType = importer.SpecTypeGraphQL
		case "asyncapi":
			specType = importer.SpecTypeAsyncAPI
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spec type: " + request.SpecType})
			return
		}

		// Create spec source
		source := importer.SpecSource{
			ID:   fmt.Sprintf("%s-%d", string(specType), time.Now().Unix()),
			Type: specType,
			Path: request.Path,
			Name: request.Name,
		}

		// Import the specification
		result, err := importerManager.ImportSpec(c.Request.Context(), source)
		if err != nil {
			logger.Error("Failed to import specification",
				zap.String("spec_type", string(specType)),
				zap.String("path", request.Path),
				zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import specification: " + err.Error()})
			return
		}

		logger.Info("Specification imported successfully",
			zap.String("spec_type", string(specType)),
			zap.String("path", request.Path),
			zap.Int("tools_count", len(result.Tools)))

		// Record import for learning
		go func(ctx context.Context) {
			if recordErr := learningEngine.RecordExecution(ctx, "import_spec", string(specType), source, result, err, result.Duration); recordErr != nil {
				logger.Warn("Failed to record spec import for learning", zap.Error(recordErr))
			}
		}(serverCtx)

		c.JSON(http.StatusOK, gin.H{
			"source":   source,
			"tools":    result.Tools,
			"errors":   result.Errors,
			"warnings": result.Warnings,
			"duration": result.Duration,
		})
	})
}
