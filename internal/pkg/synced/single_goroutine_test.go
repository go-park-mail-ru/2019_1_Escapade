package synced

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// unit
func TestInit(t *testing.T) {
	Convey("Given time duration and func", t, func() {

		var (
			d      = time.Second
			action = func() {}
		)

		Convey("When create SingleGoroutine with real time duration", func() {
			var sg = SingleGoroutine{}
			sg.Init(d, action)
			defer sg.Close()

			Convey("Then no error will happened", func() {
				So(sg.ticker, ShouldNotBeNil)
				So(sg.single, ShouldNotBeNil)
				So(sg.action, ShouldEqual, action)
			})
		})

		Convey("When create SingleGoroutine without real time duration", func() {
			var sg = SingleGoroutine{}
			sg.Init(-1, action)
			defer sg.Close()

			Convey("Then no error will happened", func() {
				So(sg.ticker, ShouldNotBeNil)
				So(sg.single, ShouldNotBeNil)
				So(sg.action, ShouldEqual, action)
			})
		})
	})
}

// integration
func TestUseCaseSingleGoroutine(t *testing.T) {
	Convey("Given action increment a", t, func() {
		var (
			a      = 0
			b      = 0
			d      = time.Millisecond * 2
			action = func() {
				a++
				time.Sleep(d)
			}
		)

		var sg = SingleGoroutine{}
		sg.Init(d, action)
		defer sg.Close()

		Convey("When action is launched s launched every 2 ms in a 50 ms interval  ", func() {
			timer1 := time.NewTimer(time.Millisecond * 10)
			var stop bool
			for i := 0; i < 100; i++ {
				go sg.Do()
			}
			for !stop {
				select {
				case <-timer1.C:
					stop = true
				case <-sg.C():
					b++
				}
			}
			c := a

			Convey("Then action will execute only 5 times", func() {
				So(c, ShouldEqual, 5)
				So(b, ShouldEqual, 5)
			})
		})
	})
}
