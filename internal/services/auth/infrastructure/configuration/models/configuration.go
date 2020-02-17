package models

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/infrastructure/configuration"
)

type Configuration struct {
	GCInterval models.Duration `json:"gc"`
	JWT        string          `json:"jwt"`
	Token      Token           `json:"token"`
	WhiteList  Whitelist       `json:"whitelist"`
}

func (co *Configuration) Get() configuration.Configuration {
	return configuration.Configuration{
		GCInterval: co.GCInterval.Duration,
		JWT:        co.JWT,
		Token:      co.Token.Get(),
		WhiteList:  co.WhiteList.Get(),
	}

}
func (co *Configuration) Set(c configuration.Configuration) {
	co.GCInterval.Duration = c.GCInterval
	co.JWT = c.JWT
	co.Token.Set(c.Token)
	co.WhiteList.Set(c.WhiteList)

}
