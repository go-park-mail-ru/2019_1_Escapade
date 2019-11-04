package engine

import (
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/metrics"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

func (lobby *Lobby) removeFromFreeRooms(roomID string) {
	defer utils.CatchPanic("lobby_metrics.go removeFromAllRooms")
	if lobby.done() {
		return
	}

	lobby.wGroup.Add(1)
	defer lobby.wGroup.Done()

	if lobby.freeRooms.Remove(roomID) && lobby.config().Metrics {
		metrics.RecruitmentRooms.Dec()
	}
}

func (lobby *Lobby) removeFromAllRooms(roomID string) {
	defer utils.CatchPanic("lobby_metrics.go removeFromAllRooms")

	if lobby.done() {
		return
	}

	lobby.wGroup.Add(1)
	defer lobby.wGroup.Done()

	if lobby.allRooms.Remove(roomID) && lobby.config().Metrics {
		metrics.ActiveRooms.Dec()
	}
}

func (lobby *Lobby) addRoomToSlice(room *Room, f func() bool) error {
	defer utils.CatchPanic("lobby_metrics.go removeFromAllRooms")
	if lobby.done() {
		return re.ErrorRoomOrLobbyDone()
	}

	lobby.wGroup.Add(1)
	defer lobby.wGroup.Done()

	var (
		done bool
		err  error
	)

	room.sync.do(func() {
		if !f() {
			err = re.ErrorLobbyCantCreateRoom()
			return
		}
		done = true
	})
	if !done {
		err = re.ErrorRoomOrLobbyDone()
	}
	return err
}

func (lobby *Lobby) addToFreeRooms(room *Room) error {
	return lobby.addRoomToSlice(room, func() bool {
		if lobby.freeRooms.Add(room) {
			if lobby.config().Metrics {
				metrics.RecruitmentRooms.Inc()
			}
			return true
		}
		return false
	})
}

func (lobby *Lobby) addToAllRooms(room *Room) error {
	return lobby.addRoomToSlice(room, func() bool {
		if lobby.allRooms.Add(room) {
			if lobby.config().Metrics {
				metrics.ActiveRooms.Inc()
			}
			return true
		}
		return false
	})
}

// m mean metrics

func (lobby *Lobby) mUserWelcome(isAnonymous bool) {
	if lobby.config().Metrics {
		metrics.Online.Inc()
		if isAnonymous {
			metrics.AnonymousOnline.Inc()
		}
	}
}

func (lobby *Lobby) mUserBye(isAnonymous bool) {
	if lobby.config().Metrics {
		metrics.Online.Dec()
		metrics.InLobby.Dec()
		if isAnonymous {
			metrics.AnonymousOnline.Dec()
		}
	}
}

func (lobby *Lobby) removeWaiter(conn *Connection) {
	if !lobby.Waiting.Remove(conn) {
		return
	}
	if lobby.config().Metrics {
		metrics.InLobby.Dec()
	}
	go lobby.sendWaiterExit(conn, All)
}

func (lobby *Lobby) addPlayer(conn *Connection) {
	lobby.Playing.Add(conn)
	if lobby.config().Metrics {
		metrics.InGame.Inc()
	}
	go lobby.sendPlayerEnter(conn, All)
}

func (lobby *Lobby) removePlayer(conn *Connection) {
	if !lobby.Playing.Remove(conn) {
		return
	}
	if lobby.config().Metrics {
		metrics.InGame.Dec()
	}
	go lobby.sendPlayerExit(conn, AllExceptThat(conn))
}

// addWaiter add connection to waiters slice and send to the connection LobbyJSON
func (lobby *Lobby) addWaiter(newConn *Connection) {

	lobby.Waiting.Add(newConn)
	if lobby.config().Metrics {
		metrics.InLobby.Inc()
	}

	lobby.greet(newConn)
	lobby.sendWaiterEnter(newConn, AllExceptThat(newConn))
}
