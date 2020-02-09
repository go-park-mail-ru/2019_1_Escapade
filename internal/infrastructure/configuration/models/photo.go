package models

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/configuration"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

//easyjson:json
type Photo struct {
	Region                string          `json:"region"`
	Endpoint              string          `json:"endpoint"`
	PlayersAvatarsStorage string          `json:"playersAvatarsStorage"`
	DefaultAvatar         string          `json:"defaultAvatar"`
	MaxFileSize           int64           `json:"maxFileSize"`
	Expire                models.Duration `json:"expire"`
	AllowedFileTypes      []string        `json:"allowedFileTypes"`

	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

func (p Photo) Get() configuration.Photo {
	return configuration.Photo{
		Region:                p.Region,
		Endpoint:              p.Endpoint,
		PlayersAvatarsStorage: p.PlayersAvatarsStorage,
		DefaultAvatar:         p.DefaultAvatar,
		MaxFileSize:           p.MaxFileSize,
		Expire:                p.Expire.Duration,
		AllowedFileTypes:      p.AllowedFileTypes,

		AccessKey: p.AccessKey,
		SecretKey: p.SecretKey,
	}
}

func (p Photo) Set(c configuration.Photo) {
	p.Region = c.Region
	p.Endpoint = c.Endpoint
	p.PlayersAvatarsStorage = c.PlayersAvatarsStorage
	p.DefaultAvatar = c.DefaultAvatar
	p.MaxFileSize = c.MaxFileSize
	p.Expire.Duration = c.Expire
	p.AllowedFileTypes = c.AllowedFileTypes

	p.AccessKey = c.AccessKey
	p.SecretKey = c.SecretKey

}
