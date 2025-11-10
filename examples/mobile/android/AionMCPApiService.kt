package com.example.aionmcp.data.api

import com.example.aionmcp.data.models.*
import retrofit2.http.*

/**
 * AionMCP REST API service interface
 */
interface AionMCPApiService {
    
    /**
     * Health check endpoint
     */
    @GET("api/v1/health")
    suspend fun getHealth(): HealthResponse
    
    /**
     * List all available tools
     */
    @GET("api/v1/tools/list")
    suspend fun listTools(): ToolsResponse
    
    /**
     * Invoke a specific tool
     */
    @POST("api/v1/tools/invoke")
    suspend fun invokeTool(
        @Body request: ToolInvokeRequest
    ): ToolInvokeResponse
    
    /**
     * Get learning statistics
     */
    @GET("api/v1/learning/stats")
    suspend fun getLearningStats(): LearningStatsResponse
    
    /**
     * Get learning insights
     */
    @GET("api/v1/learning/insights")
    suspend fun getInsights(
        @Query("type") type: String? = null,
        @Query("priority") priority: String? = null
    ): InsightsResponse
    
    /**
     * Get tool-specific insights
     */
    @GET("api/v1/learning/tools/{name}/insights")
    suspend fun getToolInsights(
        @Path("name") toolName: String
    ): InsightsResponse
    
    /**
     * Trigger manual pattern analysis
     */
    @POST("api/v1/learning/analyze")
    suspend fun analyzePatterns(): Map<String, Any>
}
