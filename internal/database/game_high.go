package database

import (
	"fmt"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	"database/sql"
)

func (db *DataBase) SaveGame(
	info models.GameInformation) (err error) {
	var (
		tx              *sql.Tx
		gameID, fieldID int
	)

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if gameID, err = db.createGame(tx, info.Game); err != nil {
		return
	}

	if err = db.createGamers(tx, gameID, info.Gamers); err != nil {
		return
	}

	if fieldID, err = db.createField(tx, gameID, info.Field); err != nil {
		return
	}

	if err = db.createActions(tx, gameID, info.Actions); err != nil {
		return
	}

	if err = db.createCells(tx, fieldID, info.Cells); err != nil {
		return
	}

	err = tx.Commit()
	fmt.Println("success save")
	return
}

func (db *DataBase) GetGame(roomID string) (
	game models.GameInformation, err error) {
	var (
		tx *sql.Tx
	)

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if game, err = db.GetGameInformation(tx, roomID); err != nil {
		return
	}

	fmt.Println("success get", game)
	err = tx.Commit()
	return
}

/*
func (db *DataBase) GetGames(userID int, page int) (
	games []*models.GameInformation, err error) {
	var (
		tx *sql.Tx
	)

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if games, err = db.GetFullGamesInformation(tx, userID, page); err != nil {
		return
	}

	err = tx.Commit()
	return
}
*/
