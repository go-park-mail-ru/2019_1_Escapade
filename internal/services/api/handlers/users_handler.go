package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	ih "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/photo"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database"
)

// UsersHandler handle requests associated with list of users
type UsersHandler struct {
	ih.Handler
	user database.UserUseCaseI
}

// Init open connections to database
func (h *UsersHandler) Init(c *config.Configuration, db *database.Input) *UsersHandler {
	h.Handler.Init(c)
	h.user = db.UserUC
	return h
}

// Close connections to database
func (h *UsersHandler) Close() error {
	return h.user.Close()
}

// HandleUsersPages process any operation associated with users
// list: receive
func (h *UsersHandler) HandleUsersPages(rw http.ResponseWriter, r *http.Request) {
	page, ok := r.URL.Query()["page"]
	if !ok {
		ih.SendError(http.StatusBadRequest, "HandleUsersPages", re.ID())
		return
	}
	ih.Route(rw, r, ih.MethodHandlers{
		http.MethodGet: h.GetUsers(r.FormValue("per_page"),
			page[0], getDifficult(r), getSort(r)),
		http.MethodOptions: nil})
}

// HandleUsersPageAmount process any operation associated with
// amount of pages in user list: receive
func (h *UsersHandler) HandleUsersPageAmount(rw http.ResponseWriter, r *http.Request) {
	ih.Route(rw, r, ih.MethodHandlers{
		http.MethodGet:     h.GetUsersPageAmount(r.FormValue("per_page")),
		http.MethodOptions: nil})
}

func (h *UsersHandler) HandleGetProfile(rw http.ResponseWriter, r *http.Request) {
	ih.Route(rw, r, ih.MethodHandlers{
		http.MethodGet:     h.GetProfile(),
		http.MethodOptions: nil})
}

func (h *UsersHandler) GetProfile() ih.HandleRequest {
	return func(rw http.ResponseWriter, r *http.Request) ih.Result {
		userID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			userID, err = getUserID(r)
			if err != nil {
				return ih.NewResult(http.StatusBadRequest, "GetProfile",
					nil, re.NoUserWrapper(err))
			}
		}
		return h.getUser(rw, r, int32(userID), "GetProfile")
	}
}

func (h *UsersHandler) GetUsersPageAmount(perPage string) ih.HandleRequest {
	return func(rw http.ResponseWriter, r *http.Request) ih.Result {
		return h.getUsersPageAmount(perPage, rw, r)
	}
}

func (h *UsersHandler) GetUsers(perPage, page, diff, sort string) ih.HandleRequest {
	return func(rw http.ResponseWriter, r *http.Request) ih.Result {
		return h.getUsers(perPage, page, diff, sort, rw, r)
	}
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
func (h *UsersHandler) getUser(rw http.ResponseWriter, r *http.Request,
	userID int32, place string) ih.Result {

	difficult, err := strconv.Atoi(getDifficult(r))
	if err != nil {
		return ih.NewResult(http.StatusBadRequest, place, nil, re.NoUserWrapper(err))
	}

	user, err := h.user.FetchOne(userID, difficult)
	if err != nil {
		return ih.NewResult(http.StatusNotFound, place, nil, re.NoUserWrapper(err))
	}

	photo.GetImages(user)

	return ih.NewResult(http.StatusOK, place, user, nil)
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
func (h *UsersHandler) getUsersPageAmount(perPage string, rw http.ResponseWriter, r *http.Request) ih.Result {
	const place = "GetUsersPageAmount"
	var pages = models.Pages{}

	perPageI, err := strconv.Atoi(perPage)
	if err != nil {
		return ih.NewResult(http.StatusBadRequest, place, nil, re.NoUserWrapper(err))
	}
	pages.Amount, err = h.user.PagesCount(perPageI)
	if err != nil {
		return ih.NewResult(http.StatusInternalServerError, place, nil, re.DatabaseWrapper(err))
	}
	return ih.NewResult(http.StatusOK, place, &pages, nil)
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
func (h *UsersHandler) getUsers(perPage, page, difficult, sort string, rw http.ResponseWriter, r *http.Request) ih.Result {
	const place = "GetUsers"
	var (
		err   error
		users []*models.UserPublicInfo
	)

	difficultI, err := strconv.Atoi(difficult)
	if err != nil {
		return ih.NewResult(http.StatusBadRequest, place, nil, re.NoUserWrapper(err))
	}
	perPageI, err := strconv.Atoi(perPage)
	if err != nil {
		return ih.NewResult(http.StatusBadRequest, place, nil, re.NoUserWrapper(err))
	}
	pageI, err := strconv.Atoi(page)
	if err != nil {
		return ih.NewResult(http.StatusBadRequest, place, nil, re.NoUserWrapper(err))
	}

	if users, err = h.user.FetchAll(difficultI, pageI, perPageI, sort); err != nil {
		return ih.NewResult(http.StatusNotFound, place, nil, re.NoUserWrapper(err))
	}

	photo.GetImages(users...)

	return ih.NewResult(http.StatusOK, place, &models.UsersPublicInfo{users}, nil)
}
