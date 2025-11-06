package agent

import (
	"encoding/json"
	"net/http"
	"strconv"

	agentpb "github.com/aionmcp/aionmcp/pkg/agent/proto"
	"github.com/aionmcp/aionmcp/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AgentAPI provides REST endpoints for agent integration
type AgentAPI struct {
	logger      *zap.Logger
	registry    types.ToolRegistry
	agentServer *AgentServer
}

// NewAgentAPI creates a new AgentAPI instance
func NewAgentAPI(logger *zap.Logger, registry types.ToolRegistry, agentServer *AgentServer) *AgentAPI {
	return &AgentAPI{
		logger:      logger,
		registry:    registry,
		agentServer: agentServer,
	}
}

// RegisterRoutes adds agent API routes to the gin router
func (api *AgentAPI) RegisterRoutes(router *gin.RouterGroup) {
	agents := router.Group("/agents")
	
	// Agent session management
	agents.POST("/register", api.registerAgent)
	agents.DELETE("/:session_id", api.unregisterAgent)
	agents.GET("/:session_id/status", api.getAgentStatus)
	agents.POST("/:session_id/heartbeat", api.heartbeat)
	
	// Tool discovery and information
	agents.GET("/:session_id/tools", api.listTools)
	agents.GET("/:session_id/tools/:tool_name", api.getTool)
	
	// Tool execution
	agents.POST("/:session_id/tools/:tool_name/invoke", api.invokeTool)
	
	// Event subscription (WebSocket would be better, but HTTP for now)
	agents.GET("/:session_id/events", api.getEvents)
	
	// Admin endpoints
	admin := agents.Group("/admin")
	admin.GET("/sessions", api.listSessions)
	admin.GET("/metrics", api.getMetrics)
}

// RegisterAgent request/response structures
type RegisterAgentRequest struct {
	AgentID              string                `json:"agent_id" binding:"required"`
	AgentName            string                `json:"agent_name" binding:"required"`
	AgentVersion         string                `json:"agent_version"`
	Capabilities         *AgentCapabilities    `json:"capabilities"`
	Metadata             map[string]string     `json:"metadata"`
	SessionTimeoutSeconds int32                `json:"session_timeout_seconds"`
}

type AgentCapabilities struct {
	SupportedProtocols    []string `json:"supported_protocols"`
	SupportedToolTypes    []string `json:"supported_tool_types"`
	SupportsStreaming     bool     `json:"supports_streaming"`
	SupportsAsyncInvocation bool   `json:"supports_async_invocation"`
	MaxConcurrentTools    int32    `json:"max_concurrent_tools"`
	PreferredFormats      []string `json:"preferred_formats"`
}

type RegisterAgentResponse struct {
	SessionID      string        `json:"session_id"`
	ExpiresAt      int64         `json:"expires_at"`
	ServerInfo     *ServerInfo   `json:"server_info"`
	AvailableTools []ToolInfo    `json:"available_tools"`
}

type ServerInfo struct {
	ServerVersion     string            `json:"server_version"`
	ProtocolVersion   string            `json:"protocol_version"`
	SupportedFeatures []string          `json:"supported_features"`
	Capabilities      map[string]string `json:"capabilities"`
}

// Tool information structures
type ToolInfo struct {
	Name          string            `json:"name"`
	DisplayName   string            `json:"display_name"`
	Description   string            `json:"description"`
	Version       string            `json:"version"`
	Type          string            `json:"type"`
	Status        string            `json:"status"`
	Tags          []string          `json:"tags"`
	Metadata      map[string]string `json:"metadata"`
	CreatedAt     int64             `json:"created_at"`
	UpdatedAt     int64             `json:"updated_at"`
	Source        *ToolSource       `json:"source"`
}

type ToolSource struct {
	SpecID      string `json:"spec_id"`
	SpecType    string `json:"spec_type"`
	SpecPath    string `json:"spec_path"`
	OperationID string `json:"operation_id"`
	QueryName   string `json:"query_name"`
}

type GetToolResponse struct {
	Tool         ToolInfo       `json:"tool"`
	InputSchema  interface{}    `json:"input_schema,omitempty"`
	OutputSchema interface{}    `json:"output_schema,omitempty"`
	Examples     []ToolExample  `json:"examples,omitempty"`
}

type ToolExample struct {
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	Input          interface{} `json:"input"`
	ExpectedOutput interface{} `json:"expected_output"`
}

// Tool invocation structures
type InvokeToolRequest struct {
	Parameters map[string]interface{} `json:"parameters"`
	Options    *ToolInvocationOptions `json:"options"`
}

type ToolInvocationOptions struct {
	TimeoutSeconds int32             `json:"timeout_seconds"`
	Async          bool              `json:"async"`
	Context        map[string]string `json:"context"`
	RetryPolicy    *ToolRetryPolicy  `json:"retry_policy"`
}

type ToolRetryPolicy struct {
	MaxRetries            int32   `json:"max_retries"`
	RetryDelaySeconds     int32   `json:"retry_delay_seconds"`
	RetryableStatusCodes  []int32 `json:"retryable_status_codes"`
}

type InvokeToolResponse struct {
	InvocationID string        `json:"invocation_id"`
	Status       string        `json:"status"`
	Result       interface{}   `json:"result,omitempty"`
	Error        *ToolError    `json:"error,omitempty"`
	Metrics      *ToolMetrics  `json:"metrics"`
	ExecutedAt   int64         `json:"executed_at"`
}

type ToolError struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Details   string      `json:"details"`
	Metadata  interface{} `json:"metadata,omitempty"`
	Retryable bool        `json:"retryable"`
}

