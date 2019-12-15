package synced

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// unit
func TestNewObserver(t *testing.T) {
	Convey("Given pairs of codes and actions", t, func() {
		pairs := []Pair{
			NewPair(1, func(Msg) { fmt.Println("1") }),
			NewPair(2, func(Msg) { fmt.Println("2") }),
		}
		Convey("When create observer without pairs", func() {
			observer := NewObserver()
			Convey("Then no error will happened", func() {
				So(observer.pairs, ShouldBeNil)
				So(observer.in, ShouldNotBeNil)
				So(observer.stop, ShouldNotBeNil)
			})
		})
		Convey("When create observer with the pairs", func() {
			observer := NewObserver(pairs...)
			Convey("Then these pairs will be set", func() {
				So(observer.pairs, ShouldResemble, pairs)
				So(observer.in, ShouldNotBeNil)
				So(observer.stop, ShouldNotBeNil)
			})
		})
	})
}

// unit
func TestAddPreAction(t *testing.T) {
	Convey("Given observer preAction func", t, func() {
		observer := NewObserver()
		p := func() { fmt.Println("1") }
		Convey("When add preAction to Observer", func() {
			o := observer.AddPreAction(p)
			Convey("Then observer's preAction will be the same as our one", func() {
				So(observer.preAction, ShouldEqual, p)
				So(observer, ShouldEqual, o)
			})
		})
	})
}

// unit
func TestAddCheck(t *testing.T) {
	Convey("Given observer check func", t, func() {
		observer := NewObserver()
		c := func() bool { return false }
		Convey("When add check to Observer", func() {
			o := observer.AddCheck(c)
			Convey("Then observer's check will be the same as our one", func() {
				So(observer.check, ShouldEqual, c)
				So(observer, ShouldEqual, o)
			})
		})
	})
}

// unit
func TestDo(t *testing.T) {
	Convey("Given observer", t, func() {
		var (
			a      = int32(50)
			result = a
			right  = int32(1)
		)
		observer := newTestObserver(&a, right, nil)

		Convey("When send message", func() {
			Convey("Then observer will process it correct", func() {
				observer.do(NewMsg(1, 2, nil))
				result += 3
				So(a, ShouldEqual, result)

				observer.do(NewMsg(1, 2, nil))
				result += 3
				So(a, ShouldEqual, result)

				observer.do(NewMsg(1, 8, nil))
				result += 3
				So(a, ShouldEqual, result)

				result = (result + 5) * 2
				observer.do(NewMsg(1, 3, nil))
				So(a, ShouldEqual, result)

				observer.do(NewMsg(1, 4, nil))
				So(a, ShouldEqual, result)
			})
		})
	})
}

// integration
func TestUseCase1(t *testing.T) {
	Convey("Given observer", t, func() {
		rand.Seed(time.Now().UnixNano())
		var (
			a        = int32(30)
			result   = a
			observer ObserverI
			right    = int32(1)
			wrong    = int32(2)
		)
		out := make(chan int32, 1)
		observer = newTestObserver(&a, right, out)
		defer close(out)

		observer.Start()

		Convey("When send message", func() {
			Convey("Then observer will process it correct", func() {
				observer.Notify(NewMsg(right, 2, nil))
				result += 3
				<-out
				So(a, ShouldEqual, result)

				observer.Notify(NewMsg(wrong, 2, nil))
				So(a, ShouldEqual, result)

				observer.Notify(NewMsg(right, 8, nil))
				result += 3
				<-out
				So(a, ShouldEqual, result)

				result = (result + 5) * 2
				observer.Notify(NewMsg(right, 3, nil))
				<-out
				So(a, ShouldEqual, result)

				observer.Notify(NewMsg(right, 4, nil))
				So(a, ShouldEqual, result)
			})
		})

		observer.Stop()
	})
}

func newTestObserver(a *int32, code int32, out chan int32) *ObserverBase {
	rand.Seed(time.Now().UnixNano())
	var (
		count = int32(10)
		pairs = make([]Pair, count)
		i     int32
	)
	for i = 0; i < count; i++ {
		if i%2 == 0 {
			pairs[i] = NewPairNoArgs(i, func() {
				*a -= 2
				if out != nil {
					out <- *a
				}
			})
		} else {
			pairs[i] = NewPairNoArgs(i, func() {
				*a *= 2
				if out != nil {
					out <- *a
				}
			})
		}
	}

	observer := NewObserver(pairs...)
	observer.AddCheck(func() bool { return *a < 80 })
	observer.AddPreAction(func() { *a += 5 })
	observer.AddPublisherCode(code)
	return observer
}
