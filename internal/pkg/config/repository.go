package config

import "io/ioutil"

//go:generate $GOPATH/bin/mockery -name "RepositoryI"
type RepositoryI interface {
	Load(path string) (*Configuration, error)
}

type RepositoryFS struct{}

func (rfs *RepositoryFS) Load(path string) (*Configuration, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	conf := new(Configuration)
	err = conf.UnmarshalJSON(data)
	return conf, err
}
