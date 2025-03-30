package main

import (
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
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
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

// Config holds the application configuration
type Config struct {
	Username      string
	Password      string
	EncryptionKey string
	Port          string
	LogDir        string
	MaxLogSize    int64
	MaxLogAge     time.Duration
	RateLimit     int
}

// Metrics holds application metrics
type Metrics struct {
	TotalRequests   uint64
	FailedRequests  uint64
	SuccessfulLogs  uint64
	FailedLogs      uint64
	LastRequestTime time.Time
	AverageResponse uint64
}

var (
	config  Config
	metrics Metrics
	logger  *log.Logger
)

// LogEntry defines the structure of the incoming JSON log data
type LogEntry struct {
	EncryptedData string `json:"encrypted_data"`
}

// LogData defines the structure of the decrypted log data
type LogData struct {
	IMEL   string `json:"imel"`
	Error  string `json:"error"`
	Level  string `json:"level,omitempty"`
	Source string `json:"source,omitempty"`
	Time   string `json:"time,omitempty"`
}

// Response represents the API response structure
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// loadConfig loads configuration from environment variables
func loadConfig() Config {
	return Config{
		Username:      getEnvOrDefault("API_USERNAME", "myuser"),
		Password:      getEnvOrDefault("API_PASSWORD", "mypassword"),
		EncryptionKey: getEnvOrDefault("API_ENCRYPTION_KEY", "your32bytekeyhere!"),
		Port:          getEnvOrDefault("PORT", "80"),
		LogDir:        getEnvOrDefault("LOG_DIR", "logs"),
		MaxLogSize:    10 * 1024 * 1024,   // 10MB
		MaxLogAge:     7 * 24 * time.Hour, // 7 days
		RateLimit:     100,                // requests per minute
	}
}

// setupLogger initializes the logger with rotation
func setupLogger() error {
	logFile := filepath.Join(config.LogDir, "api.log")
	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
	return nil
}

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, status int, data Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// healthCheckHandler handles health check requests
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	health := struct {
		Status    string    `json:"status"`
		Time      time.Time `json:"time"`
		Uptime    string    `json:"uptime"`
		GoVersion string    `json:"go_version"`
	}{
		Status:    "healthy",
		Time:      time.Now(),
		Uptime:    time.Since(startTime).String(),
		GoVersion: runtime.Version(),
	}
	writeJSON(w, http.StatusOK, Response{Success: true, Data: health})
}

// metricsHandler handles metrics requests
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, Response{Success: true, Data: metrics})
}

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
		if len(credentials) != 2 || credentials[0] != config.Username || credentials[1] != config.Password {
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

// indexHandler handles incoming log entries with improved error handling and metrics
func indexHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	atomic.AddUint64(&metrics.TotalRequests, 1)
	metrics.LastRequestTime = time.Now()

	defer func() {
		atomic.AddUint64(&metrics.AverageResponse, uint64(time.Since(start)))
	}()

	if r.URL.Path != "/" {
		writeJSON(w, http.StatusNotFound, Response{Success: false, Message: "Not found"})
		return
	}

	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, Response{Success: false, Message: "Method not allowed"})
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		atomic.AddUint64(&metrics.FailedRequests, 1)
		writeJSON(w, http.StatusBadRequest, Response{Success: false, Message: "Failed to read request body"})
		return
	}
	defer r.Body.Close()

	var entry LogEntry
	if err := json.Unmarshal(body, &entry); err != nil {
		atomic.AddUint64(&metrics.FailedRequests, 1)
		writeJSON(w, http.StatusBadRequest, Response{Success: false, Message: "Invalid JSON format"})
		return
	}

	if entry.EncryptedData == "" {
		atomic.AddUint64(&metrics.FailedRequests, 1)
		writeJSON(w, http.StatusBadRequest, Response{Success: false, Message: "Missing encrypted data"})
		return
	}

	decrypted, err := decrypt(entry.EncryptedData, config.EncryptionKey)
	if err != nil {
		atomic.AddUint64(&metrics.FailedRequests, 1)
		writeJSON(w, http.StatusInternalServerError, Response{Success: false, Message: "Decryption failed"})
		return
	}

	var logData LogData
	if err := json.Unmarshal([]byte(decrypted), &logData); err != nil {
		atomic.AddUint64(&metrics.FailedRequests, 1)
		writeJSON(w, http.StatusBadRequest, Response{Success: false, Message: "Invalid log data format"})
		return
	}

	if err := saveLog(logData); err != nil {
		atomic.AddUint64(&metrics.FailedLogs, 1)
		writeJSON(w, http.StatusInternalServerError, Response{Success: false, Message: "Failed to save log"})
		return
	}

	atomic.AddUint64(&metrics.SuccessfulLogs, 1)
	writeJSON(w, http.StatusOK, Response{Success: true, Message: "Log entry saved successfully"})
}

// saveLog saves the log entry with rotation and compression
func saveLog(logData LogData) error {
	folderPath := filepath.Join(config.LogDir, logData.IMEL)
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return fmt.Errorf("failed to create folder: %v", err)
	}

	currentTime := time.Now()
	timeString := currentTime.Format("2006-01-02_15-04-05")
	filePath := filepath.Join(folderPath, timeString+".log")

	// Check if we need to rotate logs
	if err := rotateLogs(folderPath); err != nil {
		logger.Printf("Warning: Failed to rotate logs: %v", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create log file: %v", err)
	}
	defer file.Close()

	// Write log data
	logContent := fmt.Sprintf("IMEL: %s\nTime: %s\nError: %s\nLevel: %s\nSource: %s\n",
		logData.IMEL, timeString, logData.Error, logData.Level, logData.Source)

	if _, err := file.WriteString(logContent); err != nil {
		return fmt.Errorf("failed to write to log file: %v", err)
	}

	return nil
}

// rotateLogs handles log rotation and compression
func rotateLogs(folderPath string) error {
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(folderPath, file.Name())
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			continue
		}

		// Check if file needs rotation based on size or age
		if fileInfo.Size() > config.MaxLogSize || time.Since(fileInfo.ModTime()) > config.MaxLogAge {
			// Compress old log file
			if err := compressLog(filePath); err != nil {
				logger.Printf("Failed to compress log file %s: %v", filePath, err)
				continue
			}

			// Remove original file after compression
			if err := os.Remove(filePath); err != nil {
				logger.Printf("Failed to remove original log file %s: %v", filePath, err)
			}
		}
	}

	return nil
}

// compressLog compresses a log file using gzip
func compressLog(filePath string) error {
	input, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(filePath + ".gz")
	if err != nil {
		return err
	}
	defer output.Close()

	gzWriter := gzip.NewWriter(output)
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, input)
	return err
}

var startTime = time.Now()

func main() {
	// Load configuration
	config = loadConfig()

	// Setup logger
	if err := setupLogger(); err != nil {
		log.Fatalf("Failed to setup logger: %v", err)
	}

	// Setup routes
	http.HandleFunc("/", basicAuth(indexHandler))
	http.HandleFunc("/health", healthCheckHandler)
	http.HandleFunc("/metrics", basicAuth(metricsHandler))

	// Create server with timeouts
	server := &http.Server{
		Addr:           ":" + config.Port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	logger.Printf("Server starting on port %s", config.Port)
	if err := server.ListenAndServe(); err != nil {
		logger.Fatalf("Server failed: %v", err)
	}
}
