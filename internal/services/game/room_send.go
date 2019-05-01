package game

// sendToAllInRoom send info to those in room, whose predicate
// returns true
func (room *Room) send(info interface{}, predicate SendPredicate) {
	SendToConnections(info, predicate, room.Players.Connections,
		room.Observers.Get)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *Room) sendPlayers(predicate SendPredicate) {
	get := &RoomGet{
		Players: true,
	}
	send := room.copy(get)
	room.send(send, predicate)
}

func (room *Room) sendMessage(text string, predicate SendPredicate) {
	room.send("Room("+room.Name+"):"+text, predicate)
}

func (room *Room) sendObservers(predicate SendPredicate) {
	get := &RoomGet{
		Observers: true,
	}
	send := room.copy(get)
	room.send(send, predicate)
}

// sendTAIRField send field to all in room
func (room *Room) sendField(predicate SendPredicate) {
	get := &RoomGet{
		Field: true,
	}
	send := room.copy(get)
	room.send(send, predicate)
}

// sendTAIRHistory send actions history to all in room
func (room *Room) sendHistory(predicate SendPredicate) {
	get := &RoomGet{
		History: true,
	}
	send := room.copy(get)
	room.send(send, predicate)
}

// sendTAIRAll send everything to one connection
func (room *Room) greet(conn *Connection) {
	conn.SendInformation(room)
}

// copy returns full slices of selected fields
func (room *Room) copy(get *RoomGet) *Room {
	sendRoom := &Room{
		Name:   room.Name,
		Status: room.Status,
		Date:   room.Date,
		Type:   room.Type,
	}

	if get.Players {
		sendRoom.Players = room.Players
	}
	if get.Observers {
		sendRoom.Observers = room.Observers
	}
	if get.Field {
		sendRoom.Field = room.Field
	}
	if get.History {
		sendRoom.History = room.History
	}
	return sendRoom
}

// copyLast returns last element of slices of selected fields
// func (room *Room) copyLast(get *RoomGet) *Room {
// 	sendRoom := &Room{
// 		Name:   room.Name,
// 		Status: room.Status,
// 	}

// 	if get.Players {
// 		sendRoom.Players = room.Players
// 	}
// 	if get.Observers {
// 		sendRoom.Observers = room.Observers
// 	}
// 	if get.Field {
// 		sendRoom.Field = room.Field
// 	}
// 	if get.History {
// 		sendRoom.History = room.History
// 	}
// 	return sendRoom
// }
