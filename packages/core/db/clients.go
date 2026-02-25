package db

import (
	"context"
	"fmt"

	"github.com/dcdavidev/bastion/packages/core/models"
	"github.com/google/uuid"
)

// CreateClient inserts a new client into the database.
func (db *DB) CreateClient(ctx context.Context, name string) (*models.Client, error) {
	query := `
		INSERT INTO clients (name)
		VALUES ($1)
		RETURNING id, name, created_at, updated_at
	`

	client := &models.Client{}
	err := db.Pool.QueryRow(ctx, query, name).Scan(
		&client.ID,
		&client.Name,
		&client.CreatedAt,
		&client.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return client, nil
}

// GetClients returns a list of all clients.
func (db *DB) GetClients(ctx context.Context) ([]models.Client, error) {
	query := `SELECT id, name, created_at, updated_at FROM clients ORDER BY name ASC`

	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list clients: %w", err)
	}
	defer rows.Close()

	var clients []models.Client
	for rows.Next() {
		var c models.Client
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan client: %w", err)
		}
		clients = append(clients, c)
	}

	return clients, nil
}

// GetClientByID returns a single client by its ID.
func (db *DB) GetClientByID(ctx context.Context, id uuid.UUID) (*models.Client, error) {
	query := `SELECT id, name, created_at, updated_at FROM clients WHERE id = $1`

	client := &models.Client{}
	err := db.Pool.QueryRow(ctx, query, id).Scan(
		&client.ID,
		&client.Name,
		&client.CreatedAt,
		&client.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	return client, nil
}

// DeleteClient removes a client from the database.
func (db *DB) DeleteClient(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM clients WHERE id = $1`
	_, err := db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}
	return nil
}
