Log Managing System

This project provides a comprehensive log management solution, consisting of two main components:

LogManagingApi (Go): A backend API that receives, decrypts, and stores log data.
LogManagingKotlinLib (Kotlin/Gradle): A client library for encrypting and sending log data to the API.
Project Structure

```kotlin
LogManagingSystem/
├── LogManagingApi/       (Go API project)
│   ├── logmanagingapi.go
│   ├── go.mod
│   ├── go.sum
│   └── build.bat
├── LogManagingKotlinLib/ (Kotlin/Gradle library project)
│   ├── src/
│   │   └── main/kotlin/
│   │       └── org.aria.danesh/
│   │           └── logmanagingkotlinlib/
│   │               ├── LogManagingKotlinLib.kt
│   │               ├── EncryptionUtils.kt
│   │               ├── NetworkUtils.kt
│   │               └── DataModels.kt
│   ├── build.gradle.kts
│   └── lib/              (Directory to hold the LogManagingKotlinLib.jar)
├── build/                 (Output directory for build artifacts)
│   ├── go_executables/    (Go executables)
│   └── libs/             (Gradle JAR library)
└── build.bat             (Batch script to build the entire project)
```
LogManagingApi (Go)

Description

The LogManagingApi is a Go-based HTTP server designed to receive, decrypt, and store log data.

Functionality:
Receives encrypted log data in JSON format via POST requests.
Authenticates requests using Basic Authentication.
Decrypts log data using AES-CFB encryption with a derived key.
Parses decrypted JSON data to extract log information (IMEL, error message).
Stores log information in separate log files, organized by IMEL and timestamp.
Building:
Navigate to the project root and run build.bat. This builds the API.
The build.bat will build the Go api executables to the build/go_executables folder.
Running:
Run the appropriate executable from the build/go_executables directory.
Configuration:
Username, password, and encryption key are set within logmanagingapi.go.
Port is configured via the PORT environment variable (default: 80).
LogManagingKotlinLib (Kotlin/Gradle)

Description

The LogManagingKotlinLib is a Kotlin library that provides functionality for encrypting and sending log data to the LogManagingApi.

Packages:
org.aria.danesh.logmanagingkotlinlib:
LogManagingKotlinLib.kt: Main class for interacting with the library.
EncryptionUtils.kt: Handles AES-CFB encryption and key derivation.
NetworkUtils.kt: Manages HTTP POST requests to the API.
DataModels.kt: Defines data structures for log data and API responses.
Use Cases:
Encryption: The library encrypts log data (IMEL, error message) using AES-CFB with a derived key from a shared secret.
Data Transmission: Sends the encrypted data to the LogManagingApi via HTTP POST requests.
Response Handling: Processes HTTP responses, providing success or failure indications.
Building:
Navigate to the project root and run build.bat. This builds the library.
The library jar will be placed in the build/libs folder.
Integration:
Copy the JAR from build/libs to the lib directory of your project.
Add the JAR as a dependency in your Gradle build.
Usage:
Kotlin

import kotlinx.coroutines.runBlocking
import org.aria.danesh.logmanagingkotlinlib.LogManagingKotlinLib
```kotlin
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
Build Script (build.bat)

The build.bat script automates the build process for both the Go API and the Kotlin library.

Usage: Place build.bat in the project root and run it.
The build script will build both the go api, and the Kotlin library.
How This Lib Works

Initialization: The LogManagingKotlinLib is initialized with the API endpoint, authentication credentials, and the shared encryption key.
Encryption: When sendEncryptedLog is called, the provided IMEL and error message are serialized into a JSON string. This string is then encrypted using AES-CFB with a key derived from the shared secret.
Transmission: The encrypted data is sent to the LogManagingApi via an HTTP POST request, along with Basic Authentication headers.
Decryption and Storage: The API decrypts the data, parses it, and stores the log information.
Response: The library returns a Result type, indicating success or failure.

Developed by: <img src="https://avatars.githubusercontent.com/u/23719966?v=4" width="20" height="20"> [Arya](https://github.com/arya458) [Aria](https://github.com/arya458)
