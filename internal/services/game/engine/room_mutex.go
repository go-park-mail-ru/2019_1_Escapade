package engine

// Name return the name of room given by its creator
func (room *Room) Name() string {
	room.nameM.RLock()
	v := room._name
	room.nameM.RUnlock()
	return v
}

// ID return room's unique identificator
func (room *Room) ID() string {
	room.idM.RLock()
	v := room._id
	room.idM.RUnlock()
	return v
}

func (room *Room) setName(name string) {
	room.nameM.Lock()
	room._name = name
	room.nameM.Unlock()
}

func (room *Room) setID(id string) {
	room.idM.Lock()
	room._id = id
	room.idM.Unlock()
}
