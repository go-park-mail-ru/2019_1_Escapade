package database

import (
	"os"
	"fmt"
	"context"
	"database/sql"
	"time"

	// postgresql
	_ "github.com/lib/pq"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
)

// PostgresSQL represents postgresql database
type PostgresSQL struct {
	Db *sql.DB
}

// Open connection to database
func (db *PostgresSQL) Open(CDB config.Database) error {
	var err error
	var connStr = os.Getenv("DB_CONN_STRING") //CDB.ConnectionString
	fmt.Println("connStr is", connStr)
	db.Db, err = sql.Open(CDB.DriverName, connStr)
	return err
}

// Begin transaction
func (db *PostgresSQL) Begin() (TransactionI, error) {
	return db.Db.Begin()
}

// Ping database
func (db *PostgresSQL) Ping() error {
	return db.Db.PingContext(context.Background())
}

// Close database connection
func (db *PostgresSQL) Close() error {
	return db.Db.Close()
}

// SetMaxOpenConns set max open conns
func (db *PostgresSQL) SetMaxOpenConns(n int) {
	db.Db.SetMaxOpenConns(n)
}

// SetConnMaxLifetime set max lifetime
func (db *PostgresSQL) SetConnMaxLifetime(d time.Duration) {
	db.Db.SetConnMaxLifetime(d)
}

// SetMaxIdleConns set max idle conn
func (db *PostgresSQL) SetMaxIdleConns(n int) {
	db.Db.SetMaxIdleConns(n)
}

// Exec psql
func (db *PostgresSQL) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.Db.Exec(query, args...)
}

// Query psql
func (db *PostgresSQL) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Db.Query(query, args...)
}

// QueryRow psql
func (db *PostgresSQL) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.Db.QueryRow(query, args...)
}
