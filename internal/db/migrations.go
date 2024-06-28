package db

import (
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
)

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("sqlite")
	if err != nil {
		return fmt.Errorf("migrate, set dialect = sqlite: %w", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}

func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()

	return Migrate(db, dir)
}
