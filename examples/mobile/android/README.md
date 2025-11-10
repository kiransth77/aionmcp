# Android Integration Example

This directory contains example code for integrating AionMCP into Android applications.

## Setup

### 1. Add Dependencies

Add these dependencies to your `app/build.gradle`:

```gradle
dependencies {
    // Networking - Retrofit
    implementation 'com.squareup.retrofit2:retrofit:2.9.0'
    implementation 'com.squareup.retrofit2:converter-gson:2.9.0'
    implementation 'com.squareup.okhttp3:logging-interceptor:4.11.0'
    
    // Optional: gRPC for high-performance
    implementation 'io.grpc:grpc-okhttp:1.58.0'
    implementation 'io.grpc:grpc-protobuf-lite:1.58.0'
    implementation 'io.grpc:grpc-stub:1.58.0'
    
    // Coroutines
    implementation 'org.jetbrains.kotlinx:kotlinx-coroutines-android:1.7.3'
    implementation 'org.jetbrains.kotlinx:kotlinx-coroutines-core:1.7.3'
    
    // ViewModel and Lifecycle
    implementation 'androidx.lifecycle:lifecycle-viewmodel-ktx:2.6.2'
    implementation 'androidx.lifecycle:lifecycle-runtime-ktx:2.6.2'
}
```

### 2. Add Internet Permission

Add to your `AndroidManifest.xml`:

```xml
<uses-permission android:name="android.permission.INTERNET" />
<uses-permission android:name="android.permission.ACCESS_NETWORK_STATE" />
```

### 3. Configure Network Security

Create `res/xml/network_security_config.xml`:

```xml
<?xml version="1.0" encoding="utf-8"?>
<network-security-config>
    <!-- For production with HTTPS -->
    <domain-config cleartextTrafficPermitted="false">
        <domain includeSubdomains="true">your-aionmcp-server.com</domain>
        <pin-set>
            <pin digest="SHA-256">YOUR_CERTIFICATE_PIN_HERE</pin>
        </pin-set>
    </domain-config>
    
    <!-- For development with HTTP (remove in production) -->
    <domain-config cleartextTrafficPermitted="true">
        <domain includeSubdomains="true">localhost</domain>
        <domain includeSubdomains="true">10.0.2.2</domain>
    </domain-config>
</network-security-config>
```

Update `AndroidManifest.xml`:

```xml
<application
    android:networkSecurityConfig="@xml/network_security_config"
    ...>
```

## Project Structure

```
app/src/main/java/com/example/aionmcp/
├── data/
│   ├── models/          # Data models
│   ├── api/             # API interfaces
│   └── repository/      # Repository pattern
├── ui/
│   ├── viewmodel/       # ViewModels
│   └── screens/         # UI components
└── di/                  # Dependency injection
```

## Implementation Files

### 1. Data Models (`data/models/AionMCPModels.kt`)

See [AionMCPModels.kt](./AionMCPModels.kt)

### 2. API Service (`data/api/AionMCPApiService.kt`)

See [AionMCPApiService.kt](./AionMCPApiService.kt)

### 3. Repository (`data/repository/AionMCPRepository.kt`)

See [AionMCPRepository.kt](./AionMCPRepository.kt)

### 4. ViewModel (`ui/viewmodel/AionMCPViewModel.kt`)

See [AionMCPViewModel.kt](./AionMCPViewModel.kt)

## Usage Example

### In Your Activity/Fragment

```kotlin
class MainActivity : AppCompatActivity() {
    private val viewModel: AionMCPViewModel by viewModels()
    
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
        
        // Observe health status
        viewModel.healthStatus.observe(this) { result ->
            result.onSuccess { health ->
                Log.d("AionMCP", "Server status: ${health.status}")
            }.onFailure { error ->
                Log.e("AionMCP", "Error: ${error.message}")
            }
        }
        
        // Check server health
        viewModel.checkHealth()
        
        // List available tools
        viewModel.listTools()
        
        // Observe tools
        viewModel.tools.observe(this) { result ->
            result.onSuccess { tools ->
                tools.forEach { tool ->
                    Log.d("AionMCP", "Tool: ${tool.name}")
                }
            }
        }
    }
}
```

### Invoke a Tool

```kotlin
lifecycleScope.launch {
    val args = mapOf("limit" to 10)
    viewModel.invokeTool("openapi.petstore.listPets", args)
    
    viewModel.toolResult.observe(this@MainActivity) { result ->
        result.onSuccess { response ->
            Log.d("AionMCP", "Result: ${response.result}")
        }
    }
}
```

## Advanced Features

### Caching

Implement caching for frequently accessed data:

