package api

import (
	"encoding/json"
	"errors"
	"escapade/internal/config"
	database "escapade/internal/database"
	"escapade/internal/misc"
	"escapade/internal/models"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	//"reflect"

	"github.com/gorilla/mux"
)

// Handler is struct
type Handler struct {
	DB                    database.DataBase
	PlayersAvatarsStorage string
}

// Init creates Handler
func Init(DB *database.DataBase, storage config.FileStorageConfig) (handler *Handler) {
	handler = &Handler{
		DB:                    *DB,
		PlayersAvatarsStorage: storage.PlayersAvatarsStorage,
	}
	return
}

func saveFile(path string, name string, file multipart.File) (err error) {
	var (
		data []byte
	)

	// вынести в конфиг
	const mode777 = 0777
	const mode666 = 0666

	os.MkdirAll(path, mode777)

	if data, err = ioutil.ReadAll(file); err != nil {
		return
	}

	if err = ioutil.WriteFile(path+"/"+name, data, mode666); err != nil {
		return
	}

	return
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

// catch CORS preflight
// @Summary catch CORS preflight
// @Description catch CORS preflight
// @ID OK2
// @Success 200 "successfully"
// @Router /user/login [OPTIONS]
func (h *Handler) Ok2(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, nil, "Ok")

	fmt.Println("api/ok - ok")
	return
}

// catch CORS preflight
// @Summary catch CORS preflight
// @Description catch CORS preflight
// @ID OK3
// @Success 200 "successfully"
// @Router /user/logout [OPTIONS]
func (h *Handler) Ok3(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, nil, "Ok")

	fmt.Println("api/ok - ok")
	return
}

// catch CORS preflight
// @Summary catch CORS preflight
// @Description catch CORS preflight
// @ID OK4
// @Success 200 "successfully"
// @Router /user/Avatar [OPTIONS]
func (h *Handler) Ok4(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, nil, "Ok")

	fmt.Println("api/ok - ok")
	return
}

// GetMyProfile get user profile
// @Summary get user
// @Description get public user information
// @ID GetMyProfile
// @Success 200 {object} models.UserPublicInfo "successfully"
// @Failure 500 {object} models.Result "server error"
// @Router /user [GET]
func (h *Handler) GetMyProfile(rw http.ResponseWriter, r *http.Request) {

	const place = "GetMyProfile"
	var (
		err      error
		username string
	)

	if username, err = h.getNameFromCookie(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		fmt.Println("api/Me failed")
		return
	}

	if err = sendPublicUser(h, rw, username, place); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Println("api/Me failed")
		return
	}

	rw.WriteHeader(http.StatusOK)
	fmt.Println("api/Me ok")
	return
}

// Register create new user
// @Summary create new user
// @Description create new user
// @ID Register
// @Success 200 {object} models.Result "successfully"
// @Failure 500 {object} models.Result "server error"
// @Router /user [POST]
func (h *Handler) Register(rw http.ResponseWriter, r *http.Request) {
	const place = "Register"

	var (
		user      models.UserPrivateInfo
		err       error
		sessionID string
	)

	if user, err = getUser(r); err != nil {
		rw.WriteHeader(http.StatusForbidden)
		sendErrorJSON(rw, err, place)
		fmt.Println("api/register failed")
		return
	}

	if sessionID, err = h.DB.Register(&user); err != nil {
		rw.WriteHeader(http.StatusForbidden)
		sendErrorJSON(rw, err, place)
		fmt.Println("api/register failed")
		return
	}

	misc.CreateAndSet(rw, sessionID)
	rw.WriteHeader(http.StatusCreated)
	sendSuccessJSON(rw, nil, place)

	fmt.Println("api/register ok")

	return
}

