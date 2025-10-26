package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	agentpb "github.com/aionmcp/aionmcp/pkg/agent/proto"
	"github.com/aionmcp/aionmcp/pkg/types"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AgentServer implements the gRPC AgentService interface
type AgentServer struct {
	agentpb.UnimplementedAgentServiceServer
	logger       *zap.Logger
	registry     types.ToolRegistry
	sessions     map[string]*AgentSession
	sessionsMux  sync.RWMutex
	eventStreams map[string][]chan *agentpb.Event
	streamsMux   sync.RWMutex
}

// AgentSession represents an active agent session
type AgentSession struct {
	ID            string
	AgentID       string
	AgentName     string
	AgentVersion  string
	Capabilities  *agentpb.AgentCapabilities
	Metadata      map[string]string
	CreatedAt     time.Time
	LastHeartbeat time.Time
	ExpiresAt     time.Time
	Status        agentpb.AgentStatus
	Metrics       *InternalAgentMetrics
}

// InternalAgentMetrics tracks agent usage statistics
type InternalAgentMetrics struct {
	TotalInvocations      int64
	SuccessfulInvocations int64
	FailedInvocations     int64
	TotalResponseTimeMs   int64
	LastInvocation        time.Time
	ToolUsageCount        map[string]int64
	mu                    sync.RWMutex
}

// NewAgentServer creates a new AgentServer instance
func NewAgentServer(logger *zap.Logger, registry types.ToolRegistry) *AgentServer {
	server := &AgentServer{
		logger:       logger,
		registry:     registry,
		sessions:     make(map[string]*AgentSession),
		eventStreams: make(map[string][]chan *agentpb.Event),
	}

	// Start session cleanup goroutine
	go server.sessionCleanup()

	return server
}

// RegisterAgent establishes a new agent session
func (s *AgentServer) RegisterAgent(ctx context.Context, req *agentpb.RegisterAgentRequest) (*agentpb.RegisterAgentResponse, error) {
	s.logger.Info("Agent registration request",
		zap.String("agent_id", req.AgentId),
		zap.String("agent_name", req.AgentName),
		zap.String("agent_version", req.AgentVersion))

	// Validate request
	if req.AgentId == "" {
		return nil, status.Error(codes.InvalidArgument, "agent_id is required")
	}
	if req.AgentName == "" {
		return nil, status.Error(codes.InvalidArgument, "agent_name is required")
	}

	// Generate session ID
	sessionID := uuid.New().String()

	// Set session timeout (default 300 seconds)
	timeoutSeconds := req.SessionTimeoutSeconds
	if timeoutSeconds == 0 {
		timeoutSeconds = 300
	}

	now := time.Now()
	expiresAt := now.Add(time.Duration(timeoutSeconds) * time.Second)

	// Create session
	session := &AgentSession{
		ID:            sessionID,
		AgentID:       req.AgentId,
		AgentName:     req.AgentName,
		AgentVersion:  req.AgentVersion,
		Capabilities:  req.Capabilities,
		Metadata:      req.Metadata,
		CreatedAt:     now,
		LastHeartbeat: now,
		ExpiresAt:     expiresAt,
		Status:        agentpb.AgentStatus_AGENT_STATUS_ACTIVE,
		Metrics: &InternalAgentMetrics{
			ToolUsageCount: make(map[string]int64),
		},
	}

	// Store session
	s.sessionsMux.Lock()
	s.sessions[sessionID] = session
	s.sessionsMux.Unlock()

	// Get available tools
	tools := s.getToolsForAgent(session)

	// Broadcast agent registered event
	s.broadcastEvent(&agentpb.Event{
		EventId:       uuid.New().String(),
		Type:          agentpb.EventType_EVENT_TYPE_AGENT_REGISTERED,
		TimestampUnix: now.Unix(),
		SessionId:     sessionID,
		DataJson:      fmt.Sprintf(`{"agent_id": "%s", "agent_name": "%s"}`, req.AgentId, req.AgentName),
	})

	s.logger.Info("Agent registered successfully",
		zap.String("session_id", sessionID),
		zap.String("agent_id", req.AgentId),
		zap.Int("available_tools", len(tools)))

	return &agentpb.RegisterAgentResponse{
		SessionId:     sessionID,
		ExpiresAtUnix: expiresAt.Unix(),
		ServerInfo: &agentpb.ServerInfo{
			ServerVersion:    "0.1.0",
			ProtocolVersion:  "MCP/1.0",
			SupportedFeatures: []string{"tool_execution", "event_streaming", "session_management"},
			Capabilities: map[string]string{
				"max_concurrent_tools": "10",
				"streaming_supported":  "true",
				"async_execution":      "true",
			},
		},
		AvailableTools: tools,
	}, nil
}

