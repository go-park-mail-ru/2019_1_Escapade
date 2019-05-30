package api

import (
	//"github.com/go-park-mail-ru/2019_1_Escapade/internal/clients"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	cook "github.com/go-park-mail-ru/2019_1_Escapade/internal/cookie"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	//c "github.com/smartystreets/goconvey/convey"
	"testing"
)

type TestCase struct {
	Response    string
	Body        string
	CookieValue string
	StatusCode  int
}

type TestCaseGames struct {
	Response   string
	name       string
	page       string
	StatusCode int
}

// go test -coverprofile test/cover.out
// go tool cover -html=test/cover.out -o test/coverage.html

const PATH = "../../../conf.json"
const RANDOM = "r"
const DEFAULT = "d"

var TestAPI *Handler

var replacer = strings.NewReplacer("\n", "", "\t", "")

func compareStrings(got, expected string) bool {
	got = strings.TrimSpace(got)
	expected = strings.TrimSpace(expected)
	got = replacer.Replace(got)
	expected = replacer.Replace(expected)
	return got == expected
}

func launchTests(t *testing.T, H *Handler, cases []TestCase,
	url string, af apiFunc, compare Comparator, c *http.Cookie, record bool) {

	for caseNum, item := range cases {
		w := af(H, url, strings.NewReader(item.Body), c)
		if !record {
			continue
		}

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.StatusCode)
		}

		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			break
		}

		defer resp.Body.Close()

		if !compare(string(body), item.Response) {
			t.Errorf("[%d] wrong text: got %s, expected %s",
				caseNum, string(body), item.Response)
		}
	}
}

func getUserCreateCases() []TestCase {
	const place = "CreateUser"
	name1 := utils.RandomString(16)
	email1 := utils.RandomString(16)
	return []TestCase{
		TestCase{ // correct
			Response:   createResult(place, nil),
			Body:       createPrivateUser(name1, RANDOM, email1),
			StatusCode: http.StatusCreated,
		}, // email is taken
		TestCase{
			Response:   createResult(place, re.ErrorUserIsExist()),
			Body:       createPrivateUser(RANDOM, RANDOM, email1),
			StatusCode: http.StatusBadRequest,
		}, // name is taken
		TestCase{
			Response:   createResult(place, re.ErrorUserIsExist()),
			Body:       createPrivateUser(name1, RANDOM, RANDOM),
			StatusCode: http.StatusBadRequest,
		}, // no password
		TestCase{
			Response:   createResult(place, re.ErrorInvalidPassword()),
			Body:       createPrivateUser(RANDOM, "", RANDOM),
			StatusCode: http.StatusBadRequest,
		}, // no email
		TestCase{
			Response:   createResult(place, re.ErrorInvalidPassword()),
			Body:       createPrivateUser(RANDOM, RANDOM, ""),
			StatusCode: http.StatusBadRequest,
		}, // no name
		TestCase{
			Response:   createResult(place, re.ErrorInvalidName()),
			Body:       createPrivateUser("", RANDOM, RANDOM),
			StatusCode: http.StatusBadRequest,
		}, // no anything
		TestCase{
			Response:   createResult(place, re.ErrorInvalidName()),
			Body:       createPrivateUser("", "", ""),
			StatusCode: http.StatusBadRequest,
		},
		// wrong json
		TestCase{
			Response:   createResult(place, re.ErrorInvalidName()),
			Body:       createResult("", nil),
			StatusCode: http.StatusBadRequest,
		},
		// no json
		TestCase{
			Response:   createResult(place, re.ErrorInvalidJSON()),
			Body:       "",
			StatusCode: http.StatusBadRequest,
		},
	}
}

