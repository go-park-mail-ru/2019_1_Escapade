package rerrors

import "errors"

// ErrorEmailIstaken call it, if client give you
// 	email, that exists in database
func ErrorEmailIstaken() error {
	return errors.New("Email is taken")
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

// ErrorInvalidEmail call it, if client give you
// 	invalid email
func ErrorInvalidEmail() error {
	return errors.New("Invalid email")
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
