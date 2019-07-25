package api

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/constants"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"

	"time"
)

// Init creates Handler
func Init(DB *database.DataBase, c *config.Configuration /*, authConn *grpc.ClientConn*/) (handler *Handler) {
	ws := config.WebSocketSettings{
		WriteWait:       time.Duration(c.WebSocket.WriteWait) * time.Second,
		PongWait:        time.Duration(c.WebSocket.PongWait) * time.Second,
		PingPeriod:      time.Duration(c.WebSocket.PingPeriod) * time.Second,
		MaxMessageSize:  c.WebSocket.MaxMessageSize,
		ReadBufferSize:  c.WebSocket.ReadBufferSize,
		WriteBufferSize: c.WebSocket.WriteBufferSize,
	}
	handler = &Handler{
		DB:        *DB,
		Session:   c.Session,
		Game:      c.Game,
		WebSocket: ws,
		//Clients:   clients.Init(authConn),
	}
	constants.InitField()
	constants.InitRoom()
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
