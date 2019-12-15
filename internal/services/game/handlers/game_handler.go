package handlers

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	pkgDB "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/photo"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/clients"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/constants"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine"
	gmetrics "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/metrics"
)

// GameHandler handler of game service requests
type GameHandler struct {
	c          *config.Configuration
	upgraderWS websocket.Upgrader
	db         *database.Input
	lobby      *engine.Lobby
}

// ConfigurationArgs - arguments, requeired to initialize game configuration
type ConfigurationArgs struct {
	C         *config.Configuration
	FieldPath string
	RoomPath  string
}

// InitWithPostgresql initialize handler with PostgreSQL database
func (h *GameHandler) InitWithPostgresql(chatS clients.ChatI, ca *ConfigurationArgs) error {
	var (
		rep = new(constants.RepositoryFS)
		da  = new(database.Input).InitAsPSQL()
	)
	return h.Init(rep, chatS, ca, da)
}

// Init initialize connection handler with chat service, configuration settings and
// 	database settings
func (h *GameHandler) Init(rep constants.RepositoryI,
	chatS clients.ChatI, ca *ConfigurationArgs, da *database.Input) error {

	err := h.initDB(ca.C.DataBase, da)
	if err != nil {
		return err
	}

	err = h.initConstants(rep, ca)
	if err != nil {
		return err
	}

	h.c = ca.C
	h.upgraderWS = websocket.Upgrader{
		HandshakeTimeout: h.c.WebSocket.HandshakeTimeout.Duration,
		ReadBufferSize:   h.c.WebSocket.ReadBufferSize,
		WriteBufferSize:  h.c.WebSocket.WriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	gmetrics.Init()

	h.lobby = engine.NewLobby(chatS, &h.c.Game, h.db.Game, photo.GetImages)
	go h.lobby.Run()

	return nil
}

// initDB initialize database
func (h *GameHandler) initDB(c config.Database, da *database.Input) error {
	// open database connection
	if err := da.Database.Open(c); err != nil {
		return err
	}
	h.db = da

	return pkgDB.Use(h.db.Database, h.db.User, h.db.Game)
}

// initConstants initialize game configuration
func (h *GameHandler) initConstants(rep constants.RepositoryI, ca *ConfigurationArgs) error {
	if err := constants.InitField(rep, ca.FieldPath); err != nil {
		utils.Debug(false, "Initialization error with field constants:", err.Error())
		return err
	}

	if err := constants.InitRoom(rep, ca.RoomPath); err != nil {
		utils.Debug(false, "Initialization error with room constants:", err.Error())
		return err
	}
	return nil
}

// Close all connections to database
func (h *GameHandler) Close() error {
	return re.Close(h.db.User, h.db.Game)
}

// Handle connections requests
func (h *GameHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	h.connect(rw, r)
}

func (h *GameHandler) connect(rw http.ResponseWriter, r *http.Request) api.Result {
	const place = "GameOnline"
	roomID := api.StringFromPath(r, "id", "")

	user, err := h.prepareUser(r)
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
	conn, err := engine.NewConnection(ws, user)
	if err != nil {
		utils.Debug(false, "cant create connection:", err.Error())
		return api.NewResult(0, place, nil, err)
	}
	go conn.Launch(h.c.WebSocket, roomID)

	utils.Debug(false, "hi user:", user.ID)
	// code 0 mean nothing to send to client
	return api.NewResult(0, place, nil, nil)
}

// prepare user for game service
func (h *GameHandler) prepareUser(r *http.Request) (*models.UserPublicInfo, error) {
	var (
		userID int32
		err    error
		user   *models.UserPublicInfo
	)
	if userID, err = api.GetUserIDFromAuthRequest(r); err != nil {
		userID = h.lobby.Anonymous()
	}

	if userID < 0 {
		user = h.anonymousUser(userID)
	} else {
		if user, err = h.db.User.FetchOne(userID, 0); err != nil {
			return nil, re.NoUserWrapper(err)
		}
	}
	photo.GetImages(user)
	return user, nil
}

// return not authorized user's model
func (h *GameHandler) anonymousUser(userID int32) *models.UserPublicInfo {
	return &models.UserPublicInfo{
		Name:    h.anonymousID(),
		ID:      userID,
		FileKey: photo.GetDefaultAvatar(),
	}
}

// generate anonymous user's id
func (h *GameHandler) anonymousID() string {
	min := h.c.Game.Anonymous.MinID
	max := h.c.Game.Anonymous.MaxID
	id := rand.Intn(max) + min
	return "Anonymous" + strconv.Itoa(id)
}
