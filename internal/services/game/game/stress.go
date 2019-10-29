package game

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

func (lobby *Lobby) stress(n int) {

	for i := 0; i < n; i++ {
		rs := models.NewBigRoom()
		lobby.createRoom(rs)
	}

	it := NewRoomsIterator(lobby.freeRooms)
	for it.Next() {
		room := it.Value()
		var cells []Cell
		room.Field.OpenEverything(&cells)
	}
}
