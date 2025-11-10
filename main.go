package main

import (
	"fmt"
	"log"
	"net/http" // This is used for serving the website and handling webhooks
	"os"

	"github.com/go-git/go-git/v5"
	// Rename to avoid conflict with net/http
)

const repoDir = "/app/site-repo"
const publicDir = "/app/site-repo/public"

// cloneRepo clones the HTML repo from the given URL if it's not already cloned
func cloneRepo(repoURL string) error {
	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		_, err := git.PlainClone(repoDir, false, &git.CloneOptions{
			URL:      repoURL,
			Progress: os.Stdout,
		})
		if err != nil {
			return fmt.Errorf("failed to clone repo: %w", err)
		}
		log.Println("Repository cloned successfully")
	} else {
		log.Println("Repository already exists")
	}
	return nil
}

// pullRepo pulls the latest changes from the repo
func pullRepo() error {
	r, err := git.PlainOpen(repoDir)
	if err != nil {
		return fmt.Errorf("failed to open repo: %w", err)
	}

	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	log.Println("Pulling latest changes...")
	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("failed to pull latest changes: %w", err)
	}
	log.Println("Repo updated successfully")

	return nil
}

// webhookHandler triggers a git pull when the webhook is received
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	log.Println("Webhook received, pulling updates...")
	if err := pullRepo(); err != nil {
		log.Printf("Failed to update the repo: %v\n", err)
		http.Error(w, "failed to update repo", http.StatusInternalServerError)
		return
	}

	if _, err := fmt.Fprintln(w, "Repo updated"); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

// servePublicDir starts serving the public directory as a static file server
func servePublicDir() {
	if _, err := os.Stat(publicDir); os.IsNotExist(err) {
		log.Fatalf("Error: public directory not found in %s", publicDir)
	}

	fs := http.FileServer(http.Dir(publicDir))
	http.Handle("/", fs)
}

func main() {
	// Get repo URL from environment variable
	repoURL := os.Getenv("REPO_URL")
	if repoURL == "" {
		log.Fatal("REPO_URL environment variable is not set")
	}

	// Clone repo
	if err := cloneRepo(repoURL); err != nil {
		log.Fatalf("Error cloning repo: %v", err)
	}

	// Serve the public directory
	servePublicDir()

	// Webhook listener to trigger git pull
	http.HandleFunc("/webhook", webhookHandler)

	// Start the server
	port := ":8080"
	log.Printf("Serving static site on %s and waiting for webhooks on /webhook", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
