package rerrors

import "errors"

// ErrorInvalidName call it, if client give you
// 	invalid username
<<<<<<< HEAD
func ErrorInvalidUserID() error {
	return errors.New("No such id")
=======
func ErrorHandshake() error {
	return errors.New("HandshakeError")
}

// ErrorInvalidName call it, if client give you
// 	invalid username
func ErrorNotWebsocket() error {
	return errors.New("Not a websocket")
}

// ErrorInvalidName call it, if client give you
// 	invalid username
func ErrorInvalidUserID() error {
	return errors.New("No such id")
}

// ErrorInvalidName call it, if client give you
// 	invalid username
func ErrorInvalidName() error {
	return errors.New("Invalid username")
>>>>>>> 508037185fc39abb3d6ee56a9fd2c48bac220f58
}

// ErrorInvalidNameOrEmail call it, if client give you
// 	invalid username or email
func ErrorInvalidNameOrEmail() error {
	return errors.New("Invalid username or email")
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
