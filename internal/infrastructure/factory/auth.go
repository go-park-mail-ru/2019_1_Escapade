package factory

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/domens/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	rep "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/repository/auth"
	uc "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/usecase/auth"
)

func NewAuth(c config.Configuration, authAddr string) infrastructure.AuthService {
	var r = rep.NewRepositoryConfig(c, authAddr)
	var oauth2 = uc.NewOAuth2(r)
	return oauth2
}
