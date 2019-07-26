package api

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/game"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	clients "github.com/go-park-mail-ru/2019_1_Escapade/internal/clients"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"

	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

// Handler is struct
type Handler struct {
	DB        database.DataBase
	Session   config.SessionConfig
	WebSocket config.WebSocketSettings
	Game      config.GameConfig
	Clients   *clients.Clients
}

func (h *Handler) HandleUser(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateUser(rw, r)
	case http.MethodGet:
		h.GetMyProfile(rw, r)
	case http.MethodDelete:
		h.DeleteUser(rw, r)
	case http.MethodPut:
		h.UpdateProfile(rw, r)
	case http.MethodOptions:
		return
	default:
		utils.Debug(true, "/api/user wrong request:", r.Method)
	}
	return
}

func (h *Handler) HandleSession(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Login(rw, r)
	case http.MethodDelete:
		h.Logout(rw, r)
	case http.MethodOptions:
		return
	default:
		utils.Debug(true, "/api/session wrong request:", r.Method)
	}
	return
}

func (h *Handler) HandleAvatar(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.PostImage(rw, r)
	case http.MethodGet:
		h.GetImage(rw, r)
	case http.MethodOptions:
		return
	default:
		utils.Debug(true, "/api/avatar wrong request:", r.Method)
	}
	return
}

func (h *Handler) HandleOfflineGame(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.SaveRecords(rw, r)
	case http.MethodOptions:
		return
	default:
		utils.Debug(true, "/api/game wrong request:", r.Method)
	}
	return
}

func (h *Handler) HandleUsersPages(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetUsers(rw, r)
	case http.MethodOptions:
		return
	default:
		utils.Debug(true, "/api/users/pages wrong request:", r.Method)
	}
	return
}

// GetUsersPageAmount get amount of users list page
// @Summary amount of users list page
// @Description Get amount of users list page
// @ID GetUsersPageAmount
// @Success 200 {object} models.Pages "Get successfully"
// @Failure 500 {object} models.Result "Server error"
// @Router /users/pages_amount [GET]
func (h *Handler) GetUsersPageAmount(rw http.ResponseWriter, r *http.Request) {
	const place = "GetUsersPageAmount"

	var (
		perPage int
		pages   models.Pages
		err     error
	)

	perPage = h.getPerPage(r)

	if pages.Amount, err = h.DB.GetUsersPageAmount(perPage); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.SendErrorJSON(rw, re.ErrorDataBase(), place)
		utils.PrintResult(err, http.StatusInternalServerError, place)
		return
	}

	utils.SendSuccessJSON(rw, pages, place)
	utils.PrintResult(err, http.StatusOK, place)
}

// GetUsers get users list
// @Summary Get users list
// @Description Get page of user list
// @ID GetUsers
// @Success 200 {array} models.Result "Get successfully"
// @Failure 400 {object} models.Result "Invalid pade"
// @Failure 404 {object} models.Result "Users not found"
// @Failure 500 {object} models.Result "Server error"
// @Router /users/{page} [GET]
func (h *Handler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	const place = "GetUsers"
	var (
		err       error
		users     []*models.UserPublicInfo
		page      int
		perPage   int
		difficult int
		sort      string
	)

	sort = h.getSort(r)
	perPage = h.getPerPage(r)
	page = h.getPage(r)
	difficult = h.getDifficult(r)

	if users, err = h.DB.GetUsers(difficult, page, perPage, sort); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		utils.SendErrorJSON(rw, re.ErrorUsersNotFound(), place)
		utils.PrintResult(err, http.StatusNotFound, place)
		return
	}

	photo.GetImages(users...)

	utils.SendSuccessJSON(rw, models.UsersPublicInfo{users}, place)
	utils.PrintResult(err, http.StatusOK, place)
}

// GetProfile godoc
// @Summary Get public user inforamtion
// @Description get user's best score and best time for a given difficulty, user's id, name and photo
// @ID GetProfile
// @Accept  json
// @Produce  json
// @Param name path string false "User name"
// @Success 200 {object} models.UserPublicInfo "Profile found successfully"
// @Failure 400 {object} models.Result "Invalid username"
// @Failure 404 {object} models.Result "User not found"
// @Router /users/{name}/profile [GET]
func (h *Handler) GetProfile(rw http.ResponseWriter, r *http.Request) {
	const place = "GetProfile"

	var (
		err    error
		userID int
	)

	if userID, err = h.getUserID(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	h.getUser(rw, r, userID)
	return
}

// SaveRecords save only records
func (h *Handler) SaveRecords(rw http.ResponseWriter, r *http.Request) {
	const place = "SaveRecords"
	var (
		err    error
		userID int
		record models.Record
	)
	if userID, err = h.getUserIDFromCookie(r, h.Session); err != nil {
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

func (h *Handler) getUser(rw http.ResponseWriter, r *http.Request, userID int) {
	const place = "GetProfile"

	var (
		err       error
		difficult int
		user      *models.UserPublicInfo
	)

	difficult = h.getDifficult(r)

	if user, err = h.DB.GetUser(userID, difficult); err != nil {

		rw.WriteHeader(http.StatusNotFound)
		utils.SendErrorJSON(rw, re.ErrorUserNotFound(), place)
		utils.PrintResult(err, http.StatusNotFound, place)
		return
	}
	photo.GetImages(user)

	utils.SendSuccessJSON(rw, user, place)

	rw.WriteHeader(http.StatusOK)
	utils.PrintResult(err, http.StatusOK, place)
	return
}

// GameOnline handle game online
func (h *Handler) GameOnline(rw http.ResponseWriter, r *http.Request) {
	const place = "GameOnline"
	var (
		err    error
		userID int
		ws     *websocket.Conn
		user   *models.UserPublicInfo
		roomID string
	)

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
		rw.WriteHeader(http.StatusBadRequest)
		if _, ok := err.(websocket.HandshakeError); ok {
			utils.SendErrorJSON(rw, re.ErrorHandshake(), place)
		} else {
			utils.SendErrorJSON(rw, re.ErrorNotWebsocket(), place)
		}
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	if userID, err = h.getUserIDFromCookie(r, h.Session); err != nil {
		userID = lobby.Anonymous()
		//rw.WriteHeader(http.StatusUnauthorized)
		//utils.SendErrorJSON(rw, re.ErrorAuthorization(), place)
		//utils.PrintResult(err, http.StatusUnauthorized, place)
		//return
	}

	if userID < 0 {
		user = &models.UserPublicInfo{
			Name:    "Anonymous" + strconv.Itoa(rand.Intn(10000)),
			ID:      userID,
			FileKey: "anonymous.jpg",
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
		userID int
		ws     *websocket.Conn
		user   *models.UserPublicInfo
	)

	if userID, err = h.getUserIDFromCookie(r, h.Session); err != nil {
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

// 767 - 1 51
