package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/gorilla/mux"
)

type UserI interface {
	GetName() string
	GetPassword() string
	SetName(string)
	SetPassword(string)
}

/*
IDFromPath get a parameter from a query path
name - name of parameter. For example in user/{user_id} name would be
"user_id"
if cant convert string to int, return error
*/
func (h *Handler) IDFromPath(
	r *http.Request,
	name string,
) (int32, error) {
	str := mux.Vars(r)[name]
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	if val < 0 {
		return 0, h.trace.New(ErrInvalidID)
	}
	return int32(val), nil
}

/*
 IDsFromPath get several parameters from path via IDFromPath func
 returns a map of ids, whose keys are parameter names

 if at least one value cannot be converted from a string to int32, then
 return error

 if 0 names given, no error return
*/
func (h *Handler) IDsFromPath(
	r *http.Request,
	names ...string,
) (map[string]int32, error) {
	var (
		ids = make(map[string]int32)
		err error
	)

	if len(names) == 0 {
		return ids, nil
	}
	for _, name := range names {
		ids[name], err = h.IDFromPath(r, name)
		if err != nil {
			break
		}
	}
	return ids, err
}

func (h *Handler) RequestParamsInt32(
	r *http.Request,
	withAuth bool,
	names ...string,
) (map[string]int32, error) {
	values, err := h.IDsFromPath(r, names...)
	if err == nil && withAuth {
		values[UserIDKey], err = h.GetUserIDFromAuthRequest(
			r,
		)
	}
	return values, err
}

func (h *Handler) StringFromPath(
	r *http.Request,
	name, defaultValue string,
) string {
	str := defaultValue
	vals := r.URL.Query()
	keys, ok := vals[name]
	if ok {
		if len(keys) >= 1 {
			str = keys[0]
		}
	}
	return str
}

func (h *Handler) IntFromPath(
	r *http.Request,
	name string,
	defaultVelue int,
	expected error,
) (val int, err error) {
	var str string
	if str = h.StringFromPath(r, name, ""); str == "" {
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

func (h *Handler) GetUserIDFromAuthRequest(
	r *http.Request,
) (int32, error) {
	interf := r.Context().Value(ContextUserKey)
	if interf != nil {
		i, err := strconv.Atoi(interf.(string))
		return int32(i), err
	}
	return 0, h.trace.New(ErrNoAuthFound)
}

func (h *Handler) ModelFromRequest(
	r *http.Request,
	jt models.JSONtype,
) error {
	if r.Body == nil {
		return h.trace.New(ErrNoBody)
	}
	defer r.Body.Close()

	bytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return err
	}
	return jt.UnmarshalJSON(bytes)
}

func (h *Handler) GetUser(
	r *http.Request,
	hashPassword func(string) string,
	ui UserI,
) error {

	if r.Body == nil {
		return h.trace.New(ErrNoBody)
	}
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&ui)
	if err != nil {
		return h.trace.New(ErrInvalidJSON)
	}
	ui.SetPassword(hashPassword(ui.GetPassword()))

	return nil
}

func (h *Handler) GetUserWithAllFields(
	r *http.Request,
	hashPassword func(string) string,
	ui UserI,
) error {

	err := h.GetUser(r, hashPassword, ui)
	if err != nil {
		return err
	}
	if ui.GetName() == "" {
		return h.trace.New(ErrInvalidName)
	}
	if ui.GetPassword() == "" {
		return h.trace.New(ErrInvalidPassword)
	}

	return nil
}

func (h *Handler) ValidateUser(
	user UserI,
) error {
	name := strings.TrimSpace(user.GetName())
	if name == "" || len(name) < MinNameLength {
		return h.trace.New(ErrInvalidName)
	}
	user.SetName(name)

	password := strings.TrimSpace(user.GetPassword())
	if len(password) < MinPasswordLength {
		return h.trace.New(ErrInvalidPassword)
	}
	return nil
}

// 206 -> 183
