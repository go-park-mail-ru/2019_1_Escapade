package models

// Cell send to user one cell
type Cell struct {
	X int `json:"X"`
	Y int `json:"Y"`
	// IsOpen bool
	// IsMarked bool
	Value    int `json:"Value"`
	PlayerID int `json:"PlayerID"`
}

func NewCell(x int, y int, v int) *Cell {
	cell := &Cell{
		X:     x,
		Y:     y,
		Value: v,
	}
	return cell
}

func NewCellWithID(x int, y int, v int, ID int) *Cell {
	cell := &Cell{
		X:        x,
		Y:        y,
		Value:    v,
		PlayerID: ID,
	}
	return cell
}

// Cell type
const (
	CellMine   = iota + 9
	CellOpened // for empty cells
	CellFlag
	CellFlagTaken
	CellIncrement // for id
)
<<<<<<< HEAD
=======

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
	SendAlive
)

type ClientData struct {
	Send         int           `json:"Send"`
	RoomSettings *RoomSettings `json:"RoomSettings"`
	Cell         *Cell         `json:"Cell"`
	PlayerAction int           `json:"PlayerAction"`
}

type ClientData2 struct {
	Send         int `json:"Send"`
	PlayerAction int `json:"PlayerAction"`
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
>>>>>>> 77581ff4f08b4dd3a48f7bb5d86af5aace9a5a78
