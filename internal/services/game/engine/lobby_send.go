package engine

import (
	handlers "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// all senders functions should add 1 to waitGroup!
// also all thay should be launched in goroutines and
// recover panic

func (lobby *Lobby) send(info handlers.JSONtype, predicate SendPredicate) {
	lobby.s.Do(func() {
		SendToConnections(info, predicate, lobby.Waiting)
	})
}

func (lobby *Lobby) sendToAll(info handlers.JSONtype, predicate SendPredicate) {
	lobby.s.Do(func() {
		SendToConnections(info, predicate, lobby.Waiting, lobby.Playing)
	})
}

func (lobby *Lobby) greet(conn *Connection) {
	lobby.s.DoWithOther(conn, func() {
		response := models.Response{
			Type: "Lobby",
			Value: struct {
				Lobby LobbyJSON             `json:"lobby"`
				You   models.UserPublicInfo `json:"you"`
				Room  *Room                 `json:"room,omitempty"`
			}{
				Lobby: lobby.JSON(),
				You:   *conn.User,
				Room:  conn.WaitingRoom(),
			},
		}
		conn.SendInformation(&response)
	})
}

func (lobby *Lobby) sendLobbyMessage(message string, predicate SendPredicate) {
	response := models.Response{
		Type:  "LobbyMessage",
		Value: message,
	}
	lobby.sendToAll(&response, predicate)
}

func (lobby *Lobby) sendRoomCreate(room *Room, predicate SendPredicate) {
	room.sync.Do(func() {
		response := models.Response{
			Type:  "LobbyRoomCreate",
			Value: room.models.JSON(),
		}
		lobby.send(&response, predicate)
	})
}

func (lobby *Lobby) sendRoomUpdate(room *Room, predicate SendPredicate) {
	room.sync.Do(func() {
		response := models.Response{
			Type:  "LobbyRoomUpdate",
			Value: room.models.JSON(),
		}
		lobby.send(&response, predicate)
	})
}

func (lobby *Lobby) sendRoomDelete(roomID string, predicate SendPredicate) {
	response := models.Response{
		Type:  "LobbyRoomDelete",
		Value: roomID,
	}
	lobby.send(&response, predicate)
}

func (lobby *Lobby) sendWaiterEnter(conn *Connection, predicate SendPredicate) {
	response := models.Response{
		Type:  "LobbyWaiterEnter",
		Value: conn,
	}
	lobby.send(&response, predicate)
}

func (lobby *Lobby) sendWaiterExit(conn *Connection, predicate SendPredicate) {
	response := models.Response{
		Type:  "LobbyWaiterExit",
		Value: conn,
	}
	lobby.send(&response, predicate)
}

func (lobby *Lobby) sendPlayerEnter(conn *Connection, predicate SendPredicate) {
	response := models.Response{
		Type:  "LobbyPlayerEnter",
		Value: conn,
	}
	lobby.send(&response, predicate)
}

func (lobby *Lobby) sendPlayerExit(conn *Connection, predicate SendPredicate) {
	response := models.Response{
		Type:  "LobbyPlayerExit",
		Value: conn,
	}
	lobby.send(&response, predicate)
}

func (lobby *Lobby) sendInvitation(inv *Invitation, predicate SendPredicate) {
	response := models.Response{
		Type:  "LobbyInvitation",
		Value: inv,
	}
	lobby.send(&response, predicate)
}

func (lobby *Lobby) sendInvitationCallback(conn *Connection, err error) {
	lobby.s.DoWithOther(conn, func() {
		response := models.Response{
			Type:  "LobbyInvitationCallback",
			Value: err,
		}
		conn.SendInformation(&response)
	})
}
