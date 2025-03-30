package org.aria.danesh.logmanagingkotlinlib

import org.slf4j.LoggerFactory
import java.security.SecureRandom
import javax.crypto.Cipher
import javax.crypto.spec.IvParameterSpec
import javax.crypto.spec.SecretKeySpec
import java.util.Base64
import java.security.MessageDigest
import java.nio.charset.StandardCharsets

/**
 * Utility class for handling encryption operations.
 */
object EncryptionUtils {
    private val logger = LoggerFactory.getLogger(EncryptionUtils::class.java)
    private const val ALGORITHM = "AES/CFB/NoPadding"
    private const val KEY_ALGORITHM = "AES"
    private const val SALT_LENGTH = 16
    private const val IV_LENGTH = 16

    /**
     * Encrypts data using AES/CFB with a derived key
     * @param data The data to encrypt
     * @param key The encryption key
     * @return Base64 encoded encrypted data
     */
    fun encrypt(data: String, key: String): String {
        // Generate a random salt
        val salt = ByteArray(16).apply {
            SecureRandom().nextBytes(this)
        }

        // Derive the encryption key
        val derivedKey = deriveKey(key, salt)
        val secretKeySpec = SecretKeySpec(derivedKey, "AES")

        // Generate a random IV
        val iv = ByteArray(16).apply {
            SecureRandom().nextBytes(this)
        }
        val ivParameterSpec = IvParameterSpec(iv)

        // Initialize cipher
        val cipher = Cipher.getInstance("AES/CFB/NoPadding")
        cipher.init(Cipher.ENCRYPT_MODE, secretKeySpec, ivParameterSpec)

        // Encrypt the data
        val encrypted = cipher.doFinal(data.toByteArray(StandardCharsets.UTF_8))

        // Combine salt + IV + encrypted data
        val combined = ByteArray(salt.size + iv.size + encrypted.size).apply {
            System.arraycopy(salt, 0, this, 0, salt.size)
            System.arraycopy(iv, 0, this, salt.size, iv.size)
            System.arraycopy(encrypted, 0, this, salt.size + iv.size, encrypted.size)
        }

        // Return Base64 encoded result
        return Base64.getEncoder().encodeToString(combined)
    }

    /**
     * Decrypts the given encrypted data using AES-CFB decryption with a derived key.
     * @param encryptedData The Base64 encoded encrypted data
     * @param key The encryption key
     * @return Decrypted data
     */
    fun decrypt(encryptedData: String, key: String): String {
        try {
            // Decode Base64
            val combined = Base64.getDecoder().decode(encryptedData)
            
            // Extract salt, IV, and encrypted data
            val salt = combined.copyOfRange(0, SALT_LENGTH)
            val iv = combined.copyOfRange(SALT_LENGTH, SALT_LENGTH + IV_LENGTH)
            val data = combined.copyOfRange(SALT_LENGTH + IV_LENGTH, combined.size)
            
            // Derive key from password and salt
            val derivedKey = deriveKey(key, salt)
            
            // Create cipher instance
            val cipher = Cipher.getInstance(ALGORITHM)
            val keySpec = SecretKeySpec(derivedKey, KEY_ALGORITHM)
            val ivSpec = IvParameterSpec(iv)
            
            // Initialize cipher for decryption
            cipher.init(Cipher.DECRYPT_MODE, keySpec, ivSpec)
            
            // Decrypt the data
            val decryptedData = cipher.doFinal(data)
            
            return String(decryptedData)
        } catch (e: Exception) {
            logger.error("Decryption failed", e)
            throw LogManagingException("Decryption failed: ${e.message}", e)
        }
    }

    /**
     * Derives a key from the password and salt using SHA-256
     * @param password The password to derive from
     * @param salt The salt to use
     * @return The derived key
     */
    private fun deriveKey(password: String, salt: ByteArray): ByteArray {
        val combined = password.toByteArray(StandardCharsets.UTF_8) + salt
        val sha256Hash = MessageDigest.getInstance("SHA-256").digest(combined)
        return sha256Hash.copyOfRange(0, 32)
    }
} 