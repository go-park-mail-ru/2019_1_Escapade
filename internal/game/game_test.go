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
	"github.com/gorilla/websocket"
	. "github.com/smartystreets/goconvey/convey"
)

// go test -coverprofile test/cover.out
// go tool cover -html=test/cover.out -o test/coverage.html

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

	// Only pass t into top-level Convey calls
	Convey("Given slice of connections", t, func() {

		for i := 0; i < n; i++ {
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

		time.Sleep(2 * time.Second)

		Convey("When create room", func() {
			rsCreate := `
			{
				"send":
				{
					"RoomSettings":
					{
						"name":"my best room",
						"id":"create",
						"width":12,
						"height":12,
						"players":5,
						"observers":10,
						"prepare":10, 
						"play":100, 
						"mines":5
					}
				}
			}
			`
			request := &Request{
				Connection: connections[0],
				Message:    []byte(rsCreate),
			}
			connections[0].lobby.chanBroadcast <- request

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

			request = &Request{
				Connection: connections[2],
				Message:    []byte(rsConnect),
			}
			connections[2].lobby.chanBroadcast <- request

			request = &Request{
				Connection: connections[3],
				Message:    []byte(rsConnect),
			}
			connections[3].lobby.chanBroadcast <- request

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
			//connections[3].lobby.chanBroadcast <- request

			//connections[9].cancel()
		})
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
		connections[i].Launch(wss)
		//all.Add(1)
		//go ConnectionLaunch(connections[i], settings, all)
		//}
		ready <- struct{}{}
	}
}

func ConnectionLaunch(conn *Connection, wss config.WebSocketSettings,
	wg *sync.WaitGroup) {
	defer wg.Done()
	conn.Launch(wss)
}
