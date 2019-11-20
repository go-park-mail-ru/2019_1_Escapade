package handlers

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	idb "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/photo"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/clients"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/constants"
	gdb "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine"
	gmetrics "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/metrics"
)

type GameHandler struct {
	c          *config.Configuration
	upgraderWS websocket.Upgrader
	user       database.UserUseCaseI
	game       gdb.GameUseCaseI
	lobby      *engine.Lobby
}

type DatabaseArgs struct {
	Database idb.DatabaseI
	User     database.UserRepositoryI
	Record   database.RecordRepositoryI
	Game     gdb.GameRepositoryI
}

type ConfigurationArgs struct {
	C         *config.Configuration
	FieldPath string
	RoomPath  string
}

func (h *GameHandler) InitWithPostgresql(chatS clients.Chat, ca *ConfigurationArgs) error {

	var (
		da = &DatabaseArgs{
			Database: &idb.PostgresSQL{},
			User:     &database.UserRepositoryPQ{},
			Record:   &database.RecordRepositoryPQ{},
			Game:     &gdb.GameRepositoryPQ{},
		}
	)
	return h.Init(chatS, ca, da)
}

func (h *GameHandler) Init(chatS clients.Chat, ca *ConfigurationArgs, da *DatabaseArgs) error {

	err := da.Database.Open(ca.C.DataBase)
	if err != nil {
		return err
	}

	h.c = ca.C
	h.upgraderWS = websocket.Upgrader{
		ReadBufferSize:  h.c.WebSocket.ReadBufferSize,
		WriteBufferSize: h.c.WebSocket.WriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	h.user = &database.UserUseCase{}
	h.user.Init(da.User, da.Record)

	err = h.user.Use(da.Database)
	if err != nil {
		return err
	}

	h.game = &gdb.GameUseCase{}
	h.game.Init(da.Game, chatS)

	err = h.game.Use(da.Database)
	if err != nil {
		return err
	}

	err = constants.InitField(ca.FieldPath)
	if err != nil {
		utils.Debug(false, "Initialization error with field constants:", err.Error())
		return err
	}

	err = constants.InitRoom(ca.RoomPath)
	if err != nil {
		utils.Debug(false, "Initialization error with room constants:", err.Error())
		return err
	}

	gmetrics.Init()

	h.lobby = engine.NewLobby(chatS, &h.c.Game, h.game, photo.GetImages)
	go h.lobby.Run()

	return nil
}

func (h *GameHandler) Close() {
	h.user.Close()
	h.game.Close()
	h.lobby.Close()
}

func (h *GameHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	utils.Debug(false, "catch me!")
	h.connect(rw, r)
	// ih.Route(rw, r, ih.MethodHandlers{
	// 	http.MethodGet:     h.connect,
	// 	http.MethodOptions: nil})
}

func (h *GameHandler) connect(rw http.ResponseWriter, r *http.Request) api.Result {
	const place = "GameOnline"
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
