package main

import (
	"database/sql"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	migration "github.com/nembis/bing-ruptcy/backend"
	"github.com/nembis/bing-ruptcy/backend/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", "Error", err)
		return
	}

	db, err := sql.Open("pgx", cfg.DatabaseURI)
	if err != nil {
		slog.Error("Failed to connect to postgres", "Error", err)
		return
	}

	if err = migration.HandleSqlMigrateUp(db, "pgx"); err != nil {
		slog.Error("Failed migrate up database schema", "Error", err)
		return
	}
	slog.Info("Migration completed")

	slog.Info("Hello this is a message")
}
