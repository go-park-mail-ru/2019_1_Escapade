package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"
	idb "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database"
)

// Input - arguments, requeired to initialize database
type Input struct {
	Database idb.Interface

	User   api.UserRepositoryI
	Record api.RecordRepositoryI
	Game   GameRepositoryI

	UserUC api.UserUseCaseI
	GameUC GameUseCaseI
}

// InitAsPSQL initalize database as postgresql
func (db *Input) InitAsPSQL() *Input {
	db.Database = new(idb.PostgresSQL)
	db.User = new(api.UserRepositoryPQ)
	db.Record = new(api.RecordRepositoryPQ)
	db.Game = new(GameRepositoryPQ)
	return db.Init()
}

// Init initalize database use cases
func (db *Input) Init() *Input {
	db.UserUC = new(api.UserUseCase).Init(db.User, db.Record)
	db.GameUC = new(GameUseCase).Init(db.Game)
	return db
}

// IsValid check if all interfaces are set
func (db *Input) IsValid() error {
	err := re.NoNil(db.Database, db.User, db.Record, db.User)
	if err == nil {
		return re.NoNil(db.UserUC, db.GameUC)
	}
	return err
}

// Connect open connection to database and use it in every UseCase
func (db *Input) Connect(c config.Database) error {
	return database.Open(db.Database, c, db.UserUC, db.GameUC)
}

// Close connection to database
func (db *Input) Close() error {
	return re.Close(db.Database, db.UserUC, db.UserUC, db.GameUC)
}
