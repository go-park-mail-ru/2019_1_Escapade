package eryhandlers

import (
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"net/http"
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
	userID, projectID, sceneID, objectID int32) api.Result {
	return api.UpdateModel(r, &models.ErythrocyteUpdate{}, "erythrocyteUpdate", false,
		func(userID int32) (api.JSONtype, error) {
			return H.DB.GetErythrocyte(objectID)
		},
		func(interf api.JSONtype) error {
			erythrocyte, ok := interf.(*models.Erythrocyte)
			if !ok {
				return re.NoUpdate()
			}
			utils.Debug(false, "sceling_values", erythrocyte.ScaleX, erythrocyte.ScaleY, erythrocyte.ScaleZ)
			return H.DB.UpdateErythrocyte(userID, projectID, *erythrocyte)
		})
}

// func (H *Handler) erythrocyteUpdate(rw http.ResponseWriter, r *http.Request,
// 	userID, projectID, sceneID, objectID int32) api.Result {
// 	// updated - new version of object - get it from request
// 	var place = "erythrocyteUpdate"
// 	var eu = &models.ErythrocyteUpdate{}
// 	if err := api.ModelFromRequest(r, eu); err != nil {
// 		return api.NewResult(http.StatusBadRequest, place, nil, err)
// 	}
// 	utils.Debug(false, "we get this eu", *eu.ScaleX, *eu.ScaleY, *eu.ScaleZ, *eu.SizeX, *eu.SizeY, *eu.SizeZ)

// 	// object - origin object(old version) - get it from bd
// 	object, err := H.DB.GetErythrocyte(objectID)
// 	if err != nil {
// 		return api.NewResult(http.StatusBadRequest, place, nil, err)
// 	}
// 	if eu.Update(object) {
// 		utils.Debug(false, "we had object", object.ScaleX, object.ScaleY, object.ScaleZ, object.SizeX, object.SizeY, object.SizeZ)
// 		if err = H.DB.UpdateErythrocyte(userID, projectID, *object); err != nil {
// 			return api.NewResult(http.StatusInternalServerError, place, nil, err)
// 		}
// 	}
// 	return api.NewResult(http.StatusOK, place, nil, nil)
// }

func (H *Handler) erythrocyteDelete(rw http.ResponseWriter, r *http.Request,
	userID, projectID, sceneID, objectID int32) api.Result {
	const place = "erythrocyteDelete"

	err := H.DB.DeleteErythrocyte(userID, projectID, objectID)
	if err != nil {
		return api.NewResult(http.StatusInternalServerError, place, nil, err)
	}

	return api.NewResult(http.StatusCreated, place, nil, err)
}
