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

func (H *Handler) projectEnter(rw http.ResponseWriter, r *http.Request,
	projectID, memberID int32) api.Result {
	const place = "projectEnter"

	err := H.DB.MembersWork(projectID, memberID, memberID, true)
	if err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	return api.NewResult(http.StatusCreated, place, nil, err)
}

func (H *Handler) projectExit(rw http.ResponseWriter, r *http.Request,
	projectID, memberID int32) api.Result {
	const place = "projectExit"

	err := H.DB.MembersWork(projectID, memberID, memberID, false)
	if err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	return api.NewResult(http.StatusOK, place, nil, err)
}

func (H *Handler) projectGet(rw http.ResponseWriter, r *http.Request,
	projectID, memberID int32) api.Result {
	const place = "projectGet"

	project, err := H.DB.ProjectGet(memberID, projectID)
	if err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	return api.NewResult(http.StatusOK, place, &project, err)
}

func (H *Handler) projectUpdate(rw http.ResponseWriter, r *http.Request,
	projectID, memberID int32) api.Result {
	const place = "projectUpdate"

	var project models.Project

	err := api.ModelFromRequest(r, &project)
	if err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	err = H.DB.ProjectUpdate(memberID, projectID, &project)
	if err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	return api.NewResult(http.StatusOK, place, &project, err)
}

func (H *Handler) projectAcceptUser(rw http.ResponseWriter, r *http.Request,
	projectID, goalID, memberID int32) api.Result {
	const place = "projectAcceptUser"

	err := H.DB.MembersWork(projectID, goalID, memberID, true)
	if err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	return api.NewResult(http.StatusCreated, place, nil, err)
}

func (H *Handler) projectKickUser(rw http.ResponseWriter, r *http.Request,
	projectID, goalID, memberID int32) api.Result {
	const place = "projectKickUser"

	err := H.DB.MembersWork(projectID, goalID, memberID, false)
	if err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	return api.NewResult(http.StatusOK, place, nil, err)
}

func (H *Handler) projectUpdateUser(rw http.ResponseWriter, r *http.Request,
	projectID, goalID, memberID int32) api.Result {
	const place = "projectUpdateUser"

	var user models.UserInProject

	err := api.ModelFromRequest(r, &user)
	if err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	err = H.DB.ProjectUserUpdate(memberID, goalID, projectID, &user)
	if err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	return api.NewResult(http.StatusOK, place, nil, err)
}

func (H *Handler) projectUpdateUserToken(rw http.ResponseWriter, r *http.Request,
	projectID, goalID, memberID int32) api.Result {
	const place = "projectUpdateUserToken"

	var token models.ProjectToken

	err := api.ModelFromRequest(r, &token)
	if err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	err = H.DB.ProjectTokenUpdate(memberID, goalID, projectID, &token)
	if err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	return api.NewResult(http.StatusOK, place, nil, err)
}
