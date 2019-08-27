package eryhandlers

import (
	// 	"github.com/go-park-mail-ru/2019_1_Escapade/internal/auth"
	// 	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/database"

	// 	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	// 	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/models"

	// re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	// erydb "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/database"

	"net/http"
	// "github.com/gorilla/mux"
	// mi "github.com/go-park-mail-ru/2019_1_Escapade/internal/middleware"
)

func (H *Handler) sceneCreate(rw http.ResponseWriter, r *http.Request) api.Result {
	const place = "sceneCreate"

	userID, err := api.GetUserIDFromAuthRequest(r)
	if err != nil {
		return api.NewResult(http.StatusUnauthorized, place, nil, re.AuthWrapper(err))
	}

	projectID, err := api.IDFromPath(r, "project_id")
	if err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	var scene models.Scene
	err = api.ModelFromRequest(r, &scene)
	if err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	err = H.DB.CreateScene(userID, projectID, &scene)
	if err != nil {
		return api.NewResult(http.StatusInternalServerError, place, nil, err)
	}

	return api.NewResult(http.StatusCreated, place, &scene, err)
}

func (H *Handler) sceneObjectsGet(rw http.ResponseWriter, r *http.Request,
	userID, projectID, sceneID int32) api.Result {
	const place = "sceneObjectsGet"

	scene, err := H.DB.GetSceneObjects(userID, projectID, sceneID)
	if err != nil {
		return api.NewResult(http.StatusInternalServerError, place, nil, err)
	}

	return api.NewResult(http.StatusCreated, place, &scene, err)
}

func (H *Handler) sceneUpdate(rw http.ResponseWriter, r *http.Request,
	userID, projectID, sceneID int32) api.Result {
	const place = "sceneUpdate"

	var scene models.Scene
	err := api.ModelFromRequest(r, &scene)
	if err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}
	scene.ID = sceneID

	err = H.DB.UpdateScene(userID, projectID, scene)
	if err != nil {
		return api.NewResult(http.StatusInternalServerError, place, nil, err)
	}

	return api.NewResult(http.StatusCreated, place, nil, err)
}

func (H *Handler) sceneDelete(rw http.ResponseWriter, r *http.Request,
	userID, projectID, sceneID int32) api.Result {
	const place = "sceneDelete"

	err := H.DB.DeleteScene(userID, projectID, sceneID)
	if err != nil {
		return api.NewResult(http.StatusInternalServerError, place, nil, err)
	}

	return api.NewResult(http.StatusCreated, place, nil, err)
}
