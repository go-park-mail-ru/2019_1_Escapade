package rerrors

// InvalidMessageID return error, when message id is invalid
func InvalidMessageID() error {
	return New("Invalid message id")
}

// InvalidMessage return error, when message is nil
func InvalidMessage() error {
	return New("Invalid message")
}

// InvalidChatID return error, when chat id is invalid
func InvalidChatID() error {
	return New("Invalid chat id")
}

// InvalidUser return error, when user is nil
func InvalidUser() error {
	return New("Invalid user")
}

func NoAuthFound() error {
	return New("No UserID in token")
}

func NoTokenFound() error {
	return New("No token")
}
