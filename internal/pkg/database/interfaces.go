package database

import (
	"database/sql"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
)

//go:generate $GOPATH/bin/mockery -name "UserCaseI|TransactionI|DatabaseI"

// UserCaseI interface of base user case
type UserCaseI interface {
	// open new db connection
	Open(CDB config.Database, db Interface) error
	// use existing openned connection
	Use(db Interface) error
	Get() Interface
	// close connection to db
	Close() error
}

// TransactionI interface of transaction
type TransactionI interface {
	Commit() error
	Rollback() error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// Interface interface of database
type Interface interface {
	Open(cdb config.Database) error
	Begin() (TransactionI, error)
	SetMaxOpenConns(n int)
	SetConnMaxLifetime(d time.Duration)
	SetMaxIdleConns(n int)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Ping() error
	Close() error
}
