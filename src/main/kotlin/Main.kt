import kotlinx.coroutines.runBlocking
import org.aria.danesh.logmanagingkotlinlib.LogManagingKotlinLib

fun main() = runBlocking {
    println("Hello World!") // Simple initial print statement.

    // Initialize the LogManagingKotlinLib with API details.
    val logManagingKotlinLib = LogManagingKotlinLib(
        "http://127.0.0.1:80", // API endpoint URL.
        "myuser", // Username for basic authentication.
        "mypassword", // Password for basic authentication.
        "your32bytekeyhere!" // Encryption key.
    )

    // Send an encrypted log message.
    val result = logManagingKotlinLib.sendEncryptedLog("asdasdasd", "TestLib")

    // Handle the result of the log sending operation.
    result.fold(
        onSuccess = { message ->
            println("Log sent successfully: $message") // Print success message.
        },
        onFailure = { exception ->
            println("Error sending log: ${exception.message}") // Print error message.
            exception.printStackTrace() // Print the stack trace of the exception.
        }
    )
}