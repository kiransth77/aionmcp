# iOS Integration Example

This directory contains example code for integrating AionMCP into iOS applications.

## Setup

### 1. Add Dependencies

Using Swift Package Manager, add these dependencies to your `Package.swift`:

```swift
dependencies: [
    .package(url: "https://github.com/Alamofire/Alamofire.git", from: "5.8.0"),
    .package(url: "https://github.com/grpc/grpc-swift.git", from: "1.20.0")
]
```

Or using CocoaPods, add to your `Podfile`:

```ruby
platform :ios, '14.0'
use_frameworks!

target 'YourApp' do
  # Networking
  pod 'Alamofire', '~> 5.8'
  
  # Optional: gRPC for high-performance
  pod 'gRPC-Swift', '~> 1.20'
end
```

### 2. Configure App Transport Security

Add to your `Info.plist`:

```xml
<key>NSAppTransportSecurity</key>
<dict>
    <!-- For production with HTTPS -->
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
    
    <!-- For development (remove in production) -->
    <key>NSAllowsLocalNetworking</key>
    <true/>
</dict>
```

### 3. Add Background Modes (Optional)

For background network operations:

```xml
<key>UIBackgroundModes</key>
<array>
    <string>fetch</string>
    <string>processing</string>
</array>
```

## Project Structure

```
YourApp/
├── Models/              # Data models
├── Networking/          # API client
├── Repositories/        # Repository pattern
├── ViewModels/          # ViewModels
└── Views/               # SwiftUI views
```

## Implementation Files

### 1. Data Models (`Models/AionMCPModels.swift`)

See [AionMCPModels.swift](./AionMCPModels.swift)

### 2. API Client (`Networking/AionMCPClient.swift`)

See [AionMCPClient.swift](./AionMCPClient.swift)

### 3. Repository (`Repositories/AionMCPRepository.swift`)

See [AionMCPRepository.swift](./AionMCPRepository.swift)

### 4. ViewModel (`ViewModels/AionMCPViewModel.swift`)

See [AionMCPViewModel.swift](./AionMCPViewModel.swift)

## Usage Example

### In Your SwiftUI View

```swift
import SwiftUI

struct ContentView: View {
    @StateObject private var viewModel = AionMCPViewModel(
        repository: AionMCPRepository(
            client: AionMCPClient(baseURL: "https://your-aionmcp-server.com")
        )
    )
    
    var body: some View {
        NavigationView {
            List {
                Section("Server Status") {
                    if let health = viewModel.health {
                        HStack {
                            Text("Status")
                            Spacer()
                            Text(health.status)
                                .foregroundColor(.green)
                        }
                        
                        HStack {
                            Text("Version")
                            Spacer()
                            Text(health.version)
                        }
                    } else if viewModel.isLoading {
                        ProgressView()
                    }
                }
                
                Section("Available Tools") {
                    ForEach(viewModel.tools, id: \.name) { tool in
                        VStack(alignment: .leading) {
                            Text(tool.name)
                                .font(.headline)
                            Text(tool.description)
                                .font(.caption)
                                .foregroundColor(.secondary)
                        }
                    }
                }
            }
            .navigationTitle("AionMCP")
            .refreshable {
                await viewModel.refresh()
            }
            .task {
                await viewModel.loadInitialData()
            }
        }
    }
}
```

### Invoke a Tool

```swift
struct ToolInvocationView: View {
    @StateObject private var viewModel: AionMCPViewModel
    let toolName: String
    
    @State private var limit: String = "10"
    
    var body: some View {
        Form {
            Section("Parameters") {
                TextField("Limit", text: $limit)
                    .keyboardType(.numberPad)
            }
            
            Section {
                Button("Invoke Tool") {
                    Task {
                        await invokeTool()
                    }
                }
            }
            
            if let result = viewModel.toolResult {
                Section("Result") {
                    Text(String(describing: result))
                        .font(.system(.body, design: .monospaced))
                }
            }
        }
        .navigationTitle(toolName)
    }
    
    private func invokeTool() async {
        let args: [String: Any] = ["limit": Int(limit) ?? 10]
        await viewModel.invokeTool(name: toolName, args: args)
    }
}
```

