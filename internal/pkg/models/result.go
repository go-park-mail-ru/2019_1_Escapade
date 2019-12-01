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
