package game

import (
	"escapade/internal/config"
	"escapade/internal/database"
	"escapade/internal/models"
	re "escapade/internal/return_errors"
	"escapade/internal/utils"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	clients "escapade/internal/clients"

	"github.com/gorilla/websocket"
)

// Handler is struct
type Handler struct {
	DB              database.DataBase
	Storage         config.FileStorageConfig
	Cookie          config.CookieConfig
	WebSocket       config.WebSocketSettings
	GameConfig      config.GameConfig
	AWS             config.AwsPublicConfig
	ReadBufferSize  int
	WriteBufferSize int
	Test            bool
	Clients         *clients.Clients
}

var API *Handler

// catch CORS preflight
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

// GameOnline launch multiplayer
func (h *Handler) Chat(rw http.ResponseWriter, r *http.Request) {
	const place = "Chat"
	var (
		err    error
		userID int
		ws     *websocket.Conn
		name   string
		user   *models.UserPublicInfo
	)

	userID, name, _ = h.getUserIDFromCookie(r, h.Cookie)

	lobby := GetLobby()
	if lobby == nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.SendErrorJSON(rw, re.ErrorServer(), place)
		utils.PrintResult(err, http.StatusInternalServerError, place)
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
			utils.SendErrorJSON(rw, re.ErrorHandshake(), place)
		} else {
			utils.SendErrorJSON(rw, re.ErrorNotWebsocket(), place)
		}
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	if userID <= 0 {
		rand.Seed(time.Now().UnixNano())
		userID = rand.Intn(100000) + 100000
		name = "Guest" + strconv.FormatInt(int64(userID), 10)
	}
	user = &models.UserPublicInfo{
		ID:   userID,
		Name: name,
	}

	conn := NewConnection(ws, user, lobby)
	conn.Launch(h.WebSocket)
	utils.PrintResult(err, http.StatusOK, place)
	return
}
