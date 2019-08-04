package rerrors

import "errors"

// InvalidMessageID return error, when message id is invalid
func InvalidMessageID() error {
	return errors.New("Invalid message id")
}

// InvalidMessage return error, when message is nil
func InvalidMessage() error {
	return errors.New("Invalid message")
}

// InvalidChatID return error, when chat id is invalid
func InvalidChatID() error {
	return errors.New("Invalid chat id")
}

// InvalidUser return error, when user is nil
func InvalidUser() error {
	return errors.New("Invalid user")
}
