package models

import "net/http"

// Result - every handler return it
type RequestResult struct {
	Code int
	Send JSONtype
	Err  error
}

type ResultFunc func(http.ResponseWriter, *http.Request) RequestResult
