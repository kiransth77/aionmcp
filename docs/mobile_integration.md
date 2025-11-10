# Mobile Integration Guide

## Overview

AionMCP can be integrated into Android and iOS mobile applications through its REST API and gRPC interfaces. This guide provides comprehensive instructions for implementing AionMCP clients on mobile platforms.

## Architecture

### Client-Server Model

Mobile apps communicate with AionMCP server through two primary protocols:

1. **REST API (HTTP/HTTPS)**: Best for simple integrations, standard HTTP clients
2. **gRPC**: Better for high-performance, streaming, and efficient binary communication

```
┌─────────────────┐         ┌─────────────────┐
│   Mobile App    │         │  AionMCP Server │
│  (iOS/Android)  │ ◄─────► │   (Go Backend)  │
│                 │  HTTP/  │                 │
│  - REST Client  │  gRPC   │  - Tool Registry│
│  - gRPC Client  │         │  - Learning Eng.│
│  - Auth Handler │         │  - API Specs    │
└─────────────────┘         └─────────────────┘
```

### Key Endpoints for Mobile

#### REST API Endpoints (HTTP)

- `GET /api/v1/health` - Health check
- `GET /api/v1/tools/list` - List available tools
- `POST /api/v1/tools/invoke` - Invoke a tool
- `GET /api/v1/learning/stats` - Get learning statistics
- `GET /api/v1/learning/insights` - Get system insights

#### gRPC Services

- `AgentService.RegisterAgent` - Register mobile client
- `AgentService.ListTools` - List available tools
- `AgentService.InvokeTool` - Execute tools
- `AgentService.StreamToolExecution` - Stream tool results

## Android Integration

### Prerequisites

```gradle
dependencies {
    // Networking
    implementation 'com.squareup.retrofit2:retrofit:2.9.0'
    implementation 'com.squareup.retrofit2:converter-gson:2.9.0'
    implementation 'com.squareup.okhttp3:logging-interceptor:4.11.0'
    
    // gRPC (optional, for high-performance)
    implementation 'io.grpc:grpc-okhttp:1.58.0'
    implementation 'io.grpc:grpc-protobuf-lite:1.58.0'
    implementation 'io.grpc:grpc-stub:1.58.0'
    
    // Coroutines for async operations
    implementation 'org.jetbrains.kotlinx:kotlinx-coroutines-android:1.7.3'
}
```

### REST API Client (Kotlin)

```kotlin
// AionMCPApiService.kt
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import retrofit2.http.*

interface AionMCPApiService {
    @GET("api/v1/health")
    suspend fun getHealth(): HealthResponse
    
    @GET("api/v1/tools/list")
    suspend fun listTools(): ToolsResponse
    
    @POST("api/v1/tools/invoke")
    suspend fun invokeTool(
        @Body request: ToolInvokeRequest
    ): ToolInvokeResponse
    
    @GET("api/v1/learning/stats")
    suspend fun getLearningStats(): LearningStatsResponse
}

// Data models
data class HealthResponse(
    val status: String,
    val timestamp: Long,
    val version: String
)

data class ToolsResponse(
    val protocol: String,
    val tools: List<Tool>
)

data class Tool(
    val name: String,
    val description: String,
    val metadata: Map<String, Any>
)

data class ToolInvokeRequest(
    val tool: String,
    val args: Map<String, Any>
)

data class ToolInvokeResponse(
    val tool: String,
    val result: Any
)

data class LearningStatsResponse(
    val total_executions: Int,
    val success_rate: Double,
    val avg_latency_ms: Double
)

// Client builder
class AionMCPClient(baseUrl: String) {
    private val retrofit = Retrofit.Builder()
        .baseUrl(baseUrl)
        .addConverterFactory(GsonConverterFactory.create())
        .build()
    
    val api: AionMCPApiService = retrofit.create(AionMCPApiService::class.java)
}

// Usage example
suspend fun example() {
    val client = AionMCPClient("https://your-aionmcp-server.com/")
    
    // Check server health
    val health = client.api.getHealth()
    println("Server status: ${health.status}")
    
    // List available tools
    val tools = client.api.listTools()
    tools.tools.forEach { tool ->
        println("Tool: ${tool.name} - ${tool.description}")
    }
    
    // Invoke a tool
    val request = ToolInvokeRequest(
        tool = "openapi.petstore.listPets",
        args = mapOf("limit" to 10)
    )
    val result = client.api.invokeTool(request)
    println("Result: ${result.result}")
}
```

