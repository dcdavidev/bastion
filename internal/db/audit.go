package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dcdavidev/bastion/internal/models"
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

// AuditFilter defines the available filters for audit logs.
type AuditFilter struct {
	Action     string
	TargetType string
	FromDate   *time.Time
	ToDate     *time.Time
	Limit      int
}

// GetAuditLogs returns filtered audit events.
func (db *DB) GetAuditLogs(ctx context.Context, filter AuditFilter) ([]models.AuditLog, error) {
	query := `
		SELECT id, action, target_type, target_id, metadata, created_at
		FROM audit_logs
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	if filter.Action != "" {
		query += fmt.Sprintf(" AND action = $%d", argIdx)
		args = append(args, filter.Action)
		argIdx++
	}

	if filter.TargetType != "" {
		query += fmt.Sprintf(" AND target_type = $%d", argIdx)
		args = append(args, filter.TargetType)
		argIdx++
	}

	if filter.FromDate != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argIdx)
		args = append(args, *filter.FromDate)
		argIdx++
	}

	if filter.ToDate != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argIdx)
		args = append(args, *filter.ToDate)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d", argIdx)
	args = append(args, filter.Limit)

	rows, err := db.Pool.Query(ctx, query, args...)
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
