package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
)

type UseCaseBase struct {
	Db Interface
}

func (rb *UseCaseBase) Open(CDB config.Database, db Interface) error {
	if err := db.Open(CDB); err != nil {
		return err
	}
	rb.Db = db
	rb.Db.SetMaxOpenConns(CDB.MaxOpenConns)
	rb.Db.SetMaxIdleConns(CDB.MaxIdleConns)
	rb.Db.SetConnMaxLifetime(CDB.MaxLifetime.Duration)
	return rb.Db.Ping()
}

func (rb *UseCaseBase) Use(db Interface) error {
	rb.Db = db
	return rb.Db.Ping()
}

func (rb *UseCaseBase) Close() (err error) {
	return rb.Db.Close()
}

func (rb *UseCaseBase) Get() Interface {
	return rb.Db
}
