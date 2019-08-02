package config

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// Configuration contains all types of configurations
//easyjson:json
type Configuration struct {
	Server    ServerConfig    `json:"server"`
	Cors      CORSConfig      `json:"cors"`
	DataBase  DatabaseConfig  `json:"dataBase"`
	Game      GameConfig      `json:"game"`
	Session   SessionConfig   `json:"session"`
	WebSocket WebSocketConfig `json:"websocket"`
	Services  []Client        `json:"services"`
}

// ServerConfig set host, post and buffers sizes
//easyjson:json
type ServerConfig struct {
	Host      string `json:"host"`
	PortURL   string `json:"portUrl"`
	PortValue string `json:"portValue"`
	// timeouts in seconds
	ReadTimeoutS  int `json:"readTimeoutS"`
	WriteTimeoutS int `json:"writeTimeoutS"`
	IdleTimeoutS  int `json:"idleTimeoutS"`
	WaitTimeoutS  int `json:"waitTimeoutS"`
	ExecTimeoutS  int `json:"execTimeoutS"`
}

// CORSConfig set allowable origins, headers and methods
//easyjson:json
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
//easyjson:json
type DatabaseConfig struct {
	DriverName       string `json:"driverName"`
	URL              string `json:"url"`
	ConnectionString string `json:"connectionString"`
	MaxOpenConns     int    `json:"maxOpenConns"`
	PageGames        int    `json:"pageGames"`
	PageUsers        int    `json:"pageUsers"`
}

//easyjson:json
type FieldConfig struct {
	MinAreaSize    int `json:"minAreaSize"`
	MaxAreaSize    int `json:"maxAreaSize"`
	MinProbability int `json:"minProbability"`
	MaxProbability int `json:"maxProbability"`
}

// GameConfig set, how much rooms server can create and
// how much connections can join. Also there are flags:
// can server close rooms or not(for history mode),
// metrics should be recorded or not
//easyjson:json
type GameConfig struct {
	RoomsCapacity      int32        `json:"roomsCapacity"`
	ConnectionCapacity int32        `json:"connectionCapacity"`
	Location           string       `json:"location"`
	CanClose           bool         `json:"closeRoom"`
	Metrics            bool         `json:"metrics"`
	Field              *FieldConfig `json:"field"`
}

// AuthClient client of auth microservice
//easyjson:json
type AuthClient struct {
	URL     string `json:"url"`
	Address string `json:"address"`
}

type Client struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

// SessionConfig set cookie name, path, length, expiration time
// and HTTPonly flag
//easyjson:json
type SessionConfig struct {
	Name            string `json:"name"`
	Path            string `json:"path"`
	Length          int    `json:"length"`
	LifetimeSeconds int    `json:"lifetime"`
	HTTPOnly        bool   `json:"httpOnly"`
}

// WebSocketConfig set timeouts
//easyjson:json
type WebSocketConfig struct {
	WriteWait       int   `json:"writeWait"`
	PongWait        int   `json:"pongWait"`
	PingPeriod      int   `json:"pingPeriod"`
	MaxMessageSize  int64 `json:"maxMessageSize"`
	ReadBufferSize  int   `json:"readBufferSize"`
	WriteBufferSize int   `json:"writeBufferSize"`
}

// WebSocketSettings set timeouts
//easyjson:json
type WebSocketSettings struct {
	WriteWait       time.Duration `json:"writeWait"`
	PongWait        time.Duration `json:"pongWait"`
	PingPeriod      time.Duration `json:"pingPeriod"`
	MaxMessageSize  int64         `json:"maxMessageSize"`
	ReadBufferSize  int           `json:"readBufferSize"`
	WriteBufferSize int           `json:"writeBufferSize"`
}

func set(URL, value string) {
	if URL != "" && os.Getenv(URL) == "" {
		os.Setenv(URL, value)
	}
	utils.Debug(false, "environment -", URL, " :", value)
}

// InitEnvironment set environmental variables
func InitEnvironment(c *Configuration) {

	set(c.DataBase.URL, c.DataBase.ConnectionString)
	set(c.Server.PortURL, c.Server.PortValue)
}

// Init load configuration file and put part of parameters to Environment
func Init(path string) (conf *Configuration, err error) {
	conf = &Configuration{}
	var data []byte
	if data, err = ioutil.ReadFile(path); err != nil {
		return
	}
	if err = conf.UnmarshalJSON(data); err != nil {
		return
	}
	InitEnvironment(conf)
	return
}
