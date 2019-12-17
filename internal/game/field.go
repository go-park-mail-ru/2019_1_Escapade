package game

import (
	"sync"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"math"
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
	History  []*Cell `json:"history"`
	Width    int     `json:"width"`
	Height   int     `json:"height"`

	cellsLeftM *sync.RWMutex
	CellsLeft  int `json:"-"`

	config *config.FieldConfig

	Mines     int     `json:"mines"`
	Difficult float64 `json:"difficult"`
}

// NewField create new instance of field
func NewField(rs *models.RoomSettings, config *config.FieldConfig) *Field {
	matrix := generate(rs)
	field := &Field{
		wGroup: &sync.WaitGroup{},

		doneM: &sync.RWMutex{},
		done:  false,

		matrixM: &sync.RWMutex{},
		Matrix:  matrix,

		historyM: &sync.Mutex{},
		History:  make([]*Cell, 0, rs.Width*rs.Height),
		Width:    rs.Width,
		Height:   rs.Height,
		Mines:    rs.Mines,

		config: config,

		cellsLeftM: &sync.RWMutex{},
		CellsLeft:  rs.Width*rs.Height - rs.Mines - rs.Players,
	}
	field.Difficult = float64(field.Mines) / float64(field.CellsLeft)
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
	go field.matrixFree()
	go field.historyFree()
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

	var flagsFound int

	for i := 0; i < field.Height; i++ {
		for j := 0; j < field.Width; j++ {
			v := field.getMatrixValue(i, j)
			//fmt.Println("cell ", i, j, v)
			if v > CellIncrement {
				//fmt.Println("flag!!!!", i, j, v)
				flagsFound++
			}
			//if v == CellFlagTaken {
			//	fmt.Println("CellFlagTaken!!!!", i, j, v)
			//}
			if v != CellOpened && v != CellFlagTaken {
				cell := NewCell(i, j, v, 0)
				field.saveCell(cell, cells)
				field.setCellOpen(i, j, v)
			}
		}
	}
	//fmt.Println("flagsFound", flagsFound)
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
		field.setToHistory(cell)
		*cells = append(*cells, *cell)
		field.setCellOpen(cell.X, cell.Y, cell.Value)
	}
	//fmt.Printf("save openCellArea Cell(%d/%d)=%d\n", cell.X, cell.Y, cell.Value)
}

// OpenCell open cell and return slice of opened cells
func (field *Field) OpenCell(cell *Cell) (cells []Cell) {
	if field.getDone() {
		return
	}
	field.wGroup.Add(1)
	defer field.wGroup.Done()

	cell.Value = field.getMatrixValue(cell.X, cell.Y)
	//fmt.Printf("Cell(%d/%d)=%d", cell.X, cell.Y, cell.Value)

	cells = make([]Cell, 0)
	if cell.Value < CellMine {
		field.openCellArea(cell.X, cell.Y, cell.PlayerID, &cells)
	} else {
		if cell.Value != FlagID(cell.PlayerID) {
			field.saveCell(cell, &cells)
		}
	}

	return
}

// RandomFlags create random players flags
func (field *Field) RandomFlags(players []Player) (flags []Flag) {
	if field.getDone() {
		return
	}
	field.wGroup.Add(1)
	defer field.wGroup.Done()

	flags = make([]Flag, len(players))
	for i, player := range players {
		flags[i] = Flag{
			Cell: field.CreateRandomFlag(player.ID),
			Set:  false,
		}
	}
	return flags
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
	cell = *NewCell(x, y, FlagID(playerID), playerID)

	return cell
}

// OpenSave open n(or more) cells that do not contain any mines or flags
func (field *Field) OpenSave(n int) (cells []Cell) {
	cells = make([]Cell, 0)
	size := 0

	for n > size {
		rand.Seed(time.Now().UnixNano())
		i := rand.Intn(field.Width)
		j := rand.Intn(field.Height)
		if field.lessThenMine(i, j) {
			cells = append(cells, field.OpenCell(NewCell(i, j, 0, 0))...)
			size = len(cells)
		}
	}
	return
}

func (field *Field) OpenZero() (cells []Cell) {
	cells = make([]Cell, 0)
	//size := 0

	for i := 0; i < field.Width; i++ {
		for j := 0; j < field.Height; j++ {
			if field.getMatrixValue(i, j) == 0 {
				cells = append(cells, field.OpenCell(NewCell(i, j, 0, 0))...)
			}
		}
	}
	return
}

