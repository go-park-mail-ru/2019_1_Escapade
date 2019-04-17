package rerrors

import "errors"

// ErrorNoBody call it, if client
// didnt send you body, when you need it
func ErrorNoBody() error {
	return errors.New("Cant found parameters")
}

// ErrorInvalidJSON call it, if client
// send you invalid json
func ErrorInvalidJSON() error {
	return errors.New("Found invalid json")
}
