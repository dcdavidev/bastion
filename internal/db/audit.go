package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// LogEvent records a sensitive action in the audit_logs table.
func (db *DB) LogEvent(ctx context.Context, action, targetType string, targetID uuid.UUID, metadata map[string]interface{}) error {
	query := `
		INSERT INTO audit_logs (action, target_type, target_id, metadata)
		VALUES ($1, $2, $3, $4)
	`
	
	metaJSON, err := json.Marshal(metadata)
	if err != nil {
		metaJSON = []byte("{}")
	}

	_, err = db.Pool.Exec(ctx, query, action, targetType, targetID, metaJSON)
	if err != nil {
		return fmt.Errorf("failed to log audit event: %w", err)
	}
	return nil
}
