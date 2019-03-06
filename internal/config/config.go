package config

import (
	"encoding/json"
	"io/ioutil"
)

type Configuration struct {
	Server ServerConfig `json:"server"`
	Cors   CORSConfig   `json:"cors"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type CORSConfig struct {
	Origins     []string `json:"origins"`
	Headers     []string `json:"headers"`
	Credentials string   `json:"credentials"`
	Methods     []string `json:"methods"`
}

func Init(path string) (conf Configuration, return_error error) {
	conf = Configuration{}
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return_error = err
		return
	}

	if return_error = json.Unmarshal(data, &conf); return_error != nil {
		return
	}

	return
}
