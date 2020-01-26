package database

import (
	"context"
	"time"

	idb "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/repository/database"
)

//Image implements the interface ImageUseCaseI
type Image struct {
	db             idb.Interface
	imageDB        api.ImageRepositoryI
	contextTimeout time.Duration
}

func NewImage(dbI idb.Interface, timeout time.Duration) *Image {
	return &Image{
		db:             dbI,
		imageDB:        database.NewImage(dbI),
		contextTimeout: timeout,
	}
}

// Update set filename of avatar to relation Player
func (repository *Image) Update(c context.Context, filename string, userID int32) error {
	ctx, cancel := context.WithTimeout(c, repository.contextTimeout)
	defer cancel()
	err := repository.imageDB.Update(ctx, filename, userID)
	return err
}

// FetchByName Get avatar - filename of player image by his name
func (repository *Image) FetchByName(c context.Context, name string) (string, error) {
	ctx, cancel := context.WithTimeout(c, repository.contextTimeout)
	defer cancel()
	str, err := repository.imageDB.FetchByName(ctx, name)
	return str, err
}

// FetchByID Get avatar - filename of player image by his id
func (repository *Image) FetchByID(c context.Context, id int32) (string, error) {
	ctx, cancel := context.WithTimeout(c, repository.contextTimeout)
	defer cancel()
	str, err := repository.imageDB.FetchByID(ctx, id)
	return str, err
}
