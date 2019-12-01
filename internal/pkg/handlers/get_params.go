package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/auth"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
)

type ContextKey string

const ContextUserKey ContextKey = "userID"

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
func IDFromPath(r *http.Request, name string) (int32, error) {
	str := mux.Vars(r)[name]
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	if val < 0 {
		return 0, re.ID()
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
func IDsFromPath(r *http.Request, names ...string) (map[string]int32, error) {
	var (
		ids = make(map[string]int32)
		err error
	)

	if len(names) == 0 {
		return ids, nil
	}
	for _, name := range names {
		ids[name], err = IDFromPath(r, name)
		if err != nil {
			break
		}
	}
	return ids, err
}

/*
RequestParamsInt32 get all parameters from path via IDsFromPath and UserID
via GetUserIDFromAuthRequest(if 'withAuth' true)
userID is placed in map with key set in UserIDKey
*/
const UserIDKey = "auth_user_id"

func RequestParamsInt32(r *http.Request, withAuth bool, names ...string) (map[string]int32, error) {
	values, err := IDsFromPath(r, names...)
	if err == nil && withAuth {
		values[UserIDKey], err = GetUserIDFromAuthRequest(r)
	}
	return values, err
}

func StringFromPath(r *http.Request, name string, defaultValue string) (str string) {
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

func IntFromPath(r *http.Request, name string,
	defaultVelue int, expected error) (val int, err error) {
	var str string
	if str = StringFromPath(r, name, ""); str == "" {
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

func GetUserIDFromAuthRequest(r *http.Request) (int32, error) {

	interf := r.Context().Value(ContextUserKey)
	if interf != nil {
		i, err := strconv.Atoi(interf.(string))
		return int32(i), err
	}
	return 0, re.NoAuthFound()
}

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
