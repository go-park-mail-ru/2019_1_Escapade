package engine

import "github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

type RoomConnectionEvents struct {
	r      *Room
	s      SyncI
	notify *RoomRecorder
}

// Init set Room and RoomNotifier pointers
func (room *RoomConnectionEvents) Init(r *Room, s SyncI) {
	room.r = r
	room.s = s
	var notify = &RoomRecorder{}
	notify.Init(r, s)
	room.notify = notify
}

// Timeout handle the situation, when the waiting time for the player
// to return has expired
func (room *RoomConnectionEvents) Timeout(conn *Connection, isPlayer bool) {
	room.s.do(func() {
		if isPlayer {
			if room.r.events.IsActive() {
				room.Kill(conn, ActionTimeout)
			} else {
				room.exitPlayerWhenGameNotRunning(conn)
			}
		} else {
			room.exitObsserver(conn)
		}
		room.r.lobby.greet(conn)
	})
}

// Leave handle player going back to lobby
func (room *RoomConnectionEvents) Leave(conn *Connection) {
	room.s.doWithConn(conn, func() {
		// work in rooms structs
		if room.isPlayer(conn) {
			if room.r.events.IsActive() {
				room.GiveUp(conn)
			} else {
				room.exitPlayerWhenGameNotRunning(conn)
			}
		} else {
			room.exitObsserver(conn)
		}

		// inform lobby
		go room.r.lobby.LeaveRoom(conn, ActionBackToLobby, room.r)

	})
}

// GiveUp kill connection, that call it
func (room *RoomConnectionEvents) GiveUp(conn *Connection) {
	if !room.r.events.IsActive() {
		return
	}
	room.Kill(conn, ActionGiveUp)
}

// Reconnect reconnect connection to room
func (room *RoomConnectionEvents) Reconnect(conn *Connection) {
	room.s.doWithConn(conn, func() {
		found, isPlayer := room.r.people.Search(conn)
		if found == nil {
			return
		}
		room.r.people.add(conn, isPlayer, true)
	})
}

func (room *RoomConnectionEvents) Restart(conn *Connection) {
	room.s.doWithConn(conn, func() {
		if room.r.events.Status() != StatusFinished {
			return
		}
		if err := room.r.events.Restart(conn); err != nil {
			utils.Debug(false, "cant create room for restart", err.Error())
			return
		}
		room.goToNextRoom(conn)
		room.notify.Restart(conn)
	})
}

// Enter handle user joining as player or observer
func (room *RoomConnectionEvents) Enter(conn *Connection) bool {
	var done bool
	room.s.doWithConn(conn, func() {
		if room.r.events.Status() == StatusRecruitment {
			if room.r.people.add(conn, true, false) {
				done = true
			}
		} else if room.r.people.add(conn, false, false) {
			done = true
		}
	})
	return done
}

// Kill make user die and check for finish battle
func (room *RoomConnectionEvents) Kill(conn *Connection, action int32) {
	room.s.doWithConn(conn, func() {
		if !room.r.people.isAlive(conn) {
			return
		}

		room.r.people.SetFinished(conn)
		room.notify.Kill(conn, action, room.r.Settings.Deathmatch)
		room.r.events.tryFinish()
	})
}

// Disconnect when connection has network problems
func (room *RoomConnectionEvents) Disconnect(conn *Connection) {
	room.s.doWithConn(conn, func() {
		// work in rooms structs
		if conn.PlayingRoom() == nil {
			room.Leave(conn)
			return
		}

		found, _ := room.r.people.Search(conn)
		if found == nil {
			return
		}
		found.setDisconnected()
		room.notify.Disconnect(conn)
	})
}

func (room *RoomConnectionEvents) exitPlayerWhenGameNotRunning(conn *Connection) {
	room.s.do(func() {
		room.r.people.Players.Connections.Remove(conn)
		room.notify.Leave(conn, ActionBackToLobby)
		room.r.send.PlayerExit(conn, room.r.AllExceptThat(conn))
		room.r.events.tryClose()
	})
}

func (room *RoomConnectionEvents) exitObsserver(conn *Connection) {
	room.s.do(func() {
		room.r.people.Observers.Remove(conn)
		room.notify.Leave(conn, ActionBackToLobby)
		room.r.send.ObserverExit(conn, room.r.AllExceptThat(conn))
		room.r.events.tryClose()
	})
}

func (room *RoomConnectionEvents) isPlayer(conn *Connection) bool {
	return conn.Index() >= 0
}

func (room *RoomConnectionEvents) goToNextRoom(conn *Connection) {
	room.Leave(conn)
	room.r.events.Next().connEvents.Enter(conn)
}
