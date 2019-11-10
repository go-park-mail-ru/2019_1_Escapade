package handlers

import (
	"fmt"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/clients"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	idb "github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"

	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type GameHandler struct {
	c          *config.Configuration
	upgraderWS websocket.Upgrader
	user       database.UserUseCaseI
	game       database.GameUseCaseI
	lobby      *engine.Lobby
}

func (h *GameHandler) InitWithPostgresql(chatS clients.Chat, c *config.Configuration) error {

	var (
		db     = &idb.PostgresSQL{}
		user   = &database.UserRepositoryPQ{}
		record = &database.RecordRepositoryPQ{}
		game   = &database.GameRepositoryPQ{}
	)
	return h.Init(chatS, c, db, user, record, game)
}

func (h *GameHandler) Init(chatS clients.Chat, c *config.Configuration,
	DB idb.DatabaseI,
	userDB database.UserRepositoryI,
	recordDB database.RecordRepositoryI,
	gameDB database.GameRepositoryI,
) error {

	err := DB.Open(c.DataBase)
	if err != nil {
		return err
	}

	h.c = c
	h.upgraderWS = websocket.Upgrader{
		ReadBufferSize:  c.WebSocket.ReadBufferSize,
		WriteBufferSize: c.WebSocket.WriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	h.user = &database.UserUseCase{}
	h.user.Init(userDB, recordDB)

	err = h.user.Use(DB)
	if err != nil {
		return err
	}

	h.game = &database.GameUseCase{}
	h.game.Init(gameDB, chatS)

	err = h.game.Use(DB)
	if err != nil {
		return err
	}

	h.lobby = engine.NewLobby(chatS, &c.Game, h.game, photo.GetImages)
	go h.lobby.Run()

	return nil
}

func (h *GameHandler) Close() {
	h.user.Close()
	h.game.Close()
	h.lobby.Close()
}

func (h *GameHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("catch me!")
	h.connect(rw, r)
	// ih.Route(rw, r, ih.MethodHandlers{
	// 	http.MethodGet:     h.connect,
	// 	http.MethodOptions: nil})
}

func (h *GameHandler) connect(rw http.ResponseWriter, r *http.Request) api.Result {
	const place = "GameOnline"
	fmt.Println("i want!")
	roomID := api.StringFromPath(r, "id", "")

	user, err := h.prepareUser(r, h.user, h.lobby)
	if err != nil {
		return api.NewResult(http.StatusInternalServerError, place, nil, err)
	}
	utils.Debug(false, "prepared user:", user.ID)
	ws, err := h.upgraderWS.Upgrade(rw, r, rw.Header())
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); ok {
			err = re.ErrorHandshake()
		} else {
			err = re.ErrorNotWebsocket()
		}
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}
	utils.Debug(false, "will create user:", user.ID)
	conn, err := engine.NewConnection(ws, user, h.lobby)
	if err != nil {
		utils.Debug(false, "cant create connection:", err.Error())
		return api.NewResult(0, place, nil, err)
	}
	go conn.Launch(h.c.WebSocket, roomID)
	utils.Debug(false, "hi user:", user.ID)
	// code 0 mean nothing to send to client
	return api.NewResult(0, place, nil, nil)
}

func (h *GameHandler) prepareUser(r *http.Request, userDB database.UserUseCaseI,
	lobby *engine.Lobby) (*models.UserPublicInfo, error) {
	var (
		userID int32
		err    error
		user   *models.UserPublicInfo
	)
	if userID, err = api.GetUserIDFromAuthRequest(r); err != nil {
		userID = lobby.Anonymous()
	}

	if userID < 0 {
		user = h.anonymousUser(userID)
	} else {
		if user, err = userDB.FetchOne(userID, 0); err != nil {
			return nil, re.NoUserWrapper(err)
		}
	}
	photo.GetImages(user)
	return user, nil
}

func (h *GameHandler) anonymousUser(userID int32) *models.UserPublicInfo {
	return &models.UserPublicInfo{
		Name:    h.anonymousID(),
		ID:      userID,
		FileKey: photo.GetDefaultAvatar(),
	}
}

func (h *GameHandler) anonymousID() string {
	min := h.c.Game.Anonymous.MinID
	max := h.c.Game.Anonymous.MaxID
	id := rand.Intn(max) + min
	return "Anonymous" + strconv.Itoa(id)
}
