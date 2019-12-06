package synced

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type PublisherMock struct {
	PublisherBase
	sizeChan chan int
}

func (pm *PublisherMock) init() *PublisherMock {
	pm.PublisherBase = *NewPublisher()
	pm.sizeChan = make(chan int, 1)
	return pm
}

func (pm *PublisherMock) addSub(sub ObserverI) {
	pm.PublisherBase.addSub(sub)
	pm.sizeChan <- len(pm.subscribers)
}

func (pm *PublisherMock) removeSub(sub ObserverI) {
	pm.PublisherBase.removeSub(sub)
	pm.sizeChan <- len(pm.subscribers)
}

func (pm *PublisherMock) stop() {
	pm.PublisherBase.stop()
	close(pm.sizeChan)
}

// unit
func TestNewPublisher(t *testing.T) {
	Convey("Given nothing", t, func() {

		Convey("When create publisher", func() {
			publisher := NewPublisher()
			Convey("Then there will not any nil field", func() {
				So(publisher.subscribers, ShouldNotBeNil)
				So(publisher.chanAdd, ShouldNotBeNil)
				So(publisher.chanRemove, ShouldNotBeNil)
				So(publisher.chanMsg, ShouldNotBeNil)
				So(publisher.chanStop, ShouldNotBeNil)
			})
		})
	})
}

// unit
func TestAddRemoveSub(t *testing.T) {
	Convey("Given several observers and publisher", t, func() {
		publisher := new(PublisherMock).init()
		observer1 := NewObserver(NewPairNoArgs(1, func() { fmt.Println(1) }))
		observer2 := NewObserver(NewPairNoArgs(2, func() { fmt.Println(2) }))
		observer3 := NewObserver(NewPairNoArgs(3, func() { fmt.Println(3) }))

		Convey("When subscribe observers on the publisher", func() {
			publisher.addSub(observer1)
			size := <-publisher.sizeChan
			So(size, ShouldEqual, 1)

			publisher.addSub(observer2)
			size = <-publisher.sizeChan
			So(size, ShouldEqual, 2)

			publisher.addSub(observer3)
			size = <-publisher.sizeChan
			So(size, ShouldEqual, 3)

			Convey("Then all of observers should be in publisher's subscribers list", func() {
				So(publisher.subscribers[0], ShouldEqual, observer1)
				So(publisher.subscribers[1], ShouldEqual, observer2)
				So(publisher.subscribers[2], ShouldEqual, observer3)
			})

			publisher.removeSub(observer1)
			size = <-publisher.sizeChan
			So(size, ShouldEqual, 2)

			publisher.removeSub(observer2)
			size = <-publisher.sizeChan
			So(size, ShouldEqual, 1)

			publisher.removeSub(observer3)
			size = <-publisher.sizeChan
			So(size, ShouldEqual, 0)

		})
	})
}

// functional
func TestUseCase(t *testing.T) {
	Convey("Given several observers and publisher", t, func() {
		publisher := NewPublisher()
		var (
			a, b, c               int32
			newA, newB, newC, big int32
			codeA, codeB, codeC   int32
		)
		codeA, codeB, codeC = 10, 20, 30
		newA, newB, newC, big = 30, 50, 2, 100

		observer1 := NewObserver(NewPair(1, func(msg Msg) {
			p, _ := msg.Extra.(int32)
			if a < p {
				a = p
			}
		}))
		observer1.AddPublisherCode(codeA)

		observer2 := NewObserver(NewPair(2, func(msg Msg) {
			p, _ := msg.Extra.(int32)
			if b < p {
				b = p
			}
		}))
		observer2.AddPublisherCode(codeB)

		observer3 := NewObserver(NewPair(3, func(msg Msg) {
			p, _ := msg.Extra.(int32)
			if c < p {
				c = p
			}
		}))
		observer3.AddPublisherCode(codeC)

		Convey("When subscribe observers on the publisher", func() {
			publisher.StartPublish()
			defer publisher.StopPublish()

			publisher.Observe(observer1)
			publisher.Observe(observer2)
			publisher.Observe(observer3)

			Convey("And when a publisher sends messages to subscribers", func() {
				publisher.Notify(NewMsg(codeA, 1, newA))
				publisher.Notify(NewMsg(codeB, 2, newB))
				publisher.Notify(NewMsg(codeC, 3, newC))

				publisher.Notify(NewMsg(codeB, 1, big))
				publisher.Notify(NewMsg(codeC, 2, big))
				publisher.Notify(NewMsg(codeA, 3, big))

				Convey("Then every observer is get message", func() {
					So(a, ShouldEqual, newA)
					So(b, ShouldEqual, newB)
					So(c, ShouldEqual, newC)
				})
			})
		})
	})
}
