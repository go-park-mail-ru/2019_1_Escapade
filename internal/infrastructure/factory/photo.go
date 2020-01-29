package factory

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	rep "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/repository/photo"
	uc "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/usecase/photo"
)

func NewPhotoService(publicConfigPath, privateConfigPath string) (infrastructure.PhotoServiceI, error) {
	r, err := rep.NewRepositoryIO(publicConfigPath, privateConfigPath)
	if err != nil {
		return nil, err
	}
	return uc.NewAWSService(r), nil
}
