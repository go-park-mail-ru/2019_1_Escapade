package engine

import "github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

type ConnectionEventsI interface {
	nit(s SyncI, l RoomLobbyCommunicationI,
		re *RoomRecorder, se *RoomSender, i *RoomInformation, e *RoomEvents,
		p *RoomPeople)

	Timeout(conn *Connection)
	Leave(conn *Connection)
	GiveUp(conn *Connection)
	Reconnect(conn *Connection)
	Restart(conn *Connection)
	Enter(conn *Connection)
	Disconnect(conn *Connection)
}

type RoomConnectionEvents struct {
	s  SyncI
	l  RoomLobbyCommunicationI
	re *RoomRecorder
	se *RoomSender
	i  *RoomInformation
	e  *RoomEvents
	p  *RoomPeople
}

// Init set Room and RoomNotifier pointers
func (room *RoomConnectionEvents) Init(s SyncI, l RoomLobbyCommunicationI,
	re *RoomRecorder, se *RoomSender, i *RoomInformation, e *RoomEvents,
	p *RoomPeople) {
	room.s = s
	room.l = l
	room.re = re
	room.se = se
	room.i = i
	room.e = e
	room.p = p
}

// Timeout handle the situation, when the waiting time for the player
// to return has expired
func (room *RoomConnectionEvents) Timeout(conn *Connection) {
	room.s.doWithConn(conn, func() {
		isPlayer := room.isPlayer(conn)
		if isPlayer {
			if room.e.IsActive() {
				room.Kill(conn, ActionTimeout)
			} else {
				room.exitPlayerWhenGameNotRunning(conn)
			}
		} else {
			room.exitObsserver(conn)
		}
		room.l.Greet(conn)
	})
}

// Leave handle player going back to lobby
func (room *RoomConnectionEvents) Leave(conn *Connection) {
	room.s.doWithConn(conn, func() {
		// work in rooms structs
		if room.isPlayer(conn) {
			if room.e.IsActive() {
				room.GiveUp(conn)
			} else {
				room.exitPlayerWhenGameNotRunning(conn)
			}
		} else {
			room.exitObsserver(conn)
		}

		// inform lobby
		go room.l.BackToLobby(conn)

	})
}

// GiveUp kill connection, that call it
func (room *RoomConnectionEvents) GiveUp(conn *Connection) {
	if !room.e.IsActive() {
		return
	}
	room.Kill(conn, ActionGiveUp)
}

// Reconnect reconnect connection to room
func (room *RoomConnectionEvents) Reconnect(conn *Connection) {
	room.s.doWithConn(conn, func() {
		found, isPlayer := room.p.Search(conn)
		if found == nil {
			return
		}
		room.p.add(conn, isPlayer, true)
	})
}

func (room *RoomConnectionEvents) Restart(conn *Connection) {
	room.s.doWithConn(conn, func() {
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
func (room *RoomConnectionEvents) Enter(conn *Connection) bool {
	var done bool
	room.s.doWithConn(conn, func() {
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
func (room *RoomConnectionEvents) Kill(conn *Connection, action int32) {
	room.s.doWithConn(conn, func() {
		if !room.p.isAlive(conn) {
			return
		}

		room.p.SetFinished(conn)
		room.re.Kill(conn, action, room.i.Settings.Deathmatch)
		room.e.tryFinish()
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

		found, _ := room.p.Search(conn)
		if found == nil {
			return
		}
		found.setDisconnected()
		room.re.Disconnect(conn)
	})
}

func (room *RoomConnectionEvents) exitPlayerWhenGameNotRunning(conn *Connection) {
	room.s.do(func() {
		room.p.Players.Connections.Remove(conn)
		room.re.Leave(conn, ActionBackToLobby)
		room.se.PlayerExit(conn)
		room.e.tryClose()
	})
}

func (room *RoomConnectionEvents) exitObsserver(conn *Connection) {
	room.s.do(func() {
		room.p.Observers.Remove(conn)
		room.re.Leave(conn, ActionBackToLobby)
		room.se.ObserverExit(conn)
		room.e.tryClose()
	})
}

func (room *RoomConnectionEvents) isPlayer(conn *Connection) bool {
	return conn.Index() >= 0
}

func (room *RoomConnectionEvents) goToNextRoom(conn *Connection) {
	room.Leave(conn)
	room.e.Next().connEvents.Enter(conn)
}
