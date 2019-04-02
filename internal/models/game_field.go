package models

import (
	"math/rand"
)

// Field send to user, if he disconnect and 'forgot' everything
// about map or it is his first connect
type Field struct {
	Matrix    [][]int `json:"-"`
	History   []Cell
	Width     int
	Height    int
	CellsLeft int `json:"-"`
	// Open   int
	// Mines int
	// Flags int
}

// SetFlag works only when mines not set
func (field *Field) SetFlag(x int, y int, id int) {

	field.Matrix[x][y] = CellFlag

	// To identifier which flag we see, lets set id
	// add CellIncrement to id, because if id = 3 we can think that there are 3 mines around
	// we cant use -id, becase in future there will be a lot of conditions with
	// something < 9 (to find not mine places)
	field.Matrix[x][y] = id + CellIncrement
}

func (field *Field) openCellArea(x, y, ID int) {
	if field.areCoordinatesRight(x, y) {
		v := field.Matrix[x][y]
		if v < CellMine {
			cell := NewCell(x, y, v)
			cell.PlayerID = ID
			field.History = append(field.History, *cell)
			field.Matrix[x][y] = CellOpened
		}
		if v == 0 {
			field.Matrix[x][y] = CellOpened
			field.openCellArea(x-1, y-1, ID)
			field.openCellArea(x-1, y, ID)
			field.openCellArea(x-1, y+1, ID)

			field.openCellArea(x, y+1, ID)
			field.openCellArea(x, y-1, ID)

			field.openCellArea(x+1, y-1, ID)
			field.openCellArea(x+1, y, ID)
			field.openCellArea(x+1, y+1, ID)
		}
	}
}

func (field *Field) OpenCell(cell *Cell) {
	cell.Value = field.Matrix[cell.X][cell.Y]

	if cell.Value == 0 {
		field.openCellArea(cell.X, cell.Y, cell.PlayerID)
	} else {
		field.History = append(field.History, *cell)
		if cell.Value < CellMine {
			field.Matrix[cell.X][cell.Y] = CellOpened
		} else if cell.Value == CellFlag {
			field.Matrix[cell.X][cell.Y] = CellFlagTaken
		}
	}
}

// setMine add mine to matrix and increase dangerous value in cells near mine
func (field *Field) setMine(x, y int) {

	width := field.Width
	height := field.Height
	field.Matrix[x][y] = CellMine
	for i := x - 1; i <= x+1; i++ {
		if i > 0 && i < width {
			for j := y - 1; j <= y+1; j++ {
				if j > 0 && j < height && field.Matrix[i][j] != CellMine {
					field.Matrix[i][j]++
				}
			}
		}
	}
}

// SetMines fill matrix with mines
func (field *Field) SetMines() {
	width := field.Width
	height := field.Height
	mines := height*width - field.CellsLeft

	for mines > 0 {
		i := rand.Intn(width)
		j := rand.Intn(height)
		if field.Matrix[i][j] < CellMine {
			field.setMine(i, j)
		}
	}
}

// generate matrix
func generate(rs *RoomSettings) (mines int, matrix [][]int) {
	width := rs.Width
	height := rs.Height

	matrix = make([][]int, height)
	mines = int(float32(width*height) * rs.Percent)

	return
}

// NewField create new instance of field
func NewField(rs *RoomSettings) *Field {
	mines, matrix := generate(rs)
	field := &Field{
		Matrix:    matrix,
		History:   make([]Cell, 0, rs.Width*rs.Height),
		Width:     rs.Width,
		Height:    rs.Height,
		CellsLeft: rs.Width*rs.Height - mines,
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
	return x > 0 && x < field.Width && y > 0 && y < field.Height
}

// IsInside check is cell inside fueld
func (field Field) IsInside(cell *Cell) bool {
	return field.areCoordinatesRight(cell.X, cell.Y)
}
