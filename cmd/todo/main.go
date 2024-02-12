package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"log/slog"
	"net/http"
	"os"
	"to-do-app/internal/config"
	"to-do-app/internal/http-server/handlers/todo"
	"to-do-app/internal/http-server/middleware/logger"
	"to-do-app/internal/storage/postgres"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := SetupLogger(cfg.Env)

	log.Info("App started", slog.String("env", cfg.Env))
	log.Debug("Debugging started")

	storage, err := postgres.New(
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
	)
	if err != nil {
		log.Error("failed to init storage: ", err)
		os.Exit(1)
	}

	log.Info("storage successfully initialized")

	router := chi.NewRouter()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Разрешить все домены, для разработки подойдет, для продакшена лучше указать конкретные домены
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Максимальное время предварительного запроса (preflight request) в секундах
	})

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(corsHandler.Handler)

	router.Post("/add", todo.Create(log, storage))
	router.Post("/del", todo.Delete(log, storage))
	router.Post("/complete", todo.Complete(log, storage))
	router.Post("/edit", todo.Edit(log, storage))
	router.Get("/tasks", todo.GetTasks(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
