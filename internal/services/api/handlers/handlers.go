package handlers

import (
	"net/http"
	"time"
        "fmt"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	idb "github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	mi "github.com/go-park-mail-ru/2019_1_Escapade/internal/middleware"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/router"
	server "github.com/go-park-mail-ru/2019_1_Escapade/internal/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handlers struct {
	c       *config.Configuration
	user    *UserHandler
	users   *UsersHandler
	game    *GameHandler
	session *SessionHandler
	image   *ImageHandler
}

// repositories stores all implementations of operations in the database
type repositories struct {
	user   database.UserRepositoryI
	game   database.GameRepositoryI
	record database.RecordRepositoryI
	image  database.ImageRepositoryI
}

// InitWithPostgreSQL apply postgreSQL as database
func (h *Handlers) InitWithPostgreSQL(c *config.Configuration) error {
	var (
		reps = repositories{
			user:   &database.UserRepositoryPQ{},
			game:   &database.GameRepositoryPQ{},
			record: &database.RecordRepositoryPQ{},
			image:  &database.ImageRepositoryPQ{},
		}
		database = &idb.PostgresSQL{}
	)
	return h.Init(c, database, reps)
}

// Init open connection to database and put it to all handlers
func (h *Handlers) Init(c *config.Configuration, db idb.DatabaseI, reps repositories) error {
        fmt.Println("string:", c.DataBase.ConnectionString)
	h.c = c
	err := db.Open(c.DataBase)
	if err != nil {
		return err
	}

	h.user = &UserHandler{}
	err = h.user.Init(c, db, reps.user, reps.record)
	if err != nil {
		return err
	}

	h.session = &SessionHandler{}
	err = h.session.Init(c, db, reps.user, reps.record)
	if err != nil {
		return err
	}

	h.game = &GameHandler{}
	err = h.game.Init(c, db, reps.record)
	if err != nil {
		return err
	}

	h.users = &UsersHandler{}
	err = h.users.Init(c, db, reps.user, reps.record)
	if err != nil {
		return err
	}

	h.image = &ImageHandler{}
	err = h.image.Init(c, db, reps.image)
	if err != nil {
		return err
	}
	return nil
}

// Close connections to darabase of all handlers
func (h *Handlers) Close() {
	h.user.Close()
	h.users.Close()
	h.session.Close()
	h.game.Close()
	h.image.Close()
}

// Router return router of api operations
func (h *Handlers) Router() *mux.Router {

	r := mux.NewRouter()

	r.PathPrefix("/swagger").Handler(httpSwagger.Handler(
		httpSwagger.URL("swagger/doc.json"), //The url pointing to API definition"
	))

	r.PathPrefix("/metrics").Handler(promhttp.Handler())

	var api = r.PathPrefix("/api").Subrouter()
	var apiWithAuth = r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/user", h.user.Handle).Methods("OPTIONS", "POST", "DELETE")
	apiWithAuth.HandleFunc("/user", h.user.Handle).Methods("PUT", "GET")

	api.HandleFunc("/session", h.session.Handle).Methods("POST", "OPTIONS")
	apiWithAuth.HandleFunc("/session", h.session.Handle).Methods("DELETE")

	// delete "/avatar/{name}" path
	api.HandleFunc("/avatar/{name}", h.image.Handle).Methods("GET")

	api.HandleFunc("/avatar", h.image.Handle).Methods("OPTIONS")
	apiWithAuth.HandleFunc("/avatar", h.image.Handle).Methods("POST")

	api.HandleFunc("/game", h.game.Handle).Methods("OPTIONS")
	apiWithAuth.HandleFunc("/game", h.game.Handle).Methods("POST")

	api.HandleFunc("/users/{id}", h.users.HandleGetProfile).Methods("GET", "OPTIONS")
	api.HandleFunc("/users/pages/page", h.users.HandleUsersPages).Methods("GET", "OPTIONS")
	api.HandleFunc("/users/pages/amount", h.users.HandleUsersPageAmount).Methods("GET")

	r.PathPrefix("/health").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		rw.Write([]byte("all ok " + server.GetIP()))
	})

	r.PathPrefix("/hard").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		time.Sleep(7 * time.Second)
		rw.Write([]byte("hard done " + server.GetIP()))
	})

	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	router.Use(r, mi.CORS(h.c.Cors))
	apiWithAuth.Use(mi.Auth(h.c.Cookie, h.c.Auth, h.c.AuthClient))
	return r
}

// HistoryRouter return router for history service
/*
func HistoryRouter(handler *api.Handler, cors config.CORS) *mux.Router {
	r := mux.NewRouter()

	var history = r.PathPrefix("/history").Subrouter()

	history.Use(mi.Recover, mi.CORS(cors), mux.CORSMethodMiddleware(r))

	history.HandleFunc("/ws", handler.GameOnline)
	history.Handle("/metrics", promhttp.Handler())
	return r
}
*/

// 128 -> 88