// UpdateProfile updates profile
// @Summary update user information
// @Description update public info
// @ID UpdateProfile
// @Success 200 {object} models.Result "Get successfully"
// @Failure 400 {object} models.Result "invalid info"
// @Failure 500 {object} models.Result "server error"
// @Router /user [PUT]
func (h *Handler) UpdateProfile(rw http.ResponseWriter, r *http.Request) {
	const place = "UpdateProfile"

	var (
		user models.UserPrivateInfo
		err  error
		name string
	)

	if user, err = getUser(r); err != nil {
		rw.WriteHeader(http.StatusForbidden)
		sendErrorJSON(rw, err, place)
		fmt.Println("api/UpdateProfile failed")
		return
	}

	if name, err = h.getNameFromCookie(r); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)
		fmt.Println("api/UpdateProfile failed")
		return
	}

	if err = h.DB.UpdatePlayerByName(name, &user); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)
		fmt.Println("api/UpdateProfile failed")
		return
	}

	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, nil, place)

	fmt.Println("api/UpdateProfile ok")

	return
}

// Login login
// @Summary login
// @Description login
// @ID Login
// @Success 200 {object} models.UserPublicInfo "Get successfully"
// @Failure 400 {object} models.Result "invalid name"
// @Failure 500 {object} models.Result "server error"
// @Router /user/login [POST]
func (h *Handler) Login(rw http.ResponseWriter, r *http.Request) {
	const place = "Login"
	var (
		user      models.UserPrivateInfo
		err       error
		sessionID string
		username  string
	)

	if user, err = getUser(r); err != nil {
		rw.WriteHeader(http.StatusForbidden)
		sendErrorJSON(rw, err, place)
		fmt.Println("api/Login failed")
		return
	}

	if sessionID, username, err = h.DB.Login(&user); err != nil {
		rw.WriteHeader(http.StatusForbidden)
		sendErrorJSON(rw, err, place)
		fmt.Println("api/Login failed")
		return
	}
	misc.CreateAndSet(rw, sessionID)

	if err = sendPublicUser(h, rw, username, place); err != nil {
		rw.WriteHeader(http.StatusForbidden)
		fmt.Println("api/Login failed")
		return
	}

	rw.WriteHeader(http.StatusOK)
	fmt.Println("api/Login ok")

	return
}

// Logout logout
// @Summary logout
// @Description logout
// @ID Logout
// @Success 200 {object} models.Result "Get successfully"
// @Failure 500 {object} models.Result "server error"
// @Router /user/logout [DELETE]
func (h *Handler) Logout(rw http.ResponseWriter, r *http.Request) {
	const place = "Logout"

	var (
		err       error
		sessionID string
	)

	if sessionID, err = misc.GetSessionCookie(r); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)
		return
	}

	if err = h.DB.Logout(sessionID); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)
		return
	}

	misc.CreateAndSet(rw, "")
	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, nil, place)

	fmt.Println("api/logout ok")

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
func (h *Handler) DeleteAccount(rw http.ResponseWriter, r *http.Request) {

	const place = "DeleteAccount"
	var (
		user models.UserPrivateInfo
		err  error
	)

	if user, err = getUser(r); err != nil {
		rw.WriteHeader(http.StatusForbidden)
		sendErrorJSON(rw, err, place)
		fmt.Println("api/DeleteAccount failed")
		return
	}

	if err = h.DB.DeleteAccount(&user); err != nil {
		rw.WriteHeader(http.StatusForbidden)
		sendErrorJSON(rw, err, place)
		fmt.Println("api/DeleteAccount failed")
		return
	}

	misc.CreateAndSet(rw, "")
	rw.WriteHeader(http.StatusOK)

	fmt.Println("api/DeleteAccount ok")
	return
}

func (h *Handler) DeleteAccountOptions(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("api/DeleteAccountOptions ok")
	rw.WriteHeader(http.StatusOK)
}