type ToolMetrics struct {
	ExecutionTimeMs int64              `json:"execution_time_ms"`
	MemoryUsedBytes int64              `json:"memory_used_bytes"`
	RetryCount      int32              `json:"retry_count"`
	CustomMetrics   map[string]float64 `json:"custom_metrics"`
}

// Session status structures
type AgentStatusResponse struct {
	SessionInfo     *AgentSessionInfo  `json:"session_info"`
	Metrics         *AgentMetrics      `json:"metrics"`
	RecentToolUsage []ToolUsageInfo    `json:"recent_tool_usage"`
}

type AgentSessionInfo struct {
	SessionID     string              `json:"session_id"`
	AgentID       string              `json:"agent_id"`
	AgentName     string              `json:"agent_name"`
	AgentVersion  string              `json:"agent_version"`
	CreatedAt     int64               `json:"created_at"`
	LastHeartbeat int64               `json:"last_heartbeat"`
	ExpiresAt     int64               `json:"expires_at"`
	Status        string              `json:"status"`
	Capabilities  *AgentCapabilities  `json:"capabilities"`
}

type AgentMetrics struct {
	TotalInvocations      int64             `json:"total_invocations"`
	SuccessfulInvocations int64             `json:"successful_invocations"`
	FailedInvocations     int64             `json:"failed_invocations"`
	AverageResponseTimeMs float64           `json:"average_response_time_ms"`
	LastInvocation        int64             `json:"last_invocation"`
	ToolUsageCount        map[string]int64  `json:"tool_usage_count"`
}

type ToolUsageInfo struct {
	ToolName        string `json:"tool_name"`
	InvokedAt       int64  `json:"invoked_at"`
	Status          string `json:"status"`
	ExecutionTimeMs int64  `json:"execution_time_ms"`
	ErrorMessage    string `json:"error_message,omitempty"`
}

type HeartbeatRequest struct {
	Status string `json:"status"`
}

type HeartbeatResponse struct {
	SessionValid         bool     `json:"session_valid"`
	NextHeartbeatAt      int64    `json:"next_heartbeat_at"`
	PendingNotifications []string `json:"pending_notifications"`
}

// Event structures
type Event struct {
	EventID   string      `json:"event_id"`
	Type      string      `json:"type"`
	Timestamp int64       `json:"timestamp"`
	SessionID string      `json:"session_id"`
	Data      interface{} `json:"data"`
}

