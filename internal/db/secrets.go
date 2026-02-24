package db

import (
	"context"
	"fmt"

	"github.com/dcdavidev/bastion/internal/models"
	"github.com/google/uuid"
)

// CreateSecret inserts a new encrypted secret for a project.
func (db *DB) CreateSecret(ctx context.Context, projectID uuid.UUID, key string, value string) (*models.Secret, error) {
	query := `
		INSERT INTO secrets (project_id, key, value)
		VALUES ($1, $2, $3)
		RETURNING id, project_id, key, value, version, created_at, updated_at
	`

	secret := &models.Secret{}
	err := db.Pool.QueryRow(ctx, query, projectID, key, value).Scan(
		&secret.ID,
		&secret.ProjectID,
		&secret.Key,
		&secret.Value,
		&secret.Version,
		&secret.CreatedAt,
		&secret.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create secret: %w", err)
	}

	return secret, nil
}

// GetSecretsByProject returns all the latest secrets for a specific project.
func (db *DB) GetSecretsByProject(ctx context.Context, projectID uuid.UUID) ([]models.Secret, error) {
	// Query to get the latest version of each secret key in the project
	query := `
		SELECT DISTINCT ON (key) id, project_id, key, value, version, created_at, updated_at
		FROM secrets
		WHERE project_id = $1
		ORDER BY key, version DESC
	`

	rows, err := db.Pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}
	defer rows.Close()

	var secrets []models.Secret
	for rows.Next() {
		var s models.Secret
		if err := rows.Scan(&s.ID, &s.ProjectID, &s.Key, &s.Value, &s.Version, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan secret: %w", err)
		}
		secrets = append(secrets, s)
	}

	return secrets, nil
}

// GetSecretHistory returns all versions of a specific secret.
func (db *DB) GetSecretHistory(ctx context.Context, projectID uuid.UUID, key string) ([]models.Secret, error) {
	query := `
		SELECT id, project_id, key, value, version, created_at, updated_at
		FROM secrets
		WHERE project_id = $1 AND key = $2
		ORDER BY version DESC
	`

	rows, err := db.Pool.Query(ctx, query, projectID, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret history: %w", err)
	}
	defer rows.Close()

	var history []models.Secret
	for rows.Next() {
		var s models.Secret
		if err := rows.Scan(&s.ID, &s.ProjectID, &s.Key, &s.Value, &s.Version, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan secret: %w", err)
		}
		history = append(history, s)
	}

	return history, nil
}
