package engine

import (
	"sync"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/metrics"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"time"
)

// Enter handle user joining as player or observer
func (room *Room) Enter(conn *Connection) bool {
	if room.done() {
		utils.Debug(true, "room is done")
		return false
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	var done bool

	// if room is searching new players
	if room.Status() == StatusRecruitment {
		if room.addConnection(conn, true, false) {
			utils.Debug(false, "You will be player!")
			done = true
		}
	} else if room.addConnection(conn, false, false) {
		utils.Debug(false, "You will be observer!")
		done = true
	}
	return done
}

// Free clear all resources. Call it when no
//  observers and players inside
func (room *Room) Free() {

	if room.checkAndSetCleared() {
		return
	}

	groupWaitRoom := 60 * time.Second // TODO в конфиг
	fieldWaitRoom := 40 * time.Second // TODO в конфиг
	utils.WaitWithTimeout(room.wGroup, groupWaitRoom)

	room.chanStatus <- StatusAborted

	room.setStatus(StatusFinished)
	go room.historyFree()
	go room.messagesFree()
	go room.Players.Free()
	go room.Observers.Free()
	go room.Field.Free(fieldWaitRoom)

	close(room.chanStatus)
	close(room.chanConnection)
}

// Close drives away players out of the room, free resources
// and inform lobby, that rooms closes
func (room *Room) Close() bool {
	if room.done() {
		return false
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	utils.Debug(false, "Can close?", room.lobby.config().CanClose)
	if !room.lobby.config().CanClose {
		return false
	}
	utils.Debug(false, "We closed room :С")

	room.wGroup.Add(1)
	go room.lobby.CloseRoom(room, room.wGroup)

	utils.Debug(false, "Prepare to free!")
	go room.Free()
	utils.Debug(false, "We did it")
	return true
}

// LeaveAll make every room connection to leave
func (room *Room) LeaveAll() {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	playersIterator := NewConnectionsIterator(room.Players.Connections)
	for playersIterator.Next() {
		player := playersIterator.Value()
		go room.Leave(player, true)
	}

	observersIterator := NewConnectionsIterator(room.Observers)
	for observersIterator.Next() {
		observer := observersIterator.Value()
		go room.Leave(observer, false)
	}
}

// Leave handle player going back to lobby
func (room *Room) Leave(conn *Connection, isPlayer bool) bool {
	if room.done() {
		return false
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	var (
		done   bool
		status = room.Status()
	)
	if isPlayer {
		if status == StatusRecruitment || status == StatusFinished {
			done = room.Players.Connections.Remove(conn)
		}
		utils.Debug(false, "isPlayer delete")
	} else {
		done = room.Observers.Remove(conn)
		utils.Debug(false, "notPlayer delete")
	}
	utils.Debug(false, "done", done)
	if !done {
		utils.Debug(false, "not found", conn.ID())
		return false
	}

	if status == StatusRecruitment {
		room.lobby.greet(conn)

		room.wGroup.Add(1)
		go room.lobby.sendRoomUpdate(room, All, room.wGroup)
	} else if isPlayer {
		go room.GiveUp(conn)
	}

	utils.Debug(false, "letsCheckIfNil", room.Players.Connections.len(), room.Observers.len())
	if room.Empty() {
		room.Close()
	}
	return true
}

// applyAction applies the effects of opening a cell
func (room *Room) applyAction(conn *Connection, cell *Cell) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	index := conn.Index()

	switch {
	case cell.Value < CellMine:
		room.Players.m.IncreasePlayerPoints(index, 1000*float64(cell.Value)/float64(room.Settings.Width*room.Settings.Height))
	case cell.Value == CellMine:
		room.Players.m.IncreasePlayerPoints(index, float64(-1000))
		room.Kill(conn, ActionExplode)
	case cell.Value > CellIncrement:
		room.FlagFound(*conn, cell)
	}
}

// OpenCell open cell
func (room *Room) OpenCell(conn *Connection, cell *Cell, group *sync.WaitGroup) {
	defer group.Done()

	// if user try set open cell before game launch
	if room.Status() != StatusRunning {
		return
	}

	// if wrong cell
	if !room.Field.IsInside(cell) {
		return
	}

	// if user died
	if !room.isAlive(conn) {
		return
	}

	// set who try open cell(for history)
	cell.PlayerID = conn.ID()
	cells := room.Field.OpenCell(cell)
	if len(cells) == 1 {
		newCell := cells[0]
		room.applyAction(conn, &newCell)
	} else {
		for _, foundCell := range cells {
			room.applyAction(conn, &foundCell)
		}
	}

	if len(cells) > 0 {
		go room.send.PlayerPoints(room.Players.m.Player(conn.Index()), room.All)
		go room.send.NewCells(room.All, cells...)
	}
	if room.Field.IsCleared() {
		room.updateStatus(StatusFinished)
	}
	return
}

// CellHandle processes the Cell came from the user
func (room *Room) CellHandle(conn *Connection, cell *Cell) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	utils.Debug(false, "cellHandle")
	status := room.Status()
	if status == StatusFlagPlacing {
		room.wGroup.Add(1)
		room.SetFlag(conn, cell, room.wGroup)
	} else if status == StatusRunning {
		room.wGroup.Add(1)
		room.OpenCell(conn, cell, room.wGroup)
	}
	return
}

// HandleRequest processes the equest came from the user
func (room *Room) HandleRequest(conn *Connection, rr *RoomRequest) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	if room == nil {
		return
	}

	utils.Debug(false, "room handle conn", room.Status, conn.Index(), conn.Disconnected())
	found, _ := room.Search(conn)
	if found != nil {
		utils.Debug(false, "room handle found", room.isAlive(found), found.Index(), found.Disconnected())
	}
	if rr.IsGet() {
		//go room.greet(conn)
	} else if rr.IsSend() {
		switch {
		case rr.Send.Messages != nil:
			Messages(conn, rr.Send.Messages, room.Messages())
		case rr.Send.Cell != nil:
			if room.isAlive(conn) {
				go room.CellHandle(conn, rr.Send.Cell)
			}
		case rr.Send.Action != nil:
			utils.Debug(false, "rr.Send.Action != nil ")
			room.chanConnection <- &ConnectionAction{
				conn:   conn,
				action: *rr.Send.Action,
			}
		}
	} else if rr.Message != nil {
		if conn.Index() < 0 {
			rr.Message.Status = models.StatusObserver
		} else {
			rr.Message.Status = models.StatusPlayer
		}
		Message(room.lobby, conn, rr.Message, room.appendMessage,
			room.setMessage, room.removeMessage, room.findMessage,
			room.send.sendAll, room.All, room, room.dbChatID)
	}
}

