package game

func (lobby *Lobby) greet(conn *Connection) {
	conn.SendInformation(lobby)
}

func (lobby *Lobby) send(info interface{}, predicate SendPredicate) {
	SendToConnections(info, predicate, lobby.Waiting.Get)
}

// SendMessage sends message to Connection from lobby
func (lobby *Lobby) SendMessage(conn *Connection, message string) {
	conn.SendInformation("Lobby message: " + message)
}

// send to all in lobby
func (lobby *Lobby) sendAllRooms(predicate SendPredicate) {
	get := &LobbyGet{
		AllRooms: true,
	}
	model := lobby.makeGetModel(get)
	lobby.send(model, predicate)
}

func (lobby *Lobby) sendWaiting(predicate SendPredicate) {
	get := &LobbyGet{
		Waiting: true,
	}
	model := lobby.makeGetModel(get)
	lobby.send(model, predicate)
}

func (lobby *Lobby) sendPlaying(predicate SendPredicate) {
	get := &LobbyGet{
		Waiting: true,
		Playing: true,
	}
	model := lobby.makeGetModel(get)
	lobby.send(model, predicate)
}

func (lobby *Lobby) makeGetModel(get *LobbyGet) *Lobby {
	sendLobby := &Lobby{}
	sendLobby.Type = lobby.Type
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
