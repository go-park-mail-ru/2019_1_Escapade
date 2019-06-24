package game

import (
	"encoding/json"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

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
