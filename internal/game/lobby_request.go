package game

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// LobbyRequest - client send it by websocket to
// send/get information from Lobby
type LobbyRequest struct {
	Send    *LobbySend      `json:"send"`
	Message *models.Message `json:"message"`
	Get     *LobbyGet       `json:"get"`
}

// Invitation - client can invite everybody to room
type Invitation struct {
	From    *models.UserPublicInfo `json:"from"`
	Room    *Room                  `json:"room"`
	Message *models.Message        `json:"message"`
	To      string                 `json:"to"`
	All     bool                   `json:"all"`
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
	Invitation   *Invitation
	Messages     *models.Messages
}

// LobbyGet - Information, that client can get from lobby
type LobbyGet struct {
	AllRooms  bool `json:"allRooms"`
	FreeRooms bool `json:"freeRooms"`
	Waiting   bool `json:"waiting"`
	Playing   bool `json:"playing"`
}
