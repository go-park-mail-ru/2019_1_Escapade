package loader

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/domens/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
)

// LoaderConfig load configuration file
type LoaderConfig struct {
	l   infrastructure.LoaderJSONI
	c   *config.Configuration
	err infrastructure.ErrorTrace
}

func NewLoader(
	l infrastructure.LoaderJSONI,
	trace infrastructure.ErrorTrace,
) *LoaderConfig {
	var loader LoaderConfig
	loader.l = l
	loader.c = &config.Configuration{}
	loader.l.Init(loader.c)
	return &loader
}

// Load main configuration
func (loader *LoaderConfig) Load() (*config.Configuration, error) {
	var (
		err error
		obj infrastructure.JSONtype
	)
	obj, err = loader.l.Load()
	if err != nil {
		return nil, err
	}
	c, ok := obj.(*config.Configuration)
	if !ok {
		return nil, loader.err.New(ErrCastFailed)
	}
	return c, nil
}
