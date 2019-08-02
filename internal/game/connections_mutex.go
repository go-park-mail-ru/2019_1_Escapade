package game

// RGet return connections slice only for Read!
func (conns *Connections) RGet() []*Connection {

	conns.getM.RLock()
	defer conns.getM.RUnlock()
	return conns._get
}

// Value get the current value of connection slice
func (iter *ConnectionsIterator) Value() *Connection {
	iter.conns.getM.Lock()
	defer iter.conns.getM.Unlock()
	return iter.conns._get[iter.current]
}

// Next increases the 'current' index to move through the slice
func (iter *ConnectionsIterator) Next() bool {
	iter.current++
	iter.conns.getM.Lock()
	defer iter.conns.getM.Unlock()
	if iter.current >= len(iter.conns._get) {
		return false
	}
	return true
}

// Add try add element if its possible. Return bool result
// if element not exists it will be create, otherwise it will change its value
func (conns *Connections) Add(conn *Connection) int {
	i, _ := conns.SearchByID(conn.ID())
	if i >= 0 {
		conns.set(i, conn)
	} else if conns.EnoughPlace() {
		i = conns.len()
		conns.append(conn)
	} else {
		return -1
	}
	return i
}

// Capacity return the capacity of the stored slice
func (conns *Connections) Capacity() int32 {

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
func (conns *Connections) SetCapacity(size int32) {

	conns.capacityM.Lock()
	defer conns.capacityM.Unlock()
	conns._capacity = size
}

// len return length of the slice of connections
func (conns *Connections) len() int {

	conns.getM.Lock()
	defer conns.getM.Unlock()
	return len(conns._get)
}

// remove connection with index 'i' from connection slice
func (conns *Connections) remove(i int) {

	conns.getM.Lock()
	defer conns.getM.Unlock()
	size := len(conns._get)
	if size == 0 {
		return
	}

	conns._get[i], conns._get[size-1] = conns._get[size-1], conns._get[i]
	conns._get[size-1] = nil
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
	return int32(size) < cap
}

// SearchByID find connection by connection ID
// return this connection and its index if success
// otherwise nil and -1
func (conns *Connections) SearchByID(connectionID int32) (index int, connection *Connection) {
	conns.getM.RLock()
	defer conns.getM.RUnlock()
	index = -1
	for i, c := range conns._get {
		if c.ID() == connectionID {
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
	for _, c := range conns._get {
		if c.Index() == index {
			connection = c
			return
		}
	}
	return
}
