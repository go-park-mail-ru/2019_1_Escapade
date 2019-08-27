package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/auth"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"

	"github.com/gorilla/mux"
)

type ContextKey string

const ContextUserKey ContextKey = "userID"

type UserI interface {
	GetName() string
	GetPassword() string
	SetName(string)
	SetPassword(string)
}

func IDFromPath(r *http.Request, name string) (int32, error) {
	var (
		str string
		val int
		err error
	)
	if str = GetStringFromPath(r, name, ""); str == "" {
		return 0, re.ID()
	}

	if val, err = strconv.Atoi(str); err != nil {
		return 0, err
	}
	if val < 0 {
		return 0, re.ID()
	}
	return int32(val), nil
}

func GetStringFromPath(r *http.Request, name string, defaultValue string) (str string) {
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
	if str = GetStringFromPath(r, name, ""); str == "" {
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

	return GetStringFromPath(r, "getStringFromPath", "time")
}

func (h *Handler) getName(r *http.Request) (username string, err error) {

	vars := mux.Vars(r)

	if username = vars["name"]; username == "" {
		return "", re.ErrorInvalidName()
	}

	return
}

func (h *Handler) getNameAndPage(r *http.Request) (page int, username string, err error) {
	vars := mux.Vars(r)

	if username = vars["name"]; username == "" {
		return 0, "", re.ErrorInvalidName()
	}

	if vars["page"] == "" {
		page = 1
	} else {
		if page, err = strconv.Atoi(vars["page"]); err != nil {
			return 0, username, re.ErrorInvalidPage()
		}
		if page < 1 {
			page = 1
		}
	}
	return
}

func GetUserIDFromAuthRequest(r *http.Request) (int32, error) {

	interf := r.Context().Value(ContextUserKey)
	if interf != nil {
		s := r.Context().Value(ContextUserKey).(string)
		i, err := strconv.Atoi(s)
		return int32(i), err
	}
	return 0, re.NoAuthFound()
}

func GetUser(r *http.Request, salt string, ui UserI) error {

	if r.Body == nil {
		return re.ErrorNoBody()
	}
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&ui)
	if err != nil {
		return re.ErrorInvalidJSON()
	}
	ui.SetPassword(auth.HashPassword(ui.GetPassword(), salt))

	return nil
}

// NEW - other deprecated
func ModelFromRequest(r *http.Request, jt JSONtype) error {

	if r.Body == nil {
		return re.ErrorNoBody()
	}
	defer r.Body.Close()

	bytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return err
	}
	return jt.UnmarshalJSON(bytes)
}

func getRecord(r *http.Request) (record models.Record, err error) {

	if r.Body == nil {
		return models.Record{}, re.ErrorNoBody()
	}
	defer r.Body.Close()

	if err = json.NewDecoder(r.Body).Decode(&record); err != nil {
		err = re.ErrorInvalidJSON()
	}

	return
}

func getGameInformation(r *http.Request) (info *models.GameInformation, err error) {

	if r.Body == nil {
		return nil, re.ErrorNoBody()
	}
	defer r.Body.Close()

	info = &models.GameInformation{}
	if err = json.NewDecoder(r.Body).Decode(info); err != nil {
		err = re.ErrorInvalidJSON()
	}

	return
}

func GetUserWithAllFields(r *http.Request, salt string, ui UserI) error {

	if err := GetUser(r, salt, ui); err != nil {
		return err
	}
	if ui.GetName() == "" {
		return re.ErrorInvalidName()
	}
	if ui.GetPassword() == "" {
		return re.ErrorInvalidPassword()
	}

	return nil
}

func ValidateUser(user UserI) error {
	name := strings.TrimSpace(user.GetName())
	if name == "" || len(name) < 3 {
		return re.ErrorInvalidName()
	}
	user.SetName(name)

	password := strings.TrimSpace(user.GetPassword())
	if len(password) < 3 {
		return re.ErrorInvalidPassword()
	}
	return nil
}

// 206 -> 183
