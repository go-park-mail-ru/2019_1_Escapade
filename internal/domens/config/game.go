package config

import "github.com/go-park-mail-ru/2019_1_Escapade/internal/domens"

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

// Room - configutaion of Room(engine.Room)
//easyjson:json
type Room struct {
	CanClose         bool            `json:"canClose"`
	Wait             domens.Duration `json:"wait"`
	Timeouts         GameTimeouts    `json:"timeouts"`
	Field            Field           `json:"field"`
	GarbageCollector domens.Duration `json:"garbage"`
	IDLength         int             `json:"length"`
}

// IDLength 16

// Lobby - configutaion of Lobby(engine.Lobby)
//easyjson:json
type Lobby struct {
	ConnectionsCapacity int32                `json:"connections"`
	RoomsCapacity       int32                `json:"rooms"`
	Intervals           LobbyTimersIntervals `json:"intervals"`
	ConnectionTimeout   domens.Duration      `json:"connection"`
	Wait                domens.Duration      `json:"wait"`
}

// Field - configutaion of Field(engine.Field)
//easyjson:json
type Field struct {
	MinAreaSize    int             `json:"minAreaSize"`
	MaxAreaSize    int             `json:"maxAreaSize"`
	MinProbability int             `json:"minProbability"`
	MaxProbability int             `json:"maxProbability"`
	Wait           domens.Duration `json:"wait"`
}

// Anonymous - configutaion of Anonymous users
//easyjson:json
type Anonymous struct {
	MinID int `json:"minID"`
	MaxID int `json:"maxID"`
}

// GameTimeouts - waiting time for a response from users.
//  If it is exceeded the user will be disabled
//easyjson:json
type GameTimeouts struct {
	PeopleFinding   domens.Duration `json:"peopleFinding"`
	RunningPlayer   domens.Duration `json:"runningPlayer"`
	RunningObserver domens.Duration `json:"runningObserver"`
	Finished        domens.Duration `json:"finished"`
}

// LobbyTimersIntervals intervals of launching regular actions
//  in lobby
//easyjson:json
type LobbyTimersIntervals struct {
	GarbageCollector domens.Duration `json:"garbage"`
	MessagesToDB     domens.Duration `json:"messages"`
	GamesToDB        domens.Duration `json:"games"`
}
