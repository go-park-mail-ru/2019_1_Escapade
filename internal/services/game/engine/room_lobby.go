package engine

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// LobbyProxyI control access to lobby
// Proxy Pattern
type LobbyProxyI interface {
	Start()
	Finish()
	Close()
	Notify()
	SaveGame(info models.GameInformation) error

	Greet(conn *Connection)
	BackToLobby(conn *Connection)
	WaiterToPlayer(conn *Connection)
	CreateAndAddToRoom(conn *Connection) (*Room, error)

	metricsEnabled() bool
	closeEnabled() bool

	SaveMessages(mwa *MessageWithAction)

	setWaitingRoom(conn *Connection)

	Date() time.Time
}

// RoomLobby implements LobbyProxyI
type RoomLobby struct {
	r     *Room
	s     SyncI
	i     RoomInformationI
	lobby *Lobby

	needMetrics bool
	canClose    bool
}

// Init configure dependencies with other components of the room
func (room *RoomLobby) Init(builder ComponentBuilderI, r *Room, lobby *Lobby) {
	room.lobby = lobby
	room.r = r
	builder.BuildSync(&room.s)
	builder.BuildInformation(&room.i)

	room.needMetrics = room.lobby.config().Metrics
	room.canClose = room.lobby.config().CanClose
}

func (room *RoomLobby) Finish() {
	room.s.do(func() {
		room.lobby.roomFinish(room.i.ID())
	})
}

func (room *RoomLobby) Notify() {
	room.s.do(func() {
		room.lobby.sendRoomUpdate(room.r, All)
	})
}

func (room *RoomLobby) Start() {
	room.s.do(func() {
		room.lobby.RoomStart(room.r, room.i.ID())
	})
}

func (room *RoomLobby) SaveGame(info models.GameInformation) error {
	err := room.lobby.db().Save(info)
	if err != nil {
		room.lobby.AddNotSavedGame(&info)
	}
	return err
}

func (room *RoomLobby) Close() {
	room.s.do(func() {
		room.lobby.CloseRoom(room.i.ID())
	})
}

func (room *RoomLobby) Greet(conn *Connection) {
	room.lobby.greet(conn)
}

func (room *RoomLobby) BackToLobby(conn *Connection) {
	go room.lobby.LeaveRoom(conn, ActionBackToLobby)
}

func (room *RoomLobby) WaiterToPlayer(conn *Connection) {
	room.s.doWithConn(conn, func() {
		room.lobby.waiterToPlayer(conn, room.r)
	})
}

func (room *RoomLobby) CreateAndAddToRoom(conn *Connection) (*Room, error) {
	var (
		newRoom *Room
		err     error
	)
	room.s.doWithConn(conn, func() {
		newRoom, err = room.lobby.CreateAndAddToRoom(room.i.Settings(), conn)
	})
	return newRoom, err
}

func (room *RoomLobby) Date() time.Time {
	return time.Now().In(room.lobby.location())
}

func (room *RoomLobby) metricsEnabled() bool {
	return room.needMetrics
}

func (room *RoomLobby) closeEnabled() bool {
	return room.canClose
}

func (room *RoomLobby) SaveMessages(mwa *MessageWithAction) {
	room.lobby.AddNotSavedMessage(mwa)
}

func (room *RoomLobby) setWaitingRoom(conn *Connection) {
	conn.setWaitingRoom(room.r)
}

// IsWinner is player wuth id playerID is winner

///////// sent to onlinePlayers
/*
func (room *Room) isWinner(playerIndex int, isMax *bool) func(int, Player) {
	var (
		max              = 0.
		thisPlayerPoints float64
		ignore           bool
	)

	return func(index int, player Player) {
		if ignore || player.Died || player.Points < max {
			return
		}
		max = player.Points
		if index == playerIndex {
			thisPlayerPoints = max
			*isMax = true
		} else if *isMax && max > thisPlayerPoints {
			*isMax = false
			ignore = true
		}
	}
}*/
