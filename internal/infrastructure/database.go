package infrastructure

import (
	"context"
	"database/sql"
)

//go:generate $GOPATH/bin/mockery -name "ExecerI|TransactionI|DatabaseI"

const ErrNoDatabase = "Database interface not given"

// ExecerI interface for executing queries in the database
type Execer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// TransactionI interface of transaction
type Transaction interface {
	Execer
	Commit() error
	Rollback() error
}

// DatabaseI interface of database
type Database interface {
	Execer
	Open() error
	Begin() (Transaction, error)
	PingContext(ctx context.Context) error
	Close() error
}
