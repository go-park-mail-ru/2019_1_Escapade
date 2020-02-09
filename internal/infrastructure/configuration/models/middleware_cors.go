package models

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/configuration"
)

//easyjson:json
type Cors struct {
	Origins     []string `json:"origins"`
	Headers     []string `json:"headers"`
	Methods     []string `json:"methods"`
	Credentials string   `json:"credentials"`
}

func (cors Cors) Get() configuration.Cors {
	return configuration.Cors(cors)
}

func (cors Cors) Set(c configuration.Cors) {
	cors.Origins = c.Origins
	cors.Headers = c.Headers
	cors.Methods = c.Methods
	cors.Credentials = c.Credentials
}
