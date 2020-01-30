package infrastructure

import (
	"mime/multipart"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/domens/models"
)

type PhotoServiceI interface {
	GetImages(users ...*models.UserPublicInfo)
	SaveImage(key string, file multipart.File) error
	GetImage(key string) (string, error)
	DeleteImage(key string) error

	GetDefaultAvatar() string
	MaxFileSize() int64
	AllowedFileTypes() []string
}