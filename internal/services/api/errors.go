package api

import "errors"

// ErrorInvalidName call it, if client give you
// 	invalid username as a request parameter
func ErrorInvalidName() error {
	return errors.New("Invalid username")
}

// ErrorUserNotFound call it, if you cant
// 	find user
func ErrorUserNotFound() error {
	return errors.New("User not found")
}
