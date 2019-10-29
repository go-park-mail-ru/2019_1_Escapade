package game

import (
	"math/rand"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

const RANDOMSIZE = 1000

func TestNewCell(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Given some coordinates(x,y), value, playerID and time 'before'", t, func() {
		before := time.Now()
		rand.Seed(before.UnixNano())
		x := rand.Intn(RANDOMSIZE)
		y := rand.Intn(RANDOMSIZE)
		v := rand.Intn(RANDOMSIZE)
		id := rand.Intn(RANDOMSIZE)

		Convey("When the cell is created and time 'after' set", func() {
			cell := NewCell(x, y, v, id)
			after := time.Now()

			Convey("The field 'X' should be the same as x", func() {
				So(cell.X, ShouldEqual, x)
			})
			Convey("The field 'Y' should be the same as y", func() {
				So(cell.Y, ShouldEqual, y)
			})
			Convey("The field 'Value' should be the same as v", func() {
				So(cell.Value, ShouldEqual, v)
			})
			Convey("The field 'PlayerID' should be the same as id", func() {
				So(cell.PlayerID, ShouldEqual, id)
			})
			Convey("The field 'Time' should be between 'before' and 'after'", func() {
				So(cell.Time, ShouldHappenBetween, before, after)
			})
		})
	})
}