### gRPC Client (Kotlin)

```kotlin
// AionMCPGrpcClient.kt
import io.grpc.ManagedChannel
import io.grpc.ManagedChannelBuilder
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext

class AionMCPGrpcClient(host: String, port: Int) {
    private val channel: ManagedChannel = ManagedChannelBuilder
        .forAddress(host, port)
        .usePlaintext() // Use .useTransportSecurity() for TLS
        .build()
    
    private val stub = AgentServiceGrpc.newBlockingStub(channel)
    
    suspend fun registerAgent(agentId: String, name: String): RegisterAgentResponse {
        return withContext(Dispatchers.IO) {
            val request = RegisterAgentRequest.newBuilder()
                .setAgentId(agentId)
                .setName(name)
                .build()
            stub.registerAgent(request)
        }
    }
    
    suspend fun listTools(): ListToolsResponse {
        return withContext(Dispatchers.IO) {
            val request = ListToolsRequest.newBuilder().build()
            stub.listTools(request)
        }
    }
    
    suspend fun invokeTool(toolName: String, args: Map<String, Any>): InvokeToolResponse {
        return withContext(Dispatchers.IO) {
            val request = InvokeToolRequest.newBuilder()
                .setToolName(toolName)
                .putAllArguments(args.mapValues { it.value.toString() })
                .build()
            stub.invokeTool(request)
        }
    }
    
    fun shutdown() {
        channel.shutdown()
    }
}

// Usage example
suspend fun grpcExample() {
    val client = AionMCPGrpcClient("your-server.com", 50051)
    
    try {
        // Register agent
        val registration = client.registerAgent("mobile-app-001", "Android Client")
        println("Session ID: ${registration.sessionId}")
        
        // List tools
        val tools = client.listTools()
        tools.toolsList.forEach { tool ->
            println("Tool: ${tool.name}")
        }
        
        // Invoke tool
        val result = client.invokeTool(
            "openapi.petstore.listPets",
            mapOf("limit" to 10)
        )
        println("Result: ${result.result}")
    } finally {
        client.shutdown()
    }
}
```

### Android Best Practices

1. **Network Operations on Background Thread**: Always use coroutines or RxJava
2. **Error Handling**: Implement retry logic with exponential backoff
3. **Caching**: Cache tool lists and metadata locally
4. **Connection Pooling**: Reuse HTTP clients and gRPC channels
5. **Timeouts**: Set appropriate timeouts for mobile networks

```kotlin
// Error handling example
class AionMCPRepository(private val client: AionMCPClient) {
    suspend fun invokeToolWithRetry(
        toolName: String,
        args: Map<String, Any>,
        maxRetries: Int = 3
    ): Result<ToolInvokeResponse> {
        var lastException: Exception? = null
        
        repeat(maxRetries) { attempt ->
            try {
                val request = ToolInvokeRequest(toolName, args)
                val response = client.api.invokeTool(request)
                return Result.success(response)
            } catch (e: Exception) {
                lastException = e
                if (attempt < maxRetries - 1) {
                    delay(1000L * (attempt + 1)) // Exponential backoff
                }
            }
        }
        
        return Result.failure(lastException ?: Exception("Unknown error"))
    }
}
```

## iOS Integration

### Prerequisites

```swift
// Package.swift or Podfile
dependencies: [
    .package(url: "https://github.com/Alamofire/Alamofire.git", from: "5.8.0"),
    .package(url: "https://github.com/grpc/grpc-swift.git", from: "1.20.0")
]

// Or using CocoaPods
pod 'Alamofire', '~> 5.8'
pod 'gRPC-Swift', '~> 1.20'
```

### REST API Client (Swift)