// Zero clears the entire matrix of values
func (field *Field) Zero() {
	for i := 0; i < field.Width; i++ {
		for j := 0; j < field.Height; j++ {
			field.setMatrixValue(i, j, 0)
		}
	}
}

// SetMines fill matrix with mines
func (field *Field) SetMines(flags []Flag, deathmatch bool) {
	if field.getDone() {
		return
	}
	field.wGroup.Add(1)
	defer field.wGroup.Done()

	var (
		width  = field.Width
		height = field.Height
		mines  = field.Mines
	)

	if deathmatch {
		var (
			gip      = float64(field.Width*field.Width) + float64(field.Height*field.Height)
			areaSize = gip / field.Difficult / 2000.0

			minAreaSize = float64(field.config.MinAreaSize)
			maxAreaSize = float64(field.config.MaxAreaSize)
		)

		if areaSize < minAreaSize {
			areaSize = minAreaSize
		} else if areaSize > maxAreaSize {
			areaSize = maxAreaSize
		}
		var (
			areaSizeINT    = int(areaSize)
			probability    = 5 * int(areaSize*areaSize)
			minProbability = field.config.MinProbability
			maxProbability = field.config.MaxProbability
		)
		if probability > maxProbability {
			probability = maxProbability
		} else if probability < minProbability {
			probability = minProbability
		}
		utils.Debug(false, "we have %d flags, area size %f; probability %d;;;;%f\n", len(flags), gip/field.Difficult/2000.0, probability, field.Difficult)

		for _, flag := range flags {
			x := flag.Cell.X
			y := flag.Cell.Y
			utils.Debug(false, "flag[%d, %d] - %d\n", x, y, flag.Cell.PlayerID)
			for i := x - areaSizeINT; i <= x+areaSizeINT; i++ {
				if i >= 0 && i < width {
					for j := y - areaSizeINT; j <= y+areaSizeINT; j++ {
						if j >= 0 && j < height {
							if field.getMatrixValue(i, j) != CellMine && !(x == i && y == j) {
								rand.Seed(time.Now().UnixNano())
								procent := rand.Intn(100)
								utils.Debug(false, "[%d, %d] - %d\n", i, j, procent)
								if procent > probability {
									field.setMatrixValue(i, j, CellMine)
									mines--
									if mines == 0 {
										return
									}
								}
							}
						}
					}
				}
			}
		}
	}
	for mines > 0 {

		rand.Seed(time.Now().UnixNano())
		someX := rand.Intn(width)
		someY := rand.Intn(height)

		if field.lessThenMine(someX, someY) {
			field.setMatrixValue(someX, someY, CellMine)
			mines--
		}
	}
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

// FlagID convert player ID to Flag ID
func FlagID(connID int) int {
	// if connID == 0 {
	// 	panic("connID")
	// }
	return int(math.Abs(float64(connID))) + CellIncrement
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
	utils.Debug(false, "setFlag", cell.X, cell.Y, cell.PlayerID, CellIncrement)

	field.setMatrixValue(cell.X, cell.Y, FlagID(cell.PlayerID))
}

// SetCellFlagTaken set cells flag taken
func (field *Field) SetCellFlagTaken(cell *Cell) {
	if field.getDone() {
		return
	}
	field.wGroup.Add(1)
	defer field.wGroup.Done()

	field.setMatrixValue(cell.X, cell.Y, CellFlagTaken)

	utils.Debug(false, "flag found!", cell.Value)
}

// setCellOpen set cell opened
func (field *Field) setCellOpen(x, y, v int) {
	if v < CellMine {
		field.setMatrixValue(x, y, CellOpened)
	} else {
		field.setMatrixValue(x, y, CellFlagTaken)
	}
}

func (field *Field) SetMinesLabels() {

	width := field.Width
	height := field.Height

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if field.Matrix[x][y] != 0 {
				continue
			}
			value := 0

			for i := x - 1; i <= x+1; i++ {
				if i >= 0 && i < width {
					for j := y - 1; j <= y+1; j++ {
						if j >= 0 && j < height {
							if field.Matrix[i][j] == CellMine {
								value++
							}
						}
					}
				}
			}
			field.Matrix[x][y] = value
		}
	}
}
