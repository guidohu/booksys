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

	logs, err := h.GetDB().GetLogs()
	if err != nil {
		slog.Warn("cannot get logs", slog.String("error", err.Error()))
		WriteFailureResponse("cannot get logs", w)
		return
	}

	WriteSuccessResponse("logs", &logs, w)
}
