package database

import (
	chat "github.com/go-park-mail-ru/2019_1_Escapade/chat/proto"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/clients"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	"context"
	"database/sql"
)

func (db *DataBase) CreateGame(game *models.Game) (int32, int32, error) {

	var (
		tx       *sql.Tx
		roomID   int32
		pbChatID *chat.ChatID
		err      error
	)
	if tx, err = db.Db.Begin(); err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()
	if roomID, err = db.createGame(tx, game); err != nil {
		return 0, 0, err
	}

	newChat := &chat.ChatWithUsers{
		Type:   chat.ChatType_ROOM,
		TypeId: roomID,
	}

	pbChatID, err = clients.ALL.Chat.CreateChat(context.Background(), newChat)
	if err != nil {
		return 0, 0, err
	}
	err = tx.Commit()
	return roomID, pbChatID.Value, err
}

// SaveGame save game to database
func (db *DataBase) SaveGame(
	info models.GameInformation) (err error) {
	var (
		tx              *sql.Tx
		gameID, fieldID int32
	)

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if err = db.updateGame(tx, info.Game); err != nil {
		return
	}

	/*
		msgs := MessagesToProto(info.Messages...)

		_, err = clients.ALL.Chat.AppendMessages(context.Background(), msgs)
		if err != nil {
			return
		}
	*/

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
	return
}

// GetGames get list of games
func (db *DataBase) GetGames(userID int32) (
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
		var info models.GameInformation
		if info, err = db.GetGame(URL); err != nil {
			break
		}
		games = append(games, info)
	}

	err = tx.Commit()
	return
}

// GetGamesURL get games url
func (db *DataBase) GetGamesURL(userID int32) (
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

	err = tx.Commit()
	return
}
