package engine

import (
	"math/rand"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewPlayerAction(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Given some playerID, actionID and time 'before' creating PlayerAction", t, func() {
		before := time.Now()
		rand.Seed(before.UnixNano())
		playerID := rand.Intn(10000000)
		actionID := rand.Intn(15)

		Convey("When the playerAction is created and time 'after' set", func() {
			playerAction := NewPlayerAction(playerID, actionID)
			after := time.Now()

			Convey("The field 'Player' should be the same as playerID", func() {
				So(playerAction.Player, ShouldEqual, playerID)
			})
			Convey("The field 'Action' should be the same as actionID", func() {
				So(playerAction.Action, ShouldEqual, actionID)
			})
			Convey("The field 'Time' should be between 'before' and 'after'", func() {
				So(playerAction.Time, ShouldHappenBetween, before, after)
			})
		})
	})
}
