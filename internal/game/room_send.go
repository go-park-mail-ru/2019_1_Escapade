package game

import (
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// sendToAllInRoom send info to those in room, whose predicate
// returns true
func (room *Room) send(info interface{}, predicate SendPredicate) {
	players := room.Players.Connections.RGet()
	observers := room.Observers.RGet()
	SendToConnections(info, predicate, players, observers)
}

func (room *Room) sendMessage(text string, predicate SendPredicate) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
		utils.CatchPanic("room_send.go sendMessage()")
	}()

	room.send("Room("+room.ID+"):"+text, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *Room) sendPlayerPoints(player Player, predicate SendPredicate) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
		utils.CatchPanic("room_send.go sendPlayerPoints()")
	}()

	response := models.Response{
		Type:  "RoomPlayerPoints",
		Value: player,
	}
	room.send(response, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *Room) sendGameOver(timer bool, predicate SendPredicate) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
		utils.CatchPanic("room_send.go sendGameOver()")
	}()

	cells := make([]Cell, 0)
	room.Field.OpenEverything(&cells)
	response := models.Response{
		Type: "RoomGameOver",
		Value: struct {
			Players []Player `json:"players"`
			Cells   []Cell   `json:"cells"`
			Winners []int    `json:"winners"`
			Timer   bool     `json:"timer"`
		}{
			Players: room.Players.RPlayers(),
			Cells:   cells,
			Winners: room.Winners(),
			Timer:   timer,
		},
	}
	room.send(response, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func sendAccountTaken(conn Connection) {

	response := models.Response{
		Type: "AccountTaken",
	}
	conn.SendInformation(response)
}

func (room *Room) sendNewCells(predicate SendPredicate, cells ...Cell) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
		utils.CatchPanic("room_send.go RoomNewCells()")
	}()

	response := models.Response{
		Type:  "RoomNewCells",
		Value: cells,
	}
	room.send(response, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *Room) sendPlayerEnter(conn Connection, predicate SendPredicate) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
		utils.CatchPanic("room_send.go RoomPlayerEnter()")
	}()

	response := models.Response{
		Type:  "RoomPlayerEnter",
		Value: conn,
	}
	room.send(response, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *Room) sendPlayerExit(conn Connection, predicate SendPredicate) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
		utils.CatchPanic("room_send.go RoomPlayerExit()")
	}()

	response := models.Response{
		Type:  "RoomPlayerExit",
		Value: conn,
	}
	room.send(response, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *Room) sendObserverEnter(conn Connection, predicate SendPredicate) {
	response := models.Response{
		Type:  "RoomObserverEnter",
		Value: conn,
	}
	room.send(response, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *Room) sendObserverExit(conn Connection, predicate SendPredicate) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
		utils.CatchPanic("room_send.go RoomObserverExit()")
	}()

	response := models.Response{
		Type:  "RoomObserverExit",
		Value: conn,
	}
	room.send(response, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *Room) sendStatus(predicate SendPredicate) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
		utils.CatchPanic("room_send.go RoomStatus()")
	}()

	var leftTime int
	if room.Status == StatusFlagPlacing {
		leftTime = room.Settings.TimeToPrepare - int(time.Since(room.Date).Seconds())
	}
	if room.Status == StatusRunning {
		leftTime = room.Settings.TimeToPlay - int(time.Since(room.Date).Seconds())
	}
	response := models.Response{
		Type: "RoomStatus",
		Value: struct {
			ID     string `json:"id"`
			Status int    `json:"status"`
			Time   int    `json:"time"`
		}{
			ID:     room.ID,
			Status: room.Status,
			Time:   leftTime,
		},
	}
	room.send(response, predicate)
}

func (room *Room) sendStatusOne(conn Connection) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
		utils.CatchPanic("room_send.go RoomStatus()")
	}()

	var leftTime int
	if room.Status == StatusFlagPlacing {
		leftTime = room.Settings.TimeToPrepare - int(time.Since(room.Date).Seconds())
	}
	if room.Status == StatusRunning {
		leftTime = room.Settings.TimeToPlay - int(time.Since(room.Date).Seconds())
	}
	response := models.Response{
		Type: "RoomStatus",
		Value: struct {
			ID     string `json:"id"`
			Status int    `json:"status"`
			Time   int    `json:"time"`
		}{
			ID:     room.ID,
			Status: room.Status,
			Time:   leftTime,
		},
	}
	fmt.Println("status send to ", conn.ID())
	conn.SendInformation(response)
}

// sendTAIRHistory send actions history to all in room
func (room *Room) sendAction(pa PlayerAction, predicate SendPredicate) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
		utils.CatchPanic("room_send.go sendAction()")
	}()

	response := models.Response{
		Type:  "RoomAction",
		Value: pa,
	}
	room.send(response, predicate)
}

// sendTAIRHistory send actions history to all in room
func (room *Room) sendError(err error, conn Connection) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
		utils.CatchPanic("room_send.go sendError()")
	}()

	response := models.Response{
		Type:  "RoomError",
		Value: err.Error(),
	}
	conn.SendInformation(response)
}

// sendTAIRField send field to all in room
func (room *Room) sendField(predicate SendPredicate) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
		utils.CatchPanic("room_send.go sendField()")
	}()

	response := models.Response{
		Type:  "RoomField",
		Value: room.Field,
	}
	room.send(response, predicate)
}

// sendTAIRAll send everything to one connection
func (room *Room) greet(conn *Connection, isPlayer bool) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
		utils.CatchPanic("room_send.go greet()")
	}()

	var flag Flag
	index := conn.Index()
	if index >= 0 {
		flag = room.Players.Flag(index)
	}

	copy := *conn

	//leftTime := room.Settings.TimeToPlay + room.Settings.TimeToPrepare - int(time.Since(room.Date).Seconds())

	response := models.Response{
		Type: "Room",
		Value: struct {
			Room *Room                 `json:"room"`
			You  models.UserPublicInfo `json:"you"`
			Flag Flag                  `json:"flag,omitempty"`
			//Time     int                   `json:"time"`
			IsPlayer bool `json:"isPlayer"`
		}{
			Room: room,
			You:  *copy.User,
			Flag: flag,
			//Time:     leftTime,
			IsPlayer: isPlayer,
		},
	}
	fmt.Println("room send to ", conn.ID())
	conn.SendInformation(response)
}
