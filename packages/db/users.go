package db

import (
	"context"
	"fmt"

	"github.com/dcdavidev/bastion/packages/models"
	"github.com/google/uuid"
)

// HasAdmin checks if there is at least one admin user in the database.
func (db *DB) HasAdmin(ctx context.Context) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE role = 'ADMIN')`
	var exists bool
	err := db.Pool.QueryRow(ctx, query).Scan(&exists)
	return exists, err
}

// CreateUser inserts a new user into the database.
func (db *DB) CreateUser(ctx context.Context, username, email, hash, salt, role string) (*models.User, error) {
	query := `
		INSERT INTO users (username, email, password_hash, salt, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, username, email, role, created_at, updated_at
	`

	user := &models.User{}
	err := db.Pool.QueryRow(ctx, query, username, email, hash, salt, role).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// UpdateUserPassword updates the password hash and salt for a user.
func (db *DB) UpdateUserPassword(ctx context.Context, userID uuid.UUID, hash, salt string) error {
	query := `UPDATE users SET password_hash = $1, salt = $2, updated_at = NOW() WHERE id = $3`
	_, err := db.Pool.Exec(ctx, query, hash, salt, userID)
	return err
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
	query := `SELECT id, username, email, password_hash, salt, role, created_at, updated_at FROM users WHERE username = $1`
	
	user := &models.User{}
	var hash, salt string
	err := db.Pool.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
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

// GetUserByEmail retrieves a user by email for authentication.
func (db *DB) GetUserByEmail(ctx context.Context, email string) (*models.User, string, string, error) {
	query := `SELECT id, username, email, password_hash, salt, role, created_at, updated_at FROM users WHERE email = $1`
	
	user := &models.User{}
	var hash, salt string
	err := db.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
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

// GetUserByID retrieves a user by their UUID.
func (db *DB) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `SELECT id, username, email, role, created_at, updated_at FROM users WHERE id = $1`
	user := &models.User{}
	err := db.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// AddWebAuthnCredential saves a new WebAuthn credential for a user.
func (db *DB) AddWebAuthnCredential(ctx context.Context, userID uuid.UUID, cred *models.WebAuthnCredential) error {
	query := `
		INSERT INTO webauthn_credentials (id, user_id, public_key, attestation_type, transport, sign_count, clone_warning)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := db.Pool.Exec(ctx, query, 
		cred.ID, 
		userID, 
		cred.PublicKey, 
		cred.AttestationType, 
		cred.Transport, 
		cred.SignCount, 
		cred.CloneWarning,
	)
	return err
}

// GetWebAuthnCredentials retrieves all WebAuthn credentials for a user.
func (db *DB) GetWebAuthnCredentials(ctx context.Context, userID uuid.UUID) ([]models.WebAuthnCredential, error) {
	query := `
		SELECT id, public_key, attestation_type, transport, sign_count, clone_warning, created_at, updated_at
		FROM webauthn_credentials
		WHERE user_id = $1
	`
	rows, err := db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var creds []models.WebAuthnCredential
	for rows.Next() {
		var cred models.WebAuthnCredential
		err := rows.Scan(
			&cred.ID,
			&cred.PublicKey,
			&cred.AttestationType,
			&cred.Transport,
			&cred.SignCount,
			&cred.CloneWarning,
			&cred.CreatedAt,
			&cred.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		creds = append(creds, cred)
	}
	return creds, nil
}

// UpdateWebAuthnCredential updates the sign count and other fields of a WebAuthn credential.
func (db *DB) UpdateWebAuthnCredential(ctx context.Context, cred *models.WebAuthnCredential) error {
	query := `
		UPDATE webauthn_credentials
		SET sign_count = $1, clone_warning = $2, updated_at = NOW()
		WHERE id = $3
	`
	_, err := db.Pool.Exec(ctx, query, cred.SignCount, cred.CloneWarning, cred.ID)
	return err
}
