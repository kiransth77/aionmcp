# AionMCP Android Demo App

A complete Android application demonstrating AionMCP integration.

## Features

- ✅ List available tools from AionMCP server
- ✅ Execute tools with parameters
- ✅ View execution results
- ✅ Monitor server health
- ✅ View learning statistics
- ✅ Material Design 3 UI
- ✅ Dark mode support
- ✅ Error handling and retry logic

## Screenshots

[Screenshots will be available in releases]

## Requirements

- Android Studio Hedgehog (2023.1.1) or newer
- Android SDK 24+
- Kotlin 1.9+
- AionMCP server running (see [deployment guide](../../../docs/mobile_deployment.md))

## Quick Start

### 1. Download APK (Easiest)

Download the latest APK from [Releases](https://github.com/kiransth77/aionmcp/releases):
- `aionmcp-demo-v1.0.0.apk`

Install on your Android device and configure the server URL in settings.

### 2. Build from Source

```bash
# Clone the repository
git clone https://github.com/kiransth77/aionmcp.git
cd aionmcp/examples/mobile/android-app

# Open in Android Studio
# or build from command line:
./gradlew assembleDebug

# Install on device
./gradlew installDebug
```

## Configuration

### Server URL

1. Open the app
2. Tap the settings icon (⚙️)
3. Enter your AionMCP server URL (e.g., `https://api.example.com`)
4. Optional: Enter API key if authentication is enabled
5. Tap "Save"

### Local Development

For testing with a local server:
- Use `http://10.0.2.2:8080` for Android emulator
- Use `http://<YOUR_IP>:8080` for physical devices

## Project Structure

```
android-app/
├── app/
│   ├── build.gradle.kts          # App dependencies
│   └── src/
│       └── main/
│           ├── java/com/aionmcp/demo/
│           │   ├── MainActivity.kt        # Main entry point
│           │   ├── ui/
│           │   │   ├── ToolsScreen.kt     # Tools list UI
│           │   │   ├── ExecuteScreen.kt   # Tool execution UI
│           │   │   ├── StatsScreen.kt     # Statistics UI
│           │   │   └── SettingsScreen.kt  # Settings UI
│           │   ├── data/
│           │   │   ├── models/            # Data models
│           │   │   ├── api/               # API service
│           │   │   └── repository/        # Repository
│           │   └── viewmodel/             # ViewModels
│           ├── res/                        # Resources
│           └── AndroidManifest.xml
├── build.gradle.kts              # Project build config
└── settings.gradle.kts           # Project settings
```

## Building

### Debug Build

```bash
./gradlew assembleDebug
# Output: app/build/outputs/apk/debug/app-debug.apk
```

### Release Build

```bash
# Create keystore first (one-time setup)
keytool -genkey -v -keystore aionmcp-demo.keystore \
  -alias aionmcp -keyalg RSA -keysize 2048 -validity 10000

# Build signed release
./gradlew assembleRelease
# Output: app/build/outputs/apk/release/app-release.apk
```

## Features in Detail

### Tools List
- Browse all available tools from the server
- Search and filter tools
- View tool descriptions and metadata
- Pull to refresh

### Tool Execution
- Input parameters with type validation
- Execute tools with one tap
- View formatted results (JSON, text)
- Copy results to clipboard
- Execution history

### Statistics
- Real-time server health monitoring
- Success rate metrics
- Average latency
- Total executions
- Tool-specific statistics

### Settings
- Configure server URL
- API key authentication
- Connection timeout settings
- Enable/disable debug logging
- Clear cache

## Troubleshooting

### Cannot Connect to Server

1. Check server URL is correct
2. Ensure server is running and accessible
3. For local development:
   - Emulator: Use `http://10.0.2.2:8080`
   - Physical device: Use your computer's IP
4. Check firewall settings
5. Enable debug logging in settings for more details

### SSL Certificate Errors

For development with self-signed certificates:
1. Add certificate to Android trust store, or
2. Use HTTP for local testing (not recommended for production)

### Tools Not Loading

1. Pull to refresh
2. Check network connection
3. Verify API key if authentication is enabled
4. Check server logs

## Development

### Adding New Features

1. Create UI in `ui/` package
2. Add ViewModel in `viewmodel/` package
3. Wire up in `MainActivity.kt` navigation

### Testing

```bash
# Run unit tests
./gradlew test

# Run instrumented tests
./gradlew connectedAndroidTest
```

### Code Style

Follow Kotlin coding conventions:
```bash
./gradlew ktlintFormat
```

## Technologies Used

- **Kotlin** - Primary language
- **Jetpack Compose** - Modern UI toolkit
- **Material 3** - Design system
- **Retrofit** - REST API client
- **OkHttp** - HTTP client
- **Gson** - JSON parsing
- **Kotlin Coroutines** - Async operations
- **ViewModel** - UI state management
- **DataStore** - Settings persistence

## Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

MIT License - see [LICENSE](../../../LICENSE) file

## Support

- 📖 [AionMCP Documentation](../../../docs/)
- 🐛 [Report Issues](https://github.com/kiransth77/aionmcp/issues)
- 💬 [Discussions](https://github.com/kiransth77/aionmcp/discussions)

## Related

- [iOS Demo App](../ios-app/)
- [Mobile Integration Guide](../../../docs/mobile_integration.md)
- [Mobile Deployment Guide](../../../docs/mobile_deployment.md)
