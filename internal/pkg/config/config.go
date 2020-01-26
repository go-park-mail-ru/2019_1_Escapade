package config

import (
	"os"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

// Configuration contains all types of configurations
//easyjson:json
type Configuration struct {
	Server     Server          `json:"server"`
	Cors       CORS            `json:"cors"`
	DataBase   Database        `json:"dataBase"`
	Game       Game            `json:"game"`
	Cookie     Cookie          `json:"session"`
	WebSocket  WebSocket       `json:"websocket"`
	Required   RequiredService `json:"required"`
	Auth       Auth            `json:"auth"`
	AuthClient AuthClient      `json:"authClient"`
}

// Server set host, post and buffers sizes
//easyjson:json
type Server struct {
	Name           string   `json:"name"`
	MaxConn        int      `json:"maxConn"`
	MaxHeaderBytes int      `json:"maxHeaderBytes"`
	Timeouts       Timeouts `json:"timeouts"`
	EnableTraefik  bool     `json:"enableTraefik"`
}

// RequiredService that is required for the correct working of this one
//easyjson:json
type RequiredService struct {
	Name        string   `json:"name"`
	Polling     Duration `json:"polling"`
	CounterDrop int      `json:"drop"`
	Tag         string   `json:"tag"`
}

// Timeouts of the connection to the server
//easyjson:json
type Timeouts struct {
	TTL Duration `json:"ttl"`

	Read  Duration `json:"read"`
	Write Duration `json:"write"`
	Idle  Duration `json:"idle"`
	Wait  Duration `json:"wait"`
	Exec  Duration `json:"exec"`

	Prepare Duration `json:"prepare"`
}

// CORS set allowable origins, headers and methods
//easyjson:json
type CORS struct {
	Origins     []string `json:"origins"`
	Headers     []string `json:"headers"`
	Methods     []string `json:"methods"`
	Credentials string   `json:"credentials"`
}

// Database set type of database management system
//   the url of connection string, max amount of
//   connections, tables, sizes of page  of gamers
//   and users
//easyjson:json
type Database struct {
	DriverName           string   `json:"driverName"`
	URL                  string   `json:"url"`
	ConnectionString     string   `json:"connectionString"`
	AuthConnectionString string   `json:"authConnectionString"`
	MaxOpenConns         int      `json:"maxOpenConns"`
	MaxIdleConns         int      `json:"maxIdleConns"`
	MaxLifetime          Duration `json:"maxLifetime"`
	PageGames            int      `json:"pageGames"`
	PageUsers            int      `json:"pageUsers"`
}

// WebSocket set timeouts
//easyjson:json
type WebSocket struct {
	WriteWait        Duration `json:"writeWait"`
	PongWait         Duration `json:"pongWait"`
	PingPeriod       Duration `json:"pingPeriod"`
	HandshakeTimeout Duration `json:"handshakeTimeout"`
	MaxMessageSize   int64    `json:"maxMessageSize"`
	ReadBufferSize   int      `json:"readBufferSize"`
	WriteBufferSize  int      `json:"writeBufferSize"`
}

func (conf *Configuration) Init(addr string) *Configuration {
	utils.Debug(false, " Look at config", conf.Auth.Salt, conf.Auth.AccessTokenExpire,
		conf.Auth.RefreshTokenExpire, conf.Auth.IsGenerateRefresh, conf.Auth.WithReserve,
		conf.Auth.TokenType, conf.Auth.WhiteList)
	utils.Debug(false, " Info:", conf.Server.Name, conf.Server.MaxHeaderBytes)
	utils.Debug(false, " Timeouts:", conf.Server.Timeouts.TTL,
		conf.Server.Timeouts.Read, conf.Server.Timeouts.Write, conf.Server.Timeouts.Idle,
		conf.Server.Timeouts.Wait, conf.Server.Timeouts.Exec)

	conf.setOauth2Config()
	conf.AuthClient.Address = addr
	return conf
}

// Init load configuration file and put part of parameters to Environment
func NewConfiguration(rep RepositoryI, path string) (*Configuration, error) {

	conf, err := rep.Load(path)
	if err != nil {
		return nil, err
	}

	utils.Debug(false, " Look at config", conf.Auth.Salt, conf.Auth.AccessTokenExpire,
		conf.Auth.RefreshTokenExpire, conf.Auth.IsGenerateRefresh, conf.Auth.WithReserve,
		conf.Auth.TokenType, conf.Auth.WhiteList)
	utils.Debug(false, " Info:", conf.Server.Name, conf.Server.MaxHeaderBytes)
	utils.Debug(false, " Timeouts:", conf.Server.Timeouts.TTL,
		conf.Server.Timeouts.Read, conf.Server.Timeouts.Write, conf.Server.Timeouts.Idle,
		conf.Server.Timeouts.Wait, conf.Server.Timeouts.Exec)

	conf.setOauth2Config()
	conf.AuthClient.Address = os.Getenv("AUTH_ADDRESS")
	return conf, err
}
