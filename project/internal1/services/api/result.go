package api

import (
	"encoding/json"
	"escapade/project/internal/models"
	"fmt"
	"net/http"
)

func sendErrorJSON(rw http.ResponseWriter, err error, place string) {
	result := models.Result{
		Place:   place,
		Success: false,
		Message: err.Error(),
	}

	bytes, erro := json.Marshal(result)

	if erro == nil {
		fmt.Fprintln(rw, string(bytes))
	}
}

func sendSuccessJSON(rw http.ResponseWriter, place string) {
	result := models.Result{
		Place:   place,
		Success: true,
		Message: "no error",
	}

	bytes, erro := json.Marshal(result)

	if erro == nil {
		fmt.Fprintln(rw, string(bytes))
	}
}
