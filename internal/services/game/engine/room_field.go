package engine

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
)

type RoomField struct {
	r *Room
	s SyncI

	Field *Field
}

func (room *RoomField) Init(r *Room, s SyncI, field *Field) {
	room.r = r
	room.s = s
	room.Field = field
}

func (room *RoomField) Free(wait time.Duration) {
	go room.Field.Free(wait)
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
	room.r.field.Field.Fill(flags, room.r.Settings.Deathmatch)
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

// OpenCell open cell
func (room *RoomField) OpenCell(conn *Connection, cell *Cell) {
	room.s.doWithConn(conn, func() {
		// if user try set open cell before game launch
		if room.r.events.Status() != StatusRunning {
			return
		}
		// if wrong cell
		if !room.Field.IsInside(cell) {
			return
		}
		// if user died
		if !room.r.people.isAlive(conn) {
			return
		}

		// set who try open cell(for history)
		cell.PlayerID = conn.ID()
		cells := room.Field.OpenCell(cell)
		if len(cells) == 1 {
			newCell := cells[0]
			room.r.people.OpenCell(conn, &newCell)
		} else {
			for _, foundCell := range cells {
				room.r.people.OpenCell(conn, &foundCell)
			}
		}
		if len(cells) > 0 {
			go room.r.send.PlayerPoints(room.r.people.Players.m.Player(conn.Index()), room.r.All)
			go room.r.send.NewCells(room.r.All, cells...)
		}
		if room.Field.IsCleared() {
			room.r.events.updateStatus(StatusFinished)
		}
	})
}

// SetAndSendNewCell set and send cell to conn
func (room *RoomField) SetAndSendNewCell(conn *Connection) {
	room.s.do(func() {
		found := true
		// create until it become unique
		var cell Cell
		for found {
			cell = room.Field.CreateRandomFlag(conn.ID())
			found, _ = room.r.people.flagExists(cell, nil)
		}
		room.r.people.Players.m.SetFlag(conn, cell, room.r.events.prepareOver)
		room.r.send.RandomFlagSet(conn, cell)
	})
}

// SetFlag handle user want set flag
func (room *RoomField) SetFlag(conn *Connection, cell *Cell) bool {
	var err error
	room.s.doWithConn(conn, func() {
		// if user try set flag after game launch
		if room.r.events.Status() != StatusFlagPlacing {
			err = re.ErrorBattleAlreadyBegan()
			return
		}

		if !room.Field.IsInside(cell) {
			err = re.ErrorCellOutside()
			return
		}

		if !room.r.people.isAlive(conn) {
			err = re.ErrorPlayerFinished()
			return
		}

		if found, prevConn := room.r.people.flagExists(*cell, conn); found {
			room.r.connEvents.notify.FlagСonflict(conn)

			go room.SetAndSendNewCell(conn)

			go room.SetAndSendNewCell(prevConn)
			return
		}
		room.r.people.Players.m.SetFlag(conn, *cell, room.r.events.prepareOver)
		room.r.connEvents.notify.FlagSet(conn)
	})
	if err != nil {
		room.r.send.FailFlagSet(conn, cell, err)
		return false
	}
	return true
}
