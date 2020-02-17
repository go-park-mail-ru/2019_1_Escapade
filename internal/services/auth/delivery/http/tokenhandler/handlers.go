package tokenhandler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"github.com/go-session/session"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/base/handler"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	ainfrastructure "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/infrastructure"
)

type TokenHandler struct {
	*handler.Handler

	server ainfrastructure.TokenServer
	store  ainfrastructure.TokenStore
	trace  infrastructure.ErrorTrace
	logger infrastructure.Logger
}

func New(
	server ainfrastructure.TokenServer,
	store ainfrastructure.TokenStore,
	trace infrastructure.ErrorTrace,
	logger infrastructure.Logger,
) *TokenHandler {
	return &TokenHandler{
		Handler: handler.New(logger, trace),

		server: server,
		store:  store,
		trace:  trace,
		logger: logger,
	}
}

func (th *TokenHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) models.RequestResult {
	th.logger.Println("/delete")
	token, err := th.server.ValidationBearerToken(r)
	if err != nil {
		th.logger.Println("error", err.Error())
		return th.Fail(http.StatusBadRequest, err)
	}
	if err = th.store.RemoveByAccess(token.GetAccess()); err != nil {
		th.logger.Println("error", err.Error())
		return th.Fail(http.StatusBadRequest, err)
	}

	e := json.NewEncoder(w)
	e.Encode(token)
	return th.Success(http.StatusOK, nil)
}

func (th *TokenHandler) Test(
	w http.ResponseWriter,
	r *http.Request,
) models.RequestResult {
	th.logger.Println("/test")
	token, err := th.server.ValidationBearerToken(r)
	if err != nil {
		th.logger.Println("error", err.Error())
		return th.Fail(http.StatusBadRequest, err)
	}
	th.logger.Println("token", token.GetAccessExpiresIn())

	e := json.NewEncoder(w)
	e.Encode(token)
	th.logger.Println("/dont know")
	return th.Success(http.StatusOK, nil)
}

// tokenHandler handle getting token
// @Summary Get token
// @Description Get session token for current client
// @ID tokenHandler
// @Tags auth
// @Accept  json
// @Produce  json
// @Param grant_type body string true "'password' or 'refresh_token'" default("password")
// @Param client_id body string true "client id" default("1")
// @Param client_secret body string true "client secret" default("1")
// @Param username body string false "username" default("username")
// @Param password body string false "password" default("password")
// @Param refresh_token body string false "client id" default("1")
// @Router /token [POST]
func (th *TokenHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) models.RequestResult {
	err := th.server.HandleTokenRequest(w, r)
	if err != nil {
		return th.Fail(http.StatusInternalServerError, err)
	}
	return th.Success(http.StatusOK, nil)
}

func (th *TokenHandler) Authorize(
	w http.ResponseWriter,
	r *http.Request,
) models.RequestResult {
	th.logger.Println("/authorize")
	store, err := session.Start(nil, w, r)
	if err != nil {
		return th.Fail(http.StatusInternalServerError, err)
	}

	var form url.Values
	if v, ok := store.Get("ReturnUri"); ok {
		form = v.(url.Values)
	}
	r.Form = form

	store.Delete("ReturnUri")
	store.Save()

	err = th.server.HandleAuthorizeRequest(w, r)
	if err != nil {
		return th.Fail(http.StatusBadRequest, err)
	}
	return th.Success(http.StatusOK, nil)
}

func (th *TokenHandler) Login(
	w http.ResponseWriter,
	r *http.Request,
) models.RequestResult {
	th.logger.Println("/loginHandler")
	store, err := session.Start(nil, w, r)
	if err != nil {
		return th.Fail(http.StatusInternalServerError, err)
	}

	if r.Method == "POST" {
		store.Set("LoggedInUserID", r.Form.Get("username"))
		store.Save()

		w.Header().Set("Location", "/auth")
		return th.Success(http.StatusFound, nil)
	}
	return th.outputHTML(w, r, "static/login.html")
}

func (th *TokenHandler) Auth(
	w http.ResponseWriter,
	r *http.Request,
) models.RequestResult {
	th.logger.Println("/Auth")
	store, err := session.Start(nil, w, r)
	if err != nil {
		return th.Fail(http.StatusInternalServerError, err)
	}

	if _, ok := store.Get("LoggedInUserID"); !ok {
		w.Header().Set("Location", "/login")
		return th.Success(http.StatusFound, nil)
	}

	return th.outputHTML(w, r, "static/auth.html")
}

func (th *TokenHandler) outputHTML(
	w http.ResponseWriter,
	req *http.Request,
	filename string,
) models.RequestResult {
	file, err := os.Open(filename)
	if err != nil {
		return th.Fail(http.StatusInternalServerError, err)
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
	return th.Success(handler.NoResult, nil)
}
