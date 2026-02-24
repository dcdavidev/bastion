package db

import (
	"context"
	"fmt"

	"github.com/dcdavidev/bastion/internal/models"
	"github.com/google/uuid"
)

// VaultConfig represents the global vault settings
type VaultConfig struct {
	WrappedMasterKey string
	MasterKeySalt    string
}

// GetVaultConfig retrieves the global vault configuration.
func (db *DB) GetVaultConfig(ctx context.Context) (*VaultConfig, error) {
	query := `SELECT wrapped_master_key, master_key_salt FROM vault_config LIMIT 1`
	config := &VaultConfig{}
	err := db.Pool.QueryRow(ctx, query).Scan(&config.WrappedMasterKey, &config.MasterKeySalt)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// InitializeVault sets up the master key for the first time.
func (db *DB) InitializeVault(ctx context.Context, wrappedMK, salt string) error {
	query := `INSERT INTO vault_config (wrapped_master_key, master_key_salt) VALUES ($1, $2)`
	_, err := db.Pool.Exec(ctx, query, wrappedMK, salt)
	return err
}
