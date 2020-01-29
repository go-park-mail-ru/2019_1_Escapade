package config

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/domens"
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
	Name        string          `json:"name"`
	Polling     domens.Duration `json:"polling"`
	CounterDrop int             `json:"drop"`
	Tag         string          `json:"tag"`
}

// Timeouts of the connection to the server
//easyjson:json
type Timeouts struct {
	TTL domens.Duration `json:"ttl"`

	Read  domens.Duration `json:"read"`
	Write domens.Duration `json:"write"`
	Idle  domens.Duration `json:"idle"`
	Wait  domens.Duration `json:"wait"`
	Exec  domens.Duration `json:"exec"`

	Prepare domens.Duration `json:"prepare"`
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
	DriverName           string          `json:"driverName"`
	URL                  string          `json:"url"`
	ConnectionString     string          `json:"connectionString"`
	AuthConnectionString string          `json:"authConnectionString"`
	MaxOpenConns         int             `json:"maxOpenConns"`
	MaxIdleConns         int             `json:"maxIdleConns"`
	MaxLifetime          domens.Duration `json:"maxLifetime"`
	PageGames            int             `json:"pageGames"`
	PageUsers            int             `json:"pageUsers"`
}

// WebSocket set timeouts
//easyjson:json
type WebSocket struct {
	WriteWait        domens.Duration `json:"writeWait"`
	PongWait         domens.Duration `json:"pongWait"`
	PingPeriod       domens.Duration `json:"pingPeriod"`
	HandshakeTimeout domens.Duration `json:"handshakeTimeout"`
	MaxMessageSize   int64           `json:"maxMessageSize"`
	ReadBufferSize   int             `json:"readBufferSize"`
	WriteBufferSize  int             `json:"writeBufferSize"`
}
