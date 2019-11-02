package database

import (
	//

	"database/sql"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// dataBaseI interface of database
type transactionI interface {
	Commit() error
	Rollback() error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// dataBaseI interface of database
type DatabaseI interface {
	Open(cdb config.Database) error
	Begin() (transactionI, error)
	SetMaxOpenConns(n int)
	SetConnMaxLifetime(d time.Duration)
	SetMaxIdleConns(n int)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Ping() error
	Close() error
}

// baseRepositoryI interface of base repository

type RepositoryBaseI interface {
	// open new db connection
	Open(CDB config.Database, maxIdleConns int,
		maxLifetime time.Duration, db DatabaseI) error
	// use exeisting openned connection
	Use(db DatabaseI) error
	Get() DatabaseI
	// close connection to db
	Close() error
}

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

type UserUseCaseI interface {
	RepositoryBaseI
	Init(user UserRepositoryI, record RecordRepositoryI)

	EnterAccount(name, password string) (int32, error)
	CreateAccount(user *models.UserPrivateInfo) (int, error)
	UpdateAccount(userID int32, user *models.UserPrivateInfo) (err error)
	DeleteAccount(user *models.UserPrivateInfo) error

	FetchAll(difficult int, page int, perPage int,
		sort string) (players []*models.UserPublicInfo, err error)
	FetchOne(userID int32, difficult int) (*models.UserPublicInfo, error)

	PagesCount(perPage int) (int, error)
}

type UserRepositoryI interface {
	create(tx transactionI, user *models.UserPrivateInfo) (int, error)
	delete(tx transactionI, user *models.UserPrivateInfo) error

	updateNamePassword(tx transactionI, user *models.UserPrivateInfo) error
	checkNamePassword(tx transactionI, name string,
		password string) (int32, *models.UserPublicInfo, error)
	fetchNamePassword(tx transactionI,
		userID int32) (*models.UserPrivateInfo, error)

	updateLastSeen(tx transactionI, id int) error

	fetchAll(tx transactionI, difficult int, offset int, limit int,
		sort string) ([]*models.UserPublicInfo, error)
	fetchOne(tx transactionI, userID int32,
		difficult int) (*models.UserPublicInfo, error)

	pagesCount(dbI DatabaseI, perPage int) (int, error)
}

type GameUseCaseI interface {
	RepositoryBaseI
	Init(game GameRepositoryI)

	Create(game *models.Game) (int32, int32, error)
	Save(info models.GameInformation) error

	FetchOneGame(roomID string) (models.GameInformation, error)
	FetchAllGames(userID int32) ([]models.GameInformation, error)

	FetchAllRoomsID(userID int32) ([]string, error)
}

type GameRepositoryI interface {
	createGame(tx transactionI, game *models.Game) (int32, error)
	updateGame(tx transactionI, game *models.Game) error

	createGamers(tx transactionI, GameID int32, gamers []models.Gamer) error
	createField(tx transactionI, gameID int32,
		field models.Field) (int32, error)
	createActions(tx transactionI, GameID int32, actions []models.Action) error
	createCells(tx transactionI, FieldID int32, cells []models.Cell) error

	fetchOneGame(tx transactionI, roomID string) (models.Game, error)
	fetchAllCells(tx transactionI, fieldID int) ([]models.Cell, error)
	fetchAllGamers(tx transactionI, gameID int32) ([]models.Gamer, error)
	fetchAllActions(tx transactionI, gameID int32) ([]models.Action, error)
	fetchOneField(tx transactionI, gameID int32) (int, models.Field, error)

	fetchAllRoomsID(tx transactionI, userID int32) ([]string, error)
}

type RecordUseCaseI interface {
	RepositoryBaseI
	Init(record RecordRepositoryI)

	Update(id int32, record *models.Record) error
}

type RecordRepositoryI interface {
	create(tx transactionI, id int) error
	update(tx transactionI, id int32, record *models.Record) error
}

type ImageUseCaseI interface {
	RepositoryBaseI
	Init(image ImageRepositoryI)

	Update(filename string, userID int32) error
	FetchByName(name string) (string, error)
	FetchByID(id int32) (string, error)
}

type ImageRepositoryI interface {
	update(dbI DatabaseI, filename string, userID int32) error
	fetchByName(dbI DatabaseI, name string) (string, error)
	fetchByID(dbI DatabaseI, id int32) (string, error)
}
