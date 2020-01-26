package loader

import (
	"io/ioutil"
	"os"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/photo"
)

// Loader load configuration files
type Loader struct {
	infrastructure.WithExtra

	path   string
	rep    config.RepositoryI
	config *config.Configuration

	CallExtra func() error
}

// Init initialize struct
func NewLoader(rep config.RepositoryI, path string) *Loader {
	var loader Loader
	loader.rep = rep
	loader.path = path
	return &loader
}

// Load main configuration
func (loader *Loader) Load() error {
	data, err := ioutil.ReadFile(loader.path)
	if err != nil {
		return err
	}

	loader.config = new(config.Configuration)
	err = loader.config.UnmarshalJSON(data)
	if loader.config, err = loader.rep.Load(loader.path); err != nil {
		return err
	}
	loader.config.Init(os.Getenv("AUTH_ADDRESS"))
	return nil
}

// Get main configuration
func (loader *Loader) Get() *config.Configuration {
	return loader.config
}

// LoadPhoto - extra load photo
func (loader *Loader) LoadPhoto(public, private string) error {
	return photo.Init(public, private)
}
