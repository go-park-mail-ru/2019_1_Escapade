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

// PostImage posts avatar
func (h *Handler) PostImage(rw http.ResponseWriter, r *http.Request) {
	const place = "PostImage"

	var (
		err      error
		input    multipart.File
		username string
		handle   *multipart.FileHeader
	)

	if username, err = h.getNameFromCookie(r); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)
		fmt.Println("api/PostImage failed")
		return
	}

	if input, handle, err = r.FormFile("file"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)
		fmt.Println("api/PostImage failed")
		return
	}

	if input == nil {
		err = errors.New("no input")
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)
		fmt.Println("api/PostImage failed")
		return
	}

	defer input.Close()

	if handle == nil {
		err = errors.New("no handle")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileType := handle.Header.Get("Content-Type")
	fileName := handle.Filename
	storagePath := h.PlayersAvatarsStorage
	filePath := storagePath + username

	switch fileType {
	case "image/jpeg":
		err = saveFile(filePath, fileName, input)
	case "image/png":
		err = saveFile(filePath, fileName, input)
	default:
		err = errors.New("wrong format of file:" + fileType)
	}

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)
		fmt.Println("api/PostImage failed")
		return
	}

	if err = h.DB.PostImage(fileName, username); err != nil {
		// if error then lets delete uploaded image
		_ = os.Remove(filePath)
		sendErrorJSON(rw, err, place)
		fmt.Println("api/PostImage failed")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
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

// Ok always returns StatusOk
func (h *Handler) Ok(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, nil, "Ok")

	fmt.Println("api/ok - ok")
	return
}

// Me returns my profile
func (h *Handler) Me(rw http.ResponseWriter, r *http.Request) {

	const place = "Me"
	var (
		err      error
		username string
	)

	if username, err = h.getNameFromCookie(r); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)
		fmt.Println("api/Me failed")
		return
	}

	if err = sendPublicUser(h, rw, username, place); err != nil {
		fmt.Println("api/Me failed")
		return
	}

	fmt.Println("api/Me ok")

	return
}

// Register handle registration
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

// Login handle login
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
		fmt.Println("api/Login failed")
		return
	}

	fmt.Println("api/Login ok")

	return
}

// Logout handle logout
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

// DeleteAccount deletes user
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

// DeleteAccountOptions handle preCORS request
func (h *Handler) DeleteAccountOptions(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("api/DeleteAccountOptions ok")
	rw.WriteHeader(http.StatusOK)
}

// GetPlayerGames retur
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

// GetUsersPageAmount returns amount of pages of users
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

// GetPlayerGames handle get games list
func (h *Handler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	const place = "GetUsers"

	var (
		err   error
		users []models.UserPublicInfo
		bytes []byte
		page  int
	)

	if page, err = h.getPage(r); err != nil {
		fmt.Println("No page found")

		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)
		return
	}

	if users, err = h.DB.GetUsers(page); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)
		return
	}

	if bytes, err = json.Marshal(users); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)

		fmt.Println("api/GetPlayerGames cant create json")
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(bytes)
	fmt.Println("api/GetPlayerGames ok")
}

func (h *Handler) GetImage(rw http.ResponseWriter, r *http.Request) {
	const place = "GetImage"
	var (
		err      error
		username string
		filename string
		filepath string
		file     []byte
	)

	if username, err = h.getNameFromCookie(r); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)

		fmt.Println("api/GetImage cant found file")

		return
	}

	if filename, err = h.DB.GetImage(username); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)

		fmt.Println("api/GetImage cant found file")

		return
	}

	filepath = h.PlayersAvatarsStorage + username + "/" + filename

	if file, err = ioutil.ReadFile(filepath); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)

		fmt.Println("api/GetImage cant found file")

		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(file)
	fmt.Println("api/GetImage ok")

}

func (h *Handler) GetProfile(rw http.ResponseWriter, r *http.Request) {

	const place = "GetProfile"

	var (
		err      error
		vars     map[string]string
		username string
	)

	vars = mux.Vars(r)

	if username = vars["name"]; username == "" {
		fmt.Println("No username found")

		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, errors.New("No username found"), place)
		return
	}

	if err = sendPublicUser(h, rw, username, place); err != nil {
		fmt.Println("api/GetProfile failed")
		return
	}

	fmt.Println("api/GetProfile ok")

	return
}
