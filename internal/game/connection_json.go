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
	_disconnected bool

	bothM *sync.RWMutex
	_both bool

	indexM *sync.RWMutex
	_index int

	User *models.UserPublicInfo

	ws    *websocket.Conn
	lobby *Lobby

	context context.Context
	cancel  context.CancelFunc

	actionSem chan struct{}

	send chan []byte
}

// ConnectionJSON is a wrapper for sending Connection by JSON
type ConnectionJSON struct {
	Disconnected bool `json:"disconnected"`
	Index        int  `json:"index"`

	User *models.UserPublicInfo `json:"user,omitempty"`
}

// JSON convert Connection to ConnectionJSON
func (conn *Connection) JSON() ConnectionJSON {
	return ConnectionJSON{
		Disconnected: conn.Disconnected(),
		Index:        conn.Index(),
		User:         conn.User,
	}
}

// MarshalJSON - overriding the standard method json.Marshal
func (conn *Connection) MarshalJSON() ([]byte, error) {
	return json.Marshal(conn.JSON())
}

// UnmarshalJSON - overriding the standard method json.Unmarshal
func (conn *Connection) UnmarshalJSON(b []byte) error {
	temp := &ConnectionJSON{}

	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	conn._disconnected = temp.Disconnected
	conn._index = temp.Index
	conn.User = temp.User

	return nil
}
