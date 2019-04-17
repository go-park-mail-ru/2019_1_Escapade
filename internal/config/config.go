package config

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/aws/aws-sdk-go/aws"
)

// Configuration contains all types of configurations
type Configuration struct {
	Server    ServerConfig      `json:"server"`
	Cors      CORSConfig        `json:"cors"`
	DataBase  DatabaseConfig    `json:"dataBase"`
	Storage   FileStorageConfig `json:"storage"`
	Game      GameConfig        `json:"game"`
	Cookie    CookieConfig      `json:"cookie"`
	WebSocket WebSocketConfig   `json:"websocket"`
}

// ServerConfig set host, post and buffers sizes
type ServerConfig struct {
	Host            string `json:"host"`
	Port            string `json:"port"`
	ReadBufferSize  int    `json:"readBufferSize"`
	WriteBufferSize int    `json:"writeBufferSize"`
}

// CORSConfig set allowable origins, headers and methods
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
	DriverName   string `json:"driverName"`
	URL          string `json:"url"`
	MaxOpenConns int    `json:"maxOpenConns"`
	PageGames    int    `json:"pageGames"`
	PageUsers    int    `json:"pageUsers"`
}

// FileStorageConfig set, where avatars store and
//    what mode set to files/directories
type FileStorageConfig struct {
	PlayersAvatarsStorage string `json:"playersAvatarsStorage"`
	DefaultAvatar         string `json:"defaultAvatar"`
	Region                string `json:"region"`
	Endpoint              string `json:"endpoint"`
	AwsConfig             *aws.Config
}

// GameConfig set, how much rooms server can create and
// how mich connections can join and execute together
type GameConfig struct {
	RoomsCapacity int `json:"roomsCapacity"`
	LobbyJoin     int `json:"lobbyJoin"`
	LobbyRequest  int `json:"lobbyRequest"`
}

// CookieConfig set cookie name, path, length, expiration time
// and HTTPonly flag
type CookieConfig struct {
	NameCookie     string `json:"nameCookie"`
	PathCookie     string `json:"pathCookie"`
	LengthCookie   int    `json:"lengthCookie"`
	LifetimeCookie int    `json:"lifetimeCookie"`
	HTTPOnly       bool   `json:"httpOnly"`
}

// WebSocketConfig set timeouts
type WebSocketConfig struct {
	WriteWait      int   `json:"writeWait"`
	PongWait       int   `json:"pongWait"`
	PingPeriod     int   `json:"pingPeriod"`
	MaxMessageSize int64 `json:"maxMessageSize"`
}

// WebSocketSettings set timeouts
type WebSocketSettings struct {
	WriteWait      time.Duration `json:"writeWait"`
	PongWait       time.Duration `json:"pongWait"`
	PingPeriod     time.Duration `json:"pingPeriod"`
	MaxMessageSize int64         `json:"maxMessageSize"`
}

// Init load configuration file
func Init(path string) (conf *Configuration, err error) {
	conf = &Configuration{}
	var data []byte

	if data, err = ioutil.ReadFile(path); err != nil {
		return
	}
	err = json.Unmarshal(data, conf)
	conf.Storage.AwsConfig = &aws.Config{
		Region:   aws.String("ru-msk"),
		Endpoint: aws.String("https://hb.bizmrg.com")}

	return
}
