package rerrors

import "errors"

// ErrorInvalidName call it, if client give you
// 	invalid username
func ErrorHandshake() error {
	return errors.New("HandshakeError")
}

// ErrorInvalidName call it, if client give you
// 	invalid username
func ErrorNotWebsocket() error {
	return errors.New("Not a websocket")
}
