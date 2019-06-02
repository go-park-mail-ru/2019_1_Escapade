package game

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"fmt"
	"time"
)

func (room *Room) runRoom() {

	// все в конфиг
	ticker := time.NewTicker(time.Second * 10)
	var timeoutPeopleFinding, timeoutPlayer, timeoutObserver, timeoutFinished float64
	timeoutPeopleFinding = 2
	timeoutPlayer = 60
	timeoutObserver = 5
	timeoutFinished = 20

	room.initTimers()

	for {
		select {
		case <-ticker.C:
			go room.launchGarbageCollector(timeoutPeopleFinding, timeoutPlayer, timeoutObserver, timeoutFinished)
		case conn := <-room.chanConnection:
			go room.processConnectionAction(conn)
		case newStatus := <-room.chanStatus:
			if newStatus == room.Status || newStatus > StatusFinished {
				continue
			}
			switch newStatus {
			// case StatusPeopleFinding:
			// 	//
			case StatusFlagPlacing:
				room.StartFlagPlacing()
			case StatusRunning:
				room.prepare.Stop()
				room.StartGame()
			case StatusFinished:
				ok := room.play.Stop()
				ticker.Stop()
				room.FinishGame(!ok)
				//return
			case StatusAborted:
				ticker.Stop()
				return
			}
		}
	}
}

// run - room goroutine
func (room *Room) runGame() {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("room_handle.go run()")
		room.wGroup.Done()
	}()

	ticker := time.NewTicker(time.Second * 4)

	room.initTimers()
	defer func() {
		ticker.Stop()
		room.prepare.Stop()
		room.play.Stop()
		fmt.Println("Room: Game is over!")
	}()

	loc, _ := time.LoadLocation("Europe/Moscow")
	room.Date = time.Now().In(loc)

	for {
		select {
		case <-room.chanFinish:
			return
		case <-room.prepare.C:
			room.chanStatus <- StatusRunning
		case <-room.play.C:
			room.chanStatus <- StatusFinished
			return
		case clock := <-ticker.C:
			go room.sendMessage(clock.String()+" passed", room.All)
		}
	}
}

// initTimers launch game timers. Call it when flag placement starts
func (room *Room) initTimers() {
	room.prepare = time.NewTimer(time.Second *
		time.Duration(room.Settings.TimeToPrepare))
	room.play = time.NewTimer(time.Second *
		time.Duration(room.Settings.TimeToPlay))
	return
}

func (room *Room) launchGarbageCollector(timeoutPeopleFinding, timeoutPlayer, timeoutObserver, timeoutFinished float64) {
	fmt.Println("launchGarbageCollector")
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("room_handle.go run()")
		room.wGroup.Done()
	}()

	if room.Status == StatusPeopleFinding {
		timeoutPlayer = timeoutPeopleFinding
		timeoutObserver = timeoutPeopleFinding
	}
	if room.Status == StatusFinished {
		timeoutPlayer = timeoutFinished
		timeoutObserver = timeoutFinished
	}
	i := 0
	for _, conn := range room.Players.Connections.RGet() {
		if conn == nil {
			continue
		}
		i++
		if conn.Disconnected() && time.Since(conn.time).Seconds() > timeoutPlayer {
			fmt.Println(conn.User.Name, " - bad")
			room.LeavePlayer(conn)
			//room.Leave(conn, ActionTimeOver)
		} else {
			fmt.Println(conn.User.Name, " - good", conn.Disconnected(), time.Since(conn.time).Seconds())
		}
	}
	for _, conn := range room.Observers.RGet() {
		if conn == nil {
			continue
		}
		i++
		if conn.Disconnected() && time.Since(conn.time).Seconds() > timeoutObserver {
			//fmt.Println(conn.User.Name, " - bad")
			room.LeaveObserver(conn)
			//room.Leave(conn, ActionTimeOver)
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

	playerGone := room.LeavePlayer(conn)
	observerGone := room.LeaveObserver(conn)

	fmt.Println("LeaveRoom")
	room.lobby.LeaveRoom(conn, ActionBackToLobby)
	fmt.Println("LeaveMeta")
	room.LeaveMeta(conn, ActionDisconnect)
	fmt.Println("went back", playerGone, observerGone)
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

	if conn == nil {
		return
	}

	found, _ := room.Search(conn)
	if found == nil {
		room.LeaveObserver(found)
		room.LeavePlayer(found)
	}
	fmt.Println("Disconnected")
	found.setDisconnected()
	pa := *room.addAction(found.ID(), ActionDisconnect)
	room.sendAction(pa, room.All)

	// if conn.ID() < 0 /*conn.ID() < 0*/ /*|| time.Since(conn.time).Seconds() > timeout.Seconds()*/ {

	// 	pa := *room.addAction(conn.ID(), ActionDisconnect)
	// 	room.sendAction(pa, room.All)
	// 	found.setDisconnected()
	// }
}

func (room *Room) processActionConnect(conn *Connection) {
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
	sendAccountTaken(*found)
	conn.SetConnected()
	//conn.time = time.Now()
	found = conn
	if isPlayer {
		room.RecoverPlayer(conn)
	} else {
		room.RecoverObserver(conn)
	}
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

	if room.Status == StatusRunning || room.Status == StatusFlagPlacing {
		fmt.Println("room.Status == StatusRunning || room.Status == StatusFlagPlacing")
		return
	}
	if room.Status == StatusFinished {
		pa := *room.addAction(conn.ID(), ActionRestart)
		room.sendAction(pa, room.All)
		room.Restart()
		room.lobby.addRoom(room)
	}
	fmt.Println("conn.lobby.greet(conn)")
	conn.lobby.greet(conn)
	if room.Status == StatusPeopleFinding {
		room.LeaveObserver(conn)
		room.addPlayer(conn, false)
	}
}

func (room *Room) processConnectionAction(ca ConnectionAction) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("room_handle.go run()")
		room.wGroup.Done()
	}()

	//fmt.Println("processConnectionAction start ")

	switch ca.action {
	case ActionBackToLobby:
		fmt.Println("processActionBackToLobby ------")
		room.processActionBackToLobby(ca.conn)
		fmt.Println("processActionBackToLobby ")
	case ActionDisconnect:
		fmt.Println("processActionDisconnect ------")
		room.processActionDisconnect(ca.conn)
		fmt.Println("processActionDisconnect")
	case ActionConnect:
		//fmt.Println("processActionConnect ------")
		room.processActionConnect(ca.conn)
		//fmt.Println("processActionConnect ")
	case ActionGiveUp:
		//fmt.Println("processActionGiveUp ------")
		room.processActionGiveUp(ca.conn)
		//fmt.Println("processActionGiveUp")
	case ActionRestart:
		//fmt.Println("processActionRestart ------")
		room.processActionRestart(ca.conn)
		//fmt.Println("processActionRestart ")
	}

	//fmt.Println("processConnectionAction finish ")
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
