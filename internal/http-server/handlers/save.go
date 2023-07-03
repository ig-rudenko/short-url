package handlers

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
	"net/http"
	resp "short-url/internal/http-server/dto"
	"short-url/internal/lib/logger/sl"
	"short-url/internal/lib/random"
	validatorHandler "short-url/internal/lib/validator"
	"short-url/internal/storage"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

func SaveURL(log *slog.Logger, urlManager URLManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = getHandlerLogger(log, "handlers.SaveURL", r)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		if err := validator.New().Struct(req); err != nil {
			validatorError := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, resp.Error(validatorHandler.ValidationError(validatorError)))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(8)
		}

		err = urlManager.SaveUrl(req.URL, alias)

		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				log.Info("URL под таким именем уже существует", slog.String("url", req.URL))
				render.JSON(w, r, resp.Error("URL под таким именем уже существует"))
				return
			}
			log.Error("failed to add url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to add url"))
			return
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
