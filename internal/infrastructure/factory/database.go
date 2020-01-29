package factory

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/domens/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	rep "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/repository/database"
	uc "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/usecase/database"
)

func NewDatabase(c config.Configuration) infrastructure.DatabaseI {
	var r = rep.NewRepositoryConfig(c)
	var pg = uc.NewPostgresSQL(r)
	return pg
}
