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

	RoomPlayers             *prometheus.HistogramVec
	RoomDifficult           *prometheus.HistogramVec
	RoomSize                *prometheus.HistogramVec
	RoomTime                *prometheus.HistogramVec
	RoomOpenProcent         prometheus.Histogram
	RoomMode                *prometheus.GaugeVec
	RoomAnonymous           *prometheus.GaugeVec
	RoomTimeSearchingPeople *prometheus.HistogramVec
	RoomTimePlaying         prometheus.Histogram

	Online          prometheus.Gauge
	AnonymousOnline prometheus.Gauge

	Visits prometheus.Counter

	InLobby prometheus.Gauge
	InGame  prometheus.Gauge

	LobbyMessages prometheus.Gauge
	RoomsMessages prometheus.Gauge

	RoomsReconnections prometheus.Counter
)

// Init prometheus metrics variables
func Init() {

	var (
		subsystem      = "game"
		nFinishedGames = "finished_games"
		nAllRooms      = "all_games"
		nLobby         = "lobby"
		nUsers         = "users"
	)

	// все числа в конфиг

	// Lobby characteristics
	FinishedRooms = initFinishedRooms(nLobby, subsystem)
	AbortedRooms = initAbortedRooms(nLobby, subsystem)
	ActiveRooms = initActiveRooms(nLobby, subsystem)
	RecruitmentRooms = initRecruitmentRooms(nLobby, subsystem)
	LobbyMessages = initLobbyMessages(nLobby, subsystem)

	// All rooms characteristics
	RoomsReconnections = initRoomsReconnections(nAllRooms, subsystem)
	RoomsMessages = initRoomsMessages(nAllRooms, subsystem)

	// Finished rooms characteristics
	RoomPlayers = initRoomPlayers(nFinishedGames, subsystem)
	RoomDifficult = initRoomDifficult(nFinishedGames, subsystem)
	RoomSize = initRoomSize(nFinishedGames, subsystem)
	RoomTime = initRoomTime(nFinishedGames, subsystem)
	RoomOpenProcent = initRoomOpenProcent(nFinishedGames, subsystem)
	RoomMode = initRoomMode(nFinishedGames, subsystem)
	RoomAnonymous = initRoomAnonymous(nFinishedGames, subsystem)
	RoomTimeSearchingPeople = initRoomTimeSearchingPeople(nFinishedGames, subsystem)
	RoomTimePlaying = initRoomTimePlaying(nFinishedGames, subsystem)

	// Users
	Online = initOnline(nUsers, subsystem)
	AnonymousOnline = initAnonymousOnline(nUsers, subsystem)
	Visits = initVisits(nUsers, subsystem)
	InLobby = initInLobby(nUsers, subsystem)
	InGame = initInGame(nUsers, subsystem)

	prometheus.MustRegister(FinishedRooms, AbortedRooms, ActiveRooms,
		RoomPlayers, RoomDifficult, RoomSize, RoomTime, RoomOpenProcent,
		RoomMode, RoomAnonymous, RoomTimeSearchingPeople, RoomTimePlaying,
		Online, AnonymousOnline, Visits, InLobby, InGame, LobbyMessages,
		RoomsMessages, RoomsReconnections)
}

func initFinishedRooms(nLobby, subsystem string) prometheus.Counter {

	return prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "FinishedRooms",
		Namespace: nLobby,
		Subsystem: subsystem,
		Help:      "Number of successfully completed games",
	})
}

func initAbortedRooms(nLobby, subsystem string) prometheus.Counter {
	return prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "AbortedRooms",
		Namespace: nLobby,
		Subsystem: subsystem,
		Help:      "Number of aborted games",
	})
}

func initActiveRooms(nLobby, subsystem string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "ActiveRooms",
		Namespace: nLobby,
		Subsystem: subsystem,
		Help:      "Number of playing rooms",
	})
}

func initRecruitmentRooms(nLobby, subsystem string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "RecruitmentRooms",
		Namespace: nLobby,
		Subsystem: subsystem,
		Help:      "Number of recruiting rooms",
	})
}

func initLobbyMessages(nLobby, subsystem string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "Messages",
		Namespace: nLobby,
		Subsystem: subsystem,
		Help:      "Number of sent messages",
	})
}

