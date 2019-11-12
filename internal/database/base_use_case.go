package database

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
)

type UseCaseBase struct {
	Db DatabaseI
}

func (rb *UseCaseBase) InitDBWithSQLPQ(CDB config.Database) error {
	var database = &PostgresSQL{}
	// в конфиг
	return rb.Open(CDB, 10, time.Hour, database)
}

func (rb *UseCaseBase) Open(CDB config.Database,
	maxIdleConns int, maxLifetime time.Duration, db DatabaseI) error {
	if err := db.Open(CDB); err != nil {
		return err
	}
	rb.Db = db
	db.Ping()
	rb.Db.SetMaxOpenConns(CDB.MaxOpenConns)
	rb.Db.SetMaxIdleConns(maxIdleConns)
	rb.Db.SetConnMaxLifetime(maxLifetime)
	return rb.Db.Ping()
}

func (rb *UseCaseBase) Use(db DatabaseI) error {
	rb.Db = db
	return rb.Db.Ping()
}

func (rb *UseCaseBase) Close() (err error) {
	return rb.Db.Close()
}

func (rb *UseCaseBase) Get() DatabaseI {
	return rb.Db
}
