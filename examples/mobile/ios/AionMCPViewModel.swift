import Foundation
import SwiftUI

/// ViewModel for AionMCP integration
@MainActor
class AionMCPViewModel: ObservableObject {
    private let repository: AionMCPRepository
    
    // MARK: - Published Properties
    
    @Published var health: HealthResponse?
    @Published var tools: [Tool] = []
    @Published var toolResult: ToolInvokeResponse?
    @Published var learningStats: LearningStatsResponse?
    @Published var insights: [Insight] = []
    @Published var isLoading = false
    @Published var errorMessage: String?
    
    // MARK: - Initialization
    
    init(repository: AionMCPRepository) {
        self.repository = repository
    }
    
    // MARK: - Health
    
    func checkHealth() async {
        isLoading = true
        errorMessage = nil
        
        do {
            health = try await repository.checkHealth()
        } catch {
            handleError(error)
        }
        
        isLoading = false
    }
    
    // MARK: - Tools
    
    func listTools(forceRefresh: Bool = false) async {
        isLoading = true
        errorMessage = nil
        
        do {
            tools = try await repository.listTools(forceRefresh: forceRefresh)
        } catch {
            handleError(error)
        }
        
        isLoading = false
    }
    
    func invokeTool(name: String, args: [String: Any]) async {
        isLoading = true
        errorMessage = nil
        toolResult = nil
        
        do {
            toolResult = try await repository.invokeTool(name: name, args: args)
        } catch {
            handleError(error)
        }
        
        isLoading = false
    }
    
    // MARK: - Learning
    
    func getLearningStats() async {
        isLoading = true
        errorMessage = nil
        
        do {
            learningStats = try await repository.getLearningStats()
        } catch {
            handleError(error)
        }
        
        isLoading = false
    }
    
    func getInsights(type: String? = nil, priority: String? = nil) async {
        isLoading = true
        errorMessage = nil
        
        do {
            insights = try await repository.getInsights(type: type, priority: priority)
        } catch {
            handleError(error)
        }
        
        isLoading = false
    }
    
    func getToolInsights(toolName: String) async {
        isLoading = true
        errorMessage = nil
        
        do {
            insights = try await repository.getToolInsights(toolName: toolName)
        } catch {
            handleError(error)
        }
        
        isLoading = false
    }
    
    func analyzePatterns() async {
        isLoading = true
        errorMessage = nil
        
        do {
            _ = try await repository.analyzePatterns()
        } catch {
            handleError(error)
        }
        
        isLoading = false
    }
    
    // MARK: - Cache Management
    
    func clearCache() async {
        await repository.clearCache()
        tools = []
        health = nil
        learningStats = nil
        insights = []
        toolResult = nil
    }
    
    // MARK: - Convenience Methods
    
    func refresh() async {
        await loadInitialData()
    }
    
    func loadInitialData() async {
        await checkHealth()
        await listTools()
    }
    
    // MARK: - Error Handling
    
    private func handleError(_ error: Error) {
        if let aionError = error as? AionMCPError {
            errorMessage = aionError.errorDescription
        } else {
            errorMessage = error.localizedDescription
        }
    }
}

// MARK: - Convenience Initializer

extension AionMCPViewModel {
    convenience init(baseURL: String, apiKey: String? = nil) {
        let client = AionMCPClient(baseURL: baseURL, apiKey: apiKey)
        let repository = AionMCPRepository(client: client)
        self.init(repository: repository)
    }
}
