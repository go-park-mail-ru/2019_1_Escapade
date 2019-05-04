package game

import (
	"sync"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	"fmt"
	"math/rand"
	"time"
)

// Field send to user, if he disconnect and 'forgot' everything
// about map or it is his first connect
type Field struct {
	wGroup *sync.WaitGroup

	doneM *sync.RWMutex
	done  bool

	matrixM *sync.RWMutex
	Matrix  [][]int `json:"-"`

	historyM *sync.Mutex
	History  []Cell `json:"history"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`

	cellsLeftM *sync.RWMutex
	CellsLeft  int `json:"-"`
	Mines      int `json:"mines"`
}

// NewField create new instance of field
func NewField(rs *models.RoomSettings) *Field {
	matrix := generate(rs)
	field := &Field{
		wGroup: &sync.WaitGroup{},

		doneM: &sync.RWMutex{},
		done:  false,

		matrixM: &sync.RWMutex{},
		Matrix:  matrix,

		historyM: &sync.Mutex{},
		History:  make([]Cell, 0, rs.Width*rs.Height),
		Width:    rs.Width,
		Height:   rs.Height,
		Mines:    rs.Mines,

		cellsLeftM: &sync.RWMutex{},
		CellsLeft:  rs.Width*rs.Height - rs.Mines - rs.Players,
	}
	return field
}

// place wg.add() and wg.Done() to all import functions

// Free clear matrix and history
func (field *Field) Free() {

	field.setDone()

	field.wGroup.Wait()
	for i := 0; i < len(field.Matrix); i++ {
		field.Matrix[i] = nil
	}
	field.Matrix = nil
	field.History = nil
}

// SameAs compare two fields
func (field *Field) SameAs(another *Field) bool {
	field.wGroup.Add(1)
	defer field.wGroup.Done()

	left1 := field.getCellsLeft()
	left2 := field.getCellsLeft()

	compare := field.Width == another.Width &&
		field.Height == another.Height &&
		left1 == left2

	return compare
}

// OpenEverything open all cells
func (field *Field) OpenEverything(cells *[]Cell) {
	if field.getDone() {
		return
	}
	field.wGroup.Add(1)
	defer field.wGroup.Done()

	for i := 0; i < field.Height; i++ {
		for j := 0; j < field.Width; j++ {
			v := field.getMatrixValue(i, j)
			if v != CellOpened {
				cell := NewCell(i, j, v, 0)
				field.saveCell(cell, cells)
				field.setCellOpen(i, j)
			}
		}
	}
}

// openCellArea open cell area, if there is no mines around
// in this cell
func (field *Field) openCellArea(x, y, ID int, cells *[]Cell) {
	if field.areCoordinatesRight(x, y) {
		v := field.getMatrixValue(x, y)

		if v < CellMine {
			cell := NewCell(x, y, v, ID)
			field.saveCell(cell, cells)
			field.decrementCellsLeft()
			field.setCellOpen(x, y)
		}
		if v == 0 {
			field.openCellArea(x-1, y-1, ID, cells)
			field.openCellArea(x-1, y, ID, cells)
			field.openCellArea(x-1, y+1, ID, cells)

			field.openCellArea(x, y+1, ID, cells)
			field.openCellArea(x, y-1, ID, cells)

			field.openCellArea(x+1, y-1, ID, cells)
			field.openCellArea(x+1, y, ID, cells)
			field.openCellArea(x+1, y+1, ID, cells)
		}
	}
}

// IsCleared return true if all safe cells except flags open
func (field *Field) IsCleared() bool {
	if field.getDone() {
		return true
	}
	field.wGroup.Add(1)
	defer field.wGroup.Done()

	return field.getCellsLeft() == 0
}

func (field *Field) saveCell(cell *Cell, cells *[]Cell) {
	if cell.Value != CellOpened && cell.Value != CellFlagTaken {
		cell.Time = time.Now()
		field.setToHistory(*cell)
	}
	*cells = append(*cells, *cell)
}

