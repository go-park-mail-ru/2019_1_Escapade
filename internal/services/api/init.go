package api

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/clients"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"

	"fmt"
	"google.golang.org/grpc"
	"time"
)

// Init creates Handler
func Init(DB *database.DataBase, c *config.Configuration, authConn *grpc.ClientConn) (handler *Handler) {
	ws := config.WebSocketSettings{
		WriteWait:      time.Duration(c.WebSocket.WriteWait) * time.Second,
		PongWait:       time.Duration(c.WebSocket.PongWait) * time.Second,
		PingPeriod:     time.Duration(c.WebSocket.PingPeriod) * time.Second,
		MaxMessageSize: c.WebSocket.MaxMessageSize,
	}
	handler = &Handler{
		DB:              *DB,
		Storage:         c.Storage,
		Session:         c.Session,
		GameConfig:      c.Game,
		AWS:             c.AWS,
		WebSocket:       ws,
		Clients:         clients.Init(authConn),
		WriteBufferSize: c.Server.WriteBufferSize,
		ReadBufferSize:  c.Server.ReadBufferSize,
	}
	return
}

// GetHandler return created handler with database and configuration
func InitAPI(confPath, secretPath string, authConn *grpc.ClientConn) (H *Handler, C *config.Configuration, err error) {

	var (
		db *database.DataBase
	)

	fmt.Println("InitAPI")
	if C, err = config.Init(confPath, secretPath); err != nil {
		fmt.Println("eeeer", err.Error())
		return
	}
	fmt.Println("ConfigurationAPI done")

	if db, err = database.Init(C.DataBase); err != nil {
		return
	}

	fmt.Println("Database done")
	H = Init(db, C, authConn)
	return
}
