package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// UpdateRecords update records for offline game
func (db *DataBase) UpdateRecords(id int32, record *models.Record) error {
	tx, err := db.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = db.updateRecords(tx, id, record); err != nil {
		return err
	}

	err = tx.Commit()
	return err
}
