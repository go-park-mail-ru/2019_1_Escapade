package handlers

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	idb "github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"

	ih "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"

	"net/http"
)

type UsersHandler struct {
	ih.Handler
	user database.UserUseCaseI
}

func (h *UsersHandler) Init(c *config.Configuration, DB idb.DatabaseI,
	userDB database.UserRepositoryI, recordDB database.RecordRepositoryI) error {
	h.Handler.Init(c)

	h.user = &database.UserUseCase{}
	h.user.Init(userDB, recordDB)
	err := h.user.Use(DB)
	if err != nil {
		return err
	}
	return nil
}

func (h *UsersHandler) Close() {
	h.user.Close()
}

// HandleUsersPages process any operation associated with users
// list: receive
func (h *UsersHandler) HandleUsersPages(rw http.ResponseWriter, r *http.Request) {
	ih.Route(rw, r, ih.MethodHandlers{
		http.MethodGet:     h.GetUsers,
		http.MethodOptions: nil})
}

// HandleUsersPageAmount process any operation associated with
// amount of pages in user list: receive
func (h *UsersHandler) HandleUsersPageAmount(rw http.ResponseWriter, r *http.Request) {
	ih.Route(rw, r, ih.MethodHandlers{
		http.MethodGet:     h.GetUsersPageAmount,
		http.MethodOptions: nil})
}

// GetUsersPageAmount get amount of users list page
// @Summary amount of users list page
// @Description Get amount of users list page
// @ID GetUsersPageAmount
// @Success 200 {object} models.Pages "Get successfully"
// @Failure 500 {object} models.Result "Server error"
// @Router /users/pages_amount [GET]
func (h *UsersHandler) GetUsersPageAmount(rw http.ResponseWriter, r *http.Request) ih.Result {
	const place = "GetUsersPageAmount"

	var (
		perPage int
		pages   models.Pages
		err     error
	)

	perPage = getPerPage(r)

	if pages.Amount, err = h.user.PagesCount(perPage); err != nil {
		return ih.NewResult(http.StatusInternalServerError, place, nil, re.DatabaseWrapper(err))
	}

	return ih.NewResult(http.StatusOK, place, &pages, nil)
}

// GetUsers get users list
// @Summary Get users list
// @Description Get page of user list
// @ID GetUsers
// @Success 200 {array} models.Result "Get successfully"
// @Failure 400 {object} models.Result "Invalid pade"
// @Failure 404 {object} models.Result "Users not found"
// @Failure 500 {object} models.Result "Server error"
// @Router /users/{page} [GET]
func (h *UsersHandler) GetUsers(rw http.ResponseWriter, r *http.Request) ih.Result {
	const place = "GetUsers"
	var (
		err       error
		users     []*models.UserPublicInfo
		page      int
		perPage   int
		difficult int
		sort      string
	)

	sort = getSort(r)
	perPage = getPerPage(r)
	page = getPage(r)
	difficult = getDifficult(r)

	if users, err = h.user.FetchAll(difficult, page, perPage, sort); err != nil {
		return ih.NewResult(http.StatusNotFound, place, nil, re.NoUserWrapper(err))
	}

	photo.GetImages(users...)

	return ih.NewResult(http.StatusOK, place, &models.UsersPublicInfo{users}, nil)
}
