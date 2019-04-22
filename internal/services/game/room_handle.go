package game

import (
	"encoding/json"
	"fmt"
	"time"
)

// Game status
const (
	StatusPeopleFinding = iota
	StatusAborted       // in case of error
	StatusFlagPlacing
	StatusRunning
	StatusFinished
	StatusClosed
)

// Room consist of players and observers, field and history
type Room struct {
	Name   string `json:"name"`
	Status int    `json:"status"`

	Players   *OnlinePlayers `json:"players,omitempty"`
	Observers *Connections   `json:"observers,omitempty"`

	History []*PlayerAction `json:"history,omitempty"`

	lobby *Lobby
	Field *Field `json:"field,omitempty"`

	killed int //amount of killed users
}

// SameAs compare  one room with another
func (room *Room) SameAs(another *Room) bool {
	return room.Field.SameAs(another.Field)
}

// Enter handle user joining as player or observer
func (room *Room) Enter(conn *Connection) bool {

	// if room is searching new players
	if room.Status == StatusPeopleFinding {
		conn.debug("You will be player!")
		if room.addPlayer(conn) {
			return true
		}
	} else if room.addObserver(conn) {
		conn.debug("You will be observer!")
		return true
	}
	conn.debug("No way!")

	return false
}

// Leave handle user going back to lobby
func (room *Room) Leave(conn *Connection) {

	// cant delete players, cause they always need
	// if game began
	conn.PushToLobby()
	room.lobby.playerToWaiter(conn)
	if !room.IsActive() {
		room.removeBeforeLaunch(conn)
	} else {
		room.removeDuringGame(conn)
	}
	conn.debug("Welcome back to lobby!")
}

func (room *Room) setFlag(conn *Connection, cell *Cell) bool {
	// if user try set flag after game launch
	if room.Status != StatusFlagPlacing {
		return false
	}

	if !room.Field.IsInside(cell) {
		return false
	}
	i := room.Players.Search(conn)
	if i < 0 {
		return false
	}

	room.Players.Flags[i].X = cell.X
	room.Players.Flags[i].Y = cell.Y
	return true
}

// openCell open cell
func (room *Room) openCell(conn *Connection, cell *Cell) bool {
	// if user try set open cell before game launch
	if room.Status != StatusRunning {
		return false
	}

	// if wrong cell
	if !room.Field.IsInside(cell) {
		return false
	}

	// if user died
	if conn.Player.Finished == true {
		return false
	}

	// set who try open cell(for history)
	cell.PlayerID = conn.GetPlayerID()
	room.Field.OpenCell(cell)

	room.sendField(room.all())

	if room.Field.IsCleared() {
		room.lobby.roomFinish(room)
	}
	return true
}

func (room *Room) cellHandle(conn *Connection, cell *Cell) (done bool) {
	fmt.Println("cellHandle")
	if room.Status == StatusFlagPlacing {
		done = room.setFlag(conn, cell)
	} else if room.Status == StatusRunning {
		done = room.openCell(conn, cell)
	}
	return
}

// IsActive check if game is started and results not known
func (room *Room) IsActive() bool {
	return room.Status == StatusFlagPlacing || room.Status == StatusRunning
}

func (room *Room) actionHandle(conn *Connection, action int) (done bool) {
	if room.IsActive() {
		if action == ActionGiveUp {
			conn.debug("we see you wanna give up?")
			room.GiveUp(conn)
			return true
		}
	}
	if action == ActionBackToLobby {
		conn.debug("we see you wanna back to lobby?")
		room.Leave(conn) // exit to lobby
		return true
	}
	return false
}

// handleRequest
func (room *Room) handleRequest(conn *Connection, rr *RoomRequest) {

	conn.debug("room handle conn")
	if rr.IsGet() {
		room.requestGet(conn, rr)
	} else if rr.IsSend() {
		done := false
		if rr.Send.Cell != nil {
			if conn.IsPlayerAlive() {
				done = room.cellHandle(conn, rr.Send.Cell)
			}
		} else if rr.Send.Action != nil {

			done = room.actionHandle(conn, *rr.Send.Action)
		}
		if !done {
			//conn.debug("Room cant execute request")
			Answer(conn, []byte("Room cant execute request "))
		}
	}
}

func (room *Room) startFlagPlacing() {
	room.Status = StatusFlagPlacing
	room.lobby.roomStart(room)
	room.Players.Init(room.Field)
	go room.run()
	room.sendField(room.all())
	room.sendMessage("Battle will be start soon! Set your flag!", room.all())
}

func (room *Room) startGame() {
	room.Status = StatusRunning
	room.fillField()
	room.sendMessage("Battle began! Destroy your enemy!", room.all())
}

func (room *Room) finishGame() {
	room.Status = StatusFinished
	go room.lobby.roomFinish(room)
	for _, player := range room.Players.Players {
		player.Finished = true
	}
	room.sendMessage("Battle finished!", room.all())
}

func (room *Room) run() {
	// перенести в настройки комнаты
	timerPrepare := time.NewTimer(time.Second * 20)
	timerPlaying := time.NewTimer(time.Second * 180)
	// в конфиг
	ticker := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-timerPrepare.C:
			room.startGame()
		case <-timerPlaying.C:
			room.finishGame()
			return
		case clock := <-ticker.C:
			room.sendMessage(clock.String()+" passed", room.all())
		}
	}
}

func (room *Room) requestGet(conn *Connection, rr *RoomRequest) {
	send := room.copy(rr.Get)
	bytes, _ := json.Marshal(send)
	conn.SendInformation(bytes)
}
