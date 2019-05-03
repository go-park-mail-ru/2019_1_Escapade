package database

import (
	"database/sql"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	"fmt"
)

// SaveMessage save messages to database
func (db *DataBase) SaveMessage(message *models.Message,
	inRoom bool, gameID string) (err error) {

	var (
		tx *sql.Tx
	)

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if err = db.createMessage(tx, message, inRoom, gameID); err != nil {
		return
	}

	fmt.Println("database/SaveMessage +")

	err = tx.Commit()
	return
}

// LoadMessages load messages from database
func (db *DataBase) LoadMessages(inRoom bool, gameID string) (messages []*models.Message, err error) {

	var (
		tx *sql.Tx
	)

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if messages, err = db.getMessages(tx, inRoom, gameID); err != nil {
		return
	}

	fmt.Println("database/GetMessages +")

	err = tx.Commit()
	return
}
