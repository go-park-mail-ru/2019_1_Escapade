package rerrors

import "errors"

// ErrorAuthorization call it, if client
// 	hasnt session cookie
func ErrorAuthorization() error {
	return errors.New("Required authorization")
}

// ErrorTokenType wrong type of token
func ErrorTokenType() error {
	return errors.New("Wrong type of token")
}

// ErrorPanic shows that panic happened
func ErrorPanic() error {
	return errors.New("Panic happened")
}

// ErrorCORS shows that domen cant access to server by CORS
func ErrorCORS(origin string) error {
	return errors.New("CORS couldnt find " + origin)
}
