package eryhandlers

import (
	"net/http"

	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/models"
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

	var sceneWithObjects models.SceneWithObjects
	sceneWithObjects.Scene = scene
	sceneWithObjects.Erythrocytes = make([]models.Erythrocyte, 0)
	sceneWithObjects.Files = make([]models.EryObject, 0)
	sceneWithObjects.Diseases = make([]models.Disease, 0)

	return api.NewResult(http.StatusCreated, place, &sceneWithObjects, err)
}

func (H *Handler) sceneWithObjectsGet(rw http.ResponseWriter, r *http.Request,
	userID, projectID, sceneID int32) api.Result {
	const place = "sceneWithObjectsGet"

	scene, err := H.DB.GetSceneWithObjects(userID, projectID, sceneID)
	if err != nil {
		return api.NewResult(http.StatusInternalServerError, place, nil, err)
	}

	return api.NewResult(http.StatusCreated, place, &scene, err)
}

func (H *Handler) sceneUpdate(rw http.ResponseWriter, r *http.Request,
	userID, projectID, sceneID int32) api.Result {
	return api.UpdateModel(r, &models.SceneUpdate{}, "sceneUpdate", false,
		func(userID int32) (api.JSONtype, error) {
			scene, err := H.DB.GetScene(userID, projectID, sceneID)
			return &scene, err
		},
		func(interf api.JSONtype) error {
			scene, ok := interf.(*models.Scene)
			if !ok {
				return re.NoUpdate()
			}
			return H.DB.UpdateScene(userID, projectID, *scene)
		})
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
