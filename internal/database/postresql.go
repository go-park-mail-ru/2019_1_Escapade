package database

import (
	//
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	// postgresql
	_ "github.com/lib/pq"
)

type PostgresSQL struct {
	//sql.DB
	Db *sql.DB
}

func (db *PostgresSQL) Open(CDB config.Database) error {
	var (
		err      error
		database *sql.DB
	)
	fmt.Println("try open")
	if database, err = sql.Open(CDB.DriverName, os.Getenv(CDB.URL)); err != nil {
		return err
	}
	fmt.Println("opened")

	db.Db = database
	db.Db.SetMaxOpenConns(CDB.MaxOpenConns)
	err = db.Db.Ping()
	if err == nil {
		fmt.Println("no err")
	} else {
		fmt.Println("some err", err.Error())
	}

	return err
}

func (db *PostgresSQL) Begin() (TransactionI, error) {
	return db.Db.Begin()
}

func (db *PostgresSQL) Ping() error {
	return db.Db.PingContext(context.Background())
}

func (db *PostgresSQL) Close() error {
	return db.Db.Close()
}

func (db *PostgresSQL) SetMaxOpenConns(n int) {
	db.Db.SetMaxOpenConns(n)
}

func (db *PostgresSQL) SetConnMaxLifetime(d time.Duration) {
	db.Db.SetConnMaxLifetime(d)
}

func (db *PostgresSQL) SetMaxIdleConns(n int) {
	db.Db.SetMaxIdleConns(n)
}

func (db *PostgresSQL) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.Db.Exec(query, args)
}

func (db *PostgresSQL) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Db.Query(query, args)
}

func (db *PostgresSQL) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.Db.QueryRow(query, args)
}
