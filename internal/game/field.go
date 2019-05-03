package game

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	"fmt"
	"math/rand"
	"time"
)

// Field send to user, if he disconnect and 'forgot' everything
// about map or it is his first connect
type Field struct {
	Matrix    [][]int `json:"-"`
	History   []Cell  `json:"history"`
	Width     int     `json:"width"`
	Height    int     `json:"height"`
	CellsLeft int     `json:"-"`
	Mines     int
}

// Clear clear matrix and history
func (field *Field) Clear() {
	for i := 0; i < len(field.Matrix); i++ {
		field.Matrix[i] = nil
	}
	field.Matrix = nil
	field.History = nil
}

// SameAs compare two fields
func (field *Field) SameAs(another *Field) bool {
	return field.Width == another.Width &&
		field.Height == another.Height &&
		field.CellsLeft == another.CellsLeft
}

// SetFlag works only when mines not set
func (field *Field) SetFlag(x int, y int, id int) {

	//field.Matrix[x][y] = CellFlag

	// To identifier which flag we see, lets set id
	// add CellIncrement to id, because if id = 3 we can think that there are 3 mines around
	// we cant use -id, becase in future there will be a lot of conditions with
	// something < 9 (to find not mine places)
	fmt.Println("setFlag", x, y, id, CellIncrement)
	field.Matrix[x][y] = id + CellIncrement
}

func (field *Field) openEverything(cells *[]Cell) {
	for i := 0; i < field.Height; i++ {
		for j := 0; j < field.Width; j++ {
			v := field.Matrix[i][j]
			if v != CellOpened {
				cell := NewCell(i, j, v, 0)
				field.saveCell(cell, cells)
				field.setCellOpen(i, j)
			}
		}
	}
}

func (field *Field) openCellArea(x, y, ID int, cells *[]Cell) {
	if field.areCoordinatesRight(x, y) {
		v := field.Matrix[x][y]

		if v < CellMine {
			cell := NewCell(x, y, v, ID)
			field.saveCell(cell, cells)
			field.CellsLeft--
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

func (field *Field) setCellOpen(x, y int) {
	field.Matrix[x][y] = CellOpened
}

// IsCleared return true if all safe cells except flags open
func (field *Field) IsCleared() bool {
	return field.CellsLeft == 0
}

func (field *Field) setCellFlagTaken(cell *Cell) {
	field.Matrix[cell.X][cell.Y] = CellFlagTaken
	fmt.Println("flag found!", cell.Value)
}

func (field *Field) saveCell(cell *Cell, cells *[]Cell) {
	if cell.Value != CellOpened && cell.Value != CellFlagTaken {
		cell.Time = time.Now()
		field.History = append(field.History, *cell)
	}
	*cells = append(*cells, *cell)
}

// OpenCell open cell and return slice of opened cells
func (field *Field) OpenCell(cell *Cell) (cells []Cell) {
	cell.Value = field.Matrix[cell.X][cell.Y]

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

// setMine add mine to matrix and increase dangerous value in cells near mine
func (field *Field) setMine(x, y int) {

	width := field.Width
	height := field.Height
	field.Matrix[x][y] = CellMine
	for i := x - 1; i <= x+1; i++ {
		if i >= 0 && i < width {
			for j := y - 1; j <= y+1; j++ {
				if j >= 0 && j < height && field.Matrix[i][j] < CellMine {
					field.Matrix[i][j]++
				}
			}
		}
	}
}

// RandomFlags create random players flags
func (field *Field) RandomFlags(players []Player) (cells []Cell) {
	cells = make([]Cell, len(players))
	for i, player := range players {
		cells[i] = field.CreateRandomFlag(player.ID)
	}
	return cells
}

// CreateRandomFlag create flag for player
func (field *Field) CreateRandomFlag(playerID int) (cell Cell) {
	rand.Seed(time.Now().UnixNano())
	x := rand.Intn(field.Width)
	y := rand.Intn(field.Height)
	cell = *NewCell(x, y, playerID+CellIncrement, playerID)

	return cell
}

// SetMines fill matrix with mines
func (field *Field) SetMines() {
	width := field.Width
	height := field.Height
	mines := field.Mines
	fmt.Println("begin SetMines")
	for mines > 0 {
		rand.Seed(time.Now().UnixNano())
		i := rand.Intn(width)
		j := rand.Intn(height)
		fmt.Println(i, j)
		if field.Matrix[i][j] < CellMine {
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

// NewField create new instance of field
func NewField(rs *models.RoomSettings) *Field {
	matrix := generate(rs)
	field := &Field{
		Matrix:    matrix,
		History:   make([]Cell, 0, rs.Width*rs.Height),
		Width:     rs.Width,
		Height:    rs.Height,
		Mines:     rs.Mines,
		CellsLeft: rs.Width*rs.Height - rs.Mines - rs.Players,
	}
	return field
}

// RandomCell create cell with random X,Y inside field
func (field Field) RandomCell() *Cell {
	cell := &Cell{
		X: rand.Intn(field.Width),
		Y: rand.Intn(field.Height),
	}
	return cell
}

// IsInside check if coordinates are in field
func (field Field) areCoordinatesRight(x, y int) bool {
	return x >= 0 && x < field.Width && y >= 0 && y < field.Height
}

// IsInside check is cell inside fueld
func (field Field) IsInside(cell *Cell) bool {
	return field.areCoordinatesRight(cell.X, cell.Y)
}
