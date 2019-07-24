package utils

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	"net/http"
)

// SendErrorJSON send error json
func SendErrorJSON(rw http.ResponseWriter, catched error, place string) {
	result := models.Result{
		Place:   place,
		Success: false,
		Message: catched.Error(),
	}

	if b, err := result.MarshalJSON(); err == nil {
		rw.Write(b)
	}
}

// JSONtype is interface to be sent by json
type JSONtype interface {
	MarshalJSON() ([]byte, error)
}

// SendSuccessJSON send object json
func SendSuccessJSON(rw http.ResponseWriter, result JSONtype, place string) {
	if result == nil {
		result = models.Result{
			Place:   place,
			Success: true,
			Message: "no error",
		}
	}
	if b, err := result.MarshalJSON(); err == nil {
		rw.Write(b)
	}
}
