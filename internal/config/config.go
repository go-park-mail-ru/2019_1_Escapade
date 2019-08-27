package config

import (
	"io/ioutil"
	"os"
	"time"

	"golang.org/x/oauth2"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// Configuration contains all types of configurations
//easyjson:json
type Configuration struct {
	Server     Server     `json:"server"`
	Cors       CORS       `json:"cors"`
	DataBase   Database   `json:"dataBase"`
	Game       Game       `json:"game"`
	Cookie     Cookie     `json:"cookie"`
	WebSocket  WebSocket  `json:"websocket"`
	Service    Service    `json:"service"`
	Auth       Auth       `json:"auth"`
	AuthClient AuthClient `json:"authClient"`
}

// ServerConfig set host, post and buffers sizes
//easyjson:json
type Server struct {
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
type CORS struct {
	Origins     []string `json:"origins"`
	Headers     []string `json:"headers"`
	Methods     []string `json:"methods"`
	Credentials string   `json:"credentials"`
}

// DatabaseConfig set type of database management system
//   the url of connection string, max amount of
//   connections, tables, sizes of page  of gamers
//   and users
//easyjson:json
type Database struct {
	DriverName           string `json:"driverName"`
	URL                  string `json:"url"`
	ConnectionString     string `json:"connectionString"`
	AuthConnectionString string `json:"authConnectionString"`
	MaxOpenConns         int    `json:"maxOpenConns"`
	PageGames            int    `json:"pageGames"`
	PageUsers            int    `json:"pageUsers"`
}

//easyjson:json
type Field struct {
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
type Game struct {
	RoomsCapacity      int32  `json:"roomsCapacity"`
	ConnectionCapacity int32  `json:"connectionCapacity"`
	Location           string `json:"location"`
	CanClose           bool   `json:"closeRoom"`
	Metrics            bool   `json:"metrics"`
	Field              *Field `json:"field"`
}

// AuthClient client of auth microservice
//easyjson:json
type Auth struct {
	Salt                    string       `json:"salt"`
	AccessTokenExpireHours  int          `json:"accessTokenExpireHours"`
	RefreshTokenExpireHours int          `json:"refreshTokenExpireHours"`
	IsGenerateRefresh       bool         `json:"isGenerateRefresh"`
	WithReserve             bool         `json:"withReserve"`
	TokenType               string       `json:"tokenType"`
	WhiteList               []AuthClient `json:"whiteList"`
}

type AuthClient struct {
	// address of auth service
	Address      string        `json:"address"`
	ClientID     string        `json:"id"`
	ClientSecret string        `json:"secret"`
	Scopes       []string      `json:"scopes"`
	RedirectURL  string        `json:"redirectURL"`
	Config       oauth2.Config `json:"-"`
}

//easyjson:json
type Service struct {
	ConsulID  string   `json:"-"`
	Name      string   `json:"name"`
	DependsOn []string `json:"dependsOn"`
}

// SessionConfig set cookie name, path, length, expiration time
// and HTTPonly flag
//easyjson:json
type Cookie struct {
	Path          string     `json:"path"`
	Length        int        `json:"length"`
	LifetimeHours int        `json:"lifetime"`
	HTTPOnly      bool       `json:"httpOnly"`
	Auth          AuthCookie `json:"authCookie"`
}

//easyjson:json
type AuthCookie struct {
	AccessToken   string `json:"accessToken"`
	TokenType     string `json:"rokenType"`
	RefreshToken  string `json:"refreshToken"`
	Expire        string `json:"expire"`
	ReservePrefix string `json:"reservePrefix"`
}

// WebSocketConfig set timeouts
//easyjson:json
type WebSocket struct {
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
	conf.AuthClient.Config = oauth2.Config{
		ClientID:     conf.AuthClient.ClientID,
		ClientSecret: conf.AuthClient.ClientSecret,
		Scopes:       conf.AuthClient.Scopes,
		RedirectURL:  conf.AuthClient.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  conf.AuthClient.Address + "/auth/authorize",
			TokenURL: conf.AuthClient.Address + "/auth/token",
		},
	}
	InitEnvironment(conf)
	return
}
