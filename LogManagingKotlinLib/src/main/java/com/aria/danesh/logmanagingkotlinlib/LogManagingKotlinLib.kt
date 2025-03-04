import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import java.io.BufferedReader
import java.io.DataOutputStream
import java.io.IOException
import java.io.InputStreamReader
import java.net.HttpURLConnection
import java.net.URL
import java.nio.charset.StandardCharsets
import java.security.SecureRandom
import java.util.Base64
import javax.crypto.Cipher
import javax.crypto.spec.IvParameterSpec
import javax.crypto.spec.SecretKeySpec
import org.json.JSONObject

// LogManagingKotlinLib class for handling encrypted log sending.
class LogManagingKotlinLib(
    private val baseUrl: String, // Base URL of the API.
    private val username: String, // Username for basic authentication.
    private val password: String, // Password for basic authentication.
    private val encryptionKey: String // Encryption key.
) {

    // Sends an encrypted log message to the API.
    suspend fun sendEncryptedLog(imel: String, errorMessage: String): Result<String> = withContext(Dispatchers.IO) {
        try {
            // Create a JSON object with IMEL and error message.
            val jsonData = JSONObject().apply {
                put("imel", imel)
                put("error", errorMessage)
            }

            // Encrypt the JSON data.
            val encryptedData = encrypt(jsonData.toString(), encryptionKey)

            // Create a JSON object with the encrypted data.
            val postData = JSONObject().apply {
                put("encrypted_data", encryptedData)
            }.toString()

            // Encode credentials for basic authentication.
            val credentials = "$username:$password"
            val encodedCredentials = Base64.getEncoder().encodeToString(credentials.toByteArray(StandardCharsets.UTF_8))

            // Open a connection to the API.
            val url = URL(baseUrl)
            val connection = url.openConnection() as HttpURLConnection
            connection.requestMethod = "POST"
            connection.setRequestProperty("Authorization", "Basic $encodedCredentials")
            connection.setRequestProperty("Content-Type", "application/json")
            connection.doOutput = true

            // Write the encrypted data to the request body.
            DataOutputStream(connection.outputStream).use { it.write(postData.toByteArray(StandardCharsets.UTF_8)) }

            // Get the response code.
            val responseCode = connection.responseCode

            // Handle successful response.
            if (responseCode == HttpURLConnection.HTTP_OK) {
                // Read the response from the input stream.
                BufferedReader(InputStreamReader(connection.inputStream)).use { reader ->
                    val response = reader.readText()
                    return@withContext Result.success(response) // Return success with the response.
                }
            } else {
                // Handle error response.
                BufferedReader(InputStreamReader(connection.errorStream)).use { reader ->
                    val errorResponse = reader.readText()
                    return@withContext Result.failure(IOException("HTTP error: $responseCode, $errorResponse")) // Return failure with error details.
                }
            }
            return@withContext Result.failure(IOException("Unexpected end of request")) // Return failure if no response is received.
        } catch (e: Exception) {
            return@withContext Result.failure(e) // Return failure for any exception.
        }
    }

    // Encrypts the data using AES/CFB with a derived key.
    private fun encrypt(data: String, key: String): String {
        // Generate a random salt.
        val salt = ByteArray(16)
        SecureRandom().nextBytes(salt)

        // Derive the encryption key from the provided key and salt.
        val derivedKey = deriveKey(key, salt)
        val secretKeySpec = SecretKeySpec(derivedKey, "AES")

        // Generate a random initialization vector (IV).
        val iv = ByteArray(16)
        SecureRandom().nextBytes(iv)

        val ivParameterSpec = IvParameterSpec(iv)
        val cipher = Cipher.getInstance("AES/CFB/NoPadding")
        cipher.init(Cipher.ENCRYPT_MODE, secretKeySpec, ivParameterSpec)

        // Encrypt the data.
        val encrypted = cipher.doFinal(data.toByteArray(StandardCharsets.UTF_8))
        // Combine the salt, IV, and encrypted data.
        val combined = salt + iv + encrypted

        // Base64 encode the combined data.
        return java.util.Base64.getEncoder().encodeToString(combined)
    }

    // Derives a key from the password and salt using SHA-256.
    private fun deriveKey(password: String, salt: ByteArray): ByteArray {
        // Combine the password and salt.
        val combined = password.toByteArray(StandardCharsets.UTF_8) + salt
        // Generate the SHA-256 hash.
        val sha256Hash = java.security.MessageDigest.getInstance("SHA-256").digest(combined)
        // Return the first 32 bytes of the hash as the derived key.
        return sha256Hash.copyOfRange(0, 32)
    }
}