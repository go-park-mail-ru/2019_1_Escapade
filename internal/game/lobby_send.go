package game

import "github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

func (lobby *Lobby) greet(conn *Connection) {
	response := models.Response{
		Type:  "Lobby",
		Value: lobby,
	}
	conn.SendInformation(response)
}

func (lobby *Lobby) send(info interface{}, predicate SendPredicate) {
	SendToConnections(info, predicate, lobby.Waiting.Get)
}

func (lobby *Lobby) sendRoomCreate(room Room, predicate SendPredicate) {
	response := models.Response{
		Type:  "LobbyRoomCreate",
		Value: room,
	}
	lobby.send(response, predicate)
}

func (lobby *Lobby) sendRoomUpdate(room Room, predicate SendPredicate) {
	response := models.Response{
		Type:  "LobbyRoomUpdate",
		Value: room,
	}
	lobby.send(response, predicate)
}

func (lobby *Lobby) sendRoomDelete(room Room, predicate SendPredicate) {
	response := models.Response{
		Type:  "LobbyRoomDelete",
		Value: room.ID,
	}
	lobby.send(response, predicate)
}

func (lobby *Lobby) sendWaiterEnter(conn Connection, predicate SendPredicate) {
	response := models.Response{
		Type:  "LobbyWaiterEnter",
		Value: conn,
	}
	lobby.send(response, predicate)
}

func (lobby *Lobby) sendWaiterExit(conn Connection, predicate SendPredicate) {
	response := models.Response{
		Type:  "LobbyWaiterExit",
		Value: conn,
	}
	lobby.send(response, predicate)
}

func (lobby *Lobby) sendPlayerEnter(conn Connection, predicate SendPredicate) {
	response := models.Response{
		Type:  "LobbyPlayerEnter",
		Value: conn,
	}
	lobby.send(response, predicate)
}

func (lobby *Lobby) sendPlayerExit(conn Connection, predicate SendPredicate) {
	response := models.Response{
		Type:  "LobbyPlayerExit",
		Value: conn,
	}
	lobby.send(response, predicate)
}
