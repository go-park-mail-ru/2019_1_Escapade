package rerrors

// ErrorWrongStatus return error, returns an error if
// someone tries to perform an action that is valid not with current
// room status
func ErrorWrongStatus() error {
	return New("Action refused. Too late or too early to perform your action")
}

// ErrorWrongBordersParams return error if the borders params are wrong
func ErrorWrongBordersParams(w, h, m int32) error {
	return Errorf("Wrong weight('%d') or height('%d') or mines around(%d)", w, h, m)
}

// ErrorCellOutside return error if the specified cel
//l is out of the field
func ErrorCellOutside() error {
	return New("Action refused. Cell not inside field")
}

// ErrorPlayerFinished return an error if the player, being dead,
// tries to perform an action available only to live players
func ErrorPlayerFinished() error {
	return New("Action refused. You died")
}

// ErrorLobbyCantCreateRoom return an error if lobby cant create room
func ErrorLobbyCantCreateRoom() error {
	return New("Cant create room in lobby")
}

// ErrorInvalidRoomSettings returns an error if the player tried to
// create a room with invalid characteristics
func ErrorInvalidRoomSettings() error {
	return New("Invalid roomSettings")
}

// ErrorLobbyDone return an error if lobby is turned off
func ErrorLobbyDone() error {
	return New("Lobby is turned off")
}

// NoWebSocketOrUser return an error if there is no websocket
// connection of user
func NoWebSocketOrUser() error {
	return New("Websocket connection or user is not set")
}

// ErrorRoomOrLobbyDone return an error if lobby or room is turned off
func ErrorRoomOrLobbyDone() error {
	return New("Lobby or room is turned off")
}

// ErrorConnectionDone return an error if connection died
func ErrorConnectionDone() error {
	return New("Connection died")
}

// ErrorRoomDone return an error if room deleted
func ErrorRoomDone() error {
	return New("Room deleted")
}
