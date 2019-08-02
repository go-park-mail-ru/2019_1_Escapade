package game

import "time"

// Cell send to user one cell
type Cell struct {
	X        int32     `json:"x"`
	Y        int32     `json:"y"`
	Value    int32     `json:"value"`
	PlayerID int32     `json:"playerID"`
	Time     time.Time `json:"-"`
}

// NewCell create new instance of cell
func NewCell(x int32, y int32, v int32, ID int32) *Cell {
	cell := &Cell{
		X:        x,
		Y:        y,
		Value:    v,
		PlayerID: ID,
		Time:     time.Now(),
	}
	return cell
}

// Cell type
const (
	CellMine   = iota + 9 // +9, cause <9 - amount of mines around
	CellOpened            // for empty cells
	CellFlag
	CellFlagTaken
	CellIncrement // for id
)
