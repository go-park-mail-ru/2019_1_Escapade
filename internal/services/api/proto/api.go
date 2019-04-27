package apiproto

import (
	"context"
	"escapade/internal/services/api"
)

type ApiManager struct{}

func NewApiManager() *ApiManager {
	return &ApiManager{}
}

func (am *ApiManager) User(ctx context.Context, userID int32) (um *UserModel, err error) {

	user, _ := api.API.DB.GetUser(int(userID), 0)
	um = &UserModel{
		ID:       int32(userID),
		Name:     user.Name,
		PhotoURL: user.PhotoURL,
	}
	return
}
