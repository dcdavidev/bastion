package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/dcdavidev/bastion/packages/core/db"
)

// ListAuditLogs returns filtered audit events.
func (h *Handler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 50
	}

	filter := db.AuditFilter{
		Action:     query.Get("action"),
		TargetType: query.Get("target_type"),
		Limit:      limit,
	}

	if from := query.Get("from"); from != "" {
		if t, err := time.Parse(time.RFC3339, from); err == nil {
			filter.FromDate = &t
		}
	}

	if to := query.Get("to"); to != "" {
		if t, err := time.Parse(time.RFC3339, to); err == nil {
			filter.ToDate = &t
		}
	}

	logs, err := h.DB.GetAuditLogs(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
