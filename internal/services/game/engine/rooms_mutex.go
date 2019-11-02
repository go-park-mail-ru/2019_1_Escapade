package engine

// RGet return connections slice only for Read!
func (rooms *Rooms) RGet() []*Room {

	rooms.getM.RLock()
	defer rooms.getM.RUnlock()
	return rooms._get
}

func (iter *RoomsIterator) Value() *Room {
	iter.rooms.getM.Lock()
	defer iter.rooms.getM.Unlock()
	return iter.rooms._get[iter.current]
}
func (iter *RoomsIterator) Next() bool {
	iter.current++
	iter.rooms.getM.Lock()
	defer iter.rooms.getM.Unlock()
	if iter.current >= len(iter.rooms._get) {
		return false
	}
	return true
}

// Add try add element if its possible. Return bool result
// if element not exists it will be create, otherwise it will change its value
func (rooms *Rooms) Add(room *Room) bool {
	if rooms.EnoughPlace() {
		rooms.append(room)
	} else {
		return false
	}
	return true
}

// Capacity return the capacity of the stored slice
func (rooms *Rooms) Capacity() int32 {

	rooms.capacityM.RLock()
	defer rooms.capacityM.RUnlock()
	return rooms._capacity
}

// Set set slice of connections
func (rooms *Rooms) Set(slice []*Room) {

	rooms.getM.Lock()
	defer rooms.getM.Unlock()
	rooms._get = slice
}

// SetCapacity set the capacity of slice
func (rooms *Rooms) SetCapacity(size int32) {

	rooms.capacityM.Lock()
	defer rooms.capacityM.Unlock()
	rooms._capacity = size
}

func (rooms *Rooms) len() int {

	rooms.getM.Lock()
	defer rooms.getM.Unlock()
	return len(rooms._get)
}

// append new connection to connection slice
func (rooms *Rooms) append(room *Room) {

	rooms.getM.Lock()
	defer rooms.getM.Unlock()

	rooms._get = append(rooms._get, room)

	return
}

// updates the value of an array element with an index i
func (rooms *Rooms) set(i int, room *Room) {

	rooms.getM.Lock()
	defer rooms.getM.Unlock()

	rooms._get[i] = room

	return
}

// Empty check rooms capacity is 0
// it will happen, when finish is over, cause
// when somebody explodes, the capacity decrements
func (rooms *Rooms) Empty() bool {
	rooms.getM.RLock()
	defer rooms.getM.RUnlock()
	return len(rooms._get) == 0
}

// EnoughPlace check that you can add more elements
func (rooms *Rooms) EnoughPlace() bool {

	rooms.getM.RLock()
	size := len(rooms._get)
	rooms.getM.RUnlock()

	rooms.capacityM.RLock()
	cap := rooms._capacity
	rooms.capacityM.RUnlock()
	return int32(size) < cap
}

// Free free memory
func (rooms *Rooms) Free() {
	if rooms == nil {
		return
	}
	rooms.getM.RLock()
	for _, room := range rooms._get {
		room.Free()
	}
	rooms._get = nil
	rooms.getM.RUnlock()

	rooms.capacityM.RLock()
	rooms._capacity = 0
	rooms.capacityM.RUnlock()
}

// SearchByID find connection by connection ID
// return this connection and its index if success
// otherwise nil and -1
func (rooms *Rooms) Search(roomID string) (room *Room) {
	it := NewRoomsIterator(rooms)

	for it.Next() {
		foundRoom := it.Value()
		if foundRoom.ID() == roomID {
			room = foundRoom
			return
		}
	}
	return
}

func (rooms *Rooms) Remove(roomID string) bool {
	rooms.getM.Lock()
	defer rooms.getM.Unlock()

	size := len(rooms._get)
	if size == 0 {
		return false
	}

	index := -1

	for i, foundRoom := range rooms._get {
		if foundRoom.ID() == roomID {
			index = i
			break
		}
	}

	if index == -1 {
		return false
	}

	rooms._get[index], rooms._get[size-1] = rooms._get[size-1], rooms._get[index]
	rooms._get[size-1] = nil
	rooms._get = rooms._get[:size-1]

	return true
}
