package server

import (
	"fmt"
	"os"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/photo"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
)

// ConfigutaionLoaderI interface of loading configuration
type ConfigutaionLoaderI interface {
	WithExtraI
	Load() error
	Get() *config.Configuration
}

// Loader load configuration files
type Loader struct {
	WithExtra

	path   string
	rep    config.RepositoryI
	config *config.Configuration

	CallExtra func() error
}

func (loader *Loader) InitAsFS(path string) *Loader {
	return loader.Init(new(config.RepositoryFS), path)
}

// Init initialize struct
func (loader *Loader) Init(rep config.RepositoryI, path string) *Loader {
	loader.rep = rep
	loader.path = path
	return loader
}

// Load main configuration
func (loader *Loader) Load() error {
	if loader.rep == nil {
		return re.InterfaceIsNil()
	}
	var err error
	fmt.Println("path is", loader.path)
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
