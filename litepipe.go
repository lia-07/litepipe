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

const (
	port          = 3001
	webhookSecret = ""
)

func main() {

	http.HandleFunc("/", HandleWebhook)

	fmt.Println("LitePipe version 0.0.1")
	fmt.Printf("PID: %d\n", os.Getpid())
	fmt.Printf("Listening on port %d\n", port)
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
	mac := hmac.New(sha1.New, []byte(webhookSecret))
	_, err := mac.Write(payload)
	if err != nil {
		return false
	}

	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	expectedSignature := "sha1=" + expectedMAC

	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}
