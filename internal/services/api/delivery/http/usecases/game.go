package handlers

import (
	"net/http"
	//"github.com/gorilla/websocket"

	ih "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"

	//"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/photo"
	//"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api"
)

// GameHandler handle requests associated with singleplayer
type GameHandler struct {
	record api.RecordUseCaseI
}

func NewGameHandler(record api.RecordUseCaseI) *GameHandler {
	handler := &GameHandler{record}
	return handler
}

// OfflineSave save offline game results
// @Summary Save offline game
// @Description Save offline game results of current user. The current one is the one whose token is provided.
// @ID OfflineSave
// @Security OAuth2Application[write]
// @Tags game
// @Accept  json
// @Param record body models.Record true "Results of offline game"
// @Produce  json
// @Success 200 {object} models.Result "Done"
// @Failure 400 {object} models.Result "Invalid data for save"
// @Failure 401 {object} models.Result "Required authorization"
// @Failure 500 {object} models.Result "Database error"
// @Router /game [POST]
func (h *GameHandler) OfflineSave(rw http.ResponseWriter, r *http.Request) ih.Result {
	var (
		err    error
		userID int32
		record models.Record
	)
	userID, err = ih.GetUserIDFromAuthRequest(r)
	if err != nil {
		return ih.NewResult(http.StatusUnauthorized, nil, re.AuthWrapper(err))
	}

	if err = ih.ModelFromRequest(r, &record); err != nil {
		return ih.NewResult(http.StatusBadRequest, nil, err)
	}

	err = h.record.Update(r.Context(), userID, &record)
	if err != nil {
		return ih.NewResult(http.StatusInternalServerError, nil, err)
	}
	return ih.NewResult(http.StatusOK, nil, nil)
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

// GameHistory launch local lobby only for this connection
/*
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
*/
