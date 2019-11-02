package engine

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// ConnectionJSON is a wrapper for sending Connection by JSON
//easyjson:json
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
	return conn.JSON().MarshalJSON()
}

// UnmarshalJSON - overriding the standard method json.Unmarshal
func (conn *Connection) UnmarshalJSON(b []byte) error {
	temp := &ConnectionJSON{}

	if err := temp.UnmarshalJSON(b); err != nil {
		return err
	}

	conn._disconnected = temp.Disconnected
	conn._index = temp.Index
	conn.User = temp.User

	return nil
}
