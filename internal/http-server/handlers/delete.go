package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	"net/http"
	resp "short-url/internal/http-server/dto"
	"short-url/internal/lib/logger/sl"
)

func DeleteURL(log *slog.Logger, manager URLManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = getHandlerLogger(log, "handlers.DeleteURL", r)
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			render.JSON(w, r, resp.Error("Неверная ссылка"))
			return
		}
		err := manager.DeleteUrl(alias)
		if err != nil {
			log.Error(
				"Can't delete URL for alias",
				slog.String("alias", alias),
				sl.Err(err),
			)
			render.JSON(w, r, resp.Error("Не удалось удалить URL"))
			return
		}
		render.Status(r, http.StatusNoContent)
	}
}