```kotlin
class AionMCPRepository(
    private val api: AionMCPApiService,
    private val cacheDir: File
) {
    private val toolsCache = mutableMapOf<String, List<Tool>>()
    
    suspend fun listToolsCached(): Result<List<Tool>> {
        // Check cache first
        toolsCache["all"]?.let { return Result.success(it) }
        
        // Fetch from API
        return try {
            val response = api.listTools()
            toolsCache["all"] = response.tools
            Result.success(response.tools)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
}
```

### Offline Support

```kotlin
class OfflineFirstRepository(
    private val api: AionMCPApiService,
    private val database: AionMCPDatabase
) {
    suspend fun getToolsOfflineFirst(): Flow<List<Tool>> = flow {
        // Emit cached data first
        val cached = database.toolDao().getAllTools()
        emit(cached)
        
        // Then fetch fresh data
        try {
            val response = api.listTools()
            database.toolDao().insertAll(response.tools)
            emit(response.tools)
        } catch (e: Exception) {
            // Keep cached data if fetch fails
        }
    }
}
```

### Background Sync

```kotlin
class SyncWorker(
    context: Context,
    params: WorkerParameters
) : CoroutineWorker(context, params) {
    
    override suspend fun doWork(): Result {
        val repository = AionMCPRepository(/* inject dependencies */)
        
        return try {
            repository.syncTools()
            Result.success()
        } catch (e: Exception) {
            Result.retry()
        }
    }
}

// Schedule periodic sync
val syncRequest = PeriodicWorkRequestBuilder<SyncWorker>(15, TimeUnit.MINUTES)
    .setConstraints(
        Constraints.Builder()
            .setRequiredNetworkType(NetworkType.CONNECTED)
            .build()
    )
    .build()

WorkManager.getInstance(context).enqueue(syncRequest)
```

## Testing

### Unit Tests

```kotlin
class AionMCPRepositoryTest {
    private lateinit var mockApi: AionMCPApiService
    private lateinit var repository: AionMCPRepository
    
    @Before
    fun setup() {
        mockApi = mock()
        repository = AionMCPRepository(mockApi)
    }
    
    @Test
    fun `listTools returns success`() = runBlocking {
        // Given
        val expectedTools = listOf(
            Tool("tool1", "Description 1", emptyMap())
        )
        whenever(mockApi.listTools()).thenReturn(
            ToolsResponse("MCP/1.0", expectedTools)
        )
        
        // When
        val result = repository.listTools()
        
        // Then
        assertTrue(result.isSuccess)
        assertEquals(expectedTools, result.getOrNull())
    }
}
```

### Integration Tests

```kotlin
@RunWith(AndroidJUnit4::class)
class AionMCPIntegrationTest {
    
    @Test
    fun testServerConnection() = runBlocking {
        val client = AionMCPClient("http://10.0.2.2:8080")
        val health = client.api.getHealth()
        
        assertEquals("healthy", health.status)
    }
}
```

## Proguard Rules

Add to your `proguard-rules.pro`:

```proguard
# Retrofit
-dontwarn retrofit2.**
-keep class retrofit2.** { *; }
-keepattributes Signature
-keepattributes Exceptions

# OkHttp
-dontwarn okhttp3.**
-keep class okhttp3.** { *; }

# Gson
-keep class com.google.gson.** { *; }
-keep class * implements com.google.gson.TypeAdapterFactory
-keep class * implements com.google.gson.JsonSerializer
-keep class * implements com.google.gson.JsonDeserializer

# Your models
-keep class com.example.aionmcp.data.models.** { *; }
```

## Troubleshooting

### Common Issues

1. **Network on Main Thread Exception**
   - Solution: Always use coroutines or background threads

2. **SSL/TLS Errors**
   - Solution: Check network security config and certificate pinning

3. **Timeout Errors**
   - Solution: Increase timeout or check server availability

4. **JSON Parsing Errors**
   - Solution: Verify data models match API responses

### Enable Logging

```kotlin
val logging = HttpLoggingInterceptor().apply {
    level = if (BuildConfig.DEBUG) {
        HttpLoggingInterceptor.Level.BODY
    } else {
        HttpLoggingInterceptor.Level.NONE
    }
}

val okHttpClient = OkHttpClient.Builder()
    .addInterceptor(logging)
    .build()
```

## Next Steps

1. Implement error handling UI
2. Add authentication flow
3. Implement caching strategies
4. Add unit and integration tests
5. Optimize for battery and data usage

## Resources

- [Retrofit Documentation](https://square.github.io/retrofit/)
- [Kotlin Coroutines Guide](https://kotlinlang.org/docs/coroutines-guide.html)
- [Android Architecture Components](https://developer.android.com/topic/architecture)
- [gRPC Android Tutorial](https://grpc.io/docs/platforms/android/)
