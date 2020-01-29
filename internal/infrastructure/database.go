package infrastructure

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/entity"
)

//go:generate $GOPATH/bin/mockery -name "ExecerI|TransactionI|DatabaseI"

// ExecerI interface for executing queries in the database
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

// DatabaseI interface of database
type DatabaseI interface {
	ExecerI
	Open() error
	Begin() (TransactionI, error)
	PingContext(ctx context.Context) error
	Close() error
}

// data

type DatabaseRepositoryI interface {
	Get() entity.Database
}
