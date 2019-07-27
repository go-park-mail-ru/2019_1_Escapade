package api

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"

	"net/http"
)

// GetUsersPageAmount get amount of users list page
// @Summary amount of users list page
// @Description Get amount of users list page
// @ID GetUsersPageAmount
// @Success 200 {object} models.Pages "Get successfully"
// @Failure 500 {object} models.Result "Server error"
// @Router /users/pages_amount [GET]
func (h *Handler) GetUsersPageAmount(rw http.ResponseWriter, r *http.Request) {
	const place = "GetUsersPageAmount"

	var (
		perPage int
		pages   models.Pages
		err     error
	)

	perPage = h.getPerPage(r)

	if pages.Amount, err = h.DB.GetUsersPageAmount(perPage); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.SendErrorJSON(rw, re.ErrorDataBase(), place)
		utils.PrintResult(err, http.StatusInternalServerError, place)
		return
	}

	utils.SendSuccessJSON(rw, pages, place)
	utils.PrintResult(err, http.StatusOK, place)
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
func (h *Handler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	const place = "GetUsers"
	var (
		err       error
		users     []*models.UserPublicInfo
		page      int
		perPage   int
		difficult int
		sort      string
	)

	sort = h.getSort(r)
	perPage = h.getPerPage(r)
	page = h.getPage(r)
	difficult = h.getDifficult(r)

	if users, err = h.DB.GetUsers(difficult, page, perPage, sort); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		utils.SendErrorJSON(rw, re.ErrorUsersNotFound(), place)
		utils.PrintResult(err, http.StatusNotFound, place)
		return
	}

	photo.GetImages(users...)

	utils.SendSuccessJSON(rw, models.UsersPublicInfo{users}, place)
	utils.PrintResult(err, http.StatusOK, place)
}
