package configuration

import "time"

type PhotoRepository interface {
	Get() Photo
	Set(Photo)
}

type Photo struct {
	Region, Endpoint                     string
	PlayersAvatarsStorage, DefaultAvatar string // TODO убрать отсюда бизнес логику
	MaxFileSize                          int64
	Expire                               time.Duration
	AllowedFileTypes                     []string

	AccessKey, SecretKey string
}
