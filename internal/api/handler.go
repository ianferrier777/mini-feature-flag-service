package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ianferrier777/mini-feature-flag-service/internal/flags"
)

const adminToken = "super-secret-token"

// Response format for flag checks
type FlagResponse struct {
	Enabled bool   `json:"enabled"`
	Reason  string `json:"reason"`
}

// RegisterRoutes registers all HTTP routes for the API
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/flags", handleAllFlags)
	mux.HandleFunc("/flags/", handleFlagRoutes)
}

func isAuthorized(r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	return auth == "Bearer "+adminToken
}

func handleFlagRoutes(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		http.Error(w, "Missing flag name", http.StatusBadRequest)
		return
	}
	flagName := parts[2]

	switch r.Method {
		case http.MethodGet:
			handleGetFlag(w, r, flagName)
		case http.MethodPost:
			handleCreateFlag(w, r, flagName)
		case http.MethodPut:
			handleUpdateFlag(w, r, flagName)
		case http.MethodDelete:
			handleDeleteFlag(w, r, flagName)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGetFlag processes a request to check a flag status
func handleGetFlag(w http.ResponseWriter, r *http.Request, flagName string) {
	userId := r.URL.Query().Get("userId")
	region := r.URL.Query().Get("region")

	enabled, reason, found := flags.EvaluateFlag(flagName, userId, region)
	if !found {
		http.Error(w, "Flag not found", http.StatusNotFound)
		return
	}

	response := FlagResponse{
		Enabled: enabled,
		Reason:  reason,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleAllFlags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	allFlags := flags.GetAllFlags()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allFlags)
}

func handleCreateFlag(w http.ResponseWriter, r *http.Request, flagName string) {
	var body struct {
		Enabled       bool     `json:"enabled"`
		TargetUsers   []string `json:"targetUsers"`
		TargetRegions []string `json:"targetRegions"`
	}

	if !isAuthorized(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	flags.SetFlag(flagName, body.Enabled, body.TargetUsers, body.TargetRegions)

	err := flags.SaveToFile()
	if err != nil {
		http.Error(w, "Failed to save flag to file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(FlagResponse{
		Enabled: body.Enabled,
		Reason:  "Flag created",
	})
}

func handleUpdateFlag(w http.ResponseWriter, r *http.Request, flagName string) {
	var body struct {
		Enabled       bool     `json:"enabled"`
		TargetUsers   []string `json:"targetUsers"`
		TargetRegions []string `json:"targetRegions"`
	}

	if !isAuthorized(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	updated := flags.UpdateFlag(flagName, body.Enabled, body.TargetUsers, body.TargetRegions)
	if !updated {
		http.Error(w, "Flag not found. Use POST to create.", http.StatusNotFound)
		return
	}

	err := flags.SaveToFile()
	if err != nil {
		http.Error(w, "Failed to save flag to file", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(FlagResponse{
		Enabled: body.Enabled,
		Reason:  "Flag updated",
	})
}

func handleDeleteFlag(w http.ResponseWriter, r *http.Request, flagName string) {
	if !isAuthorized(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	deleted := flags.DeleteFlag(flagName)
	if !deleted {
		http.Error(w, "Flag not found", http.StatusNotFound)
		return
	}

	err := flags.SaveToFile()
	if err != nil {
		http.Error(w, "Failed to save flag changes to file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 - success, no body
}