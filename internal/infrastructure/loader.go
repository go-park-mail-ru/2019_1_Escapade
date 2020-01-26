package infrastructure

import "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"

// ConfigutaionLoaderI interface of loading configuration
type ConfigutaionLoaderI interface {
	WithExtraI
	Load() error
	Get() *config.Configuration
}
