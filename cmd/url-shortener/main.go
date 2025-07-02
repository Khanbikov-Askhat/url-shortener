package main

import (
	"fmt"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	fmt.Println(os.Getwd())

	//Init config
	cfg := config.MustLoad()

	fmt.Println(cfg)

	//init logger: slog
	log := setupLogger(cfg.Env)
	log.Info("starting url-shortener", "env", cfg.Env)
	log.Debug("debug messages are enabled")

	//init storage: sqlite
	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		log.Error("failed to init storage", "error", sl.Error(err))
		os.Exit(1)
	}

	// id, err := storage.SaveUrl("https://google.com", "google")
	// if err != nil {
	// 	log.Error("failed to save url", sl.Error(err))
	// 	os.Exit(1)
	// }

	// log.Info("url saved", slog.Int64("id", id))

	// id, err = storage.SaveUrl("https://google.com", "google")
	// if err != nil {
	// 	log.Error("failed to save url", sl.Error(err))
	// 	os.Exit(1)
	// }

	// log.Info("url saved", slog.Int64("id", id))

	//init router: chi, "chi render"
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	//TODO: run server:
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
