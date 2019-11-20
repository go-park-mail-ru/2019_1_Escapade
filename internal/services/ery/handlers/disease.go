package eryhandlers

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/models"

	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
)

func (H *Handler) diseaseCreate(rw http.ResponseWriter, r *http.Request,
	userID, projectID, sceneID int32) api.Result {
	const place = "diseaseCreate"

	var disease models.Disease
	err := api.ModelFromRequest(r, &disease)
	if err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	err = H.DB.CreateDisease(userID, projectID, sceneID, &disease)
	if err != nil {
		return api.NewResult(http.StatusInternalServerError, place, nil, err)
	}

	return api.NewResult(http.StatusCreated, place, &disease, err)
}

func (H *Handler) diseaseUpdate(rw http.ResponseWriter, r *http.Request,
	userID, projectID, sceneID, objectID int32) api.Result {
	return api.UpdateModel(r, &models.DiseaseUpdate{}, "diseaseUpdate", false,
		func(userID int32) (api.JSONtype, error) {
			return H.DB.GetDisease(objectID)
		},
		func(interf api.JSONtype) error {
			disease, ok := interf.(*models.Disease)
			if !ok {
				return re.NoUpdate()
			}
			return H.DB.UpdateDisease(userID, projectID, *disease)
		})
}

func (H *Handler) diseaseDelete(rw http.ResponseWriter, r *http.Request,
	userID, projectID, sceneID, objectID int32) api.Result {
	const place = "diseaseDelete"

	err := H.DB.DeleteErythrocyte(userID, projectID, objectID)
	if err != nil {
		return api.NewResult(http.StatusInternalServerError, place, nil, err)
	}

	return api.NewResult(http.StatusCreated, place, nil, err)
}
