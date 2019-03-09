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

// Вызывать с defer в начале функций
func errorSheduler(rw http.ResponseWriter, err error, who string) {
	if err != nil {
		sendErrorJSON(rw, err, who)
		fmt.Println(who+" failed:", err.Error())
	}
}

func (h *Handler) getNameFromCookie(r *http.Request) (username string, err error) {
	var sessionID string

	if sessionID, err = misc.GetSessionCookie(r); err != nil {
		return
	}

	if username, err = h.DB.GetNameBySessionID(sessionID); err != nil {
		return
	}

	return
}

// UploadAvatar uploads avatar
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
	}

	if input, handle, err = r.FormFile("file"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)
		fmt.Println("api/PostImage failed")
		return
	}

	if input == nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)
		fmt.Println("api/PostImage failed")
		return
	}

	defer input.Close()

	if handle == nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)
		fmt.Println("api/PostImage failed")
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

		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)

		fmt.Println("api/PostImage failed")
		return
	}

	rw.WriteHeader(http.StatusCreated)
	sendSuccessJSON(rw, place)

	fmt.Println("api/PostImage ok")
}

func saveFile(path string, name string, file multipart.File) (err error) {
	var (
		data []byte
	)

	os.MkdirAll(path, 0777)

	if data, err = ioutil.ReadAll(file); err != nil {
		return
	}

	if err = ioutil.WriteFile(path+"/"+name, data, 0666); err != nil {
		return
	}

	return
}

// Ok always returns StatusOk
func (h *Handler) Ok(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, "Ok")

	fmt.Println("api/ok - ok")
	return
}

func (h *Handler) Me(rw http.ResponseWriter, r *http.Request) {

	const place = "Me"

	var (
		err       error
		sessionID string
		username  string
	)

	if sessionID, err = misc.GetSessionCookie(r); err != nil {
		rw.WriteHeader(http.StatusForbidden)
		sendErrorJSON(rw, err, place)

		fmt.Println("api/Me failed")
		return
	}

	if username, err = h.DB.GetNameBySessionID(sessionID); err != nil {
		rw.WriteHeader(http.StatusForbidden)
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

	sessionCookie := misc.CreateCookie(sessionID)
	http.SetCookie(rw, sessionCookie)
	rw.WriteHeader(http.StatusCreated)
	sendSuccessJSON(rw, place)

	fmt.Println("api/register ok")

	return
}

// Login handle login
func (h *Handler) Login(rw http.ResponseWriter, r *http.Request) {
	const place = "Login"
	var (
		user      models.UserPrivateInfo
		err       error
		username  string
		sessionID string
	)

	if user, err = getUser(r); err != nil {
		rw.WriteHeader(http.StatusForbidden)
		sendErrorJSON(rw, err, place)
		return
	}

	if sessionID, err = h.DB.Login(&user); err != nil {
		rw.WriteHeader(http.StatusForbidden)
		sendErrorJSON(rw, err, place)
		return
	}

	sessionCookie := misc.CreateCookie(sessionID)
	http.SetCookie(rw, sessionCookie)

	if username, err = h.DB.GetNameBySessionID(sessionID); err != nil {
		rw.WriteHeader(http.StatusForbidden)
		sendErrorJSON(rw, err, place)

		fmt.Println("api/Login failed")
		return
	}

	if err = sendPublicUser(h, rw, username, place); err != nil {
		fmt.Println("api/Login failed")
		return
	}

	fmt.Println("api/Login ok")

	return
}

// Login handle logout
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

	http.SetCookie(rw, misc.CreateCookie(""))
	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, place)

	fmt.Println("api/logout ok")

	return
}

// DeleteAccount handle registration
func (h *Handler) DeleteAccount(rw http.ResponseWriter, r *http.Request) {

	const place = "DeleteAccount"
	var (
		user      models.UserPrivateInfo
		err       error
		sessionID string
	)

	if user, err = getUser(r); err != nil {
		rw.WriteHeader(http.StatusForbidden)
		sendErrorJSON(rw, err, place)

		fmt.Println("api/DeleteAccount failed")
		return
	}

	if sessionID, err = misc.GetSessionCookie(r); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)

		return
	}

	if sessionID, err = h.DB.DeleteAccount(&user, sessionID); err != nil {
		rw.WriteHeader(http.StatusForbidden)
		sendErrorJSON(rw, err, place)

		fmt.Println("api/DeleteAccount failed")
		return
	}

	http.SetCookie(rw, misc.CreateCookie(""))
	rw.WriteHeader(http.StatusOK)

	fmt.Println("api/DeleteAccount ok")
	return
}

// DeleteAccountOptions handle preCORS request
func (h *Handler) DeleteAccountOptions(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("api/DeleteAccountOptions ok")
	rw.WriteHeader(http.StatusOK)
}

// GetPlayerGames handle get games list
func (h *Handler) GetPlayerGames(rw http.ResponseWriter, r *http.Request) {
	const place = "GetPlayerGames"

	var (
		err      error
		games    []models.Game
		bytes    []byte
		vars     map[string]string
		username string
		page     int
	)

	vars = mux.Vars(r)

	if username = vars["name"]; username == "" {
		fmt.Println("No username found")

		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, errors.New("No username found"), place)
		return
	}

	if vars["page"] == "" {
		page = 1
	} else {
		if page, err = strconv.Atoi(vars["page"]); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			sendErrorJSON(rw, err, place)
			return
		}
		if page < 1 {
			page = 1
		}
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
