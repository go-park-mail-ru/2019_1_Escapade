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

func ErrorLobbyDone() error {
	return errors.New("Lobby is turned off")
}

func ErrorConnectionDone() error {
	return errors.New("Connection died")
}

func ErrorRoomDone() error {
	return errors.New("Room deleted")
}