## Advanced Features

### Caching

Implement caching with UserDefaults or Core Data:

```swift
class AionMCPCache {
    private let defaults = UserDefaults.standard
    private let cacheKey = "aionmcp_tools_cache"
    private let timestampKey = "aionmcp_cache_timestamp"
    private let cacheValiditySeconds: TimeInterval = 300 // 5 minutes
    
    func cacheTools(_ tools: [Tool]) {
        if let encoded = try? JSONEncoder().encode(tools) {
            defaults.set(encoded, forKey: cacheKey)
            defaults.set(Date(), forKey: timestampKey)
        }
    }
    
    func getCachedTools() -> [Tool]? {
        guard let timestamp = defaults.object(forKey: timestampKey) as? Date,
              Date().timeIntervalSince(timestamp) < cacheValiditySeconds,
              let data = defaults.data(forKey: cacheKey),
              let tools = try? JSONDecoder().decode([Tool].self, from: data) else {
            return nil
        }
        return tools
    }
    
    func clearCache() {
        defaults.removeObject(forKey: cacheKey)
        defaults.removeObject(forKey: timestampKey)
    }
}
```

### Offline Support with Core Data

```swift
import CoreData

@MainActor
class OfflineRepository: ObservableObject {
    private let viewContext: NSManagedObjectContext
    private let client: AionMCPClient
    
    init(viewContext: NSManagedObjectContext, client: AionMCPClient) {
        self.viewContext = viewContext
        self.client = client
    }
    
    func fetchTools() async throws -> [Tool] {
        // Try to fetch from network
        do {
            let tools = try await client.listTools().tools
            await saveToCache(tools)
            return tools
        } catch {
            // Fall back to cached data
            return try fetchFromCache()
        }
    }
    
    private func saveToCache(_ tools: [Tool]) async {
        // Save to Core Data
        for tool in tools {
            let entity = ToolEntity(context: viewContext)
            entity.name = tool.name
            entity.toolDescription = tool.description
            entity.lastUpdated = Date()
        }
        
        try? viewContext.save()
    }
    
    private func fetchFromCache() throws -> [Tool] {
        let request = ToolEntity.fetchRequest()
        let entities = try viewContext.fetch(request)
        return entities.map { Tool(name: $0.name ?? "", description: $0.toolDescription ?? "") }
    }
}
```

### Background Refresh

```swift
import BackgroundTasks

class BackgroundSyncManager {
    static let shared = BackgroundSyncManager()
    private let taskIdentifier = "com.example.aionmcp.refresh"
    
    func registerBackgroundTasks() {
        BGTaskScheduler.shared.register(
            forTaskWithIdentifier: taskIdentifier,
            using: nil
        ) { task in
            self.handleBackgroundRefresh(task: task as! BGAppRefreshTask)
        }
    }
    
    func scheduleBackgroundRefresh() {
        let request = BGAppRefreshTaskRequest(identifier: taskIdentifier)
        request.earliestBeginDate = Date(timeIntervalSinceNow: 15 * 60) // 15 minutes
        
        try? BGTaskScheduler.shared.submit(request)
    }
    
    private func handleBackgroundRefresh(task: BGAppRefreshTask) {
        let repository = AionMCPRepository(client: AionMCPClient(baseURL: "..."))
        
        Task {
            do {
                _ = try await repository.listTools(forceRefresh: true)
                task.setTaskCompleted(success: true)
            } catch {
                task.setTaskCompleted(success: false)
            }
            
            scheduleBackgroundRefresh()
        }
    }
}

// In AppDelegate or App struct
func application(_ application: UIApplication, didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?) -> Bool {
    BackgroundSyncManager.shared.registerBackgroundTasks()
    return true
}

func applicationDidEnterBackground(_ application: UIApplication) {
    BackgroundSyncManager.shared.scheduleBackgroundRefresh()
}
```

### Combine Integration

