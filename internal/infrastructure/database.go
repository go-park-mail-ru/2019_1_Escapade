package infrastructure

import (
	"context"
	"database/sql"
)

//go:generate $GOPATH/bin/mockery -name "TransactionI|DatabaseI"

type ExecerI interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// TransactionI interface of transaction
type TransactionI interface {
	ExecerI
	Commit() error
	Rollback() error
}

// Interface interface of database
type DatabaseI interface {
	ExecerI
	Open() error
	Begin() (TransactionI, error)
	PingContext(ctx context.Context) error
	Close() error
}
