package game

import (
	"sync"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/metrics"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

func (lobby *Lobby) removeFromFreeRooms(roomID string, group *sync.WaitGroup) {
	defer group.Done()
	defer utils.CatchPanic("lobby_metrics.go removeFromFreeRooms")

	if lobby.done() {
		return
	}

	lobby.wGroup.Add(1)
	defer lobby.wGroup.Done()

	if lobby.freeRooms.Remove(roomID) && lobby.config.Metrics {
		metrics.RecruitmentRooms.Dec()
	}
}

func (lobby *Lobby) removeFromAllRooms(roomID string, group *sync.WaitGroup) {
	defer group.Done()
	defer utils.CatchPanic("lobby_metrics.go removeFromAllRooms")

	if lobby.done() {
		return
	}

	lobby.wGroup.Add(1)
	defer lobby.wGroup.Done()

	if lobby.allRooms.Remove(roomID) && lobby.config.Metrics {
		metrics.ActiveRooms.Dec()
	}
}

func (lobby *Lobby) addToFreeRooms(room *Room, group *sync.WaitGroup) error {
	defer group.Done()
	defer utils.CatchPanic("lobby_metrics.go removeFromAllRooms")

	if lobby.done() || room.done() {
		return re.ErrorRoomOrLobbyDone()
	}

	lobby.wGroup.Add(1)
	defer lobby.wGroup.Done()

	room.wGroup.Add(1)
	defer room.wGroup.Done()

	if lobby.freeRooms.Add(room) {
		if lobby.config.Metrics {
			metrics.RecruitmentRooms.Inc()
		}
	} else {
		return re.ErrorLobbyCantCreateRoom()
	}
	return nil
}

func (lobby *Lobby) addToAllRooms(room *Room, group *sync.WaitGroup) error {
	defer group.Done()
	defer utils.CatchPanic("lobby_metrics.go removeFromAllRooms")

	if lobby.done() || room.done() {
		return re.ErrorRoomOrLobbyDone()
	}

	lobby.wGroup.Add(1)
	defer lobby.wGroup.Done()

	room.wGroup.Add(1)
	defer room.wGroup.Done()

	if lobby.allRooms.Add(room) {
		if lobby.config.Metrics {
			metrics.ActiveRooms.Inc()
		}
	} else {
		return re.ErrorLobbyCantCreateRoom()
	}
	return nil
}
