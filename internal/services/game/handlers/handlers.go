package handlers

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine"

	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"

	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

func gameOnline(db *database.DataBase, c *config.Configuration, upgraderWS websocket.Upgrader) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}
		result := connect(rw, r, db, c, upgraderWS)
		api.SendResult(rw, result)
	}
}

func connect(rw http.ResponseWriter, r *http.Request, db *database.DataBase,
	c *config.Configuration, upgraderWS websocket.Upgrader) api.Result {
	const place = "GameOnline"
	var (
		ws   *websocket.Conn
		user *models.UserPublicInfo
	)

	roomID := api.StringFromPath(r, "id", "")

	lobby := engine.GetLobby()
	if lobby == nil {
		return api.NewResult(http.StatusInternalServerError, place, nil, re.ErrorServer())
	}

	err := prepareUser(r, db, lobby)
	if err != nil {
		return api.NewResult(http.StatusInternalServerError, place, nil, err)
	}

	if ws, err = upgraderWS.Upgrade(rw, r, rw.Header()); err != nil {
		if _, ok := err.(websocket.HandshakeError); ok {
			err = re.ErrorHandshake()
		} else {
			err = re.ErrorNotWebsocket()
		}
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	conn := engine.NewConnection(ws, user, lobby)
	conn.Launch(c.WebSocket, roomID)
	// code 0 mean nothing to send to client
	return api.NewResult(0, place, nil, nil)
}

func prepareUser(r *http.Request, db *database.DataBase, lobby *engine.Lobby) error {
	var (
		userID int32
		err    error
		user   *models.UserPublicInfo
	)
	if userID, err = api.GetUserIDFromAuthRequest(r); err != nil {
		userID = lobby.Anonymous()
	}

	if userID < 0 {
		user = anonymousUser(userID)
	} else {
		if user, err = db.GetUser(userID, 0); err != nil {
			return re.NoUserWrapper(err)
		}
	}
	photo.GetImages(user)
	return nil
}

func anonymousUser(userID int32) *models.UserPublicInfo {
	return &models.UserPublicInfo{
		Name:    "Anonymous" + strconv.Itoa(rand.Intn(10000)), // в конфиг
		ID:      int32(userID),
		FileKey: photo.GetDefaultAvatar(),
	}
}
