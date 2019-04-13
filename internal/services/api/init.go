package api

import (
	"escapade/internal/config"
	"escapade/internal/database"
	"escapade/internal/services/game"
)

// Init creates Handler
func Init(DB *database.DataBase, config *config.Configuration) (handler *Handler) {
	lobby := game.NewLobby(config.Game.RoomsCapacity,
		config.Game.LobbyJoin, config.Game.LobbyRequest)
	handler = &Handler{
		DB:              *DB,
		Storage:         config.Storage,
		WriteBufferSize: config.Server.WriteBufferSize,
		ReadBufferSize:  config.Server.ReadBufferSize,
		Lobby:           lobby,
	}
	go handler.Lobby.Run()
	return
}

func GetHandler(confPath string) (handler *Handler, conf *config.Configuration, err error) {

	var (
		db *database.DataBase
	)

	if conf, err = config.Init(confPath); err != nil {
		return
	}

	if db, err = database.Init(conf.DataBase); err != nil {
		return
	}

	handler = Init(db, conf)
	return
}
