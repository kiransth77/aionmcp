# AionMCP Demo Applications

Complete, downloadable mobile applications demonstrating AionMCP integration.

## 📱 Available Apps

### Android Demo App
- **Location**: `android-app/`
- **Platform**: Android 7.0+ (API 24+)
- **Technology**: Kotlin, Jetpack Compose, Retrofit
- **Download**: [APK Releases](https://github.com/kiransth77/aionmcp/releases)
- **Source**: Full project with Gradle build files

### iOS Demo App
- **Location**: `ios-app/`
- **Platform**: iOS 16.0+
- **Technology**: Swift, SwiftUI, Alamofire
- **Download**: TestFlight Beta (Coming Soon) | App Store (Coming Soon)
- **Source**: Full Xcode project

## 🚀 Quick Start

### For End Users (Testing)

#### Android
1. Download APK from [Releases](https://github.com/kiransth77/aionmcp/releases)
2. Install on your Android device
3. Open app and configure server URL in Settings
4. Start exploring AionMCP tools!

#### iOS
**Note:** iOS app is currently available as source code only. TestFlight beta coming soon.

To build from source:
```bash
cd ios-app
open AionMCPDemo.xcodeproj
# Build and run in Xcode (Cmd+R)
```

### For Developers (Building from Source)

#### Android
```bash
cd android-app
./gradlew assembleDebug
./gradlew installDebug
```

#### iOS
```bash
cd ios-app
open AionMCPDemo.xcodeproj
# Build and run in Xcode (Cmd+R)
```

## 📦 Pre-built Binaries

### Download from Releases

Visit [Releases Page](https://github.com/kiransth77/aionmcp/releases) to download:

- **Android**: `aionmcp-demo-android-v1.0.0.apk` (~8 MB)
- **iOS**: TestFlight Beta (Coming Soon) | App Store (Coming Soon)

### Release Schedule

- **v1.0.0** (Current) - Initial release with core features
- **v1.1.0** (Planned) - Advanced features, offline mode improvements
- **v2.0.0** (Future) - gRPC support, background sync

## ✨ Features

Both apps include:

### Core Functionality
- ✅ Browse available tools from AionMCP server
- ✅ Execute tools with parameter input
- ✅ View formatted execution results
- ✅ Real-time server health monitoring
- ✅ Learning statistics and insights
- ✅ Execution history

### User Experience
- ✅ Modern UI (Material 3 for Android, native iOS design)
- ✅ Dark mode support
- ✅ Pull to refresh
- ✅ Search and filter
- ✅ Error handling with retry
- ✅ Offline caching
- ✅ Settings persistence

### Technical Features
- ✅ REST API integration
- ✅ Authentication support (API keys)
- ✅ Automatic retries with exponential backoff
- ✅ Response caching
- ✅ Debug logging

## 🔧 Configuration

### Server Setup

Before using the apps, ensure your AionMCP server is running and accessible:

```bash
# Start AionMCP server
./bin/aionmcp --config config.yaml

# Or using Docker
docker run -p 8080:8080 aionmcp:latest
```

See [Mobile Deployment Guide](../../docs/mobile_deployment.md) for detailed server configuration.

### App Configuration

#### First Time Setup
1. Open the app
2. Navigate to Settings
3. Enter server URL:
   - Production: `https://api.yourdomain.com`
   - Local (Android Emulator): `http://10.0.2.2:8080`
   - Local (iOS Simulator): `http://localhost:8080`
   - Local (Physical Device): `http://<YOUR_IP>:8080`
4. Optional: Enter API key if authentication is enabled
5. Tap Save

#### Advanced Settings
- **Connection Timeout**: Adjust for slow networks (default: 30s)
- **Retry Attempts**: Number of automatic retries (default: 3)
- **Cache Duration**: How long to cache data (default: 5 min)
- **Debug Logging**: Enable for troubleshooting

## 📸 Screenshots

### Android App
[Screenshots to be added in releases]
- Home screen with tools list
- Tool execution screen
- Statistics dashboard
- Settings screen

### iOS App
[Screenshots to be added in releases]
- Tools list view
- Tool detail and execution
- Statistics with charts
- Settings interface

## 🏗️ Building from Source

### Prerequisites

**Android:**
- Android Studio Hedgehog (2023.1.1) or newer
- Android SDK 24+
- Gradle 8.0+

**iOS:**
- Xcode 15.0 or newer
- macOS 13.0+ (Ventura)
- CocoaPods or Swift Package Manager

### Build Instructions

#### Android (Command Line)

```bash
cd examples/mobile/android-app

# Debug build
./gradlew assembleDebug
# Output: app/build/outputs/apk/debug/app-debug.apk

# Release build (requires signing)
./gradlew assembleRelease
# Output: app/build/outputs/apk/release/app-release.apk
```

#### Android (Android Studio)

1. Open Android Studio
2. File → Open → Select `android-app` folder
3. Wait for Gradle sync
4. Run → Run 'app' (Shift+F10)

#### iOS (Command Line)

```bash
cd examples/mobile/ios-app

# Build
xcodebuild -project AionMCPDemo.xcodeproj \
           -scheme AionMCPDemo \
           -configuration Debug \
           -destination 'platform=iOS Simulator,name=iPhone 15'

# Run tests
xcodebuild test -project AionMCPDemo.xcodeproj \
                -scheme AionMCPDemo \
                -destination 'platform=iOS Simulator,name=iPhone 15'
```

#### iOS (Xcode)

1. Open Xcode
2. File → Open → Select `ios-app/AionMCPDemo.xcodeproj`
3. Select target device/simulator
4. Product → Run (Cmd+R)

## 🧪 Testing

### Test Accounts

For testing, you can use:
- **Server**: Demo server at `https://demo.aionmcp.dev` (when available)
- **API Key**: Use demo key from documentation

### Sample Data

The apps include sample data for offline testing:
- Mock tool definitions
- Example execution results
- Cached statistics

### Test Scenarios

1. **Basic Connectivity**
   - Configure server URL
   - Check health status
   - List available tools

2. **Tool Execution**
   - Select a tool (e.g., "openapi.petstore.listPets")
   - Enter parameters
   - Execute and view results

3. **Error Handling**
   - Test with invalid server URL
   - Test with incorrect API key
   - Test network timeout scenarios

4. **Offline Mode**
   - Load tools while online
   - Disconnect network
   - Verify cached data is accessible

## 📚 Documentation

- [Android App README](android-app/README.md) - Detailed Android documentation
- [iOS App README](ios-app/README.md) - Detailed iOS documentation
- [Mobile Integration Guide](../../docs/mobile_integration.md) - API integration details
- [Mobile Deployment Guide](../../docs/mobile_deployment.md) - Server deployment

## 🤝 Contributing

We welcome contributions to improve the demo apps!

### Development Workflow

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly on both platforms (if applicable)
5. Submit a pull request

### Areas for Contribution

- **UI/UX improvements**: Better layouts, animations, user experience
- **Features**: Additional functionality, settings, integrations
- **Bug fixes**: Report and fix issues
- **Tests**: Unit tests, UI tests, integration tests
- **Documentation**: Improve guides, add examples
- **Localization**: Translate to other languages

### Coding Standards

**Android:**
- Follow Kotlin coding conventions
- Use Jetpack Compose best practices
- Run `ktlint` before committing

**iOS:**
- Follow Swift API Design Guidelines
- Use SwiftUI best practices
- Run `swiftlint` before committing

## 🐛 Troubleshooting

### Common Issues

#### "Cannot connect to server"
- Verify server is running
- Check server URL is correct
- Ensure network connectivity
- For local dev, use correct IP address
- Check firewall settings

#### "SSL Certificate error"
- For development with self-signed certs, configure trust settings
- Use HTTP for local testing (not recommended for production)

#### "Tools not loading"
- Pull to refresh
- Check API key if authentication is enabled
- Verify server has tools registered
- Check server logs

#### Android-specific: "Clear text traffic not permitted"
- Update `network_security_config.xml`
- Add domain to allowed cleartext traffic (dev only)

#### iOS-specific: "App Transport Security blocked"
- Update `Info.plist` ATS settings
- Enable local networking for dev

### Debug Logging

Enable debug logging in Settings to see detailed logs:
- Network requests/responses
- API calls
- Error messages
- Cache operations

## 📋 Release Checklist

Before creating a new release:

- [ ] Update version numbers
- [ ] Test on physical devices (Android & iOS)
- [ ] Test on various screen sizes
- [ ] Verify all features work
- [ ] Update CHANGELOG.md
- [ ] Create release notes
- [ ] Build signed APK (Android)
- [ ] Archive and export IPA (iOS)
- [ ] Upload to GitHub Releases
- [ ] Update TestFlight (iOS)
- [ ] Update documentation

## 📞 Support

Having issues? Need help?

- 📖 [Read the Documentation](../../docs/)
- 🐛 [Report a Bug](https://github.com/kiransth77/aionmcp/issues/new?template=bug_report.md)
- 💡 [Request a Feature](https://github.com/kiransth77/aionmcp/issues/new?template=feature_request.md)
- 💬 [Join Discussions](https://github.com/kiransth77/aionmcp/discussions)

## 📄 License

Both demo apps are licensed under the MIT License. See [LICENSE](../../LICENSE) for details.

## 🙏 Acknowledgments

Built with:
- Android: Kotlin, Jetpack Compose, Retrofit, Material 3
- iOS: Swift, SwiftUI, Alamofire, iOS native frameworks

## 🔗 Related Resources

- [AionMCP Main Repository](https://github.com/kiransth77/aionmcp)
- [API Documentation](../../docs/architecture.md)
- [Integration Examples](../)
- [Server Source Code](../../cmd/server/)
