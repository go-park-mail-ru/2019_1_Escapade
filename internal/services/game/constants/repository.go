package constants

import "io/ioutil"

type RepositoryI interface {
	getRoom(path string) (roomConfiguration, error)
	getField(path string) (fieldConfiguration, error)
}

// RepositoryFS get Room and Field configuration from file system
type RepositoryFS struct{}

// getRoom load Room config(.json file) from FS by its path
func (rfs *RepositoryFS) getRoom(path string) (roomConfiguration, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return roomConfiguration{}, err
	}

	var tmp *roomConfiguration
	if err = tmp.UnmarshalJSON(data); err != nil {
		return roomConfiguration{}, err
	}
	return roomConfiguration{}, nil
}

// getField load field config(.json file) from FS by its path
func (rfs *RepositoryFS) getField(path string) (fieldConfiguration, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fieldConfiguration{}, err
	}

	var tmp *fieldConfiguration
	if err = tmp.UnmarshalJSON(data); err != nil {
		return fieldConfiguration{}, err
	}
	return fieldConfiguration{}, nil
}
