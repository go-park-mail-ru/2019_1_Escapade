package rerrors

// ErrorNoBody call it, if client
// didnt send you body, when you need it
func ErrorNoBody() error {
	return New("Cant found parameters")
}

// ErrorInvalidJSON call it, if client
// send you invalid json
func ErrorInvalidJSON() error {
	return New("Found invalid json")
}

// ErrorMethodNotAllowed call it,
//  if that method not allowed
func ErrorMethodNotAllowed() error {
	return New("Method not allowed")
}
