package engine

import (
	"sync"

	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
)

// SendPredicate - returns true if the parcel send to that conn
type SendPredicate func(conn *Connection) bool

// SendToConnections send 'info' to everybody,  whose predicate
// returns true
func SendToConnections(info api.JSONtype,
	predicate SendPredicate, groups ...*Connections) {

	waitJobs := &sync.WaitGroup{}
	for _, group := range groups {

		it := NewConnectionsIterator(group)
		for it.Next() {
			member := it.Value()
			if member == nil {
				panic("why nill in send")
			}
			if predicate(member) {
				waitJobs.Add(1)
				go member.sendGroupInformation(info, waitJobs)
			}
		}

		// for _, connection := range group {
		// 	if connection == nil {
		// 		continue
		// 	}
		// 	if predicate(connection) {
		// 		waitJobs.Add(1)
		// 		go connection.sendGroupInformation(info, waitJobs)
		// 	}
		// }
	}
	waitJobs.Wait()
}

// AllExceptThat is SendPredicate to SendToConnections
// it will send everybody except selected one and disconnected
func AllExceptThat(me *Connection) func(*Connection) bool {
	return func(conn *Connection) bool {
		return !conn.done() && !me.done() && conn.ID() != me.ID()
	}
}

// All is SendPredicate to SendToConnections
// it will send everybody, who is connected
func All(conn *Connection) bool {
	return !conn.done()
}

// Me is SendPredicate to SendToConnections
// it will send only to selected connection
func Me(me *Connection) func(*Connection) bool {
	return func(conn *Connection) bool {
		return !conn.done() && !me.done() && conn.ID() == me.ID()
	}
}

// All is SendPredicate to SendToConnections
// it will send everybody in room, who is connected
func (room *RoomSender) All(conn *Connection) bool {
	return !conn.done()
}

// AllExceptThat is SendPredicate to SendToConnections
// it will send everybody in room, except selected one
func (room *RoomSender) AllExceptThat(me *Connection) func(*Connection) bool {
	return func(conn *Connection) bool {
		return !conn.done() && !me.done() && conn != me
	}
}
