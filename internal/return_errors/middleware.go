package rerrors

import "errors"

// DELETE!!!!! IT IS COPY OF ErrorAuthorization
// ErrorNoCookie shows that there is no session cookie
func ErrorNoCookie() error {
	return errors.New("Not found session cookie")
}

// ErrorAuthorization call it, if client
// 	hasnt session cookie
func ErrorAuthorization() error {
	return errors.New("Required authorization")
}

// ErrorPanic shows that panic happened
func ErrorPanic() error {
	return errors.New("Panic happened")
}

// ErrorCORS shows that domen cant access to server by CORS
func ErrorCORS(origin string) error {
	return errors.New("CORS couldnt find " + origin)
}
