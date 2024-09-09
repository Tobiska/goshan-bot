package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"goshan-bot/internal/config"
)

type SQLDatabase interface {
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type Client struct {
	sql.DB
}

func New(cfg *config.Database) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.Dsn)
	if err != nil {
		return nil, fmt.Errorf("error open sql conn: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error ping database: %w", err)
	}

	return db, nil
}
