package eryhandlers

import (
	"net/http"
	
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/models"
)

func (H *Handler) projectsCreate(rw http.ResponseWriter, r *http.Request, userID int32) api.Result {
	const place = "projectsCreate"
	var project models.Project

	err := api.ModelFromRequest(r, &project)
	if err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	projectInfo, err := H.DB.ProjectListCreate(&project, userID)
	if err != nil {
		return api.NewResult(http.StatusInternalServerError, place, nil, err)
	}

	return api.NewResult(http.StatusCreated, place, &projectInfo, err)
}

func (H *Handler) projectsGet(rw http.ResponseWriter, r *http.Request, userID int32) api.Result {
	const place = "projectsGet"

	projectsList, err := H.DB.ProjectListGet(userID)
	if err != nil {
		return api.NewResult(http.StatusInternalServerError, place, nil, err)
	}

	return api.NewResult(http.StatusCreated, place, &projectsList, err)
}

func (H *Handler) projectsSearch(rw http.ResponseWriter, r *http.Request) api.Result {

	const place = "projectsSearch"
	var (
		projects models.Projects
		err      error
		name     string
		userID   int32
	)

	name = r.FormValue("name")

	if userID, err = api.GetUserIDFromAuthRequest(r); err != nil {
		return api.NewResult(http.StatusUnauthorized, place, nil, re.AuthWrapper(err))
	}

	if projects, err = H.DB.GetProjects(userID, name); err != nil {
		return api.NewResult(http.StatusNotFound, place, nil, err)
	}

	return api.NewResult(http.StatusOK, place, &projects, nil)
}
