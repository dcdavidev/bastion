package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// ListAuditLogs returns the latest audit events.
func (h *Handler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 50
	}

	logs, err := h.DB.GetAuditLogs(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
