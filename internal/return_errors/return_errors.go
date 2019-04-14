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

// ErrorInvalidName call it, if client give you
// 	invalid username
func ErrorInvalidUserID() error {
	return errors.New("No such id")
}

// ErrorInvalidName call it, if client give you
// 	invalid username
func ErrorInvalidName() error {
	return errors.New("Invalid username")
}

// ErrorInvalidNameOrEmail call it, if client give you
// 	invalid username or email
func ErrorInvalidNameOrEmail() error {
	return errors.New("Invalid username or email")
}

// ErrorInvalidPassword call it, if client give you
// 	invalid password
func ErrorInvalidPassword() error {
	return errors.New("Invalid password")
}

// ErrorInvalidEmail call it, if client give you
// 	invalid email
func ErrorInvalidEmail() error {
	return errors.New("Invalid email")
}

// ErrorInvalidPage call it, if client give you
// 	invalid page
func ErrorInvalidPage() error {
	return errors.New("Invalid page")
}

// ErrorNameIstaken call it, if client give you
// 	name, that exists in database
func ErrorNameIstaken() error {
	return errors.New("Username is taken")
}

// ErrorWrongPassword call it, if client give you
// 	wrong password+name/email
func ErrorWrongPassword() error {
	return errors.New("Password is taken")
}

// ErrorEmailIstaken call it, if client give you
// 	email, that exists in database
func ErrorEmailIstaken() error {
	return errors.New("Email is taken")
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

// ErrorGamesNotFound call it, if you cant
// 	find games
func ErrorGamesNotFound() error {
	return errors.New("Games not found")
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

// ErrorNoBody call it, if client
// didnt send you body, when you need it
func ErrorNoBody() error {
	return errors.New("Cant found parameters")
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
