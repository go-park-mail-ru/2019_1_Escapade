package api

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/game"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"

	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

// OfflineSave save only records
func (h *Handler) OfflineSave(rw http.ResponseWriter, r *http.Request) {
	const place = "OfflineSave"
	var (
		err    error
		userID int32
		record models.Record
	)
	if userID, err = h.getUserIDFromAuthRequest(r); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		utils.SendErrorJSON(rw, re.ErrorAuthorization(), place)
		utils.PrintResult(err, http.StatusUnauthorized, place)
		return
	}
	if record, err = getRecord(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}
	if err = h.DB.UpdateRecords(userID, &record); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}
}

// GameOnline launch multiplayer
/*
func (h *Handler) SaveGame(rw http.ResponseWriter, r *http.Request) {
	const place = "SaveOfflineGame"
	var (
		err             error
		userID          int
		gameInformation *models.GameInformation
	)
	if userID, err = h.getUserIDFromCookie(r, h.Session); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		utils.SendErrorJSON(rw, re.ErrorAuthorization(), place)
		utils.PrintResult(err, http.StatusUnauthorized, place)
		return
	}
	if gameInformation, err = getGameInformation(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}
	if err = h.DB.SaveGame(userID, gameInformation); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}
}
*/

// GameOnline handle game online
func (h *Handler) GameOnline(rw http.ResponseWriter, r *http.Request) {
	const place = "GameOnline"
	var (
		err    error
		userID int32
		ws     *websocket.Conn
		user   *models.UserPublicInfo
		roomID string
	)

	utils.Debug(false, "GameOnline")
	lobby := game.GetLobby()
	if lobby == nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.SendErrorJSON(rw, re.ErrorServer(), place)
		utils.PrintResult(err, http.StatusInternalServerError, place)
	}

	roomID = getStringFromPath(r, "id", "")

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  h.WebSocket.ReadBufferSize,
		WriteBufferSize: h.WebSocket.WriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	if ws, err = upgrader.Upgrade(rw, r, rw.Header()); err != nil {
		utils.Debug(false, "cant upgrade", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		if _, ok := err.(websocket.HandshakeError); ok {
			utils.SendErrorJSON(rw, re.ErrorHandshake(), place)
		} else {
			utils.SendErrorJSON(rw, re.ErrorNotWebsocket(), place)
		}
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	if userID, err = h.getUserIDFromAuthRequest(r); err != nil {
		userID = lobby.Anonymous()
		//rw.WriteHeader(http.StatusUnauthorized)
		//utils.SendErrorJSON(rw, re.ErrorAuthorization(), place)
		//utils.PrintResult(err, http.StatusUnauthorized, place)
		//return
	}

	if userID < 0 {
		user = &models.UserPublicInfo{
			Name:    "Anonymous" + strconv.Itoa(rand.Intn(10000)),
			ID:      int32(userID),
			FileKey: photo.GetDefaultAvatar(),
		}
	} else {
		if user, err = h.DB.GetUser(userID, 0); err != nil {
			rw.WriteHeader(http.StatusNotFound)
			utils.SendErrorJSON(rw, re.ErrorUserNotFound(), place)
			utils.PrintResult(err, http.StatusNotFound, place)
			return
		}
	}
	photo.GetImages(user)

	conn := game.NewConnection(ws, user, lobby)
	conn.Launch(h.WebSocket, roomID)

	utils.PrintResult(err, http.StatusOK, place)
	return
}

// GameHistory launch local lobby only for this connection
func (h *Handler) GameHistory(rw http.ResponseWriter, r *http.Request) {
	const place = "GameHistory:"
	var (
		err    error
		userID int32
		ws     *websocket.Conn
		user   *models.UserPublicInfo
	)

	if userID, err = h.getUserIDFromAuthRequest(r); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		utils.SendErrorJSON(rw, re.ErrorAuthorization(), place)
		utils.PrintResult(err, http.StatusUnauthorized, place)
		return
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  h.WebSocket.ReadBufferSize,
		WriteBufferSize: h.WebSocket.WriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	if ws, err = upgrader.Upgrade(rw, r, rw.Header()); err != nil {
		utils.Debug(true, place, "can't upgrade the http connection to websockets")
		rw.WriteHeader(http.StatusBadRequest)
		if _, ok := err.(websocket.HandshakeError); ok {
			utils.SendErrorJSON(rw, re.ErrorHandshake(), place)
		} else {
			utils.SendErrorJSON(rw, re.ErrorNotWebsocket(), place)
		}
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	if user, err = h.DB.GetUser(userID, 0); err != nil {
		utils.Debug(false, place, "cant get user from database")
		rw.WriteHeader(http.StatusNotFound)
		utils.SendErrorJSON(rw, re.ErrorUserNotFound(), place)
		utils.PrintResult(err, http.StatusNotFound, place)
		return
	}
	photo.GetImages(user)

	game.LaunchLobbyHistory(&h.DB, ws, user, h.WebSocket, &h.Game, photo.GetImages)
	return
}
