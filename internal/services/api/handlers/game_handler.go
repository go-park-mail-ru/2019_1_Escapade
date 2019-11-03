package handlers

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	ih "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"

	//"github.com/go-park-mail-ru/2019_1_Escapade/internal/game"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	//"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"
	//"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"

	"net/http"
	//"github.com/gorilla/websocket"
)

type GameHandler struct {
	ih.Handler
	record database.RecordUseCaseI
}

func (h *GameHandler) Init(c *config.Configuration, DB database.DatabaseI,
	recordDB database.RecordRepositoryI) error {
	h.Handler.Init(c)

	h.record = &database.RecordUseCase{}
	h.record.Init(recordDB)
	err := h.record.Use(DB)
	if err != nil {
		return err
	}
	return nil
}

func (h *GameHandler) Close() {
	h.record.Close()
}

// TODO add fetch
// HandleOfflineGame process any operation associated with offline
// games: save
func (h *GameHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	ih.Route(rw, r, ih.MethodHandlers{
		http.MethodPost:    h.OfflineSave,
		http.MethodOptions: nil})
}

// OfflineSave save offline game results
// @Summary Save offline game
// @Description Save offline game results
// @Accept  json
// @Param record body models.Record true "Results of offline game"
// @Produce  json
// @ID OfflineSave
// @Success 200 {object} models.Result "Done"
// @Failure 400 {object} models.Result "Wrong input data"
// @Failure 401 {object} models.Result "Required authorization"
// @Failure 404 {object} models.Result "Avatar not found"
// @Router /game [POST]
func (h *GameHandler) OfflineSave(rw http.ResponseWriter, r *http.Request) ih.Result {
	const place = "OfflineSave"
	var (
		err    error
		userID int32
		record models.Record
	)
	if userID, err = ih.GetUserIDFromAuthRequest(r); err != nil {
		return ih.NewResult(http.StatusUnauthorized, place, nil, re.AuthWrapper(err))
	}

	if err = ih.ModelFromRequest(r, &record); err != nil {
		return ih.NewResult(http.StatusBadRequest, place, nil, err)
	}
	if err = h.record.Update(userID, &record); err != nil {
		return ih.NewResult(http.StatusInternalServerError, place, nil, err)
	}
	return ih.NewResult(http.StatusOK, place, nil, nil)
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

// 199
