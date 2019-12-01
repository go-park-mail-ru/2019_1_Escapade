package rerrors

import "errors"

// ErrorAuthorization call it, if client
// 	hasnt session cookie
func ErrorAuthorization() error {
	return errors.New("Required authorization")
}

// AuthWrapper - wrapper for authorization-related error
func AuthWrapper(err error) error {
	return errors.New("Required authorization. More: " + err.Error())
}

// AuthWrapper - wrapper for authorization-related error
func NoHeaders() error {
	return errors.New("No token headers")
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
func CORS(origin string) error {
	return errors.New("CORS couldnt find " + origin)
}
