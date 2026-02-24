package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type CreateSecretRequest struct {
	ProjectID uuid.UUID `json:"project_id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"` // Already encrypted
}

// CreateSecret handles the creation of a new secret version.
func (h *Handler) CreateSecret(w http.ResponseWriter, r *http.Request) {
	var req CreateSecretRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ProjectID == uuid.Nil || req.Key == "" || req.Value == "" {
		http.Error(w, "project_id, key and value are required", http.StatusBadRequest)
		return
	}

	secret, err := h.DB.CreateSecret(r.Context(), req.ProjectID, req.Key, req.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(secret)
}

// ListSecretsByProject returns the latest versions of secrets for a project.
func (h *Handler) ListSecretsByProject(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	if projectIDStr == "" {
		http.Error(w, "project_id query parameter is required", http.StatusBadRequest)
		return
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project_id", http.StatusBadRequest)
		return
	}

	secrets, err := h.DB.GetSecretsByProject(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(secrets)
}

// GetSecretHistory returns all versions of a specific secret.
func (h *Handler) GetSecretHistory(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	key := r.URL.Query().Get("key")

	if projectIDStr == "" || key == "" {
		http.Error(w, "project_id and key query parameters are required", http.StatusBadRequest)
		return
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project_id", http.StatusBadRequest)
		return
	}

	history, err := h.DB.GetSecretHistory(r.Context(), projectID, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}