// StartFlagPlacing prepare field, players and observers
func (room *Room) StartFlagPlacing() {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	room.FillField()

	room.setStatus(StatusFlagPlacing)

	playersIterator := NewConnectionsIterator(room.Players.Connections)
	for playersIterator.Next() {
		player := playersIterator.Value()
		go room.send.greet(player, true)
		room.lobby.waiterToPlayer(player, room)
	}

	observersIterator := NewConnectionsIterator(room.Observers)
	for playersIterator.Next() {
		observer := observersIterator.Value()
		go room.send.greet(observer, true)
		room.lobby.waiterToPlayer(observer, room)
	}
	room.Players.Init(room.Field)

	room.wGroup.Add(1)
	room.lobby.RoomStart(room, room.wGroup)

	go room.send.StatusToAll(room.All, StatusFlagPlacing, nil)
	go room.send.Field(room.All)
}

// StartGame start game
func (room *Room) StartGame() {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()
	utils.Debug(false, "StartGame")

	room.FillField()

	open := float64(room.Settings.Mines) / float64(room.Settings.Width*room.Settings.Height) * float64(100)
	utils.Debug(false, "opennn", open, room.Settings.Width*room.Settings.Height)

	cells := room.Field.OpenZero() //room.Field.OpenSave(int(open))
	go room.send.NewCells(room.All, cells...)
	room.setStatus(StatusRunning)
	room.setDate(time.Now().In(room.lobby.location()))
	go room.send.StatusToAll(room.All, StatusRunning, nil)
	go room.send.Message("Battle began! Destroy your enemy!", room.All)
}

// FinishGame finish game
func (room *Room) FinishGame(timer bool) {
	if room == nil {
		utils.Debug(true, "room nil")
	}
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	room.setStatus(StatusFinished)

	// save Group
	saveAndSendGroup := &sync.WaitGroup{}

	cells := make([]Cell, 0)
	room.Field.OpenEverything(&cells)

	saveAndSendGroup.Add(1)
	go room.send.GameOver(timer, room.All, cells, saveAndSendGroup)

	saveAndSendGroup.Add(1)
	go room.models.Save(saveAndSendGroup)

	saveAndSendGroup.Add(1)
	go room.Players.m.Finish(saveAndSendGroup)
	saveAndSendGroup.Wait()

	go room.metricsRoom(room.lobby.config().Metrics, false)

	room.wGroup.Add(1)
	room.lobby.roomFinish(room, room.wGroup)
}

func (room *Room) metricsRoom(needMetrics bool, cancel bool) {
	if !needMetrics {
		return
	}
	var (
		roomType        string
		anonymous, mode int
	)
	if cancel {
		roomType = "aborted"
		metrics.AbortedRooms.Inc()
	} else {
		roomType = "finished"
		metrics.FinishedRooms.Inc()
	}
	if !room.Settings.NoAnonymous {
		anonymous = 1
	}
	if room.Settings.Deathmatch {
		mode = 1
	}

	size := float64(room.Settings.Width * room.Settings.Height)

	utils.Debug(false, "metrics RoomPlayers", room.Settings.Players)
	metrics.RoomPlayers.WithLabelValues(roomType).Observe(float64(room.Settings.Players))
	utils.Debug(false, "metrics difficult", room.Field.Difficult)
	metrics.RoomDifficult.WithLabelValues(roomType).Observe(float64(room.Field.Difficult))
	utils.Debug(false, "metrics size", size)
	metrics.RoomSize.WithLabelValues(roomType).Observe(size)
	utils.Debug(false, "metrics TimeToPlay", room.Settings.TimeToPlay)
	metrics.RoomTime.WithLabelValues(roomType).Observe(float64(room.Settings.TimeToPlay))
	if !cancel {
		openProcent := 1 - float64(float64(room.Field.cellsLeft())/size)
		utils.Debug(false, "metrics openProcent", openProcent)
		metrics.RoomOpenProcent.Observe(openProcent)

		utils.Debug(false, "metrics playing time", room.playingTime().Seconds())
		metrics.RoomTimePlaying.Observe(room.playingTime().Seconds())
	}
	metrics.RoomMode.WithLabelValues(roomType, utils.String(mode)).Inc()
	metrics.RoomAnonymous.WithLabelValues(roomType, utils.String(anonymous)).Inc()
	utils.Debug(false, "metrics recruitmentTime", room.recruitmentTime().Seconds())
	metrics.RoomTimeSearchingPeople.WithLabelValues(roomType).Observe(room.recruitmentTime().Seconds())
}

func (room *Room) CancelGame() {
	if room == nil {
		utils.Debug(true, "room nil")
	}
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	room.setStatus(StatusFinished)

	go room.metricsRoom(room.lobby.config().Metrics, true)

	room.wGroup.Add(1)
	room.lobby.roomFinish(room, room.wGroup)
}
