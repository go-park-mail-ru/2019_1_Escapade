package models

// Result is the query result with detailed explanation
//easyjson:json
type Result struct {
	Place   string `json:"place"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Response is game struct, that frontend is waiting from backend
//easyjson:json
type Response struct {
	Type    string      `json:"type"`
	Message string      `json:"message,omitempty"`
	Value   interface{} `json:"value"`
}

// FailFlagSet is called when room cant set flag
func FailFlagSet(value interface{}, err error) Response {
	return Response{
		Type:    "FailFlagSet",
		Message: err.Error(),
		Value:   value,
	}
}

// RandomFlagSet is called when any player set his flag at the same as any other
func RandomFlagSet(value interface{}) Response {
	return Response{
		Type:    "ChangeFlagSet",
		Message: "The cell you have selected is chosen by another person.",
		Value:   value,
	}
}
