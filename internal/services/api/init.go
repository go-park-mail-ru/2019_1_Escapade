package api

import (
	"escapade/internal/config"
	"escapade/internal/database"
	"fmt"
	"time"
)

// Init creates Handler
func Init(DB *database.DataBase, c *config.Configuration) (handler *Handler) {
	ws := config.WebSocketSettings{
		WriteWait:      time.Duration(c.WebSocket.WriteWait) * time.Second,
		PongWait:       time.Duration(c.WebSocket.PongWait) * time.Second,
		PingPeriod:     time.Duration(c.WebSocket.PingPeriod) * time.Second,
		MaxMessageSize: c.WebSocket.MaxMessageSize,
	}
	handler = &Handler{
		DB:              *DB,
		Storage:         c.Storage,
		Cookie:          c.Cookie,
		GameConfig:      c.Game,
		WebSocket:       ws,
		WriteBufferSize: c.Server.WriteBufferSize,
		ReadBufferSize:  c.Server.ReadBufferSize,
	}
	return
}

// GetHandler return created handler with database and configuration
func GetHandler(confPath string) (handler *Handler,
	conf *config.Configuration, err error) {

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