// OpenCell open cell and return slice of opened cells
func (field *Field) OpenCell(cell *Cell) (cells []Cell) {
	if field.getDone() {
		return
	}
	field.wGroup.Add(1)
	defer field.wGroup.Done()

	cell.Value = field.getMatrixValue(cell.X, cell.Y)

	cells = make([]Cell, 0)
	if cell.Value < CellMine {
		field.openCellArea(cell.X, cell.Y, cell.PlayerID, &cells)
	} else {
		if cell.Value != cell.PlayerID+CellIncrement {
			field.saveCell(cell, &cells)
		}
	}

	return
}

// RandomFlags create random players flags
func (field *Field) RandomFlags(players []Player) (cells []Cell) {
	if field.getDone() {
		return
	}
	field.wGroup.Add(1)
	defer field.wGroup.Done()

	cells = make([]Cell, len(players))
	for i, player := range players {
		cells[i] = field.CreateRandomFlag(player.ID)
	}
	return cells
}

// CreateRandomFlag create flag for player
func (field *Field) CreateRandomFlag(playerID int) (cell Cell) {
	if field.getDone() {
		return
	}
	field.wGroup.Add(1)
	defer field.wGroup.Done()

	rand.Seed(time.Now().UnixNano())
	x := rand.Intn(field.Width)
	y := rand.Intn(field.Height)
	cell = *NewCell(x, y, playerID+CellIncrement, playerID)

	return cell
}

// SetMines fill matrix with mines
func (field *Field) SetMines() {
	if field.getDone() {
		return
	}
	field.wGroup.Add(1)
	defer field.wGroup.Done()

	width := field.Width
	height := field.Height
	mines := field.Mines
	fmt.Println("begin SetMines")
	for mines > 0 {
		rand.Seed(time.Now().UnixNano())
		i := rand.Intn(width)
		j := rand.Intn(height)
		fmt.Println(i, j)
		if field.lessThenMine(i, j) {
			field.setMine(i, j)
			mines--
		}
	}
	fmt.Println("end SetMines")
}

// generate matrix
func generate(rs *models.RoomSettings) (matrix [][]int) {
	width := rs.Width
	height := rs.Height

	matrix = [][]int{}
	for i := 0; i < height; i++ {
		matrix = append(matrix, make([]int, width))
	}
	return
}

// IsInside check if coordinates are in field
func (field Field) areCoordinatesRight(x, y int) bool {
	return x >= 0 && x < field.Width && y >= 0 && y < field.Height
}

// IsInside check is cell inside fueld
func (field Field) IsInside(cell *Cell) bool {
	return field.areCoordinatesRight(cell.X, cell.Y)
}

///////////////////// Set cells func //////////

// SetFlag works only when mines not set
func (field *Field) SetFlag(cell *Cell) {
	if field.getDone() {
		return
	}
	field.wGroup.Add(1)
	defer field.wGroup.Done()

	// To identifier which flag we see, lets set id
	// add CellIncrement to id, because if id = 3 we can think that there are 3 mines around
	// we cant use -id, becase in future there will be a lot of conditions with
	// something < 9 (to find not mine places)
	fmt.Println("setFlag", cell.X, cell.Y, cell.PlayerID, CellIncrement)

	field.setMatrixValue(cell.X, cell.Y, cell.PlayerID+CellIncrement)
}

// SetCellFlagTaken set cells flag taken
func (field *Field) SetCellFlagTaken(cell *Cell) {
	if field.getDone() {
		return
	}
	field.wGroup.Add(1)
	defer field.wGroup.Done()

	field.setMatrixValue(cell.X, cell.Y, CellFlagTaken)

	fmt.Println("flag found!", cell.Value)
}

// setCellOpen set cell opened
func (field *Field) setCellOpen(x, y int) {
	field.setMatrixValue(x, y, CellOpened)
}

// setMine add mine to matrix and increase dangerous value in cells near mine
func (field *Field) setMine(x, y int) {

	width := field.Width
	height := field.Height
	field.setMatrixValue(x, y, CellMine)
	for i := x - 1; i <= x+1; i++ {
		if i >= 0 && i < width {
			for j := y - 1; j <= y+1; j++ {
				if j >= 0 && j < height {
					if field.lessThenMine(i, j) {
						field.incrementMatrixValue(i, j)
					}
				}
			}
		}
	}
}