package rerrors

// ErrorUserIsExist call it, if client give you
// 	name, that exists in database
func ErrorUserIsExist() error {
	return New("Username with such name exists")
}

func UserExistWrapper(err error) error {
	return New("Username with such name exists. More: " + err.Error())
}

// ErrorNameIstaken call it, if client give you
// 	name, that exists in database
func ErrorNameIstaken() error {
	return New("Username is taken")
}

// ErrorInvalidPassword call it, if client give you
// 	invalid password
func ErrorInvalidPassword() error {
	return New("Invalid password")
}

// ErrorInvalidName call it, if client give you
// 	invalid username
func ErrorInvalidName() error {
	return New("Invalid username")
}

// ErrorUserNotFound call it, if you cant
// 	find user
func ErrorUserNotFound() error {
	return New("User not found")
}

// NoUserWrapper - wrapper for an no-user error
func NoUserWrapper(err error) error {
	return New("User not found. More: " + err.Error())
}
