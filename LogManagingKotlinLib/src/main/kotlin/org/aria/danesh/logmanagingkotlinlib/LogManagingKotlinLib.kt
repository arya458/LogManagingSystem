package org.aria.danesh.logmanagingkotlinlib

import kotlinx.coroutines.*
import kotlinx.serialization.*
import kotlinx.serialization.json.Json
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody.Companion.toRequestBody
import org.slf4j.LoggerFactory
import java.io.IOException
import java.util.concurrent.TimeUnit
import kotlin.time.Duration.Companion.seconds

/**
 * Data class representing the log data structure
 */
@Serializable
data class LogData(
    val imel: String,
    val error: String,
    val level: String,
    val source: String?,
    val timestamp: String = System.currentTimeMillis().toString()
)

/**
 * Data class representing a log entry
 */
@Serializable
data class LogEntry(
    val encrypted_data: String
)

/**
 * Custom exception for LogManagingKotlinLib errors
 */
class LogManagingException(message: String, cause: Throwable? = null) : Exception(message, cause)

/**
 * Main class for the Log Managing Kotlin Library.
 * Provides functionality for encrypting and sending log data to the LogManagingApi.
 */
class LogManagingKotlinLib(
    private val apiUrl: String,
    private val username: String,
    private val password: String,
    private val encryptionKey: String,
    private val maxRetries: Int = 3,
    private val retryDelay: Long = 1,
    private val timeout: Long = 30,
    private val client: OkHttpClient = createDefaultClient()
) {
    private val logger = LoggerFactory.getLogger(LogManagingKotlinLib::class.java)
    private val jsonMediaType = "application/json; charset=utf-8".toMediaType()
    private val json = Json { 
        ignoreUnknownKeys = true 
        encodeDefaults = true
    }

    companion object {
        private fun createDefaultClient(): OkHttpClient {
            return OkHttpClient.Builder()
                .connectTimeout(30, TimeUnit.SECONDS)
                .readTimeout(30, TimeUnit.SECONDS)
                .writeTimeout(30, TimeUnit.SECONDS)
                .retryOnConnectionFailure(true)
                .build()
        }
    }

    /**
     * Sends an encrypted log entry to the API.
     * @param imel The IMEL identifier
     * @param error The error message
     * @param level Optional log level (default: "ERROR")
     * @param source Optional source identifier
     * @return Result containing either a success message or an exception
     */
    suspend fun sendEncryptedLog(
        imel: String,
        error: String,
        level: String = "ERROR",
        source: String? = null
    ): Result<String> = withContext(Dispatchers.IO) {
        try {
            val logData = LogData(
                imel = imel,
                error = error,
                level = level,
                source = source
            )

            val encryptedData = EncryptionUtils.encrypt(
                json.encodeToString(LogData.serializer(), logData),
                encryptionKey
            )

            val logEntry = LogEntry(encryptedData)
            val requestBody = json.encodeToString(LogEntry.serializer(), logEntry)
                .toRequestBody(jsonMediaType)

            val request = Request.Builder()
                .url(apiUrl)
                .post(requestBody)
                .addHeader("Authorization", createBasicAuthHeader())
                .build()

            var retryCount = 0
            var lastException: Exception? = null

            while (retryCount < maxRetries) {
                try {
                    val response = client.newCall(request).execute()
                    if (response.isSuccessful) {
                        return@withContext Result.success("Log sent successfully")
                    }

                    val errorBody = response.body?.string() ?: "No error body"
                    logger.error("Failed to send log. Status: ${response.code}, Body: $errorBody")
                    lastException = LogManagingException("Failed to send log: $errorBody")
                    break
                } catch (e: IOException) {
                    lastException = e
                    retryCount++
                    if (retryCount < maxRetries) {
                        logger.warn("Retry attempt $retryCount after ${retryDelay.seconds}")
                        delay(retryDelay * 1000)
                    }
                }
            }

            Result.failure(lastException ?: LogManagingException("Unknown error occurred"))
        } catch (e: Exception) {
            logger.error("Error sending log", e)
            Result.failure(e)
        }
    }

    private fun createBasicAuthHeader(): String {
        val credentials = "$username:$password"
        val encodedCredentials = java.util.Base64.getEncoder().encodeToString(credentials.toByteArray())
        return "Basic $encodedCredentials"
    }
}

/**
 * Extension function to convert objects to JSON
 */
private fun Any.toJson(): String {
    return kotlinx.serialization.json.Json.encodeToString(this)
} 