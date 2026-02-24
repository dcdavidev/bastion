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

// GetAuditLogs returns the latest audit events.
func (db *DB) GetAuditLogs(ctx context.Context, limit int) ([]models.AuditLog, error) {
	query := `
		SELECT id, action, target_type, target_id, metadata, created_at
		FROM audit_logs
		ORDER BY created_at DESC
		LIMIT $1
	`
	
	rows, err := db.Pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var l models.AuditLog
		var metaRaw []byte
		if err := rows.Scan(&l.ID, &l.Action, &l.TargetType, &l.TargetID, &metaRaw, &l.CreatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal(metaRaw, &l.Metadata)
		logs = append(logs, l)
	}

	return logs, nil
}
