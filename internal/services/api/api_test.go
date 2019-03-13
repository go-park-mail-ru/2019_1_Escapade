package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type TestCase struct {
	Response   string
	Body       string
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
	H.DeleteAccount(httptest.NewRecorder(), preq)
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
