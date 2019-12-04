package engine

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
	action_ "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine/action"
	room_ "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine/room"
)

type SubscriberRoom struct {
	synced.SubscriberBase
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
			room.people.add(conn, true, false)
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

/////////////////////// callbacks

func (lobby *Lobby) EventsSub(r *Room) synced.SubscriberI {
	var sub = SubscriberRoom{
		lobby: lobby,
		room:  r,
	}
	sub.SubscriberBase = synced.NewSubscriber(sub.eventsCallback)
	return sub
}

func (sub *SubscriberRoom) eventsCallback(msg synced.Msg) {
	sub.lobby.s.DoWithOther(sub.room, func() {
		if msg.Code != room_.UpdateStatus {
			return
		}
		code, ok := msg.Content.(int)
		if !ok {
			return
		}
		sub.lobby.sendRoomUpdate(sub.room, All)
		switch code {
		case room_.StatusFlagPlacing:
			sub.lobby.roomStart(sub.room, sub.room.info.ID())
		case room_.StatusFinished:
			sub.lobby.roomFinish(sub.room.info.ID())
		case room_.StatusAborted:
			sub.lobby.roomClose(sub.room.info.ID())
		}
	})
}

func (lobby *Lobby) ConnectionSub(room *Room) synced.SubscriberI {
	var sub = SubscriberRoom{
		lobby: lobby,
		room:  room,
	}
	sub.SubscriberBase = synced.NewSubscriber(sub.connectionCallback)
	return sub
}

func (sub *SubscriberRoom) connectionCallback(msg synced.Msg) {
	if msg.Code != room_.UpdateConnection {
		return
	}
	action, ok := msg.Content.(ConnectionMsg)
	if !ok {
		return
	}
	sub.lobby.sendRoomUpdate(sub.room, All)
	switch action.code {
	case action_.BackToLobby:
		sub.lobby.LeaveRoom(action.connection, action.code)
	}
}

// 109
