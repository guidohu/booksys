package handlers

import (
	"net/http"
	"server/database"

	"golang.org/x/exp/slog"
)

type GetLogsResponse []database.Log

func (h *Handler) GetLogs(w http.ResponseWriter, r *http.Request) {
	hCtx, err := GetHandlerContext(w, r)
	if err != nil {
		slog.Warn("Cannot get handler context", slog.String("error", err.Error()))
		return
	}
	dbh := hCtx.Database
	currency, _ := h.config.GetString("currency")
	if currency == "" {
		slog.Warn("Cannot get currency from configuration")
		WriteFailureResponse("cannot get currency", w)
		return
	}
	logs, err := dbh.GetLogs(currency)
	if err != nil {
		slog.Warn("cannot get logs", slog.String("error", err.Error()))
		WriteFailureResponse("cannot get logs", w)
		return
	}

	WriteSuccessResponse("logs", &logs, w)
}
