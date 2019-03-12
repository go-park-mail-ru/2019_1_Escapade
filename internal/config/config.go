package config

import (
	"encoding/json"
	"io/ioutil"
)

type Configuration struct {
	Server   ServerConfig      `json:"server"`
	Cors     CORSConfig        `json:"cors"`
	DataBase DatabaseConfig    `json:"dataBase"`
	Storage  FileStorageConfig `json:"storage"`
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

// DatabaseConfig set type of database management system
//   the url of connection string, max amount of
//   connections, tables, sizes of page  of gamers
//   and users
type DatabaseConfig struct {
	DriverName   string   `json:"driverName"`
	URL          string   `json:"url"`
	MaxOpenConns int      `json:"maxOpenConns"`
	Tables       []string `json:"tables"`
	PageGames    int      `json:"pageGames"`
	PageUsers    int      `json:"pageUsers"`
}

// FileStorageConfig set, where avatars store and
//    what mode set to files/directories
type FileStorageConfig struct {
	PlayersAvatarsStorage string `json:"playersAvatarsStorage"`
	FileMode              int    `json:"fileMode"`
}

// Init load configuration file
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
