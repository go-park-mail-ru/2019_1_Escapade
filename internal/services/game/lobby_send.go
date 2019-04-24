package game

import (
	"encoding/json"
)

func (lobby *Lobby) greet(conn *Connection) {
	bytes, _ := json.Marshal(lobby)
	conn.SendInformation(bytes)
}

func (lobby *Lobby) send(info interface{}, predicate SendPredicate) {
	SendToConnections(info, predicate, lobby.Waiting.Get)
}

// SendMessage sends message to Connection from lobby
func (lobby *Lobby) SendMessage(conn *Connection, message string) {
	conn.SendInformation([]byte("Lobby message: " + message))
}

// send to all in lobby
func (lobby *Lobby) sendTAILRooms() {
	get := &LobbyGet{
		AllRooms:  true,
		FreeRooms: true,
	}
	model := lobby.makeGetModel(get)
	lobby.send(model, All)
}

func (lobby *Lobby) sendTAILPeople() {
	get := &LobbyGet{
		Waiting: true,
		Playing: true,
	}
	model := lobby.makeGetModel(get)
	lobby.send(model, All)
}

func (lobby *Lobby) makeGetModel(get *LobbyGet) *Lobby {
	sendLobby := &Lobby{}
	if get.AllRooms {
		sendLobby.AllRooms = lobby.AllRooms
	}
	if get.FreeRooms {
		sendLobby.FreeRooms = lobby.FreeRooms
	}
	if get.Waiting {
		sendLobby.Waiting = lobby.Waiting
	}
	if get.Playing {
		sendLobby.Playing = lobby.Playing
	}
	return sendLobby
}
