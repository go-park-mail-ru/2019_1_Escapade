package handlers

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	ih "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"

	clients "github.com/go-park-mail-ru/2019_1_Escapade/internal/clients"

	"net/http"
)

/*
Handler contains all API operations
DB - database, where api work with information
Cookie - cookie settings, more in structure config.Cookie
Clients - grps.Clients, need to connect to Auth server
Auth - auth settinds, more in structure config.Auth
*/
type Handler struct {
	//DB         database.DataBase
	Cookie     config.Cookie
	Clients    *clients.Clients
	AuthClient config.AuthClient
	Auth       config.Auth
	Db         *useCases
}

type useCases struct {
	user   database.UserUseCaseI
	game   database.GameUseCaseI
	record database.RecordUseCaseI
	image  database.ImageUseCaseI
}

type repositories struct {
	user   database.UserRepositoryI
	game   database.GameRepositoryI
	record database.RecordRepositoryI
	image  database.ImageRepositoryI
}

func (cases *useCases) initAsPG(CDB config.Database) error {
	var (
		reps = repositories{
			user:   &database.UserRepositoryPQ{},
			game:   &database.GameRepositoryPQ{},
			record: &database.RecordRepositoryPQ{},
			image:  &database.ImageRepositoryPQ{},
		}
		database = &database.PostgresSQL{}
	)
	return cases.init(CDB, database, reps)
}

func (cases *useCases) init(CDB config.Database, db database.DatabaseI, reps repositories) error {
	err := db.Open(CDB)
	if err != nil {
		return err
	}

	cases.user = &database.UserUseCase{}
	cases.user.Open(CDB, 10, time.Hour, db)
	cases.user.Init(reps.user, reps.record)
	cases.user.Use(db)
	if err != nil {
		return err
	}

	cases.game = &database.GameUseCase{}
	cases.game.Init(reps.game)
	cases.game.Use(db)
	if err != nil {
		return err
	}

	cases.record = &database.RecordUseCase{}
	cases.record.Init(reps.record)
	cases.record.Use(db)
	if err != nil {
		return err
	}

	cases.image = &database.ImageUseCase{}
	cases.image.Init(reps.image)
	cases.image.Use(db)
	if err != nil {
		return err
	}
	return nil
}

// GetHandler init handler and configuration for api service
func (H *Handler) NEW_Init(c *config.Configuration) {
	H.Cookie = c.Cookie
	H.AuthClient = c.AuthClient
	H.Auth = c.Auth
	H.Db = &useCases{}
}

func (H *Handler) NEW_SetDb(CDB config.Database, db database.DatabaseI,
	reps repositories) error {
	H.Db = &useCases{}
	return H.Db.init(CDB, db, reps)
}

func (H *Handler) NEW_SetPostreSQL(CDB config.Database) error {
	H.Db = &useCases{}
	return H.Db.initAsPG(CDB)
}

// GetHandler init handler and configuration for api service
/*
func GetHandler(c *config.Configuration) (*Handler, error) {

	var (
		db  *database.DataBase
		err error
	)

	if db, err = database.Init(c.DataBase); err != nil {
		return nil, err
	}

	handler := &Handler{
		DB:         *db,
		Cookie:     c.Cookie,
		AuthClient: c.AuthClient,
		Auth:       c.Auth,
	}
	return handler, nil
}*/

// TODO добавить группу ожидания и здесь ждать ее, не забыть использовтаь ustils.WaitWithTimeout
func (h *Handler) Close() {
	//h.DB.Db.Close()
	h.Db.user.Close()
	h.Db.game.Close()
	h.Db.record.Close()
	h.Db.image.Close()
}

// HandleUser process any operation associated with user
// profile: create, receive, update, and delete
func (h *Handler) HandleUser(rw http.ResponseWriter, r *http.Request) {
	ih.Route(rw, r, ih.MethodHandlers{
		http.MethodPost:    h.CreateUser,
		http.MethodGet:     h.GetMyProfile,
		http.MethodDelete:  h.DeleteUser,
		http.MethodPut:     h.UpdateProfile,
		http.MethodOptions: nil})
}

// HandleSession process any operation associated with user
// authorization: enter and exit
func (h *Handler) HandleSession(rw http.ResponseWriter, r *http.Request) {
	ih.Route(rw, r, ih.MethodHandlers{
		http.MethodPost:    h.Login,
		http.MethodDelete:  h.Logout,
		http.MethodOptions: nil})
}

// TODO add deleting
// HandleAvatar process any operation associated with user
// avatar: load and get
func (h *Handler) HandleAvatar(rw http.ResponseWriter, r *http.Request) {
	ih.Route(rw, r, ih.MethodHandlers{
		http.MethodPost:    h.PostImage,
		http.MethodGet:     h.GetImage,
		http.MethodOptions: nil})
}

// TODO add receipt
// HandleOfflineGame process any operation associated with offline
// games: save
func (h *Handler) HandleOfflineGame(rw http.ResponseWriter, r *http.Request) {
	ih.Route(rw, r, ih.MethodHandlers{
		http.MethodPost:    h.OfflineSave,
		http.MethodOptions: nil})
}

// HandleUsersPages process any operation associated with users
// list: receive
func (h *Handler) HandleUsersPages(rw http.ResponseWriter, r *http.Request) {
	ih.Route(rw, r, ih.MethodHandlers{
		http.MethodGet:     h.GetUsers,
		http.MethodOptions: nil})
}

// HandleUsersPageAmount process any operation associated with
// amount of pages in user list: receive
func (h *Handler) HandleUsersPageAmount(rw http.ResponseWriter, r *http.Request) {
	ih.Route(rw, r, ih.MethodHandlers{
		http.MethodGet:     h.GetUsersPageAmount,
		http.MethodOptions: nil})
}

// 128 -> 88