func getUserDeleteCases() []TestCase {
	const place = "DeleteUser"
	name1 := utils.RandomString(16)
	email1 := utils.RandomString(16)
	return []TestCase{
		TestCase{ // correct
			Response:   createResult(place, nil),
			Body:       createPrivateUser(name1, RANDOM, email1),
			StatusCode: http.StatusOK,
		}, // user almost deleted
		TestCase{
			Response:   createResult(place, re.ErrorUserNotFound()),
			Body:       createPrivateUser(RANDOM, RANDOM, email1),
			StatusCode: http.StatusBadRequest,
		}, // user almost deleted
		TestCase{
			Response:   createResult(place, re.ErrorUserNotFound()),
			Body:       createPrivateUser(name1, RANDOM, RANDOM),
			StatusCode: http.StatusBadRequest,
		}, // no password
		TestCase{
			Response:   createResult(place, re.ErrorInvalidPassword()),
			Body:       createPrivateUser(RANDOM, "", RANDOM),
			StatusCode: http.StatusBadRequest,
		}, // no email
		TestCase{
			Response:   createResult(place, re.ErrorInvalidPassword()),
			Body:       createPrivateUser(RANDOM, RANDOM, ""),
			StatusCode: http.StatusBadRequest,
		}, // no name
		TestCase{
			Response:   createResult(place, re.ErrorInvalidName()),
			Body:       createPrivateUser("", RANDOM, RANDOM),
			StatusCode: http.StatusBadRequest,
		}, // no anything
		TestCase{
			Response:   createResult(place, re.ErrorInvalidName()),
			Body:       createPrivateUser("", "", ""),
			StatusCode: http.StatusBadRequest,
		},
		// wrong json
		TestCase{
			Response:   createResult(place, re.ErrorInvalidName()),
			Body:       createResult("", nil),
			StatusCode: http.StatusBadRequest,
		},
		// no json
		TestCase{
			Response:   createResult(place, re.ErrorInvalidJSON()),
			Body:       "",
			StatusCode: http.StatusBadRequest,
		},
	}
}

// func TestCreateUserConvey(t *testing.T) {
// 	c.Convey("Create user", t, func() {
// 		url := "/user"
// 		cases := getUserCreateCases()
// 		c.Convey("Should create user", func() {
// 			launchTests(t, cases, url, createUser, nil, true)
// 			launchTests(t, cases, url, deleteUser, nil, false)
// 			c.So(nil, c.ShouldBeNil)
// 		})
// 	})

// }

func TestAll(t *testing.T) {

	const (
		place      = "main"
		confPath   = "../../conf.json"
		secretPath = "../../secret.json"
	)

	var (
		H             *Handler
		configuration *config.Configuration
		err           error
	)
	if configuration, err = config.InitPublic(confPath); err != nil {
		t.Error("eeeer", err.Error())
		return
	}

	// fmt.Println("launchTests")
	// authConn, err := clients.ServiceConnectionsInit(configuration.AuthClient)
	// if err != nil {
	// 	t.Error("serviceConnectionsInit error:", err)
	// }
	// defer authConn.Close()

	H, err = GetAPIHandler(configuration) // init.go
	if err != nil {
		t.Error("serviceConnectionsInit error:", err)
	}
	if H == nil {
		fmt.Println("launchTests failes")
	}

	TCreateUser(t, H)
	TDeleteUser(t, H)
	TUpdateProfile(t, H)
	TGetMyProfile(t, H)
	TGetProfile(t, H)
	TLogin(t, H)
	TLogout(t, H)
	t.Error("serviceConnectionsInit error:")
	// delete everything in database after tests
}

func TCreateUser(t *testing.T, H *Handler) {

	url := "/user"

	cases := getUserCreateCases()

	launchTests(t, H, cases, url, createUser, compareStrings, nil, true)

	launchTests(t, H, cases, url, deleteUser, compareStrings, nil, false)
}

func TDeleteUser(t *testing.T, H *Handler) {
	url := "/user"

	cases := getUserDeleteCases()

	launchTests(t, H, cases, url, createUser, compareStrings, nil, false)
	launchTests(t, H, cases, url, deleteUser, compareStrings, nil, true)
}

