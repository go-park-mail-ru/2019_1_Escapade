package database

import (
	"database/sql"
	"escapade/internal/models"
	"fmt"
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
		fmt.Println("db/createGame err:", err.Error())
		return
	}

	if err = db.createGamers(tx, userID, info.Gamers); err != nil {
		fmt.Println("db/createGamers err:", err.Error())
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
		fmt.Println("db/createGame err:", err.Error())
		return
	}

	err = tx.Commit()
	return
}
