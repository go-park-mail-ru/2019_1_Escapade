package game

import (
	"fmt"
	"sync"

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
		db, gameSettings.CanClose)
	all := &sync.WaitGroup{}
	all.Add(1)
	go lobby.Run(all)
	defer func() {
		fmt.Println("stop lobby!")
		lobby.Stop()
	}()

	all.Wait()
	err = lobby.LoadRooms(urls)

	if err != nil {
		fmt.Println("LoadRooms", err.Error())
		return
	}

	conn := NewConnection(ws, user, lobby)
	conn.Launch(WSsettings)

}