// UnregisterAgent terminates an agent session
func (s *AgentServer) UnregisterAgent(ctx context.Context, req *agentpb.UnregisterAgentRequest) (*agentpb.UnregisterAgentResponse, error) {
	s.logger.Info("Agent unregistration request", zap.String("session_id", req.SessionId))

	session, exists := s.getSession(req.SessionId)
	if !exists {
		return nil, status.Error(codes.NotFound, "session not found")
	}

	// Remove session
	s.sessionsMux.Lock()
	delete(s.sessions, req.SessionId)
	s.sessionsMux.Unlock()

	// Close event streams for this session
	s.closeEventStreams(req.SessionId)

	// Broadcast agent unregistered event
	s.broadcastEvent(&agentpb.Event{
		EventId:       uuid.New().String(),
		Type:          agentpb.EventType_EVENT_TYPE_AGENT_UNREGISTERED,
		TimestampUnix: time.Now().Unix(),
		SessionId:     req.SessionId,
		DataJson:      fmt.Sprintf(`{"agent_id": "%s"}`, session.AgentID),
	})

	s.logger.Info("Agent unregistered successfully",
		zap.String("session_id", req.SessionId),
		zap.String("agent_id", session.AgentID))

	return &agentpb.UnregisterAgentResponse{
		Success: true,
		Message: "Agent session terminated successfully",
	}, nil
}

// ListTools returns available tools for the agent
func (s *AgentServer) ListTools(ctx context.Context, req *agentpb.ListToolsRequest) (*agentpb.ListToolsResponse, error) {
	session, exists := s.getSession(req.SessionId)
	if !exists {
		return nil, status.Error(codes.Unauthenticated, "invalid session")
	}

	// Update last heartbeat
	s.updateHeartbeat(req.SessionId)

	tools := s.getToolsForAgent(session)

	// Apply filtering if specified
	if req.Filter != nil {
		tools = s.applyToolFilter(tools, req.Filter)
	}

	// Apply pagination
	totalCount := len(tools)
	if req.Pagination != nil {
		tools = s.applyPagination(tools, req.Pagination)
	}

	s.logger.Debug("Listed tools for agent",
		zap.String("session_id", req.SessionId),
		zap.Int("total_tools", totalCount),
		zap.Int("returned_tools", len(tools)))

	return &agentpb.ListToolsResponse{
		Tools:      tools,
		TotalCount: int32(totalCount),
		Pagination: &agentpb.PaginationMetadata{
			CurrentPage:  1,
			PageSize:     int32(len(tools)),
			TotalPages:   1,
			HasNext:      false,
			HasPrevious:  false,
		},
	}, nil
}

// GetTool returns detailed information about a specific tool
func (s *AgentServer) GetTool(ctx context.Context, req *agentpb.GetToolRequest) (*agentpb.GetToolResponse, error) {
	_, exists := s.getSession(req.SessionId)
	if !exists {
		return nil, status.Error(codes.Unauthenticated, "invalid session")
	}

	// Update last heartbeat
	s.updateHeartbeat(req.SessionId)

	tool, err := s.registry.Get(req.ToolName)
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("tool not found: %s", req.ToolName))
	}

	toolInfo := s.convertToToolInfo(tool)

	var inputSchema, outputSchema string
	var examples []*agentpb.ToolExample

	if req.IncludeSchema {
		// Get schema information if available
		inputSchema = `{"type": "object", "properties": {}}`  // Placeholder
		outputSchema = `{"type": "object", "properties": {}}` // Placeholder

		// Add example usage
		examples = []*agentpb.ToolExample{
			{
				Name:                "Basic Usage",
				Description:         fmt.Sprintf("Example usage of %s tool", req.ToolName),
				InputJson:           `{"parameter": "example_value"}`,
				ExpectedOutputJson:  `{"result": "example_result"}`,
			},
		}
	}

	s.logger.Debug("Retrieved tool details",
		zap.String("session_id", req.SessionId),
		zap.String("tool_name", req.ToolName),
		zap.Bool("include_schema", req.IncludeSchema))

	return &agentpb.GetToolResponse{
		Tool:             toolInfo,
		InputSchemaJson:  inputSchema,
		OutputSchemaJson: outputSchema,
		Examples:         examples,
	}, nil
}

