package game

// Cell send to user one cell
type Cell struct {
	X        int `json:"x"`
	Y        int `json:"y"`
	Value    int `json:"value"`
	PlayerID int `json:"playerID"`
}

// NewCell create new instance of cell
func NewCell(x int, y int, v int, ID int) *Cell {
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
	CellMine   = iota + 9 // +9, cause <9 - amount of mines around
	CellOpened            // for empty cells
	CellFlag
	CellFlagTaken
	CellIncrement // for id
)
