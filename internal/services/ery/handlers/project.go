package eryhandlers

import (
	// 	"github.com/go-park-mail-ru/2019_1_Escapade/internal/auth"
	// 	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/database"

	// 	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	// 	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
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
	var err error

	var delete = r.FormValue("delete")
	if delete == "" {
		err = H.DB.MembersWork(projectID, memberID, memberID, false)
	} else {
		err = H.DB.ProjectDelete(memberID, projectID)
	}
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
	for i := range project.Members {
		project.Members[i].PhotoTitle, _ = photo.GetImageFromS3(project.Members[i].PhotoTitle)
	}

	return api.NewResult(http.StatusOK, place, &project, err)
}

func (H *Handler) projectUpdate(rw http.ResponseWriter, r *http.Request,
	projectID, memberID int32) api.Result {
	return api.UpdateModel(r, &models.ProjectUpdate{}, "projectUpdate", false,
		func(userID int32) (api.JSONtype, error) {
			return H.DB.GetProject(projectID)
		},
		func(interf api.JSONtype) error {
			project, ok := interf.(*models.Project)
			if !ok {
				return re.NoUpdate()
			}
			return H.DB.ProjectUpdate(memberID, projectID, project)
		})
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
	return api.UpdateModel(r, &models.UserInProjectUpdate{}, "projectUpdateUser", false,
		func(userID int32) (api.JSONtype, error) {
			return H.DB.GetUserInProject(projectID, goalID)
		},
		func(interf api.JSONtype) error {
			user, ok := interf.(*models.UserInProject)
			if !ok {
				return re.NoUpdate()
			}
			return H.DB.ProjectUserUpdate(memberID, goalID, projectID, user)
		})
}

func (H *Handler) projectUpdateUserToken(rw http.ResponseWriter, r *http.Request,
	projectID, goalID, memberID int32) api.Result {
	return api.UpdateModel(r, &models.ProjectTokenUpdate{}, "projectUpdateUserToken", false,
		func(userID int32) (api.JSONtype, error) {
			token, err := H.DB.GetProjectToken(goalID, projectID)
			return &token, err
		},
		func(interf api.JSONtype) error {
			token, ok := interf.(*models.ProjectToken)
			if !ok {
				return re.NoUpdate()
			}
			return H.DB.ProjectTokenUpdate(memberID, goalID, projectID, token)
		})
}
