import Foundation
import Alamofire

/// AionMCP REST API Client
class AionMCPClient {
    private let baseURL: String
    private let session: Session
    private let apiKey: String?
    
    init(baseURL: String, apiKey: String? = nil) {
        self.baseURL = baseURL
        self.apiKey = apiKey
        
        // Configure session
        let configuration = URLSessionConfiguration.default
        configuration.timeoutIntervalForRequest = 30
        configuration.timeoutIntervalForResource = 60
        
        // Add interceptor for API key if provided
        var interceptor: RequestInterceptor?
        if let apiKey = apiKey {
            interceptor = APIKeyInterceptor(apiKey: apiKey)
        }
        
        self.session = Session(
            configuration: configuration,
            interceptor: interceptor
        )
    }
    
    // MARK: - Health
    
    func getHealth() async throws -> HealthResponse {
        return try await request(
            endpoint: "/api/v1/health",
            method: .get
        )
    }
    
    // MARK: - Tools
    
    func listTools() async throws -> ToolsResponse {
        return try await request(
            endpoint: "/api/v1/tools/list",
            method: .get
        )
    }
    
    func invokeTool(name: String, args: [String: Any]) async throws -> ToolInvokeResponse {
        let request = ToolInvokeRequest(
            tool: name,
            args: args.mapValues { AnyCodable($0) }
        )
        
        return try await request(
            endpoint: "/api/v1/tools/invoke",
            method: .post,
            parameters: request
        )
    }
    
    // MARK: - Learning
    
    func getLearningStats() async throws -> LearningStatsResponse {
        return try await request(
            endpoint: "/api/v1/learning/stats",
            method: .get
        )
    }
    
    func getInsights(type: String? = nil, priority: String? = nil) async throws -> InsightsResponse {
        var queryParams: [String: String] = [:]
        if let type = type { queryParams["type"] = type }
        if let priority = priority { queryParams["priority"] = priority }
        
        return try await request(
            endpoint: "/api/v1/learning/insights",
            method: .get,
            queryParameters: queryParams
        )
    }
    
    func getToolInsights(toolName: String) async throws -> InsightsResponse {
        return try await request(
            endpoint: "/api/v1/learning/tools/\(toolName)/insights",
            method: .get
        )
    }
    
    func analyzePatterns() async throws -> [String: Any] {
        return try await request(
            endpoint: "/api/v1/learning/analyze",
            method: .post
        )
    }
    
    // MARK: - Private Request Handler
    
    private func request<T: Decodable>(
        endpoint: String,
        method: HTTPMethod,
        parameters: Encodable? = nil,
        queryParameters: [String: String]? = nil
    ) async throws -> T {
        let url = baseURL + endpoint
        
        do {
            let dataRequest: DataRequest
            
            if let parameters = parameters {
                dataRequest = session.request(
                    url,
                    method: method,
                    parameters: parameters,
                    encoder: JSONParameterEncoder.default
                )
            } else if let queryParams = queryParameters {
                dataRequest = session.request(
                    url,
                    method: method,
                    parameters: queryParams
                )
            } else {
                dataRequest = session.request(url, method: method)
            }
            
            let response = await dataRequest
                .validate()
                .serializingDecodable(T.self)
                .response
            
            switch response.result {
            case .success(let value):
                return value
                
            case .failure(let error):
                // Try to decode error response
                if let data = response.data,
                   let errorResponse = try? JSONDecoder().decode(ErrorResponse.self, from: data) {
                    throw AionMCPError.apiError(errorResponse.error)
                }
                
                // Handle HTTP status codes
                if let statusCode = response.response?.statusCode {
                    switch statusCode {
                    case 401:
                        throw AionMCPError.unauthorized
                    case 500...599:
                        throw AionMCPError.serverError(statusCode)
                    default:
                        throw AionMCPError.networkError(error)
                    }
                }
                
                throw AionMCPError.networkError(error)
            }
        } catch let error as AionMCPError {
            throw error
        } catch {
            throw AionMCPError.networkError(error)
        }
    }
}

// MARK: - API Key Interceptor

private class APIKeyInterceptor: RequestInterceptor {
    private let apiKey: String
    
    init(apiKey: String) {
        self.apiKey = apiKey
    }
    
    func adapt(_ urlRequest: URLRequest, for session: Session, completion: @escaping (Result<URLRequest, Error>) -> Void) {
        var urlRequest = urlRequest
        urlRequest.headers.add(name: "X-API-Key", value: apiKey)
        completion(.success(urlRequest))
    }
}
