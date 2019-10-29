package game

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/metrics"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"time"
)

func (room *Room) runRoom() {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("room_handle.go run()")
		room.wGroup.Done()
	}()
	// все в конфиг
	ticker := time.NewTicker(time.Second * 10)
	//var timeoutPeopleFinding, timeoutPlayer, timeoutObserver, timeoutFinished float64
	// timeoutPeopleFinding = 2
	// timeoutPlayer = 60
	// timeoutObserver = 5
	// timeoutFinished = 20

	room.initTimers(true)
	defer func() {
		ticker.Stop()
		room.prepare.Stop()
		room.play.Stop()
	}()

	var beginGame, timeOut bool

	for {
		select {
		//go room.launchGarbageCollector(timeoutPeopleFinding, timeoutPlayer, timeoutObserver, timeoutFinished)
		case <-room.prepare.C:
			if beginGame {
				room.prepareOver()
			}
		case <-room.play.C:
			if beginGame {
				timeOut = true
				room.playingOver()
			}
		case conn := <-room.chanConnection:
			go room.processConnectionAction(conn)
		case newStatus := <-room.chanStatus:
			oldStatus := room.Status()

			if newStatus == oldStatus || newStatus > StatusFinished {
				continue
			}
			if oldStatus == StatusRecruitment {
				room.setRecruitmentTime()
			} else if newStatus == StatusFinished {
				room.setPlayingTime()
			}
			switch newStatus {
			case StatusFlagPlacing:
				beginGame = true
				room.StartFlagPlacing()
			case StatusRunning:
				room.StartGame()
			case StatusFinished:
				if oldStatus == StatusRecruitment {
					room.CancelGame()
				} else {
					room.FinishGame(timeOut)
				}
			//return
			case StatusAborted:
				ticker.Stop()
				return
			}
		}
	}
}

// initTimers launch game timers. Call it when flag placement starts
func (room *Room) initTimers(first bool) {
	if first {
		room.prepare = time.NewTimer(time.Millisecond)
		room.play = time.NewTimer(time.Millisecond)
	} else {
		if room.Settings.Deathmatch {
			room.prepare.Reset(time.Second *
				time.Duration(room.Settings.TimeToPrepare))
		} else {
			room.prepare.Reset(time.Millisecond)
		}
		room.play.Reset(time.Second *
			time.Duration(room.Settings.TimeToPlay))
	}
	return
}

func (room *Room) launchGarbageCollector(timeoutPeopleFinding, timeoutPlayer, timeoutObserver, timeoutFinished float64) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("room_handle.go run()")
		room.wGroup.Done()
	}()

	status := room.Status()
	if status == StatusRecruitment {
		timeoutPlayer = timeoutPeopleFinding
		timeoutObserver = timeoutPeopleFinding
	}
	if status == StatusFinished {
		timeoutPlayer = timeoutFinished
		timeoutObserver = timeoutFinished
	}
	i := 0

	playersIterator := NewConnectionsIterator(room.Players.Connections)
	for playersIterator.Next() {
		player := playersIterator.Value()
		if player == nil {
			utils.Debug(true, "found nil player")
		}

		i++
		t := player.Time()
		if player.Disconnected() && time.Since(t).Seconds() > timeoutPlayer {
			//fmt.Println(player.User.Name, " - bad")
			room.Kill(player, ActionTimeout)
			room.Leave(player, true)
		} else {
			//fmt.Println(player.User.Name, " - good", player.Disconnected(), time.Since(t).Seconds())
		}
	}

	observerIterator := NewConnectionsIterator(room.Observers)
	for observerIterator.Next() {
		observer := observerIterator.Value()
		if observer == nil {
			utils.Debug(true, "found nil observer")
		}

		i++
		t := observer.Time()
		if observer.Disconnected() && time.Since(t).Seconds() > timeoutObserver {
			room.Leave(observer, false)
		} else {
			//fmt.Println(conn.User.Name, " - good", conn)
		}
	}
}

