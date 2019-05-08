package game

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"net/http"

	"net/http/httptest"
	"strings"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
	"github.com/gorilla/websocket"

	. "github.com/smartystreets/goconvey/convey"
)

var upgrader = websocket.Upgrader{}
var ready = make(chan struct{})

var (
	TestConnection *Connection
)

// connection.go tests

func TestNewConnectionWithoutDatabase(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Given websocket and some user", t, func() {

		s := httptest.NewServer(http.HandlerFunc(echo))
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http")

		ws, _, err := websocket.DefaultDialer.Dial(u, nil)
		Convey("When websocket dials, the error should be nil", func() {
			So(err, ShouldBeNil)
		})

		defer ws.Close()

		id := rand.Intn(10000)
		user := createRandomUser(id)

		Convey("When create new connection", func() {
			lobby := NewLobby(RANDOMSIZE, RANDOMSIZE, nil, true)
			conn := NewConnection(ws, user, lobby)

			Convey("All pointers fields should be not nil", func() {
				So(conn.wGroup, ShouldNotBeNil)
				So(conn.doneM, ShouldNotBeNil)
				So(conn.roomM, ShouldNotBeNil)
				So(conn.disconnectedM, ShouldNotBeNil)
				So(conn.bothM, ShouldNotBeNil)
				So(conn.indexM, ShouldNotBeNil)
				So(conn.context, ShouldNotBeNil)
				So(conn.cancel, ShouldNotBeNil)
				So(conn.actionSem, ShouldNotBeNil)
				So(conn.send, ShouldNotBeNil)
			})
			Convey("All not pointers fields should be default", func() {
				So(conn._done, ShouldBeFalse)
				So(conn._room, ShouldBeNil)
				So(conn._Disconnected, ShouldBeFalse)
				So(conn._both, ShouldBeFalse)
				So(conn._Index, ShouldEqual, -1)
				So(conn.ws, ShouldEqual, ws)
			})
		})
	})
}

func TestNewConnectionWithoutInputParameters(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Given websocket and some user", t, func() {

		s := httptest.NewServer(http.HandlerFunc(echo))
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http")

		ws, _, err := websocket.DefaultDialer.Dial(u, nil)
		Convey("When websocket dials, the error should be nil", func() {
			So(err, ShouldBeNil)
		})

		defer ws.Close()

		id := rand.Intn(10000)
		user := createRandomUser(id)
		lobby := NewLobby(RANDOMSIZE, RANDOMSIZE, nil, true)

		Convey("When create new connection without lobby", func() {
			conn := NewConnection(ws, user, nil)

			Convey("connection should be nil", func() {
				So(conn, ShouldBeNil)
			})
		})
		Convey("When create new connection without user", func() {
			conn := NewConnection(ws, nil, lobby)

			Convey("connection should be nil", func() {
				So(conn, ShouldBeNil)
			})
		})
		Convey("When create new connection without ws", func() {
			conn := NewConnection(nil, user, lobby)

			Convey("connection should be nil", func() {
				So(conn, ShouldBeNil)
			})
		})
	})
}

func TestPushToRoom(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Given connection", t, func() {

		s := httptest.NewServer(http.HandlerFunc(echo))
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http")

		ws, _, err := websocket.DefaultDialer.Dial(u, nil)
		Convey("When websocket dials, the error should be nil", func() {
			So(err, ShouldBeNil)
		})

		defer ws.Close()

		id := rand.Intn(10000)
		user := createRandomUser(id)
		lobby := NewLobby(RANDOMSIZE, RANDOMSIZE, nil, true)

		conn := NewConnection(ws, user, lobby)

		roomID := utils.RandomString(16)
		settings := models.NewSmallRoom()
		room, err := NewRoom(settings, roomID, lobby)
		Convey("When create room, the error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("When push to room", func() {
			conn.PushToRoom(room)

			Convey("connection's room should be the room", func() {
				So(conn.Room(), ShouldEqual, room)
			})
		})

		Convey("When done and push to room", func() {
			conn.PushToLobby()
			conn.setDone()
			conn.PushToRoom(room)

			Convey("connection's room should be nil", func() {
				So(conn.Room(), ShouldBeNil)
			})
		})
	})
}

func TestPushToLobby(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Given connection", t, func() {

		s := httptest.NewServer(http.HandlerFunc(echo))
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http")

		ws, _, err := websocket.DefaultDialer.Dial(u, nil)
		Convey("When websocket dials, the error should be nil", func() {
			So(err, ShouldBeNil)
		})

		defer ws.Close()

		id := rand.Intn(10000)
		user := createRandomUser(id)
		lobby := NewLobby(RANDOMSIZE, RANDOMSIZE, nil, true)

		conn := NewConnection(ws, user, lobby)

		roomID := utils.RandomString(16)
		settings := models.NewSmallRoom()
		room, err := NewRoom(settings, roomID, lobby)
		Convey("When create room, the error should be nil", func() {
			So(err, ShouldBeNil)
		})
		Convey("When push to room", func() {
			conn.PushToRoom(room)

			Convey("connection's room should be the room", func() {
				So(conn.Room(), ShouldEqual, room)
			})
		})
		Convey("When push to lobby", func() {
			conn.PushToLobby()

			Convey("connection's room should be nil", func() {
				So(conn.Room(), ShouldBeNil)
			})
		})

		Convey("When done and push to lobby", func() {
			conn.PushToRoom(room)

			Convey("connection's room should be room", func() {
				So(conn.Room(), ShouldEqual, room)
			})
		})
	})
}

