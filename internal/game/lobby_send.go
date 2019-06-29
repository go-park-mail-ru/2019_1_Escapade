package game

import (
	"sync"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// all senders functions should add 1 to waitGroup!
// also all thay should be launched in goroutines and
// recover panic

func (lobby *Lobby) send(info interface{}, predicate SendPredicate) {
	SendToConnections(info, predicate, lobby.Waiting)
}

func (lobby *Lobby) sendToAll(info interface{}, predicate SendPredicate) {
	SendToConnections(info, predicate, lobby.Waiting, lobby.Playing)
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
		Type: "Lobby",
		Value: struct {
			Lobby Lobby                 `json:"lobby"`
			You   models.UserPublicInfo `json:"you"`
			Room  *Room                 `json:"room,omitempty"`
		}{
			Lobby: *lobby,
			You:   *conn.User,
			Room:  conn.WaitingRoom(),
		},
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

func (lobby *Lobby) sendRoomCreate(room *Room, predicate SendPredicate, group *sync.WaitGroup) {
	defer group.Done()
	defer utils.CatchPanic("lobby_send.go sendRoomCreate")

	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer lobby.wGroup.Done()

	response := models.Response{
		Type:  "LobbyRoomCreate",
		Value: room.JSON(),
	}
	lobby.send(response, predicate)
}

func (lobby *Lobby) sendRoomUpdate(room *Room, predicate SendPredicate, group *sync.WaitGroup) {

	utils.Debug(false, "here")
	defer group.Done()
	utils.Debug(false, "we")
	defer utils.CatchPanic("lobby_send.go sendRoomUpdate")

	if lobby.done() {
		utils.Debug(false, "aaaaaa")
		return
	}
	utils.Debug(false, "wtf")
	lobby.wGroup.Add(1)
	defer lobby.wGroup.Done()

	response := models.Response{
		Type:  "LobbyRoomUpdate",
		Value: room.JSON(),
	}
	utils.Debug(false, "why")
	lobby.send(response, predicate)
}

func (lobby *Lobby) sendRoomDelete(roomID string, predicate SendPredicate) {
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
		Value: roomID,
	}
	lobby.send(response, predicate)
}

func (lobby *Lobby) sendWaiterEnter(conn *Connection, predicate SendPredicate) {
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

func (lobby *Lobby) sendWaiterExit(conn *Connection, predicate SendPredicate) {
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

func (lobby *Lobby) sendPlayerEnter(conn *Connection, predicate SendPredicate) {
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

func (lobby *Lobby) sendPlayerExit(conn *Connection, predicate SendPredicate) {
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

func (lobby *Lobby) sendInvitation(inv *Invitation, predicate SendPredicate) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("lobby sendInvitation")
	}()

	response := models.Response{
		Type:  "LobbyInvitation",
		Value: inv,
	}
	lobby.send(response, predicate)
}

func (lobby *Lobby) sendInvitationCallback(conn *Connection, err error) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("lobby sendCallback")
	}()

	response := models.Response{
		Type:  "LobbyInvitationCallback",
		Value: err,
	}
	conn.SendInformation(response)
}
