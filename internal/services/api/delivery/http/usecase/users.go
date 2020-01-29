package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/domens/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	ih "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api"
	delivery "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/delivery/http"
)

// UsersHandler handle requests associated with list of users
type UsersHandler struct {
	user api.UserUseCaseI
	rep  delivery.RepositoryI

	photo infrastructure.PhotoServiceI
	trace infrastructure.ErrorTrace
}

func NewUsersHandler(
	user api.UserUseCaseI,
	rep delivery.RepositoryI,
	photo infrastructure.PhotoServiceI,
	trace infrastructure.ErrorTrace,
) *UsersHandler {
	return &UsersHandler{
		user: user,
		rep:  rep,

		photo: photo,
		trace: trace,
	}
}

func (h *UsersHandler) GetOneUser(
	rw http.ResponseWriter,
	r *http.Request,
) ih.Result {
	userID, err := h.rep.GetUserID(r)
	if err != nil {
		return ih.NewResult(
			http.StatusBadRequest,
			nil,
			h.trace.WrapWithText(err, ErrUserNotFound),
		)
	}
	return h.GetUser(rw, r, int32(userID))
}

// GetMyProfile get user public information
// @Summary get user public information
// @Description  get user's best score and best time for a given difficulty, user's id, name and photo of current user. The current one is the one whose token is provided.
// @ID getProfile
// @Tags users
// @Accept json
// @Param id path int true "user's id"
// @Produce  json
// @Success 200 {object} models.UserPublicInfo "Get user successfully"
// @Failure 400 {object} models.Result "Wrong input data"
// @Failure 404 {object} models.Result "Not found"
// @Router /users/{id} [GET]
func (h *UsersHandler) GetUser(
	rw http.ResponseWriter,
	r *http.Request,
	userID int32,
) ih.Result {

	difficult, err := strconv.Atoi(h.rep.GetDifficult(r))
	if err != nil {
		return ih.NewResult(
			http.StatusBadRequest,
			nil,
			h.trace.WrapWithText(err, ErrUserNotFound),
		)
	}

	user, err := h.user.FetchOne(
		r.Context(),
		userID,
		difficult,
	)
	if err != nil {
		return ih.NewResult(
			http.StatusNotFound,
			nil,
			h.trace.WrapWithText(err, ErrUserNotFound),
		)
	}

	h.photo.GetImages(user)

	return ih.NewResult(http.StatusOK, user, nil)
}

// GetUsersPageAmount get number of pages with users
// @Summary get number of pages with users
// @Description You pass how many users should be placed on one page, and in return you get how many pages with users you can get.
// @ID GetUsersPageAmount
// @Tags users
// @Accept json
// @Param id query int true "number of users in one page"
// @Produce  json
// @Success 200 {object} models.Pages "Get successfully"
// @Failure 400 {object} models.Result "Invalid path parameter"
// @Failure 500 {object} models.Result "Database error"
// @Router /users/pages/amount [GET]
func (h *UsersHandler) GetUsersPageAmount(
	rw http.ResponseWriter,
	r *http.Request,
) ih.Result {
	var pages = models.Pages{}

	perPage := h.rep.GetPerPage(r)
	perPageI, err := strconv.Atoi(perPage)
	if err != nil {
		return ih.NewResult(
			http.StatusBadRequest,
			nil,
			h.trace.WrapWithText(err, ErrUserNotFound),
		)
	}
	pages.Amount, err = h.user.PagesCount(
		r.Context(),
		perPageI,
	)
	if err != nil {
		return ih.NewResult(
			http.StatusInternalServerError,
			nil,
			h.trace.WrapWithText(err, ErrFailedPageCountGet),
		)
	}
	return ih.NewResult(http.StatusOK, &pages, nil)
}

// GetUsers get users list
// @Summary Get users list
// @Description Get one page of users with selected size.
// @ID GetUsers
// @Tags users
// @Accept json
// @Param page path int true "the offset of users list" default(0)
// @Param per_page query int true "the limit of users page"" default(0)
// @Param difficult query int false "which difficult records will be given" default(0)
// @Param sort query string false "sort list by 'score' or by 'time'" default("time")
// @Produce  json
// @Success 200 {array} models.UserPublicInfo "Get successfully"
// @Failure 400 {object} models.Result "Invalid pade"
// @Failure 404 {object} models.Result "Users not found"
// @Failure 500 {object} models.Result "Server error"
// @Router /users/pages/{page} [GET]
func (h *UsersHandler) GetUsers(
	w http.ResponseWriter,
	r *http.Request,
) ih.Result {
	var (
		err       error
		perPage   = h.rep.GetPerPage(r)
		page      = h.rep.GetPage(r)
		difficult = h.rep.GetDifficult(r)
		sort      = h.rep.GetSort(r)
	)

	difficultI, err := strconv.Atoi(difficult)
	if err != nil {
		return ih.NewResult(
			http.StatusBadRequest,
			nil,
			h.trace.WrapWithText(err, ErrUserNotFound),
		)
	}
	perPageI, err := strconv.Atoi(perPage)
	if err != nil {
		return ih.NewResult(
			http.StatusBadRequest,
			nil,
			h.trace.WrapWithText(err, ErrUserNotFound),
		)
	}
	pageI, err := strconv.Atoi(page)
	if err != nil {
		return ih.NewResult(
			http.StatusBadRequest,
			nil,
			h.trace.WrapWithText(err, ErrUserNotFound),
		)
	}

	users, err := h.user.FetchAll(
		r.Context(),
		difficultI,
		pageI,
		perPageI,
		sort,
	)
	if err != nil {
		return ih.NewResult(
			http.StatusNotFound,
			nil,
			h.trace.WrapWithText(err, ErrUserNotFound),
		)
	}

	h.photo.GetImages(users...)

	return ih.NewResult(
		http.StatusOK,
		&models.UsersPublicInfo{users},
		nil,
	)
}
