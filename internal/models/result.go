package models

// Result is the query result with detailed explanation
type Result struct {
	Type    string `json:"type"`
	Place   string `json:"place"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type Response struct {
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Value   interface{} `json:"value"`
}

func FailFlagSet(value interface{}, err error) Response {
	return Response{
		Type:    "FailFlagSet",
		Message: err.Error(),
		Value:   value,
	}
}

func RandomFlagSet(value interface{}) Response {
	return Response{
		Type:    "ChangeFlagSet",
		Message: "The cell you have selected is chosen by another person.",
		Value:   value,
	}
}
