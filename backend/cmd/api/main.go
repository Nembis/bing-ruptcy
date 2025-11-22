package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	migration "github.com/nembis/bing-ruptcy/backend"
	"github.com/nembis/bing-ruptcy/backend/internal/config"
	"github.com/nembis/bing-ruptcy/backend/internal/otel"
	"github.com/nembis/bing-ruptcy/backend/internal/server"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Failed to start server", "Error", err)
	}
}

func run() (err error) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	if err = handleMigration(cfg.DatabaseURI); err != nil {
		return
	}

	conn, err := pgx.Connect(ctx, cfg.DatabaseURI)
	if err != nil {
		return
	}
	defer conn.Close(ctx)

	otelShutdown, err := otel.SetupOTelSDK(ctx)
	if err != nil {
		return
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	server := server.NewServer(conn, cfg)
	httpServer := server.GetHttpServelr(ctx)

	errChan := make(chan error, 1)
	go func() {
		slog.Info("Starting server", "Port", "4000")
		errChan <- httpServer.ListenAndServe()
	}()

	select {
	case err = <-errChan:
		break

	case <-ctx.Done():
		stop()
	}

	err = httpServer.Shutdown(ctx)
	return
}

func handleMigration(dbURI string) error {
	db, err := sql.Open("pgx", dbURI)
	if err != nil {
		return err
	}
	defer db.Close()

	if err = migration.HandleSqlMigrateUp(db, "pgx"); err != nil {
		return err
	}

	return nil
}
