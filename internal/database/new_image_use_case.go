package database

// ImageUseCase implements the interface ImageUseCaseI
type ImageUseCase struct {
	UseCaseBase
	image ImageRepositoryI
}

func (db *ImageUseCase) Init(image ImageRepositoryI) {
	db.image = image
}

// Update set filename of avatar to relation Player
func (db *ImageUseCase) Update(filename string, userID int32) error {
	return db.image.update(db.Db, filename, userID)
}

// FetchByName Get avatar - filename of player image by his name
func (db *ImageUseCase) FetchByName(name string) (string, error) {
	return db.image.fetchByName(db.Db, name)
}

// FetchByID Get avatar - filename of player image by his id
func (db *ImageUseCase) FetchByID(id int32) (string, error) {
	return db.image.fetchByID(db.Db, id)
}
