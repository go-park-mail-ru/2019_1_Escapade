package api

import (
	database "escapade/internal/database"
	"escapade/internal/misc"
	"escapade/internal/models"
	re "escapade/internal/return_errors"
	"escapade/internal/services/game"
	"fmt"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/websocket"

	//"reflect"

	"github.com/gorilla/mux"
)

// Handler is struct
type Handler struct {
	DB                    database.DataBase
	PlayersAvatarsStorage string
	FileMode              int
	ReadBufferSize        int
	WriteBufferSize       int
	Lobby                 *game.Lobby
	Test                  bool
}

// catch CORS preflight
// @Summary catch CORS preflight
// @Description catch CORS preflight
// @ID OK1
// @Success 200 "successfully"
// @Router /user [OPTIONS]
func (h *Handler) Ok(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, nil, "Ok")

	fmt.Println("api/ok - ok")
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
		err      error
		username string
	)

	if username, err = h.getNameFromCookie(r); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		sendErrorJSON(rw, re.ErrorAuthorization(), place)
		printResult(err, http.StatusUnauthorized, place)
		return
	}

	if err = sendPublicUser(h, rw, username, place); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, re.ErrorServer(), place)
		printResult(err, http.StatusUnauthorized, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	printResult(err, http.StatusOK, place)
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
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if sessionID, err = h.DB.Register(&user); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	misc.CreateAndSet(rw, sessionID)
	rw.WriteHeader(http.StatusCreated)
	sendSuccessJSON(rw, nil, place)
	printResult(err, http.StatusCreated, place)
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
		user models.UserPrivateInfo
		err  error
		name string
	)

	if user, err = getUser(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if name, err = h.getNameFromCookie(r); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		sendErrorJSON(rw, re.ErrorAuthorization(), place)
		printResult(err, http.StatusUnauthorized, place)
		return
	}

	if err = h.DB.UpdatePlayerByName(name, &user); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, nil, place)
	printResult(err, http.StatusOK, place)
	return
}

// Login login
// @Summary login
// @Description login
// @ID Login
// @Success 200 {object} models.UserPublicInfo "Get successfully"
// @Failure 400 {object} models.Result "invalid name/email or password"
// @Failure 500 {object} models.Result "server error"
// @Router /session [POST]
func (h *Handler) Login(rw http.ResponseWriter, r *http.Request) {
	const place = "Login"
	var (
		user      models.UserPrivateInfo
		err       error
		sessionID string
		username  string
	)

	if user, err = getUser(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if sessionID, username, err = h.DB.Login(&user); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, re.ErrorUserNotFound(), place)
		printResult(err, http.StatusBadRequest, place)
		return
	}
	misc.CreateAndSet(rw, sessionID)

	if err = sendPublicUser(h, rw, username, place); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, re.ErrorDataBase(), place)
		printResult(err, http.StatusInternalServerError, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	printResult(err, http.StatusOK, place)
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

	if sessionID, err = misc.GetSessionCookie(r); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		sendErrorJSON(rw, re.ErrorAuthorization(), place)
		printResult(err, http.StatusUnauthorized, place)
		return
	}

	if err = h.DB.Logout(sessionID); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, re.ErrorDataBase(), place)
		printResult(err, http.StatusInternalServerError, place)
		return
	}

	misc.CreateAndSet(rw, "")
	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, nil, place)
	printResult(err, http.StatusOK, place)
	return
}

// DeleteAccount delete account
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
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if err = h.DB.DeleteAccount(&user); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, re.ErrorUserNotFound(), place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	misc.CreateAndSet(rw, "")
	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, nil, place)
	printResult(err, http.StatusOK, place)
	return
}

