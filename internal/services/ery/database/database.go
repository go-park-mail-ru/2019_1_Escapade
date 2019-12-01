package database

import (
	"time"

	//
	_ "github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
)

type DB struct{ db *sqlx.DB }

func Init(link string, maxOpen, maxIdle int, ttl time.Duration) (*DB, error) {
	db, err := sqlx.Connect("postgres", link)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(ttl)

	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.Close()
}
