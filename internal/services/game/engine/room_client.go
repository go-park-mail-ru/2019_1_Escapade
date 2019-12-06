package engine

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	action_ "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine/action"
	room_ "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine/room"
)

// RClientI specifies the actions that a room can perform on a connection
// Room Client Interface - strategy pattern
type RClientI interface {
	synced.PublisherI

	Timeout(conn *Connection)
	BackToLobby(conn *Connection)
	GiveUp(conn *Connection)
	Reconnect(conn *Connection)
	Restart(conn *Connection)
	//Enter(conn *Connection)
	Disconnect(conn *Connection)
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

// ConnectionMsg the message, that RClient send to observers in
//   Extra field of sync.Msh
type ConnectionMsg struct {
	connection *Connection
	content    interface{}
}

// Init configure dependencies with other components of the room
func (room *RClient) Init(builder RBuilderI, isDeathMatch bool) {
	room.init(isDeathMatch)
	room.build(builder)
	room.subscribe()
}

// Timeout handle the situation, when the waiting time for the player
// to return has expired
func (room *RClient) Timeout(conn *Connection) {
	room.notify(conn, action_.Timeout, conn.IsPlayer())
}

// BackToLobby handle player going back to lobby
func (room *RClient) BackToLobby(conn *Connection) {
	room.notify(conn, action_.BackToLobby, nil)
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
		if err := room.e.Restart(conn); err != nil {
			return
		}
		room.notify(conn, action_.Restart, nil)
	})
}

// Enter handle user joining as player or observer
// func (room *RClient) Enter(conn *Connection) {
// 	room.s.DoWithOther(conn, func() {
// 		if room.e.Status() == room_.StatusRecruitment {
// 			room.notify(conn, action_.ConnectAsPlayer, nil)
// 		} else {
// 			room.notify(conn, action_.ConnectAsObserver, nil)
// 		}
// 	})
// }

// Disconnect when connection has network problems
func (room *RClient) Disconnect(conn *Connection) {
	room.notify(conn, action_.Disconnect, nil)
}

// init struct's values
func (room *RClient) init(isDeathMatch bool) {
	room.PublisherBase = *synced.NewPublisher()
	room.isDeathMatch = isDeathMatch
}

// build components
func (room *RClient) build(builder RBuilderI) {
	builder.BuildSync(&room.s)
	builder.BuildLobby(&room.l)
	builder.BuildRecorder(&room.re)
	builder.BuildSender(&room.se)
	builder.BuildEvents(&room.e)
	builder.BuildPeople(&room.p)
}

// subscribe to room events
func (room *RClient) subscribe() {
	room.e.SubscribeRunnable(room)
}

// notify subscribers of the connection event
func (room *RClient) notify(conn *Connection, code int32, content interface{}) {
	room.Notify(synced.Msg{
		Publisher: room_.UpdateConnection,
		Action:    code,
		Extra: ConnectionMsg{
			connection: conn,
			content:    content,
		},
	})
}

// start RClient goroutines
//   implements RunnableI interface
func (room *RClient) start() {
	room.StartPublish()
}

// stop RClient goroutines
//   implements RunnableI interface
func (room *RClient) stop() {
	room.StopPublish()
}