// GetPlayerGames get games
// @Summary get users game
// @Description Get amount of users list page
// @ID GetPlayerGames
// @Success 200 {array} models.Game "Get successfully"
// @Failure 400 {object} models.Result "invalid username or page"
// @Failure 404 {object} models.Result "games not found"
// @Failure 500 {object} models.Result "server error"
// @Router /users/{name}/games/{page} [GET]
func (h *Handler) GetPlayerGames(rw http.ResponseWriter, r *http.Request) {
	const place = "GetPlayerGames"

	var (
		err      error
		games    []models.Game
		bytes    []byte
		username string
		page     int
	)

	if page, username, err = h.getNameAndPage(r); err != nil {
		fmt.Println("No username found")

		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, errors.New("No username found"), place)
		return
	}

	if games, err = h.DB.GetGames(username, page); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)
		return
	}

	if bytes, err = json.Marshal(games); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)

		fmt.Println("api/GetPlayerGames cant create json")
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(bytes)
	fmt.Println("api/GetPlayerGames ok")
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
		bytes []byte
	)

	if pages.Amount, err = h.DB.GetUsersPageAmount(); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)
		fmt.Println("api/GetUsersPageAmount cant work with DB")
		return
	}

	if bytes, err = json.Marshal(pages); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)

		fmt.Println("api/GetUsersAmount cant create json")
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(bytes)
	fmt.Println("api/GetUsersAmount ok")
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
		bytes []byte
		page  int
	)

	if page, err = h.getPage(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, ErrorInvalidPage(), place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if users, err = h.DB.GetUsers(page); err != nil {
		rw.WriteHeader(http.StatusNoContent)
		sendErrorJSON(rw, ErrorUsersNotFound(), place)
		printResult(err, http.StatusNoContent, place)
		return
	}

	if bytes, err = json.Marshal(users); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, ErrorServer(), place)
		printResult(err, http.StatusInternalServerError, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(bytes)
	printResult(err, http.StatusOK, place)
}

// GetImage returns user avatar
// @Summary Get user avatar
// @Description Get user avatar
// @ID GetImage
// @Success 200 {object} models.Result "Avatar found successfully"
// @Failure 401 {object} models.Result "Required authorization"
// @Failure 404 {object} models.Result "Avatar not found"
// @Router /user/Avatar [GET]
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
		sendErrorJSON(rw, ErrorAuthorization(), place)
		printResult(err, http.StatusUnauthorized, place)
		return
	}

	if filename, err = h.DB.GetImage(userID); err != nil {
		rw.WriteHeader(http.StatusNoContent)
		sendErrorJSON(rw, ErrorAvatarNotFound(), place)
		printResult(err, http.StatusNoContent, place)
		return
	}

	filepath = h.PlayersAvatarsStorage + strconv.Itoa(userID) + "/" + filename

	if file, err = ioutil.ReadFile(filepath); err != nil {
		rw.WriteHeader(http.StatusNoContent)
		sendErrorJSON(rw, ErrorAvatarNotFound(), place)
		printResult(err, http.StatusNoContent, place)
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
// @Router /user/Avatar [POST]
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
		sendErrorJSON(rw, ErrorAuthorization(), place)
		printResult(err, http.StatusUnauthorized, place)
		return
	}

	if input, handle, err = r.FormFile("file"); err != nil || input == nil || handle == nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, ErrorInvalidFile(), place)
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
		err = saveFile(filePath, fileName, input)
	case "image/png":
		err = saveFile(filePath, fileName, input)
	default:
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, ErrorInvalidFileFormat(), place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, ErrorServer(), place)
		printResult(err, http.StatusInternalServerError, place)
		return
	}

	if err = h.DB.PostImage(fileName, userID); err != nil {
		_ = os.Remove(filePath) // if error then lets delete uploaded image
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, ErrorDataBase(), place)
		printResult(err, http.StatusInternalServerError, place)
		return
	}

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
		err = ErrorInvalidName()
		rw.WriteHeader(http.StatusBadGateway)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadGateway, place)
		return
	}

	if err = sendPublicUser(h, rw, username, place); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		sendErrorJSON(rw, ErrorUserNotFound(), place)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	printResult(err, http.StatusOK, place)

	return
}

// 536
