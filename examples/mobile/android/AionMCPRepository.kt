package com.example.aionmcp.data.repository

import com.example.aionmcp.data.api.AionMCPApiService
import com.example.aionmcp.data.models.*
import kotlinx.coroutines.delay
import java.io.IOException

/**
 * Repository for AionMCP API operations
 * Implements error handling, retry logic, and data caching
 */
class AionMCPRepository(
    private val api: AionMCPApiService
) {
    // Simple in-memory cache
    private var toolsCache: List<Tool>? = null
    private var toolsCacheTimestamp: Long = 0
    private val cacheValidityMs = 5 * 60 * 1000L // 5 minutes
    
    /**
     * Check server health
     */
    suspend fun checkHealth(): Result<HealthResponse> = safeApiCall {
        api.getHealth()
    }
    
    /**
     * List all available tools with caching
     */
    suspend fun listTools(forceRefresh: Boolean = false): Result<List<Tool>> {
        // Return cached data if valid
        if (!forceRefresh && isCacheValid()) {
            toolsCache?.let { return Result.success(it) }
        }
        
        return safeApiCall {
            val response = api.listTools()
            toolsCache = response.tools
            toolsCacheTimestamp = System.currentTimeMillis()
            response.tools
        }
    }
    
    /**
     * Invoke a tool with retry logic
     */
    suspend fun invokeTool(
        toolName: String,
        args: Map<String, Any>,
        maxRetries: Int = 3
    ): Result<ToolInvokeResponse> {
        return retryApiCall(maxRetries) {
            api.invokeTool(ToolInvokeRequest(toolName, args))
        }
    }
    
    /**
     * Get learning statistics
     */
    suspend fun getLearningStats(): Result<LearningStatsResponse> = safeApiCall {
        api.getLearningStats()
    }
    
    /**
     * Get learning insights
     */
    suspend fun getInsights(
        type: String? = null,
        priority: String? = null
    ): Result<List<Insight>> = safeApiCall {
        api.getInsights(type, priority).insights
    }
    
    /**
     * Get insights for a specific tool
     */
    suspend fun getToolInsights(toolName: String): Result<List<Insight>> = safeApiCall {
        api.getToolInsights(toolName).insights
    }
    
    /**
     * Trigger pattern analysis
     */
    suspend fun analyzePatterns(): Result<Map<String, Any>> = safeApiCall {
        api.analyzePatterns()
    }
    
    /**
     * Clear the tools cache
     */
    fun clearCache() {
        toolsCache = null
        toolsCacheTimestamp = 0
    }
    
    /**
     * Check if cached data is still valid
     */
    private fun isCacheValid(): Boolean {
        return toolsCache != null && 
               (System.currentTimeMillis() - toolsCacheTimestamp) < cacheValidityMs
    }
    
    /**
     * Safe API call wrapper with error handling
     */
    private suspend fun <T> safeApiCall(
        apiCall: suspend () -> T
    ): Result<T> {
        return try {
            Result.success(apiCall())
        } catch (e: IOException) {
            Result.failure(NetworkException("Network error: ${e.message}", e))
        } catch (e: Exception) {
            Result.failure(ApiException("API error: ${e.message}", e))
        }
    }
    
    /**
     * API call with retry logic
     */
    private suspend fun <T> retryApiCall(
        maxRetries: Int,
        apiCall: suspend () -> T
    ): Result<T> {
        var lastException: Exception? = null
        
        repeat(maxRetries) { attempt ->
            try {
                return Result.success(apiCall())
            } catch (e: IOException) {
                lastException = NetworkException("Network error: ${e.message}", e)
                // Exponential backoff
                if (attempt < maxRetries - 1) {
                    delay(1000L * (attempt + 1))
                }
            } catch (e: Exception) {
                // Don't retry on non-network errors
                return Result.failure(ApiException("API error: ${e.message}", e))
            }
        }
        
        return Result.failure(lastException ?: Exception("Unknown error"))
    }
}

/**
 * Custom exception for network errors
 */
class NetworkException(message: String, cause: Throwable? = null) : Exception(message, cause)

/**
 * Custom exception for API errors
 */
class ApiException(message: String, cause: Throwable? = null) : Exception(message, cause)
