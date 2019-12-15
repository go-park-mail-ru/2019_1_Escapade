package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"
	idb "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
)

//go:generate $GOPATH/bin/mockery -name "GameUseCaseI|GameRepositoryI"

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

//mockgen -source=interfaces.go -destination=mock_database.go -package=database -aux_files=github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database=interfaces.go

type GameUseCaseI interface {
	database.UserCaseI
	Init(game GameRepositoryI) GameUseCaseI

	Create(game *models.Game) (int32, error)
	Save(info models.GameInformation) error

	FetchOneGame(roomID string) (models.GameInformation, error)
	FetchAllGames(userID int32) ([]models.GameInformation, error)

	FetchAllRoomsID(userID int32) ([]string, error)
}

type GameRepositoryI interface {
	CreateGame(tx idb.TransactionI, game *models.Game) (int32, error)
	UpdateGame(tx idb.TransactionI, game *models.Game) error

	CreateGamers(tx idb.TransactionI, GameID int32, gamers []models.Gamer) error
	CreateField(tx idb.TransactionI, gameID int32,
		field models.Field) (int32, error)
	CreateActions(tx idb.TransactionI, GameID int32, actions []models.Action) error
	CreateCells(tx idb.TransactionI, FieldID int32, cells []models.Cell) error

	FetchOneGame(tx idb.TransactionI, roomID string) (models.Game, error)
	FetchAllCells(tx idb.TransactionI, fieldID int) ([]models.Cell, error)
	FetchAllGamers(tx idb.TransactionI, gameID int32) ([]models.Gamer, error)
	FetchAllActions(tx idb.TransactionI, gameID int32) ([]models.Action, error)
	FetchOneField(tx idb.TransactionI, gameID int32) (int, models.Field, error)

	FetchAllRoomsID(tx idb.TransactionI, userID int32) ([]string, error)
}
