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
func Init(path string) (conf *Configuration, err error) {
	conf = &Configuration{}
	var data []byte

	if data, err = ioutil.ReadFile(path); err != nil {
		return
	}
	err = json.Unmarshal(data, conf)

	return
}
