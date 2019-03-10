package api

import (
	"encoding/json"
	"escapade/internal/models"
	"fmt"
	"net/http"
)

func sendErrorJSON(rw http.ResponseWriter, catched error, place string) {
	var (
		result models.Result
		bytes  []byte
		err    error
	)

	result = models.Result{
		Place:   place,
		Success: false,
		Message: catched.Error(),
	}

	if bytes, err = json.Marshal(result); err != nil {
		fmt.Println("sendErrorJSON cant create json")
		return
	}

	fmt.Fprintln(rw, string(bytes))
	fmt.Println("sendErrorJSON sent:" + result.Message)
}

func sendSuccessJSON(rw http.ResponseWriter, result interface{}, place string) {
	var (
		bytes []byte
		err   error
	)
	if result == nil {
		result = models.Result{
			Place:   place,
			Success: true,
			Message: "no error",
		}
	}

	if bytes, err = json.Marshal(result); err != nil {
		fmt.Println("sendSuccessJSON failed")
		return
	}

	fmt.Println("sendSuccessJSON success")
	fmt.Fprintln(rw, string(bytes))
}
