package eryhandlers

import (
	// 	"github.com/go-park-mail-ru/2019_1_Escapade/internal/auth"
	// 	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/database"

	// 	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	// 	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"

	// re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	// erydb "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/models"

	"net/http"
	// "github.com/gorilla/mux"
	// mi "github.com/go-park-mail-ru/2019_1_Escapade/internal/middleware"
)

func (H *Handler) projectsCreate(rw http.ResponseWriter, r *http.Request, userID int32) api.Result {
	const place = "projectsCreate"
	var project models.Project

	err := api.ModelFromRequest(r, &project)
	if err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	projectInfo, err := H.DB.ProjectListCreate(project, userID)
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

func (h *Handler) projectsSearch(rw http.ResponseWriter, r *http.Request) api.Result {

	const place = "projectsSearch"
	var (
		projects models.Projects
		err      error
		name     string
	)

	name = r.FormValue("name")

	if projects, err = h.DB.GetProjects(name); err != nil {
		api.NewResult(http.StatusNotFound, place, nil, err)
	}

	return api.NewResult(http.StatusOK, place, &projects, nil)
}
