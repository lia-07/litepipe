package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {

	// read the flags
	configPath := flag.String("config", "config.json", "The path to the config file")
	flag.Parse()

	loadConfig(*configPath)

	http.HandleFunc("/", HandleWebhook)

	fmt.Println("\n\033[30;46m LitePipe \033[0m version 0.1.13")
	fmt.Printf("PID: %d\n", os.Getpid())
	fmt.Printf("Listening on port %d\n\n", config.Port)

	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)

}

var config struct {
	Port               int      `json:"port"`
	WebhookSecret      string   `json:"webhookSecret"`
	TriggerDirectories []string `json:"triggerDirectories"`
	Tasks              []string `json:"tasks"`
	TasksDirectory     string   `json:"tasksDirectory"`
}

func loadConfig(path string) {
	file, err := os.Open(path)
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

	validateConfig()
}

func validateConfig() {
	if config.Port == 0 {
		config.Port = 3001
	}

	if config.WebhookSecret == "" {
		fmt.Println("You need to provide your webhook secret")
		os.Exit(1)
	}

	if config.TriggerDirectories == nil || len(config.TriggerDirectories) == 0 {
		fmt.Println("No trigger directories specified, defaulting to *")
		config.TriggerDirectories = []string{"*"}
	}

	if config.Tasks == nil || len(config.Tasks) == 0 {
		fmt.Println("There needs to be at least one task")
		os.Exit(1)
	}

	if config.TasksDirectory == "" {
		fmt.Println("No directory for tasks to be executed in specified, defaulting to current directory")
	}
}

type commit struct {
	Added    []string `json:"added"`
	Modified []string `json:"modified"`
	Removed  []string `json:"removed"`
	ID       string   `json:"id"`
	Message  string   `json:"message"`
	Author   struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"author"`
	Timestamp string `json:"timestamp"`
}

type webhookBody struct {
	Commit commit `json:"head_commit"`
}

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	signature := r.Header.Get("X-Hub-Signature")
	if !validWebhookSignature(signature, body) {
		fmt.Println("Invalid Webhook")
		http.Error(w, "Invalid Webhook", http.StatusForbidden)
		return
	}

	var payload webhookBody
	err = json.Unmarshal(body, &payload)
	if err != nil {
		fmt.Println("Failed to parse JSON payload:", err)
		http.Error(w, "Failed to parse JSON payload", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	processWebhookPayload(payload)
}

func processWebhookPayload(payload webhookBody) {
	var commit commit = payload.Commit
	fmt.Println("---------------")
	fmt.Printf("\x1b[1mReceived webhook for commit:\x1b[0m \n%s \"%s\" \nat %s\n", commit.ID, commit.Message, time.Now().UTC().Format("2006-01-02 15:04:05 MST"))

	triggerChanged := false

	if len(commit.Added) > 0 {
		for _, file := range commit.Added {
			if pathsMatch(file) {
				triggerChanged = true
			}
		}
	}

	if len(commit.Modified) > 0 {
		for _, file := range commit.Modified {
			if pathsMatch(file) {
				triggerChanged = true
			}
		}
	}

	if len(commit.Removed) > 0 {
		for _, file := range commit.Removed {
			if pathsMatch(file) {
				triggerChanged = true
			}
		}
	}

	if triggerChanged {
		fmt.Printf("\n\x1b[1mOne or more changes in trigger directory/ies, running tasks...\x1b[0m\n")
		for i, task := range config.Tasks {
			fmt.Printf("\n\x1b[1m(%d/%d): %s\n\x1b[0m", i+1, len(config.Tasks), task)

			start := time.Now()

			cmd := exec.Command("bash", "-c", task)
			cmd.Dir = config.TasksDirectory

			// output command outputs - maybe add a flag to show?
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				fmt.Printf("\x1b[101;30m Task failed: \x1b[0m %s", err)
			} else {
				fmt.Print("\x1b[1mTask completed\x1b[0m")
			}

			elapsed := time.Since(start)
			fmt.Printf(" in %s\n", elapsed)
		}
	} else {
		fmt.Printf("\n\x1b[1mNo changes in trigger directory/ies\x1b[0m\n")

	}

	fmt.Print("\nStill listening...\n\n")

}

func validWebhookSignature(signature string, payload []byte) bool {
	mac := hmac.New(sha1.New, []byte(config.WebhookSecret))
	_, err := mac.Write(payload)
	if err != nil {
		return false
	}

	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	expectedSignature := "sha1=" + expectedMAC

	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

func pathsMatch(path string) bool {
	for _, pattern := range config.TriggerDirectories {
		matched, err := filepath.Match(pattern, path)
		if err != nil {
			fmt.Println("Error:", err)
			return false
		}
		if matched {
			return true
		}
	}
	return false
}
