package configuration

import "time"

type ConfigurationRepository interface {
	Get() Configuration
	Set(Configuration)
}

type Configuration struct {
	GCInterval time.Duration
	JWT        string
	Token      Token
	WhiteList  Whitelist
}
