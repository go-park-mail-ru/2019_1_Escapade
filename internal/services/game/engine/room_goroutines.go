package engine

/*
func (room *Room) runHistory(conn *Connection) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("room_handle.go run()")
		room.wGroup.Done()
	}()

	//players := *room.Players
	actions := room.history()
	cells := room.Field.History
	actionsSize := len(actions)
	cellsSize := len(cells)
	actionsI := 0
	cellsI := 0
	actionTime := time.Now()
	cellTime := time.Now()
	if actionsSize > 0 {
		actionTime = actions[0].Time
	}
	if cellsSize > 0 {
		cellTime = cells[0].Time
	}

	// offline, err := room.lobby.createRoom(room.Settings)
	// if err != nil {
	// 	panic("offline doesnt work")
	// }
	// room.Leave(conn, ActionBackToLobby)
	// offline.Enter(conn)
	// offline.MakeObserver(conn, true)

	// offline.StartFlagPlacing()

	ticker := time.NewTicker(time.Millisecond * 10)
	defer func() {
		room.Status = StatusHistory
		ticker.Stop()
	}()

	for {
		select {
		case value := <-ticker.C:
			if (actionsSize + cellsSize) == (actionsI + cellsI) {
				return
			}
			for actionsI < actionsSize && value > time.Since() {

			}

		}
	}
}*/
