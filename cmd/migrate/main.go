package main

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/Strangebrewer/go-service-template/config"
)

func main() {
	if len(os.Args) < 2 {
		slog.Error("usage: migrate [up|down]")
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()

	dbURL := strings.NewReplacer("postgres://", "pgx5://", "postgresql://", "pgx5://").Replace(cfg.DatabaseURL)

	m, err := migrate.New("file://db/migrations", dbURL)
	if err != nil {
		slog.Error("failed to create migrator", "error", err)
		os.Exit(1)
	}
	defer m.Close()

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			slog.Error("migration up failed", "error", err)
			os.Exit(1)
		}
		slog.Info("migrations applied")
	case "down":
		if err := m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			slog.Error("migration down failed", "error", err)
			os.Exit(1)
		}
		slog.Info("migration rolled back one step")
	default:
		slog.Error("unknown command, use up or down", "command", os.Args[1])
		os.Exit(1)
	}
}
