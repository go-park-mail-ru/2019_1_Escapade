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

	EventsSub() synced.SubscriberI
}

// RoomField implements FieldProxyI
type RoomField struct {
	//r  *Room
	s  synced.SyncI
	re ActionRecorderI
	se RSendI
	e  EventsI
	p  PeopleI

	field        *Field
	isDeathmatch bool
}

// Init configure dependencies with other components of the room
func (room *RoomField) Init(builder RBuilderI, field *Field,
	isDeathmatch bool) {
	builder.BuildSync(&room.s)
	builder.BuildRecorder(&room.re)
	builder.BuildSender(&room.se)
	builder.BuildEvents(&room.e)
	builder.BuildPeople(&room.p)

	room.isDeathmatch = isDeathmatch
	room.field = field
}

func (room *RoomField) Field() FieldI {
	return room.field
}

func (room *RoomField) SaveCell(cell *Cell) []Cell {
	cells := make([]Cell, 0)
	room.field.saveCell(cell, &cells)
	return cells
}

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

func (room *RoomField) check(conn *Connection, cell *Cell, setFlag bool) error {

	var rightStatus = room_.StatusRunning
	if setFlag {
		rightStatus = room_.StatusFlagPlacing
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
		if room.check(conn, cell, false) != nil {
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

		room.e.tryFinish()
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
		if err = room.check(conn, cell, true); err != nil {
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

func (room *RoomField) Configure(info models.GameInformation) {
	room.configureField(info.Field)
	room.configureHistory(info.Cells)
}

func (room *RoomField) configureField(info models.Field) {
	room.field.Width = info.Width
	room.field.Height = info.Height
	room.field.setCellsLeft(info.CellsLeft)
	room.field.Mines = info.Mines
}

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

func (room *RoomField) EventsSub() synced.SubscriberI {
	return synced.NewSubscriber(room.eventsCallback)
}

func (room *RoomField) eventsCallback(msg synced.Msg) {
	if msg.Code != room_.UpdateStatus {
		return
	}
	code, ok := msg.Content.(int)
	if !ok {
		return
	}
	switch code {
	case room_.StatusRunning:
		cells := room.Field().OpenZero() //room.Field.OpenSave(int(open))
		room.se.NewCells(cells...)
	case room_.StatusFinished:
		cells := make([]Cell, 0)
		room.Field().OpenEverything(cells)
		room.se.GameOver(room.e.Timeout(), room.se.All, cells)
	case room_.StatusAborted:
		room.Field().Free()
	}

}

// 234 -> 192
