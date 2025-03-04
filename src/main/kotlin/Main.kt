import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.runBlocking
import kotlinx.coroutines.withContext
import org.json.JSONObject
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

fun main() = runBlocking {
    println("Hello World!")
    val logManagingKotlinLib = LogManagingKotlinLib(
        "http://127.0.0.1:80",
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

//class LogManagingKotlinLib(
//    private val baseUrl: String,
//    private val username: String,
//    private val password: String,
//    private val encryptionKey: String
//) {
//
//    suspend fun sendEncryptedLog(imel: String, errorMessage: String): Result<String> =
//        withContext(Dispatchers.IO) {
//            try {
//                val jsonData = JSONObject().apply {
//                    put("imel", imel)
//                    put("error", errorMessage)
//                }
//
//                val encryptedData = encrypt(jsonData.toString(), encryptionKey)
//
//                val postData = JSONObject().apply {
//                    put("encrypted_data", encryptedData)
//                }.toString()
//
//                val credentials = "$username:$password"
//                val encodedCredentials = Base64.getEncoder()
//                    .encodeToString(credentials.toByteArray(StandardCharsets.UTF_8))
//
//                val url = URL(baseUrl)
//                val connection = url.openConnection() as HttpURLConnection
//                connection.requestMethod = "POST"
//                connection.setRequestProperty("Authorization", "Basic $encodedCredentials")
//                connection.setRequestProperty("Content-Type", "application/json")
//                connection.doOutput = true
//
//                DataOutputStream(connection.outputStream).use {
//                    it.write(postData.toByteArray(StandardCharsets.UTF_8))
//                }
//
//                val responseCode = connection.responseCode
//                if (responseCode == HttpURLConnection.HTTP_OK) {
//                    BufferedReader(InputStreamReader(connection.inputStream)).use { reader ->
//                        val response = reader.readText()
//                        return@withContext Result.success(response)
//                    }
//                } else {
//                    BufferedReader(InputStreamReader(connection.errorStream)).use { reader ->
//                        val errorResponse = reader.readText()
//                        return@withContext Result.failure(
//                            IOException("HTTP error: $responseCode, $errorResponse")
//                        )
//                    }
//                }
//                return@withContext Result.failure(IOException("Unexpected end of request"))
//            } catch (e: Exception) {
//                return@withContext Result.failure(e)
//            }
//        }
//
//    private fun encrypt(data: String, key: String): String {
//        val salt = ByteArray(16)
//        SecureRandom().nextBytes(salt)
//
//        val derivedKey = deriveKey(key, salt)
//        val secretKeySpec = SecretKeySpec(derivedKey, "AES")
//
//        val iv = ByteArray(16)
//        SecureRandom().nextBytes(iv)
//
//        val ivParameterSpec = IvParameterSpec(iv)
//        val cipher = Cipher.getInstance("AES/CFB/NoPadding")
//        cipher.init(Cipher.ENCRYPT_MODE, secretKeySpec, ivParameterSpec)
//
//        val encrypted = cipher.doFinal(data.toByteArray(StandardCharsets.UTF_8))
//        val combined = salt + iv + encrypted
//
//        return java.util.Base64.getEncoder().encodeToString(combined)
//    }
//
//    private fun deriveKey(password: String, salt: ByteArray): ByteArray {
//        val combined = password.toByteArray(StandardCharsets.UTF_8) + salt
//        val sha256Hash = java.security.MessageDigest.getInstance("SHA-256").digest(combined)
//        return sha256Hash.copyOfRange(0, 32)
//    }
//}