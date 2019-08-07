package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	me "github.com/go-park-mail-ru/2019_1_Escapade/internal/middleware"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"github.com/gorilla/mux"
)

func getStringFromPath(r *http.Request, name string, defaultValue string) (str string) {
	str = defaultValue

	vals := r.URL.Query()
	keys, ok := vals[name]
	if ok {
		if len(keys) >= 1 {
			str = keys[0]
		}
	}
	return
}

func getIntFromPath(r *http.Request, name string,
	defaultVelue int, expected error) (val int, err error) {
	var str string
	if str = getStringFromPath(r, name, ""); str == "" {
		err = expected
		return
	}
	val = defaultVelue

	if val, err = strconv.Atoi(str); err != nil {
		err = expected
		return
	}
	if val < 0 {
		err = expected
		return
	}
	return
}

func (h *Handler) getUserID(r *http.Request) (id int, err error) {

	id, err = getIntFromPath(r, "id", 1, re.ErrorInvalidUserID())
	return
}

func (h *Handler) getPage(r *http.Request) int {

	page, _ := getIntFromPath(r, "page", 1, nil)
	return page
}

func (h *Handler) getPerPage(r *http.Request) int {

	page, _ := getIntFromPath(r, "per_page", 100, nil)
	return page
}

func (h *Handler) getDifficult(r *http.Request) int {

	diff, _ := getIntFromPath(r, "difficult", 0, nil)
	if diff > 3 {
		diff = 3
	}
	return diff
}

func (h *Handler) getSort(r *http.Request) string {

	return getStringFromPath(r, "getStringFromPath", "time")
}

func (h *Handler) getName(r *http.Request) (username string, err error) {
	var (
		vars map[string]string
	)

	vars = mux.Vars(r)

	if username = vars["name"]; username == "" {
		err = re.ErrorInvalidName()
		return
	}

	return
}

func (h *Handler) getNameAndPage(r *http.Request) (page int, username string, err error) {
	var (
		vars map[string]string
	)

	vars = mux.Vars(r)

	if username = vars["name"]; username == "" {
		err = re.ErrorInvalidName()
		return
	}

	if vars["page"] == "" {
		page = 1
	} else {
		if page, err = strconv.Atoi(vars["page"]); err != nil {
			err = re.ErrorInvalidPage()
			return
		}
		if page < 1 {
			page = 1
		}

	}
	return
}

func (h *Handler) getUserIDFromAuthRequest(r *http.Request) (int32, error) {

	interf := r.Context().Value(me.ContextUserKey)
	if interf != nil {
		s := r.Context().Value(me.ContextUserKey).(string)
		i, err := strconv.Atoi(s)
		utils.Debug(false, "id found as", i)
		return int32(i), err
	}
	utils.Debug(false, "id not found")
	return 0, re.NoAuthFound()
}

func getUser(r *http.Request) (user models.UserPrivateInfo, err error) {

	if r.Body == nil {
		err = re.ErrorNoBody()

		return
	}
	defer r.Body.Close()

	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		err = re.ErrorInvalidJSON()
	}

	return
}

func getRecord(r *http.Request) (record models.Record, err error) {

	if r.Body == nil {
		err = re.ErrorNoBody()

		return
	}
	defer r.Body.Close()

	if err = json.NewDecoder(r.Body).Decode(&record); err != nil {
		err = re.ErrorInvalidJSON()
	}

	return
}

func getGameInformation(r *http.Request) (info *models.GameInformation, err error) {

	if r.Body == nil {
		err = re.ErrorNoBody()

		return
	}
	defer r.Body.Close()

	info = &models.GameInformation{}
	if err = json.NewDecoder(r.Body).Decode(info); err != nil {
		err = re.ErrorInvalidJSON()
	}

	return
}

func getUserWithAllFields(r *http.Request) (user models.UserPrivateInfo, err error) {

	if user, err = getUser(r); err != nil {
		return
	}
	if user.Name == "" {
		err = re.ErrorInvalidName()
		return
	}

	if user.Password == "" {
		err = re.ErrorInvalidPassword()
		return
	}

	return
}
