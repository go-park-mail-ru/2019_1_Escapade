package api

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/cookie"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/game"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"

	clients "github.com/go-park-mail-ru/2019_1_Escapade/internal/clients"

	uuid "github.com/satori/go.uuid"
)

// Handler is struct
type Handler struct {
	DB        database.DataBase
	Storage   config.FileStorageConfig
	Session   config.SessionConfig
	WebSocket config.WebSocketSettings
	Game      config.GameConfig
	AWS       config.AwsPublicConfig
	Clients   *clients.Clients
}

// Ok catch CORS preflight
// @Summary catch CORS preflight
// @Description catch CORS preflight
// @ID OK1
// @Success 200 "successfully"
// @Router /user [OPTIONS]
func (h *Handler) Ok(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	const place = "api/Ok"
	utils.SendSuccessJSON(rw, nil, place)
	utils.PrintResult(nil, http.StatusOK, place)
	return
}

// GetMyProfile get public user information
// @Summary get user
// @Description get public user information
// @ID GetMyProfile
// @Success 200 {object} models.UserPublicInfo "successfully"
// @Failure 401 {object} models.Result "Required authorization"
// @Failure 500 {object} models.Result "server error"
// @Router /user [GET]
func (h *Handler) GetMyProfile(rw http.ResponseWriter, r *http.Request) {

	const place = "GetMyProfile"
	var (
		err    error
		userID int
	)

	if userID, err = h.getUserIDFromCookie(r, h.Session); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		utils.SendErrorJSON(rw, re.ErrorAuthorization(), place)
		utils.PrintResult(err, http.StatusUnauthorized, place)
		return
	}

	h.getUser(rw, r, userID)

	return
}

// CreateUser create new user
// @Summary create new user
// @Description create new user
// @ID Register
// @Success 201 {object} models.Result "Create user successfully"
// @Failure 400 {object} models.Result "Invalid information"
// @Router /user [POST]
func (h *Handler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	const place = "CreateUser"
	var (
		user      models.UserPrivateInfo
		err       error
		sessionID string
	)

	if user, err = getUserWithAllFields(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	if err = validateUser(&user); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusBadRequest, place)
	}

	if _, sessionID, err = h.register(r.Context(), user); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	cookie.CreateAndSet(rw, h.Session, sessionID)
	rw.WriteHeader(http.StatusCreated)
	utils.SendSuccessJSON(rw, nil, place)
	utils.PrintResult(err, http.StatusCreated, place)
	return
}

// UpdateProfile updates profile
// @Summary update user information
// @Description update public info
// @ID UpdateProfile
// @Success 200 {object} models.Result "Get successfully"
// @Failure 400 {object} models.Result "invalid info"
// @Failure 401 {object} models.Result "need authorization"
// @Router /user [PUT]
func (h *Handler) UpdateProfile(rw http.ResponseWriter, r *http.Request) {
	const place = "UpdateProfile"

	var (
		user   models.UserPrivateInfo
		err    error
		userID int
	)

	if user, err = getUser(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	if userID, err = h.getUserIDFromCookie(r, h.Session); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		utils.SendErrorJSON(rw, re.ErrorAuthorization(), place)
		utils.PrintResult(err, http.StatusUnauthorized, place)
		return
	}

	if err = h.DB.UpdatePlayerPersonalInfo(userID, &user); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	utils.SendSuccessJSON(rw, nil, place)
	utils.PrintResult(err, http.StatusOK, place)
	return
}

