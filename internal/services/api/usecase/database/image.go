package database

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/repository/database"
)

//Image implements the interface ImageUseCaseI
type Image struct {
	db             infrastructure.DatabaseI
	imageDB        api.ImageRepositoryI
	contextTimeout time.Duration
}

func NewImage(
	dbI infrastructure.DatabaseI,
	timeout time.Duration,
) *Image {
	return &Image{
		db:             dbI,
		imageDB:        database.NewImage(dbI),
		contextTimeout: timeout,
	}
}

// Update set filename of avatar to relation Player
func (usecase *Image) Update(
	c context.Context,
	filename string,
	userID int32,
) error {
	ctx, cancel := context.WithTimeout(
		c,
		usecase.contextTimeout,
	)
	defer cancel()
	return usecase.imageDB.Update(ctx, filename, userID)
}

// FetchByName Get avatar - filename of player image by his name
func (usecase *Image) FetchByName(
	c context.Context,
	name string,
) (string, error) {
	ctx, cancel := context.WithTimeout(
		c,
		usecase.contextTimeout,
	)
	defer cancel()
	return usecase.imageDB.FetchByName(ctx, name)
}

// FetchByID Get avatar - filename of player image by his id
func (usecase *Image) FetchByID(
	c context.Context,
	id int32,
) (string, error) {
	ctx, cancel := context.WithTimeout(
		c,
		usecase.contextTimeout,
	)
	defer cancel()
	return usecase.imageDB.FetchByID(ctx, id)
}
