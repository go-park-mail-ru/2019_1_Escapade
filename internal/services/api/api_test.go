package api

import (
	misc "escapade/internal/misc"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
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

const confPath = "../../../conf.json"

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

func TestCreateUser(t *testing.T) {
	cases := []TestCase{
		TestCase{
			Response: `{
				"place":"CreateUser",
				"success":true,
				"message":"no error"}
			`,
			Body: `{
				"name": "username111",
				"password": "1454543",
				"email": "123123"
			}`,
			StatusCode: http.StatusCreated,
		},
		TestCase{
			Response: `{
				"place":"CreateUser",
				"success":false,
				"message":"Invalid password"}`,
			Body: `{
				"name": "TestCase2",
				"email": "123123"
			}`,
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Response: `{
				"place":"CreateUser",
				"success":false,
				"message":"Invalid username"}`,
			Body: `{
				"password": "username111",
				"email": "123123"
			}`,
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Response: `{
				"place":"CreateUser",
				"success":false,
				"message":"Invalid email"}`,
			Body: `{
				"name": "username111",
				"password": "123123"
			}`,
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Response: `{
				"place":"CreateUser",
				"success":false,
				"message":"Invalid username"}`,
			Body: `{
			}`,
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Response: `{
				"place":"CreateUser",
				"success":false,
				"message":"Username is taken"}
			`,
			Body: `{
				"name": "username111",
				"password": "1454543",
				"email": "emailtest"
			}`,
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Response: `{
				"place":"CreateUser",
				"success":false,
				"message":"Email is taken"}
			`,
			Body: `{
				"name": "usertest",
				"password": "1454543",
				"email": "123123"
			}`,
			StatusCode: http.StatusBadRequest,
		},
	}

	H, _, err := GetHandler(confPath)
	if err != nil || H == nil {
		fmt.Println("TestCreateUser catched error:", err.Error())
		return
	}

	url := "/user"
	preq := httptest.NewRequest("DELETE", url, strings.NewReader(cases[0].Body))
	H.DeleteUser(httptest.NewRecorder(), preq)
	for caseNum, item := range cases {
		req := httptest.NewRequest("POST", url, strings.NewReader(item.Body))
		w := httptest.NewRecorder()

		H.CreateUser(w, req)

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

func TestDeleteUser(t *testing.T) {
	cases := []TestCase{
		TestCase{
			Response: `{
				"place":"DeleteUser",
				"success":true,
				"message":"no error"}
			`,
			Body: `{
				"name": "username111",
				"password": "1454543",
				"email": "123123"
			}`,
			StatusCode: http.StatusOK,
		},
		TestCase{
			Response: `{
				"place":"DeleteUser",
				"success":false,
				"message":"Invalid password"}`,
			Body: `{
				"name": "TestCase2",
				"email": "123123"
			}`,
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Response: `{
				"place":"DeleteUser",
				"success":false,
				"message":"Invalid username"}`,
			Body: `{
				"password": "username111",
				"email": "123123"
			}`,
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Response: `{
				"place":"DeleteUser",
				"success":false,
				"message":"Invalid email"}`,
			Body: `{
				"name": "username111",
				"password": "123123"
			}`,
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Response: `{
				"place":"DeleteUser",
				"success":false,
				"message":"Invalid username"}`,
			Body: `{
			}`,
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Response: `{
				"place":"DeleteUser",
				"success":false,
				"message":"User not found"}
			`,
			Body: `{
				"name": "username111",
				"password": "1454543",
				"email": "emailtest"
			}`,
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Response: `{
				"place":"DeleteUser",
				"success":false,
				"message":"User not found"}
			`,
			Body: `{
				"name": "usertest",
				"password": "1454543",
				"email": "123123"
			}`,
			StatusCode: http.StatusBadRequest,
		},
	}

	H, _, err := GetHandler(confPath)
	if err != nil || H == nil {
		fmt.Println("TestDeleteUser catched error:", err.Error())
		return
	}

	url := "/user"
	preq := httptest.NewRequest("POST", url, strings.NewReader(cases[0].Body))
	H.CreateUser(httptest.NewRecorder(), preq)

	for caseNum, item := range cases {
		req := httptest.NewRequest("DELETE", url, strings.NewReader(item.Body))
		w := httptest.NewRecorder()

		H.DeleteUser(w, req)

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

func TestUpdateProfile(t *testing.T) {
	users := []string{
		`{
				"name": "TestUpdateUser",
				"password": "TestUpdateUser",
				"email": "TestUpdateUser"
			}`,
		`{
				"name": "taken",
				"password": "taken",
				"email": "taken"
			}`,
	}
	cases := []TestCase{
		TestCase{
			Response: `{
				"place":"UpdateProfile",
				"success":false,
				"message":"Required authorization"
			}
			`,
			Body:       ``,
			StatusCode: http.StatusUnauthorized,
		},
		TestCase{
			Response: `{
				"place":"UpdateProfile",
				"success":false,
				"message":"Cant found parameters"
			}
			`,
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Response: `{
				"place":"UpdateProfile",
				"success":true,
				"message":"no error"}
			`,
			Body: `{
				"name": "TestUpdateUser1"
			}`,
			StatusCode: http.StatusOK,
		},
		TestCase{
			Response: `{
				"place":"UpdateProfile",
				"success":true,
				"message":"no error"}
			`,
			Body: `{
				"email": "TestUpdateUser2"
			}`,
			StatusCode: http.StatusOK,
		},
		TestCase{
			Response: `{
				"place":"UpdateProfile",
				"success":true,
				"message":"no error"}
			`,
			Body: `{
				"password": "TestUpdateUser3"
			}`,
			StatusCode: http.StatusOK,
		},
		TestCase{
			Response: `{
				"place":"UpdateProfile",
				"success":true,
				"message":"no error"}
			`,
			Body: `{
				"name": "TestUpdateUser",
				"password": "TestUpdateUser",
				"email": "TestUpdateUser"
			}`,
			StatusCode: http.StatusOK,
		},
		TestCase{
			Response: `{
				"place":"UpdateProfile",
				"success":false,
				"message":"Invalid email"}
			`,
			Body: `{
				"email": "taken"
			}`,
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Response: `{
				"place":"UpdateProfile",
				"success":false,
				"message":"Invalid username"}
			`,
			Body: `{
				"name": "taken"
			}`,
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Response: `{
				"place":"UpdateProfile",
				"success":true,
				"message":"no error"}
			`,
			Body:       ``,
			StatusCode: http.StatusOK,
		},
	}

	var (
		H         *Handler
		err       error
		url       string
		createReq *http.Request
		deleteReq *http.Request
		cookiestr string
		cookie    *http.Cookie
	)

	if H, _, err = GetHandler(confPath); err != nil || H == nil {
		fmt.Println("TestUpdateUser catched error:", err.Error())
		return
	}

	url = "/user"

	// create the user, which we will update
	createReq = httptest.NewRequest("POST", url, strings.NewReader(users[0]))
	w := httptest.NewRecorder()
	H.CreateUser(w, createReq)

	if cookiestr, err = H.DB.GetSessionByName("TestUpdateUser"); err != nil {
		fmt.Println("TestUpdateUser cant get cookie:", err.Error())
		return
	}
	cookie = misc.CreateCookie(cookiestr)

	// create user, which name/email we try to take(expected catch error)
	createReq = httptest.NewRequest("POST", url, strings.NewReader(users[1]))
	H.CreateUser(httptest.NewRecorder(), createReq)

	for caseNum, item := range cases {
		url := "/user"
		var req *http.Request

		req = httptest.NewRequest("PUT", url, strings.NewReader(item.Body))

		if caseNum == 1 {
			req.Body = nil
		}

		w := httptest.NewRecorder()
		if caseNum != 0 {
			req.AddCookie(cookie)
		}
		H.UpdateProfile(w, req)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.StatusCode)
		}

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		defer resp.Body.Close()

		checkStrings(t, caseNum, string(body), item.Response)
	}

	// delete created users
	deleteReq = httptest.NewRequest("DELETE", url, strings.NewReader(users[0]))
	H.DeleteUser(httptest.NewRecorder(), deleteReq)
	deleteReq = httptest.NewRequest("DELETE", url, strings.NewReader(users[1]))
	H.DeleteUser(httptest.NewRecorder(), deleteReq)
}

func TestGetProfile(t *testing.T) {
	users := []string{
		`{
				"name": "TestGetProfile",
				"password": "TestGetProfile",
				"email": "TestGetProfile"
			}`,
	}
	cases := []TestCase{
		TestCase{
			Response: `{
				"place":"GetMyProfile",
				"success":false,
				"message":"Required authorization"
			}
			`,
			Body:       ``,
			StatusCode: http.StatusUnauthorized,
		},
		TestCase{
			Response: `{
				"name":"TestGetProfile",
				"email":"TestGetProfile",
				"bestScore":{
					"String":"0",
					"Valid":true
					},
					"bestTime":{
						"String":"0",
						"Valid":true
					}
				}
			`,
			StatusCode: http.StatusOK,
		},
	}

	var (
		H         *Handler
		err       error
		url       string
		createReq *http.Request
		deleteReq *http.Request
		cookiestr string
		cookie    *http.Cookie
	)

	if H, _, err = GetHandler(confPath); err != nil || H == nil {
		fmt.Println("TestGetUser catched error:", err.Error())
		return
	}

	url = "/user"

	// create the user, which we will update
	createReq = httptest.NewRequest("POST", url, strings.NewReader(users[0]))
	H.CreateUser(httptest.NewRecorder(), createReq)

	okreq := httptest.NewRequest("GET", url, nil)
	H.Ok(httptest.NewRecorder(), okreq)

	if cookiestr, err = H.DB.GetSessionByName("TestGetProfile"); err != nil {
		fmt.Println("TestGetUser cant get cookie:", err.Error())
		return
	}
	cookie = misc.CreateCookie(cookiestr)

	for caseNum, item := range cases {
		url := "/user"
		var req *http.Request

		req = httptest.NewRequest("GET", url, strings.NewReader(item.Body))

		w := httptest.NewRecorder()
		if caseNum != 0 {
			req.AddCookie(cookie)
		}
		H.GetMyProfile(w, req)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.StatusCode)
		}

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		defer resp.Body.Close()

		checkStrings(t, caseNum, string(body), item.Response)
	}

	// delete created users
	deleteReq = httptest.NewRequest("DELETE", url, strings.NewReader(users[0]))
	H.DeleteUser(httptest.NewRecorder(), deleteReq)
}

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
			Response: `{
				"name":"username",
				"email":"test@mail.ru",
				"bestScore":{"String":"0","Valid":true},
				"bestTime":{"String":"0","Valid":true}}
				`,
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

	H, _, err := GetHandler(confPath)
	if err != nil || H == nil {
		fmt.Println("TestCreateUser catched error:", err.Error())
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

	H, _, err := GetHandler(confPath)
	if err != nil || H == nil {
		fmt.Println("TestCreateUser catched error:", err.Error())
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
		t.Errorf("aaaaaaaaaaaaaa")
		return
	}

	reqLogin.AddCookie(misc.CreateCookie(str))
	H.Login(wLogin, reqLogin)

	//body, _ := ioutil.ReadAll(wLogin.Result().Body)
	//t.Errorf(string(body))

	var cookie *http.Cookie
	if cookie, err = reqLogin.Cookie(misc.NameCookie); err != nil {
		fmt.Println("TestUpdateUser cant get cookie:", err.Error())
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

	H, _, err := GetHandler(confPath)
	if err != nil || H == nil {
		fmt.Println("TestCreateUser catched error:", err.Error())
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

func TestGetUsersPageAmount(t *testing.T) {
	cases := []TestCaseGames{
		TestCaseGames{
			Response: `{"amount":1}
			`,
			StatusCode: http.StatusOK,
		},
	}

	H, _, err := GetHandler(confPath)
	if err != nil || H == nil {
		fmt.Println("TestCreateUser catched error:", err.Error())
		return
	}

	url := "/users/pages_amount"

	/* 1 test */
	req1 := httptest.NewRequest("GET", url, nil)
	w1 := httptest.NewRecorder()
	H.GetUsersPageAmount(w1, req1)

	if w1.Code != cases[0].StatusCode {
		t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
			1, w1.Code, cases[0].StatusCode)
	}

	resp1 := w1.Result()
	body1, _ := ioutil.ReadAll(resp1.Body)
	defer resp1.Body.Close()

	checkStrings(t, 1, string(body1), cases[0].Response)
}