func TUpdateProfile(t *testing.T, H *Handler) {
	const place = "UpdateProfile"
	takenName := "takenName"
	takenEmail := "takenEmail"

	updates := []TestCase{
		TestCase{
			Response:   createResult(place, re.ErrorAuthorization()),
			Body:       createPrivateUser(RANDOM, RANDOM, RANDOM),
			StatusCode: http.StatusUnauthorized,
		},
		TestCase{
			Response:   createResult(place, nil),
			Body:       createPrivateUser(RANDOM, RANDOM, RANDOM),
			StatusCode: http.StatusOK,
		},
		TestCase{
			Response:   createResult(place, nil),
			Body:       createPrivateUser("", RANDOM, RANDOM),
			StatusCode: http.StatusOK,
		},
		TestCase{
			Response:   createResult(place, nil),
			Body:       createPrivateUser(RANDOM, "", RANDOM),
			StatusCode: http.StatusOK,
		},
		TestCase{
			Response:   createResult(place, nil),
			Body:       createPrivateUser(RANDOM, RANDOM, ""),
			StatusCode: http.StatusOK,
		},
		TestCase{
			Response:   createResult(place, nil),
			Body:       createPrivateUser("", "", RANDOM),
			StatusCode: http.StatusOK,
		},
		TestCase{
			Response:   createResult(place, nil),
			Body:       createPrivateUser("", RANDOM, ""),
			StatusCode: http.StatusOK,
		},
		TestCase{
			Response:   createResult(place, nil),
			Body:       createPrivateUser(RANDOM, "", ""),
			StatusCode: http.StatusOK,
		},
		TestCase{
			Response:   createResult(place, re.ErrorUserIsExist()),
			Body:       createPrivateUser(takenName, "", takenEmail),
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Response:   createResult(place, re.ErrorUserIsExist()),
			Body:       createPrivateUser(takenName, "", ""),
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Response:   createResult(place, re.ErrorUserIsExist()),
			Body:       createPrivateUser("", "", takenEmail),
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Response:   createResult(place, nil),
			Body:       createPrivateUser("", "", ""),
			StatusCode: http.StatusOK,
		},
		TestCase{
			Response:   createResult(place, re.ErrorInvalidJSON()),
			Body:       "not json",
			StatusCode: http.StatusBadRequest,
		},
	}

	var (
		cookie *http.Cookie
	)

	url := "/user"

	user := models.UserPrivateInfo{
		Name:     takenName,
		Password: utils.RandomString(16),
	}

	H.register(context.Background(), user)

	user = *createRandomUser()
	_, cookiestr, err := H.register(context.Background(), user)
	if err != nil {
		t.Error(" error:", err.Error())
		return
	}
	cookie = cook.CreateCookie(cookiestr, H.Session)

	launchTests(t, H, updates[:1], url, updateUser, compareStrings, nil, true)
	launchTests(t, H, updates[1:], url, updateUser, compareStrings, cookie, true)

}

func TGetMyProfile(t *testing.T, H *Handler) {
	const place = "GetMyProfile"
	name := utils.RandomString(16)
	email := utils.RandomString(16)

	gets := []TestCase{
		TestCase{
			Response:   createResult(place, re.ErrorAuthorization()),
			Body:       createPrivateUser(RANDOM, RANDOM, RANDOM),
			StatusCode: http.StatusUnauthorized,
		},
		TestCase{
			Response: createPublicUser(name, email, "",
				DEFAULT, DEFAULT, DEFAULT, DEFAULT, DEFAULT),
			StatusCode: http.StatusOK,
		},
	}

	url := "/user"

	user := models.UserPrivateInfo{
		Name: name,
	}
	_, cookiestr, err := H.register(context.Background(), user)
	if err != nil {
		t.Error("TestUpdateUser catched error:", err.Error())
		return
	}
	cookie := cook.CreateCookie(cookiestr, H.Session)

	launchTests(t, H, gets[:1], url, getMyProfile, models.ComparePublicUsers, nil, true)
	launchTests(t, H, gets[1:], url, getMyProfile, models.ComparePublicUsers, cookie, true)

	if err = H.deleteAccount(context.Background(), &user, cookiestr); err != nil {
		t.Error("TestUpdateUser catched error:", err.Error())
		return
	}
}

func TGetProfile(t *testing.T, H *Handler) {
	const place = "GetProfile"
	name := utils.RandomString(16)
	email := utils.RandomString(16)

	gets := []TestCase{
		TestCase{
			Response:   createResult(place, re.ErrorInvalidUserID()),
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Response:   createResult(place, re.ErrorUserNotFound()),
			StatusCode: http.StatusNotFound,
		},
		TestCase{
			Response: createPublicUser(name, email, "",
				DEFAULT, DEFAULT, DEFAULT, DEFAULT, DEFAULT),
			StatusCode: http.StatusOK,
		},
	}

	user := models.UserPrivateInfo{
		Name:     name,
		Password: utils.RandomString(16),
	}

	id, _, err := H.register(context.Background(), user)
	if err != nil {
		t.Error(" error:", err.Error())
		return
	}

	url := "/user"
	launchTests(t, H, gets[:1], url, getProfile, models.ComparePublicUsers, nil, true)
	url = "/user?id=dfdf"
	launchTests(t, H, gets[:1], url, getProfile, models.ComparePublicUsers, nil, true)
	url = "/user?id=100000"
	launchTests(t, H, gets[1:2], url, getProfile, models.ComparePublicUsers, nil, true)
	url = "/user?id=" + strconv.FormatInt(int64(id), 10)
	launchTests(t, H, gets[2:], url, getProfile, models.ComparePublicUsers, nil, true)
}

