package rerrors

import "errors"

// ErrorUserIsExist call it, if client give you
// 	name, that exists in database
func ErrorUserIsExist() error {
	return errors.New("Username with such name exists")
}

// ErrorNameIstaken call it, if client give you
// 	name, that exists in database
func ErrorNameIstaken() error {
	return errors.New("Username is taken")
}

// ErrorInvalidPassword call it, if client give you
// 	invalid password
func ErrorInvalidPassword() error {
	return errors.New("Invalid password")
}

// ErrorInvalidName call it, if client give you
// 	invalid username
func ErrorInvalidName() error {
	return errors.New("Invalid username")
}

// ErrorUserNotFound call it, if you cant
// 	find user
func ErrorUserNotFound() error {
	return errors.New("User not found")
}
