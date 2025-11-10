import Foundation

// MARK: - Health Response

struct HealthResponse: Codable {
    let status: String
    let timestamp: Int64
    let version: String
    let iteration: String?
}

// MARK: - Tool

struct Tool: Codable, Identifiable {
    var id: String { name }
    
    let name: String
    let description: String
    let metadata: [String: AnyCodable]?
    let source: String?
    let type: String?
}

// MARK: - Tools Response

struct ToolsResponse: Codable {
    let `protocol`: String
    let tools: [Tool]
}

// MARK: - Tool Invoke Request

struct ToolInvokeRequest: Codable {
    let tool: String
    let args: [String: AnyCodable]
}

// MARK: - Tool Invoke Response

struct ToolInvokeResponse: Codable {
    let tool: String
    let result: AnyCodable
}

// MARK: - Learning Stats Response

struct LearningStatsResponse: Codable {
    let totalExecutions: Int
    let successRate: Double
    let avgLatencyMs: Double
    let toolStats: [String: ToolStats]?
    
    enum CodingKeys: String, CodingKey {
        case totalExecutions = "total_executions"
        case successRate = "success_rate"
        case avgLatencyMs = "avg_latency_ms"
        case toolStats = "tool_stats"
    }
}

// MARK: - Tool Stats

struct ToolStats: Codable {
    let executions: Int
    let successRate: Double
    let avgLatencyMs: Double
    
    enum CodingKeys: String, CodingKey {
        case executions
        case successRate = "success_rate"
        case avgLatencyMs = "avg_latency_ms"
    }
}

// MARK: - Insight

struct Insight: Codable, Identifiable {
    let id: String
    let type: String
    let priority: String
    let title: String
    let description: String
    let recommendation: String?
    let createdAt: String
    
    enum CodingKeys: String, CodingKey {
        case id, type, priority, title, description, recommendation
        case createdAt = "created_at"
    }
}

// MARK: - Insights Response

struct InsightsResponse: Codable {
    let insights: [Insight]
}

// MARK: - Error Response

struct ErrorResponse: Codable {
    let error: String
    let code: String?
    let details: [String: AnyCodable]?
}

// MARK: - AnyCodable Helper

/// Helper type for encoding/decoding Any values
struct AnyCodable: Codable {
    let value: Any
    
    init(_ value: Any) {
        self.value = value
    }
    
    init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()
        
        if let int = try? container.decode(Int.self) {
            value = int
        } else if let double = try? container.decode(Double.self) {
            value = double
        } else if let string = try? container.decode(String.self) {
            value = string
        } else if let bool = try? container.decode(Bool.self) {
            value = bool
        } else if let array = try? container.decode([AnyCodable].self) {
            value = array.map { $0.value }
        } else if let dict = try? container.decode([String: AnyCodable].self) {
            value = dict.mapValues { $0.value }
        } else {
            value = NSNull()
        }
    }
    
    func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()
        
        switch value {
        case let int as Int:
            try container.encode(int)
        case let int64 as Int64:
            try container.encode(int64)
        case let double as Double:
            try container.encode(double)
        case let float as Float:
            try container.encode(float)
        case let string as String:
            try container.encode(string)
        case let bool as Bool:
            try container.encode(bool)
        case let array as [Any]:
            try container.encode(array.map { AnyCodable($0) })
        case let dict as [String: Any]:
            try container.encode(dict.mapValues { AnyCodable($0) })
        case is NSNull:
            try container.encodeNil()
        default:
            try container.encodeNil()
        }
    }
}

// MARK: - Custom Errors

enum AionMCPError: LocalizedError {
    case networkError(Error)
    case apiError(String)
    case decodingError(Error)
    case invalidResponse
    case unauthorized
    case serverError(Int)
    
    var errorDescription: String? {
        switch self {
        case .networkError(let error):
            return "Network error: \(error.localizedDescription)"
        case .apiError(let message):
            return "API error: \(message)"
        case .decodingError(let error):
            return "Failed to decode response: \(error.localizedDescription)"
        case .invalidResponse:
            return "Invalid response from server"
        case .unauthorized:
            return "Unauthorized access"
        case .serverError(let code):
            return "Server error with code: \(code)"
        }
    }
}
