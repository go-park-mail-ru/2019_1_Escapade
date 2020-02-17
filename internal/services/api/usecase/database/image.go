package database

import (
	"context"
	"errors"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/repository/database"
)

//Image implements the interface ImageUseCaseI
type Image struct {
	db             infrastructure.Database
	trace          infrastructure.ErrorTrace
	imageDB        api.ImageRepositoryI
	contextTimeout time.Duration
}

// NewImage create new instance of Image
func NewImage(
	dbI infrastructure.Database,
	trace infrastructure.ErrorTrace,
	timeout time.Duration,
) (*Image, error) {
	if dbI == nil {
		return nil, errors.New(ErrNoDatabase)
	}
	imageRep, err := database.NewImage(dbI)
	if err != nil {
		return nil, err
	}
	if trace == nil {
		trace = new(infrastructure.ErrorTraceNil)
	}
	return &Image{
		db:             dbI,
		trace:          trace,
		imageDB:        imageRep,
		contextTimeout: timeout,
	}, nil
}

// Update set filename of avatar to relation Player
func (usecase *Image) Update(
	c context.Context,
	filename string,
	userID int32,
) error {
	if filename == "" {
		return usecase.trace.New(InvalidID)
	}
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
	if name == "" {
		return "", usecase.trace.New(InvalidUserName)
	}
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
	if id <= 0 {
		return "", usecase.trace.New(InvalidID)
	}
	ctx, cancel := context.WithTimeout(
		c,
		usecase.contextTimeout,
	)
	defer cancel()
	return usecase.imageDB.FetchByID(ctx, id)
}
