package main

import (
	"context"
	"errors"
	"goth/internal/config"
	"goth/internal/handlers"
	// "goth/internal/hash/passwordhash"
	database "goth/internal/store/db"
	"goth/internal/store/dbstore"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	m "goth/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

/*
* Set to production at build time
* used to determine what assets to load
 */
var Environment = "development"

func init() {
	os.Setenv("env", Environment)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	r := chi.NewRouter()

	cfg := config.MustLoadConfig()

	db := database.MustOpen(cfg.DatabaseName)

	taskStore := dbstore.NewTaskStore(
		dbstore.NewTaskStoreParams{
			DB:           db,
		},
	)

	sessionStore := dbstore.NewSessionStore(
		dbstore.NewSessionStoreParams{
			DB: db,
		},
	)

	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	authMiddleware := m.NewAuthMiddleware(sessionStore, cfg.SessionCookieName)

	r.Group(func(r chi.Router) {
		r.Use(
			middleware.Logger,
			m.TextHTMLMiddleware,
			m.CSPMiddleware,
			authMiddleware.AddSessionToContext,
		)

		r.Get("/", handlers.NewHomeHandler(handlers.HomeHandlerParams{
            TaskStore: taskStore,
            SessionStore: sessionStore,
        }).ServeHTTP)

        r.Get("/todos/{id}", handlers.NewGetTaskHandler(handlers.GetTaskHandlerParams{
			TaskStore: taskStore,
		}).ServeHTTP)

        r.Get("/todos/{id}/edit", handlers.NewGetTaskEditorHandler(handlers.GetTaskEditorHandlerParams{
			TaskStore: taskStore,
		}).ServeHTTP)

		r.Post("/todos", handlers.NewPostTaskHandler(handlers.PostTaskHandlerParams{
			TaskStore: taskStore,
		}).ServeHTTP)

        r.Delete("/todos/{id}", handlers.NewDeleteTaskHandler(handlers.DeleteTaskHandlerParams{
			TaskStore: taskStore,
		}).ServeHTTP)

		r.Put("/todos/{id}", handlers.NewPutTaskHandler(handlers.PutTaskHandlerParams{
		    TaskStore: taskStore,
		}).ServeHTTP)

	})

	killSig := make(chan os.Signal, 1)

	signal.Notify(killSig, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: r,
	}

	go func() {
		err := srv.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			logger.Info("Server shutdown complete")
		} else if err != nil {
			logger.Error("Server error", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	logger.Info("Server started", slog.String("port", cfg.Port), slog.String("env", Environment))
	<-killSig

	logger.Info("Shutting down server")

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", slog.Any("err", err))
		os.Exit(1)
	}

	logger.Info("Server shutdown complete")
}
