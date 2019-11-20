package database

import (
	"database/sql"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
)

// UserCaseI interface of base user case
type UserCaseI interface {
	// open new db connection
	Open(CDB config.Database, maxIdleConns int,
		maxLifetime time.Duration, db DatabaseI) error
	// use exeisting openned connection
	Use(db DatabaseI) error
	Get() DatabaseI
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

// DatabaseI interface of database
type DatabaseI interface {
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
