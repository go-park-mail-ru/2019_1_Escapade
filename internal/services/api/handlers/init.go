package handlers

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
)

// Init creates Handler
func Init(DB *database.DataBase, c *config.Configuration) (handler *Handler) {
	handler = &Handler{
		DB:         *DB,
		Cookie:     c.Cookie,
		AuthClient: c.AuthClient,
		Auth:       c.Auth,
	}
	/*
		oauth2.Config{
				ClientID:     "1",
				ClientSecret: "1",
				Scopes:       []string{"all"},
				RedirectURL:  "http://auth:3001/oauth2",
				Endpoint: oauth2.Endpoint{
					AuthURL:  auth.AuthServerURL + "/authorize",
					TokenURL: auth.AuthServerURL + "/token",
				},
	*/
	return
}

// GetHandler init handler and configuration for api service
func GetHandler(C *config.Configuration /*, authConn *grpc.ClientConn*/) (H *Handler, err error) {

	var (
		db *database.DataBase
	)

	if db, err = database.Init(C.DataBase); err != nil {
		return
	}

	H = Init(db, C /*authConn*/)
	return
}
