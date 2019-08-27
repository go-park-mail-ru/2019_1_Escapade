package eryhandlers

import (
	// 	"github.com/go-park-mail-ru/2019_1_Escapade/internal/auth"
	// 	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/database"

	// 	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	// 	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/models"

	// re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	// erydb "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/database"

	"net/http"
	// "github.com/gorilla/mux"
	// mi "github.com/go-park-mail-ru/2019_1_Escapade/internal/middleware"
)

func (H *Handler) erythrocyteCreate(rw http.ResponseWriter, r *http.Request,
	userID, projectID, sceneID int32) api.Result {
	const place = "erythrocyteCreate"

	var erythrocyte models.Erythrocyte
	err := api.ModelFromRequest(r, &erythrocyte)
	if err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	err = H.DB.CreateErythrocyte(userID, projectID, sceneID, &erythrocyte)
	if err != nil {
		return api.NewResult(http.StatusInternalServerError, place, nil, err)
	}

	return api.NewResult(http.StatusCreated, place, &erythrocyte, err)
}

func (H *Handler) erythrocyteUpdate(rw http.ResponseWriter, r *http.Request,
	userID, projectID, sceneID int32) api.Result {
	const place = "erythrocyteUpdate"

	var erythrocyte models.Erythrocyte
	err := api.ModelFromRequest(r, &erythrocyte)
	if err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	err = H.DB.UpdateErythrocyte(userID, projectID, erythrocyte)
	if err != nil {
		return api.NewResult(http.StatusInternalServerError, place, nil, err)
	}

	return api.NewResult(http.StatusCreated, place, nil, err)
}

func (H *Handler) erythrocyteDelete(rw http.ResponseWriter, r *http.Request,
	userID, projectID, sceneID int32) api.Result {
	const place = "erythrocyteDelete"

	var erythrocyte models.Erythrocyte
	err := api.ModelFromRequest(r, &erythrocyte)
	if err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	err = H.DB.DeleteErythrocyte(userID, projectID, erythrocyte.ID)
	if err != nil {
		return api.NewResult(http.StatusInternalServerError, place, nil, err)
	}

	return api.NewResult(http.StatusCreated, place, nil, err)
}
