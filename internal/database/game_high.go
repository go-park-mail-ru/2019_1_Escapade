package database

import (
	"fmt"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	"database/sql"
)

// SaveGame save game to database
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

// GetGames get list of games
func (db *DataBase) GetGames(userID int) (
	games []models.GameInformation, err error) {
	var (
		tx   *sql.Tx
		URLs []string
	)

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if URLs, err = db.getGamesURL(tx, userID); err != nil {
		return
	}

	games = make([]models.GameInformation, 0)
	for _, URL := range URLs {
		info := models.GameInformation{}
		if info, err = db.GetGame(URL); err != nil {
			break
		}
		games = append(games, info)
	}

	err = tx.Commit()
	return
}

// GetGamesURL get games url
func (db *DataBase) GetGamesURL(userID int) (
	URLs []string, err error) {
	var (
		tx *sql.Tx
	)

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if URLs, err = db.getGamesURL(tx, userID); err != nil {
		return
	}

	err = tx.Commit()
	return
}

// GetGame get game
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
