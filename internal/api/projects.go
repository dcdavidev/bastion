package api

import (
	"encoding/json"
	"net/http"

	"github.com/dcdavidev/bastion/internal/auth"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
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

// GetProject returns a single project by ID.
func (h *Handler) GetProject(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	project, err := h.DB.GetProjectByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

// GetProjectKey returns the wrapped data key for the authenticated user.
func (h *Handler) GetProjectKey(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	// Extract claims from context (added by JWTMiddleware)
	claims, ok := r.Context().Value(auth.AdminContextKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userIDStr, _ := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)
	isAdmin, _ := claims["admin"].(bool)

	wrappedKey, err := h.DB.GetProjectKeyForUser(r.Context(), projectID, userID, isAdmin)
	if err != nil {
		http.Error(w, "Access denied or project not found", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"wrapped_data_key": wrappedKey})
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
