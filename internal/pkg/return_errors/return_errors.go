package rerrors

func ID() error {
	return New("No such id")
}

func ErrorInvalidUserID() error {
	return New("No such id")
}

// ErrorInvalidPage call it, if client give you
// 	invalid page
func ErrorInvalidPage() error {
	return New("Invalid page")
}

// ErrorUsersNotFound call it, if you cant
// 	find users
func ErrorUsersNotFound() error {
	return New("Users not found")
}

// NoUpdate godoc
func NoUpdate() error {
	return New("No updated fields")
}

// ErrorGamesNotFound call it, if you cant
// 	find games
func ErrorGamesNotFound() error {
	return New("Games not found")
}

// ErrorAvatarNotFound call it, if you cant
// find avatar
func ErrorAvatarNotFound() error {
	return New("Avatar not found")
}

func NoAvatarWrapper(err error) error {
	return New("Avatar not found. More: " + err.Error())
}

// ErrorDataBase  call it, if error in database
func ErrorDataBase() error {
	return New("DataBase error")
}

func DatabaseWrapper(err error) error {
	return New("Database error. More: " + err.Error())
}

// ErrorServer  call it, if error internal
func ErrorServer() error {
	return New("Server error")
}

func ServerWrapper(err error) error {
	return New("Server error. More: " + err.Error())
}
