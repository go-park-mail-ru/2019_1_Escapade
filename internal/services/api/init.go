package api

import (
	"escapade/internal/config"
	"escapade/internal/database"
	"fmt"
	"time"

	clients "escapade/internal/clients"

	"google.golang.org/grpc"
)

// Init creates Handler
func Init(DB *database.DataBase, c *config.Configuration, cl *clients.Clients) (handler *Handler) {

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
		AWS:             c.AWS,
		WebSocket:       ws,
		WriteBufferSize: c.Server.WriteBufferSize,
		ReadBufferSize:  c.Server.ReadBufferSize,
		Clients:         cl,
	}
	return
}

// GetHandler return created handler with database and configuration
func GetHandler(confPath, secretPath string, authConn *grpc.ClientConn) (handler *Handler,
	conf *config.Configuration, err error) {

	var (
		db          *database.DataBase
		servClients *clients.Clients
	)

	if conf, err = config.Init(confPath, secretPath); err != nil {
		return
	}
	fmt.Println("confPath done")

	if db, err = database.Init(conf.DataBase); err != nil {
		return
	}
	fmt.Println("database done")
	servClients = clients.Init(authConn)
	fmt.Println("clients done")
	handler = Init(db, conf, servClients)
	return
}
