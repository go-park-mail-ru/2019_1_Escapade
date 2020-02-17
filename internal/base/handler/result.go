package handler

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

func (h *Handler) HandleFunc(
	rf models.ResultFunc,
) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		result := rf(rw, r)
		h.Send(rw, r, &result)
	}
}

func (h *Handler) Fail(
	code int,
	err error,
) models.RequestResult {
	return models.RequestResult{
		Code: code,
		Err:  err,
	}
}

func (h *Handler) Success(
	code int,
	send models.JSONtype,
) models.RequestResult {
	return models.RequestResult{
		Code: code,
		Send: send,
	}
}

// SendResult send result to client
func (h *Handler) SendResult(
	rw http.ResponseWriter,
	result models.RequestResult,
) {
	if result.Code == NoResult {
		return
	}

	h.logger.Println("WriteHeader!", result.Code)
	rw.WriteHeader(result.Code)
	if result.Err != nil {
		http.Error(rw, result.Err.Error(), result.Code)
		h.sendErrorJSON(rw, result.Err)
	} else {
		rw.WriteHeader(result.Code)
		h.sendSuccessJSON(rw, result.Send)
	}
	h.Debug(result.Err, result.Code, "")
}

// Debug use utils.debug to log information
func (h *Handler) Debug(
	catched error,
	number int,
	place string,
) {
	if catched != nil {
		h.logger.Println("api/"+place+" failed(code:",
			number, "). Error message:"+catched.Error())
	} else {
		h.logger.Println("api/"+place+" success(code:",
			number, ")")
	}
}

// Warning use utils.debug to log warnings
func (h *Handler) Warning(
	err error,
	text string,
) {
	if err != nil {
		h.logger.Println("Warning.", text, "More:", err.Error())
	} else {
		h.logger.Println("Warning.", text)
	}
}

// SendErrorJSON send error json
func (h *Handler) sendErrorJSON(
	rw http.ResponseWriter,
	catched error,
) {
	result := models.Result{
		Success: false,
		Message: catched.Error(),
	}

	if b, err := result.MarshalJSON(); err == nil {
		rw.Write(b)
	}
}

// SendSuccessJSON send object json
func (h *Handler) sendSuccessJSON(
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
func (h *Handler) Route(
	rw http.ResponseWriter,
	r *http.Request,
	mHr MethodHandlers,
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
		h.logger.Println(
			place+" method not allowed:",
			r.Method,
		)
		result = &models.RequestResult{
			Code: http.StatusMethodNotAllowed,
			Err:  h.trace.New(ErrMethodNotAllowed),
		}
	}
	h.SendResult(rw, *result)
}

func (h *Handler) Send(
	rw http.ResponseWriter,
	request *http.Request,
	result *models.RequestResult,
) {
	if result == nil {
		place := request.URL.Path
		h.logger.Println(
			place+" method not allowed:",
			request.Method,
		)
		result = &models.RequestResult{
			Code: http.StatusMethodNotAllowed,
			Err:  h.trace.New(ErrMethodNotAllowed),
		}
	}
	h.SendResult(rw, *result)
}

func (h *Handler) SendError(code int, err error) func(
	rw http.ResponseWriter,
	r *http.Request,
) models.RequestResult {
	return func(
		rw http.ResponseWriter,
		r *http.Request,
	) models.RequestResult {
		return h.Fail(code, err)
	}
}

func (h *Handler) OPTIONS() func(
	rw http.ResponseWriter,
	r *http.Request,
) models.RequestResult {
	return func(
		rw http.ResponseWriter,
		r *http.Request,
	) models.RequestResult {
		return h.Success(NoResult, nil)
	}
}
