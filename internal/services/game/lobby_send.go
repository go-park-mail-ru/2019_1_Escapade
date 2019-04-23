package game

import (
	"encoding/json"
)

func (lobby *Lobby) greet(conn *Connection) {
	bytes, _ := json.Marshal(lobby)
	conn.SendInformation(bytes)
}

func (lobby *Lobby) sendToWaiters(info interface{}, predicate SendPredicate) {
	SendToConnections(info, predicate, lobby.Waiting.Get)
}

func (lobby *Lobby) sendToAllInLobby(info interface{}) {
	SendToConnections(info, All(), lobby.Waiting.Get)
}

// send to all in lobby
func (lobby *Lobby) sendTAILRooms() {
	get := &LobbyGet{
		AllRooms:  true,
		FreeRooms: true,
	}
	send := lobby.makeGetModel(get)
	lobby.sendToAllInLobby(send)
}

func (lobby *Lobby) sendTAILPeople() {
	get := &LobbyGet{
		Waiting: true,
		Playing: true,
	}
	send := lobby.makeGetModel(get)
	lobby.sendToAllInLobby(send)
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
