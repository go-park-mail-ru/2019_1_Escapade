package handlers

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

// Result - every handler return it
type Result struct {
	code  int
	place string
	send  JSONtype
	err   error
}

// NewResult construct instance of result
func NewResult(code int, place string, send JSONtype, err error) Result {
	return Result{
		code:  code,
		place: place,
		send:  send,
		err:   err,
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
		sendErrorJSON(rw, result.err, result.place)
	} else {
		sendSuccessJSON(rw, result.send, result.place)
	}
	Debug(result.err, result.code, result.place)
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
func Warning(err error, text string, place string) {
	if err != nil {
		utils.Debug(false, "Warning in "+place+".", text, "More:", err.Error())
	} else {
		utils.Debug(false, "Warning in "+place+".", text)
	}
}

// SendErrorJSON send error json
func sendErrorJSON(rw http.ResponseWriter, catched error, place string) {
	result := models.Result{
		Place:   place,
		Success: false,
		Message: catched.Error(),
	}

	if b, err := result.MarshalJSON(); err == nil {
		rw.Write(b)
	}
}

// SendSuccessJSON send object json
func sendSuccessJSON(rw http.ResponseWriter, result JSONtype, place string) {
	if result == nil {
		result = &models.Result{
			Place:   place,
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
			code:  http.StatusMethodNotAllowed,
			err:   re.ErrorMethodNotAllowed(),
			place: place,
		}
	}
	SendResult(rw, *result)
}

func SendError(code int, place string, err error) func(rw http.ResponseWriter, r *http.Request) Result {
	return func(rw http.ResponseWriter, r *http.Request) Result {
		return NewResult(code, place, nil, err)
	}
}
