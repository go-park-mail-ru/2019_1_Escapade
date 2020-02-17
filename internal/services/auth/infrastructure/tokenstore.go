package infrastructure

import (
	"gopkg.in/oauth2.v3"
)

type TokenStore interface {
	Create(info oauth2.TokenInfo) error

	RemoveByCode(code string) error
	RemoveByAccess(access string) error
	RemoveByRefresh(refresh string) error

	GetByCode(code string) (oauth2.TokenInfo, error)
	GetByAccess(access string) (oauth2.TokenInfo, error)
	GetByRefresh(refresh string) (oauth2.TokenInfo, error)

	Close() error
}
