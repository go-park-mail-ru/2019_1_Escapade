package oauth

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
	"github.com/go-session/session"
	"github.com/jackc/pgx"
	pg "github.com/vgarvardt/go-oauth2-pg"
	"github.com/vgarvardt/go-pg-adapter/pgxadapter"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
)

func Init(accessTokenExp, refreshTokenExp time.Duration,
	isGenerateRefresh bool, jwtSecret, pgHost string, pgPort int,
	clients []*models.Client) (*manage.Manager, *pg.TokenStore, error) {
	manager := manage.NewDefaultManager()
	cfg := &manage.Config{
		AccessTokenExp:    accessTokenExp,
		RefreshTokenExp:   refreshTokenExp,
		IsGenerateRefresh: isGenerateRefresh,
	}

	manager.SetPasswordTokenCfg(cfg)
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte(jwtSecret), jwt.SigningMethodHS512))

	//link := "dbname=escabase user=rolepade password=escapade sslmode=disable"
	link := "postgres://rolepade:escapade@localhost:5432/escabase?sslmode=disable"
	pgxConnConfig, _ := pgx.ParseURI(link)
	//pgxConnConfig.Host = e_server.FixHost(pgHost)
	//pgxConnConfig.Port = uint16(pgPort)
	pgxConn, err := pgx.Connect(pgxConnConfig)
	if err != nil {
		return nil, nil, err
	}

	adapter := pgxadapter.NewConn(pgxConn)
	tokenStore, err := pg.NewTokenStore(adapter, pg.WithTokenStoreGCInterval(time.Minute))
	if err != nil {
		return nil, nil, err
	}
	//defer tokenStore.Close()

	manager.MapTokenStorage(tokenStore)

	clientStore, err := pg.NewClientStore(adapter)
	if err != nil {
		tokenStore.Close()
		return nil, nil, err
	}

	if err = addClients(clientStore, clients); err != nil {
		tokenStore.Close()
		return nil, nil, err
	}
	manager.MapClientStorage(clientStore)

	return manager, tokenStore, err
}

func addClients(store *pg.ClientStore, clients []*models.Client) error {
	if len(clients) == 0 {
		return store.Create(&models.Client{
			ID:     "1",
			Secret: "1",
			Domain: "http://localhost:3001",
		})
	} else {
		var err error
		for _, client := range clients {
			err = store.Create(client)
			if err != nil {
				utils.Debug(false, "Warning:", err.Error())
				//break
			}
		}
		return nil //err
	}
}

func Server(db *database.DataBase, manager *manage.Manager) *server.Server {
	srv := server.NewServer(
		&server.Config{
			TokenType:            "Bearer",
			AllowedResponseTypes: []oauth2.ResponseType{oauth2.Code, oauth2.Token},
			AllowedGrantTypes: []oauth2.GrantType{
				oauth2.PasswordCredentials,
				oauth2.Refreshing,
			},
		}, manager)
	srv.SetAllowGetAccessRequest(true)
	// allow the grant types model:AuthorizationCode,PasswordCredentials,ClientCredentials,Refreshing
	//srv.SetAllowedGrantType("password_credentials", "refreshing")

	srv.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		var intUserID int32
		if intUserID, err = db.Login(username, password); err != nil {
			err = re.ErrorUserNotFound()
			return
		}
		userID = utils.String32(intUserID)
		utils.Debug(false, "userID", userID, intUserID)
		return
	})

	srv.SetUserAuthorizationHandler(userAuthorizeHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		utils.Debug(false, "Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		utils.Debug(false, "Response Error:", re.Error.Error())
	})

	return srv
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
