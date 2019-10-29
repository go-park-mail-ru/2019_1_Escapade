package game

// LeaveMeta update metainformation about user leaving room
func (room *Room) LeaveMeta(conn *Connection, action int32) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	pa := *room.addAction(conn.ID(), action)
	if !room.Empty() {
		room.sendAction(pa, room.AllExceptThat(conn))
	}

	return
}