// InvokeTool executes a tool with given parameters
func (s *AgentServer) InvokeTool(ctx context.Context, req *agentpb.InvokeToolRequest) (*agentpb.InvokeToolResponse, error) {
	session, exists := s.getSession(req.SessionId)
	if !exists {
		return nil, status.Error(codes.Unauthenticated, "invalid session")
	}

	// Update last heartbeat
	s.updateHeartbeat(req.SessionId)

	startTime := time.Now()
	
	s.logger.Info("Tool invocation request",
		zap.String("session_id", req.SessionId),
		zap.String("tool_name", req.ToolName),
		zap.String("invocation_id", req.InvocationId))

	// Get tool from registry
	tool, err := s.registry.Get(req.ToolName)
	if err != nil {
		s.updateMetrics(session, req.ToolName, false, time.Since(startTime))
		return nil, status.Error(codes.NotFound, fmt.Sprintf("tool not found: %s", req.ToolName))
	}

	// Parse parameters from JSON
	// For now, we'll use a simple map[string]interface{} approach
	// In a full implementation, you'd want proper JSON unmarshaling
	var parameters map[string]interface{}
	if req.ParametersJson != "" {
		// Placeholder: would use json.Unmarshal here
		parameters = make(map[string]interface{})
	}

	// Execute tool
	result, err := tool.Execute(parameters)
	executionTime := time.Since(startTime)

	var toolError *agentpb.ToolError
	var resultJson string
	var status agentpb.ToolInvocationStatus

	if err != nil {
		status = agentpb.ToolInvocationStatus_TOOL_INVOCATION_STATUS_FAILED
		toolError = &agentpb.ToolError{
			Code:      agentpb.ErrorCode_ERROR_CODE_EXECUTION_FAILED,
			Message:   err.Error(),
			Details:   fmt.Sprintf("Tool execution failed: %v", err),
			Retryable: true,
		}
		s.updateMetrics(session, req.ToolName, false, executionTime)
		
		s.logger.Error("Tool execution failed",
			zap.String("session_id", req.SessionId),
			zap.String("tool_name", req.ToolName),
			zap.String("invocation_id", req.InvocationId),
			zap.Error(err))
	} else {
		status = agentpb.ToolInvocationStatus_TOOL_INVOCATION_STATUS_SUCCESS
		resultJson = fmt.Sprintf(`{"result": %v}`, result) // Placeholder JSON serialization
		s.updateMetrics(session, req.ToolName, true, executionTime)
		
		s.logger.Info("Tool executed successfully",
			zap.String("session_id", req.SessionId),
			zap.String("tool_name", req.ToolName),
			zap.String("invocation_id", req.InvocationId),
			zap.Duration("execution_time", executionTime))
	}

	// Broadcast tool invocation event
	s.broadcastEvent(&agentpb.Event{
		EventId:       uuid.New().String(),
		Type:          agentpb.EventType_EVENT_TYPE_TOOL_INVOCATION,
		TimestampUnix: time.Now().Unix(),
		SessionId:     req.SessionId,
		DataJson:      fmt.Sprintf(`{"tool_name": "%s", "status": "%s", "execution_time_ms": %d}`, req.ToolName, status.String(), executionTime.Milliseconds()),
	})

	return &agentpb.InvokeToolResponse{
		InvocationId:    req.InvocationId,
		Status:          status,
		ResultJson:      resultJson,
		Error:           toolError,
		Metrics: &agentpb.ToolMetrics{
			ExecutionTimeMs: executionTime.Milliseconds(),
			RetryCount:      0,
			CustomMetrics: map[string]float64{
				"execution_timestamp": float64(time.Now().Unix()),
			},
		},
		ExecutedAtUnix: time.Now().Unix(),
	}, nil
}

// StreamEvents provides real-time events to agents
func (s *AgentServer) StreamEvents(req *agentpb.StreamEventsRequest, stream agentpb.AgentService_StreamEventsServer) error {
	session, exists := s.getSession(req.SessionId)
	if !exists {
		return status.Error(codes.Unauthenticated, "invalid session")
	}

	s.logger.Info("Starting event stream",
		zap.String("session_id", req.SessionId),
		zap.String("agent_id", session.AgentID))

	// Create event channel for this stream
	eventChan := make(chan *agentpb.Event, 100)

	// Register the stream
	s.streamsMux.Lock()
	if s.eventStreams[req.SessionId] == nil {
		s.eventStreams[req.SessionId] = make([]chan *agentpb.Event, 0)
	}
	s.eventStreams[req.SessionId] = append(s.eventStreams[req.SessionId], eventChan)
	s.streamsMux.Unlock()

	// Send initial connection event
	connectEvent := &agentpb.Event{
		EventId:       uuid.New().String(),
		Type:          agentpb.EventType_EVENT_TYPE_SERVER_STATUS,
		TimestampUnix: time.Now().Unix(),
		SessionId:     req.SessionId,
		DataJson:      `{"status": "connected", "message": "Event stream established"}`,
	}

	if err := stream.Send(connectEvent); err != nil {
		s.logger.Error("Failed to send connection event", zap.Error(err))
		return err
	}

	// Stream events until context is done or client disconnects
	for {
		select {
		case <-stream.Context().Done():
			s.logger.Info("Event stream closed by client",
				zap.String("session_id", req.SessionId))
			s.removeEventStream(req.SessionId, eventChan)
			return nil

		case event := <-eventChan:
			if err := stream.Send(event); err != nil {
				s.logger.Error("Failed to send event",
					zap.String("session_id", req.SessionId),
					zap.Error(err))
				s.removeEventStream(req.SessionId, eventChan)
				return err
			}
		}
	}
}

