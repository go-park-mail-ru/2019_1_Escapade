package handlers

import (
	"net/http"

	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
)

// JSONtype is interface to be sent by json
type JSONtype interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

// ModelUpdate is interface to update model
type ModelUpdate interface {
	JSONtype

	Update(JSONtype) bool
}

// UpdateModel update any object in DB
// dont use for updating passwords(or hash it before set to BD)
// needAuth - if true, then userID will be taken from auth request
func UpdateModel(r *http.Request, updated ModelUpdate,
	place string, needAuth bool,
	getFromDB func(userID int32) (JSONtype, error),
	setToDB func(JSONtype) error) Result {
	var (
		userID int32
		err    error
	)

	// if need userID from auth middleware
	if needAuth {
		if userID, err = GetUserIDFromAuthRequest(r); err != nil {
			return NewResult(http.StatusUnauthorized, place, nil, re.AuthWrapper(err))
		}
	}

	// updated - new version of object - get it from request
	if err = ModelFromRequest(r, updated); err != nil {
		return NewResult(http.StatusBadRequest, place, nil, err)
	}

	// object - origin object(old version) - get it from bd
	object, err := getFromDB(userID)
	if err != nil {
		return NewResult(http.StatusBadRequest, place, nil, err)
	}

	// update origin object to new version (taking into account empty fields)
	if !updated.Update(object) {
		return NewResult(http.StatusBadRequest, place, nil, re.NoUpdate())
	}

	// try to set updated object to database
	if err = setToDB(object); err != nil {
		return NewResult(http.StatusInternalServerError, place, nil, err)
	}

	return NewResult(http.StatusOK, place, object, nil)
}
