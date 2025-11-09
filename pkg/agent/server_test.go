package agent

import (
	"context"
	"fmt"
	"testing"
	"time"

	agentpb "github.com/aionmcp/aionmcp/pkg/agent/proto"
	"github.com/aionmcp/aionmcp/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockTool implements the types.Tool interface for testing
type MockTool struct {
	mock.Mock
}

func (m *MockTool) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockTool) Description() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockTool) Execute(input any) (any, error) {
	args := m.Called(input)
	return args.Get(0), args.Error(1)
}

func (m *MockTool) Metadata() types.ToolMetadata {
	args := m.Called()
	return args.Get(0).(types.ToolMetadata)
}

// MockToolRegistry implements the types.ToolRegistry interface for testing
type MockToolRegistry struct {
	mock.Mock
}

func (m *MockToolRegistry) Get(name string) (types.Tool, error) {
	args := m.Called(name)
	return args.Get(0).(types.Tool), args.Error(1)
}

func (m *MockToolRegistry) ListTools() []types.ToolMetadata {
	args := m.Called()
	return args.Get(0).([]types.ToolMetadata)
}

func (m *MockToolRegistry) Count() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockToolRegistry) Register(tool types.Tool) error {
	args := m.Called(tool)
	return args.Error(0)
}

func (m *MockToolRegistry) RegisterWithSource(tool types.Tool, sourceID, version string) error {
	args := m.Called(tool, sourceID, version)
	return args.Error(0)
}

func (m *MockToolRegistry) RegisterBatch(tools []types.Tool, sourceID string) error {
	args := m.Called(tools, sourceID)
	return args.Error(0)
}

func (m *MockToolRegistry) Unregister(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockToolRegistry) UnregisterBySource(sourceID string) error {
	args := m.Called(sourceID)
	return args.Error(0)
}

func (m *MockToolRegistry) GetVersion(name string) (string, error) {
	args := m.Called(name)
	return args.String(0), args.Error(1)
}

func (m *MockToolRegistry) GetSource(name string) (string, error) {
	args := m.Called(name)
	return args.String(0), args.Error(1)
}

func (m *MockToolRegistry) ListToolsBySource(sourceID string) []types.ToolMetadata {
	args := m.Called(sourceID)
	return args.Get(0).([]types.ToolMetadata)
}

func (m *MockToolRegistry) GetToolSources() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockToolRegistry) GetRegistryStats() map[string]interface{} {
	args := m.Called()
	return args.Get(0).(map[string]interface{})
}

