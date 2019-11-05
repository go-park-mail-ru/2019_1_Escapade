package engine

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

type RoomLifecycle struct {
	r *Room
}

type RoomLobbyCommunicationI interface {
	Init(r *Room, s SyncI, i *RoomInformation, lobby *Lobby)

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

type RoomLobbyCommunication struct {
	r     *Room
	s     SyncI
	i     *RoomInformation
	lobby *Lobby

	needMetrics bool
	canClose    bool
}

func (room *RoomLobbyCommunication) Init(r *Room, s SyncI,
	i *RoomInformation, lobby *Lobby) {
	room.lobby = lobby
	room.r = r
	room.s = s
	room.i = i

	room.needMetrics = room.lobby.config().Metrics
	room.canClose = room.lobby.config().CanClose
}

func (room *RoomLobbyCommunication) Finish() {
	room.s.do(func() {
		room.lobby.roomFinish(room.i.ID())
	})
}

func (room *RoomLobbyCommunication) Notify() {
	room.s.do(func() {
		room.lobby.sendRoomUpdate(room.r, All)
	})
}

func (room *RoomLobbyCommunication) Start() {
	room.s.do(func() {
		room.lobby.RoomStart(room.r, room.i.ID())
	})
}

func (room *RoomLobbyCommunication) SaveGame(info models.GameInformation) error {
	err := room.lobby.db().Save(info)
	if err != nil {
		room.lobby.AddNotSavedGame(&info)
	}
	return err
}

func (room *RoomLobbyCommunication) Close() {
	room.s.do(func() {
		room.lobby.CloseRoom(room.i.ID())
	})
}

func (room *RoomLobbyCommunication) Greet(conn *Connection) {
	room.lobby.greet(conn)
}

func (room *RoomLobbyCommunication) BackToLobby(conn *Connection) {
	go room.lobby.LeaveRoom(conn, ActionBackToLobby)
}

func (room *RoomLobbyCommunication) WaiterToPlayer(conn *Connection) {
	room.s.doWithConn(conn, func() {
		room.lobby.waiterToPlayer(conn, room.r)
	})
}

func (room *RoomLobbyCommunication) CreateAndAddToRoom(conn *Connection) (*Room, error) {
	var (
		newRoom *Room
		err     error
	)
	room.s.doWithConn(conn, func() {
		newRoom, err = room.lobby.CreateAndAddToRoom(room.i.Settings, conn)
	})
	return newRoom, err
}

func (room *RoomLobbyCommunication) Date() time.Time {
	return time.Now().In(room.lobby.location())
}

func (room *RoomLobbyCommunication) metricsEnabled() bool {
	return room.needMetrics
}

func (room *RoomLobbyCommunication) closeEnabled() bool {
	return room.canClose
}

func (room *RoomLobbyCommunication) SaveMessages(mwa *MessageWithAction) {
	room.lobby.AddNotSavedMessage(mwa)
}

func (room *RoomLobbyCommunication) setWaitingRoom(conn *Connection) {
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
