package engine

import (
	"testing"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
	. "github.com/smartystreets/goconvey/convey"
)

// connection.go unit tests

var (
	userDummy = &models.UserPublicInfo{
		ID:   1,
		Name: "test",
	}
	ws = &WebsocketConnStub{}
)

func TestNewConnectionCorrect(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Create new connection", t, func() {
		var (
			ws        = &WebsocketConnStub{}
			conn, err = NewConnection(ws, userDummy)
		)
		So(err, ShouldBeNil)
		So(conn, ShouldNotBeNil)
		So(conn._ws, ShouldEqual, ws)
		So(conn.User, ShouldEqual, userDummy)
	})
}

func TestNewConnectionWrong(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Create new connection", t, func() {
		var (
			conn, err = NewConnection(nil, userDummy)
		)
		So(err, ShouldResemble, re.NoWebSocketOrUser())
		So(conn, ShouldBeNil)

		conn, err = NewConnection(ws, nil)

		So(err, ShouldResemble, re.NoWebSocketOrUser())
		So(conn, ShouldBeNil)
	})
}

func TestRestore(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Given two connections", t, func() {

		var (
			first, second      *Connection
			ws                 = &WebsocketConnStub{}
			pRoom1, pRoom2     = &Room{}, &Room{}
			wRoom1, wRoom2     = &Room{}, &Room{}
			oldIndex, newIndex = 1, 2
			err                error
		)
		first, err = NewConnection(ws, userDummy)
		So(err, ShouldBeNil)

		first.setPlayingRoom(pRoom1)
		first.setWaitingRoom(wRoom1)
		first.SetIndex(oldIndex)

		So(first.PlayingRoom(), ShouldEqual, pRoom1)
		So(first.WaitingRoom(), ShouldEqual, wRoom1)
		So(first.Index(), ShouldEqual, oldIndex)

		second = newConnection()
		second.setPlayingRoom(pRoom2)
		second.setWaitingRoom(wRoom2)
		second.SetIndex(newIndex)

		Convey("When restore first connection from second", func() {

			first.Restore(second)
			Convey("first conn's index, playing and waiting room index must be the same as second", func() {
				So(first.PlayingRoom(), ShouldEqual, pRoom2)
				So(first.WaitingRoom(), ShouldEqual, wRoom2)
				So(first.Index(), ShouldEqual, newIndex)
			})
		})
	})
}

func TestIsAnonymous(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Given connection", t, func() {

		var (
			conn *Connection
			err  error
		)
		conn, err = NewConnection(ws, userDummy)
		So(err, ShouldBeNil)

		conn.User.ID = 10
		So(conn.IsAnonymous(), ShouldBeFalse)

		conn.User.ID = -2
		So(conn.IsAnonymous(), ShouldBeTrue)
	})
}

func TestPushToRoom(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Test PushToRoom", t, func() {
		var (
			conn                   *Connection
			ws                     = &WebsocketConnStub{}
			pRoom1, wRoom1, pRoom2 = &Room{}, &Room{}, &Room{}
			err                    error
		)
		conn, err = NewConnection(ws, userDummy)
		So(err, ShouldBeNil)

		conn.setPlayingRoom(pRoom1)
		conn.setWaitingRoom(wRoom1)

		So(conn.PlayingRoom(), ShouldEqual, pRoom1)
		So(conn.WaitingRoom(), ShouldEqual, wRoom1)

		conn.PushToRoom(pRoom2)

		So(conn.PlayingRoom(), ShouldEqual, pRoom2)
		So(conn.WaitingRoom(), ShouldBeNil)
	})
}

func TestPushToLobby(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Test PushToLobby", t, func() {
		var (
			conn           *Connection
			pRoom1, wRoom1 = &Room{}, &Room{}
			err            error
		)
		conn, err = NewConnection(ws, userDummy)
		So(err, ShouldBeNil)

		conn.setPlayingRoom(pRoom1)
		conn.setWaitingRoom(wRoom1)

		So(conn.PlayingRoom(), ShouldEqual, pRoom1)
		So(conn.WaitingRoom(), ShouldEqual, wRoom1)

		conn.PushToLobby()

		So(conn.PlayingRoom(), ShouldBeNil)
		So(conn.WaitingRoom(), ShouldBeNil)
	})
}

func TestFree(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Test TestFree", t, func() {
		var (
			conn *Connection
			err  error
		)
		conn, err = NewConnection(ws, userDummy)
		So(err, ShouldBeNil)
		conn.Free()
	})
}
