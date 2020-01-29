package rerrors

func ErrorHandshake() error {
	return New("HandshakeError")
}

func ErrorNotWebsocket() error {
	return New("Not a websocket")
}
