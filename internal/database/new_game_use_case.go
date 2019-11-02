package database

import (
	"context"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/clients"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	chat "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat"
)

// GameuseCasePQ implements the interface GameuseCaseI
type GameUseCase struct {
	UseCaseBase
	game GameRepositoryI
}

func (db *GameUseCase) Init(game GameRepositoryI) {
	db.game = game
}

func (db *GameUseCase) Create(game *models.Game) (int32, int32, error) {

	var (
		tx       transactionI
		roomID   int32
		pbChatID *chat.ChatID
		err      error
		id       int32
	)
	if tx, err = db.Db.Begin(); err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()
	if roomID, err = db.game.createGame(tx, game); err != nil {
		return 0, 0, err
	}

	newChat := &chat.ChatWithUsers{
		Type:   chat.ChatType_ROOM,
		TypeId: roomID,
	}

	pbChatID, err = clients.ALL.Chat().CreateChat(context.Background(), newChat)
	if err == nil {
		id = pbChatID.Value
	}
	err = tx.Commit()
	return roomID, id, err
}

// SaveGame save game to database
func (db *GameUseCase) Save(info models.GameInformation) error {
	var (
		tx              transactionI
		gameID, fieldID int32
		err             error
	)

	if tx, err = db.Db.Begin(); err != nil {
		return err
	}
	defer tx.Rollback()

	if err = db.game.updateGame(tx, &info.Game); err != nil {
		return err
	}

	/*
		msgs := MessagesToProto(info.Messages...)

		_, err = clients.ALL.Chat.AppendMessages(context.Background(), msgs)
		if err != nil {
			return
		}
	*/

	if err = db.game.createGamers(tx, gameID, info.Gamers); err != nil {
		return err
	}

	if fieldID, err = db.game.createField(tx, gameID, info.Field); err != nil {
		return err
	}

	if err = db.game.createActions(tx, gameID, info.Actions); err != nil {
		return err
	}

	if err = db.game.createCells(tx, fieldID, info.Cells); err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

// GetGames get list of games
func (db *GameUseCase) FetchAllGames(userID int32) ([]models.GameInformation, error) {
	var (
		tx    transactionI
		URLs  []string
		games []models.GameInformation
		err   error
	)

	if tx, err = db.Db.Begin(); err != nil {
		return games, err
	}
	defer tx.Rollback()

	if URLs, err = db.game.fetchAllRoomsID(tx, userID); err != nil {
		return games, err
	}

	games = make([]models.GameInformation, 0)
	for _, URL := range URLs {
		var info models.GameInformation
		if info, err = db.FetchOneGame(URL); err != nil {
			break
		}
		games = append(games, info)
	}

	err = tx.Commit()
	return games, err
}

// GetGamesURL get games url
func (db *GameUseCase) FetchAllRoomsID(userID int32) ([]string, error) {
	var (
		tx   transactionI
		URLs []string
		err  error
	)

	if tx, err = db.Db.Begin(); err != nil {
		return URLs, err
	}
	defer tx.Rollback()

	if URLs, err = db.game.fetchAllRoomsID(tx, userID); err != nil {
		return URLs, err
	}

	err = tx.Commit()
	return URLs, err
}

// GetGame get all information about game:
// game, gamers, field, history of cells and actions
func (db *GameUseCase) FetchOneGame(roomID string) (models.GameInformation, error) {
	var (
		tx              transactionI
		gameInformation models.GameInformation
		err             error
	)

	if tx, err = db.Db.Begin(); err != nil {
		return gameInformation, err
	}
	defer tx.Rollback()

	game, err := db.game.fetchOneGame(tx, roomID)
	if err != nil {
		return gameInformation, err
	}
	gameID := game.ID

	gamers, err := db.game.fetchAllGamers(tx, gameID)
	if err != nil {
		return gameInformation, err
	}

	fieldID, field, err := db.game.fetchOneField(tx, gameID)
	if err != nil {
		return gameInformation, err
	}

	actions, err := db.game.fetchAllActions(tx, gameID)
	if err != nil {
		return gameInformation, err
	}

	cells, err := db.game.fetchAllCells(tx, fieldID)
	if err != nil {
		return gameInformation, err
	}
	/*

		var (
			chat = &pChat.Chat{
				Type:   pChat.ChatType_ROOM,
				TypeId: gameID,
			}
			chatID    *pChat.ChatID
			pMessages *pChat.Messages
		)

		chatID, err = clients.ALL.Chat.GetChat(context.Background(), chat)

		if err != nil {
			utils.Debug(true, "cant access to chat service", err.Error())
		}
		pMessages, err = clients.ALL.Chat.GetChatMessages(context.Background(), chatID)

		var messages []*models.Message
		messages = MessagesFromProto(pMessages.Messages...)
		//db.getMessages(tx, true, game.RoomID)
		if err != nil {
			utils.Debug(true, "cant get messages!", err.Error())
		}
	*/

	gameInformation = models.GameInformation{
		Game:    game,
		Gamers:  gamers,
		Field:   field,
		Actions: actions,
		Cells:   cells,
		//Messages: messages,
	}

	err = tx.Commit()
	return gameInformation, err
}