// HeartBeat maintains agent session liveness
func (s *AgentServer) HeartBeat(ctx context.Context, req *agentpb.HeartBeatRequest) (*agentpb.HeartBeatResponse, error) {
	session, exists := s.getSession(req.SessionId)
	if !exists {
		return &agentpb.HeartBeatResponse{
			SessionValid: false,
		}, nil
	}

	// Update heartbeat and status
	s.sessionsMux.Lock()
	session.LastHeartbeat = time.Now()
	if req.Status != agentpb.AgentStatus_AGENT_STATUS_UNSPECIFIED {
		session.Status = req.Status
	}
	s.sessionsMux.Unlock()

	nextHeartbeat := time.Now().Add(30 * time.Second) // 30 second heartbeat interval

	return &agentpb.HeartBeatResponse{
		SessionValid:           true,
		NextHeartbeatAtUnix:    nextHeartbeat.Unix(),
		PendingNotifications:   []string{}, // Placeholder for future notifications
	}, nil
}

// GetAgentStatus returns current agent session information
func (s *AgentServer) GetAgentStatus(ctx context.Context, req *agentpb.GetAgentStatusRequest) (*agentpb.GetAgentStatusResponse, error) {
	session, exists := s.getSession(req.SessionId)
	if !exists {
		return nil, status.Error(codes.NotFound, "session not found")
	}

	s.updateHeartbeat(req.SessionId)

	sessionInfo := &agentpb.AgentSessionInfo{
		SessionId:          session.ID,
		AgentId:            session.AgentID,
		AgentName:          session.AgentName,
		AgentVersion:       session.AgentVersion,
		CreatedAtUnix:      session.CreatedAt.Unix(),
		LastHeartbeatUnix:  session.LastHeartbeat.Unix(),
		ExpiresAtUnix:      session.ExpiresAt.Unix(),
		Status:             session.Status,
		Capabilities:       session.Capabilities,
	}

	session.Metrics.mu.RLock()
	metrics := &agentpb.AgentMetrics{
		TotalInvocations:      session.Metrics.TotalInvocations,
		SuccessfulInvocations: session.Metrics.SuccessfulInvocations,
		FailedInvocations:     session.Metrics.FailedInvocations,
		ToolUsageCount:        session.Metrics.ToolUsageCount,
		LastInvocationUnix:    session.Metrics.LastInvocation.Unix(),
	}
	
	if session.Metrics.TotalInvocations > 0 {
		metrics.AverageResponseTimeMs = float64(session.Metrics.TotalResponseTimeMs) / float64(session.Metrics.TotalInvocations)
	}
	session.Metrics.mu.RUnlock()

	return &agentpb.GetAgentStatusResponse{
		SessionInfo: sessionInfo,
		Metrics:     metrics,
		RecentToolUsage: []*agentpb.ToolUsageInfo{}, // Placeholder for recent usage history
	}, nil
}

// Helper methods

func (s *AgentServer) getSession(sessionID string) (*AgentSession, bool) {
	s.sessionsMux.RLock()
	defer s.sessionsMux.RUnlock()
	session, exists := s.sessions[sessionID]
	return session, exists
}

func (s *AgentServer) updateHeartbeat(sessionID string) {
	s.sessionsMux.Lock()
	defer s.sessionsMux.Unlock()
	if session, exists := s.sessions[sessionID]; exists {
		session.LastHeartbeat = time.Now()
	}
}

func (s *AgentServer) updateMetrics(session *AgentSession, toolName string, success bool, duration time.Duration) {
	session.Metrics.mu.Lock()
	defer session.Metrics.mu.Unlock()

	session.Metrics.TotalInvocations++
	session.Metrics.TotalResponseTimeMs += duration.Milliseconds()
	session.Metrics.LastInvocation = time.Now()
	
	if success {
		session.Metrics.SuccessfulInvocations++
	} else {
		session.Metrics.FailedInvocations++
	}

	if session.Metrics.ToolUsageCount[toolName] == 0 {
		session.Metrics.ToolUsageCount[toolName] = 0
	}
	session.Metrics.ToolUsageCount[toolName]++
}

