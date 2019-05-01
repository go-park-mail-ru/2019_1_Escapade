package game

import (
	"escapade/internal/models"
)

// LobbyRequest - client send it by websocket to
// send/get information from Lobby
type LobbyRequest struct {
	Send *LobbySend `json:"send"`
	Get  *LobbyGet  `json:"get"`
}

// NewLobbyRequest creates Lobby instance
func NewLobbyRequest(s *LobbySend, g *LobbyGet) *LobbyRequest {
	return &LobbyRequest{
		Send: s,
		Get:  g,
	}
}

// IsGet checks, if client want get info
func (lr *LobbyRequest) IsGet() bool {
	return lr.Get != nil
}

// IsSend checks, if client want send info
func (lr *LobbyRequest) IsSend() bool {
	return lr.Send != nil
}

// LobbySend - Information, that client can send to lobby
type LobbySend struct {
	RoomSettings *models.RoomSettings
}

// LobbyGet - Information, that client can get from lobby
type LobbyGet struct {
	AllRooms  bool `json:"allRooms"`
	FreeRooms bool `json:"freeRooms"`
	Waiting   bool `json:"waiting"`
	Playing   bool `json:"playing"`
}
