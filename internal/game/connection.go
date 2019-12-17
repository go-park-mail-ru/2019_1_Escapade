package game

import (
	"encoding/json"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"fmt"
	"sync"
	"time"

	"context"

	"github.com/gorilla/websocket"
)

// Connection is a websocket of a player, that belongs to room
type Connection struct {
	wGroup *sync.WaitGroup

	doneM *sync.RWMutex
	_done bool

	playingRoomM *sync.RWMutex
	_playingRoom *Room

	disconnectedM *sync.RWMutex
	_disconnected bool

	waitingRoomM *sync.RWMutex
	_waitingRoom *Room

	indexM *sync.RWMutex
	_index int

	timeM *sync.RWMutex
	_time time.Time

	UUID string
	User *models.UserPublicInfo

	wsM *sync.Mutex
	_ws *websocket.Conn

	lobby *Lobby

	context context.Context
	cancel  context.CancelFunc

	actionSem chan struct{}

	send chan []byte
}

// NewConnection creates a new connection
func NewConnection(ws *websocket.Conn, user *models.UserPublicInfo, lobby *Lobby) *Connection {
	if ws == nil || user == nil || lobby == nil {
		return nil
	}

	context, cancel := context.WithCancel(lobby.context)

	return &Connection{
		wGroup: &sync.WaitGroup{},

		doneM: &sync.RWMutex{},
		_done: false,

		playingRoomM: &sync.RWMutex{},
		_playingRoom: nil,

		disconnectedM: &sync.RWMutex{},
		_disconnected: false,

		waitingRoomM: &sync.RWMutex{},
		_waitingRoom: nil,

		indexM: &sync.RWMutex{},
		_index: -1,

		UUID: utils.RandomString(16),
		User: user,

		wsM: &sync.Mutex{},
		_ws: ws,

		lobby: lobby,

		context: context,
		cancel:  cancel,

		timeM: &sync.RWMutex{},
		_time: time.Now(),

		send:      make(chan []byte),
		actionSem: make(chan struct{}, 1),
	}
}

// new!!!
// Restore
// it calls in lobby restore
func (conn *Connection) Restore(copy *Connection) {
	if conn.done() {
		return
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
	}()

	//fmt.Println("copy info", copy.PlayingRoom(), copy.Both(), copy.Index())
	conn.setPlayingRoom(copy.PlayingRoom())
	conn.setWaitingRoom(copy.WaitingRoom())
	conn.SetIndex(copy.Index())
}

// PushToRoom set field 'room' to real room
func (conn *Connection) PushToRoom(room *Room) {
	if conn.done() {
		return
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
	}()

	conn.setPlayingRoom(room)
	conn.setWaitingRoom(nil)
}

// PushToLobby set field 'room' to nil
func (conn *Connection) PushToLobby() {
	if conn.done() {
		return
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
	}()

	conn.setPlayingRoom(nil)
	conn.setWaitingRoom(nil)
}

// IsConnected check player isnt disconnected
func (conn *Connection) IsConnected() bool {
	if conn.done() {
		return false
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
	}()
	return conn.Disconnected() == false
}

// Free free memory, if flag disconnect true then connection and player will not become nil
func (conn *Connection) Free() {

	if conn.done() {
		return
	}
	conn.setDone()

	conn.wGroup.Wait()

	// dont delete. conn = nil make pointer nil, but other pointers
	// arent nil. If conn.disconnected = true it is mean that all
	// resources are cleared, but pointer alive, so we only make pointer = nil
	if conn.lobby == nil {
		return
	}

	conn.setDisconnected()

	conn.wsClose()
	close(conn.send)
	close(conn.actionSem)
	// dont delete. conn = nil make pointer nil, but other pointers
	// arent nil and we make 'conn.disconnected = true' for them

	conn.lobby = nil
	conn.setPlayingRoom(nil)
	conn.setWaitingRoom(nil)

	//fmt.Println("conn free memory")
}

