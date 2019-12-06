package engine

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
	action_ "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine/action"
	room_ "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine/room"
)

type RoomInLobby struct {
	room  *Room
	lobby *Lobby
}

// RoomStart - room remove from free
func (lobby *Lobby) roomStart(room *Room, roomID string) {
	lobby.s.DoWithOther(room, func() {
		lobby.removeFromFreeRooms(roomID)
		lobby.sendRoomUpdate(room, All)
		room.people.ForEach(func(c *Connection, isPlayer bool) {
			lobby.waiterToPlayer(c, room)
		})
	})
}

// roomFinish - room remove from all
func (lobby *Lobby) roomFinish(roomID string) {
	lobby.s.Do(func() {
		lobby.removeFromAllRooms(roomID)
		lobby.sendRoomDelete(roomID, All)
	})
}

// CloseRoom free room resources
func (lobby *Lobby) roomClose(roomID string) {
	lobby.s.Do(func() {
		lobby.removeFromFreeRooms(roomID)
		lobby.removeFromAllRooms(roomID)
		lobby.sendRoomDelete(roomID, All)
	})
}

// CreateAndAddToRoom create room and add player to it
func (lobby *Lobby) CreateAndAddToRoom(rs *models.RoomSettings, conn *Connection) (*Room, error) {
	var (
		room *Room
		err  error
	)
	lobby.s.DoWithOther(conn, func() {
		room, err = lobby.createRoom(rs)
		if err == nil {
			utils.Debug(false, "We create your own room, cool!", conn.ID())
			room.people.Enter(conn, true, false)
		} else {
			utils.Debug(true, "cant create. Why?", conn.ID(), err.Error())
		}
	})
	return room, err
}

// createRoom create room, add to all and free rooms
// and run it
func (lobby *Lobby) createRoom(rs *models.RoomSettings) (*Room, error) {
	var (
		room *Room
		err  error
	)
	lobby.s.Do(func() {
		id := utils.RandomString(lobby.rconfig().IDLength)
		game := &models.Game{Settings: rs}
		room, err = NewRoom(lobby.rconfig(), lobby, game, id)
		if err != nil {
			return
		}
		if err = lobby.addRoom(room); err != nil {
			return
		}
	})
	return room, err
}

func (lobby *Lobby) subscribeToRoom(room *Room) {
	var ril = &RoomInLobby{
		room:  room,
		lobby: lobby,
	}
	ril.eventsSubscribe(room.events)
	ril.connectionSubscribe(room.client)
}

// LoadRooms load rooms from database
func (lobby *Lobby) LoadRooms(URLs []string) error {
	var err = re.ErrorLobbyDone()
	lobby.s.Do(func() {
		err = nil
		for _, URL := range URLs {
			room, err := lobby.Load(URL)
			if err != nil {
				return
			}
			if err = lobby.addRoom(room); err != nil {
				return
			}
		}
	})
	return err
}

// addRoom add room to slice of all and free lobby rooms
func (lobby *Lobby) addRoom(room *Room) error {
	var err = re.ErrorLobbyDone()
	lobby.s.Do(func() {
		if err = lobby.addToAllRooms(room); err != nil {
			return
		}
		if err = lobby.addToFreeRooms(room); err != nil {
			return
		}
		lobby.sendRoomCreate(room, All) // inform all about new room
	})
	return err
}

/////////////////////// subsribe

// eventsSubscribe subscibe to events associated with room's status
func (sub *RoomInLobby) eventsSubscribe(e EventsI) {
	observer := synced.NewObserver(
		synced.NewPairNoArgs(room_.StatusFlagPlacing, func() {
			sub.lobby.roomStart(sub.room, sub.room.info.ID())
		}),
		synced.NewPairNoArgs(room_.StatusFinished, func() {
			sub.lobby.roomFinish(sub.room.info.ID())
		}),
		synced.NewPairNoArgs(room_.StatusAborted, func() {
			sub.lobby.roomClose(sub.room.info.ID())
		}))
	e.Observe(observer.AddPublisherCode(room_.UpdateStatus))
}

// connectionBackToLobby is called when user came to lobby from room
func (sub *RoomInLobby) connectionBackToLobby(msg synced.Msg) {
	action, ok := msg.Extra.(ConnectionMsg)
	if !ok {
		return
	}
	sub.lobby.LeaveRoom(action.connection, msg.Action)
}

// connectionSubscribe subscibe to events associated with connection's events
func (sub *RoomInLobby) connectionSubscribe(c RClientI) {
	observer := synced.NewObserver(
		synced.NewPair(action_.BackToLobby, sub.connectionBackToLobby))
	c.Observe(observer.AddPublisherCode(room_.UpdateConnection).
		AddPreAction(func() {
			sub.lobby.sendRoomUpdate(sub.room, All)
		}))
}

// 109 -> 154
