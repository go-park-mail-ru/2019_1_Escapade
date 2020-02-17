package configuration

import "gopkg.in/oauth2.v3/models"

// AuthRepository manage getting and setting the configuration of AuthService
type WhitelistRepository interface {
	Get() Whitelist
	Set(Whitelist)
}

type Whitelist []*models.Client
