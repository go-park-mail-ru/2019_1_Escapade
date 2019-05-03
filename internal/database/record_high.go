package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	"database/sql"
	"fmt"
)

// UpdateRecords update records for offline game
func (db *DataBase) UpdateRecords(id int,
	record *models.Record) (err error) {
	var (
		tx *sql.Tx
	)

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if err = db.updateRecords(tx, id, record); err != nil {
		fmt.Println("updateRecords err:", err.Error())
		return
	}

	err = tx.Commit()
	return
}
