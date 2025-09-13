package handlers

import (
	"net/http"
	"server/database"

	"golang.org/x/exp/slog"
)

type GetLogsResponse []database.Log

func (h *Handler) GetLogs(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}

	dbh, done := h.Database.GetHandler()
	defer done()
	if dbh == nil {
		slog.Warn("No database connection is available.")
		WriteFailureResponse("Operation cannot be performed. Database connection is not established properly.", w)
		return
	}
	logs, err := dbh.GetLogs()
	if err != nil {
		slog.Warn("cannot get logs", slog.String("error", err.Error()))
		WriteFailureResponse("cannot get logs", w)
		return
	}

	WriteSuccessResponse("logs", &logs, w)
}
