package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LogEntry defines the structure of the incoming JSON log data.
type LogEntry struct {
	EncryptedData string `json:"encrypted_data"`
}

// Global variables for authentication and encryption.
var (
	username      = "myuser"             // Username for basic authentication.
	password      = "mypassword"         // Password for basic authentication.
	encryptionKey = "your32bytekeyhere!" // Encryption key for log data.
)

// basicAuth middleware for HTTP request authentication.
func basicAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}

		authParts := strings.Split(auth, " ")
		if len(authParts) != 2 || authParts[0] != "Basic" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(authParts[1])
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}

		credentials := strings.Split(string(decoded), ":")
		if len(credentials) != 2 || credentials[0] != username || credentials[1] != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}

		handler(w, r)
	}
}

// decrypt decrypts the encrypted log data.
func decrypt(encryptedData string, key string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	salt := decoded[0:16]
	iv := decoded[16:32]
	ciphertext := decoded[32:]

	derivedKey := deriveKey(key, salt)

	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return "", err
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

// deriveKey derives a key from the password and salt using SHA256.
func deriveKey(password string, salt []byte) []byte {
	combined := []byte(password)
	combined = append(combined, salt...)
	hash := sha256.Sum256(combined)
	return hash[:32]
}

// indexHandler handles incoming log entries, decrypts them, and saves them to files.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request path is root.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// Check if the request method is POST.
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// Read the request body.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	// Unmarshal the JSON request body into a LogEntry struct.
	var entry LogEntry
	err = json.Unmarshal(body, &entry)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	// Check if the encrypted data is present.
	if entry.EncryptedData == "" {
		http.Error(w, "Missing encrypted data", http.StatusBadRequest)
		return
	}
	// Decrypt the encrypted log data.
	decrypted, err := decrypt(entry.EncryptedData, encryptionKey)
	if err != nil {
		http.Error(w, "Decryption failed", http.StatusInternalServerError)
		return
	}
	// Unmarshal the decrypted JSON data into a map.
	var logData map[string]string
	err = json.Unmarshal([]byte(decrypted), &logData)
	if err != nil {
		http.Error(w, "Decrypted data is not valid JSON", http.StatusBadRequest)
		return
	}
	// Extract IMEL and error message from the decrypted data.
	imel, okImel := logData["imel"]
	errorMessage, okError := logData["error"]
	// Check if IMEL and error message are present in the decrypted data.
	if !okImel || !okError {
		http.Error(w, "Missing imel or error in decrypted data", http.StatusBadRequest)
		return
	}
	// Create a folder for the IMEL if it doesn't exist.
	folderPath := filepath.Join("logs", imel)
	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		http.Error(w, "Failed to create folder", http.StatusInternalServerError)
		return
	}
	// Create a log file with a timestamp.
	currentTime := time.Now()
	timeString := currentTime.Format("2006-01-02_15-04-05")
	filePath := filepath.Join(folderPath, timeString+".log")
	file, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to create log file", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	// Write the log data to the file.
	_, err = file.WriteString(fmt.Sprintf("IMEL: %s\nTime: %s\nError: %s\n", imel, timeString, errorMessage))
	if err != nil {
		http.Error(w, "Failed to write to log file", http.StatusInternalServerError)
		return
	}
	// Send a success response.
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Log entry saved successfully")
}

// encrypt encrypts the given data with the provided key.
func encrypt(data string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(data))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// main function to start the HTTP server.
func main() {
	// Handle all requests to the root path with basic authentication.
	http.HandleFunc("/", basicAuth(indexHandler))
	// Get the port from the environment variable or use the default port 80.
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
		log.Printf("Defaulting to port %s", port)
	}
	// Encode the username and password for curl command example.
	encodedCreds := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	// Create example JSON data for encryption and curl command.
	jsonData := map[string]string{"imel": "123456789012345", "error": "Test error"}
	jsonBytes, _ := json.Marshal(jsonData)
	derivedKey := deriveKey(encryptionKey, make([]byte, 16))
	encrypted, _ := encrypt(string(jsonBytes), derivedKey)

	// Create a curl command example for testing.
	curlCommand := fmt.Sprintf("curl -X POST -H \"Authorization: Basic %s\" -H \"Content-Type: application/json\" -d '{\"encrypted_data\":\"%s\"}' http://localhost:%s/", encodedCreds, encrypted, port)

	// Log the listening port and curl command example.
	log.Printf("Listening on port %s", port)
	log.Printf("Open http://localhost:%s in the browser", port)
	log.Printf("Use this curl command to test the endpoint:\n%s", curlCommand)

	// Start the HTTP server.
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
