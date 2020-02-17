package models

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/infrastructure/configuration"
)

//easyjson:json
type Token struct {
	AccessExpire      models.Duration `json:"access_expire"`
	RefreshExpire     models.Duration `json:"refresh_expire"`
	IsGenerateRefresh bool            `json:"is_generate_refresh"`
	Type              string          `json:"type"`
}

// Get configuration.TokenGeneration from json model
// implementation of TokenGenerationRepository
func (a *Token) Get() configuration.Token {
	return configuration.Token{
		AccessExpire:      a.AccessExpire.Duration,
		RefreshExpire:     a.RefreshExpire.Duration,
		IsGenerateRefresh: a.IsGenerateRefresh,
		Type:              a.Type,
	}
}

// Set data from configuration.TokenGeneration
// implementation of TokenGenerationRepository
func (a *Token) Set(c configuration.Token) {
	a.AccessExpire.Init(c.AccessExpire)
	a.RefreshExpire.Init(c.RefreshExpire)
	a.IsGenerateRefresh = c.IsGenerateRefresh
	a.Type = c.Type
}
