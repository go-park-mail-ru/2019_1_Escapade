package game

import (
	"testing"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFieldAll(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Launch everything", t, func() {

		rs := models.NewSmallRoom()
		//field.go
		field := NewField(rs)
		generate(rs)
		field.Free()
		field1 := NewField(rs)
		field2 := NewField(rs)
		field1.SameAs(field2)

		var cells []Cell
		field2.OpenEverything(&cells)
		field2.Matrix[1][1] = 0
		field2.openCellArea(1, 1, 1, &cells)
		field2.Matrix[1][1] = 3
		field2.openCellArea(1, 1, 1, &cells)
		field2.IsCleared()
		cell := NewCell(1, 1, 1, 1)
		field2.saveCell(cell, &cells)
		field2.OpenCell(cell)
		players := make([]Player, 2)
		players[0] = *NewPlayer(1)
		players[1] = *NewPlayer(2)
		field2.RandomFlags(players)
		field2.CreateRandomFlag(1)
		//field2.SetMines()
		field2.SetFlag(cell)
		field2.SetCellFlagTaken(cell)
		field2.setCellOpen(1, 2)

		field.SetMines()
		field.OpenEverything(&cells)
		field.IsCleared()
		field.OpenCell(cell)
		field.SetFlag(cell)
		field.SetCellFlagTaken(cell)
	})
}
