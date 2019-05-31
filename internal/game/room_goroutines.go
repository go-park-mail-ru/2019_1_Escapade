package game

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"fmt"
	"time"
)

func (room *Room) runRoom() {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		fmt.Println("Sorry it is end")
		utils.CatchPanic("room_handle.go run()")
		room.wGroup.Done()
	}()

	// все в конфиг
	ticker := time.NewTicker(time.Second * 1)
	var timeoutPeopleFinding, timeoutPlayer, timeoutObserver, timeoutFinished float64
	timeoutPeopleFinding = 2
	timeoutPlayer = 60
	timeoutObserver = 5
	timeoutFinished = 20

	for {
		select {
		case <-ticker.C:
			go room.launchGarbageCollector(timeoutPeopleFinding, timeoutPlayer, timeoutObserver, timeoutFinished)
		case conn := <-room.chanConnection:
			fmt.Println("we process it")
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

	room.Date = time.Now()

	for {
		select {
		case <-room.chanFinish:
			return
		case <-room.prepare.C:
			room.chanStatus <- StatusRunning
		case <-room.play.C:
			room.chanStatus <- StatusFinished
			//fmt.Println("finish!")
			//room.FinishGame(true)
			return
		case clock := <-ticker.C:
			//fmt.Println("clock!", room.ID)
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
	//fmt.Println("launchGarbageCollector")

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
	playerGone := room.LeavePlayer(conn)
	observerGone := room.LeaveObserver(conn)

	fmt.Println("look", playerGone, observerGone)
	if playerGone {
		fmt.Println("go away")
		room.lobby.LeaveRoom(conn, ActionBackToLobby)
		room.LeaveMeta(conn, ActionDisconnect)
	}
	if observerGone {
		fmt.Println("go away")
		room.lobby.LeaveRoom(conn, ActionBackToLobby)
		room.LeaveMeta(conn, ActionDisconnectObserver)
	}
}

func (room *Room) processActionDisconnect(conn *Connection) {
	found, _ := room.Search(conn)
	var refreshSeconds = 1
	fmt.Println("tiiiiiime ", conn.time, time.Since(conn.time).Seconds(), float64(refreshSeconds))
	if conn.ID() < 0 || time.Since(conn.time).Seconds() > float64(refreshSeconds) {

		pa := *room.addAction(conn.ID(), ActionTaken)
		room.sendAction(pa, room.AllExceptThat(found))
		found.setDisconnected()
	}
}

func (room *Room) processActionConnect(conn *Connection) {
	found, isPlayer := room.Search(conn)
	if found == nil {
		return
	}
	room.sendAccountTaken(*found)
	conn.time = time.Now()
	found = conn
	if isPlayer {
		room.RecoverPlayer(conn)
	} else {
		room.RecoverObserver(conn)
	}
}

func (room *Room) processActionGiveUp(conn *Connection) {
	if room.IsActive() {
		go room.GiveUp(conn)
	}
}

func (room *Room) processActionRestart(conn *Connection) {
	if room.Status == StatusRunning || room.Status == StatusFlagPlacing {
		return
	}
	conn.lobby.greet(conn)
	if room.Status == StatusFinished {
		fmt.Println("goood")
		room.Restart()
		room.lobby.addRoom(room)
	}
	if room.Status == StatusPeopleFinding {
		room.addPlayer(conn, false)
	}
}

func (room *Room) processConnectionAction(ca ConnectionAction) {
	fmt.Println("processConnectionAction start ")

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
		fmt.Println("processActionConnect ------")
		room.processActionConnect(ca.conn)
		fmt.Println("processActionConnect ")
	case ActionGiveUp:
		fmt.Println("processActionGiveUp ------")
		room.processActionGiveUp(ca.conn)
		fmt.Println("processActionGiveUp")
	case ActionRestart:
		fmt.Println("processActionRestart ------")
		room.processActionRestart(ca.conn)
		fmt.Println("processActionRestart ")
	}

	fmt.Println("processConnectionAction finish ")
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
