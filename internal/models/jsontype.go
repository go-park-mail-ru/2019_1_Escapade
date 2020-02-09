package models

// JSONtype is interface to be sent by json
type JSONtype interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

// ModelUpdate is interface to update model
type ModelUpdate interface {
	JSONtype

	Update(JSONtype) bool
}
