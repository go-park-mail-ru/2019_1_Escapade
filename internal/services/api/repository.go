package api

import (
	"context"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
)

//go:generate $GOPATH/bin/mockery -name "UserRepositoryI|RecordRepositoryI|ImageRepositoryI"

// UsersSelectParams parameters to select user
type UsersSelectParams struct {
	Difficult int
	Offset    int
	Limit     int
	Sort      string
}

type UserRepositoryI interface {
	Create(ctx context.Context,
		user *models.UserPrivateInfo) (int, error)
	Delete(ctx context.Context,
		user *models.UserPrivateInfo) error

	UpdateNamePassword(ctx context.Context,
		user *models.UserPrivateInfo) error
	CheckNamePassword(ctx context.Context,
		name string, password string) (int32, *models.UserPublicInfo, error)
	FetchNamePassword(ctx context.Context,
		userID int32) (*models.UserPrivateInfo, error)

	UpdateLastSeen(ctx context.Context, id int) error

	FetchAll(ctx context.Context,
		params UsersSelectParams) ([]*models.UserPublicInfo, error)
	FetchOne(ctx context.Context,
		userID int32, difficult int) (*models.UserPublicInfo, error)

	PagesCount(ctx context.Context, perPage int) (int, error)
}

type RecordRepositoryI interface {
	Create(ctx context.Context, id int) error
	Update(ctx context.Context,
		id int32, record *models.Record) error
}

type ImageRepositoryI interface {
	Update(ctx context.Context,
		filename string, userID int32) error
	FetchByName(ctx context.Context,
		name string) (string, error)
	FetchByID(ctx context.Context, id int32) (string, error)
}