// GetPlayerGames get games
// @Summary get users game
// @Description Get amount of users list page
// @ID GetPlayerGames
// @Success 200 {array} models.Game "Get successfully"
// @Failure 400 {object} models.Result "invalid username or page"
// @Failure 404 {object} models.Result "games not found"
// @Failure 500 {object} models.Result "Databse error"
// @Router /users/{name}/games/{page} [GET]
func (h *Handler) GetPlayerGames(rw http.ResponseWriter, r *http.Request) {
	const place = "GetPlayerGames"

	var (
		err      error
		games    []models.Game
		username string
		page     int
	)

	if page, username, err = h.getNameAndPage(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if games, err = h.DB.GetGames(username, page); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		sendErrorJSON(rw, re.ErrorGamesNotFound(), place)
		printResult(err, http.StatusNotFound, place)
		return
	}

	sendSuccessJSON(rw, games, place)
	rw.WriteHeader(http.StatusOK)
	printResult(err, http.StatusOK, place)
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
		pages models.Pages
		err   error
	)

	if pages.Amount, err = h.DB.GetUsersPageAmount(); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, re.ErrorDataBase(), place)
		printResult(err, http.StatusInternalServerError, place)
		return
	}

	sendSuccessJSON(rw, pages, place)
	rw.WriteHeader(http.StatusOK)
	printResult(err, http.StatusOK, place)
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
		err   error
		users []models.UserPublicInfo
		page  int
	)

	if page, err = h.getPage(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, re.ErrorInvalidPage(), place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if users, err = h.DB.GetUsers(page); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		sendErrorJSON(rw, re.ErrorUsersNotFound(), place)
		printResult(err, http.StatusNotFound, place)
		return
	}

	sendSuccessJSON(rw, users, place)
	rw.WriteHeader(http.StatusOK)
	printResult(err, http.StatusOK, place)
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
		err      error
		userID   int
		filename string
		filepath string
		file     []byte
	)

	if userID, err = h.getUserIDFromCookie(r); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		sendErrorJSON(rw, re.ErrorAuthorization(), place)
		printResult(err, http.StatusUnauthorized, place)
		return
	}

	if filename, err = h.DB.GetImage(userID); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		sendErrorJSON(rw, re.ErrorAvatarNotFound(), place)
		printResult(err, http.StatusNotFound, place)
		return
	}

	filepath = h.PlayersAvatarsStorage + strconv.Itoa(userID) + "/" + filename

	if file, err = ioutil.ReadFile(filepath); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		sendErrorJSON(rw, re.ErrorAvatarNotFound(), place)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(file)
	printResult(err, http.StatusOK, place)
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
	)

	if userID, err = h.getUserIDFromCookie(r); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		sendErrorJSON(rw, re.ErrorAuthorization(), place)
		printResult(err, http.StatusUnauthorized, place)
		return
	}

	if input, handle, err = r.FormFile("file"); err != nil || input == nil || handle == nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, re.ErrorInvalidFile(), place)
		printResult(err, http.StatusInternalServerError, place)
		return
	}

	defer input.Close()

	fileType := handle.Header.Get("Content-Type")
	fileName := handle.Filename
	storagePath := h.PlayersAvatarsStorage
	filePath := storagePath + strconv.Itoa(userID)

	switch fileType {
	case "image/jpeg":
		err = saveFile(filePath, fileName, input, os.FileMode(h.FileMode))
	case "image/png":
		err = saveFile(filePath, fileName, input, os.FileMode(h.FileMode))
	default:
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, re.ErrorInvalidFileFormat(), place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, re.ErrorServer(), place)
		printResult(err, http.StatusInternalServerError, place)
		return
	}

	if err = h.DB.PostImage(fileName, userID); err != nil {
		_ = os.Remove(filePath) // if error then lets delete uploaded image
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, re.ErrorDataBase(), place)
		printResult(err, http.StatusInternalServerError, place)
		return
	}

	sendSuccessJSON(rw, nil, place)
	rw.WriteHeader(http.StatusCreated)
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
		err      error
		username string
	)

	vars := mux.Vars(r)
	username = vars["name"]

	if username == "" {
		err = re.ErrorInvalidName()
		rw.WriteHeader(http.StatusBadGateway)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadGateway, place)
		return
	}

	if err = sendPublicUser(h, rw, username, place); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		sendErrorJSON(rw, re.ErrorUserNotFound(), place)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	printResult(err, http.StatusOK, place)
	return
}

// GameOnline launch multiplayer
func (h *Handler) GameOnline(rw http.ResponseWriter, r *http.Request) {
	const place = "GameOnline"
	var (
		err      error
		userID   int
		userName string
		ws       *websocket.Conn
	)

	if !h.Test {
		if userName, err = h.getNameFromCookie(r); err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			sendErrorJSON(rw, re.ErrorAuthorization(), place)
			printResult(err, http.StatusUnauthorized, place)
			return
		}

		if userID, err = h.getUserIDFromCookie(r); err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			sendErrorJSON(rw, re.ErrorAuthorization(), place)
			printResult(err, http.StatusUnauthorized, place)
			return
		}
	} else {
		userName = game.RandString(16)
		userID = rand.Intn(10000)
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  h.ReadBufferSize,
		WriteBufferSize: h.WriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	if ws, err = upgrader.Upgrade(rw, r, rw.Header()); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		if _, ok := err.(websocket.HandshakeError); ok {
			sendErrorJSON(rw, re.ErrorHandshake(), place)
		} else {
			sendErrorJSON(rw, re.ErrorNotWebsocket(), place)
		}
		printResult(err, http.StatusBadRequest, place)
		return
	}

	player := game.NewPlayer(userName, userID)
	conn := game.NewConnection(ws, player, h.Lobby)
	// Join Player to lobby
	h.Lobby.ChanJoin <- conn

	fmt.Printf("Player: %d has joined \n", conn.GetPlayerID())

	//rw.WriteHeader(http.StatusOK)
	printResult(err, http.StatusOK, place)
	return
}