func initRoomsReconnections(nAllRooms, subsystem string) prometheus.Counter {
	return prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "Reconnections",
		Namespace: nAllRooms,
		Subsystem: subsystem,
		Help:      "Number of reconnections in game",
	})
}

func initRoomsMessages(nAllRooms, subsystem string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "Messages",
		Namespace: nAllRooms,
		Subsystem: subsystem,
		Help:      "Number of sent messages",
	})
}

func initRoomPlayers(nFinishedGames, subsystem string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "Players",
		Namespace: nFinishedGames,
		Subsystem: subsystem,
		Buckets:   prometheus.ExponentialBuckets(2, 2, 7),
		Help:      "Number of players who played the game",
	}, []string{"room_type"})
}

func initRoomDifficult(nFinishedGames, subsystem string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "Difficult",
		Namespace: nFinishedGames,
		Subsystem: subsystem,
		Buckets:   []float64{0.2, 0.5, 0.9},
		Help:      "Complexity of the game",
	}, []string{"room_type"})
}
func initRoomSize(nFinishedGames, subsystem string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "Size",
		Namespace: nFinishedGames,
		Subsystem: subsystem,
		Buckets:   prometheus.ExponentialBuckets(100, 4, 5),
		Help:      "Size of the field",
	}, []string{"room_type"})
}

func initRoomTime(nFinishedGames, subsystem string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "Time",
		Namespace: nFinishedGames,
		Subsystem: subsystem,
		Buckets:   prometheus.ExponentialBuckets(1, 10, 5),
		Help:      "The most time allotted for the game",
	}, []string{"room_type"})
}

func initRoomOpenProcent(nFinishedGames, subsystem string) prometheus.Histogram {
	return prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:      "Procent",
		Namespace: nFinishedGames,
		Subsystem: subsystem,
		Buckets:   []float64{0.3, 0.5, 0.8, 0.9, 1},
		Help:      "The percentage opening of the field",
	})
}

func initRoomMode(nFinishedGames, subsystem string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "Mode",
		Namespace: nFinishedGames,
		Subsystem: subsystem,
		Help:      "Deathmatch or not",
	}, []string{"room_type", "deathmatch"})
}

func initRoomAnonymous(nFinishedGames, subsystem string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "Anonymous",
		Namespace: nFinishedGames,
		Subsystem: subsystem,
		Help:      "Anonymous disable[1]/anonymous enable(and they are in game)[2]//anonymous enable(but they are not in game)[3]",
	}, []string{"room_type", "anonymous"})
}

func initRoomTimeSearchingPeople(nFinishedGames, subsystem string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "TimeSearchingPeople",
		Namespace: nFinishedGames,
		Subsystem: subsystem,
		Buckets:   prometheus.ExponentialBuckets(1, 10, 5),
		Help:      "Time spent recruiting people",
	}, []string{"room_type"})
}

func initRoomTimePlaying(nFinishedGames, subsystem string) prometheus.Histogram {
	return prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:      "TimePlaying",
		Namespace: nFinishedGames,
		Subsystem: subsystem,
		Buckets:   prometheus.ExponentialBuckets(1, 10, 5),
		Help:      "Time spent playing",
	})
}

func initOnline(nUsers, subsystem string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "Online",
		Namespace: nUsers,
		Subsystem: subsystem,
		Help:      "Users online at one moment",
	})
}

func initAnonymousOnline(nUsers, subsystem string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "AnonymousOnline",
		Namespace: nUsers,
		Subsystem: subsystem,
		Help:      "Anonymous users online at one moment",
	})
}

func initVisits(nUsers, subsystem string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "Visits",
		Namespace: nUsers,
		Subsystem: subsystem,
		Help:      "Number of visits",
	})
}

func initInLobby(nUsers, subsystem string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "InLobby",
		Namespace: nUsers,
		Subsystem: subsystem,
		Help:      "Users in lobby at one moment",
	})
}

func initInGame(nUsers, subsystem string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "InGame",
		Namespace: nUsers,
		Subsystem: subsystem,
		Help:      "Users in games at one moment",
	})
}
