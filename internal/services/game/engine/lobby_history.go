package engine

import (
	"github.com/gorilla/websocket"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// LaunchLobbyHistory launch local lobby with rooms from database
func LaunchLobbyHistory(db database.GameUseCaseI,
	ws *websocket.Conn, user *models.UserPublicInfo,
	cw config.WebSocket, gameSettings *config.Game,
	si SetImage) {

	urls, err := db.FetchAllRoomsID(user.ID)

	if err != nil {
		utils.Debug(false, "GetGamesURL", err.Error())
		return
	}

	gameSettings.RoomsCapacity = int32(len(urls) * 2)
	lobby := NewLobby(gameSettings, db, si)

	go lobby.Run()
	defer func() {
		lobby.Stop()
	}()

	if len(urls) > 0 {
		err = lobby.LoadRooms(urls)
		if err != nil {
			utils.Debug(false, "LoadRooms", err.Error())
			return
		}
	}

	conn := NewConnection(ws, user, lobby)
	conn.Launch(cw, "")
}
