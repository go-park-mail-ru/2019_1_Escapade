package database

import (
	"database/sql"
	"escapade/internal/models"
	"fmt"
)

// Register check sql-injections and are email and name unique
// Then add cookie to database and returns session_id
func (db *DataBase) SaveMessage(message *models.Message) (userID int, err error) {

	var (
		tx *sql.Tx
	)

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if userID, err = db.createMessage(tx, message); err != nil {
		return
	}

	fmt.Println("database/SaveMessage +")

	err = tx.Commit()
	return
}

func (db *DataBase) GetMessages() (messages []*models.Message, err error) {

	var (
		tx *sql.Tx
	)

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if messages, err = db.getMessages(tx); err != nil {
		return
	}

	fmt.Println("database/GetMessages +")

	err = tx.Commit()
	return
}
