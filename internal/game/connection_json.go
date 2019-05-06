package game

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/gorilla/websocket"
)

// Connection is a websocket of a player, that belongs to room
type Connection struct {
	wGroup *sync.WaitGroup

	doneM *sync.RWMutex
	_done bool

	roomM *sync.RWMutex
	_room *Room

	disconnectedM *sync.RWMutex
	_Disconnected bool

	bothM *sync.RWMutex
	_both bool

	indexM *sync.RWMutex
	_Index int

	User *models.UserPublicInfo

	ws    *websocket.Conn
	lobby *Lobby

	context context.Context
	cancel  context.CancelFunc

	actionSem chan struct{}

	send chan []byte
}

type ConnectionJSON struct {
	Disconnected bool `json:"disconnected"`
	Index        int  `json:"index"`

	User *models.UserPublicInfo `json:"user,omitempty"`
}

func (conn *Connection) JSON() ConnectionJSON {
	return ConnectionJSON{
		Disconnected: conn.Disconnected(),
		Index:        conn.Index(),
		User:         conn.User,
	}
}

func (conn *Connection) MarshalJSON() ([]byte, error) {
	return json.Marshal(conn.JSON())
}

func (conn *Connection) UnmarshalJSON(b []byte) error {
	temp := &ConnectionJSON{}

	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	conn._Disconnected = temp.Disconnected
	conn._Index = temp.Index
	conn.User = temp.User

	return nil
}
