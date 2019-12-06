package engine

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
	room_ "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine/room"
)

// FieldProxyI control access to field
// Proxy Pattern
type FieldProxyI interface {
	// for models
	ModelCells() []models.Cell
	Field() FieldI

	// for events
	SetFlag(conn *Connection, cell *Cell)
	OpenCell(conn *Connection, cell *Cell)

	SaveCell(cell *Cell) []Cell

	Configure(info models.GameInformation)
}

// RoomField implements FieldProxyI
type RoomField struct {
	s  synced.SyncI
	re ActionRecorderI
	se RSendI
	e  EventsI
	p  PeopleI

	field        *Field
	isDeathmatch bool
}

// init struct's values
func (room *RoomField) init(field *Field, isDeathmatch bool) {
	room.isDeathmatch = isDeathmatch
	room.field = field
}

// build components
func (room *RoomField) build(builder RBuilderI) {
	builder.BuildSync(&room.s)
	builder.BuildRecorder(&room.re)
	builder.BuildSender(&room.se)
	builder.BuildEvents(&room.e)
	builder.BuildPeople(&room.p)
}

// subscribe to room events
func (room *RoomField) subscribe() {
	room.eventsSubscribe()
}

// Init configure dependencies with other components of the room
func (room *RoomField) Init(builder RBuilderI, field *Field, isDeathmatch bool) {
	room.init(field, isDeathmatch)
	room.build(builder)
	room.subscribe()
}

// Field return instance of Field
func (room *RoomField) Field() FieldI {
	return room.field
}

// SaveCell save cell to history
func (room *RoomField) SaveCell(cell *Cell) []Cell {
	cells := make([]Cell, 0)
	room.field.saveCell(cell, &cells)
	return cells
}

// ModelCells cast []Cell to models.Cell
func (room *RoomField) ModelCells() []models.Cell {
	cells := make([]models.Cell, 0)
	for _, cellHistory := range room.field.History() {
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

// verify that this cell and the action on it are valid for execution
func (room *RoomField) verify(conn *Connection, cell *Cell, setFlag bool) error {

	var rightStatus int32
	if setFlag {
		rightStatus = room_.StatusFlagPlacing
	} else {
		rightStatus = room_.StatusRunning
	}
	if room.e.Status() != rightStatus {
		return re.ErrorWrongStatus()
	}
	// if wrong cell
	if !room.field.IsInside(cell) {
		return re.ErrorCellOutside()
	}
	// if user died
	if !room.p.isAlive(conn) {
		return re.ErrorPlayerFinished()
	}
	return nil
}

// OpenCell open cell
func (room *RoomField) OpenCell(conn *Connection, cell *Cell) {
	room.s.DoWithOther(conn, func() {
		if room.verify(conn, cell, false) != nil {
			return
		}
		utils.Debug(false, "open cell func")
		// set who try open cell(for history)
		cell.PlayerID = conn.ID()
		cells := room.field.OpenCell(cell)
		utils.Debug(false, "try to open", len(cells))
		if len(cells) == 0 {
			return
		}
		for _, foundCell := range cells {
			room.p.OpenCell(conn, &foundCell)
		}
		utils.Debug(false, "send me", len(cells))
		room.se.PlayerPoints(room.p.getPlayer(conn))
		room.se.NewCells(cells...)
	})
}

// SetAndSendNewCell set and send cell to conn
func (room *RoomField) SetAndSendNewCell(conn *Connection) {
	room.s.Do(func() {
		found := true
		// create until it become unique
		var cell Cell
		for found {
			cell = room.field.RandomFlag(conn.ID())
			found, _ = room.p.flagExists(cell, nil)
		}
		room.p.setFlag(conn, cell)
		room.se.RandomFlagSet(conn, &cell)
	})
}

// SetFlag handle user want set flag
func (room *RoomField) SetFlag(conn *Connection, cell *Cell) {
	var err error
	room.s.DoWithOther(conn, func() {
		if err = room.verify(conn, cell, true); err != nil {
			return
		}
		if found, prevConn := room.p.flagExists(*cell, conn); found {
			room.re.FlagСonflict(conn)
			room.SetAndSendNewCell(conn)
			room.SetAndSendNewCell(prevConn)
			return
		}
		room.p.setFlag(conn, *cell)
		room.re.FlagSet(conn)
	})
	if err != nil {
		room.se.FailFlagSet(conn, cell, err)
	}
}

// Configure set matrix of cells and history of open cells
func (room *RoomField) Configure(info models.GameInformation) {
	room.configureField(info.Field)
	room.configureHistory(info.Cells)
}

// configureField set matrix of cells
func (room *RoomField) configureField(info models.Field) {
	room.field.Width = info.Width
	room.field.Height = info.Height
	room.field.setCellsLeft(info.CellsLeft)
	room.field.Mines = info.Mines
}

// configureHistory set history of open cells
func (room *RoomField) configureHistory(info []models.Cell) {
	room.field.setHistory(make([]*Cell, 0))
	for _, cellDB := range info {
		cell := &Cell{
			X:        cellDB.X,
			Y:        cellDB.Y,
			Value:    cellDB.Value,
			PlayerID: cellDB.PlayerID,
			Time:     cellDB.Date,
		}
		room.field.setToHistory(cell)
	}
}

// eventsRunning is called when game begins
func (room *RoomField) eventsRunning(synced.Msg) {
	cells := room.Field().OpenZero() //room.Field.OpenSave(int(open))
	room.se.NewCells(cells...)
}

// eventsFinished is called when game finished
func (room *RoomField) eventsFinished(msg synced.Msg) {
	cells := make([]Cell, 0)
	room.Field().OpenEverything(cells)
	results, ok := msg.Extra.(room_.FinishResults)
	if !ok {
		return
	}
	room.se.GameOver(results.Timeout, room.se.All, cells)
}

// eventsAborted is called when room wanna to clear resources
func (room *RoomField) eventsAborted(synced.Msg) {
	room.Field().Free()
}

// eventsSubscribe subscibe to events associated with room's status
func (room *RoomField) eventsSubscribe() {
	observer := synced.NewObserver(
		synced.NewPair(room_.StatusRunning, room.eventsRunning),
		synced.NewPair(room_.StatusFinished, room.eventsFinished),
		synced.NewPair(room_.StatusAborted, room.eventsAborted))
	room.e.Observe(observer.AddPublisherCode(room_.UpdateStatus))
}

// 234 -> 192
