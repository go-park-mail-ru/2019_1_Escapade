package api

import (
	cook "escapade/internal/cookie"
	re "escapade/internal/return_errors"
	"escapade/internal/utils"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
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

var replacer = strings.NewReplacer("\n", "", "\t", "")

func checkStrings(t *testing.T, caseNum int, got string, expected string) {
	got = strings.TrimSpace(got)
	expected = strings.TrimSpace(expected)
	got = replacer.Replace(got)
	expected = replacer.Replace(expected)
	if got != expected {
		t.Errorf("[%d] wrong Response: got \n%+v\nexpected\n%+v",
			caseNum, got, expected)
	}
}

func checkCookies(t *testing.T, caseNum int, got http.Cookie, expected http.Cookie) {

	if got.Name != expected.Name {
		t.Errorf("[%d] wrong cookie name: got \n%+v\nexpected\n%+v",
			caseNum, got.Name, expected.Name)
	}

	if got.Value != expected.Value {
		t.Errorf("[%d] wrong cookie value: got \n%+v\nexpected\n%+v",
			caseNum, got.Value, expected.Value)
	}

	if got.Path != expected.Path {
		t.Errorf("[%d] wrong cookie path: got \n%+v\nexpected\n%+v",
			caseNum, got.Path, expected.Path)
	}

	if got.HttpOnly != expected.HttpOnly {
		t.Errorf("[%d] wrong cookie HttpOnly flag: got \n%+v\nexpected\n%+v",
			caseNum, got.HttpOnly, expected.HttpOnly)
	}
}

func launchTests(t *testing.T, H *Handler, cases []TestCase,
	url string, af apiFunc, c *http.Cookie, record bool) {
	for caseNum, item := range cases {
		w := af(H, url, strings.NewReader(item.Body), c)
		fmt.Println("look at", item.Body)
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

		checkStrings(t, caseNum, string(body), item.Response)
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
			Response:   createResult(place, re.ErrorEmailIstaken()),
			Body:       createPrivateUser(RANDOM, RANDOM, email1),
			StatusCode: http.StatusBadRequest,
		}, // name is taken
		TestCase{
			Response:   createResult(place, re.ErrorNameIstaken()),
			Body:       createPrivateUser(name1, RANDOM, RANDOM),
			StatusCode: http.StatusBadRequest,
		}, // no password
		TestCase{
			Response:   createResult(place, re.ErrorInvalidPassword()),
			Body:       createPrivateUser(RANDOM, "", RANDOM),
			StatusCode: http.StatusBadRequest,
		}, // no email
		TestCase{
			Response:   createResult(place, re.ErrorInvalidEmail()),
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
			Response:   createResult(place, re.ErrorInvalidEmail()),
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

func TestCreateUser(t *testing.T) {
	H, _, err := GetHandler(PATH)
	if err != nil || H == nil {
		t.Error("TestCreateUser catched error:", err.Error())
		return
	}
	url := "/user"

	cases := getUserCreateCases()

	launchTests(t, H, cases, url, createUser, nil, true)

	launchTests(t, H, cases, url, deleteUser, nil, false)
}

func TestDeleteUser(t *testing.T) {
	H, _, err := GetHandler(PATH)
	if err != nil || H == nil {
		t.Error("TestCreateUser catched error:", err.Error())
		return
	}
	url := "/user"

	cases := getUserDeleteCases()

	launchTests(t, H, cases, url, createUser, nil, false)
	launchTests(t, H, cases, url, deleteUser, nil, true)
}
func TestUpdateProfile(t *testing.T) {
	const place = "UpdateProfile"
	firstName := utils.RandomString(16)
	takenName := utils.RandomString(16)
	takenEmail := utils.RandomString(16)

	users := []TestCase{
		TestCase{
			Body:       createPrivateUser(firstName, RANDOM, RANDOM),
			Response:   ``,
			StatusCode: http.StatusUnauthorized,
		},
		TestCase{
			Body:       createPrivateUser(takenName, RANDOM, takenEmail),
			Response:   ``,
			StatusCode: http.StatusUnauthorized,
		},
	}
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
			Response:   createResult(place, re.ErrorEmailIstaken()),
			Body:       createPrivateUser(takenName, "", takenEmail),
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Response:   createResult(place, re.ErrorNameIstaken()),
			Body:       createPrivateUser(takenName, "", ""),
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Response:   createResult(place, re.ErrorEmailIstaken()),
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
		cookiestr string
		cookie    *http.Cookie
	)

	H, _, err := GetHandler(PATH)

	if err != nil || H == nil {
		t.Error("TestUpdateUser catched error:", err.Error())
		return
	}
	url := "/user"

	launchTests(t, H, users, url, createUser, nil, false)

	if cookiestr, err = H.DB.GetSessionByName(firstName); err != nil {
		t.Error("TestUpdateUser cant get cookie:", err.Error())
		return
	}

	cookie = cook.CreateCookie(cookiestr, H.Cookie)

	launchTests(t, H, updates[:1], url, updateUser, nil, true)
	launchTests(t, H, updates[1:], url, updateUser, cookie, true)
	launchTests(t, H, users, url, deleteUser, nil, false)
}

func TestGetMyProfile(t *testing.T) {
	const place = "GetMyProfile"
	name := utils.RandomString(16)
	email := utils.RandomString(16)

	users := []TestCase{
		TestCase{
			Body:       createPrivateUser(name, RANDOM, email),
			Response:   ``,
			StatusCode: http.StatusUnauthorized,
		},
	}
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

	H, _, err := GetHandler(PATH)

	if err != nil || H == nil {
		t.Error("TestUpdateUser catched error:", err.Error())
		return
	}
	url := "/user"

	launchTests(t, H, users, url, createUser, nil, false)

	var cookiestr string
	if cookiestr, err = H.DB.GetSessionByName(name); err != nil {
		t.Error("TestUpdateUser cant get cookie:", err.Error())
		return
	}

	cookie := cook.CreateCookie(cookiestr, H.Cookie)

	launchTests(t, H, gets[:1], url, getMyProfile, nil, true)
	launchTests(t, H, gets[1:], url, getMyProfile, cookie, true)
	launchTests(t, H, users, url, deleteUser, nil, false)
}

func TestGetProfile(t *testing.T) {
	const place = "GetProfile"
	name := utils.RandomString(16)
	email := utils.RandomString(16)

	users := []TestCase{
		TestCase{
			Body:       createPrivateUser(name, RANDOM, email),
			Response:   ``,
			StatusCode: http.StatusUnauthorized,
		},
	}
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

	H, _, err := GetHandler(PATH)

	if err != nil || H == nil {
		t.Error("TestUpdateUser catched error:", err.Error())
		return
	}

	url := "/user"
	launchTests(t, H, users, url, createUser, nil, false)

	launchTests(t, H, gets[:1], url, getProfile, nil, true)
	url = "/user?id=dfdf"
	launchTests(t, H, gets[:1], url, getProfile, nil, true)
	url = "/user?id=10"
	launchTests(t, H, gets[1:2], url, getProfile, nil, true)
	url = "/user?id=1"
	launchTests(t, H, gets[2:], url, getProfile, nil, true)
	url = "/user"
	launchTests(t, H, users, url, deleteUser, nil, false)
}

// old tests

// +
func TestLogin(t *testing.T) {
	cases := []TestCase{
		TestCase{
			Response: `
				{
					"place":"Login",
					"success":false,
					"message":"Cant found parameters"
				}
				`,
			Body: `{
				"name": "username",
				"password": "1454543",
				"email": "test@mail.ru"
			}`,
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Response: createPublicUser("username", "test@mail.ru", "",
				DEFAULT, DEFAULT, DEFAULT, DEFAULT, DEFAULT),
			Body: `{
				"name": "username",
				"password": "1454543",
				"email": "test@mail.ru"
			}`,
			StatusCode: http.StatusOK,
		},
		TestCase{
			Response: `{
				"place":"Login",
				"success":false,
				"message":"User not found"}
			`,
			Body: `{
				"name": "username",
				"password": "145sdsw4543",
				"email": "test121@mail.ru"
			}`,
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Response: `{
				"place":"Login",
				"success":false,
				"message":"User not found"}
			`,
			Body: `{
				"password": "",
				"email": ""
			}`,
			StatusCode: http.StatusBadRequest,
		},
	}

	urlSignUp := "/user"

	H, _, err := GetHandler(PATH)
	if err != nil || H == nil {
		t.Error("TestCreateUser catched error:", err.Error())
		return
	}
	preq := httptest.NewRequest("DELETE", urlSignUp, strings.NewReader(cases[0].Body))
	H.DeleteUser(httptest.NewRecorder(), preq)

	req := httptest.NewRequest("POST", urlSignUp, strings.NewReader(cases[0].Body))
	w := httptest.NewRecorder()
	H.CreateUser(w, req)
	//body, _ := ioutil.ReadAll(w.Result().Body)
	//t.Errorf(string(body))

	url := "/session"
	for caseNum, item := range cases {
		req := httptest.NewRequest("POST", url, strings.NewReader(item.Body))
		if caseNum == 0 {
			req.Body = nil
		}
		w := httptest.NewRecorder()

		H.Login(w, req)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.StatusCode)
		}

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		defer resp.Body.Close()

		checkStrings(t, caseNum, string(body), item.Response)
	}
}

