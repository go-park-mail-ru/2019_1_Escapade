package engine

import (
	"github.com/gorilla/websocket"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/clients"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/database"
)

// LaunchLobbyHistory launch local lobby with rooms from database
func LaunchLobbyHistory(chatS clients.Chat, db database.GameUseCaseI,
	ws *websocket.Conn, user *models.UserPublicInfo,
	cw config.WebSocket, gameSettings *config.Game,
	si SetImage) {

	urls, err := db.FetchAllRoomsID(user.ID)

	if err != nil {
		utils.Debug(false, "GetGamesURL", err.Error())
		return
	}

	gameSettings.Lobby.RoomsCapacity = int32(len(urls) * 2)
	lobby := NewLobby(chatS, gameSettings, db, si)

	go lobby.Run()
	defer func() {
		lobby.Close()
	}()

	if len(urls) > 0 {
		err = lobby.LoadRooms(urls)
		if err != nil {
			utils.Debug(false, "LoadRooms", err.Error())
			return
		}
	}

	conn, err := NewConnection(ws, user, lobby)
	if err != nil {
		utils.Debug(false, "cant create connection")
		return
	}
	conn.Launch(cw, "")
}
