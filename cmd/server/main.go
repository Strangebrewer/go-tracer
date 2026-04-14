package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Strangebrewer/go-service-template/app"
	"github.com/Strangebrewer/go-service-template/config"
	"github.com/Strangebrewer/go-service-template/db_connection"
	"github.com/Strangebrewer/go-service-template/example"
	"github.com/Strangebrewer/go-service-template/middleware"
	"github.com/Strangebrewer/go-service-template/server"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()

	pool, err := db_connection.NewPool(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	authMiddleware, err := middleware.RequireAuth(cfg.JWTPublicKey)
	if err != nil {
		slog.Error("failed to parse JWT public key", "error", err)
		os.Exit(1)
	}

	application := &app.Application{
		ExampleStore: example.NewStore(pool),
	}

	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	srv := server.New(":"+port, cfg.AllowedOrigins, application, authMiddleware)

	go func() {
		slog.Info("server starting", "port", port)
		if err := srv.HTTPServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.HTTPServer.Shutdown(ctx); err != nil {
		slog.Error("server shutdown failed", "error", err)
		os.Exit(1)
	}
	slog.Info("server stopped")
}
