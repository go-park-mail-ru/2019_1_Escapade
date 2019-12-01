package database

import (
	idb "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/clients"
)

/*
	[ObjectName]UseCase opens a connection to the database, but cannot perform
	operations on it. All operations are performed using [ObjectName]Repository
	classes. It([ObjectName]Repository) cant open a database connection itself
	and all it's methods require a database connection as input data.
	[ObjectName]UseCase is responsible for closing the connection,
	opening transactions and other actions related directly to the connection.
	[ObjectName]UseCase call methods of [ObjectName]Repository to perform
	operations in database
*/

type GameUseCaseI interface {
	idb.UserCaseI
	Init(game GameRepositoryI, chatS clients.Chat)

	Create(game *models.Game) (int32, int32, error)
	Save(info models.GameInformation) error

	FetchOneGame(roomID string) (models.GameInformation, error)
	FetchAllGames(userID int32) ([]models.GameInformation, error)

	FetchAllRoomsID(userID int32) ([]string, error)
}

type GameRepositoryI interface {
	createGame(tx idb.TransactionI, game *models.Game) (int32, error)
	updateGame(tx idb.TransactionI, game *models.Game) error

	createGamers(tx idb.TransactionI, GameID int32, gamers []models.Gamer) error
	createField(tx idb.TransactionI, gameID int32,
		field models.Field) (int32, error)
	createActions(tx idb.TransactionI, GameID int32, actions []models.Action) error
	createCells(tx idb.TransactionI, FieldID int32, cells []models.Cell) error

	fetchOneGame(tx idb.TransactionI, roomID string) (models.Game, error)
	fetchAllCells(tx idb.TransactionI, fieldID int) ([]models.Cell, error)
	fetchAllGamers(tx idb.TransactionI, gameID int32) ([]models.Gamer, error)
	fetchAllActions(tx idb.TransactionI, gameID int32) ([]models.Action, error)
	fetchOneField(tx idb.TransactionI, gameID int32) (int, models.Field, error)

	fetchAllRoomsID(tx idb.TransactionI, userID int32) ([]string, error)
}