// InRoom check is player in room
func (conn *Connection) InPlayingRoom() bool {
	if conn.done() {
		return false
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
	}()
	return conn.PlayingRoom() != nil
}

// Launch run the writer and reader goroutines and wait them to free memory
func (conn *Connection) Launch(ws config.WebSocketSettings, roomID string) {

	// dont place there conn.wGroup.Add(1)
	if conn.lobby == nil || conn.lobby.context == nil {
		fmt.Println("lobby nil or hasnt context!")
		return
	}

	all := &sync.WaitGroup{}

	fmt.Println("JoinConn!")
	conn.lobby.JoinConn(conn, 3)
	all.Add(1)
	go conn.WriteConn(conn.context, ws, all)
	all.Add(1)
	go conn.ReadConn(conn.context, ws, all)

	conn.SetConnected()

	//fmt.Println("Wait!")
	if roomID != "" {
		rs := &models.RoomSettings{}
		rs.ID = roomID
		conn.lobby.EnterRoom(conn, rs)
	}
	all.Wait()

	conn.setDisconnected()
	fmt.Println("conn finished")
	conn.lobby.Leave(conn, "finished")
	//conn.Free()
}

// ReadConn connection goroutine to read messages from websockets
func (conn *Connection) ReadConn(parent context.Context, wsc config.WebSocketSettings, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	if conn.done() {
		return
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
		utils.CatchPanic("connection.go WriteConn()")
	}()

	conn.wsInit(wsc)
	for {
		select {
		case <-parent.Done():
			fmt.Println("ReadConn done catched")
			return
		default:
			_, message, err := conn.wsReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					fmt.Println("IsUnexpectedCloseError:" + err.Error())
				} else {
					fmt.Println("expected error:" + err.Error())
				}
				if conn.lobby != nil {
					conn.lobby.Leave(conn, "err.Error()")
				}
				return
			}
			fmt.Println("#", conn.ID(), "read from conn:", string(message))
			conn.SetConnected()
			conn.lobby.chanBroadcast <- &Request{
				Connection: conn,
				Message:    message,
			}
		}
	}
}

// WriteConn connection goroutine to write messages to websockets
// dont put conn.debug here
func (conn *Connection) WriteConn(parent context.Context, wsc config.WebSocketSettings, wg *sync.WaitGroup) {
	defer wg.Done()

	if conn.done() {
		return
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
		utils.CatchPanic("connection.go WriteConn()")
	}()

	ticker := time.NewTicker(wsc.PingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-parent.Done():
			fmt.Println("WriteConn done catched")
			return
		case message, ok := <-conn.send:

			if !ok {
				conn.wsWriteMessage(websocket.CloseMessage, []byte{}, wsc)
				return
			}

			utils.ShowWebsocketMessage(message, conn.ID())

			if err := conn.wsWriteInWriter(message, wsc); err != nil {
				return
			}

		case <-ticker.C:
			if err := conn.wsWriteMessage(websocket.PingMessage, []byte{}, wsc); err != nil {
				return
			}
		}
	}
}

// SendInformation send info
func (conn *Connection) SendInformation(value interface{}) {
	if conn.done() {
		return
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
	}()

	if conn.Disconnected() {
		return
	}

	var (
		bytes []byte
		err   error
	)

	bytes, err = json.Marshal(value)

	if err != nil {
		utils.Debug(true, "cant send information")
	} else {
		conn.send <- bytes
	}
}

// sendGroupInformation send info with WaitGroup
func (conn *Connection) sendGroupInformation(value interface{}, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		utils.CatchPanic("connection.go sendGroupInformation()")
	}()
	conn.SendInformation(value)
}

// ID return players id
func (conn *Connection) ID() int {
	if conn.done() {
		return conn.User.ID
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
	}()
	if conn.User == nil {
		return -1
	}
	return conn.User.ID
}

func sendAccountTaken(conn *Connection) {

	response := models.Response{
		Type: "AccountTaken",
	}
	if conn == nil {
		panic("sendAccountTaken")
	}
	fmt.Println("send sendAccountTaken")
	conn.SendInformation(response)
}