```swift
// AionMCPClient.swift
import Foundation
import Alamofire

struct HealthResponse: Codable {
    let status: String
    let timestamp: Int64
    let version: String
}

struct Tool: Codable {
    let name: String
    let description: String
    let metadata: [String: AnyCodable]
}

struct ToolsResponse: Codable {
    let protocol: String
    let tools: [Tool]
}

struct ToolInvokeRequest: Codable {
    let tool: String
    let args: [String: AnyCodable]
}

struct ToolInvokeResponse: Codable {
    let tool: String
    let result: AnyCodable
}

struct LearningStatsResponse: Codable {
    let totalExecutions: Int
    let successRate: Double
    let avgLatencyMs: Double
    
    enum CodingKeys: String, CodingKey {
        case totalExecutions = "total_executions"
        case successRate = "success_rate"
        case avgLatencyMs = "avg_latency_ms"
    }
}

class AionMCPClient {
    private let baseURL: String
    private let session: Session
    
    init(baseURL: String) {
        self.baseURL = baseURL
        
        let configuration = URLSessionConfiguration.default
        configuration.timeoutIntervalForRequest = 30
        configuration.timeoutIntervalForResource = 60
        
        self.session = Session(configuration: configuration)
    }
    
    // Health check
    func getHealth() async throws -> HealthResponse {
        return try await session.request("\(baseURL)/api/v1/health")
            .validate()
            .serializingDecodable(HealthResponse.self)
            .value
    }
    
    // List tools
    func listTools() async throws -> ToolsResponse {
        return try await session.request("\(baseURL)/api/v1/tools/list")
            .validate()
            .serializingDecodable(ToolsResponse.self)
            .value
    }
    
    // Invoke tool
    func invokeTool(name: String, args: [String: Any]) async throws -> ToolInvokeResponse {
        let request = ToolInvokeRequest(
            tool: name,
            args: args.mapValues { AnyCodable($0) }
        )
        
        return try await session.request(
            "\(baseURL)/api/v1/tools/invoke",
            method: .post,
            parameters: request,
            encoder: JSONParameterEncoder.default
        )
        .validate()
        .serializingDecodable(ToolInvokeResponse.self)
        .value
    }
    
    // Get learning stats
    func getLearningStats() async throws -> LearningStatsResponse {
        return try await session.request("\(baseURL)/api/v1/learning/stats")
            .validate()
            .serializingDecodable(LearningStatsResponse.self)
            .value
    }
}

// Helper for encoding/decoding Any type
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
        case let double as Double:
            try container.encode(double)
        case let string as String:
            try container.encode(string)
        case let bool as Bool:
            try container.encode(bool)
        case let array as [Any]:
            try container.encode(array.map { AnyCodable($0) })
        case let dict as [String: Any]:
            try container.encode(dict.mapValues { AnyCodable($0) })
        default:
            try container.encodeNil()
        }
    }
}

// Usage example
func example() async {
    let client = AionMCPClient(baseURL: "https://your-aionmcp-server.com")
    
    do {
        // Check server health
        let health = try await client.getHealth()
        print("Server status: \(health.status)")
        
        // List available tools
        let tools = try await client.listTools()
        for tool in tools.tools {
            print("Tool: \(tool.name) - \(tool.description)")
        }
        
        // Invoke a tool
        let result = try await client.invokeTool(
            name: "openapi.petstore.listPets",
            args: ["limit": 10]
        )
        print("Result: \(result.result)")
    } catch {
        print("Error: \(error)")
    }
}
```

### gRPC Client (Swift)

