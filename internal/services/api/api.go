package api

import (
	"encoding/json"
	"errors"
	"escapade/internal/config"
	database "escapade/internal/database"
	"escapade/internal/misc"
	"escapade/internal/models"
	"fmt"
	"io"
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

// UploadAvatar uploads avatar
func (h *Handler) PostImage(rw http.ResponseWriter, r *http.Request) {
	const place = "PostImage"

	var (
		err     error
		input   multipart.File
		created *os.File

		sessionID string
		username  string
	)

	input, _, err = r.FormFile("avatar")

	if err != nil || input == nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)
		return
	}

	defer input.Close()

	if sessionID, err = misc.GetSessionCookie(r); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)

		fmt.Println("api/PostImage failed")
		return
	}

	if username, err = h.DB.GetNameBySessionID(sessionID); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)

		fmt.Println("api/PostImage failed")
		return
	}

	imageName := misc.CreateImageName()
	path := h.PlayersAvatarsStorage + username + "/" + imageName

	if created, err = os.Create(path); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)

		fmt.Println("api/PostImage failed")
		return
	}

	defer created.Close()

	if _, err = io.Copy(created, input); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)

		fmt.Println("api/PostImage failed")
		return
	}

	if err = h.DB.PostImage(imageName, username); err != nil {
		// if error then lets delete uploaded image
		_ = os.Remove(path)

		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)

		fmt.Println("api/PostImage failed")
		return
	}
	rw.WriteHeader(http.StatusCreated)
	sendSuccessJSON(rw, place)

	fmt.Println("api/PostImage ok")
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
		err   error
		image models.Image
		file  []byte //*os.File
	)

	if r.Body == nil {
		err = errors.New("JSON not found")
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)

		fmt.Println("api/GetImage doesnt recieve json")

		return
	}
	_ = json.NewDecoder(r.Body).Decode(&image)

	filename := h.PlayersAvatarsStorage + image.Path
	file, err = ioutil.ReadFile(filename)

	if err != nil {
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
