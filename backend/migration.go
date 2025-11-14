package migration

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed sql/schema/*.sql
var embedMigrations embed.FS

func HandleSqlMigrateUp(db *sql.DB, dialect string) error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect(dialect); err != nil {
		return err
	}

	return goose.Up(db, "sql/schema")
}
