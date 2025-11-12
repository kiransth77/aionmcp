# AionMCP iOS Demo App

A complete iOS application demonstrating AionMCP integration.

## Features

- ✅ List available tools from AionMCP server
- ✅ Execute tools with parameters
- ✅ View execution results
- ✅ Monitor server health
- ✅ View learning statistics
- ✅ Native iOS design with SwiftUI
- ✅ Dark mode support
- ✅ Error handling and retry logic
- ✅ Offline caching

## Screenshots

[Screenshots will be available in releases]

## Requirements

- Xcode 15.0 or newer
- iOS 16.0+
- Swift 5.9+
- AionMCP server running (see [deployment guide](../../../docs/mobile_deployment.md))

## Quick Start

### 1. Download from App Store (Coming Soon)

The app will be available on the App Store.

### 2. Install via TestFlight (Beta)

**Coming Soon**: TestFlight beta will be available in the future. For now, please build from source.

### 3. Build from Source

```bash
# Clone the repository
git clone https://github.com/kiransth77/aionmcp.git
cd aionmcp/examples/mobile/ios-app

# Open in Xcode
open AionMCPDemo.xcodeproj

# Or use command line:
xcodebuild -project AionMCPDemo.xcodeproj \
           -scheme AionMCPDemo \
           -configuration Debug \
           -destination 'platform=iOS Simulator,name=iPhone 15'
```

## Configuration

### Server URL

1. Open the app
2. Navigate to Settings tab (⚙️)
3. Enter your AionMCP server URL (e.g., `https://api.example.com`)
4. Optional: Enter API key if authentication is enabled
5. Tap "Save"

### Local Development

For testing with a local server:
- Use `http://localhost:8080` for iOS Simulator
- Use `http://<YOUR_IP>:8080` for physical devices

## Project Structure

```
ios-app/
├── AionMCPDemo.xcodeproj      # Xcode project
├── AionMCPDemo/
│   ├── AionMCPDemoApp.swift   # App entry point
│   ├── Views/
│   │   ├── ContentView.swift       # Main navigation
│   │   ├── ToolsView.swift         # Tools list
│   │   ├── ToolDetailView.swift    # Tool execution
│   │   ├── StatsView.swift         # Statistics
│   │   └── SettingsView.swift      # Settings
│   ├── ViewModels/
│   │   ├── ToolsViewModel.swift    # Tools logic
│   │   └── StatsViewModel.swift    # Stats logic
│   ├── Models/
│   │   ├── Tool.swift              # Data models
│   │   ├── HealthResponse.swift
│   │   └── LearningStats.swift
│   ├── Networking/
│   │   ├── AionMCPClient.swift     # API client
│   │   └── AionMCPRepository.swift # Repository
│   ├── Storage/
│   │   └── SettingsManager.swift   # UserDefaults wrapper
│   └── Assets.xcassets             # Images and colors
├── Package.swift                    # Swift Package Manager
└── README.md
```

## Building

### Debug Build

```bash
# Build for simulator
xcodebuild -project AionMCPDemo.xcodeproj \
           -scheme AionMCPDemo \
           -configuration Debug \
           -sdk iphonesimulator

# Run on simulator
xcrun simctl boot "iPhone 15"
xcrun simctl install booted ./build/Debug-iphonesimulator/AionMCPDemo.app
xcrun simctl launch booted com.aionmcp.demo
```

### Release Build

```bash
# Archive for distribution
xcodebuild -project AionMCPDemo.xcodeproj \
           -scheme AionMCPDemo \
           -configuration Release \
           -archivePath ./build/AionMCPDemo.xcarchive \
           archive

# Export IPA
xcodebuild -exportArchive \
           -archivePath ./build/AionMCPDemo.xcarchive \
           -exportPath ./build \
           -exportOptionsPlist ExportOptions.plist
```

## Features in Detail

### Tools List
- Browse all available tools from the server
- Search and filter tools
- View tool descriptions and metadata
- Pull to refresh
- Offline caching

### Tool Execution
- Input parameters with validation
- Execute tools with one tap
- View formatted results (JSON, text)
- Share results
- Execution history

### Statistics
- Real-time server health monitoring
- Success rate metrics
- Average latency
- Total executions
- Tool-specific statistics
- Animated charts

### Settings
- Configure server URL
- API key authentication
- Connection timeout settings
- Clear cache
- About information

## Troubleshooting

### Cannot Connect to Server

1. Check server URL is correct
2. Ensure server is running and accessible
3. For local development, use your computer's IP address
4. Check App Transport Security settings
5. Enable debug logging for more details

### App Transport Security Errors

Add to `Info.plist` for development:
```xml
<key>NSAppTransportSecurity</key>
<dict>
    <key>NSAllowsLocalNetworking</key>
    <true/>
</dict>
```

### Tools Not Loading

1. Pull to refresh
2. Check network connection
3. Verify API key if authentication is enabled
4. Check settings

## Development

### Adding New Features

1. Create SwiftUI View in `Views/` folder
2. Add ViewModel in `ViewModels/` folder
3. Wire up in `ContentView.swift` navigation

### Testing

```bash
# Run unit tests
xcodebuild test -project AionMCPDemo.xcodeproj \
                -scheme AionMCPDemo \
                -destination 'platform=iOS Simulator,name=iPhone 15'

# Run UI tests
xcodebuild test -project AionMCPDemo.xcodeproj \
                -scheme AionMCPDemoUITests \
                -destination 'platform=iOS Simulator,name=iPhone 15'
```

### Code Style

Use SwiftLint:
```bash
swiftlint lint
swiftlint autocorrect
```

## Technologies Used

- **Swift** - Primary language
- **SwiftUI** - Modern UI framework
- **Combine** - Reactive programming
- **Alamofire** - REST API client
- **UserDefaults** - Settings persistence
- **Async/Await** - Concurrency
- **Charts** - Data visualization (iOS 16+)

## Dependencies

Managed via Swift Package Manager:
- Alamofire 5.8+

Add in Xcode:
1. File → Add Package Dependencies
2. Enter package URL: `https://github.com/Alamofire/Alamofire.git`

## App Store Submission

### Required Assets

- App icons (all sizes)
- Launch screen
- Screenshots for all device sizes
- Privacy policy
- App description

### Privacy

The app requires:
- Network access (to communicate with AionMCP server)
- No personal data collection

Add to `Info.plist`:
```xml
<key>NSAppTransportSecurity</key>
<dict>
    <key>NSAllowsArbitraryLoads</key>
    <false/>
</dict>
```

## Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test on simulator and device
5. Submit a pull request

## License

MIT License - see [LICENSE](../../../LICENSE) file

## Support

- 📖 [AionMCP Documentation](../../../docs/)
- 🐛 [Report Issues](https://github.com/kiransth77/aionmcp/issues)
- 💬 [Discussions](https://github.com/kiransth77/aionmcp/discussions)

## Related

- [Android Demo App](../android-app/)
- [Mobile Integration Guide](../../../docs/mobile_integration.md)
- [Mobile Deployment Guide](../../../docs/mobile_deployment.md)
