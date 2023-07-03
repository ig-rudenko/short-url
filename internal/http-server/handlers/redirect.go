package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	"net/http"
	resp "short-url/internal/http-server/dto"
	"short-url/internal/lib/logger/sl"
	"short-url/internal/storage"
)

func Redirect(log *slog.Logger, manager URLManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = getHandlerLogger(log, "handlers.Redirect", r)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			render.JSON(w, r, resp.Error("Неверная ссылка"))
			return
		}

		resURL, err := manager.GetUrl(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("URL не был найден", slog.String("alias", alias))
			}
			log.Error(
				"Can't get URL for alias",
				slog.String("alias", alias),
				sl.Err(err),
			)
			render.JSON(w, r, resp.Error("Неверный URL"))
			return
		}
		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