func TLogin(t *testing.T, H *Handler) {
	const place = "Login"
	name := utils.RandomString(16)
	password := utils.RandomString(16)
	email := utils.RandomString(16)

	cases := []TestCase{
		TestCase{
			Response: createPublicUser(name, email, "",
				DEFAULT, DEFAULT, DEFAULT, DEFAULT, DEFAULT),
			Body:       createPrivateUser(name, password, email),
			StatusCode: http.StatusOK,
		},
		TestCase{
			Response:   createResult(place, re.ErrorUserNotFound()),
			Body:       createPrivateUser(RANDOM, RANDOM, RANDOM),
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Response:   createResult(place, re.ErrorUserNotFound()),
			Body:       createPrivateUser("", "", ""),
			StatusCode: http.StatusBadRequest,
		},
	}

	user := models.UserPrivateInfo{
		Name:     name,
		Password: password,
	}

	_, cookiestr, err := H.register(context.Background(), user)
	if err != nil {
		t.Error(" error:", err.Error())
		return
	}

	cookie := cook.CreateCookie(cookiestr, H.Session)

	url := "/session"
	launchTests(t, H, cases, url, login, compareStrings, nil, true)
	launchTests(t, H, cases, url, login, compareStrings, cookie, true)
}

// old tests

// +
func TLogout(t *testing.T, H *Handler) {
	const place = "Logout"
	name := utils.RandomString(16)
	password := utils.RandomString(16)
	email := utils.RandomString(16)
	cases := []TestCase{
		TestCase{
			Response:   createResult(place, nil),
			Body:       "",
			StatusCode: http.StatusOK,
		},
		TestCase{
			Response:   createResult(place, re.ErrorAuthorization()),
			Body:       "",
			StatusCode: http.StatusUnauthorized,
		},
	}

	cookie, err := getCookie(H, name, password, email)
	if err != nil {
		t.Error(" error:", err.Error())
		return
	}

	url := "/session"
	launchTests(t, H, cases[:1], url, logout, compareStrings, cookie, true)
	launchTests(t, H, cases[1:], url, logout, compareStrings, nil, true)
}

/*
func TestGetPlayerGames(t *testing.T) {
	cases := []TestCaseGames{
		TestCaseGames{
			Response: `[
				{
					"fieldWidth":25,
					"fieldHeight":25,
					"minsTotal":80,
					"minsFound":30,
					"finihsed":false,
					"exploded":false
				},
				{
					"fieldWidth":25,
					"fieldHeight":25,
					"minsTotal":70,
					"minsFound":70,
					"finihsed":true,
					"exploded":false
					}
				]`,
			name:       `panda`,
			page:       `1`,
			StatusCode: http.StatusOK,
		},
		TestCaseGames{
			Response:   `[]`,
			name:       `dfsdgf`,
			page:       `4545444`,
			StatusCode: http.StatusOK,
		},
		TestCaseGames{
			Response: `{
				"place":"GetPlayerGames",
				"success":false,
				"message":"Invalid page"}`,
			name:       `panda`,
			page:       `sdsd`,
			StatusCode: http.StatusBadRequest,
		},
	}

	H, _, err := GetHandler(PATH, "")
	if err != nil || H == nil {
		t.Error("TestCreateUser catched error:", err.Error())
		return
	}

	url := "/users"

	for caseNum, item := range cases {
		req1 := httptest.NewRequest("GET", url, nil)
		req1 = mux.SetURLVars(req1, map[string]string{"name": item.name, "page": item.page})
		w1 := httptest.NewRecorder()
		H.GetPlayerGames(w1, req1)

		if w1.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				1, w1.Code, item.StatusCode)
		}

		resp1 := w1.Result()
		body1, _ := ioutil.ReadAll(resp1.Body)
		defer resp1.Body.Close()

		checkStrings(t, caseNum, string(body1), item.Response)
	}
}
*/

