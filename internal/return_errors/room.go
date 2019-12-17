package rerrors

import "errors"

// ErrorRoomIsFull room is full
func ErrorRoomIsFull() error {
	return errors.New("Room is full")
}
