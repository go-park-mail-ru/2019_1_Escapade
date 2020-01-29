package auth

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/domens/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/entity"
)

type RepositoryConfig struct {
	entity.Auth
}

func NewRepositoryConfig(c config.Configuration, addr string) *RepositoryConfig {
	var dc RepositoryConfig
	dc.Salt = c.Auth.Salt
	dc.Cookie = c.Cookie
	dc.Auth.Auth = c.Auth
	dc.Client = c.AuthClient
	dc.Client.Address = addr
	return &dc
}

func (dc *RepositoryConfig) Get() entity.Auth {
	return dc.Auth
}
