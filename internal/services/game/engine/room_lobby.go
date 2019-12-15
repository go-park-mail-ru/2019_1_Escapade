package engine

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/clients"
)

// LobbyProxyI control access to lobby
// Proxy Pattern
type LobbyProxyI interface {
	SaveGame(info models.GameInformation) error

	WaiterToPlayer(conn *Connection)
	CreateAndAddToRoom(conn *Connection) (*Room, error)

	SaveMessages(mwa *MessageWithAction)

	setWaitingRoom(conn *Connection)

	config() *config.Room

	ChatService() clients.ChatI
}

// RoomLobby implements LobbyProxyI
type RoomLobby struct {
	r     *Room
	s     synced.SyncI
	i     RoomInformationI
	lobby *Lobby

	needMetrics bool
	canClose    bool
}

// Init configure dependencies with other components of the room
func (room *RoomLobby) Init(builder RBuilderI, r *Room, lobby *Lobby) {
	room.lobby = lobby
	room.r = r
	builder.BuildSync(&room.s)
	builder.BuildInformation(&room.i)

	room.needMetrics = room.lobby.config().Metrics
	room.canClose = room.lobby.config().Room.CanClose
}

func (room *RoomLobby) ChatService() clients.ChatI {
	return room.lobby.ChatService
}

func (room *RoomLobby) SaveGame(info models.GameInformation) error {
	err := room.lobby.db().Save(info)
	if err != nil {
		room.lobby.AddNotSavedGame(&info)
	}
	return err
}

func (room *RoomLobby) WaiterToPlayer(conn *Connection) {
	room.s.DoWithOther(conn, func() {
		room.lobby.waiterToPlayer(conn, room.r)
	})
}

func (room *RoomLobby) CreateAndAddToRoom(conn *Connection) (*Room, error) {
	var (
		newRoom *Room
		err     error
	)
	room.s.DoWithOther(conn, func() {
		newRoom, err = room.lobby.CreateAndAddToRoom(room.i.Settings(), conn)
	})
	return newRoom, err
}

func (room *RoomLobby) config() *config.Room {
	return room.lobby.rconfig()
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
