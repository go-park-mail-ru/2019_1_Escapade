package oauth2server

import (
	"context"
	"net/http"
	"time"

	err "errors"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/infrastructure/configuration"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/usecase/user"
	"github.com/go-session/session"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
)

type Oauth2Server struct {
	srv            *server.Server
	logger         infrastructure.Logger
	trace          infrastructure.ErrorTrace
	contextTimeout time.Duration
	userDB         *user.UserDB
}

func New(
	tokenConfig configuration.TokenRepository,
	db infrastructure.Database,
	logger infrastructure.Logger,
	trace infrastructure.ErrorTrace,
	manager *manage.Manager,
	contextTimeout time.Duration,
) (*Oauth2Server, error) {
	// check configuration repository given
	if tokenConfig == nil {
		return nil, err.New(ErrNoTokenConfiguration)
	}

	if manager == nil {
		return nil, err.New(ErrNoManager)
	}

	//overriding the nil value of ErrorTrace
	if trace == nil {
		trace = new(infrastructure.ErrorTraceNil)
	}

	//overriding the nil value of Logger
	if logger == nil {
		logger = new(infrastructure.LoggerNil)
	}

	var ts = &Oauth2Server{
		srv: server.NewServer(
			&server.Config{
				TokenType:             tokenConfig.Get().Type,
				AllowedResponseTypes:  AllowedResponseTypes,
				AllowedGrantTypes:     AllowedGrantTypes,
				AllowGetAccessRequest: AllowGetAccessRequest,
			}, manager),
		logger:         logger,
		trace:          trace,
		contextTimeout: contextTimeout,
		userDB:         user.NewUserDB(db),
	}
	ts.srv.SetPasswordAuthorizationHandler(ts.passwordCheck)
	ts.srv.SetUserAuthorizationHandler(ts.userAuthorize)
	ts.srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		ts.logger.Println("Internal Error:", err.Error())
		return
	})
	ts.srv.SetResponseErrorHandler(func(re *errors.Response) {
		ts.logger.Println("Response Error:", re.Error.Error())
	})
	return ts, nil
}

func (ts *Oauth2Server) passwordCheck(
	username, password string,
) (string, error) {
	var (
		intUserID int32
		err       error
	)
	ctx, cancel := context.WithTimeout(
		context.Background(), ts.contextTimeout,
	)
	defer cancel()
	intUserID, err = ts.userDB.CheckNamePassword(
		ctx,
		username,
		password,
	)
	if intUserID == 0 || err != nil {
		return "", ts.trace.New(ErrUserNotFound)
	}

	stringUserID := utils.String32(intUserID)
	ts.logger.Println("userID", stringUserID, intUserID)
	return stringUserID, nil
}

func (ts *Oauth2Server) userAuthorize(
	w http.ResponseWriter,
	r *http.Request,
) (userID string, err error) {
	ts.logger.Println("/userAuthorizeHandler")
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

func (ts *Oauth2Server) ValidationBearerToken(
	r *http.Request,
) (ti oauth2.TokenInfo, err error) {
	return ts.ValidationBearerToken(r)
}

func (ts *Oauth2Server) HandleAuthorizeRequest(
	w http.ResponseWriter,
	r *http.Request,
) error {
	return ts.srv.HandleAuthorizeRequest(w, r)
}

func (ts *Oauth2Server) HandleTokenRequest(
	w http.ResponseWriter,
	r *http.Request,
) error {
	return ts.srv.HandleTokenRequest(w, r)
}