type GetEventsResponse struct {
	Events []Event `json:"events"`
}

// Admin structures
type ListSessionsResponse struct {
	Sessions []AgentSessionInfo `json:"sessions"`
}

type MetricsResponse struct {
	TotalSessions     int                `json:"total_sessions"`
	ActiveSessions    int                `json:"active_sessions"`
	TotalInvocations  int64              `json:"total_invocations"`
	ToolUsageStats    map[string]int64   `json:"tool_usage_stats"`
	SessionMetrics    map[string]interface{} `json:"session_metrics"`
}

// registerAgent handles agent registration
func (api *AgentAPI) registerAgent(c *gin.Context) {
	var req RegisterAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to gRPC request
	grpcReq := &agentpb.RegisterAgentRequest{
		AgentId:              req.AgentID,
		AgentName:            req.AgentName,
		AgentVersion:         req.AgentVersion,
		SessionTimeoutSeconds: req.SessionTimeoutSeconds,
		Metadata:             req.Metadata,
	}

	if req.Capabilities != nil {
		grpcReq.Capabilities = &agentpb.AgentCapabilities{
			SupportedProtocols:      req.Capabilities.SupportedProtocols,
			SupportedToolTypes:      req.Capabilities.SupportedToolTypes,
			SupportsStreaming:       req.Capabilities.SupportsStreaming,
			SupportsAsyncInvocation: req.Capabilities.SupportsAsyncInvocation,
			MaxConcurrentTools:      req.Capabilities.MaxConcurrentTools,
			PreferredFormats:        req.Capabilities.PreferredFormats,
		}
	}

	// Call gRPC method
	grpcResp, err := api.agentServer.RegisterAgent(c.Request.Context(), grpcReq)
	if err != nil {
		api.logger.Error("Failed to register agent", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert response
	resp := RegisterAgentResponse{
		SessionID: grpcResp.SessionId,
		ExpiresAt: grpcResp.ExpiresAtUnix,
		ServerInfo: &ServerInfo{
			ServerVersion:     grpcResp.ServerInfo.ServerVersion,
			ProtocolVersion:   grpcResp.ServerInfo.ProtocolVersion,
			SupportedFeatures: grpcResp.ServerInfo.SupportedFeatures,
			Capabilities:      grpcResp.ServerInfo.Capabilities,
		},
		AvailableTools: make([]ToolInfo, len(grpcResp.AvailableTools)),
	}

	for i, tool := range grpcResp.AvailableTools {
		resp.AvailableTools[i] = api.convertToolInfo(tool)
	}

	api.logger.Info("Agent registered via REST API",
		zap.String("agent_id", req.AgentID),
		zap.String("session_id", resp.SessionID))

	c.JSON(http.StatusCreated, resp)
}

// unregisterAgent handles agent unregistration
func (api *AgentAPI) unregisterAgent(c *gin.Context) {
	sessionID := c.Param("session_id")

	grpcReq := &agentpb.UnregisterAgentRequest{
		SessionId: sessionID,
	}

	grpcResp, err := api.agentServer.UnregisterAgent(c.Request.Context(), grpcReq)
	if err != nil {
		api.logger.Error("Failed to unregister agent", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	api.logger.Info("Agent unregistered via REST API", zap.String("session_id", sessionID))

	c.JSON(http.StatusOK, gin.H{
		"success": grpcResp.Success,
		"message": grpcResp.Message,
	})
}

// listTools handles tool listing for an agent
func (api *AgentAPI) listTools(c *gin.Context) {
	sessionID := c.Param("session_id")

	// Parse query parameters for filtering and pagination
	grpcReq := &agentpb.ListToolsRequest{
		SessionId: sessionID,
	}

	// Add basic pagination if requested
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			grpcReq.Pagination = &agentpb.PaginationOptions{
				Page: int32(page),
			}
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
			if grpcReq.Pagination == nil {
				grpcReq.Pagination = &agentpb.PaginationOptions{}
			}
			grpcReq.Pagination.PageSize = int32(pageSize)
		}
	}

	grpcResp, err := api.agentServer.ListTools(c.Request.Context(), grpcReq)
	if err != nil {
		api.logger.Error("Failed to list tools", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tools := make([]ToolInfo, len(grpcResp.Tools))
	for i, tool := range grpcResp.Tools {
		tools[i] = api.convertToolInfo(tool)
	}

	c.JSON(http.StatusOK, gin.H{
		"tools":       tools,
		"total_count": grpcResp.TotalCount,
		"pagination":  grpcResp.Pagination,
	})
}

// getTool handles getting detailed tool information
func (api *AgentAPI) getTool(c *gin.Context) {
	sessionID := c.Param("session_id")
	toolName := c.Param("tool_name")
	includeSchema := c.Query("include_schema") == "true"

	grpcReq := &agentpb.GetToolRequest{
		SessionId:     sessionID,
		ToolName:      toolName,
		IncludeSchema: includeSchema,
	}

	grpcResp, err := api.agentServer.GetTool(c.Request.Context(), grpcReq)
	if err != nil {
		api.logger.Error("Failed to get tool", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := GetToolResponse{
		Tool: api.convertToolInfo(grpcResp.Tool),
	}

	if includeSchema {
		// Parse JSON schema strings back to objects
		// Parse input schema if available
		if grpcResp.InputSchemaJson != "" {
			var inputSchema map[string]interface{}
			if err := json.Unmarshal([]byte(grpcResp.InputSchemaJson), &inputSchema); err != nil {
				api.logger.Warn("Failed to parse input schema JSON", 
					zap.String("schema", grpcResp.InputSchemaJson), 
					zap.Error(err))
				// Fallback to placeholder
				resp.InputSchema = map[string]interface{}{"type": "object"}
			} else {
				resp.InputSchema = inputSchema
			}
		} else {
			resp.InputSchema = map[string]interface{}{"type": "object"}
		}
		
		// Parse output schema if available
		if grpcResp.OutputSchemaJson != "" {
			var outputSchema map[string]interface{}
			if err := json.Unmarshal([]byte(grpcResp.OutputSchemaJson), &outputSchema); err != nil {
				api.logger.Warn("Failed to parse output schema JSON", 
					zap.String("schema", grpcResp.OutputSchemaJson), 
					zap.Error(err))
				// Fallback to placeholder
				resp.OutputSchema = map[string]interface{}{"type": "object"}
			} else {
				resp.OutputSchema = outputSchema
			}
		} else {
			resp.OutputSchema = map[string]interface{}{"type": "object"}
		}
		
		resp.Examples = make([]ToolExample, len(grpcResp.Examples))
		for i, example := range grpcResp.Examples {
			var inputMap map[string]interface{}
			var outputMap map[string]interface{}

			// Parse example input JSON
			if example.InputJson != "" {
				if err := json.Unmarshal([]byte(example.InputJson), &inputMap); err != nil {
					api.logger.Warn("Failed to parse example input JSON", 
						zap.String("input", example.InputJson), 
						zap.Error(err))
					inputMap = map[string]interface{}{}
				}
			} else {
				inputMap = map[string]interface{}{}
			}
			
			// Parse example expected output JSON
			if example.ExpectedOutputJson != "" {
				if err := json.Unmarshal([]byte(example.ExpectedOutputJson), &outputMap); err != nil {
					api.logger.Warn("Failed to parse example expected output JSON", 
						zap.String("expected_output", example.ExpectedOutputJson), 
						zap.Error(err))
					outputMap = map[string]interface{}{}
				}
			} else {
				outputMap = map[string]interface{}{}
			}

			resp.Examples[i] = ToolExample{
				Name:           example.Name,
				Description:    example.Description,
				Input:          inputMap,
				ExpectedOutput: outputMap,
			}
		}
	}

	c.JSON(http.StatusOK, resp)
}

// invokeTool handles tool execution
func (api *AgentAPI) invokeTool(c *gin.Context) {
	sessionID := c.Param("session_id")
	toolName := c.Param("tool_name")

	var req InvokeToolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invocationID := uuid.New().String()

	// Serialize parameters to JSON
	parametersJSON := "{}"
	if req.Parameters != nil {
		paramsBytes, err := json.Marshal(req.Parameters)
		if err != nil {
			api.logger.Error("Failed to marshal parameters", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameters format"})
			return
		}
		parametersJSON = string(paramsBytes)
	}

	// Convert to gRPC request
	grpcReq := &agentpb.InvokeToolRequest{
		SessionId:      sessionID,
		ToolName:       toolName,
		InvocationId:   invocationID,
		ParametersJson: parametersJSON,
	}

	if req.Options != nil {
		grpcReq.Options = &agentpb.ToolInvocationOptions{
			TimeoutSeconds: req.Options.TimeoutSeconds,
			Async:          req.Options.Async,
			Context:        req.Options.Context,
		}

		if req.Options.RetryPolicy != nil {
			grpcReq.Options.RetryPolicy = &agentpb.ToolRetryPolicy{
				MaxRetries:           req.Options.RetryPolicy.MaxRetries,
				RetryDelaySeconds:    req.Options.RetryPolicy.RetryDelaySeconds,
				RetryableStatusCodes: req.Options.RetryPolicy.RetryableStatusCodes,
			}
		}
	}

	grpcResp, err := api.agentServer.InvokeTool(c.Request.Context(), grpcReq)
	if err != nil {
		api.logger.Error("Failed to invoke tool", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := InvokeToolResponse{
		InvocationID: grpcResp.InvocationId,
		Status:       grpcResp.Status.String(),
		ExecutedAt:   grpcResp.ExecutedAtUnix,
	}

	// Parse result from JSON
	if grpcResp.ResultJson != "" {
		var result interface{}
		if err := json.Unmarshal([]byte(grpcResp.ResultJson), &result); err != nil {
			api.logger.Error("Failed to parse tool result JSON", 
				zap.Error(err), 
				zap.String("result_json", grpcResp.ResultJson))
			resp.Result = map[string]interface{}{"_error": "Failed to parse result JSON"}
		} else {
			resp.Result = result
		}
	}

	if grpcResp.Error != nil {
		resp.Error = &ToolError{
			Code:      grpcResp.Error.Code.String(),
			Message:   grpcResp.Error.Message,
			Details:   grpcResp.Error.Details,
			Retryable: grpcResp.Error.Retryable,
		}
	}

	if grpcResp.Metrics != nil {
		resp.Metrics = &ToolMetrics{
			ExecutionTimeMs: grpcResp.Metrics.ExecutionTimeMs,
			MemoryUsedBytes: grpcResp.Metrics.MemoryUsedBytes,
			RetryCount:      grpcResp.Metrics.RetryCount,
			CustomMetrics:   grpcResp.Metrics.CustomMetrics,
		}
	}

	statusCode := http.StatusOK
	if grpcResp.Status == agentpb.ToolInvocationStatus_TOOL_INVOCATION_STATUS_FAILED {
		statusCode = http.StatusInternalServerError
	}

	api.logger.Info("Tool invoked via REST API",
		zap.String("session_id", sessionID),
		zap.String("tool_name", toolName),
		zap.String("invocation_id", invocationID),
		zap.String("status", resp.Status))

	c.JSON(statusCode, resp)
}

// getAgentStatus handles getting agent session status
func (api *AgentAPI) getAgentStatus(c *gin.Context) {
	sessionID := c.Param("session_id")

	grpcReq := &agentpb.GetAgentStatusRequest{
		SessionId: sessionID,
	}

	grpcResp, err := api.agentServer.GetAgentStatus(c.Request.Context(), grpcReq)
	if err != nil {
		api.logger.Error("Failed to get agent status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := AgentStatusResponse{
		SessionInfo: &AgentSessionInfo{
			SessionID:     grpcResp.SessionInfo.SessionId,
			AgentID:       grpcResp.SessionInfo.AgentId,
			AgentName:     grpcResp.SessionInfo.AgentName,
			AgentVersion:  grpcResp.SessionInfo.AgentVersion,
			CreatedAt:     grpcResp.SessionInfo.CreatedAtUnix,
			LastHeartbeat: grpcResp.SessionInfo.LastHeartbeatUnix,
			ExpiresAt:     grpcResp.SessionInfo.ExpiresAtUnix,
			Status:        grpcResp.SessionInfo.Status.String(),
		},
		Metrics: &AgentMetrics{
			TotalInvocations:      grpcResp.Metrics.TotalInvocations,
			SuccessfulInvocations: grpcResp.Metrics.SuccessfulInvocations,
			FailedInvocations:     grpcResp.Metrics.FailedInvocations,
			AverageResponseTimeMs: grpcResp.Metrics.AverageResponseTimeMs,
			LastInvocation:        grpcResp.Metrics.LastInvocationUnix,
			ToolUsageCount:        grpcResp.Metrics.ToolUsageCount,
		},
		RecentToolUsage: make([]ToolUsageInfo, len(grpcResp.RecentToolUsage)),
	}

	if grpcResp.SessionInfo.Capabilities != nil {
		resp.SessionInfo.Capabilities = &AgentCapabilities{
			SupportedProtocols:      grpcResp.SessionInfo.Capabilities.SupportedProtocols,
			SupportedToolTypes:      grpcResp.SessionInfo.Capabilities.SupportedToolTypes,
			SupportsStreaming:       grpcResp.SessionInfo.Capabilities.SupportsStreaming,
			SupportsAsyncInvocation: grpcResp.SessionInfo.Capabilities.SupportsAsyncInvocation,
			MaxConcurrentTools:      grpcResp.SessionInfo.Capabilities.MaxConcurrentTools,
			PreferredFormats:        grpcResp.SessionInfo.Capabilities.PreferredFormats,
		}
	}

	for i, usage := range grpcResp.RecentToolUsage {
		resp.RecentToolUsage[i] = ToolUsageInfo{
			ToolName:        usage.ToolName,
			InvokedAt:       usage.InvokedAtUnix,
			Status:          usage.Status.String(),
			ExecutionTimeMs: usage.ExecutionTimeMs,
			ErrorMessage:    usage.ErrorMessage,
		}
	}

	c.JSON(http.StatusOK, resp)
}

// heartbeat handles agent heartbeat
func (api *AgentAPI) heartbeat(c *gin.Context) {
	sessionID := c.Param("session_id")

	var req HeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := agentpb.AgentStatus_AGENT_STATUS_ACTIVE
	if req.Status != "" {
		// Would parse status string to enum
		switch req.Status {
		case "idle":
			status = agentpb.AgentStatus_AGENT_STATUS_IDLE
		case "busy":
			status = agentpb.AgentStatus_AGENT_STATUS_BUSY
		case "error":
			status = agentpb.AgentStatus_AGENT_STATUS_ERROR
		}
	}

	grpcReq := &agentpb.HeartBeatRequest{
		SessionId: sessionID,
		Status:    status,
	}

	grpcResp, err := api.agentServer.HeartBeat(c.Request.Context(), grpcReq)
	if err != nil {
		api.logger.Error("Failed to process heartbeat", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := HeartbeatResponse{
		SessionValid:         grpcResp.SessionValid,
		NextHeartbeatAt:      grpcResp.NextHeartbeatAtUnix,
		PendingNotifications: grpcResp.PendingNotifications,
	}

	c.JSON(http.StatusOK, resp)
}

// getEvents handles getting recent events (placeholder for real-time events)
func (api *AgentAPI) getEvents(c *gin.Context) {
	sessionID := c.Param("session_id")

	// Validate session exists
	if _, exists := api.agentServer.getSession(sessionID); !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	// Placeholder - in a real implementation, this would return recent events
	// For now, return empty events list
	resp := GetEventsResponse{
		Events: []Event{},
	}

	c.JSON(http.StatusOK, resp)
}

// listSessions handles listing all active sessions (admin)
func (api *AgentAPI) listSessions(c *gin.Context) {
	api.agentServer.sessionsMux.RLock()
	sessions := make([]AgentSessionInfo, 0, len(api.agentServer.sessions))
	
	for _, session := range api.agentServer.sessions {
		sessionInfo := AgentSessionInfo{
			SessionID:     session.ID,
			AgentID:       session.AgentID,
			AgentName:     session.AgentName,
			AgentVersion:  session.AgentVersion,
			CreatedAt:     session.CreatedAt.Unix(),
			LastHeartbeat: session.LastHeartbeat.Unix(),
			ExpiresAt:     session.ExpiresAt.Unix(),
			Status:        session.Status.String(),
		}

		if session.Capabilities != nil {
			sessionInfo.Capabilities = &AgentCapabilities{
				SupportedProtocols:      session.Capabilities.SupportedProtocols,
				SupportedToolTypes:      session.Capabilities.SupportedToolTypes,
				SupportsStreaming:       session.Capabilities.SupportsStreaming,
				SupportsAsyncInvocation: session.Capabilities.SupportsAsyncInvocation,
				MaxConcurrentTools:      session.Capabilities.MaxConcurrentTools,
				PreferredFormats:        session.Capabilities.PreferredFormats,
			}
		}

		sessions = append(sessions, sessionInfo)
	}
	api.agentServer.sessionsMux.RUnlock()

	resp := ListSessionsResponse{
		Sessions: sessions,
	}

	c.JSON(http.StatusOK, resp)
}

// getMetrics handles getting server metrics (admin)
func (api *AgentAPI) getMetrics(c *gin.Context) {
	api.agentServer.sessionsMux.RLock()
	totalSessions := len(api.agentServer.sessions)
	activeSessions := 0
	
	var totalInvocations int64
	toolUsageStats := make(map[string]int64)

	for _, session := range api.agentServer.sessions {
		if session.Status == agentpb.AgentStatus_AGENT_STATUS_ACTIVE {
			activeSessions++
		}

		session.Metrics.mu.RLock()
		totalInvocations += session.Metrics.TotalInvocations
		for tool, count := range session.Metrics.ToolUsageCount {
			toolUsageStats[tool] += count
		}
		session.Metrics.mu.RUnlock()
	}
	api.agentServer.sessionsMux.RUnlock()

	resp := MetricsResponse{
		TotalSessions:    totalSessions,
		ActiveSessions:   activeSessions,
		TotalInvocations: totalInvocations,
		ToolUsageStats:   toolUsageStats,
		SessionMetrics:   map[string]interface{}{},
	}

	c.JSON(http.StatusOK, resp)
}

// Helper methods

func (api *AgentAPI) convertToolInfo(grpcTool *agentpb.ToolInfo) ToolInfo {
	tool := ToolInfo{
		Name:        grpcTool.Name,
		DisplayName: grpcTool.DisplayName,
		Description: grpcTool.Description,
		Version:     grpcTool.Version,
		Type:        grpcTool.Type.String(),
		Status:      grpcTool.Status.String(),
		Tags:        grpcTool.Tags,
		Metadata:    grpcTool.Metadata,
		CreatedAt:   grpcTool.CreatedAtUnix,
		UpdatedAt:   grpcTool.UpdatedAtUnix,
	}

	if grpcTool.Source != nil {
		tool.Source = &ToolSource{
			SpecID:      grpcTool.Source.SpecId,
			SpecType:    grpcTool.Source.SpecType,
			SpecPath:    grpcTool.Source.SpecPath,
			OperationID: grpcTool.Source.OperationId,
			QueryName:   grpcTool.Source.QueryName,
		}
	}

	return tool
}