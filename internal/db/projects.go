package db

import (
	"context"
	"fmt"

	"github.com/dcdavidev/bastion/internal/models"
	"github.com/google/uuid"
)

// CreateProject inserts a new project for a specific client.
func (db *DB) CreateProject(ctx context.Context, clientID uuid.UUID, name string, wrappedKey string) (*models.Project, error) {
	query := `
		INSERT INTO projects (client_id, name, wrapped_data_key)
		VALUES ($1, $2, $3)
		RETURNING id, client_id, name, wrapped_data_key, created_at, updated_at
	`

	project := &models.Project{}
	err := db.Pool.QueryRow(ctx, query, clientID, name, wrappedKey).Scan(
		&project.ID,
		&project.ClientID,
		&project.Name,
		&project.WrappedDataKey,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return project, nil
}

// GetProjectsByClient returns all projects belonging to a specific client.
func (db *DB) GetProjectsByClient(ctx context.Context, clientID uuid.UUID) ([]models.Project, error) {
	query := `
		SELECT id, client_id, name, wrapped_data_key, created_at, updated_at 
		FROM projects 
		WHERE client_id = $1 
		ORDER BY name ASC
	`

	rows, err := db.Pool.Query(ctx, query, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.ClientID, &p.Name, &p.WrappedDataKey, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, p)
	}

	return projects, nil
}

// GetProjectByID returns a single project by its ID.
func (db *DB) GetProjectByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	query := `SELECT id, client_id, name, wrapped_data_key, created_at, updated_at FROM projects WHERE id = $1`

	project := &models.Project{}
	err := db.Pool.QueryRow(ctx, query, id).Scan(
		&project.ID,
		&project.ClientID,
		&project.Name,
		&project.WrappedDataKey,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return project, nil
}

// DeleteProject removes a project from the database.
func (db *DB) DeleteProject(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	return nil
}
