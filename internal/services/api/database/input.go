package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"
	idb "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
)

// Input - arguments, requeired to initialize database
type Input struct {
	Database idb.Interface

	// Repository
	User   UserRepositoryI
	Record RecordRepositoryI
	Image  ImageRepositoryI

	// UseCase
	UserUC   UserUseCaseI
	RecordUC RecordUseCaseI
	ImageUC  ImageUseCaseI
}

// InitAsPSQL initalize database as postgresql
func (db *Input) InitAsPSQL() *Input {
	db.User = new(UserRepositoryPQ)
	db.Record = new(RecordRepositoryPQ)
	db.Image = new(ImageRepositoryPQ)
	db.Database = new(idb.PostgresSQL)
	return db
}

// Init initalize database use cases
func (db *Input) Init() *Input {
	db.UserUC = new(UserUseCase).Init(db.User, db.Record)
	db.RecordUC = new(RecordUseCase).Init(db.Record)
	db.ImageUC = new(ImageUseCase).Init(db.Image)
	return db
}

// IsValid check if all interfaces are set
func (db *Input) IsValid() error {
	err := re.NoNil(db.Database, db.User, db.Record, db.Image)
	if err == nil {
		return re.NoNil(db.UserUC, db.RecordUC, db.ImageUC)
	}
	return err
}

// Connect open connection to database and use it in every UseCase
func (db *Input) Open(c config.Database) error {
	return database.Open(db.Database, c, db.UserUC, db.RecordUC, db.ImageUC)
}

// Close connection to database
func (db *Input) Close() error {
	return re.Close(db.Database, db.UserUC, db.RecordUC, db.ImageUC)
}