/*
func TestGetUsersPageAmount(t *testing.T) {
	cases := []TestCase{
		TestCase{
			Response: `{"amount":0}
			`,
			StatusCode: http.StatusOK,
		},
		TestCase{
			Response: `{"amount":3}
			`,
			StatusCode: http.StatusOK,
		},
	}

	users := []TestCase{
		TestCase{
			Body:       createPrivateUser(RANDOM, RANDOM, RANDOM),
			Response:   ``,
			StatusCode: http.StatusUnauthorized,
		},
		TestCase{
			Body:       createPrivateUser(RANDOM, RANDOM, RANDOM),
			Response:   ``,
			StatusCode: http.StatusUnauthorized,
		},
		TestCase{
			Body:       createPrivateUser(RANDOM, RANDOM, RANDOM),
			Response:   ``,
			StatusCode: http.StatusUnauthorized,
		},
	}

	authConn, err := clients.ServiceConnectionsInit()
	if err != nil {
		log.Fatal("serviceConnectionsInit error:", err)
	}
	defer authConn.Close()

	H, _, err := GetHandler(PATH, "", authConn)

	if err != nil || H == nil {
		t.Error("catched error:", err.Error())
		return
	}
	url := "/users/pages_amount"

	launchTests(t, H, cases[:1], url, getUsersPageAmount, nil, true)
	launchTests(t, H, users, url, createUser, nil, false)
	launchTests(t, H, cases[1:], url, getUsersPageAmount, nil, true)
	launchTests(t, H, users, url, deleteUser, nil, false)
}

func TestGetUsers(t *testing.T) {
	authConn, err := clients.ServiceConnectionsInit()
	if err != nil {
		log.Fatal("serviceConnectionsInit error:", err)
	}
	defer authConn.Close()
	H, _, err := GetHandler(PATH, "", authConn)
	if err != nil || H == nil {
		t.Error("Catched error:", err.Error())
		return
	}
	url := "/users/pages"

	cases := getUserCreateCases()

	launchTests(t, H, cases, url, getUsers, nil, false)

	launchTests(t, H, cases, url, deleteUser, nil, false)
}

func TestPostImage(t *testing.T) {
	authConn, err := clients.ServiceConnectionsInit()
	if err != nil {
		log.Fatal("serviceConnectionsInit error:", err)
	}
	defer authConn.Close()

	H, _, err := GetHandler(PATH, "", authConn)
	if err != nil || H == nil {
		t.Error("Catched error:", err.Error())
		return
	}
	url := "/avatar"

	cases := getUserCreateCases()

	launchTests(t, H, cases, url, postImage, nil, false)

	launchTests(t, H, cases, url, deleteUser, nil, false)
}

func TestGetImage(t *testing.T) {
	authConn, err := clients.ServiceConnectionsInit()
	if err != nil {
		log.Fatal("serviceConnectionsInit error:", err)
	}
	defer authConn.Close()
	H, _, err := GetHandler(PATH, "", authConn)
	if err != nil || H == nil {
		t.Error("Catched error:", err.Error())
		return
	}
	url := "/avatar"

	cases := getUserCreateCases()

	const place = "UpdateProfile"

	user := *createRandomUser()
	_, cookiestr, err := H.register(context.Background(), user)
	if err != nil {
		t.Error(" error:", err.Error())
		return
	}
	cookie := cook.CreateCookie(cookiestr, H.Session)

	launchTests(t, H, cases, url, getImage, cookie, false)
	launchTests(t, H, cases, url, getImage, nil, false)

	if err = H.deleteAccount(context.Background(), &user, cookiestr); err != nil {
		t.Error(" error:", err.Error())
		return
	}
}
*/

/////////// object comparators /////

type Comparator func(a, b string) bool

/////////// api handlers //////

type apiFunc func(H *Handler, url string, r *strings.Reader, c *http.Cookie) *httptest.ResponseRecorder

func createUser(H *Handler, url string, r *strings.Reader, c *http.Cookie) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", url, r)
	w := httptest.NewRecorder()
	if c != nil {
		req.AddCookie(c)
	}
	H.CreateUser(w, req)
	return w
}

func deleteUser(H *Handler, url string, r *strings.Reader, c *http.Cookie) *httptest.ResponseRecorder {
	req := httptest.NewRequest("DELETE", url, r)
	w := httptest.NewRecorder()
	if c != nil {
		req.AddCookie(c)
	}
	H.DeleteUser(w, req)
	return w
}

func updateUser(H *Handler, url string, r *strings.Reader, c *http.Cookie) *httptest.ResponseRecorder {
	req := httptest.NewRequest("PUT", url, r)
	w := httptest.NewRecorder()
	if c != nil {
		req.AddCookie(c)
	}
	H.UpdateProfile(w, req)
	return w
}

