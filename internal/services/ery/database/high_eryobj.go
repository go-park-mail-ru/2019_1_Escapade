package database

import (
	//"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/models"
	"github.com/jmoiron/sqlx"

	//"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	//
	_ "github.com/jackc/pgx"
)

func (db *DB) CreateEryObject(userID, projectID, sceneID int32, obj *models.EryObject) error {
	return db.workInScene(userID, projectID,
		func(tx *sqlx.Tx) error {
			return db.createEryObject(tx, sceneID, obj)
		})
}

func (db *DB) UpdateEryObject(userID, projectID int32, obj models.EryObject) error {
	return db.workInScene(userID, projectID,
		func(tx *sqlx.Tx) error {
			return db.updateEryObject(tx, &obj)
		})
}

func (db *DB) DeleteEryObject(userID, projectID, objID int32) error {
	return db.workInScene(userID, projectID,
		func(tx *sqlx.Tx) error {
			return db.deleteScene(tx, objID)
		})
}

/*
func (db *DB) CreateEryobj(userID, projectID, sceneID int32, obj models.EryObject) (models.EryObject, error) {
	tx, err := db.db.Beginx()
	if err != nil {
		return obj, err
	}
	defer tx.Rollback()

	// токен пользователя, который инициализировал действие
	mainToken, err := db.getProjectToken(tx, userID, projectID)
	if err != nil {
		return obj, re.UserInProjectNotFoundWrapper(err)
	}

	// имеет ли данный пользователь право управлять сценами
	if !mainToken.CanEditScene() {
		return obj, re.ProjectNotAllowed()
	}

	obj.ID, err = db.createEryObject(tx, sceneID, obj)
	if err != nil {
		return obj, err
	}

	err = tx.Commit()
	return obj, err
}

func (db *DB) UpdateEryobj(userID, projectID, sceneID int32, obj models.EryObject) error {
	tx, err := db.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// токен пользователя, который инициализировал действие
	mainToken, err := db.getProjectToken(tx, userID, projectID)
	if err != nil {
		return re.UserInProjectNotFoundWrapper(err)
	}

	// имеет ли данный пользователь право управлять сценами
	if !mainToken.CanEditScene() {
		return re.ProjectNotAllowed()
	}

	err = db.updateEryObject(tx, &obj)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func (db *DB) DeleteEryobj(userID, projectID, objID int32) error {
	tx, err := db.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// токен пользователя, который инициализировал действие
	mainToken, err := db.getProjectToken(tx, userID, projectID)
	if err != nil {
		return re.UserInProjectNotFoundWrapper(err)
	}

	// имеет ли данный пользователь право управлять сценами
	if !mainToken.CanEditScene() {
		return re.ProjectNotAllowed()
	}

	err = db.deleteEryObject(tx, objID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}
*/
