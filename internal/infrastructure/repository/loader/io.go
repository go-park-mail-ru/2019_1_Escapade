package loader

import (
	"io/ioutil"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
)

type RepositoryIO struct {
	Object infrastructure.JSONtype
	Path   string
}

func NewRepositoryIO(path string) *RepositoryIO {
	return &RepositoryIO{
		Path: path,
	}
}

func (rep *RepositoryIO) Load() (infrastructure.JSONtype, error) {
	data, err := ioutil.ReadFile(rep.Path)
	if err != nil {
		return nil, err
	}

	err = rep.Object.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}

	return rep.Object, nil
}

func (rep *RepositoryIO) Init(object infrastructure.JSONtype) {
	rep.Object = object
}
