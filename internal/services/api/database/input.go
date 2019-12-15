package database

import idb "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"

type Input struct {
	User   UserRepositoryI
	Record RecordRepositoryI
	Image  ImageRepositoryI

	Database idb.Interface
}

func (input *Input) InitAsPSQL() *Input {
	input.User = new(UserRepositoryPQ)
	input.Record = new(RecordRepositoryPQ)
	input.Image = new(ImageRepositoryPQ)
	input.Database = new(idb.PostgresSQL)
	return input
}
