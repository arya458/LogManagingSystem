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

type LogEntry struct {
	EncryptedData string `json:"encrypted_data"`
}

var (
	username      = "myuser"
	password      = "mypassword"
	encryptionKey = "your32bytekeyhere!"
)

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

func deriveKey(password string, salt []byte) []byte {
	combined := []byte(password)
	combined = append(combined, salt...)
	hash := sha256.Sum256(combined)
	return hash[:32]
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var entry LogEntry
	err = json.Unmarshal(body, &entry)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if entry.EncryptedData == "" {
		http.Error(w, "Missing encrypted data", http.StatusBadRequest)
		return
	}

	decrypted, err := decrypt(entry.EncryptedData, encryptionKey)
	if err != nil {
		http.Error(w, "Decryption failed", http.StatusInternalServerError)
		return
	}

	var logData map[string]string
	err = json.Unmarshal([]byte(decrypted), &logData)
	if err != nil {
		http.Error(w, "Decrypted data is not valid JSON", http.StatusBadRequest)
		return
	}

	imel, okImel := logData["imel"]
	errorMessage, okError := logData["error"]

	if !okImel || !okError {
		http.Error(w, "Missing imel or error in decrypted data", http.StatusBadRequest)
		return
	}

	folderPath := filepath.Join("logs", imel)
	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		http.Error(w, "Failed to create folder", http.StatusInternalServerError)
		return
	}
	currentTime := time.Now()
	timeString := currentTime.Format("2006-01-02_15-04-05")
	filePath := filepath.Join(folderPath, timeString+".log")
	file, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to create log file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("IMEL: %s\nTime: %s\nError: %s\n", imel, timeString, errorMessage))
	if err != nil {
		http.Error(w, "Failed to write to log file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Log entry saved successfully")
}

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

func main() {
	http.HandleFunc("/", basicAuth(indexHandler))
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
		log.Printf("Defaulting to port %s", port)
	}

	encodedCreds := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	jsonData := map[string]string{"imel": "123456789012345", "error": "Test error"}
	jsonBytes, _ := json.Marshal(jsonData)
	derivedKey := deriveKey(encryptionKey, make([]byte, 16))
	encrypted, _ := encrypt(string(jsonBytes), derivedKey)

	curlCommand := fmt.Sprintf("curl -X POST -H \"Authorization: Basic %s\" -H \"Content-Type: application/json\" -d '{\"encrypted_data\":\"%s\"}' http://localhost:%s/", encodedCreds, encrypted, port)

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://localhost:%s in the browser", port)
	log.Printf("Use this curl command to test the endpoint:\n%s", curlCommand)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
