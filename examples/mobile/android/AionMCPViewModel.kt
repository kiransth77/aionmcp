package com.example.aionmcp.ui.viewmodel

import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.example.aionmcp.data.models.*
import com.example.aionmcp.data.repository.AionMCPRepository
import kotlinx.coroutines.launch

/**
 * ViewModel for AionMCP integration
 * Manages UI state and coordinates repository operations
 */
class AionMCPViewModel(
    private val repository: AionMCPRepository
) : ViewModel() {
    
    // Health status
    private val _healthStatus = MutableLiveData<Result<HealthResponse>>()
    val healthStatus: LiveData<Result<HealthResponse>> = _healthStatus
    
    // Tools list
    private val _tools = MutableLiveData<Result<List<Tool>>>()
    val tools: LiveData<Result<List<Tool>>> = _tools
    
    // Tool invocation result
    private val _toolResult = MutableLiveData<Result<ToolInvokeResponse>>()
    val toolResult: LiveData<Result<ToolInvokeResponse>> = _toolResult
    
    // Learning statistics
    private val _learningStats = MutableLiveData<Result<LearningStatsResponse>>()
    val learningStats: LiveData<Result<LearningStatsResponse>> = _learningStats
    
    // Insights
    private val _insights = MutableLiveData<Result<List<Insight>>>()
    val insights: LiveData<Result<List<Insight>>> = _insights
    
    // Loading state
    private val _isLoading = MutableLiveData<Boolean>()
    val isLoading: LiveData<Boolean> = _isLoading
    
    /**
     * Check server health
     */
    fun checkHealth() {
        viewModelScope.launch {
            _isLoading.value = true
            val result = repository.checkHealth()
            _healthStatus.value = result
            _isLoading.value = false
        }
    }
    
    /**
     * List all available tools
     */
    fun listTools(forceRefresh: Boolean = false) {
        viewModelScope.launch {
            _isLoading.value = true
            val result = repository.listTools(forceRefresh)
            _tools.value = result
            _isLoading.value = false
        }
    }
    
    /**
     * Invoke a tool
     */
    fun invokeTool(toolName: String, args: Map<String, Any>) {
        viewModelScope.launch {
            _isLoading.value = true
            val result = repository.invokeTool(toolName, args)
            _toolResult.value = result
            _isLoading.value = false
        }
    }
    
    /**
     * Get learning statistics
     */
    fun getLearningStats() {
        viewModelScope.launch {
            _isLoading.value = true
            val result = repository.getLearningStats()
            _learningStats.value = result
            _isLoading.value = false
        }
    }
    
    /**
     * Get insights
     */
    fun getInsights(type: String? = null, priority: String? = null) {
        viewModelScope.launch {
            _isLoading.value = true
            val result = repository.getInsights(type, priority)
            _insights.value = result
            _isLoading.value = false
        }
    }
    
    /**
     * Get insights for a specific tool
     */
    fun getToolInsights(toolName: String) {
        viewModelScope.launch {
            _isLoading.value = true
            val result = repository.getToolInsights(toolName)
            _insights.value = result
            _isLoading.value = false
        }
    }
    
    /**
     * Trigger pattern analysis
     */
    fun analyzePatterns() {
        viewModelScope.launch {
            _isLoading.value = true
            repository.analyzePatterns()
            _isLoading.value = false
        }
    }
    
    /**
     * Clear cache
     */
    fun clearCache() {
        repository.clearCache()
    }
}
