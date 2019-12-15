package constants

import (
	"io/ioutil"
)

type RepositoryI interface {
	GetRoom(path string) (RoomConfiguration, error)
	GetField(path string) (FieldConfiguration, error)
}

// RepositoryFS get Room and Field configuration from file system
type RepositoryFS struct{}

// getRoom load Room config(.json file) from FS by its path
func (rfs *RepositoryFS) GetRoom(path string) (RoomConfiguration, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return RoomConfiguration{}, err
	}

	var tmp *RoomConfiguration
	if err = tmp.UnmarshalJSON(data); err != nil {
		return RoomConfiguration{}, err
	}
	return RoomConfiguration{}, nil
}

// getField load field config(.json file) from FS by its path
func (rfs *RepositoryFS) GetField(path string) (FieldConfiguration, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return FieldConfiguration{}, err
	}

	var tmp *FieldConfiguration
	if err = tmp.UnmarshalJSON(data); err != nil {
		return FieldConfiguration{}, err
	}
	return FieldConfiguration{}, nil
}
