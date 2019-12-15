package database

import (
	idb "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
)

//go:generate $GOPATH/bin/mockery -name "UserUseCaseI|UserRepositoryI|RecordUseCaseI|RecordRepositoryI|ImageUseCaseI|ImageRepositoryI"

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
	idb.UserCaseI
	Init(user UserRepositoryI, record RecordRepositoryI) UserUseCaseI

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
	Create(tx idb.TransactionI, user *models.UserPrivateInfo) (int, error)
	Delete(tx idb.TransactionI, user *models.UserPrivateInfo) error

	UpdateNamePassword(tx idb.TransactionI, user *models.UserPrivateInfo) error
	CheckNamePassword(tx idb.TransactionI, name string,
		password string) (int32, *models.UserPublicInfo, error)
	FetchNamePassword(tx idb.TransactionI,
		userID int32) (*models.UserPrivateInfo, error)

	UpdateLastSeen(tx idb.TransactionI, id int) error

	FetchAll(tx idb.TransactionI, params UsersSelectParams) ([]*models.UserPublicInfo, error)
	FetchOne(tx idb.TransactionI, userID int32,
		difficult int) (*models.UserPublicInfo, error)

	PagesCount(dbI idb.Interface, perPage int) (int, error)
}

type RecordUseCaseI interface {
	idb.UserCaseI
	Init(record RecordRepositoryI) RecordUseCaseI

	Update(id int32, record *models.Record) error
}

type RecordRepositoryI interface {
	Create(tx idb.TransactionI, id int) error
	Update(tx idb.TransactionI, id int32, record *models.Record) error
}

type ImageUseCaseI interface {
	idb.UserCaseI
	Init(image ImageRepositoryI) ImageUseCaseI

	Update(filename string, userID int32) error
	FetchByName(name string) (string, error)
	FetchByID(id int32) (string, error)
}

type ImageRepositoryI interface {
	Update(dbI idb.Interface, filename string, userID int32) error
	FetchByName(dbI idb.Interface, name string) (string, error)
	FetchByID(dbI idb.Interface, id int32) (string, error)
}
