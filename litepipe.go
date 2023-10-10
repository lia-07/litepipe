package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

// Baked in values for now
const (
	port          = 3001
	webhookSecret = ""
)

func main() {
	router := httprouter.New()

	router.POST("/", HandleWebhook)

	fmt.Println("LitePipe version 0.0.1")
	fmt.Printf("PID: %d\n", os.Getpid())
	fmt.Printf("Listening on port %d", port)
	addr := fmt.Sprintf(":%d", port)
	http.ListenAndServe(addr, router)

}

func HandleWebhook(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	signature := r.Header.Get("X-Hub-Signature")
	if !verifyWebhookSignature(signature, body) {
		fmt.Println("Invalid Webhook")
		http.Error(w, "Invalid Webhook", http.StatusForbidden)
		return
	}

	fmt.Println("Valid Webhook")

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