```swift
// AionMCPGrpcClient.swift
import GRPC
import NIO

class AionMCPGrpcClient {
    private let group: EventLoopGroup
    private let channel: GRPCChannel
    private let client: AgentService_AgentServiceNIOClient
    
    init(host: String, port: Int) {
        self.group = PlatformSupport.makeEventLoopGroup(loopCount: 1)
        
        self.channel = try! GRPCChannelPool.with(
            target: .host(host, port: port),
            transportSecurity: .plaintext, // Use .tls() for secure connections
            eventLoopGroup: group
        )
        
        self.client = AgentService_AgentServiceNIOClient(channel: channel)
    }
    
    func registerAgent(agentId: String, name: String) async throws -> AgentService_RegisterAgentResponse {
        var request = AgentService_RegisterAgentRequest()
        request.agentID = agentId
        request.name = name
        
        let call = client.registerAgent(request)
        return try await call.response.get()
    }
    
    func listTools() async throws -> AgentService_ListToolsResponse {
        let request = AgentService_ListToolsRequest()
        let call = client.listTools(request)
        return try await call.response.get()
    }
    
    func invokeTool(toolName: String, args: [String: String]) async throws -> AgentService_InvokeToolResponse {
        var request = AgentService_InvokeToolRequest()
        request.toolName = toolName
        request.arguments = args
        
        let call = client.invokeTool(request)
        return try await call.response.get()
    }
    
    func shutdown() throws {
        try channel.close().wait()
        try group.syncShutdownGracefully()
    }
    
    deinit {
        try? shutdown()
    }
}

// Usage example
func grpcExample() async {
    let client = AionMCPGrpcClient(host: "your-server.com", port: 50051)
    
    do {
        // Register agent
        let registration = try await client.registerAgent(
            agentId: "ios-app-001",
            name: "iOS Client"
        )
        print("Session ID: \(registration.sessionID)")
        
        // List tools
        let tools = try await client.listTools()
        for tool in tools.tools {
            print("Tool: \(tool.name)")
        }
        
        // Invoke tool
        let result = try await client.invokeTool(
            toolName: "openapi.petstore.listPets",
            args: ["limit": "10"]
        )
        print("Result: \(result.result)")
    } catch {
        print("Error: \(error)")
    }
}
```

### iOS Best Practices

1. **Use async/await**: Modern Swift concurrency for cleaner code
2. **App Transport Security**: Configure ATS for your server domain
3. **Background Tasks**: Handle network operations during app backgrounding
4. **Memory Management**: Properly manage gRPC channels and connections
5. **Error Handling**: Implement comprehensive error handling with user feedback

```swift
// Info.plist configuration for ATS
<key>NSAppTransportSecurity</key>
<dict>
    <key>NSExceptionDomains</key>
    <dict>
        <key>your-aionmcp-server.com</key>
        <dict>
            <key>NSExceptionAllowsInsecureHTTPLoads</key>
            <false/>
            <key>NSExceptionRequiresForwardSecrecy</key>
            <true/>
            <key>NSIncludesSubdomains</key>
            <true/>
        </dict>
    </dict>
</dict>

// Error handling with retry
actor AionMCPRepository {
    private let client: AionMCPClient
    
    init(client: AionMCPClient) {
        self.client = client
    }
    
    func invokeToolWithRetry(
        name: String,
        args: [String: Any],
        maxRetries: Int = 3
    ) async throws -> ToolInvokeResponse {
        var lastError: Error?
        
        for attempt in 0..<maxRetries {
            do {
                return try await client.invokeTool(name: name, args: args)
            } catch {
                lastError = error
                if attempt < maxRetries - 1 {
                    try await Task.sleep(nanoseconds: UInt64(1_000_000_000 * (attempt + 1)))
                }
            }
        }
        
        throw lastError ?? NSError(domain: "AionMCP", code: -1)
    }
}
```

## Security Considerations

### Authentication

AionMCP supports API key-based authentication. Include the API key in request headers:

**Android:**
```kotlin
val interceptor = Interceptor { chain ->
    val request = chain.request().newBuilder()
        .addHeader("X-API-Key", "your-api-key")
        .build()
    chain.proceed(request)
}

val okHttpClient = OkHttpClient.Builder()
    .addInterceptor(interceptor)
    .build()

val retrofit = Retrofit.Builder()
    .client(okHttpClient)
    .build()
```

**iOS:**
```swift
let interceptor = RequestInterceptor()

class RequestInterceptor: Alamofire.RequestInterceptor {
    func adapt(_ urlRequest: URLRequest, for session: Session, completion: @escaping (Result<URLRequest, Error>) -> Void) {
        var urlRequest = urlRequest
        urlRequest.headers.add(name: "X-API-Key", value: "your-api-key")
        completion(.success(urlRequest))
    }
}

let session = Session(interceptor: interceptor)
```

