package handlers

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

func HandleFunc(
	rf models.ResultFunc,
	trace infrastructure.ErrorTrace,
	logger infrastructure.Logger,
) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		result := rf(rw, r)
		Send(rw, r, &result, trace, logger)
	}
}

// NewResult construct instance of result
func NewResult(
	code int,
	send models.JSONtype,
	err error,
) models.RequestResult {
	return models.RequestResult{
		Code: code,
		Send: send,
		Err:  err,
	}
}

// SendResult send result to client
func SendResult(
	rw http.ResponseWriter,
	result models.RequestResult,
	logger infrastructure.Logger,
) {
	if result.Code == 0 {
		return
	}
	logger.Println("WriteHeader!", result.Code)
	rw.WriteHeader(result.Code)
	if result.Err != nil {
		sendErrorJSON(rw, result.Err)
	} else {
		sendSuccessJSON(rw, result.Send)
	}
	Debug(logger, result.Err, result.Code, "")
}

// Debug use utils.debug to log information
func Debug(
	logger infrastructure.Logger,
	catched error,
	number int,
	place string,
) {
	if catched != nil {
		logger.Println("api/"+place+" failed(code:",
			number, "). Error message:"+catched.Error())
	} else {
		logger.Println("api/"+place+" success(code:",
			number, ")")
	}
}

// Warning use utils.debug to log warnings
func Warning(
	logger infrastructure.Logger,
	err error,
	text string,
) {
	if err != nil {
		logger.Println("Warning.", text, "More:", err.Error())
	} else {
		logger.Println("Warning.", text)
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
func sendSuccessJSON(
	rw http.ResponseWriter,
	result models.JSONtype,
) {
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
type HandleRequest func(
	rw http.ResponseWriter,
	r *http.Request,
) models.RequestResult

// MethodHandlers - map, where key - method, value - HandleRequest
type MethodHandlers map[string]HandleRequest

// Route direct the request depending on its method
// mHr - map, where key - method, value - HandleRequest
func Route(
	rw http.ResponseWriter,
	r *http.Request,
	mHr MethodHandlers,
	trace infrastructure.ErrorTrace,
	logger infrastructure.Logger,
) {
	var result *models.RequestResult
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
		logger.Println(
			place+" method not allowed:",
			r.Method,
		)
		result = &models.RequestResult{
			Code: http.StatusMethodNotAllowed,
			Err:  trace.New(ErrMethodNotAllowed),
		}
	}
	SendResult(rw, *result, logger)
}

func Send(
	rw http.ResponseWriter,
	request *http.Request,
	result *models.RequestResult,
	trace infrastructure.ErrorTrace,
	logger infrastructure.Logger,
) {
	if result == nil {
		place := request.URL.Path
		logger.Println(
			place+" method not allowed:",
			request.Method,
		)
		result = &models.RequestResult{
			Code: http.StatusMethodNotAllowed,
			Err:  trace.New(ErrMethodNotAllowed),
		}
	}
	SendResult(rw, *result, logger)
}

func SendError(
	code int,
	err error,
) func(
	rw http.ResponseWriter,
	r *http.Request,
) models.RequestResult {
	return func(
		rw http.ResponseWriter,
		r *http.Request,
	) models.RequestResult {
		return NewResult(code, nil, err)
	}
}

func OPTIONS() func(
	rw http.ResponseWriter,
	r *http.Request,
) models.RequestResult {
	return func(
		rw http.ResponseWriter,
		r *http.Request,
	) models.RequestResult {
		return NewResult(0, nil, nil)
	}
}
