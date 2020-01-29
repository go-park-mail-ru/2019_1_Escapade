package infrastructure

// LoaderI interface of loading struct
type LoaderJSONI interface {
	Init(object JSONtype)
	Load() (JSONtype, error)
}

// JSONtype is interface to be sent by json
type JSONtype interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}
