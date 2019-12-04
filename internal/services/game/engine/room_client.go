package engine

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
	action_ "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine/action"
	room_ "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine/room"
)

// RClientI specifies the actions that a room can perform on a connection
// Room Client Interface - strategy pattern
type RClientI interface {
	Timeout(conn *Connection)
	Leave(conn *Connection)
	GiveUp(conn *Connection)
	Reconnect(conn *Connection)
	Restart(conn *Connection)
	Enter(conn *Connection)
	Disconnect(conn *Connection)
	isPlayer(conn *Connection) bool
}

// RClient implements RClientI
type RClient struct {
	synced.PublisherBase

	s  synced.SyncI
	l  LobbyProxyI
	re ActionRecorderI
	se RSendI
	e  EventsI
	p  PeopleI

	isDeathMatch bool
}

type ConnectionMsg struct {
	connection *Connection
	code       int32
	content    interface{}
}

// Init configure dependencies with other components of the room
func (room *RClient) Init(builder RBuilderI, isDeathMatch bool) {
	builder.BuildSync(&room.s)
	builder.BuildLobby(&room.l)
	builder.BuildRecorder(&room.re)
	builder.BuildSender(&room.se)
	builder.BuildEvents(&room.e)
	builder.BuildPeople(&room.p)

	room.isDeathMatch = isDeathMatch

	room.PublisherBase = *synced.NewPublisher()
	room.Start(room.l.ConnectionSub(),
		room.re.ConnectionSub(), room.p.ConnectionSub())
}

func (room *RClient) free() {
	room.Stop()
}

func (room *RClient) notify(conn *Connection, code int32, content interface{}) {
	room.Notify(synced.Msg{
		Code: room_.UpdateConnection,
		Content: ConnectionMsg{
			connection: conn,
			code:       code,
			content:    content,
		},
	})
}

// Timeout handle the situation, when the waiting time for the player
// to return has expired
func (room *RClient) Timeout(conn *Connection) {
	room.leave(conn, action_.Timeout)
}

func (room *RClient) leave(conn *Connection, action int32) {
	room.s.DoWithOther(conn, func() {
		isPlayer := room.isPlayer(conn)
		if isPlayer {
			room.notify(conn, action, nil)
		}
		room.notify(conn, action_.BackToLobby, isPlayer)
	})
}

// Leave handle player going back to lobby
func (room *RClient) Leave(conn *Connection) {
	room.GiveUp(conn)
}

// GiveUp kill connection, that call it
func (room *RClient) GiveUp(conn *Connection) {
	room.notify(conn, action_.GiveUp, nil)
}

// Reconnect connection to room
func (room *RClient) Reconnect(conn *Connection) {
	room.s.DoWithOther(conn, func() {
		found, isPlayer := room.p.Search(conn)
		if found == nil {
			return
		}
		room.notify(conn, action_.Reconnect, isPlayer)
	})
}

// Restart marks the connection as wanting to restart and informs
// 	the room of this intention
func (room *RClient) Restart(conn *Connection) {
	room.s.DoWithOther(conn, func() {
		if room.e.Status() != room_.StatusFinished {
			return
		}
		if err := room.e.Restart(conn); err != nil {
			utils.Debug(false, "cant create room for restart", err.Error())
			return
		}
		room.goToNextRoom(conn)
		room.notify(conn, action_.Restart, nil)
	})
}

// Enter handle user joining as player or observer
func (room *RClient) Enter(conn *Connection) {
	room.s.DoWithOther(conn, func() {
		if room.e.Status() == room_.StatusRecruitment {
			room.notify(conn, action_.ConnectAsPlayer, nil)
		} else {
			room.notify(conn, action_.ConnectAsObserver, nil)
		}
	})
}

// Disconnect when connection has network problems
func (room *RClient) Disconnect(conn *Connection) {
	room.s.DoWithOther(conn, func() {
		// work in rooms structs
		if conn.PlayingRoom() == nil {
			room.Leave(conn)
			return
		}

		found, _ := room.p.Search(conn)
		if found == nil {
			return
		}
		found.setDisconnected()
		room.re.Disconnect(conn)
	})
}

func (room *RClient) isPlayer(conn *Connection) bool {
	return conn.Index() >= 0
}

func (room *RClient) goToNextRoom(conn *Connection) {
	room.Leave(conn)
	room.e.Next().client.Enter(conn)
}