```swift
import Combine

class AionMCPCombineClient {
    private let client: AionMCPClient
    
    init(client: AionMCPClient) {
        self.client = client
    }
    
    func healthPublisher() -> AnyPublisher<HealthResponse, Error> {
        Future { promise in
            Task {
                do {
                    let health = try await self.client.getHealth()
                    promise(.success(health))
                } catch {
                    promise(.failure(error))
                }
            }
        }
        .eraseToAnyPublisher()
    }
    
    func toolsPublisher() -> AnyPublisher<[Tool], Error> {
        Future { promise in
            Task {
                do {
                    let response = try await self.client.listTools()
                    promise(.success(response.tools))
                } catch {
                    promise(.failure(error))
                }
            }
        }
        .eraseToAnyPublisher()
    }
}
```

## Testing

### Unit Tests

```swift
import XCTest
@testable import YourApp

class AionMCPRepositoryTests: XCTestCase {
    var mockClient: MockAionMCPClient!
    var repository: AionMCPRepository!
    
    override func setUp() {
        super.setUp()
        mockClient = MockAionMCPClient()
        repository = AionMCPRepository(client: mockClient)
    }
    
    func testListTools() async throws {
        // Given
        let expectedTools = [
            Tool(name: "tool1", description: "Description 1")
        ]
        mockClient.toolsToReturn = expectedTools
        
        // When
        let tools = try await repository.listTools()
        
        // Then
        XCTAssertEqual(tools.count, 1)
        XCTAssertEqual(tools.first?.name, "tool1")
    }
}

class MockAionMCPClient: AionMCPClient {
    var toolsToReturn: [Tool] = []
    
    override func listTools() async throws -> ToolsResponse {
        return ToolsResponse(protocol: "MCP/1.0", tools: toolsToReturn)
    }
}
```

### UI Tests

```swift
import XCTest

class AionMCPUITests: XCTestCase {
    var app: XCUIApplication!
    
    override func setUp() {
        super.setUp()
        continueAfterFailure = false
        app = XCUIApplication()
        app.launch()
    }
    
    func testToolsList() {
        // Wait for tools to load
        let toolsList = app.tables.firstMatch
        XCTAssertTrue(toolsList.waitForExistence(timeout: 5))
        
        // Verify tools are displayed
        XCTAssertTrue(toolsList.cells.count > 0)
    }
}
```

## Troubleshooting

### Common Issues

1. **App Transport Security Errors**
   - Solution: Configure ATS properly in Info.plist

2. **Network Timeout**
   - Solution: Increase timeout or check server availability

3. **JSON Decoding Errors**
   - Solution: Verify models match API responses

4. **Background Task Not Running**
   - Solution: Check background modes and task registration

### Enable Logging

```swift
import os.log

extension Logger {
    static let aionmcp = Logger(subsystem: "com.example.aionmcp", category: "networking")
}

// Usage
Logger.aionmcp.info("Fetching tools from server")
Logger.aionmcp.error("Failed to fetch tools: \(error.localizedDescription)")
```

### Network Debugging

```swift
let monitor = EventMonitor()

class EventMonitor: Alamofire.EventMonitor {
    func requestDidFinish(_ request: Request) {
        print("📡 Request: \(request.description)")
    }
    
    func request<Value>(_ request: DataRequest, didParseResponse response: DataResponse<Value, AFError>) {
        print("📥 Response: \(response.debugDescription)")
    }
}

let session = Session(eventMonitors: [monitor])
```

## Next Steps

1. Implement error handling UI
2. Add authentication flow
3. Implement caching strategies
4. Add comprehensive tests
5. Optimize for battery efficiency
6. Add analytics and monitoring

## Resources

- [Alamofire Documentation](https://github.com/Alamofire/Alamofire)
- [Swift Concurrency Guide](https://docs.swift.org/swift-book/LanguageGuide/Concurrency.html)
- [SwiftUI Documentation](https://developer.apple.com/documentation/swiftui)
- [gRPC Swift Tutorial](https://grpc.io/docs/languages/swift/)
- [Core Data Guide](https://developer.apple.com/documentation/coredata)
