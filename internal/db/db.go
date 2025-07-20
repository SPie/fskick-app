package db

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

var (
	ErrNotFound = sql.ErrNoRows
)

type DbConfig struct {
	database  string
	withDebug bool
	withLog   bool
}

func CreateDbConfig(database string, withDebug bool, withLog bool) DbConfig {
	return DbConfig{
		database:  database,
		withDebug: withDebug,
		withLog:   withLog,
	}
}

type Connection interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
	Begin() (*sql.Tx, error)
	Close() error
}

func OpenDbConnection(cfg DbConfig) (*sql.DB, error) {
	conn, err := sql.Open("sqlite3", cfg.database)
	if err != nil {
		return nil, fmt.Errorf("open db connection: %w", err)
	}

	return conn, nil
}

func MigrateFS(conn *sql.DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()

	err := goose.SetDialect("sqlite")
	if err != nil {
		return fmt.Errorf("migrate, set dialect = sqlite: %w", err)
	}

	err = goose.Up(conn, dir)
	if err != nil {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}
