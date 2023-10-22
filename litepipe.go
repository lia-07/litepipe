package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Config struct {
	Port               int      `json:"port"`
	WebhookSecret      string   `json:"webhookSecret"`
	TriggerDirectories []string `json:"triggerDirectories"`
	Tasks              []string `json:"tasks"`
}

var config Config

func main() {
	// try to read the config file
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Printf("Error opening config file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	reader := io.Reader(file)
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&config); err != nil {
		fmt.Printf("Error parsing config file: %v\n", err)
		os.Exit(1)
	}

	// defaults and validations
	if config.Port == 0 {
		config.Port = 3001
	}
	if config.WebhookSecret == "" {
		fmt.Println("You need to include your webhook secret")
		os.Exit(1)
	}
	if config.TriggerDirectories == nil || len(config.TriggerDirectories) == 0 {
		fmt.Println("There needs to be at least one trigger directory")
		os.Exit(1)
	}
	if config.Tasks == nil || len(config.Tasks) == 0 {
		fmt.Println("There needs to be at least one task")
		os.Exit(1)
	}

	http.HandleFunc("/", HandleWebhook)

	fmt.Println("LitePipe version 0.0.1")
	fmt.Printf("PID: %d\n", os.Getpid())
	fmt.Printf("Listening on port %d\n", config.Port)
	http.ListenAndServe(":3001", nil)
}

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// verify webhook signature
	signature := r.Header.Get("X-Hub-Signature")
	if !verifyWebhookSignature(signature, body) {
		fmt.Println("Invalid Webhook")
		http.Error(w, "Invalid Webhook", http.StatusForbidden)
		return
	}

	// try and parse webhook payload
	var payload map[string]interface{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		fmt.Println("Failed to parse JSON payload:", err)
		http.Error(w, "Failed to parse JSON payload", http.StatusInternalServerError)
		return
	}

	payloadJSON, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		fmt.Println("Failed to marshal payload to JSON:", err)
	} else {
		fmt.Printf("Payload as JSON:\n%s\n", payloadJSON)
	}

	w.WriteHeader(http.StatusOK)
}

func verifyWebhookSignature(signature string, payload []byte) bool {
	mac := hmac.New(sha1.New, []byte(config.WebhookSecret))
	_, err := mac.Write(payload)
	if err != nil {
		return false
	}

	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	expectedSignature := "sha1=" + expectedMAC

	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}
