package database

import (
	idb "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database"
)

// Input - arguments, requeired to initialize database
type Input struct {
	Database idb.Interface
	User     api.UserUseCaseI
	Game     GameUseCaseI
}

func (input *Input) InitAsPSQL() *Input {
	var (
		user   = &api.UserRepositoryPQ{}
		record = &api.RecordRepositoryPQ{}
		game   = &GameRepositoryPQ{}
	)
	input.Database = &idb.PostgresSQL{}
	input.User = new(api.UserUseCase).Init(user, record)
	input.Game = new(GameUseCase).Init(game)
	return input
}
