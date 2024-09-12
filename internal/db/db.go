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

type Row interface {
	Scan(dest ...any) error
}

type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
}

type Transaction interface {
	Commit() error
	Exec(query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
	Rollback() error
}

type Handler struct {
	conn *sql.DB
}

func OpenDbHandler(cfg DbConfig) (Handler, error) {
	conn, err := sql.Open("sqlite3", cfg.database)
	if err != nil {
		return Handler{}, fmt.Errorf("open db connection: %w", err)
	}

	return Handler{conn: conn}, nil
}

func (handler Handler) Query(query string, args ...any) (Rows, error) {
	return handler.conn.Query(query, args...)
}

func (handler Handler) QueryRow(query string, args ...any) Row {
	return handler.conn.QueryRow(query, args...)
}

func (handler Handler) Begin() (Transaction, error) {
	return handler.conn.Begin()
}

func (handler Handler) MigrateFS(migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()

	return handler.migrate(dir)
}

func (handler Handler) migrate(dir string) error {
	err := goose.SetDialect("sqlite")
	if err != nil {
		return fmt.Errorf("migrate, set dialect = sqlite: %w", err)
	}

	err = goose.Up(handler.conn, dir)
	if err != nil {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}

func (handler Handler) Close() error {
	return handler.conn.Close()
}