func TestIsConnected(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Given connection", t, func() {

		s := httptest.NewServer(http.HandlerFunc(echo))
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http")

		ws, _, err := websocket.DefaultDialer.Dial(u, nil)
		Convey("When websocket dials, the error should be nil", func() {
			So(err, ShouldBeNil)
		})

		defer ws.Close()

		id := rand.Intn(10000)
		user := createRandomUser(id)
		lobby := NewLobby(RANDOMSIZE, RANDOMSIZE, nil, true)

		conn := NewConnection(ws, user, lobby)

		Convey("When is connected", func() {
			v := conn.IsConnected()

			Convey("'disconnected' should be false", func() {
				So(conn.Disconnected(), ShouldNotEqual, v)
			})
		})
		conn.setDisconnected()
		Convey("When is disconnected", func() {
			v := conn.IsConnected()

			Convey("'disconnected' should be true", func() {
				So(conn.Disconnected(), ShouldNotEqual, v)
			})
		})

		Convey("When done and connected", func() {
			conn._Disconnected = false
			conn.setDone()
			conn.setDisconnected()
			v := conn.IsConnected()

			Convey("'disconnected' should be false", func() {
				So(conn.Disconnected(), ShouldEqual, v)
			})
		})

		Convey("When done and disconnected", func() {
			conn.setDisconnected()
			conn.setDone()
			v := conn.IsConnected()

			Convey("'disconnected' should be false", func() {
				So(conn.Disconnected(), ShouldEqual, v)
			})
		})
	})
}

func TestDirty(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Given connection", t, func() {

		s := httptest.NewServer(http.HandlerFunc(echo))
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http")

		ws, _, err := websocket.DefaultDialer.Dial(u, nil)
		Convey("When websocket dials, the error should be nil", func() {
			So(err, ShouldBeNil)
		})

		defer ws.Close()

		id := rand.Intn(10000)
		user := createRandomUser(id)
		lobby := NewLobby(RANDOMSIZE, RANDOMSIZE, nil, true)

		conn := NewConnection(ws, user, lobby)

		Convey("the id should be equal id", func() {
			So(conn.User.ID, ShouldEqual, id)
		})
		Convey("When is made dirty", func() {
			conn.Dirty()

			Convey("the id should be -1", func() {
				So(conn.User.ID, ShouldEqual, -1)
			})
		})
		Convey("When done and disconnected", func() {
			conn.User.ID = id
			conn.setDone()
			conn.Dirty()

			Convey("the id should be id", func() {
				So(conn.User.ID, ShouldEqual, id)
			})
		})
	})
}

func TestKill(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Given connection", t, func() {

		s := httptest.NewServer(http.HandlerFunc(echo))
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http")

		ws, _, err := websocket.DefaultDialer.Dial(u, nil)

		So(err, ShouldBeNil)

		defer ws.Close()

		time.Sleep(2 * time.Second)
		<-ready

		id := TestConnection.User.ID

		Convey("When push to room", func() {
			roomID := utils.RandomString(16)
			settings := models.NewSmallRoom()
			room, err := NewRoom(settings, roomID, lobby)

			So(err, ShouldBeNil)
			TestConnection.PushToRoom(room)

			Convey("connection's room should be the room", func() {
				So(TestConnection.Room(), ShouldEqual, room)
			})
		})
		Convey("When connection is killed without dirty", func() {
			fmt.Println("kill start")
			TestConnection.Kill(utils.RandomString(16), false)
			fmt.Println("kill finish")
			Convey("the disconnected should be true and id not -1", func() {
				So(TestConnection.Disconnected(), ShouldBeTrue)
				So(TestConnection.User.ID, ShouldEqual, id)
			})
		})
	})
}

func TestAll(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Launch everything", t, func() {

		s := httptest.NewServer(http.HandlerFunc(echo))
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http")

		ws, _, err := websocket.DefaultDialer.Dial(u, nil)

		So(err, ShouldBeNil)

		defer ws.Close()

		time.Sleep(2 * time.Second)
		<-ready

		// connection.go
		TestConnection.InRoom()
		// connection_json.go
		TestConnection.JSON()
		TestConnection.MarshalJSON()
		TestConnection.UnmarshalJSON([]byte(utils.RandomString(16)))
		// connection_mutex.go
		TestConnection.setDone()
		TestConnection.done()
		TestConnection.Disconnected()
		TestConnection.setDisconnected()
		TestConnection.Room()
		TestConnection.RoomID()
		TestConnection.Both()
		TestConnection.Index()
		TestConnection.SetIndex(1)
		TestConnection.setRoom(nil)
		TestConnection.setBoth(true)

	})
}

func echo(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("errrrrrrr", err.Error())
		return
	}
	defer ws.Close()
	id := rand.Intn(10000)
	user := createRandomUser(id)
	lobby := NewLobby(RANDOMSIZE, RANDOMSIZE, nil, true)

	TestConnection = NewConnection(ws, user, lobby)
	settings := config.WebSocketSettings{
		WriteWait:       time.Duration(60) * time.Second,
		PongWait:        time.Duration(10) * time.Second,
		PingPeriod:      time.Duration(9) * time.Second,
		MaxMessageSize:  512,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	fmt.Println("seeeend")
	ready <- struct{}{}
	TestConnection.Launch(settings)
}

func createRandomUser(id int) *models.UserPublicInfo {
	rand.Seed(time.Now().UnixNano())
	dif := rand.Intn(4)
	return &models.UserPublicInfo{
		ID:        id,
		Name:      utils.RandomString(16),
		Email:     utils.RandomString(16),
		PhotoURL:  utils.RandomString(16),
		FileKey:   utils.RandomString(16),
		Difficult: dif,
	}
}
