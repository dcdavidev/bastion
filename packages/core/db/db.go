package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dcdavidev/bastion/packages/core/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq" // Required for golang-migrate postgres driver
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Database defines the interface for database operations.
type Database interface {
	Close()
	Ping(ctx context.Context) error
	RunMigrations() error
	GetMigrationStatus() (uint, bool, error)

	// Auth & Users
	HasAdmin(ctx context.Context) (bool, error)
	CreateUser(ctx context.Context, username, email, hash, salt, role string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, string, string, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, string, string, error)
	GrantProjectAccess(ctx context.Context, userID, projectID uuid.UUID, wrappedKey string) error

	// Vault
	GetVaultConfig(ctx context.Context) (*VaultConfig, error)
	InitializeVault(ctx context.Context, wrappedMK, salt string) error

	// Clients
	CreateClient(ctx context.Context, name string) (*models.Client, error)
	GetClients(ctx context.Context) ([]models.Client, error)
	GetClientByID(ctx context.Context, id uuid.UUID) (*models.Client, error)
	DeleteClient(ctx context.Context, id uuid.UUID) error

	// Projects
	CreateProject(ctx context.Context, clientID uuid.UUID, name string, wrappedKey string) (*models.Project, error)
	GetProjectsByClient(ctx context.Context, clientID uuid.UUID) ([]models.Project, error)
	GetProjectByID(ctx context.Context, id uuid.UUID) (*models.Project, error)
	DeleteProject(ctx context.Context, id uuid.UUID) error
	GetProjectKeyForUser(ctx context.Context, projectID, userID uuid.UUID, isAdmin bool) (string, error)

	// Secrets
	CreateSecret(ctx context.Context, projectID uuid.UUID, key string, value string) (*models.Secret, error)
	GetSecretsByProject(ctx context.Context, projectID uuid.UUID) ([]models.Secret, error)
	GetSecretHistory(ctx context.Context, projectID uuid.UUID, key string) ([]models.Secret, error)

	// Audit
	LogEvent(ctx context.Context, action, targetType string, targetID uuid.UUID, metadata map[string]interface{}) error
	GetAuditLogs(ctx context.Context, filter AuditFilter) ([]models.AuditLog, error)
}

// DB wrap the pgxpool.Pool to provide database access.
type DB struct {
	Pool *pgxpool.Pool
}

// NewConnection initializes a new PostgreSQL connection pool.
func NewConnection() (*DB, error) {
	connStr := os.Getenv("BASTION_DATABASE_URL")
	if connStr == "" {
		connStr = os.Getenv("DATABASE_URL")
	}
	if connStr == "" {
		return nil, fmt.Errorf("BASTION_DATABASE_URL or DATABASE_URL environment variable is not set")
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse DATABASE_URL: %w", err)
	}

	// Default pool settings
	config.MaxConns = 20
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Ping the database to ensure connection is working
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL")
	return &DB{Pool: pool}, nil
}

// Close closes the connection pool.
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// Ping checks if the database connection is alive.
func (db *DB) Ping(ctx context.Context) error {
	if db.Pool == nil {
		return fmt.Errorf("database pool is nil")
	}
	return db.Pool.Ping(ctx)
}

// RunMigrations applies all pending migrations.
func (db *DB) RunMigrations() error {
	connStr := os.Getenv("BASTION_DATABASE_URL")
	if connStr == "" {
		connStr = os.Getenv("DATABASE_URL")
	}
	
	// We need a standard sql.DB for golang-migrate
	importDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("could not open sql.DB for migrations: %w", err)
	}
	defer importDB.Close()

	driver, err := postgres.WithInstance(importDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %w", err)
	}

	migrationPath := db.getMigrationPath()

	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("could not create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Database migrations applied successfully")
	return nil
}

func (db *DB) getMigrationPath() string {
	migrationPath := "file://internal/db/migrations"
	if _, err := os.Stat("packages/core/db/migrations"); err == nil {
		migrationPath = "file://packages/core/db/migrations"
	} else if _, err := os.Stat("../../packages/core/db/migrations"); err == nil {
		migrationPath = "file://../../packages/core/db/migrations"
	}
	return migrationPath
}

// GetMigrationStatus returns the current migration version and whether there are pending migrations.
func (db *DB) GetMigrationStatus() (uint, bool, error) {
	connStr := os.Getenv("BASTION_DATABASE_URL")
	if connStr == "" {
		connStr = os.Getenv("DATABASE_URL")
	}
	importDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return 0, false, err
	}
	defer importDB.Close()

	driver, err := postgres.WithInstance(importDB, &postgres.Config{})
	if err != nil {
		return 0, false, err
	}

	migrationPath := db.getMigrationPath()
	m, err := migrate.NewWithDatabaseInstance(migrationPath, "postgres", driver)
	if err != nil {
		return 0, false, err
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return 0, false, err
	}

	if dirty {
		return version, true, fmt.Errorf("database is in a dirty state at version %d", version)
	}

	// Check if there are any migrations newer than 'version'
	err = m.Up()
	if err == migrate.ErrNoChange {
		return version, false, nil
	} else if err == nil {
		// If Up() succeeded, it means there WERE pending migrations. 
		// But Up() also applied them. For a "status" check we might want something non-destructive.
		// However, golang-migrate doesn't have a simple "check" without applying or looking at the source.
		return version, true, nil
	}

	return version, true, err
}