func (room *Room) processActionBackToLobby(conn *Connection) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("room_handle.go run()")
		room.wGroup.Done()
	}()

	room.wGroup.Add(1)
	go room.lobby.LeaveRoom(conn, ActionBackToLobby, room, room.wGroup)

	room.Leave(conn, conn.Index() >= 0)
	if conn.Index() >= 0 {
		room.Kill(conn, ActionBackToLobby)
	}

	room.LeaveMeta(conn, ActionDisconnect)
}

func (room *Room) processActionDisconnect(conn *Connection) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("room_handle.go run()")
		room.wGroup.Done()
	}()

	if conn.PlayingRoom() == nil {
		room.processActionBackToLobby(conn)
		return
	}

	found, _ := room.Search(conn)
	if found == nil {
		return
	}

	found.setDisconnected()
	pa := *room.addAction(found.ID(), ActionDisconnect)
	room.sendAction(pa, room.All)

	// if conn.ID() < 0 /*conn.ID() < 0*/ /*|| time.Since(conn.time).Seconds() > timeout.Seconds()*/ {

	// 	pa := *room.addAction(conn.ID(), ActionDisconnect)
	// 	room.sendAction(pa, room.All)
	// 	found.setDisconnected()
	// }
}

func (room *Room) processActionReconnect(conn *Connection) {
	if room.lobby.config().Metrics {
		metrics.RoomsReconnections.Inc()
	}

	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("room_handle.go run()")
		room.wGroup.Done()
	}()

	found, isPlayer := room.Search(conn)
	if found == nil {
		return
	}
	room.addConnection(conn, isPlayer, true)
}

func (room *Room) processActionGiveUp(conn *Connection) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("room_handle.go run()")
		room.wGroup.Done()
	}()

	if room.IsActive() {
		go room.GiveUp(conn)
	}
}

func (room *Room) processActionRestart(conn *Connection) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("room_handle.go run()")
		room.wGroup.Done()
	}()

	status := room.Status()
	// if status == StatusRunning || status == StatusFlagPlacing {
	// 	fmt.Println("room.Status == StatusRunning || room.Status == StatusFlagPlacing")
	// 	return
	// }
	if status == StatusFinished {
		room.Restart(conn)
	}
}

func (room *Room) processConnectionAction(ca *ConnectionAction) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("room_handle.go run()")
		room.wGroup.Done()
	}()

	switch ca.action {
	case ActionBackToLobby:
		room.processActionBackToLobby(ca.conn)
	case ActionDisconnect:
		room.processActionDisconnect(ca.conn)
	case ActionReconnect:
		room.processActionReconnect(ca.conn)
	case ActionGiveUp:
		room.processActionGiveUp(ca.conn)
	case ActionRestart:
		room.processActionRestart(ca.conn)
	}
}

/*
func (room *Room) runHistory(conn *Connection) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("room_handle.go run()")
		room.wGroup.Done()
	}()

	//players := *room.Players
	actions := room.history()
	cells := room.Field.History
	actionsSize := len(actions)
	cellsSize := len(cells)
	actionsI := 0
	cellsI := 0
	actionTime := time.Now()
	cellTime := time.Now()
	if actionsSize > 0 {
		actionTime = actions[0].Time
	}
	if cellsSize > 0 {
		cellTime = cells[0].Time
	}

	// offline, err := room.lobby.createRoom(room.Settings)
	// if err != nil {
	// 	panic("offline doesnt work")
	// }
	// room.Leave(conn, ActionBackToLobby)
	// offline.Enter(conn)
	// offline.MakeObserver(conn, true)

	// offline.StartFlagPlacing()

	ticker := time.NewTicker(time.Millisecond * 10)
	defer func() {
		room.Status = StatusHistory
		ticker.Stop()
	}()

	for {
		select {
		case value := <-ticker.C:
			if (actionsSize + cellsSize) == (actionsI + cellsI) {
				return
			}
			for actionsI < actionsSize && value > time.Since() {

			}

		}
	}
}*/
