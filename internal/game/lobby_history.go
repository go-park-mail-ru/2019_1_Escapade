package game

import (
	"github.com/gorilla/websocket"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// LaunchLobbyHistory launch local lobby with rooms from database
func LaunchLobbyHistory(db *database.DataBase,
	ws *websocket.Conn, user *models.UserPublicInfo,
	WSsettings config.WebSocketSettings, gameSettings *config.GameConfig,
	si SetImage) {

	urls, err := db.GetGamesURL(user.ID)

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
	conn.Launch(WSsettings, "")
}
