package cleanenv

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Loader struct{}

func (l *Loader) Load(path string, cfg interface{}) error {
	return cleanenv.ReadConfig(path, cfg)
}
