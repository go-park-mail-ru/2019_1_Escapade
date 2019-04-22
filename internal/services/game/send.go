package game

import (
	"encoding/json"
	"sync"
)

// SendPredicate - returns true if the parcel send to that conn
type SendPredicate func(conn *Connection) bool

// SendToConnections send 'info' to everybody,  whose predicate
// returns true
func SendToConnections(info interface{},
	predicate SendPredicate, groups ...[]*Connection) {
	waitJobs := &sync.WaitGroup{}
	bytes, _ := json.Marshal(info)
	for _, group := range groups {
		for _, connection := range group {
			if predicate(connection) {
				waitJobs.Add(1)
				go connection.sendGroupInformation(bytes, waitJobs)
			}
		}
	}
	waitJobs.Wait()
}

// allExceptThat is predicate to sendToAllInRoom
// it will send everybody except selected one and disconnected
func AllExceptThat(me *Connection) func(*Connection) bool {
	return func(conn *Connection) bool {
		return conn != me && conn.disconnected == false
	}
}

// allExceptThat is predicate to sendToAllInRoom
// it will send everybody except selected one and disconnected
func All() func(*Connection) bool {
	return func(conn *Connection) bool {
		return conn.IsConnected()
	}
}

// all is predicate to sendToAllInRoom
// it will send everybody except disconnected
func (room *Room) InThisRoom() func(conn *Connection) bool {
	return func(conn *Connection) bool {
		return conn.room == room
	}
}
