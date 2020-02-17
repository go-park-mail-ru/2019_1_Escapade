package infrastructure

import (
	"mime/multipart"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

type PhotoService interface {
	GetImages(users ...*models.UserPublicInfo)
	SaveImage(key string, file multipart.File) error
	GetImage(key string) (string, error)
	DeleteImage(key string) error

	GetDefaultAvatar() string
	MaxFileSize() int64
	AllowedFileTypes() []string
}

type PhotoServiceNil struct{}

func (pn *PhotoServiceNil) GetImages(
	users ...*models.UserPublicInfo,
) {
}
func (pn *PhotoServiceNil) SaveImage(
	key string,
	file multipart.File,
) error {
	return nil
}
func (pn *PhotoServiceNil) GetImage(key string) (string, error) {
	return "PhotoServiceNil", nil
}
func (pn *PhotoServiceNil) DeleteImage(key string) error {
	return nil
}

func (pn *PhotoServiceNil) GetDefaultAvatar() string {
	return "PhotoServiceNil"
}
func (pn *PhotoServiceNil) MaxFileSize() int64 {
	return 0
}
func (pn *PhotoServiceNil) AllowedFileTypes() []string {
	return nil
}
