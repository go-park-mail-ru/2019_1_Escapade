package game

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/metrics"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
	"github.com/gorilla/websocket"
	. "github.com/smartystreets/goconvey/convey"
)

// go test -coverprofile cover.out
// go tool cover -html=cover.out -o test/coverage.html

var n = 10
var connections = make([]*Connection, n)

var wss = config.WebSocketSettings{
	WriteWait:       time.Duration(60) * time.Second,
	PongWait:        time.Duration(10) * time.Second,
	PingPeriod:      time.Duration(9) * time.Second,
	MaxMessageSize:  512,
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var gc = &config.GameConfig{
	ConnectionCapacity: 500,
	RoomsCapacity:      500,
	CanClose:           true,
}

func TestCreateRoom(t *testing.T) {

	metrics.InitRoomMetric("game")
	metrics.InitPlayersMetric("game")
	// Only pass t into top-level Convey calls
	Convey("Given slice of connections", t, func() {

		for i := 0; i < 5; i++ {
			s := httptest.NewServer(http.HandlerFunc(testGame(i)))
			defer s.Close()

			u := "ws" + strings.TrimPrefix(s.URL, "http")

			ws, _, err := websocket.DefaultDialer.Dial(u, nil)
			if err != nil {
				Convey("When websocket dials, the error should be nil", func() {
					So(err, ShouldBeNil)
				})

				return
			}
			defer ws.Close()
			<-ready
		}

		time.Sleep(1 * time.Second)

		Convey("When create room", func() {
			rsCreate := `{"send":{"RoomSettings":{"name":"my best room","id":"create","width":12,"height":12,"players":5,"observers":10,"prepare":10, "play":100, "mines":5}},"get":null}`
			request := &Request{
				Connection: connections[0],
				Message:    []byte(rsCreate),
			}
			connections[0].lobby.chanBroadcast <- request
			time.Sleep(1 * time.Second)
			rsConnect := `
			{
				"send":
				{
					"RoomSettings":
					{
						"id":""
					}
				}
			}
			`
			request = &Request{
				Connection: connections[1],
				Message:    []byte(rsConnect),
			}
			connections[1].lobby.chanBroadcast <- request
			time.Sleep(1 * time.Second)
			request = &Request{
				Connection: connections[2],
				Message:    []byte(rsConnect),
			}
			connections[2].lobby.chanBroadcast <- request
			time.Sleep(1 * time.Second)
			request = &Request{
				Connection: connections[3],
				Message:    []byte(rsConnect),
			}
			connections[3].lobby.chanBroadcast <- request
			time.Sleep(1 * time.Second)
			actionBackToLobby := `
			{
				"send":
				{
					"action": 14
				}
			}
			`

			request = &Request{
				Connection: connections[3],
				Message:    []byte(actionBackToLobby),
			}

			connections[3].lobby.chanBroadcast <- request

			// rooms := connections[3].lobby._AllRooms.Get
			// if len(rooms) == 0 {
			// 	panic("paaaaaaaaaaanic")
			// }
			// id := rooms[0].ID
			// rsToRoom := `{"send":{"RoomSettings":{"id":"` + id + `"}}}`
			// request = &Request{
			// 	Connection: connections[3],
			// 	Message:    []byte(rsToRoom),
			// }
			// connections[3].lobby.chanBroadcast <- request

			//connections[9].cancel()

			//connections[0].lobby.chanLeave <- connections[0]
			//connections[0].lobby.chanJoin <- connections[0]
			//	time.Sleep(5 * time.Second)

			roomID := utils.RandomString(16)
			settings := models.NewSmallRoom()
			_, err := NewRoom(settings, roomID, GetLobby())
			if err != nil {
				panic(111111)
			}

			//connections[0].lobby.sendLobbyMessage("no", All)
			// connections[0].lobby.sendRoomCreate(*room, All)
			// connections[0].lobby.sendRoomUpdate(*room, All)
			// connections[0].lobby.sendRoomDelete(*room, All)
			// connections[0].lobby.sendWaiterEnter(*connections[0], All)
			// connections[0].lobby.sendPlayerEnter(*connections[0], All)
			// connections[0].lobby.sendPlayerExit(*connections[0], All)
			//lobby := GetLobby()
			//lobby.CreateAndAddToRoom(settings, connections[4])
			//lobby.RoomStart(room)
			//room.FinishGame(false)

		})
	})
}

func TestLobby(t *testing.T) {

	metrics.InitRoomMetric("game")
	metrics.InitPlayersMetric("game")
	// Only pass t into top-level Convey calls
	Convey("Given slice of connections", t, func() {

		s := httptest.NewServer(http.HandlerFunc(testGame(7)))
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http")

		ws, _, err := websocket.DefaultDialer.Dial(u, nil)
		if err != nil {
			Convey("When websocket dials, the error should be nil", func() {
				So(err, ShouldBeNil)
			})

			return
		}
		defer ws.Close()

		roomID := utils.RandomString(16)
		settings := models.NewSmallRoom()
		room, err := NewRoom(settings, roomID, GetLobby())
		if err != nil {
			panic(111111)
		}

		time.Sleep(3 * time.Second)

		l := connections[7].lobby
		l.sendLobbyMessage("no", All)
		l.sendRoomCreate(*room, All)
		l.sendRoomUpdate(*room, All)
		l.sendRoomDelete(*room, All)
		l.sendWaiterEnter(*connections[7], All)
		l.sendPlayerEnter(*connections[7], All)
		l.sendPlayerExit(*connections[7], All)
		l.sendWaiterExit(*connections[7], All)

		//l.CreateAndAddToRoom(settings, connections[7])
		//l.RoomStart(room)
		//room.FinishGame(false)

	})
}

func CatchPanic(place string) {
	// panic doesnt recover

	if r := recover(); r != nil {
		fmt.Println("Panic recovered in", place)
		fmt.Println("More", r)
	}
}

func TestRoom(t *testing.T) {

	metrics.InitRoomMetric("game")
	metrics.InitPlayersMetric("game")
	defer CatchPanic("TestRoom")
	// Only pass t into top-level Convey calls
	Convey("Given slice of connections", t, func() {

		s := httptest.NewServer(http.HandlerFunc(testGame(7)))
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http")

		ws, _, err := websocket.DefaultDialer.Dial(u, nil)
		if err != nil {
			Convey("When websocket dials, the error should be nil", func() {
				So(err, ShouldBeNil)
			})

			return
		}
		defer ws.Close()

		roomID := utils.RandomString(16)
		settings := models.NewSmallRoom()
		room, err := NewRoom(settings, roomID, GetLobby())
		if err != nil {
			panic(111111)
		}

		time.Sleep(3 * time.Second)
		room.playersCapacity()
		room.player(0)
		room.players()
		room.sendPlayerEnter(*connections[7], All)
		room.playerFlag(0)
		room.sendField(All)
		room.sendObserverEnter(*connections[7], All)
		room.sendObserverExit(*connections[7], All)
		room.sendPlayerExit(*connections[7], All)
		room.playerFinished(0)
		room.sendStatus(All)

		room.RecoverPlayer(connections[7])
		room.RecoverObserver(connections[6], connections[7])
		room.RemoveFromGame(connections[7], true)
		room.SetFinished(connections[7])

		//l.CreateAndAddToRoom(settings, connections[7])
		//l.RoomStart(room)
		//room.FinishGame(false)

	})
}

func testGame(i int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("errrrrrrr", err.Error())
			return
		}
		defer ws.Close()
		Launch(gc, nil)

		//TestConnection = NewConnection(ws, user, lobby)

		fmt.Println("seeeend")

		//all := &sync.WaitGroup{}

		fmt.Println("JoinConn!")

		//var users = make([]*models.UserPublicInfo, n)
		//for i := 0; i < n; i++ {
		//users[i] = createRandomUser(i + 1)
		user := createRandomUser(i + 1)
		connections[i] = NewConnection(ws, user, GetLobby())
		ready <- struct{}{}
		connections[i].Launch(wss)
		//all.Add(1)
		//go ConnectionLaunch(connections[i], settings, all)
		//}
	}
}

func testLobby(i int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("errrrrrrr", err.Error())
			return
		}
		defer ws.Close()

		tl := NewLobby(500, 500, nil, true)
		user := createRandomUser(i + 1)
		connections[i] = NewConnection(ws, user, tl)
		connections[i].Launch(wss)
		//all.Add(1)
		//go ConnectionLaunch(connections[i], settings, all)
		//}
	}
}

func ConnectionLaunch(conn *Connection, wss config.WebSocketSettings,
	wg *sync.WaitGroup) {
	defer wg.Done()
	conn.Launch(wss)
}
