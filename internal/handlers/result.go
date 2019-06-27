package api

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	"encoding/json"
	"net/http"
)

func sendErrorJSON(rw http.ResponseWriter, catched error, place string) {
	result := models.Result{
		Place:   place,
		Success: false,
		Message: catched.Error(),
	}

	json.NewEncoder(rw).Encode(result)
}

func sendSuccessJSON(rw http.ResponseWriter, result interface{}, place string) {
	if result == nil {
		result = models.Result{
			Place:   place,
			Success: true,
			Message: "no error",
		}
	}
	json.NewEncoder(rw).Encode(result)
}
