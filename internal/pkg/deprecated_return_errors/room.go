package rerrors

// ErrorRoomIsFull room is full
func ErrorRoomIsFull() error {
	return New("Room is full")
}
