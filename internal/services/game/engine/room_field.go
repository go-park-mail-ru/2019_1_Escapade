package engine

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
)

type RoomField struct {
	//r  *Room
	s  SyncI
	re *RoomRecorder
	se *RoomSender
	e  *RoomEvents
	p  *RoomPeople

	Field        *Field
	isDeathmatch bool
}

func (room *RoomField) Init(s SyncI, re *RoomRecorder, se *RoomSender,
	e *RoomEvents, p *RoomPeople, field *Field, isDeathmatch bool) {
	room.s = s
	room.re = re
	room.se = se
	room.e = room.e
	room.p = room.p
	room.isDeathmatch = isDeathmatch
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
	room.Field.Fill(flags, room.isDeathmatch)
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
		if room.e.Status() != StatusRunning {
			return
		}
		// if wrong cell
		if !room.Field.IsInside(cell) {
			return
		}
		// if user died
		if !room.p.isAlive(conn) {
			return
		}

		// set who try open cell(for history)
		cell.PlayerID = conn.ID()
		cells := room.Field.OpenCell(cell)
		if len(cells) == 1 {
			newCell := cells[0]
			room.p.OpenCell(conn, &newCell)
		} else {
			for _, foundCell := range cells {
				room.p.OpenCell(conn, &foundCell)
			}
		}
		if len(cells) > 0 {
			go room.se.PlayerPoints(room.p.Players.m.Player(conn.Index()))
			go room.se.NewCells(cells...)
		}
		if room.Field.IsCleared() {
			room.e.updateStatus(StatusFinished)
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
			found, _ = room.p.flagExists(cell, nil)
		}
		room.p.Players.m.SetFlag(conn, cell, room.e.prepareOver)
		room.se.RandomFlagSet(conn, cell)
	})
}

// SetFlag handle user want set flag
func (room *RoomField) SetFlag(conn *Connection, cell *Cell) bool {
	var err error
	room.s.doWithConn(conn, func() {
		// if user try set flag after game launch
		if room.e.Status() != StatusFlagPlacing {
			err = re.ErrorBattleAlreadyBegan()
			return
		}

		if !room.Field.IsInside(cell) {
			err = re.ErrorCellOutside()
			return
		}

		if !room.p.isAlive(conn) {
			err = re.ErrorPlayerFinished()
			return
		}

		if found, prevConn := room.p.flagExists(*cell, conn); found {
			room.re.FlagСonflict(conn)

			go room.SetAndSendNewCell(conn)

			go room.SetAndSendNewCell(prevConn)
			return
		}
		room.p.Players.m.SetFlag(conn, *cell, room.e.prepareOver)
		room.re.FlagSet(conn)
	})
	if err != nil {
		room.se.FailFlagSet(conn, cell, err)
		return false
	}
	return true
}