// Login login
// @Summary login
// @Description login
// @ID Login
// @Success 200 {object} models.UserPublicInfo "Get successfully"
// @Failure 400 {object} models.Result "invalid name or password"
// @Failure 500 {object} models.Result "server error"
// @Router /session [POST]
func (h *Handler) Login(rw http.ResponseWriter, r *http.Request) {
	const place = "Login"
	var (
		user        models.UserPrivateInfo
		err         error
		found       *models.UserPublicInfo
		sessionName string
	)

	if user, err = getUser(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	sessionName = utils.RandomString(16)
	if found, err = h.DB.Login(&user, sessionName); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, re.ErrorUserNotFound(), place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	/*
		ctx := r.Context()

		sessionID, err := h.Clients.Session.Create(ctx,
			&session.Session{
				UserID: int32(user.ID),
				Login:  user.Name,
			})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("cookie set ", sessionID.ID)
	*/

	cookie.CreateAndSet(rw, h.Session, sessionName)

	utils.SendSuccessJSON(rw, found, place)

	rw.WriteHeader(http.StatusOK)
	utils.PrintResult(err, http.StatusOK, place)
	return
}

// Logout logout
// @Summary logout
// @Description logout
// @ID Logout
// @Success 200 {object} models.Result "Get successfully"
// @Success 401 {object} models.Result "Require authorization"
// @Failure 500 {object} models.Result "server error"
// @Router /session [DELETE]
func (h *Handler) Logout(rw http.ResponseWriter, r *http.Request) {
	const place = "Logout"
	var (
		err       error
		sessionID string
	)

	if sessionID, err = cookie.GetSessionCookie(r, h.Session); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		utils.SendErrorJSON(rw, re.ErrorAuthorization(), place)
		utils.PrintResult(err, http.StatusUnauthorized, place)
		return
	}
	h.DB.DeleteSession(sessionID)
	/*
		ctx := context.Background()
		_, err = h.Clients.Session.Delete(ctx,
			&session.SessionID{
				ID: sessionID,
			})
		if err != nil {
			fmt.Println(err)
			return
		}*/

	cookie.CreateAndSet(rw, h.Session, "")
	rw.WriteHeader(http.StatusOK)
	utils.SendSuccessJSON(rw, nil, place)
	utils.PrintResult(err, http.StatusOK, place)
	return
}

// DeleteUser delete account
// @Summary delete account
// @Description delete account
// @ID DeleteAccount
// @Success 200 {object} models.Result "Get successfully"
// @Failure 400 {object} models.Result "invalid input"
// @Failure 500 {object} models.Result "server error"
// @Router /user [DELETE]
func (h *Handler) DeleteUser(rw http.ResponseWriter, r *http.Request) {

	const place = "DeleteUser"
	var (
		user models.UserPrivateInfo
		err  error
	)

	if user, err = getUserWithAllFields(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	if err = h.deleteAccount(context.Background(), &user, ""); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, re.ErrorUserNotFound(), place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	cookie.CreateAndSet(rw, h.Session, "")
	rw.WriteHeader(http.StatusOK)
	utils.SendSuccessJSON(rw, nil, place)
	utils.PrintResult(err, http.StatusOK, place)
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

	if err = h.Setfiles(users...); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusNotFound, place)
		return
	}

	utils.SendSuccessJSON(rw, users, place)
	utils.PrintResult(err, http.StatusOK, place)
}

// GetImage returns user avatar
// @Summary Get user avatar
// @Description Get user avatar
// @ID GetImage
// @Success 200 {object} models.Result "Avatar found successfully"
// @Failure 401 {object} models.Result "Required authorization"
// @Failure 404 {object} models.Result "Avatar not found"
// @Router /avatar [GET]
func (h *Handler) GetImage(rw http.ResponseWriter, r *http.Request) {
	const place = "GetImage"
	var (
		err     error
		userID  int
		fileKey string
		url     models.Avatar
	)

	//
	if userID, err = h.getUserIDFromCookie(r, h.Session); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		utils.SendErrorJSON(rw, re.ErrorAuthorization(), place)
		utils.PrintResult(err, http.StatusUnauthorized, place)
		return
	}

	if fileKey, err = h.DB.GetImage(userID); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		utils.SendErrorJSON(rw, re.ErrorAvatarNotFound(), place)
		utils.PrintResult(err, http.StatusNotFound, place)
		return
	}

	url.URL, err = h.getURLToAvatar(fileKey)
	if err != nil {
		log.Println("Failed to sign request", err)
		rw.WriteHeader(http.StatusNotFound)
		utils.SendErrorJSON(rw, re.ErrorAvatarNotFound(), place)
		utils.PrintResult(err, http.StatusNotFound, place)
		return
	}

	utils.SendSuccessJSON(rw, url, place)
	utils.PrintResult(err, http.StatusOK, place)
}

