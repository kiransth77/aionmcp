import Foundation

/// Repository for AionMCP operations with error handling and caching
actor AionMCPRepository {
    private let client: AionMCPClient
    private var toolsCache: [Tool]?
    private var cacheTimestamp: Date?
    private let cacheValiditySeconds: TimeInterval = 300 // 5 minutes
    
    init(client: AionMCPClient) {
        self.client = client
    }
    
    // MARK: - Health
    
    func checkHealth() async throws -> HealthResponse {
        return try await client.getHealth()
    }
    
    // MARK: - Tools
    
    func listTools(forceRefresh: Bool = false) async throws -> [Tool] {
        // Return cached data if valid
        if !forceRefresh, let cached = toolsCache, isCacheValid() {
            return cached
        }
        
        // Fetch from API
        let response = try await client.listTools()
        toolsCache = response.tools
        cacheTimestamp = Date()
        return response.tools
    }
    
    func invokeTool(name: String, args: [String: Any], maxRetries: Int = 3) async throws -> ToolInvokeResponse {
        return try await retryOperation(maxRetries: maxRetries) {
            try await self.client.invokeTool(name: name, args: args)
        }
    }
    
    // MARK: - Learning
    
    func getLearningStats() async throws -> LearningStatsResponse {
        return try await client.getLearningStats()
    }
    
    func getInsights(type: String? = nil, priority: String? = nil) async throws -> [Insight] {
        let response = try await client.getInsights(type: type, priority: priority)
        return response.insights
    }
    
    func getToolInsights(toolName: String) async throws -> [Insight] {
        let response = try await client.getToolInsights(toolName: toolName)
        return response.insights
    }
    
    func analyzePatterns() async throws -> [String: Any] {
        return try await client.analyzePatterns()
    }
    
    // MARK: - Cache Management
    
    func clearCache() {
        toolsCache = nil
        cacheTimestamp = nil
    }
    
    private func isCacheValid() -> Bool {
        guard let timestamp = cacheTimestamp else { return false }
        return Date().timeIntervalSince(timestamp) < cacheValiditySeconds
    }
    
    // MARK: - Retry Logic
    
    private func retryOperation<T>(
        maxRetries: Int,
        operation: @escaping () async throws -> T
    ) async throws -> T {
        var lastError: Error?
        
        for attempt in 0..<maxRetries {
            do {
                return try await operation()
            } catch let error as AionMCPError {
                lastError = error
                
                // Don't retry on certain errors
                switch error {
                case .unauthorized:
                    throw error
                case .networkError:
                    // Only retry network errors
                    if attempt < maxRetries - 1 {
                        // Exponential backoff
                        let delay = UInt64(1_000_000_000 * (attempt + 1))
                        try await Task.sleep(nanoseconds: delay)
                    }
                default:
                    throw error
                }
            } catch {
                lastError = error
                if attempt < maxRetries - 1 {
                    let delay = UInt64(1_000_000_000 * (attempt + 1))
                    try await Task.sleep(nanoseconds: delay)
                }
            }
        }
        
        throw lastError ?? AionMCPError.invalidResponse
    }
}