### TLS/SSL

Always use HTTPS in production:

1. **Certificate Pinning**: Pin your server's certificate for enhanced security
2. **Self-Signed Certificates**: Configure trust for development environments
3. **Transport Security**: Enable strict transport security on server

### Data Privacy

1. **Sensitive Data**: Never log sensitive information
2. **Secure Storage**: Use Keychain (iOS) or EncryptedSharedPreferences (Android)
3. **GDPR Compliance**: Handle user data according to regulations

## Network Optimization

### Connection Management

1. **Connection Pooling**: Reuse connections when possible
2. **Keep-Alive**: Enable HTTP keep-alive
3. **Compression**: Enable gzip compression

### Offline Support

```kotlin
// Android offline cache
val cacheSize = 10 * 1024 * 1024L // 10 MB
val cache = Cache(context.cacheDir, cacheSize)

val okHttpClient = OkHttpClient.Builder()
    .cache(cache)
    .addNetworkInterceptor(CacheInterceptor())
    .build()
```

```swift
// iOS offline cache
let configuration = URLSessionConfiguration.default
configuration.requestCachePolicy = .returnCacheDataElseLoad
configuration.urlCache = URLCache(
    memoryCapacity: 10 * 1024 * 1024,
    diskCapacity: 50 * 1024 * 1024
)
```

### Request Batching

Combine multiple tool invocations when possible to reduce network overhead.

## Testing

### Mock Server for Development

Use mock servers during development:

**Android:**
```kotlin
class MockAionMCPApi : AionMCPApiService {
    override suspend fun getHealth() = HealthResponse("healthy", System.currentTimeMillis(), "0.1.0")
    override suspend fun listTools() = ToolsResponse("MCP/1.0", listOf())
    // ... implement other methods
}
```

**iOS:**
```swift
class MockAionMCPClient: AionMCPClient {
    override func getHealth() async throws -> HealthResponse {
        return HealthResponse(status: "healthy", timestamp: Int64(Date().timeIntervalSince1970), version: "0.1.0")
    }
    // ... implement other methods
}
```

## Deployment Considerations

### Server Configuration

1. **CORS**: Enable CORS for web-based mobile frameworks
2. **Rate Limiting**: Implement per-client rate limits
3. **Load Balancing**: Use load balancers for scalability
4. **Monitoring**: Monitor mobile client connections

### Example CORS Configuration

Add to your AionMCP server configuration:

```yaml
server:
  port: 8080
  cors:
    enabled: true
    allowed_origins:
      - "https://your-mobile-app.com"
    allowed_methods:
      - GET
      - POST
      - PUT
      - DELETE
    allowed_headers:
      - Content-Type
      - Authorization
      - X-API-Key
```

## Troubleshooting

### Common Issues

1. **Connection Timeout**: Increase timeout values for slow networks
2. **SSL Errors**: Verify certificate configuration
3. **JSON Parsing**: Ensure data models match server responses
4. **Authentication Failures**: Verify API key and headers

### Debugging

Enable network logging:

**Android:**
```kotlin
val logging = HttpLoggingInterceptor().apply {
    level = HttpLoggingInterceptor.Level.BODY
}

val okHttpClient = OkHttpClient.Builder()
    .addInterceptor(logging)
    .build()
```

**iOS:**
```swift
let monitor = EventMonitor()

class EventMonitor: Alamofire.EventMonitor {
    func requestDidFinish(_ request: Request) {
        print("Request: \(request.description)")
    }
}

let session = Session(eventMonitors: [monitor])
```

## Next Steps

1. Review the [API Documentation](./architecture.md)
2. Explore [Example Specifications](../examples/specs/)
3. Check [Self-Learning Features](./self_learning_usage.md)
4. Join our community for support

## Resources

- [AionMCP GitHub Repository](https://github.com/kiransth77/aionmcp)
- [Protocol Buffers Documentation](https://protobuf.dev/)
- [gRPC Mobile Documentation](https://grpc.io/docs/platforms/android/)
- [Retrofit Documentation](https://square.github.io/retrofit/)
- [Alamofire Documentation](https://github.com/Alamofire/Alamofire)
