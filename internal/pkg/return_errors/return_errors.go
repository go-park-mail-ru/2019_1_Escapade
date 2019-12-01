package rerrors

import "errors"

func ID() error {
	return errors.New("No such id")
}

// ErrorInvalidName call it, if client give you
// 	invalid username
func ErrorInvalidUserID() error {
	return errors.New("No such id")
}

// ErrorInvalidPage call it, if client give you
// 	invalid page
func ErrorInvalidPage() error {
	return errors.New("Invalid page")
}

// ErrorUsersNotFound call it, if you cant
// 	find users
func ErrorUsersNotFound() error {
	return errors.New("Users not found")
}

// NoUpdate godoc
func NoUpdate() error {
	return errors.New("No updated fields")
}

// ErrorGamesNotFound call it, if you cant
// 	find games
func ErrorGamesNotFound() error {
	return errors.New("Games not found")
}

// ErrorAvatarNotFound call it, if you cant
// find avatar
func ErrorAvatarNotFound() error {
	return errors.New("Avatar not found")
}

func NoAvatarWrapper(err error) error {
	return errors.New("Avatar not found. More: " + err.Error())
}

// ErrorDataBase  call it, if error in database
func ErrorDataBase() error {
	return errors.New("DataBase error")
}

func DatabaseWrapper(err error) error {
	return errors.New("Database error. More: " + err.Error())
}

// ErrorServer  call it, if error internal
func ErrorServer() error {
	return errors.New("Server error")
}

func ServerWrapper(err error) error {
	return errors.New("Server error. More: " + err.Error())
}
