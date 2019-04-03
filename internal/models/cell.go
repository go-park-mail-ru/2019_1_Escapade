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
