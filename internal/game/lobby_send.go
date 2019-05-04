package game

import "github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
import "github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

// all senders functions should add 1 to waitGroup!
// also all thay should be launched in goroutines and
// recover panic

func (lobby *Lobby) send(info interface{}, predicate SendPredicate) {
	SendToConnections(info, predicate, lobby.waiting())
}

func (lobby *Lobby) sendToAll(info interface{}, predicate SendPredicate) {
	SendToConnections(info, predicate, lobby.waiting(), lobby.playing())
}

func (lobby *Lobby) greet(conn *Connection) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("lobby greet")
	}()

	response := models.Response{
		Type:  "Lobby",
		Value: lobby,
	}
	conn.SendInformation(response)
}

func (lobby *Lobby) sendLobbyMessage(message string, predicate SendPredicate) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("lobby sendLobbyMessage")
	}()

	response := models.Response{
		Type:  "LobbyMessage",
		Value: message,
	}
	lobby.sendToAll(response, predicate)
}

func (lobby *Lobby) sendRoomCreate(room Room, predicate SendPredicate) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("lobby sendRoomCreate")
	}()

	response := models.Response{
		Type:  "LobbyRoomCreate",
		Value: room,
	}
	lobby.send(response, predicate)
}

func (lobby *Lobby) sendRoomUpdate(room Room, predicate SendPredicate) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("lobby sendRoomUpdate")
	}()

	response := models.Response{
		Type:  "LobbyRoomUpdate",
		Value: room,
	}
	lobby.send(response, predicate)
}

func (lobby *Lobby) sendRoomDelete(room Room, predicate SendPredicate) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("lobby sendRoomDelete")
	}()

	response := models.Response{
		Type:  "LobbyRoomDelete",
		Value: room.ID,
	}
	lobby.send(response, predicate)
}

func (lobby *Lobby) sendWaiterEnter(conn Connection, predicate SendPredicate) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("lobby sendWaiterEnter")
	}()

	response := models.Response{
		Type:  "LobbyWaiterEnter",
		Value: conn,
	}
	lobby.send(response, predicate)
}

func (lobby *Lobby) sendWaiterExit(conn Connection, predicate SendPredicate) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("lobby sendWaiterExit")
	}()

	response := models.Response{
		Type:  "LobbyWaiterExit",
		Value: conn,
	}
	lobby.send(response, predicate)
}

func (lobby *Lobby) sendPlayerEnter(conn Connection, predicate SendPredicate) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("lobby sendPlayerEnter")
	}()

	response := models.Response{
		Type:  "LobbyPlayerEnter",
		Value: conn,
	}
	lobby.send(response, predicate)
}

func (lobby *Lobby) sendPlayerExit(conn Connection, predicate SendPredicate) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("lobby sendPlayerExit")
	}()

	response := models.Response{
		Type:  "LobbyPlayerExit",
		Value: conn,
	}
	lobby.send(response, predicate)
}