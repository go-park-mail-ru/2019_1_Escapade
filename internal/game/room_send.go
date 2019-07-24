package game

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// sendToAllInRoom send info to those in room, whose predicate
// returns true
func (room *Room) send(info utils.JSONtype, predicate SendPredicate) {
	players := room.Players.Connections
	observers := room.Observers
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

	room.send(models.Result{
		Message: "Room(" + room.ID() + "):" + text}, predicate)
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
func (room *Room) sendGameOver(timer bool, predicate SendPredicate,
	cells []Cell, wg *sync.WaitGroup) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
		utils.CatchPanic("room_send.go sendGameOver()")
	}()

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
func (room *Room) sendPlayerEnter(conn *Connection, predicate SendPredicate) {
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
func (room *Room) sendPlayerExit(conn *Connection, predicate SendPredicate) {
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
func (room *Room) sendObserverEnter(conn *Connection, predicate SendPredicate) {
	response := models.Response{
		Type:  "RoomObserverEnter",
		Value: conn,
	}
	room.send(response, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *Room) sendObserverExit(conn *Connection, predicate SendPredicate) {
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
func (room *Room) sendStatus(predicate SendPredicate, status int, wg *sync.WaitGroup) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
		utils.CatchPanic("room_send.go RoomStatus()")
	}()

	var leftTime int
	since := int(time.Since(room.Date()).Seconds())
	if status == StatusFlagPlacing {
		leftTime = room.Settings.TimeToPrepare - since
	}
	if status == StatusRunning {
		leftTime = room.Settings.TimeToPlay - since
	}
	response := models.Response{
		Type: "RoomStatus",
		Value: struct {
			ID     string `json:"id"`
			Status int    `json:"status"`
			Time   int    `json:"time"`
		}{
			ID:     room.ID(),
			Status: status,
			Time:   leftTime,
		},
	}
	fmt.Println("!!!!!!!leftTime ", leftTime)
	room.send(response, predicate)
}

func (room *Room) sendStatusOne(conn *Connection) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
		utils.CatchPanic("room_send.go RoomStatus()")
	}()

	var leftTime int
	status := room.Status()
	since := int(time.Since(room.Date()).Seconds())
	if status == StatusFlagPlacing {
		leftTime = room.Settings.TimeToPrepare - since
	}
	if status == StatusRunning {
		leftTime = room.Settings.TimeToPlay - since
	}
	response := models.Response{
		Type: "RoomStatus",
		Value: struct {
			ID     string `json:"id"`
			Status int    `json:"status"`
			Time   int    `json:"time"`
		}{
			ID:     room.ID(),
			Status: status,
			Time:   leftTime,
		},
	}
	fmt.Println("leftTime ", leftTime)
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
func (room *Room) sendError(err error, conn *Connection) {
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
	if conn.done() {
		return
	}
	conn.wGroup.Add(1)
	defer conn.wGroup.Done()

	var flag Flag
	if room.Settings.Deathmatch {
		index := conn.Index()
		if index >= 0 {
			flag = room.Players.Flag(index)
		}
	} else {
		flag = Flag{Cell: *NewCell(-1, -1, 0, 0)}
	}

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
			You:  *conn.User,
			Flag: flag,
			//Time:     leftTime,
			IsPlayer: isPlayer,
		},
	}
	fmt.Println("room send to ", conn.ID())
	conn.SendInformation(response)
}
