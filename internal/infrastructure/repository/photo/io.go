package photo

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/entity"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/repository/loader"
)

type RepositoryIO struct {
	publicConfig  *entity.PhotoPublicConfig
	privateConfig *entity.PhotoPrivateConfig
}

func NewRepositoryIO(publicConfigPath, privateConfigPath string) (*RepositoryIO, error) {
	var rep RepositoryIO
	rep.publicConfig = new(entity.PhotoPublicConfig)
	rep.privateConfig = new(entity.PhotoPrivateConfig)
	var (
		publicLoader  = loader.NewRepositoryIO(publicConfigPath)
		privateLoader = loader.NewRepositoryIO(privateConfigPath)
	)
	publicLoader.Init(rep.publicConfig)
	privateLoader.Init(rep.privateConfig)

	public, err := publicLoader.Load()
	if err != nil {
		return nil, err
	}
	private, err := privateLoader.Load()
	if err != nil {
		return nil, err
	}

	rep.publicConfig = public.(*entity.PhotoPublicConfig)
	rep.privateConfig = private.(*entity.PhotoPrivateConfig)
	return &rep, nil
}

func (rep *RepositoryIO) GetPublic() entity.PhotoPublicConfig {
	return *rep.publicConfig
}

func (rep *RepositoryIO) GetPrivate() entity.PhotoPrivateConfig {
	return *rep.privateConfig
}
