package api

import (
	"context"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
)

//go:generate $GOPATH/bin/mockery -name "UserUseCaseI|RecordUseCaseI|ImageUseCaseI"

type UserUseCaseI interface {
	EnterAccount(ctx context.Context,
		name, password string) (int32, error)
	CreateAccount(ctx context.Context,
		user *models.UserPrivateInfo) (int, error)
	UpdateAccount(ctx context.Context,
		userID int32, user *models.UserPrivateInfo) (err error)
	DeleteAccount(ctx context.Context,
		user *models.UserPrivateInfo) error

	FetchAll(ctx context.Context,
		difficult int, page int, perPage int,
		sort string) (players []*models.UserPublicInfo, err error)
	FetchOne(ctx context.Context,
		userID int32, difficult int) (*models.UserPublicInfo, error)

	PagesCount(ctx context.Context,
		perPage int) (int, error)
}

type RecordUseCaseI interface {
	Update(ctx context.Context,
		id int32, record *models.Record) error
}

type ImageUseCaseI interface {
	Update(ctx context.Context,
		filename string, userID int32) error
	FetchByName(ctx context.Context,
		name string) (string, error)
	FetchByID(ctx context.Context,
		id int32) (string, error)
}
