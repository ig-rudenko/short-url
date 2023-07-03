package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"short-url/internal/config"
	"short-url/internal/http-server/handlers"
	middlewareLogger "short-url/internal/http-server/middleware/logger"
	"short-url/internal/lib/logger/sl"
	"short-url/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.Load()

	log := setupLogger(cfg.Env)

	log.Debug(fmt.Sprintf("configs: %s", *cfg))

	storage, err := sqlite.New(cfg.DSN)
	if err != nil {
		log.Error("Ошибка инициализации БД", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middlewareLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", handlers.SaveURL(log, storage))
	router.Get("/{alias}", handlers.Redirect(log, storage))
	router.Delete("/{alias}", handlers.DeleteURL(log, storage))

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	log.Info("Запуск сервера", slog.String("address", cfg.Address))
	err = server.ListenAndServe()
	if err != nil {
		log.Error("Ошибка запуска сервера")
		os.Exit(1)
	}

}

func setupLogger(env string) *slog.Logger {

	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log

}
