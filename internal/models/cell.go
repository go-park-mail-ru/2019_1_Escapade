package models

// Cell send to user one cell
type Cell struct {
	X int
	Y int
	// IsOpen bool
	// IsMarked bool
	Value    int
	PlayerID int
}

func NewCell(x int, y int, v int) *Cell {
	cell := &Cell{
		X:     x,
		Y:     y,
		Value: v,
	}
	return cell
}

// SameAs compare two cells
// func (cell *Cell) SameAs(another *Cell) bool {
// 	return (cell.X == another.X) && (cell.Y == another.Y)
// }

/*
0-8- amount of mines
9 - mine
10 - flag
11 - almost open
*/

// RoomJoining joining to room
// type RoomJoining struct {
// 	Success bool
// 	Message string
// }

// func RoomJoiningFault(reason string) *RoomJoining {
// 	rj := &RoomJoining{
// 		Success: false,
// 		Message: reason,
// 	}
// 	return rj
// }

// func RoomJoiningSuccess() *RoomJoining {
// 	rj := &RoomJoining{
// 		Success: true,
// 	}
// 	return rj
// }

type Status struct {
	Ready       bool
	PeopleFound int
}

func NewStatus(peopleFound int, ready bool) *Status {
	status := &Status{
		Ready:       ready,
		PeopleFound: peopleFound,
	}
	return status
}

// Players send to user, if he disconnect and 'forgot' everything
// about users or it is his first connect
type People struct {
	PlayersCapacity int
	PlayersSize     int
	Players         []Player

	ObserversCapacity int
	ObserversSize     int
	Observers         []Player
}

type Player struct {
	ID         int
	Name       string
	Points     int
	LastAction int
}

func NewPlayer(name string, id int) *Player {
	player := &Player{
		ID:         id,
		Name:       name,
		Points:     0,
		LastAction: ActionNo,
	}
	return player
}

// Reset - call it in every game beginning
func (player *Player) Reset() {
	player.LastAction = ActionNo
	player.Points = 0
}

// Cell type
const (
	CellMine   = iota + 9
	CellOpened // for empty cells
	CellFlag
	CellFlagTaken
	CellIncrement // for id
)

// Player actions
const (
	ActionError = iota - 1
	ActionNo
	ActionConnectAsPlayer
	ActionConnectAsObserver
	ActionReconnect
	ActionDisconnect
	ActionStop
	ActionContinue
	ActionExplode
	ActionWin
	ActionLose
	ActionGetPoints
	ActionFlagSet
	ActionGiveUp
)

// What to send to user
const (
	SendPlayerAction = iota
	SendGameStatus   // in case of error
	SendRoomSettings
	SendCells
)

type ClientData struct {
	Send         int
	RoomSettings *RoomSettings
	Cell         *Cell
	PlayerAction int
}

type GameInfo struct {
	// show, what of fields is filled
	Send int
	// any player action.
	PlayerAction Player
	// game status
	Status int
	//cell information. Send during game for all players
	Cells []Cell
}

type Flag struct {
	ID   int
	Cell *Cell
}

func NewFlag(cell *Cell, id int) *Flag {
	flag := &Flag{
		ID:   id,
		Cell: cell,
	}
	return flag
}
