package selflearn

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Collector handles the collection of execution feedback
type Collector struct {
	config  CollectionConfig
	storage Storage
	logger  *zap.Logger
}

// NewCollector creates a new feedback collector
func NewCollector(config CollectionConfig, storage Storage, logger *zap.Logger) *Collector {
	return &Collector{
		config:  config,
		storage: storage,
		logger:  logger,
	}
}

// ExecutionContext holds context information for a tool execution
type ExecutionContext struct {
	ToolName   string
	SourceType string
	UserAgent  string
	SessionID  string
	RequestID  string
	Metadata   map[string]interface{}
}

// CollectExecution captures feedback for a tool execution
func (c *Collector) CollectExecution(ctx context.Context, execCtx ExecutionContext, input interface{}, output interface{}, err error, duration time.Duration) error {
	if !c.config.Enabled {
		return nil
	}

	// Apply sampling rate
	if !c.shouldSample() {
		return nil
	}

	// Don't collect successful executions if configured not to
	if err == nil && !c.config.IncludeSuccessful {
		return nil
	}

	record := c.createExecutionRecord(execCtx, input, output, err, duration)

	if c.config.AsyncProcessing {
		// Process asynchronously to avoid blocking tool execution
		go func() {
			if storeErr := c.storage.StoreExecution(context.Background(), record); storeErr != nil {
				c.logger.Error("Failed to store execution record",
					zap.String("record_id", record.ID),
					zap.Error(storeErr))
			}
		}()
		return nil
	}

	// Synchronous processing
	return c.storage.StoreExecution(ctx, record)
}

// createExecutionRecord creates an execution record from the provided data
func (c *Collector) createExecutionRecord(execCtx ExecutionContext, input interface{}, output interface{}, err error, duration time.Duration) ExecutionRecord {
	recordID := c.generateID()
	
	record := ExecutionRecord{
		ID:         recordID,
		ToolName:   execCtx.ToolName,
		Timestamp:  time.Now().UTC(),
		Duration:   duration,
		Success:    err == nil,
		SourceType: execCtx.SourceType,
		Context:    make(map[string]interface{}),
	}

	// Add context metadata
	if execCtx.SessionID != "" {
		record.Context["session_id"] = execCtx.SessionID
	}
	if execCtx.RequestID != "" {
		record.Context["request_id"] = execCtx.RequestID
	}
	if execCtx.UserAgent != "" {
		record.Context["user_agent"] = execCtx.UserAgent
	}
	if execCtx.Metadata != nil {
		for k, v := range execCtx.Metadata {
			record.Context[k] = v
		}
	}

	// Process input/output if enabled
	if c.config.IncludeInputOutput {
		record.Input = c.sanitizeData(input, c.config.MaxInputSize)
		if output != nil {
			record.Output = c.sanitizeData(output, c.config.MaxOutputSize)
		}
	}

	// Process error information
	if err != nil {
		record.Error = err.Error()
		record.ErrorType = c.classifyError(err)
	}

	return record
}

// shouldSample determines if this execution should be sampled based on the sample rate
func (c *Collector) shouldSample() bool {
	if c.config.SampleRate >= 1.0 {
		return true
	}
	if c.config.SampleRate <= 0.0 {
		return false
	}

	// Simple random sampling
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	randomValue := float64(randomBytes[0]) / 255.0
	return randomValue < c.config.SampleRate
}

// classifyError attempts to classify the error into predefined types
func (c *Collector) classifyError(err error) ErrorType {
	if err == nil {
		return ""
	}

	errMsg := strings.ToLower(err.Error())

	// Network-related errors
	networkPatterns := []string{
		"connection", "network", "dns", "http", "tcp", "socket",
		"refused", "unreachable", "reset", "broken pipe", "no route",
	}
	for _, pattern := range networkPatterns {
		if strings.Contains(errMsg, pattern) {
			return ErrorTypeNetwork
		}
	}

	// Validation errors
	validationPatterns := []string{
		"validation", "invalid", "required", "missing", "format", "schema",
		"constraint", "malformed", "parse", "unmarshal", "decode",
	}
	for _, pattern := range validationPatterns {
		if strings.Contains(errMsg, pattern) {
			return ErrorTypeValidation
		}
	}

	// Configuration errors
	configPatterns := []string{
		"config", "credential", "auth", "permission", "unauthorized", "forbidden",
		"access denied", "not found", "endpoint", "url", "path",
	}
	for _, pattern := range configPatterns {
		if strings.Contains(errMsg, pattern) {
			return ErrorTypeConfiguration
		}
	}

	// Performance errors
	performancePatterns := []string{
		"timeout", "deadline", "slow", "rate limit", "throttle", "capacity",
		"overload", "busy", "queue", "limit exceeded",
	}
	for _, pattern := range performancePatterns {
		if strings.Contains(errMsg, pattern) {
			return ErrorTypePerformance
		}
	}

	// Logic errors
	logicPatterns := []string{
		"logic", "business", "rule", "condition", "state", "workflow",
		"sequence", "order", "dependency",
	}
	for _, pattern := range logicPatterns {
		if strings.Contains(errMsg, pattern) {
			return ErrorTypeLogic
		}
	}

	return ErrorTypeUnknown
}

// sanitizeData sanitizes and truncates data for storage
func (c *Collector) sanitizeData(data interface{}, maxSize int) interface{} {
	if data == nil {
		return nil
	}

	// Apply PII filtering if enabled
	if c.config.PIIFilterEnabled {
		data = c.filterPII(data)
	}

	// Convert to string for size checking and truncation
	dataStr := fmt.Sprintf("%v", data)
	if len(dataStr) > maxSize {
		return dataStr[:maxSize] + "... [truncated]"
	}

	return data
}

// filterPII applies basic PII filtering to the data
func (c *Collector) filterPII(data interface{}) interface{} {
	if data == nil {
		return nil
	}

	// Define PII patterns (simplified)
	piiPatterns := []*regexp.Regexp{
		regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`), // email
		regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`),                                    // SSN
		regexp.MustCompile(`\b\d{4}\s?\d{4}\s?\d{4}\s?\d{4}\b`),                       // credit card
		regexp.MustCompile(`\b\d{3}-\d{3}-\d{4}\b`),                                   // phone
	}

	// Convert to string for pattern matching
	dataStr := fmt.Sprintf("%v", data)
	
	// Apply PII masking
	for _, pattern := range piiPatterns {
		dataStr = pattern.ReplaceAllString(dataStr, "[REDACTED]")
	}

	// Try to maintain original data type if possible
	typ := reflect.TypeOf(data)
	if typ == nil {
		// If data is a typed nil, return the filtered string representation
		return dataStr
	}
	switch typ.Kind() {
	case reflect.String:
		return dataStr
	case reflect.Map, reflect.Slice, reflect.Struct:
		// For complex types, return the filtered string representation
		return dataStr
	default:
		return data
	}
}

// generateID generates a unique ID for execution records
func (c *Collector) generateID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		c.logger.Error("Failed to generate random ID", zap.Error(err))
		return ""
	}
	return hex.EncodeToString(bytes)
}

// UpdateConfig updates the collector configuration
func (c *Collector) UpdateConfig(config CollectionConfig) {
	c.config = config
	c.logger.Info("Collector configuration updated",
		zap.Bool("enabled", config.Enabled),
		zap.Float64("sample_rate", config.SampleRate))
}

// GetConfig returns the current collector configuration
func (c *Collector) GetConfig() CollectionConfig {
	return c.config
}