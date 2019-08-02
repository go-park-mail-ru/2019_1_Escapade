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
	WSsettings config.WebSocketSettings, gameSettings *config.GameConfig,
	si SetImage) {

	urls, err := db.GetGamesURL(user.ID)

	if err != nil {
		fmt.Println("GetGamesURL", err.Error())
		return
	}

	gameSettings.RoomsCapacity = int32(len(urls) * 2)
	lobby := NewLobby(gameSettings, db, si)

	go lobby.Run()
	defer func() {
		fmt.Println("stop lobby!")
		lobby.Stop()
	}()

	if len(urls) > 0 {
		err = lobby.LoadRooms(urls)
		if err != nil {
			fmt.Println("LoadRooms", err.Error())
			return
		}
	}

	fmt.Println("connection create!")
	conn := NewConnection(ws, user, lobby)
	conn.Launch(WSsettings, "")
	fmt.Println("conn launch")
}
