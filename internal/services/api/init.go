package api

import (
	"escapade/internal/config"
	"escapade/internal/database"
	//"reflect"
)

// Init creates Handler
func Init(DB *database.DataBase, storage config.FileStorageConfig, server config.ServerConfig) (handler *Handler) {
	handler = &Handler{
		DB:                    *DB,
		PlayersAvatarsStorage: storage.PlayersAvatarsStorage,
		FileMode:              storage.FileMode,
		WriteBufferSize:       server.WriteBufferSize,
		ReadBufferSize:        server.ReadBufferSize,
	}
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

	handler = Init(db, conf.Storage, conf.Server)
	return
}
