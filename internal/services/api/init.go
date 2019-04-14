package api

import (
	"escapade/internal/config"
	"escapade/internal/database"
	"escapade/internal/services/game"
	"fmt"
)

// Init creates Handler
func Init(DB *database.DataBase, config *config.Configuration) (handler *Handler) {
	lobby := game.NewLobby(config.Game.RoomsCapacity,
		config.Game.LobbyJoin, config.Game.LobbyRequest)
	handler = &Handler{
		DB:                    *DB,
		PlayersAvatarsStorage: config.Storage.PlayersAvatarsStorage,
		FileMode:              config.Storage.FileMode,
		WriteBufferSize:       config.Server.WriteBufferSize,
		ReadBufferSize:        config.Server.ReadBufferSize,
		Lobby:                 lobby,
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
	fmt.Println("confPath done")

	if db, err = database.Init(conf.DataBase); err != nil {
		return
	}

	fmt.Println("database done")
	handler = Init(db, conf)
	return
}
