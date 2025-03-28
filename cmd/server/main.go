package main

import (
	"log"
	"net/http"

	"github.com/ianferrier777/mini-feature-flag-service/internal/api"
	"github.com/ianferrier777/mini-feature-flag-service/internal/flags"
)

func main() {
	log.Println("🚀 Starting Mini Feature Flag Service on :8080")

	err := flags.LoadFromFile()
	if err != nil {
		log.Fatalf("Failed to load flags from file: %v", err)
	}

	flags.InitFlags()

	// Initialize router
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Register feature flag routes from internal/api
	api.RegisterRoutes(mux)

	// Start HTTP server
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
