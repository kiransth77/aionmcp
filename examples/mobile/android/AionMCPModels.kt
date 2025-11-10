package com.example.aionmcp.data.models

import com.google.gson.annotations.SerializedName

/**
 * Health check response from AionMCP server
 */
data class HealthResponse(
    @SerializedName("status")
    val status: String,
    
    @SerializedName("timestamp")
    val timestamp: Long,
    
    @SerializedName("version")
    val version: String,
    
    @SerializedName("iteration")
    val iteration: String? = null
)

/**
 * Tool metadata
 */
data class Tool(
    @SerializedName("name")
    val name: String,
    
    @SerializedName("description")
    val description: String,
    
    @SerializedName("metadata")
    val metadata: Map<String, Any>? = null,
    
    @SerializedName("source")
    val source: String? = null,
    
    @SerializedName("type")
    val type: String? = null
)

/**
 * Response when listing available tools
 */
data class ToolsResponse(
    @SerializedName("protocol")
    val protocol: String,
    
    @SerializedName("tools")
    val tools: List<Tool>
)

/**
 * Request to invoke a tool
 */
data class ToolInvokeRequest(
    @SerializedName("tool")
    val tool: String,
    
    @SerializedName("args")
    val args: Map<String, Any>
)

/**
 * Response from tool invocation
 */
data class ToolInvokeResponse(
    @SerializedName("tool")
    val tool: String,
    
    @SerializedName("result")
    val result: Any
)

/**
 * Learning statistics from the server
 */
data class LearningStatsResponse(
    @SerializedName("total_executions")
    val totalExecutions: Int,
    
    @SerializedName("success_rate")
    val successRate: Double,
    
    @SerializedName("avg_latency_ms")
    val avgLatencyMs: Double,
    
    @SerializedName("tool_stats")
    val toolStats: Map<String, ToolStats>? = null
)

/**
 * Statistics for individual tools
 */
data class ToolStats(
    @SerializedName("executions")
    val executions: Int,
    
    @SerializedName("success_rate")
    val successRate: Double,
    
    @SerializedName("avg_latency_ms")
    val avgLatencyMs: Double
)

/**
 * Learning insight
 */
data class Insight(
    @SerializedName("id")
    val id: String,
    
    @SerializedName("type")
    val type: String,
    
    @SerializedName("priority")
    val priority: String,
    
    @SerializedName("title")
    val title: String,
    
    @SerializedName("description")
    val description: String,
    
    @SerializedName("recommendation")
    val recommendation: String? = null,
    
    @SerializedName("created_at")
    val createdAt: String
)

/**
 * Insights response
 */
data class InsightsResponse(
    @SerializedName("insights")
    val insights: List<Insight>
)

/**
 * API error response
 */
data class ErrorResponse(
    @SerializedName("error")
    val error: String,
    
    @SerializedName("code")
    val code: String? = null,
    
    @SerializedName("details")
    val details: Map<String, Any>? = null
)
