package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	"database/sql"
)

func (db *DataBase) SaveGame(userID int,
	info *models.GameInformation) (err error) {
	var (
		tx *sql.Tx
	)

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if err = db.createGame(tx, userID, info.Game); err != nil {
		return
	}

	if err = db.createGamers(tx, userID, info.Gamers); err != nil {
		return
	}

	err = tx.Commit()
	return
}

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
