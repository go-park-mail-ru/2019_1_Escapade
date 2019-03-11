package api

import "errors"

// ErrorInvalidName call it, if client give you
// 	invalid username as a request parameter
func ErrorInvalidName() error {
	return errors.New("Invalid username")
}

// ErrorInvalidPage call it, if client give you
// 	invalid username as a request parameter
func ErrorInvalidPage() error {
	return errors.New("Invalid page")
}

// ErrorUserNotFound call it, if you cant
// 	find user
func ErrorUserNotFound() error {
	return errors.New("User not found")
}

// ErrorUsersNotFound call it, if you cant
// 	find users
func ErrorUsersNotFound() error {
	return errors.New("Users not found")
}

// ErrorAuthorization call it, if client
// 	hasnt session cookie
func ErrorAuthorization() error {
	return errors.New("Required authorization")
}

// ErrorAvatarNotFound call it, if you cant
// find avatar
func ErrorAvatarNotFound() error {
	return errors.New("Avatar not found")
}

// ErrorInvalidFile call it, if client give you
// 	invalid file as a request parameter
func ErrorInvalidFile() error {
	return errors.New("Invalid file")
}

// ErrorInvalidFileFormat call it, if client give you
// 	invalid file as a request parameter
func ErrorInvalidFileFormat() error {
	return errors.New("Invalid file format. Use .png or .jpg only")
}

// ErrorDataBase  call it, if error in database
func ErrorDataBase() error {
	return errors.New("DataBase error")
}

// ErrorServer  call it, if error internal
func ErrorServer() error {
	return errors.New("Server error")
}
