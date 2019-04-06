package api

import (
	"fmt"
	"math/rand"
	"time"

	"encoding/json"
	"escapade/internal/models"
	"escapade/internal/services/game"
	"flag"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:3001", "http service address")

func TestExample(t *testing.T) {

	H, _, err := GetHandler(confPath)
	H.Test = true
	if err != nil || H == nil {
		t.Error("TestDeleteUser catched error:", err.Error())
		return
	}

	// Create test server with the echo handler.
	s := httptest.NewServer(http.HandlerFunc(H.GameOnline))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	// Send message to server, read response and check to see if it's what we expect.
	for i := 0; i < 10; i++ {
		if err := ws.WriteMessage(websocket.TextMessage, []byte("hello")); err != nil {
			t.Fatalf("%v", err)
		}
		_, p, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("%v", err)
		}
		if string(p) != "hello" {
			fmt.Println(string(p))
			t.Fatalf("bad message")
		}
	}
}

func TestCreateRoom(t *testing.T) {

	n := 5
	s := launchServer(t, n)
	for i := 0; i < n; i++ {
		defer s[i].Close()
	}

	ws := launchWS(t, s, n)
	for i := 0; i < n; i++ {
		defer ws[i].Close()
	}

	for i := 0; i < n; i++ {
		getRooms(t, ws[i])
	}
	sendLR(t, ws[0])
	getRooms(t, ws[1])
	askAllFromLobby(t, ws[1])
	getLobby(t, ws[1])
	askAllFromLobby(t, ws[2])
	getLobby(t, ws[2])
	askAllFromLobby(t, ws[3])
	getLobby(t, ws[3])
	askAllFromLobby(t, ws[0])
	getLobby(t, ws[0])
	askAllFromLobby(t, ws[1])
	getLobby(t, ws[1])

	// askAllFromLobby(t, ws[4])
	// for {
	// 	if err := getLobby(t, ws[4]); err != nil {
	// 		break
	// 	}
	// }

	time.Sleep(2 * time.Second)
	//t.Fatalf("stop")
}

// send lobby request
func sendLR(t *testing.T, ws *websocket.Conn) {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(3)
	var rs *models.RoomSettings
	switch i {
	case 0:
		rs = models.NewSmallRoom()
	case 1:
		rs = models.NewUsualRoom()
	case 2:
		rs = models.NewBigRoom()
	}

	ls := &game.LobbySend{rs}
	lr := game.NewLobbyRequest(ls, nil)
	bytes, err := json.Marshal(lr)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if err = ws.WriteMessage(websocket.TextMessage, bytes); err != nil {
		t.Fatalf("%v", err)
	}
}

func askAllFromLobby(t *testing.T, ws *websocket.Conn) {

	lg := &game.LobbyGet{true, true, true, true}
	lr := game.NewLobbyRequest(nil, lg)
	bytes, err := json.Marshal(lr)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if err = ws.WriteMessage(websocket.TextMessage, bytes); err != nil {
		t.Fatalf("%v", err)
	}
}

func getLobby(t *testing.T, ws *websocket.Conn) error {

	_, r, err := ws.ReadMessage()
	//if err != nil {
	//	t.Fatalf("%v", err)
	//}
	real := string(r)
	//expected := `{"Capacity":500,"Size":0,"Rooms":{}}`
	fmt.Println("GOT", real)
	// if real != expected {
	// 	t.Fatalf("Expected %v, got %v", expected, real)
	// } else {
	// 	fmt.Println("getRooms done")
	// }
	return err
}

func launchServer(t *testing.T, n int) []*httptest.Server {
	H, _, err := GetHandler(confPath)
	if err != nil {
		t.Fatalf("%v", err)
	}
	H.Test = true
	if err != nil || H == nil {
		t.Fatalf("%v", err)
		return nil
	}

	servers := make([]*httptest.Server, n)
	for i := 0; i < n; i++ {
		servers[i] = httptest.NewServer(http.HandlerFunc(H.GameOnline))
	}

	return servers
}

func launchWS(t *testing.T, s []*httptest.Server, n int) []*websocket.Conn {
	// Convert http://127.0.0.1 to ws://127.0.0.

	// Connect to the server
	ws := make([]*websocket.Conn, n)
	var err error
	for i := 0; i < n; i++ {
		u := "ws" + strings.TrimPrefix(s[i].URL, "http")
		ws[i], _, err = websocket.DefaultDialer.Dial(u, nil)
		if err != nil {
			t.Fatalf("%v", err)
		}
	}
	return ws
}

func getRooms(t *testing.T, ws *websocket.Conn) {
	_, r, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}
	real := string(r)
	expected := `{"capacity":500,"get":[]}`
	if real != expected {
		t.Fatalf("Expected %v, got %v", expected, real)
	} else {
		fmt.Println("getRooms done")
	}
}
