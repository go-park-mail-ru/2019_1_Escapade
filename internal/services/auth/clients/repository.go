package clients

import "gopkg.in/oauth2.v3/models"

type RepositoryI interface {
	Get() []*models.Client
}

type RepositoryHC struct{}

func (rep *RepositoryHC) Get() []*models.Client {
	return []*models.Client{&models.Client{
		ID:     "1",
		Secret: "1",
		Domain: "api.consul.localhost",
	}}
}
