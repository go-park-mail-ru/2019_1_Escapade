package rerrors

import "errors"

// ErrorWrongStatus return error, returns an error if
// someone tries to perform an action that is valid not with current
// room status
func ErrorWrongStatus() error {
	return errors.New("Action refused. Too late or too early to perform your action")
}

// ErrorCellOutside return error if the specified cel
//l is out of the field
func ErrorCellOutside() error {
	return errors.New("Action refused. Cell not inside field")
}

// ErrorPlayerFinished return an error if the player, being dead,
// tries to perform an action available only to live players
func ErrorPlayerFinished() error {
	return errors.New("Action refused. You died")
}

// ErrorLobbyCantCreateRoom return an error if lobby cant create room
func ErrorLobbyCantCreateRoom() error {
	return errors.New("Cant create room in lobby")
}

// ErrorInvalidRoomSettings returns an error if the player tried to
// create a room with invalid characteristics
func ErrorInvalidRoomSettings() error {
	return errors.New("Invalid roomSettings")
}

// ErrorLobbyDone return an error if lobby is turned off
func ErrorLobbyDone() error {
	return errors.New("Lobby is turned off")
}

// ErrorRoomOrLobbyDone return an error if lobby or room is turned off
func ErrorRoomOrLobbyDone() error {
	return errors.New("Lobby or room is turned off")
}

// ErrorConnectionDone return an error if connection died
func ErrorConnectionDone() error {
	return errors.New("Connection died")
}

// ErrorRoomDone return an error if room deleted
func ErrorRoomDone() error {
	return errors.New("Room deleted")
}
