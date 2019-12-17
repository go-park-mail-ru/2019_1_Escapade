package game

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// EnterPlayer handle player try to enter room
func (room *Room) addConnection(conn *Connection, isPlayer bool, needRecover bool) bool {
	//fmt.Println("addPlayer", recover)
	if room.done() {
		return false
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	// later return back!!!!!
	// if room.lobby.Metrics() {
	// 	metrics.Players.WithLabelValues(room.ID, conn.User.Name).Inc()
	// }

	utils.Debug(false, "Room("+room.ID()+") wanna connect you mr ", conn.ID())

	// if room hasnt got places
	if !room.Push(conn, isPlayer, needRecover) {
		return false
	}

	var pa *PlayerAction
	if needRecover {
		pa = room.addAction(conn.ID(), ActionReconnect)
	} else if isPlayer {
		if !room.Players.EnoughPlace() {
			return false
		}
		pa = room.addAction(conn.ID(), ActionConnectAsPlayer)
		// maybe delete it?
		//go room.sendPlayerEnter(*conn, room.AllExceptThat(conn))
	} else {
		if !room.Observers.EnoughPlace() {
			return false
		}
		pa = room.addAction(conn.ID(), ActionConnectAsObserver)
		// maybe delete it?
		go room.sendObserverEnter(conn, room.AllExceptThat(conn))
	}
	go room.sendAction(*pa, room.AllExceptThat(conn))

	if !needRecover {
		room.lobby.sendRoomUpdate(room, All)
	}
	room.sendStatusOne(conn)

	return true
}

func (room *Room) Push(conn *Connection, isPlayer bool, needRecover bool) bool {
	if room.done() {
		return false
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	if isPlayer {
		if !needRecover && !room.Players.EnoughPlace() {
			return false
		}
		room.Players.Add(conn, room.Field.CreateRandomFlag(conn.ID()), false, needRecover)
		if !needRecover && !room.Players.EnoughPlace() {
			room.chanStatus <- StatusFlagPlacing
		}
	} else {
		if !needRecover && !room.Observers.EnoughPlace() {
			return false
		}
		room.Observers.Add(conn)
	}

	room.greet(conn, isPlayer)
	if room.Status() != StatusPeopleFinding {
		room.lobby.waiterToPlayer(conn, room)
	} else {
		conn.setWaitingRoom(room)
	}

	return true

}

// Search search connection in players and observers of room
// return connection and flag isPlayer
func (room *Room) Search(find *Connection) (*Connection, bool) {
	i, found := room.Players.SearchConnection(find)
	if i >= 0 {
		return found, true
	}
	i, found = room.Observers.SearchByID(find.ID())
	if i >= 0 {
		return found, false
	}
	return nil, true
}
