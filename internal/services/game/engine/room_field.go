package engine

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
)

// FieldProxyI control access to field
// Proxy Pattern
type FieldProxyI interface {
	// for models
	Model() models.Field
	ModelCells() []models.Cell
	JSON() FieldJSON

	// for events
	Fill(flags []Flag)
	OpenZero() []Cell
	OpenEverything(cells []Cell)
	Free(wait time.Duration)
	SetFlag(conn *Connection, cell *Cell)
	OpenCell(conn *Connection, cell *Cell)

	saveCell(cell *Cell) []Cell

	RandomFlag(conn *Connection) Cell

	IsCleared() bool
	cellsLeft() int32
	difficult() float64

	Configure(info models.GameInformation)
}

// RoomField implements FieldProxyI
type RoomField struct {
	//r  *Room
	s  SyncI
	re ActionRecorderProxyI
	se SendStrategyI
	e  EventsI
	p  PeopleI

	Field        *Field
	isDeathmatch bool
}

// Init configure dependencies with other components of the room
func (room *RoomField) Init(builder ComponentBuilderI, field *Field,
	isDeathmatch bool) {
	builder.BuildSync(&room.s)
	builder.BuildRecorder(&room.re)
	builder.BuildSender(&room.se)
	builder.BuildEvents(&room.e)
	builder.BuildPeople(&room.p)

	room.isDeathmatch = isDeathmatch
	room.Field = field
}

func (room *RoomField) Free(wait time.Duration) {
	go room.Field.Free(wait)
}

func (room *RoomField) saveCell(cell *Cell) []Cell {
	cells := make([]Cell, 0)
	room.Field.saveCell(cell, cells)
	return cells
}

func (room *RoomField) ModelCells() []models.Cell {
	cells := make([]models.Cell, 0)
	for _, cellHistory := range room.Field.History() {
		cell := models.Cell{
			PlayerID: cellHistory.PlayerID,
			X:        cellHistory.X,
			Y:        cellHistory.Y,
			Value:    cellHistory.Value,
			Date:     cellHistory.Time,
		}
		cells = append(cells, cell)
	}
	return cells
}

func (room *RoomField) Fill(flags []Flag) {
	room.Field.Fill(flags, room.isDeathmatch)
}

func (room *RoomField) OpenZero() []Cell {
	return room.Field.OpenZero()
}

func (room *RoomField) OpenEverything(cells []Cell) {
	room.Field.OpenEverything(cells)
}

func (room *RoomField) Model() models.Field {
	return models.Field{
		Width:     room.Field.Width,
		Height:    room.Field.Height,
		CellsLeft: room.Field._cellsLeft,
		Difficult: 0,
		Mines:     room.Field.Mines,
	}
}

func (room *RoomField) JSON() FieldJSON {
	return room.Field.JSON()
}

func (room *RoomField) RandomFlag(conn *Connection) Cell {
	return room.Field.CreateRandomFlag(conn.ID())
}

func (room *RoomField) check(conn *Connection, cell *Cell, setFlag bool) error {

	var rightStatus = StatusRunning
	if setFlag {
		rightStatus = StatusFlagPlacing
	}
	if room.e.Status() != rightStatus {
		return re.ErrorWrongStatus()
	}
	// if wrong cell
	if !room.Field.IsInside(cell) {
		return re.ErrorCellOutside()
	}
	// if user died
	if !room.p.isAlive(conn) {
		return re.ErrorPlayerFinished()
	}
	return nil
}

func (room *RoomField) cellsLeft() int32 {
	return room.Field.cellsLeft()
}

func (room *RoomField) difficult() float64 {
	return room.Field.Difficult
}

// OpenCell open cell
func (room *RoomField) OpenCell(conn *Connection, cell *Cell) {
	room.s.doWithConn(conn, func() {
		if room.check(conn, cell, false) != nil {
			return
		}

		// set who try open cell(for history)
		cell.PlayerID = conn.ID()
		cells := room.Field.OpenCell(cell)
		if len(cells) == 0 {
			return
		}
		for _, foundCell := range cells {
			room.p.OpenCell(conn, &foundCell)
		}

		go room.se.PlayerPoints(room.p.getPlayer(conn))
		go room.se.NewCells(cells...)

		room.e.tryFinish()
	})
}

func (room *RoomField) IsCleared() bool {
	return room.Field.IsCleared()
}

// SetAndSendNewCell set and send cell to conn
func (room *RoomField) SetAndSendNewCell(conn *Connection) {
	room.s.do(func() {
		found := true
		// create until it become unique
		var cell Cell
		for found {
			cell = room.Field.CreateRandomFlag(conn.ID())
			found, _ = room.p.flagExists(cell, nil)
		}
		room.p.setFlag(conn, cell)
		room.se.RandomFlagSet(conn, &cell)
	})
}

// SetFlag handle user want set flag
func (room *RoomField) SetFlag(conn *Connection, cell *Cell) {
	var err error
	room.s.doWithConn(conn, func() {
		if err = room.check(conn, cell, true); err != nil {
			return
		}
		if found, prevConn := room.p.flagExists(*cell, conn); found {
			room.re.FlagСonflict(conn)
			go room.SetAndSendNewCell(conn)
			go room.SetAndSendNewCell(prevConn)
			return
		}
		room.p.setFlag(conn, *cell)
		room.re.FlagSet(conn)
	})
	if err != nil {
		room.se.FailFlagSet(conn, cell, err)
	}
}

func (room *RoomField) Configure(info models.GameInformation) {
	room.configureField(info.Field)
	room.configureHistory(info.Cells)
}

func (room *RoomField) configureField(info models.Field) {
	room.Field.Width = info.Width
	room.Field.Height = info.Height
	room.Field.setCellsLeft(info.CellsLeft)
	room.Field.Mines = info.Mines
}

func (room *RoomField) configureHistory(info []models.Cell) {
	room.Field.setHistory(make([]*Cell, 0))
	for _, cellDB := range info {
		cell := &Cell{
			X:        cellDB.X,
			Y:        cellDB.Y,
			Value:    cellDB.Value,
			PlayerID: cellDB.PlayerID,
			Time:     cellDB.Date,
		}
		room.Field.setToHistory(cell)
	}
}