// PostImage create avatar
// @Summary Create user avatar
// @Description Create user avatar
// @ID PostImage
// @Success 201 {object} models.Result "Avatar created successfully"
// @Failure 401 {object} models.Result "Required authorization"
// @Failure 500 {object} models.Result "Avatar not found"
// @Router /avatar [POST]
func (h *Handler) PostImage(rw http.ResponseWriter, r *http.Request) {
	const place = "PostImage"

	var (
		err    error
		input  multipart.File
		userID int
		handle *multipart.FileHeader
		url    models.Avatar
	)

	if userID, err = h.getUserIDFromCookie(r, h.Session); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		utils.SendErrorJSON(rw, re.ErrorAuthorization(), place)
		utils.PrintResult(err, http.StatusUnauthorized, place)
		return
	}

	if input, handle, err = r.FormFile("file"); err != nil || input == nil || handle == nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.SendErrorJSON(rw, re.ErrorInvalidFile(), place)
		utils.PrintResult(err, http.StatusInternalServerError, place)
		return
	}

	defer input.Close()

	fileType := handle.Header.Get("Content-Type")
	//Генерация уник.ключа для хранения картинки
	fileKey := uuid.NewV4()

	switch fileType {
	case "image/jpeg":
		err = h.saveFile(fileKey.String(), input)
	case "image/png":
		err = h.saveFile(fileKey.String(), input)
	default:
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, re.ErrorInvalidFileFormat(), place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.SendErrorJSON(rw, re.ErrorServer(), place)
		utils.PrintResult(err, http.StatusInternalServerError, place)
		return
	}

	if err = h.DB.PostImage(fileKey.String(), userID); err != nil {
		h.deleteAvatar(fileKey.String())
		rw.WriteHeader(http.StatusInternalServerError)
		utils.SendErrorJSON(rw, re.ErrorDataBase(), place)
		utils.PrintResult(err, http.StatusInternalServerError, place)
		return
	}

	if url.URL, err = h.getURLToAvatar(fileKey.String()); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		utils.SendErrorJSON(rw, re.ErrorAvatarNotFound(), place)
		utils.PrintResult(err, http.StatusNotFound, place)
		return
	}

	utils.SendSuccessJSON(rw, url, place)
	utils.PrintResult(err, http.StatusCreated, place)
}

// GetProfile returns model UserPublicInfo
// @Summary Get some of user fields
// @Description return public information, such as name or best_score
// @ID GetProfile
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
	if err = h.Setfiles(user); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusNotFound, place)
		return
	}

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
		h.Setfiles(user)
	} else {
		if user, err = h.DB.GetUser(userID, 0); err != nil {
			rw.WriteHeader(http.StatusNotFound)
			utils.SendErrorJSON(rw, re.ErrorUserNotFound(), place)
			utils.PrintResult(err, http.StatusNotFound, place)
			return
		}

		if err = h.Setfiles(user); err != nil {
			rw.WriteHeader(http.StatusNotFound)
			utils.SendErrorJSON(rw, err, place)
			utils.PrintResult(err, http.StatusNotFound, place)
			return
		}
	}

	conn := game.NewConnection(ws, user, lobby)
	conn.Launch(h.WebSocket, roomID)

	utils.PrintResult(err, http.StatusOK, place)
	return
}

// GameHistory launch local lobby only for this connection
func (h *Handler) GameHistory(rw http.ResponseWriter, r *http.Request) {
	const place = "GameHistory"
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
		fmt.Println("err689", err.Error())
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
		fmt.Println("err700", err.Error())
		rw.WriteHeader(http.StatusNotFound)
		utils.SendErrorJSON(rw, re.ErrorUserNotFound(), place)
		utils.PrintResult(err, http.StatusNotFound, place)
		return
	}

	if err = h.Setfiles(user); err != nil {
		fmt.Println("err707", err.Error())
		rw.WriteHeader(http.StatusNotFound)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusNotFound, place)
		return
	}

	game.LaunchLobbyHistory(&h.DB, ws, user, h.WebSocket, h.Game, h.Setfiles)
	return
}
