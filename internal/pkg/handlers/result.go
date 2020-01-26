package handlers

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

// Result - every handler return it
type Result struct {
	code int
	send JSONtype
	err  error
}

type ResultFunc func(http.ResponseWriter, *http.Request) Result

func HandleFunc(rf ResultFunc) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		result := rf(rw, r)
		Send(rw, r, &result)
	}

}

// NewResult construct instance of result
func NewResult(code int, send JSONtype, err error) Result {
	return Result{
		code: code,
		send: send,
		err:  err,
	}
}

// SendResult send result to client
func SendResult(rw http.ResponseWriter, result Result) {
	if result.code == 0 {
		return
	}
	utils.Debug(false, "WriteHeader!", result.code)
	rw.WriteHeader(result.code)
	if result.err != nil {
		sendErrorJSON(rw, result.err)
	} else {
		sendSuccessJSON(rw, result.send)
	}
	Debug(result.err, result.code, "")
}

// Debug use utils.debug to log information
func Debug(catched error, number int, place string) {
	if catched != nil {
		utils.Debug(false, "api/"+place+" failed(code:", number, "). Error message:"+catched.Error())
	} else {
		utils.Debug(false, "api/"+place+" success(code:", number, ")")
	}
}

// Warning use utils.debug to log warnings
func Warning(err error, text string) {
	if err != nil {
		utils.Debug(false, "Warning.", text, "More:", err.Error())
	} else {
		utils.Debug(false, "Warning.", text)
	}
}

// SendErrorJSON send error json
func sendErrorJSON(rw http.ResponseWriter, catched error) {
	result := models.Result{
		Success: false,
		Message: catched.Error(),
	}

	if b, err := result.MarshalJSON(); err == nil {
		rw.Write(b)
	}
}

// SendSuccessJSON send object json
func sendSuccessJSON(rw http.ResponseWriter, result JSONtype) {
	if result == nil {
		result = &models.Result{
			Success: true,
			Message: "no error",
		}
	}
	if b, err := result.MarshalJSON(); err == nil {
		utils.Debug(false, string(b))
		rw.Write(b)
	}
}

// HandleRequest handle request and return Result
type HandleRequest func(rw http.ResponseWriter, r *http.Request) Result

// MethodHandlers - map, where key - method, value - HandleRequest
type MethodHandlers map[string]HandleRequest

// Route direct the request depending on its method
// mHr - map, where key - method, value - HandleRequest
func Route(rw http.ResponseWriter, r *http.Request, mHr MethodHandlers) {
	var result *Result
	for k, v := range mHr {
		if k == r.Method {
			if v == nil {
				return
			}
			r := v(rw, r)
			result = &r
			break
		}
	}

	if result == nil {
		place := r.URL.Path
		utils.Debug(false, place+" method not allowed:", r.Method)
		result = &Result{
			code: http.StatusMethodNotAllowed,
			err:  re.ErrorMethodNotAllowed(),
		}
	}
	SendResult(rw, *result)
}

func Send(rw http.ResponseWriter, request *http.Request, result *Result) {
	if result == nil {
		place := request.URL.Path
		utils.Debug(false, place+" method not allowed:", request.Method)
		result = &Result{
			code: http.StatusMethodNotAllowed,
			err:  re.ErrorMethodNotAllowed(),
		}
	}
	SendResult(rw, *result)
}

func SendError(code int, err error) func(rw http.ResponseWriter, r *http.Request) Result {
	return func(rw http.ResponseWriter, r *http.Request) Result {
		return NewResult(code, nil, err)
	}
}

func OPTIONS() func(rw http.ResponseWriter, r *http.Request) Result {
	return func(rw http.ResponseWriter, r *http.Request) Result {
		return NewResult(0, nil, nil)
	}
}
