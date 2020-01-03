package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
)

// Input - arguments, requeired to initialize database
type Input struct {
	Database database.Interface

	// Repository
	User    UserRepositoryI
	Message MessageRepositoryI
	Chat    ChatRepositoryI

	// UseCase
	UserUC    UserUseCaseI
	ChatUC    ChatUseCaseI
	MessageUC MessageUseCaseI
}

// InitAsPSQL initalize database as postgresql
func (db *Input) InitAsPSQL() *Input {
	db.Database = new(database.PostgresSQL)
	db.User = new(UserRepositoryPQ)
	db.Message = new(MessageRepositoryPQ)
	db.Chat = new(ChatRepositoryPQ)
	return db.Init()
}

// Init initalize database use cases
func (db *Input) Init() *Input {
	db.UserUC = new(UserUseCase).Init(db.User)
	db.MessageUC = new(MessageUseCase).Init(db.Message)
	db.ChatUC = new(ChatUseCase).Init(db.Chat, db.User)
	return db
}

// IsValid check if all interfaces are set
func (db *Input) IsValid() error {
	err := re.NoNil(db.Database, db.User, db.Chat, db.Message)
	if err == nil {
		return re.NoNil(db.UserUC, db.ChatUC, db.MessageUC)
	}
	return err
}

// Connect open connection to database and use it in every UseCase
func (db *Input) Connect(c config.Database) error {
	return database.Open(db.Database, c, db.UserUC, db.MessageUC, db.ChatUC)
}

// Close connection to database
func (db *Input) Close() error {
	return re.Close(db.Database, db.UserUC, db.MessageUC, db.ChatUC)
}
