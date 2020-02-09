package postgresql

import (
	"context"
	"database/sql"
	"errors"

	// postgresql
	_ "github.com/lib/pq"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/configuration"
)

// PostgresSQL represents postgresql database
type PostgresSQL struct {
	c  configuration.Database
	Db *sql.DB
}

// NewPostgresSQL create new instance of PostgresSQL
func New(conf configuration.DatabaseRepository) (*PostgresSQL, error) {
	if conf == nil {
		return nil, errors.New(ErrorNoConfiguration)
	}
	return &PostgresSQL{c: conf.Get()}, nil
}

func (db *PostgresSQL) Run() error {
	return db.Open()
}

// Open connection to psql and set:
//	- amount of max open connections
//  - amount of max idle connections
//  - max lifetime of connection
// return result of ping
func (db *PostgresSQL) Open() error {
	str := db.c.ConnectionString
	connectionString := DRIVER + "://" +
		str.User + ":" +
		str.Password + "@" + str.Address +
		"/" + str.Database + "?sslmode=disable"
	database, err := sql.Open(DRIVER, connectionString)
	if err != nil {
		return err
	}
	db.Db = database
	db.Db.SetMaxOpenConns(db.c.MaxOpenConns)
	db.Db.SetMaxIdleConns(db.c.MaxIdleConns)
	db.Db.SetConnMaxLifetime(db.c.MaxLifetime)
	return db.Db.Ping()
}

// Begin starts a transaction. The default isolation level is dependent on
// the driver.
func (db *PostgresSQL) Begin() (infrastructure.Transaction, error) {
	return db.Db.Begin()
}

// PingContext verifies a connection to the database is still alive,
// establishing a connection if necessary.
func (db *PostgresSQL) PingContext(
	ctx context.Context,
) error {
	return db.Db.PingContext(context.Background())
}

// Close closes the database and prevents new queries from starting.
// Close then waits for all queries that have started processing on the server
// to finish.
//
// It is rare to Close a DB, as the DB handle is meant to be
// long-lived and shared between many goroutines.
func (db *PostgresSQL) Close() error {
	return db.Db.Close()
}

// ExecContext executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (db *PostgresSQL) ExecContext(
	ctx context.Context,
	query string,
	args ...interface{},
) (sql.Result, error) {
	return db.Db.ExecContext(ctx, query, args...)
}

// QueryContext executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (db *PostgresSQL) QueryContext(
	ctx context.Context,
	query string,
	args ...interface{},
) (*sql.Rows, error) {
	return db.Db.QueryContext(ctx, query, args...)
}

// QueryRowContext executes a query that is expected to return at most one row.
// QueryRowContext always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
func (db *PostgresSQL) QueryRowContext(
	ctx context.Context,
	query string,
	args ...interface{}) *sql.Row {
	return db.Db.QueryRowContext(ctx, query, args...)
}
