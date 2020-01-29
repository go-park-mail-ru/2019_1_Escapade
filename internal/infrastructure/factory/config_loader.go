package factory

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/domens/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"

	rep "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/repository/loader"
	uc "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/usecase/loader"
)

func LoadConfiguration(trace infrastructure.ErrorTrace, path string) (*config.Configuration, error) {
	var r = rep.NewRepositoryIO(path)
	var loader = uc.NewLoader(r, trace)
	return loader.Load()
}
