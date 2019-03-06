package models

type Result struct {
	Place   string `json:"place"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}
