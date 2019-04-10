package database

import (
	"database/sql"
	"escapade/internal/models"
	"fmt"
)

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
	fmt.Println("done")
	return
}
