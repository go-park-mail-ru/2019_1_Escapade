package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	a_handlers "github.com/go-park-mail-ru/2019_1_Escapade/auth/internal/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	e_server "github.com/go-park-mail-ru/2019_1_Escapade/internal/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"github.com/go-session/session"
	"github.com/jackc/pgx"
	pg "github.com/vgarvardt/go-oauth2-pg"
	"github.com/vgarvardt/go-pg-adapter/pgxadapter"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
)

func main() {
	manager := manage.NewDefaultManager()
	cfg := &manage.Config{
		// access token expiration time
		AccessTokenExp: time.Hour * 2,
		// refresh token expiration time
		RefreshTokenExp: time.Hour * 24 * 14,
		// whether to generate the refreshing token
		IsGenerateRefresh: true,
	}

	manager.SetPasswordTokenCfg(cfg)
	// token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// generate jwt access token
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte("00000000"), jwt.SigningMethodHS512))

	link := "postgres://rolepade:escapade@localhost:5432/escabase?sslmode=disable"
	pgxConnConfig, _ := pgx.ParseURI(link)
	pgxConn, _ := pgx.Connect(pgxConnConfig)

	// use PostgreSQL token store with pgx.Connection adapter
	adapter := pgxadapter.NewConn(pgxConn)
	tokenStore, _ := pg.NewTokenStore(adapter, pg.WithTokenStoreGCInterval(time.Minute))
	defer tokenStore.Close()

	clientStore, _ := pg.NewClientStore(adapter)

	manager.MapTokenStorage(tokenStore)

	clientStore.Create(&models.Client{
		ID:     "1",
		Secret: "1",
		Domain: "http://localhost:3001",
	})
	manager.MapClientStorage(clientStore)

	/*clientStore := store.NewClientStore()
	clientStore.Set("1", &models.Client{
		ID:     "1",
		Secret: "1",
		Domain: "http://localhost:3001",
	})
	manager.MapClientStorage(clientStore)*/

	srv := server.NewServer(server.NewConfig(), manager)
	srv.SetAllowGetAccessRequest(true)
	// allow the grant types model:AuthorizationCode,PasswordCredentials,ClientCredentials,Refreshing
	//srv.SetAllowedGrantType("password_credentials", "refreshing")

	var (
		configuration     *config.Configuration
		err               error
		db                *database.DataBase
		configurationPath = "auth.json"
	)

	configuration, err = config.Init(configurationPath)
	if err != nil {
		utils.Debug(false, "Initialization error with main configuration:", err.Error())
		return
	}

	db, err = database.Init(configuration.DataBase)
	if err != nil {
		utils.Debug(false, "Initialization error with database:", err.Error())
		return
	}

	srv.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		var intUserID int32
		utils.Debug(false, "look at", username, password)
		if intUserID, err = db.Login(username, password); err != nil {
			utils.Debug(false, "whaaaat", err.Error())
			err = re.ErrorUserNotFound()
			return
		}
		userID = utils.String32(intUserID)
		utils.Debug(false, "userID", userID, intUserID)
		return
	})

	//srv.SetAllowedGrantType("password", "refresh_token")
	srv.SetUserAuthorizationHandler(userAuthorizeHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	r := a_handlers.Router(srv, tokenStore)
	port := e_server.Port(configuration)
	server := e_server.Server(r, configuration.Server, true, port)

	e_server.LaunchHTTP(server, configuration.Server, func() { db.Db.Close() })

}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	utils.Debug(false, "/userAuthorizeHandler")
	store, err := session.Start(nil, w, r)
	if err != nil {
		return
	}

	uid, ok := store.Get("LoggedInUserID")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}

		store.Set("ReturnUri", r.Form)
		store.Save()

		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	userID = uid.(string)
	store.Delete("LoggedInUserID")
	store.Save()
	return
}
