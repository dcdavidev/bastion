package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type CreateCollaboratorRequest struct {
	Username       string    `json:"username"`
	Email          string    `json:"email,omitempty"`
	PasswordHash   string    `json:"password_hash"`
	Salt           string    `json:"salt"`
	ProjectID      uuid.UUID `json:"project_id"`
	WrappedDataKey string    `json:"wrapped_data_key"`
}

// CreateCollaborator handles the creation of a new collaborator.
func (h *Handler) CreateCollaborator(w http.ResponseWriter, r *http.Request) {
	// Check if the current user is an admin (later integration with JWT claims)
	// For now, we assume the middleware protects it or we check manually
	
	var req CreateCollaboratorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 1. Create User
	user, err := h.DB.CreateUser(r.Context(), req.Username, req.Email, req.PasswordHash, req.Salt, "COLLABORATOR")
	if err != nil {
		http.Error(w, "Could not create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 2. Grant Project Access
	err = h.DB.GrantProjectAccess(r.Context(), user.ID, req.ProjectID, req.WrappedDataKey)
	if err != nil {
		http.Error(w, "Could not grant project access: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

	// Log audit event
	h.DB.LogEvent(r.Context(), "CREATE_COLLABORATOR", "USER", user.ID, map[string]interface{}{
		"username":   user.Username,
		"project_id": req.ProjectID,
		"ip":         r.RemoteAddr,
	})
}
