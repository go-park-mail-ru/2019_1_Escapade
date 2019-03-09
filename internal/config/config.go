package config

import (
	"encoding/json"
	"io/ioutil"
)

type Configuration struct {
	Server   ServerConfig   `json:"server"`
	Cors     CORSConfig     `json:"cors"`
	DataBase DatabaseConfig `json:"dataBase"`
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

type DatabaseConfig struct {
	DriverName   string   `json:"driverName"`
	URL          string   `json:"url"`
	MaxOpenConns int      `json:"maxOpenConns"`
	Tables       []string `json:"tables"`
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
