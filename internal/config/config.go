package config

import (
	json "encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"golang.org/x/oauth2"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
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

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' {
		sd := string(b[1 : len(b)-1])
		d.Duration, err = time.ParseDuration(sd)
		return
	}

	var id int64
	id, err = json.Number(string(b)).Int64()
	d.Duration = time.Duration(id)

	return
}

func (d Duration) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}

// Server set host, post and buffers sizes
//easyjson:json
type Server struct {
	Name           string   `json:"name"`
	MaxConn        int      `json:"maxConn"`
	MaxHeaderBytes int      `json:"maxHeaderBytes"`
	Timeouts       Timeouts `json:"timeouts"`
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
type Timeouts struct {
	TTL Duration `json:"ttl"`

	Read  Duration `json:"read"`
	Write Duration `json:"write"`
	Idle  Duration `json:"idle"`
	Wait  Duration `json:"wait"`
	Exec  Duration `json:"exec"`
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
	MinAreaSize    int      `json:"minAreaSize"`
	MaxAreaSize    int      `json:"maxAreaSize"`
	MinProbability int      `json:"minProbability"`
	MaxProbability int      `json:"maxProbability"`
	Wait           Duration `json:"wait"`
}

//groupWaitRoom := 60 * time.Second // TODO в конфиг

//easyjson:json
type Room struct {
	CanClose         bool         `json:"canClose"`
	Wait             Duration     `json:"wait"`
	Timeouts         GameTimeouts `json:"timeouts"`
	Field            Field        `json:"field"`
	GarbageCollector Duration     `json:"garbage"`
	IDLength         int          `json:"length"`
}

// IDLength 16

//easyjson:json
type Anonymous struct {
	MinID int `json:"minID"`
	MaxID int `json:"maxID"`
}

// Timeouts of the connection to the server
type GameTimeouts struct {
	PeopleFinding   Duration `json:"peopleFinding"`
	RunningPlayer   Duration `json:"runningPlayer"`
	RunningObserver Duration `json:"runningObserver"`
	Finished        Duration `json:"finished"`
}

type LobbyTimersIntervals struct {
	GarbageCollector Duration `json:"garbage"`
	MessagesToDB     Duration `json:"messages"`
	GamesToDB        Duration `json:"games"`
}

type Lobby struct {
	ConnectionsCapacity int32                `json:"connections"`
	RoomsCapacity       int32                `json:"rooms"`
	Intervals           LobbyTimersIntervals `json:"intervals"`
	ConnectionTimeout   Duration             `json:"connection"`
	Wait                Duration             `json:"wait"`
}

// conectionTimeout = 10s

// Game set, how much rooms server can create and
// how much connections can join. Also there are flags:
// can server close rooms or not(for history mode),
// metrics should be recorded or not
//easyjson:json
type Game struct {
	Lobby     Lobby     `json:"lobby"`
	Room      Room      `json:"room"`
	Anonymous Anonymous `json:"anonymous"`
	Location  string    `json:"location"`
	Metrics   bool      `json:"metrics"`
}

// groupWaitTimeout := 80 * time.Second // TODO в конфиг

// Auth client of auth microservice
//easyjson:json
type Auth struct {
	Salt               string       `json:"salt"`
	AccessTokenExpire  Duration     `json:"accessTokenExpire"`
	RefreshTokenExpire Duration     `json:"refreshTokenExpire"`
	IsGenerateRefresh  bool         `json:"isGenerateRefresh"`
	WithReserve        bool         `json:"withReserve"`
	TokenType          string       `json:"tokenType"`
	WhiteList          []AuthClient `json:"whiteList"`
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

// Cookie set cookie name, path, length, expiration time
// and HTTPonly flag
//easyjson:json
type Cookie struct {
	Path          string     `json:"path"`
	LifetimeHours int        `json:"lifetime_hours"`
	HTTPOnly      bool       `json:"httpOnly"`
	Auth          AuthCookie `json:"keys"`
}

//easyjson:json
type AuthCookie struct {
	AccessToken   string `json:"accessToken"`
	TokenType     string `json:"tokenType"`
	RefreshToken  string `json:"refreshToken"`
	Expire        string `json:"expire"`
	ReservePrefix string `json:"reservePrefix"`
}

// WebSocket set timeouts
//easyjson:json
type WebSocket struct {
	WriteWait       Duration `json:"writeWait"`
	PongWait        Duration `json:"pongWait"`
	PingPeriod      Duration `json:"pingPeriod"`
	MaxMessageSize  int64    `json:"maxMessageSize"`
	ReadBufferSize  int      `json:"readBufferSize"`
	WriteBufferSize int      `json:"writeBufferSize"`
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
	utils.Debug(false, " Look at config", conf.Auth.Salt, conf.Auth.AccessTokenExpire,
		conf.Auth.RefreshTokenExpire, conf.Auth.IsGenerateRefresh, conf.Auth.WithReserve,
		conf.Auth.TokenType, conf.Auth.WhiteList)
	utils.Debug(false, " Info:", conf.Server.Name, conf.Server.MaxHeaderBytes)
	utils.Debug(false, " Timeouts:", conf.Server.Timeouts.TTL,
		conf.Server.Timeouts.Read, conf.Server.Timeouts.Write, conf.Server.Timeouts.Idle,
		conf.Server.Timeouts.Wait, conf.Server.Timeouts.Exec)

	conf.setOauth2Config()
	conf.AuthClient.Address = os.Getenv("AUTH_ADDRESS")
	//InitEnvironment(conf)
	return
}

func (conf *Configuration) setOauth2Config() {
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
}
