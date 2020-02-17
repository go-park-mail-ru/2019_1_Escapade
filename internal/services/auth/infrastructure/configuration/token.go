package configuration

import "time"

type TokenRepository interface {
	Get() Token
	Set(Token)
}

type Token struct {
	AccessExpire, RefreshExpire time.Duration
	IsGenerateRefresh           bool
	Type                        string
}
