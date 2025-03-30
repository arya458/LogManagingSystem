package org.aria.danesh.logmanagingkotlinlib

import org.slf4j.LoggerFactory
import java.security.SecureRandom
import javax.crypto.Cipher
import javax.crypto.spec.IvParameterSpec
import javax.crypto.spec.SecretKeySpec
import java.util.Base64
import java.security.MessageDigest

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
     * Encrypts the given data using AES-CFB encryption with a derived key.
     * @param data The data to encrypt
     * @param key The encryption key
     * @return Base64 encoded encrypted data
     */
    fun encrypt(data: String, key: String): String {
        try {
            // Generate random salt
            val salt = generateRandomBytes(SALT_LENGTH)
            
            // Generate random IV
            val iv = generateRandomBytes(IV_LENGTH)
            
            // Derive key from password and salt
            val derivedKey = deriveKey(key, salt)
            
            // Create cipher instance
            val cipher = Cipher.getInstance(ALGORITHM)
            val keySpec = SecretKeySpec(derivedKey, KEY_ALGORITHM)
            val ivSpec = IvParameterSpec(iv)
            
            // Initialize cipher for encryption
            cipher.init(Cipher.ENCRYPT_MODE, keySpec, ivSpec)
            
            // Encrypt the data
            val encryptedData = cipher.doFinal(data.toByteArray())
            
            // Combine salt, IV, and encrypted data
            val combined = ByteArray(salt.size + iv.size + encryptedData.size)
            System.arraycopy(salt, 0, combined, 0, salt.size)
            System.arraycopy(iv, 0, combined, salt.size, iv.size)
            System.arraycopy(encryptedData, 0, combined, salt.size + iv.size, encryptedData.size)
            
            // Encode as Base64
            return Base64.getEncoder().encodeToString(combined)
        } catch (e: Exception) {
            logger.error("Encryption failed", e)
            throw LogManagingException("Encryption failed: ${e.message}", e)
        }
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
     * Generates random bytes using SecureRandom.
     * @param length The length of the random bytes to generate
     * @return Array of random bytes
     */
    private fun generateRandomBytes(length: Int): ByteArray {
        val bytes = ByteArray(length)
        SecureRandom().nextBytes(bytes)
        return bytes
    }

    /**
     * Derives a key from the password and salt using SHA-256.
     * @param password The password
     * @param salt The salt
     * @return Derived key
     */
    private fun deriveKey(password: String, salt: ByteArray): ByteArray {
        val md = MessageDigest.getInstance("SHA-256")
        val combined = password.toByteArray() + salt
        return md.digest(combined)
    }
} 