func getMyProfile(H *Handler, url string, r *strings.Reader, c *http.Cookie) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", url, r)
	w := httptest.NewRecorder()
	if c != nil {
		req.AddCookie(c)
	}
	H.GetMyProfile(w, req)
	return w
}

func getProfile(H *Handler, url string, r *strings.Reader, c *http.Cookie) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", url, r)
	w := httptest.NewRecorder()
	if c != nil {
		req.AddCookie(c)
	}
	H.GetProfile(w, req)
	return w
}

func getUsersPageAmount(H *Handler, url string, r *strings.Reader, c *http.Cookie) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", url, r)
	w := httptest.NewRecorder()
	if c != nil {
		req.AddCookie(c)
	}
	H.GetUsersPageAmount(w, req)
	return w
}

func postImage(H *Handler, url string, r *strings.Reader, c *http.Cookie) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", url, r)
	w := httptest.NewRecorder()
	if c != nil {
		req.AddCookie(c)
	}
	H.PostImage(w, req)
	return w
}

func getImage(H *Handler, url string, r *strings.Reader, c *http.Cookie) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", url, r)
	w := httptest.NewRecorder()
	if c != nil {
		req.AddCookie(c)
	}
	H.GetImage(w, req)
	return w
}

func getUsers(H *Handler, url string, r *strings.Reader, c *http.Cookie) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", url, r)
	w := httptest.NewRecorder()
	if c != nil {
		req.AddCookie(c)
	}
	H.GetUsers(w, req)
	return w
}

func login(H *Handler, url string, r *strings.Reader, c *http.Cookie) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", url, r)
	w := httptest.NewRecorder()
	if c != nil {
		req.AddCookie(c)
	}
	H.Login(w, req)
	return w
}

func logout(H *Handler, url string, r *strings.Reader, c *http.Cookie) *httptest.ResponseRecorder {
	req := httptest.NewRequest("DELETE", url, r)
	w := httptest.NewRecorder()
	if c != nil {
		req.AddCookie(c)
	}
	H.Logout(w, req)
	return w
}

/////// create string models ///////

func createString(key, value string) string {
	if value == "" {
		return value
	}
	return `"` + key + `":"` + value + `"`
}

func createNotString(key, value string) string {
	return `"` + key + `":` + value
}

func createModel(fields ...string) string {
	var (
		returnString string
		count        int
	)
	returnString = "{"
	for _, field := range fields {
		if field != "" {
			if count != 0 {
				if field[0] != '{' {
					returnString += ","
				}
			}
			returnString += field
			count++
		}
	}
	return returnString + "}"
}

func createResult(place string, err error) string {
	errText := "no error"
	success := "true"
	if err != nil {
		errText = err.Error()
		success = "false"
	}
	return createModel(
		createString("place", place),
		createNotString("success", success),
		createString("message", errText))
}

func createPrivateUser(name, password, email string) string {
	if name == RANDOM {
		name = utils.RandomString(16)
	}
	if password == RANDOM {
		password = utils.RandomString(16)
	}
	if email == RANDOM {
		email = utils.RandomString(16)
	}
	return createModel(
		createString("name", name),
		createString("password", password),
		createString("email", email))
}

// Only server can send it, so Random not available
func createPublicUser(name, email, photo, bestScore, valid1, bestTime, valid2, difficult string) string {
	if bestScore == DEFAULT {
		bestScore = "0"
		valid1 = "true"
	}

	if bestTime == DEFAULT {
		bestTime = "24:00:00"
		valid2 = "true"
	}

	if difficult == DEFAULT {
		difficult = "0"
	}

	return createModel(
		createString("name", name),
		createString("email", email),
		createString("photo", photo),
		`"bestScore":`,
		createModel(
			createString("String", bestScore),
			createNotString("Valid", valid1),
		),
		`"bestTime":`,
		createModel(
			createString("String", bestTime),
			createNotString("Valid", valid2)),
		createNotString("difficult", difficult))
}

func createRandomUser() *models.UserPrivateInfo {
	return &models.UserPrivateInfo{
		Name:     utils.RandomString(16),
		Password: utils.RandomString(16),
	}
}

func getCookie(H *Handler, name, password, email string) (cookie *http.Cookie, err error) {
	user := models.UserPrivateInfo{
		Name:     name,
		Password: password,
	}

	var cookiestr string
	if _, cookiestr, err = H.register(context.Background(), user); err != nil {
		return
	}

	cookie = cook.CreateCookie(cookiestr, H.Session)
	return
}

//00:03 1069