// +
func TestLogout(t *testing.T) {
	cases := []TestCase{
		TestCase{
			Response: `{
				"place":"Logout",
				"success":true,
				"message":"no error"}
			`,
			Body: `{
				"name": "username1",
				"password": "1454543",
				"email": "test1@mail.ru"
			}`,
			StatusCode: http.StatusOK,
		},
		TestCase{
			Response: `{
				"place":"Logout",
				"success":false,
				"message":"Required authorization"}
			`,
			Body: `{
				"name": "username",
				"password": "145sdsw4543",
				"email": "test121@mail.ru"
			}`,
			StatusCode: http.StatusUnauthorized,
		},
	}

	urlSignUp := "/user"

	H, _, err := GetHandler(PATH)
	if err != nil || H == nil {
		t.Error("TestCreateUser catched error:", err.Error())
		return
	}
	preq := httptest.NewRequest("DELETE", urlSignUp, strings.NewReader(cases[0].Body))
	H.DeleteUser(httptest.NewRecorder(), preq)

	req := httptest.NewRequest("POST", urlSignUp, strings.NewReader(cases[0].Body))
	w := httptest.NewRecorder()
	H.CreateUser(w, req)

	urlLogin := "/session"
	reqLogin := httptest.NewRequest("POST", urlLogin, strings.NewReader(cases[0].Body))
	wLogin := httptest.NewRecorder()
	str, err := H.DB.GetSessionByName("username1")
	if err != nil {
		return
	}

	reqLogin.AddCookie(cook.CreateCookie(str, H.Cookie))
	H.Login(wLogin, reqLogin)

	//body, _ := ioutil.ReadAll(wLogin.Result().Body)
	//t.Errorf(string(body))

	var cookie *http.Cookie
	if cookie, err = reqLogin.Cookie(H.Cookie.NameCookie); err != nil {
		t.Error("TestUpdateUser cant get cookie:", err.Error())
		return
	}

	url := "/session"

	/* 1 test */
	req1 := httptest.NewRequest("DELETE", url, strings.NewReader(cases[0].Body))
	w1 := httptest.NewRecorder()
	req1.AddCookie(cookie)
	H.Logout(w1, req1)

	if w1.Code != cases[0].StatusCode {
		t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
			1, w1.Code, cases[0].StatusCode)
	}

	resp1 := w1.Result()
	body1, _ := ioutil.ReadAll(resp1.Body)
	defer resp1.Body.Close()

	checkStrings(t, 1, string(body1), cases[0].Response)
	/* 2 test */
	req2 := httptest.NewRequest("DELETE", url, strings.NewReader(cases[1].Body))
	w2 := httptest.NewRecorder()
	H.Logout(w2, req2)

	if w2.Code != cases[1].StatusCode {
		t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
			2, w2.Code, cases[1].StatusCode)
	}

	resp2 := w2.Result()
	body2, _ := ioutil.ReadAll(resp2.Body)
	defer resp2.Body.Close()

	checkStrings(t, 2, string(body2), cases[1].Response)

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

	H, _, err := GetHandler(PATH)
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

	H, _, err := GetHandler(PATH)

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
	H, _, err := GetHandler(PATH)
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
	H, _, err := GetHandler(PATH)
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
	H, _, err := GetHandler(PATH)
	if err != nil || H == nil {
		t.Error("Catched error:", err.Error())
		return
	}
	url := "/avatar"

	cases := getUserCreateCases()

	const place = "UpdateProfile"
	firstName := utils.RandomString(16)

	users := []TestCase{
		TestCase{
			Body:       createPrivateUser(firstName, RANDOM, RANDOM),
			Response:   ``,
			StatusCode: http.StatusUnauthorized,
		},
	}

	launchTests(t, H, users, url, createUser, nil, false)

	var cookiestr string
	if cookiestr, err = H.DB.GetSessionByName(firstName); err != nil {
		t.Error("TestUpdateUser cant get cookie:", err.Error())
		return
	}

	cookie := cook.CreateCookie(cookiestr, H.Cookie)

	launchTests(t, H, cases, url, getImage, cookie, false)
	launchTests(t, H, cases, url, getImage, nil, false)
	launchTests(t, H, users, url, deleteUser, nil, false)
}

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
