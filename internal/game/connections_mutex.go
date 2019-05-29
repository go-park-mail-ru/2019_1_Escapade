package game

// RGet return connections slice only for Read!
func (conns *Connections) RGet() []*Connection {

	conns.getM.RLock()
	defer conns.getM.RUnlock()
	return conns._get
}

// Capacity return the capacity of the stored slice
func (conns *Connections) Capacity() int {

	conns.capacityM.RLock()
	defer conns.capacityM.RUnlock()
	return conns._capacity
}

// Set set slice of connections
func (conns *Connections) Set(slice []*Connection) {

	conns.getM.Lock()
	defer conns.getM.Unlock()
	conns._get = slice
}

// SetCapacity set the capacity of slice
func (conns *Connections) SetCapacity(size int) {

	conns.capacityM.Lock()
	defer conns.capacityM.Unlock()
	conns._capacity = size
}

// remove connection with index 'i' from connection slice
func (conns *Connections) remove(i int) {

	conns.getM.Lock()
	defer conns.getM.Unlock()
	size := len(conns._get)

	conns._get[i], conns._get[size-1] = conns._get[size-1], conns._get[i]
	conns._get[size-1] = nil // now not pointer
	conns._get = conns._get[:size-1]
	return
}

// append new connection to connection slice
func (conns *Connections) append(conn *Connection) {

	conns.getM.Lock()
	defer conns.getM.Unlock()

	conns._get = append(conns._get, conn)

	return
}

// updates the value of an array element with an index i
func (conns *Connections) set(i int, conn *Connection) {

	conns.getM.Lock()
	defer conns.getM.Unlock()

	conns._get[i] = conn

	return
}

// Empty check rooms capacity is 0
// it will happen, when finish is over, cause
// when somebody explodes, the capacity decrements
func (conns *Connections) Empty() bool {
	conns.getM.RLock()
	defer conns.getM.RUnlock()
	return len(conns._get) == 0
}

// EnoughPlace check that you can add more elements
func (conns *Connections) EnoughPlace() bool {

	conns.getM.RLock()
	size := len(conns._get)
	conns.getM.RUnlock()

	conns.capacityM.RLock()
	cap := conns._capacity
	conns.capacityM.RUnlock()
	return size < cap
}

// SearchByID find connection by connection ID
// return this connection and its index if success
// otherwise nil and -1
func (conns *Connections) SearchByID(connectionID int) (connection *Connection, index int) {
	conns.getM.RLock()
	defer conns.getM.RUnlock()
	index = -1
	for i, c := range conns._get {
		if conns._get[i].ID() == connectionID {
			connection = c
			index = i
			return
		}
	}
	return
}

// SearchByIndex find connection by its Index
// return this connection and its index if success
// otherwise nil and -1
func (conns Connections) SearchByIndex(index int) (connection *Connection) {
	conns.getM.RLock()
	defer conns.getM.RUnlock()
	for i, c := range conns._get {
		if conns._get[i].Index() == index {
			connection = c
			return
		}
	}
	return
}
