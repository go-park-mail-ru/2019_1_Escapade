package rerrors

import "errors"

func ErrorBattleAlreadyBegan() error {
	return errors.New("Action refused. Game also began")
}

func ErrorCellOutside() error {
	return errors.New("Action refused. Cell not inside field")
}

func ErrorPlayerFinished() error {
	return errors.New("Action refused. You died")
}

func ErrorLobbyCantCreateRoom() error {
	return errors.New("Cant create room in lobby")
}

func ErrorInvalidRoomSettings() error {
	return errors.New("Invalid roomSettings")
}
