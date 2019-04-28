package api

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	"encoding/json"
	"fmt"
	"net/http"
)

func printResult(catched error, number int, place string) {
	if catched != nil {
		fmt.Println("api/"+place+" failed(code:", number, "). Error message:"+catched.Error())
	} else {
		fmt.Println("api/"+place+" success(code:", number, ")")
	}
}

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
