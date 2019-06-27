package metrics

import "github.com/prometheus/client_golang/prometheus"

// game vars
var (
	// FinishedRooms - rooms, that finished
	FinishedRooms prometheus.Counter
	// AbortedRooms - rooms, that were deleted before launch
	AbortedRooms prometheus.Counter
	// ActiveRooms - rooms, that are active now
	ActiveRooms prometheus.Gauge
	// RecruitmentRooms - rooms, that are recruiting
	RecruitmentRooms prometheus.Gauge

	FinishedRoomPlayers             prometheus.Histogram
	FinishedRoomDifficult           prometheus.Histogram
	FinishedRoomSize                prometheus.Histogram
	FinishedRoomTime                prometheus.Histogram
	FinishedRoomOpenProcent         prometheus.Histogram
	FinishedRoomMode                prometheus.Histogram
	FinishedRoomAnonymous           prometheus.Histogram
	FinishedRoomTimeSearchingPeople prometheus.Histogram
	FinishedRoomTimePlaying         prometheus.Histogram

	AbortedRoomTimeSearchingPeople prometheus.Histogram

	Online          prometheus.Gauge
	AnonymousOnline prometheus.Gauge

	Visits prometheus.Counter

	InLobby prometheus.Gauge
	InGame  prometheus.Gauge

	LobbyMessages prometheus.Gauge
	RoomsMessages prometheus.Gauge

	RoomsReconnections prometheus.Counter
)

// api vars
var (
	// Hits hits
	Hits *prometheus.CounterVec
	// Users - registered users
	Users *prometheus.GaugeVec
)

func InitApi() {
	var (
		subsystem = "api"
	)
	Hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "FinishedRooms",
		Subsystem: subsystem,
	}, []string{"status", "path", "method"})

	prometheus.MustRegister(Hits)
}

func InitGame() {

	var (
		subsystem      = "game"
		nFinishedGames = "finished_games"
		nAbortedGames  = "aborted_games"
		nAllRooms      = "all_games"
		nLobby         = "lobby"
		nUsers         = "users"
	)

	// Lobby characteristics
	FinishedRooms = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "FinishedRooms",
		Namespace: nLobby,
		Subsystem: subsystem,
		Help:      "Number of successfully completed games",
	})
	AbortedRooms = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "AbortedRooms",
		Namespace: nLobby,
		Subsystem: subsystem,
		Help:      "Number of aborted games",
	})
	ActiveRooms = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "ActiveRooms",
		Namespace: nLobby,
		Subsystem: subsystem,
		Help:      "Number of playing rooms",
	})
	RecruitmentRooms = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "RecruitmentRooms",
		Namespace: nLobby,
		Subsystem: subsystem,
		Help:      "Number of recruiting rooms",
	})
	LobbyMessages = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "Messages",
		Namespace: nLobby,
		Subsystem: subsystem,
		Help:      "Number of sent messages",
	})

	// Lobby characteristics
	RoomsReconnections = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "Reconnections",
		Namespace: nAllRooms,
		Subsystem: subsystem,
		Help:      "Number of reconnections in game",
	})

	RoomsMessages = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "Messages",
		Namespace: nAllRooms,
		Subsystem: subsystem,
		Help:      "Number of sent messages",
	})

	// Finished rooms characteristics
	FinishedRoomPlayers = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:      "Players",
		Namespace: nFinishedGames,
		Subsystem: subsystem,
		Help:      "Number of players who played the game",
	})
	FinishedRoomDifficult = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:      "Difficult",
		Namespace: nFinishedGames,
		Subsystem: subsystem,
		Help:      "Complexity of the game",
	})
	FinishedRoomSize = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:      "Size",
		Namespace: nFinishedGames,
		Subsystem: subsystem,
		Help:      "Size of the game",
	})
	FinishedRoomTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:      "Time",
		Namespace: nFinishedGames,
		Subsystem: subsystem,
		Help:      "The most time allotted for the game",
	})
	FinishedRoomOpenProcent = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:      "Procent",
		Namespace: nFinishedGames,
		Subsystem: subsystem,
		Help:      "The percentage opening of the field",
	})
	FinishedRoomMode = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:      "Mode",
		Namespace: nFinishedGames,
		Subsystem: subsystem,
		Help:      "Deathmatch or not",
	})
	FinishedRoomAnonymous = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:      "Anonymous",
		Namespace: nFinishedGames,
		Subsystem: subsystem,
		Help:      "Anonymous disable[AD]/anonymous enable(and they are in game)[AEY]//anonymous enable(but they are not in game)[AEN]",
	})
	FinishedRoomTimeSearchingPeople = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:      "TimeSearchingPeople",
		Namespace: nFinishedGames,
		Subsystem: subsystem,
		Help:      "Time spent recruiting people",
	})
	FinishedRoomTimePlaying = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:      "TimePlaying",
		Namespace: nFinishedGames,
		Subsystem: subsystem,
		Help:      "Time spent playing",
	})

	// Aborted rooms characteristics
	AbortedRoomTimeSearchingPeople = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:      "TimePlaying",
		Namespace: nAbortedGames,
		Subsystem: subsystem,
		Help:      "Time spent recruiting people",
	})

	// Users
	Online = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "Online",
		Namespace: nUsers,
		Subsystem: subsystem,
		Help:      "Users online at one moment",
	})
	AnonymousOnline = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "AnonymousOnline",
		Namespace: nUsers,
		Subsystem: subsystem,
		Help:      "Anonymous users online at one moment",
	})
	Visits = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "Visits",
		Namespace: nUsers,
		Subsystem: subsystem,
		Help:      "Number of visits",
	})
	InLobby = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "InLobby",
		Namespace: nUsers,
		Subsystem: subsystem,
		Help:      "Users in lobby at one moment",
	})
	InGame = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "InGame",
		Namespace: nUsers,
		Subsystem: subsystem,
		Help:      "Users in games at one moment",
	})

	prometheus.MustRegister(FinishedRooms, AbortedRooms, ActiveRooms,
		FinishedRoomPlayers, FinishedRoomDifficult, FinishedRoomSize,
		FinishedRoomTime, FinishedRoomOpenProcent, FinishedRoomMode,
		FinishedRoomAnonymous, FinishedRoomTimeSearchingPeople,
		FinishedRoomTimePlaying, AbortedRoomTimeSearchingPeople,
		Online, AnonymousOnline, Visits, InLobby, InGame, LobbyMessages,
		RoomsMessages, RoomsReconnections)
}