func (s *AgentServer) getToolsForAgent(session *AgentSession) []*agentpb.ToolInfo {
	toolMetadata := s.registry.ListTools()
	result := make([]*agentpb.ToolInfo, 0, len(toolMetadata))

	for _, metadata := range toolMetadata {
		result = append(result, s.convertToolMetadataToToolInfo(metadata))
	}

	return result
}

func (s *AgentServer) convertToToolInfo(tool types.Tool) *agentpb.ToolInfo {
	metadata := tool.Metadata()
	return s.convertToolMetadataToToolInfo(metadata)
}

func (s *AgentServer) convertToolMetadataToToolInfo(metadata types.ToolMetadata) *agentpb.ToolInfo {
	return &agentpb.ToolInfo{
		Name:            metadata.Name,
		DisplayName:     metadata.Name,
		Description:     metadata.Description,
		Version:         metadata.Version,
		Type:           agentpb.ToolType_TOOL_TYPE_FUNCTION, // Default type
		Status:         agentpb.ToolStatus_TOOL_STATUS_AVAILABLE,
		Tags:           metadata.Tags,
		Metadata:       make(map[string]string),
		CreatedAtUnix:  metadata.CreatedAt.Unix(),
		UpdatedAtUnix:  metadata.UpdatedAt.Unix(),
		Source: &agentpb.ToolSource{
			SpecId:   metadata.Source,
			SpecType: metadata.Source,
		},
	}
}

func (s *AgentServer) applyToolFilter(tools []*agentpb.ToolInfo, filter *agentpb.ToolFilter) []*agentpb.ToolInfo {
	// Placeholder implementation - would include actual filtering logic
	return tools
}

func (s *AgentServer) applyPagination(tools []*agentpb.ToolInfo, pagination *agentpb.PaginationOptions) []*agentpb.ToolInfo {
	// Placeholder implementation - would include actual pagination logic
	return tools
}

func (s *AgentServer) broadcastEvent(event *agentpb.Event) {
	s.streamsMux.RLock()
	defer s.streamsMux.RUnlock()

	for sessionID, streams := range s.eventStreams {
		for _, stream := range streams {
			select {
			case stream <- event:
				// Event sent successfully
			default:
				// Channel is full, skip this stream
				s.logger.Warn("Event stream channel full",
					zap.String("session_id", sessionID),
					zap.String("event_type", event.Type.String()))
			}
		}
	}
}

func (s *AgentServer) removeEventStream(sessionID string, targetChan chan *agentpb.Event) {
	s.streamsMux.Lock()
	defer s.streamsMux.Unlock()

	if streams, exists := s.eventStreams[sessionID]; exists {
		for i, stream := range streams {
			if stream == targetChan {
				// Remove this stream from the slice
				s.eventStreams[sessionID] = append(streams[:i], streams[i+1:]...)
				close(targetChan)
				break
			}
		}

		// Clean up empty session entries
		if len(s.eventStreams[sessionID]) == 0 {
			delete(s.eventStreams, sessionID)
		}
	}
}

func (s *AgentServer) closeEventStreams(sessionID string) {
	s.streamsMux.Lock()
	defer s.streamsMux.Unlock()

	if streams, exists := s.eventStreams[sessionID]; exists {
		for _, stream := range streams {
			close(stream)
		}
		delete(s.eventStreams, sessionID)
	}
}

func (s *AgentServer) sessionCleanup() {
	ticker := time.NewTicker(1 * time.Minute) // Check every minute
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		s.sessionsMux.Lock()

		for sessionID, session := range s.sessions {
			if now.After(session.ExpiresAt) {
				s.logger.Info("Session expired, cleaning up",
					zap.String("session_id", sessionID),
					zap.String("agent_id", session.AgentID))

				delete(s.sessions, sessionID)
				
				// Close event streams for expired session
				go s.closeEventStreams(sessionID)

				// Broadcast session expired event
				go s.broadcastEvent(&agentpb.Event{
					EventId:       uuid.New().String(),
					Type:          agentpb.EventType_EVENT_TYPE_SESSION_EXPIRED,
					TimestampUnix: now.Unix(),
					SessionId:     sessionID,
					DataJson:      fmt.Sprintf(`{"agent_id": "%s", "reason": "expired"}`, session.AgentID),
				})
			}
		}

		s.sessionsMux.Unlock()
	}
}