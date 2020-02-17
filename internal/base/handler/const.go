package handler

type ContextKey string

const (
	/*
		RequestParamsInt32 get all parameters from path via IDsFromPath and UserID
		via GetUserIDFromAuthRequest(if 'withAuth' true)
		userID is placed in map with key set in UserIDKey
	*/
	ContextUserKey ContextKey = "userID"
	UserIDKey                 = "auth_user_id"

	ErrInvalidID        = "Invalid id"
	ErrInvalidName      = "Invalid user's name"
	ErrInvalidPassword  = "Invalid user's password"
	ErrNoBody           = "No request body"
	ErrNoAuthFound      = "NO auth found"
	ErrInvalidJSON      = "Cant decode object as json"
	ErrMethodNotAllowed = "Method not allowed"
	ErrNoUpdate         = "Cant update object"

	MinNameLength     = 3
	MinPasswordLength = 3

	NoResult = 0
)
