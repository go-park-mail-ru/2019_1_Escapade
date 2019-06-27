package utils

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	"encoding/json"
	"fmt"
	"net/http"
)

// SendErrorJSON send error json
func SendErrorJSON(rw http.ResponseWriter, catched error, place string) {
	result := models.Result{
		Place:   place,
		Success: false,
		Message: catched.Error(),
	}

	json.NewEncoder(rw).Encode(result)
}

// SendSuccessJSON send object json
func SendSuccessJSON(rw http.ResponseWriter, result interface{}, place string) {
	if result == nil {
		result = models.Result{
			Place:   place,
			Success: true,
			Message: "no error",
		}
	}
	fmt.Println("result:", result)
	json.NewEncoder(rw).Encode(result)
}
