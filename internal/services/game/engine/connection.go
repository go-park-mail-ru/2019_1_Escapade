package engine

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

// Connection is a websocket of a player, that belongs to room
type Connection struct {
	s synced.SyncI

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
	_ws WebsocketConnI

	//lobby LobbyConnectionI //*Lobby
	Events *ConnEvents

	context context.Context
	cancel  context.CancelFunc

	actionSem chan struct{}

	send chan []byte
}

func newConnection() *Connection {
	var s = &synced.SyncWgroup{}
	s.Init(0)
	context, cancel := context.WithCancel(context.Background())

	return &Connection{
		s: s,

		playingRoomM: &sync.RWMutex{},
		_playingRoom: nil,

		disconnectedM: &sync.RWMutex{},
		_disconnected: false,

		waitingRoomM: &sync.RWMutex{},
		_waitingRoom: nil,

		indexM: &sync.RWMutex{},
		_index: -1,

		UUID: utils.RandomString(16),

		wsM: &sync.Mutex{},

		timeM: &sync.RWMutex{},
		_time: time.Now(),

		Events: NewConnEvent(),

		context: context,
		cancel:  cancel,

		send:      make(chan []byte),
		actionSem: make(chan struct{}, 1),
	}
}

// NewConnection creates a new connection
func NewConnection(ws WebsocketConnI, user *models.UserPublicInfo) (*Connection, error) {
	if ws == nil || user == nil {
		return nil, re.NoWebSocketOrUser()
	}
	conn := newConnection()
	conn.User = user
	conn.setWs(ws)

	return conn, nil
}

// Restore set restored playing and waiting rooms, conn's index
// in Players slice
// It calls in lobby restore
func (conn *Connection) Restore(copy *Connection) {
	conn.s.Do(func() {
		conn.setPlayingRoom(copy.PlayingRoom())
		conn.setWaitingRoom(copy.WaitingRoom())
		conn.SetIndex(copy.Index())
	})
}

// IsAnonymous return true if user not registered
func (conn *Connection) IsAnonymous() bool {
	return conn.ID() < 0
}

// PushToRoom set field 'room' to real room
func (conn *Connection) PushToRoom(room *Room) {
	conn.s.Do(func() {
		conn.setPlayingRoom(room)
		conn.setWaitingRoom(nil)
	})
}

// PushToLobby set field 'room' to nil
func (conn *Connection) PushToLobby() {
	conn.s.Do(func() {
		conn.setPlayingRoom(nil)
		conn.setWaitingRoom(nil)
	})
}

// IsConnected check player isnt disconnected
func (conn *Connection) IsConnected() bool {
	return conn.Disconnected() == false
}

// Free free memory, if flag disconnect true then connection and player will not become nil
func (conn *Connection) Free() {

	conn.s.Clear(func() {
		conn.setDisconnected()

		conn.wsClose()
		conn.Events.Close()
		close(conn.send)
		close(conn.actionSem)
		// dont delete. conn = nil make pointer nil, but other pointers
		// arent nil and we make 'conn.disconnected = true' for them

		conn.setPlayingRoom(nil)
		conn.setWaitingRoom(nil)
	})
}

// InPlayingRoom check is player in playing room
func (conn *Connection) InPlayingRoom() bool {
	return conn.PlayingRoom() != nil
}

// Launch run the writer and reader goroutines and wait them to free memory
func (conn *Connection) Launch(cw config.WebSocket, roomID string) {
	conn.s.Do(func() {
		conn.Events.Join()

		all := &sync.WaitGroup{}

		all.Add(1)
		go conn.WriteConn(conn.context, cw, all)
		all.Add(1)
		go conn.ReadConn(conn.context, cw, all)

		conn.SetConnected()

		if roomID != "" {
			conn.Events.EnterRoom(roomID)
		}
		all.Wait()

		conn.setDisconnected()
		if conn == nil {
			utils.Debug(true, "conn nil")
		}
		conn.Events.Leave()
	})
	//conn.Free()
}

// ReadConn connection goroutine to read messages from websockets
func (conn *Connection) ReadConn(parent context.Context, wsc config.WebSocket, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	conn.s.Do(func() {
		conn.setPlayingRoom(nil)
		conn.setWaitingRoom(nil)

		conn.wsInit(wsc)
		for {
			select {
			case <-parent.Done():
				utils.Debug(false, "ReadConn done catched")
				return
			default:
				_, message, err := conn.wsReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
						utils.Debug(false, "IsUnexpectedCloseError:"+err.Error())
					} else {
						utils.Debug(false, "expected error:"+err.Error())
					}
					return
				}
				utils.Debug(false, "#", conn.ID(), "read from conn:", string(message))
				conn.SetConnected()
				conn.Events.Request(message)
			}
		}
	})
}

// WriteConn connection goroutine to write messages to websockets
// dont put conn.debug here
func (conn *Connection) WriteConn(parent context.Context, wsc config.WebSocket, wg *sync.WaitGroup) {
	defer wg.Done()

	conn.s.Do(func() {
		ping := time.NewTicker(wsc.PingPeriod.Duration)
		defer ping.Stop()

		for {
			select {
			case <-parent.Done():
				utils.Debug(false, "WriteConn done catched")
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
			case <-ping.C:
				if err := conn.wsWriteMessage(websocket.PingMessage, []byte{}, wsc); err != nil {
					return
				}
			}
		}
	})
}

// SendInformation send info
func (conn *Connection) SendInformation(value handlers.JSONtype) {
	conn.s.Do(func() {
		if conn.Disconnected() {
			return
		}
		fmt.Println("sended value:", value)
		if bytes, err := value.MarshalJSON(); err != nil {
			utils.Debug(true, "cant send information:", err.Error())
		} else {
			conn.send <- bytes
		}
	})
}

// sendGroupInformation send info with WaitGroup
func (conn *Connection) sendGroupInformation(value handlers.JSONtype, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		utils.CatchPanic("connection.go sendGroupInformation()")
	}()
	conn.SendInformation(value)
}

// ID return player's id
func (conn *Connection) ID() int32 {
	var userID int32
	conn.s.Do(func() {
		userID = conn.User.ID
	})
	return userID
}

// sendAccountTaken send the message 'AccountTaken' to the connection
func sendAccountTaken(conn *Connection) {
	conn.s.Do(func() {
		response := models.Response{
			Type: "AccountTaken",
		}
		conn.SendInformation(&response)
	})
}

func (conn *Connection) GetSync() synced.SyncI {
	return conn.s
}

// 363 -> 305
