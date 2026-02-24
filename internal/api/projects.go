package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CreateProjectRequest struct {
	ClientID       uuid.UUID `json:"client_id"`
	Name           string    `json:"name"`
	WrappedDataKey string    `json:"wrapped_data_key"`
}

// CreateProject handles the creation of a new project for a client.
func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ClientID == uuid.Nil || req.Name == "" || req.WrappedDataKey == "" {
		http.Error(w, "client_id, name and wrapped_data_key are required", http.StatusBadRequest)
		return
	}

	project, err := h.DB.CreateProject(r.Context(), req.ClientID, req.Name, req.WrappedDataKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

// ListProjectsByClient returns all projects for a specific client.
func (h *Handler) ListProjectsByClient(w http.ResponseWriter, r *http.Request) {
	clientIDStr := r.URL.Query().Get("client_id")
	if clientIDStr == "" {
		http.Error(w, "client_id query parameter is required", http.StatusBadRequest)
		return
	}

	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		http.Error(w, "Invalid client_id", http.StatusBadRequest)
		return
	}

	projects, err := h.DB.GetProjectsByClient(r.Context(), clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// DeleteProject removes a project by ID.
func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	if err := h.DB.DeleteProject(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
