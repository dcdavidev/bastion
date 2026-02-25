package db

import (
	"context"
	"fmt"

	"github.com/dcdavidev/bastion/packages/core/models"
	"github.com/google/uuid"
)

// CreateUser inserts a new user into the database.
func (db *DB) CreateUser(ctx context.Context, username, hash, salt, role string) (*models.User, error) {
	query := `
		INSERT INTO users (username, password_hash, salt, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, role, created_at, updated_at
	`

	user := &models.User{}
	err := db.Pool.QueryRow(ctx, query, username, hash, salt, role).Scan(
		&user.ID,
		&user.Username,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GrantProjectAccess links a user to a project with a specific wrapped data key.
func (db *DB) GrantProjectAccess(ctx context.Context, userID, projectID uuid.UUID, wrappedKey string) error {
	query := `
		INSERT INTO user_project_access (user_id, project_id, wrapped_data_key)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, project_id) DO UPDATE SET wrapped_data_key = EXCLUDED.wrapped_data_key
	`
	_, err := db.Pool.Exec(ctx, query, userID, projectID, wrappedKey)
	return err
}

// GetUserByUsername retrieves a user for authentication.
func (db *DB) GetUserByUsername(ctx context.Context, username string) (*models.User, string, string, error) {
	query := `SELECT id, username, password_hash, salt, role, created_at, updated_at FROM users WHERE username = $1`
	
	user := &models.User{}
	var hash, salt string
	err := db.Pool.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&hash,
		&salt,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		return nil, "", "", err
	}
	
	return user, hash, salt, nil
}
