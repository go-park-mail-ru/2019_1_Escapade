package cleanenv

import (
	"flag"

	"github.com/ilyakaznacheev/cleanenv"
)

type Loader struct{}

func (l *Loader) Load(path string, cfg interface{}) error {
	return cleanenv.ReadConfig(path, cfg)
}

func (l *Loader) FUsage(fset *flag.FlagSet, cfg interface{}, wrap func()) func() {
	return cleanenv.FUsage(fset.Output(), cfg, nil, wrap)
}
