package handlers

import (
	"log/slog"
	"net/http"

	"server/database"
)

// GetLogsResponse is the payload returned by GetLogs.
type GetLogsResponse []database.Log

// GetLogs serves the activity log, with amounts rendered in the configured
// currency.
func (h *Handler) GetLogs(w http.ResponseWriter, r *http.Request) {
	hCtx := GetHandlerContext(r)
	dbh := hCtx.Database
	currency, _ := h.config.GetString("currency")
	if currency == "" {
		slog.Warn("Cannot get currency from configuration")
		WriteFailureResponse("cannot get currency", w)
		return
	}
	logs, err := dbh.GetLogs(currency)
	if err != nil {
		slog.Warn("cannot get logs", slog.Any("error", err))
		WriteFailureResponse("cannot get logs", w)
		return
	}

	WriteSuccessResponse("logs", &logs, w)
}
