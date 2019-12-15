package database

import idb "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"

// ImageUseCase implements the interface ImageUseCaseI
type ImageUseCase struct {
	idb.UseCaseBase
	image ImageRepositoryI
}

func (db *ImageUseCase) Init(image ImageRepositoryI) ImageUseCaseI {
	db.image = image
	return db
}

// Update set filename of avatar to relation Player
func (db *ImageUseCase) Update(filename string, userID int32) error {
	return db.image.Update(db.Db, filename, userID)
}

// FetchByName Get avatar - filename of player image by his name
func (db *ImageUseCase) FetchByName(name string) (string, error) {
	return db.image.FetchByName(db.Db, name)
}

// FetchByID Get avatar - filename of player image by his id
func (db *ImageUseCase) FetchByID(id int32) (string, error) {
	return db.image.FetchByID(db.Db, id)
}
