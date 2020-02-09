package handlers

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// UpdateModel update any object in DB
// dont use for updating passwords(or hash it before set to BD)
// needAuth - if true, then userID will be taken from auth request
func UpdateModel(
	r *http.Request,
	updated models.ModelUpdate,
	needAuth bool,
	getFromDB func(userID int32) (models.JSONtype, error),
	setToDB func(models.JSONtype) error,
	trace infrastructure.ErrorTrace,
) models.RequestResult {
	var (
		userID int32
		err    error
	)

	// if need userID from auth middleware
	if needAuth {
		userID, err = GetUserIDFromAuthRequest(r, trace)
		if err != nil {
			return NewResult(
				http.StatusUnauthorized,
				nil,
				err,
			)
		}
	}

	// updated - new version of object - get it from request
	err = ModelFromRequest(r, trace, updated)
	if err != nil {
		return NewResult(http.StatusBadRequest, nil, err)
	}

	// object - origin object(old version) - get it from bd
	object, err := getFromDB(userID)
	if err != nil {
		return NewResult(http.StatusBadRequest, nil, err)
	}

	// update origin object to new version (taking into account empty fields)
	if !updated.Update(object) {
		return NewResult(
			http.StatusBadRequest,
			nil,
			trace.New(ErrNoUpdate),
		)
	}

	// try to set updated object to database
	if err = setToDB(object); err != nil {
		return NewResult(
			http.StatusInternalServerError,
			nil,
			err,
		)
	}

	return NewResult(http.StatusOK, object, nil)
}
