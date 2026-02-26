package models

import (
	"time"

	"github.com/google/uuid"
)

// Client represents a customer who owns multiple projects.
type Client struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Project represents a group of secrets for a specific client.
type Project struct {
	ID             uuid.UUID `json:"id"`
	ClientID       uuid.UUID `json:"client_id"`
	Name           string    `json:"name"`
	WrappedDataKey string    `json:"wrapped_data_key,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Secret represents an encrypted secret stored in the vault.
type Secret struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"` // This will be the encrypted payload
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// User represents an administrator or a collaborator.
type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email,omitempty"`
	PasswordHash string    `json:"-"`
	Salt         string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// AuditLog tracks sensitive operations in the vault.
type AuditLog struct {
	ID         uuid.UUID              `json:"id"`
	Action     string                 `json:"action"`
	TargetType string                 `json:"target_type"`
	TargetID   uuid.UUID              `json:"target_id,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
}
