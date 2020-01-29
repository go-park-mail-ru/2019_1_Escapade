package rerrors

// ErrorAuthorization call it, if client
// 	hasnt session cookie
func ErrorAuthorization() error {
	return New("Required authorization")
}

// AuthWrapper - wrapper for authorization-related error
func AuthWrapper(err error) error {
	return New("Required authorization. More: " + err.Error())
}

// NoHeaders - no headers were given to authorize
func NoHeaders() error {
	return New("No token headers")
}

// ErrorTokenType wrong type of token
func ErrorTokenType() error {
	return New("Wrong type of token")
}

// CORS shows that domen cant access to server by CORS
func CORS(origin string) error {
	return New("CORS couldnt find " + origin)
}
