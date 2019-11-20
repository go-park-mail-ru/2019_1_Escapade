package engine

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

// RClientI specifies the actions that a room can perform on a connection
// Room Client Interface - strategy pattern
type RClientI interface {
	Timeout(conn *Connection)
	Leave(conn *Connection)
	GiveUp(conn *Connection)
	Reconnect(conn *Connection)
	Restart(conn *Connection)
	Enter(conn *Connection) bool
	Disconnect(conn *Connection)
	isPlayer(conn *Connection) bool

	Kill(conn *Connection, action int32)
}

// RClient implements RClientI
type RClient struct {
	s  synced.SyncI
	l  LobbyProxyI
	re ActionRecorderI
	se RSendI
	e  EventsI
	p  PeopleI

	isDeathMatch bool
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
}

// Timeout handle the situation, when the waiting time for the player
// to return has expired
func (room *RClient) Timeout(conn *Connection) {
	room.leave(conn, ActionTimeout)
}

func (room *RClient) leave(conn *Connection, action int32) {
	room.s.DoWithOther(conn, func() {
		isPlayer := room.isPlayer(conn)
		if isPlayer {
			if room.e.IsActive() {
				room.Kill(conn, action)
			}
		}
		room.p.Remove(conn)
		room.re.Leave(conn, ActionBackToLobby, isPlayer)
		room.e.tryClose()
		room.l.Greet(conn)
	})
}

// Leave handle player going back to lobby
func (room *RClient) Leave(conn *Connection) {
	room.leave(conn, ActionGiveUp)
}

// GiveUp kill connection, that call it
func (room *RClient) GiveUp(conn *Connection) {
	if !room.e.IsActive() {
		return
	}
	room.Kill(conn, ActionGiveUp)
}

// Reconnect connection to room
func (room *RClient) Reconnect(conn *Connection) {
	room.s.DoWithOther(conn, func() {
		found, isPlayer := room.p.Search(conn)
		if found == nil {
			return
		}
		room.p.add(conn, isPlayer, true)
	})
}

// Restart marks the connection as wanting to restart and informs
// 	the room of this intention
func (room *RClient) Restart(conn *Connection) {
	room.s.DoWithOther(conn, func() {
		if room.e.Status() != StatusFinished {
			return
		}
		if err := room.e.Restart(conn); err != nil {
			utils.Debug(false, "cant create room for restart", err.Error())
			return
		}
		room.goToNextRoom(conn)
		room.re.Restart(conn)
	})
}

// Enter handle user joining as player or observer
func (room *RClient) Enter(conn *Connection) bool {
	var done bool
	room.s.DoWithOther(conn, func() {
		if room.e.Status() == StatusRecruitment {
			if room.p.add(conn, true, false) {
				done = true
			}
		} else if room.p.add(conn, false, false) {
			done = true
		}
	})
	return done
}

// Kill make user die and check for finish battle
func (room *RClient) Kill(conn *Connection, action int32) {
	room.s.DoWithOther(conn, func() {
		if !room.p.isAlive(conn) {
			return
		}

		room.p.SetFinished(conn)
		room.re.Kill(conn, action, room.isDeathMatch)
		room.e.tryFinish()
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
