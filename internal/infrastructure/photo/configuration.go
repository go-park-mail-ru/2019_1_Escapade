package photo

import "github.com/go-park-mail-ru/2019_1_Escapade/internal/domens"

import "time"

// AwsPublicConfig public aws information as region and endpoint
type PhotoPublic struct {
	Region                string
	Endpoint              string
	PlayersAvatarsStorage string // TODO убрать отсюда бизнес логику
	DefaultAvatar         string // TODO убрать отсюда бизнес логику
	MaxFileSize           int64
	Expire                time.Duration
	AllowedFileTypes      []string
}

//easyjson:json
type PhotoPublicJSON struct {
	Region                string          `json:"region"`
	Endpoint              string          `json:"endpoint"`
	PlayersAvatarsStorage string          `json:"playersAvatarsStorage"`
	DefaultAvatar         string          `json:"defaultAvatar"`
	MaxFileSize           int64           `json:"maxFileSize"`
	Expire                domens.Duration `json:"expire"`
	AllowedFileTypes      []string        `json:"allowedFileTypes"`
}

func (p PhotoPublicJSON) Get() PhotoPublic {
	return PhotoPublic{
		Region:                p.Region,
		Endpoint:              p.Endpoint,
		PlayersAvatarsStorage: p.PlayersAvatarsStorage,
		DefaultAvatar:         p.DefaultAvatar,
		MaxFileSize:           p.MaxFileSize,
		Expire:                p.Expire.Duration,
		AllowedFileTypes:      p.AllowedFileTypes,
	}
}

// AwsPrivateConfig private aws information. Need another json.
//easyjson:json
type PhotoPrivate struct {
	AccessURL string
	AccessKey string
	SecretURL string
	SecretKey string
}

type PhotoPrivateJSON struct {
	AccessURL string `json:"accessUrl"`
	AccessKey string `json:"accessKey"`
	SecretURL string `json:"secretUrl"`
	SecretKey string `json:"secretKey"`
}

func (p PhotoPrivateJSON) Get() PhotoPrivate {
	return PhotoPrivate(p)
}
