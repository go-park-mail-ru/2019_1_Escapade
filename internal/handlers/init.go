package api

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	"golang.org/x/oauth2"
)

const (
	authServerURL = "http://localhost:9096"
)

var (
	oauth2Config = oauth2.Config{
		ClientID:     "222222",
		ClientSecret: "22222222",
		Scopes:       []string{"all"},
		RedirectURL:  "http://localhost:9094/oauth2",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authServerURL + "/authorize",
			TokenURL: authServerURL + "/token",
		},
	}
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
		Oauth: oauth2.Config{
			ClientID:     "1",
			ClientSecret: "1",
			Scopes:       []string{"all"},
			RedirectURL:  "http://localhost:3001/oauth2",
			Endpoint: oauth2.Endpoint{
				AuthURL:  authServerURL + "/authorize",
				TokenURL: authServerURL + "/token",
			},
		},
		//Clients:   clients.Init(authConn),
	}
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
