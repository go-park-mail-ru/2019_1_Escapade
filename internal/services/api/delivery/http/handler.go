package api

import (
	/*
		"net/http"

		"github.com/gorilla/mux"
		"github.com/prometheus/client_golang/prometheus/promhttp"

		"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
		mi "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/middleware"
		"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"
		"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database"
	*/
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
	mi "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/router"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handlers struct {
	Game    GameUseCaseI
	Session SessionUseCaseI
	Image   ImageUseCaseI
	User    UserUseCaseI
	Users   UsersUseCaseI
}

type UseCases struct {
	Image  api.ImageUseCaseI
	User   api.UserUseCaseI
	Record api.RecordUseCaseI
}

func NewHandler(r router.Interface,
	h *Handlers,
	uc *UseCases,
	config *config.AuthToken,
	cors config.CORS,
	subnet string) http.Handler {

	r.PathHandler("/swagger", httpSwagger.Handler(
		httpSwagger.URL("swagger/doc.json"), //The url pointing to API definition"
	))

	mi.Init()

	http.Handle("/metrics", promhttp.Handler())
	r.PathHandler("/metrics", promhttp.Handler())

	var api = r.PathSubrouter("/api")
	var apiWithAuth = r.PathSubrouter("/api")

	apiWithAuth.POST("/game", h.Game.OfflineSave)
	api.OPTIONS("/game", handlers.OPTIONS())

	apiWithAuth.DELETE("/session", h.Session.Logout)
	api.POST("/session", h.Session.Login)
	api.OPTIONS("/session", handlers.OPTIONS())

	api.POST("/user", h.User.CreateUser)
	api.DELETE("/user", h.User.DeleteUser)
	apiWithAuth.PUT("/user", h.User.UpdateProfile)
	apiWithAuth.GET("/user", h.User.GetMyProfile)
	api.OPTIONS("/user", handlers.OPTIONS())

	// delete "/avatar/{name}" path
	api.GET("/avatar/{name}", h.Image.GetImage)
	api.OPTIONS("/avatar/{name}", handlers.OPTIONS())

	apiWithAuth.POST("/avatar", h.Image.PostImage)
	api.OPTIONS("/avatar", handlers.OPTIONS())

	api.GET("/users/{id}", h.Users.GetOneUser)
	api.OPTIONS("/users/{id}", handlers.OPTIONS())

	api.GET("/users/pages/page", h.Users.GetUsers)
	api.OPTIONS("/users/pages/page", handlers.OPTIONS())

	api.GET("/users/pages/amount", h.Users.GetUsersPageAmount)
	api.OPTIONS("/users/pages/amount", handlers.OPTIONS())

	api.PathHandlerFunc("/health", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("all ok " + server.GetIP(nil)))
	})

	api.Use(mi.CORS(cors), mi.Metrics(subnet))
	apiWithAuth.Use(mi.Auth(config.Cookie, config.Auth, config.AuthClient))
	return r.Router()
}

/*
type Handlers struct {
	c  *config.Configuration
	db *database.Input

	user    *UserHandler
	users   *UsersHandler
	game    *GameHandler
	session *SessionHandler
	image   *ImageHandler

	subnet string
}

// InitWithPostgreSQL apply postgreSQL as database
func (h *Handlers) InitWithPostgreSQL(subnet string, c *config.Configuration) error {
	return h.OpenDB(subnet, c, new(database.Input).InitAsPSQL())
}

// Init open connection to database and put it to all handlers
func (h *Handlers) OpenDB(subnet string, c *config.Configuration, input *database.Input) error {
	input.Init()
	if err := input.IsValid(); err != nil {
		return err
	}

	h.subnet = subnet
	h.c = c
	h.db = input

	err := input.Open(c.DataBase)
	if err != nil {
		return err
	}

	h.user = new(UserHandler).Init(c, input)
	h.session = new(SessionHandler).Init(c, input)
	h.game = new(GameHandler).Init(c, input)
	h.users = new(UsersHandler).Init(c, input)
	h.image = new(ImageHandler).Init(c, input)
	return nil
}

// Close connections to darabase of all handlers
func (h *Handlers) Close() error {
	return h.db.Close()
}

func some() {
	r := mux.NewRouter()

	var ewe router.Interface
	ewe = r
}

// Router return router of api operations
func (h *Handlers) Router() *mux.Router {

	r := mux.NewRouter()

	r.PathPrefix("/swagger").Handler(httpSwagger.Handler(
		httpSwagger.URL("swagger/doc.json"), //The url pointing to API definition"
	))

	mi.Init()

	http.Handle("/metrics", promhttp.Handler())
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

	api.PathPrefix("/health").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("all ok " + server.GetIP(nil)))
	})

	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	router.Use(r, mi.CORS(h.c.Cors), mi.Metrics(h.subnet))
	//router.Use(api, mi.Metrics)
	apiWithAuth.Use(mi.Auth(h.c.Cookie, h.c.Auth, h.c.AuthClient))
	return r
}*/

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