func TestAgentServer_RegisterAgent(t *testing.T) {
	logger := zap.NewNop()
	mockRegistry := &MockToolRegistry{}
	server := NewAgentServer(logger, mockRegistry)

	// Setup mock expectations
	mockTools := []types.ToolMetadata{
		{
			Name:        "test-tool",
			Description: "A test tool",
			Version:     "1.0.0",
			Source:      "test",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	mockRegistry.On("ListTools").Return(mockTools)

	req := &agentpb.RegisterAgentRequest{
		AgentId:      "test-agent-1",
		AgentName:    "Test Agent",
		AgentVersion: "1.0.0",
		Capabilities: &agentpb.AgentCapabilities{
			SupportedProtocols: []string{"mcp/1.0"},
			SupportsStreaming:  true,
		},
		SessionTimeoutSeconds: 600,
	}

	resp, err := server.RegisterAgent(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.SessionId)
	assert.Greater(t, resp.ExpiresAtUnix, time.Now().Unix())
	assert.NotNil(t, resp.ServerInfo)
	assert.Equal(t, "0.1.0", resp.ServerInfo.ServerVersion)
	assert.Equal(t, "MCP/1.0", resp.ServerInfo.ProtocolVersion)
	assert.Len(t, resp.AvailableTools, 1)
	assert.Equal(t, "test-tool", resp.AvailableTools[0].Name)

	// Verify session was created
	session, exists := server.getSession(resp.SessionId)
	assert.True(t, exists)
	assert.Equal(t, "test-agent-1", session.AgentID)
	assert.Equal(t, "Test Agent", session.AgentName)

	mockRegistry.AssertExpectations(t)
}

func TestAgentServer_RegisterAgent_ValidationErrors(t *testing.T) {
	logger := zap.NewNop()
	mockRegistry := &MockToolRegistry{}
	server := NewAgentServer(logger, mockRegistry)

	tests := []struct {
		name          string
		req           *agentpb.RegisterAgentRequest
		expectedError string
	}{
		{
			name: "missing agent ID",
			req: &agentpb.RegisterAgentRequest{
				AgentName: "Test Agent",
			},
			expectedError: "agent_id is required",
		},
		{
			name: "missing agent name",
			req: &agentpb.RegisterAgentRequest{
				AgentId: "test-agent",
			},
			expectedError: "agent_name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := server.RegisterAgent(context.Background(), tt.req)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

func TestAgentServer_UnregisterAgent(t *testing.T) {
	logger := zap.NewNop()
	mockRegistry := &MockToolRegistry{}
	server := NewAgentServer(logger, mockRegistry)

	// First register an agent
	mockRegistry.On("ListTools").Return([]types.ToolMetadata{})
	registerReq := &agentpb.RegisterAgentRequest{
		AgentId:   "test-agent-1",
		AgentName: "Test Agent",
	}
	registerResp, err := server.RegisterAgent(context.Background(), registerReq)
	assert.NoError(t, err)

	// Now unregister the agent
	unregisterReq := &agentpb.UnregisterAgentRequest{
		SessionId: registerResp.SessionId,
	}
	unregisterResp, err := server.UnregisterAgent(context.Background(), unregisterReq)

	assert.NoError(t, err)
	assert.True(t, unregisterResp.Success)
	assert.Equal(t, "Agent session terminated successfully", unregisterResp.Message)

	// Verify session was removed
	_, exists := server.getSession(registerResp.SessionId)
	assert.False(t, exists)

	mockRegistry.AssertExpectations(t)
}

func TestAgentServer_ListTools(t *testing.T) {
	logger := zap.NewNop()
	mockRegistry := &MockToolRegistry{}
	server := NewAgentServer(logger, mockRegistry)

	// Register an agent first
	mockTools := []types.ToolMetadata{
		{
			Name:        "tool1",
			Description: "Tool 1",
			Version:     "1.0.0",
			Source:      "test",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "tool2",
			Description: "Tool 2",
			Version:     "1.1.0",
			Source:      "test",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	mockRegistry.On("ListTools").Return(mockTools)

	registerReq := &agentpb.RegisterAgentRequest{
		AgentId:   "test-agent-1",
		AgentName: "Test Agent",
	}
	registerResp, err := server.RegisterAgent(context.Background(), registerReq)
	assert.NoError(t, err)

	// List tools
	listReq := &agentpb.ListToolsRequest{
		SessionId: registerResp.SessionId,
	}
	listResp, err := server.ListTools(context.Background(), listReq)

	assert.NoError(t, err)
	assert.Len(t, listResp.Tools, 2)
	assert.Equal(t, int32(2), listResp.TotalCount)
	assert.Equal(t, "tool1", listResp.Tools[0].Name)
	assert.Equal(t, "tool2", listResp.Tools[1].Name)

	mockRegistry.AssertExpectations(t)
}

func TestAgentServer_InvokeTool(t *testing.T) {
	logger := zap.NewNop()
	mockRegistry := &MockToolRegistry{}
	mockTool := &MockTool{}
	server := NewAgentServer(logger, mockRegistry)

	// Register an agent first
	mockRegistry.On("ListTools").Return([]types.ToolMetadata{})
	registerReq := &agentpb.RegisterAgentRequest{
		AgentId:   "test-agent-1",
		AgentName: "Test Agent",
	}
	registerResp, err := server.RegisterAgent(context.Background(), registerReq)
	assert.NoError(t, err)

	// Setup tool execution
	mockRegistry.On("Get", "test-tool").Return(mockTool, nil)
	mockTool.On("Execute", mock.Anything).Return(map[string]interface{}{
		"result": "success",
		"value":  42,
	}, nil)

	// Invoke tool
	invokeReq := &agentpb.InvokeToolRequest{
		SessionId:      registerResp.SessionId,
		ToolName:       "test-tool",
		InvocationId:   "test-invocation-1",
		ParametersJson: `{"param1": "value1"}`,
	}
	invokeResp, err := server.InvokeTool(context.Background(), invokeReq)

	assert.NoError(t, err)
	assert.Equal(t, "test-invocation-1", invokeResp.InvocationId)
	assert.Equal(t, agentpb.ToolInvocationStatus_TOOL_INVOCATION_STATUS_SUCCESS, invokeResp.Status)
	assert.NotEmpty(t, invokeResp.ResultJson)
	assert.Nil(t, invokeResp.Error)
	assert.NotNil(t, invokeResp.Metrics)
	assert.GreaterOrEqual(t, invokeResp.Metrics.ExecutionTimeMs, int64(0))

	mockRegistry.AssertExpectations(t)
	mockTool.AssertExpectations(t)
}

func TestAgentServer_InvokeTool_NotFound(t *testing.T) {
	logger := zap.NewNop()
	mockRegistry := &MockToolRegistry{}
	server := NewAgentServer(logger, mockRegistry)

	// Register an agent first
	mockRegistry.On("ListTools").Return([]types.ToolMetadata{})
	registerReq := &agentpb.RegisterAgentRequest{
		AgentId:   "test-agent-1",
		AgentName: "Test Agent",
	}
	registerResp, err := server.RegisterAgent(context.Background(), registerReq)
	assert.NoError(t, err)

	// Setup tool not found
	mockRegistry.On("Get", "nonexistent-tool").Return((*MockTool)(nil), assert.AnError)

	// Invoke nonexistent tool
	invokeReq := &agentpb.InvokeToolRequest{
		SessionId:      registerResp.SessionId,
		ToolName:       "nonexistent-tool",
		InvocationId:   "test-invocation-1",
		ParametersJson: `{}`,
	}
	_, err = server.InvokeTool(context.Background(), invokeReq)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tool not found")

	mockRegistry.AssertExpectations(t)
}

func TestAgentServer_HeartBeat(t *testing.T) {
	logger := zap.NewNop()
	mockRegistry := &MockToolRegistry{}
	server := NewAgentServer(logger, mockRegistry)

	// Register an agent first
	mockRegistry.On("ListTools").Return([]types.ToolMetadata{})
	registerReq := &agentpb.RegisterAgentRequest{
		AgentId:   "test-agent-1",
		AgentName: "Test Agent",
	}
	registerResp, err := server.RegisterAgent(context.Background(), registerReq)
	assert.NoError(t, err)

	// Send heartbeat
	heartbeatReq := &agentpb.HeartBeatRequest{
		SessionId: registerResp.SessionId,
		Status:    agentpb.AgentStatus_AGENT_STATUS_ACTIVE,
	}
	heartbeatResp, err := server.HeartBeat(context.Background(), heartbeatReq)

	assert.NoError(t, err)
	assert.True(t, heartbeatResp.SessionValid)
	assert.Greater(t, heartbeatResp.NextHeartbeatAtUnix, time.Now().Unix())

	mockRegistry.AssertExpectations(t)
}

func TestAgentServer_HeartBeat_InvalidSession(t *testing.T) {
	logger := zap.NewNop()
	mockRegistry := &MockToolRegistry{}
	server := NewAgentServer(logger, mockRegistry)

	// Send heartbeat with invalid session
	heartbeatReq := &agentpb.HeartBeatRequest{
		SessionId: "invalid-session-id",
		Status:    agentpb.AgentStatus_AGENT_STATUS_ACTIVE,
	}
	heartbeatResp, err := server.HeartBeat(context.Background(), heartbeatReq)

	assert.NoError(t, err)
	assert.False(t, heartbeatResp.SessionValid)
}

func TestAgentServer_GetAgentStatus(t *testing.T) {
	logger := zap.NewNop()
	mockRegistry := &MockToolRegistry{}
	server := NewAgentServer(logger, mockRegistry)

	// Register an agent first
	mockRegistry.On("ListTools").Return([]types.ToolMetadata{})
	registerReq := &agentpb.RegisterAgentRequest{
		AgentId:      "test-agent-1",
		AgentName:    "Test Agent",
		AgentVersion: "1.0.0",
		Capabilities: &agentpb.AgentCapabilities{
			SupportedProtocols: []string{"mcp/1.0"},
		},
	}
	registerResp, err := server.RegisterAgent(context.Background(), registerReq)
	assert.NoError(t, err)

	// Get agent status
	statusReq := &agentpb.GetAgentStatusRequest{
		SessionId: registerResp.SessionId,
	}
	statusResp, err := server.GetAgentStatus(context.Background(), statusReq)

	assert.NoError(t, err)
	assert.NotNil(t, statusResp.SessionInfo)
	assert.Equal(t, registerResp.SessionId, statusResp.SessionInfo.SessionId)
	assert.Equal(t, "test-agent-1", statusResp.SessionInfo.AgentId)
	assert.Equal(t, "Test Agent", statusResp.SessionInfo.AgentName)
	assert.Equal(t, "1.0.0", statusResp.SessionInfo.AgentVersion)
	assert.NotNil(t, statusResp.Metrics)
	assert.Equal(t, int64(0), statusResp.Metrics.TotalInvocations)
	assert.NotNil(t, statusResp.RecentToolUsage)

	mockRegistry.AssertExpectations(t)
}

// Benchmark tests
func BenchmarkAgentServer_RegisterAgent(b *testing.B) {
	logger := zap.NewNop()
	mockRegistry := &MockToolRegistry{}
	server := NewAgentServer(logger, mockRegistry)

	mockRegistry.On("ListTools").Return([]types.ToolMetadata{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := &agentpb.RegisterAgentRequest{
			AgentId:   fmt.Sprintf("agent-%d", i),
			AgentName: "Benchmark Agent",
		}
		_, _ = server.RegisterAgent(context.Background(), req)
	}
}

func BenchmarkAgentServer_ListTools(b *testing.B) {
	logger := zap.NewNop()
	mockRegistry := &MockToolRegistry{}
	server := NewAgentServer(logger, mockRegistry)

	// Create many tools for benchmark
	tools := make([]types.ToolMetadata, 1000)
	for i := 0; i < 1000; i++ {
		tools[i] = types.ToolMetadata{
			Name:        fmt.Sprintf("tool-%d", i),
			Description: fmt.Sprintf("Tool %d", i),
			Version:     "1.0.0",
		}
	}

	mockRegistry.On("ListTools").Return(tools)

	// Register an agent
	registerReq := &agentpb.RegisterAgentRequest{
		AgentId:   "benchmark-agent",
		AgentName: "Benchmark Agent",
	}
	registerResp, _ := server.RegisterAgent(context.Background(), registerReq)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := &agentpb.ListToolsRequest{
			SessionId: registerResp.SessionId,
		}
		_, _ = server.ListTools(context.Background(), req)
	}
}
