package models

// Result is the query result with detailed explanation
type Result struct {
	Place   string `json:"place"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}
