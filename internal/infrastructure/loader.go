package infrastructure

// Loader interface of loading struct
type Loader interface {
	Load(path string, cfg interface{}) error
}
