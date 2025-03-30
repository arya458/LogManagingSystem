# Log Managing System

A secure, efficient, and scalable log management solution built with Kotlin and Go. This system provides a robust way to collect, encrypt, and store application logs with high performance and security.

## 🌟 Features

- 🔒 **Secure Transmission**: AES encryption with CFB mode and key derivation
- 🔑 **Authentication**: Basic authentication for API endpoints
- 📁 **Organized Storage**: Logs stored by device identifier (IMEL)
- ⚡ **High Performance**: Go backend for fast processing
- 🌐 **Cross-Platform**: Kotlin client for easy integration
- 🔄 **Retry Mechanism**: Built-in retry logic for reliability

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or later
- Kotlin 1.9.21 or later
- Java 17 or later
- Gradle 8.0 or later

### Building the System

1. Clone the repository:
   ```bash
   git clone https://github.com/arya458/LogManagingSystem.git
   cd LogManagingSystem
   ```

2. Build the Go API:
   ```bash
   ./build.bat
   ```
   The executable will be created in `build/api/LogManagingApi-windows-64.exe`

3. Build the Kotlin library:
   ```bash
   cd LogManagingKotlinLib
   ./gradlew build
   ```
   The JAR will be created in `build/libs/`

## 📚 Usage

### Kotlin Client

Add the following dependencies to your Kotlin project:

```kotlin
dependencies {
    implementation("org.aria.danesh:logmanagingkotlinlib:1.0.0")
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.7.3")
    implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:1.6.0")
    implementation("com.squareup.okhttp3:okhttp:4.12.0")
}
```

Example usage:

```kotlin
import kotlinx.coroutines.runBlocking
import org.aria.danesh.logmanagingkotlinlib.LogManagingKotlinLib

fun main() = runBlocking {
    val logManager = LogManagingKotlinLib(
        apiUrl = "http://your-api-host:port",
        username = "your-username",
        password = "your-password",
        encryptionKey = "your-32-byte-encryption-key",
        maxRetries = 3,
        retryDelay = 1,
        timeout = 30
    )

    val result = logManager.sendEncryptedLog(
        imel = "device-identifier",
        error = "Error message",
        level = "ERROR",
        source = "application-name"
    )

    result.fold(
        onSuccess = { message -> println("Log sent successfully: $message") },
        onFailure = { exception -> println("Error sending log: ${exception.message}") }
    )
}
```

### Go API

The Go API provides the following endpoints:

- `POST /` - Send encrypted logs
- `GET /health` - Check server health

Example using curl:
```bash
curl -X POST \
  -H "Authorization: Basic $(echo -n 'username:password' | base64)" \
  -H "Content-Type: application/json" \
  -d '{"encrypted_data":"your-encrypted-data"}' \
  http://localhost:8080/
```

## 🔒 Security

- All log data is encrypted using AES encryption with CFB mode
- Key derivation using SHA-256 for enhanced security
- Basic authentication for API access
- Secure credential handling
- No sensitive data in version control

## 🛠️ Development

### Project Structure

```
LogManagingSystem/
├── LogManagingApi/          # Go backend API
│   ├── main.go             # Main API implementation
│   └── build.bat           # Build script
├── LogManagingKotlinLib/   # Kotlin client library
│   ├── src/                # Source code
│   └── build.gradle.kts    # Build configuration
└── docs/                   # Documentation
    └── index.html         # Documentation website
```

### Building from Source

1. Go API:
   ```bash
   cd LogManagingApi
   go build -o LogManagingApi.exe
   ```

2. Kotlin Library:
   ```bash
   cd LogManagingKotlinLib
   ./gradlew build
   ```

## 🤝 Contributing

We welcome contributions! Please feel free to submit pull requests or report issues.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

- **Aria Danesh** - [GitHub](https://github.com/arya458)

## 🙏 Acknowledgments

- Go team for the excellent standard library
- Kotlin team for the amazing language and ecosystem
- All contributors who help improve this project
