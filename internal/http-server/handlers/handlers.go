package handlers

import (
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLManager
type URLManager interface {
	SaveUrl(url_ string, alias string) error
	GetUrl(alias string) (string, error)
	DeleteUrl(alias string) error
}

func getHandlerLogger(logger *slog.Logger, operation string, r *http.Request) *slog.Logger {
	return logger.With(
		slog.String("op", operation),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
}
