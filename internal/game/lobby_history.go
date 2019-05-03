package game

import (
	"fmt"

	"github.com/gorilla/websocket"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// LaunchLobbyHistory launch local lobby with rooms from database
func LaunchLobbyHistory(db *database.DataBase,
	ws *websocket.Conn, user *models.UserPublicInfo,
	WSsettings config.WebSocketSettings, gameSettings config.GameConfig) {

	urls, err := db.GetGamesURL(user.ID)

	if err != nil {
		fmt.Println("GetGamesURL", err.Error())
		return
	}

	lobby := NewLobby(gameSettings.ConnectionCapacity, len(urls),
		gameSettings.LobbyJoin, gameSettings.LobbyRequest, db,
		gameSettings.CanClose)
	go lobby.Run(false)
	defer func() {
		fmt.Println("stop lobby!")
		lobby.Stop()
	}()

	err = lobby.LoadRooms(urls)

	if err != nil {
		fmt.Println("LoadRooms", err.Error())
		return
	}

	conn := NewConnection(ws, user, lobby)
	conn.Launch(WSsettings)

}
