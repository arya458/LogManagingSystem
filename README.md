# Log Managing System

This project consists of two main components:

1.  **LogManagingApi (Go):** A Go-based API that receives encrypted log data, decrypts it, and stores it in log files.
2.  **LogManagingKotlinLib (Kotlin/Gradle):** A Kotlin library that encrypts log data before sending it to the Go API.

## Project Structure
```kotlin
LogManagingSystem/
├── LogManagingApi/       (Go API project)
│   ├── logmanagingapi.go
│   ├── go.mod
│   ├── go.sum
│   ├── build.bat 
│   └── ...
├── LogManagingKotlinLib/ (Kotlin/Gradle library project)
│   ├── src/
│   │   └── main/kotlin/
│   │       └── LogManagingKotlinLib.kt
│   ├── build.gradle.kts
│   └── ...
├── build/                 (Output directory for build artifacts)
│   ├── go_executables/    (Go executables)
│   └── libs/             (Gradle JAR library)
└── build_all.bat         (Batch script to build the entire project) NOT YET
└── lib/                  (Directory to hold the LogManagingKotlinLib.jar)
```

## LogManagingApi (Go)

### Description

The `LogManagingApi` is a Go-based HTTP server that:

* Receives encrypted log data in JSON format via POST requests.
* Authenticates requests using Basic Authentication.
* Decrypts the log data using AES-CFB encryption with a derived key from a shared secret.
* Parses the decrypted JSON data to extract log information (IMEL, error message).
* Stores the log information in separate log files, organized by IMEL and timestamp.

### Prerequisites

* Go (version 1.16 or later)

### Building

1.  Navigate to the project root directory.
2.  Run `build_all.bat`. This will build executables for multiple OS in the `build/go_executables` directory.
3.  Alternatively, you can navigate to the `LogManagingApi` directory and run `go build` to build for your current OS.

### Running

1.  Navigate to the `build/go_executables` directory.
2.  Run the appropriate executable for your operating system (e.g., `myprogram.exe` on Windows, `myprogram_linux_amd64` on Linux).

### Configuration

* **Username and Password:** Set in the `logmanagingapi.go` file.
* **Encryption Key:** Set in the `logmanagingapi.go` file.
* **Port:** Configured through the `PORT` environment variable (defaults to 80).

## LogManagingKotlinLib (Kotlin/Gradle)

### Description

The `LogManagingKotlinLib` is a Kotlin library that:

* Encrypts log data (IMEL, error message) using AES-CFB encryption with a derived key.
* Sends the encrypted data to the `LogManagingApi` via HTTP POST requests.
* Handles HTTP responses and reports success or failure.

### Prerequisites

* Java Development Kit (JDK)
* Gradle

### Building

1.  Navigate to the project root directory.
2.  Run `build_all.bat`. This will build a jar file to the `build/libs` directory.
3.  Alternatively, you can navigate to the `LogManagingKotlinLib` directory and run `./gradlew jar` to build the JAR file.

### Integration

1.  Copy the generated JAR file from `build/libs` into the `lib` directory of your project.
2.  Add the following dependency to your `build.gradle.kts` file:

    ```kotlin
    dependencies {
        implementation(files("./lib/LogManagingKotlinLib.jar"))
        // ... other dependencies
    }
    ```

3.  Sync your Gradle project.

### Usage

```kotlin
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val logManagingKotlinLib = LogManagingKotlinLib(
        "http://<api_host>:<api_port>",
        "myuser",
        "mypassword",
        "your32bytekeyhere!"
    )

    val result = logManagingKotlinLib.sendEncryptedLog("asdasdasd", "TestLib")

    result.fold(
        onSuccess = { message ->
            println("Log sent successfully: $message")
        },
        onFailure = { exception ->
            println("Error sending log: ${exception.message}")
            exception.printStackTrace()
        }
    )
}
```
Build Script (build_all.bat)
The build.bat script automates the build process for both the Go API.

Usage
Place the build_all.bat file in the project's root directory.
Double-click the build_all.bat file.
