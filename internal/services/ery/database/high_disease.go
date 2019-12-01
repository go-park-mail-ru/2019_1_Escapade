package database

import (
	"github.com/jmoiron/sqlx"
	//
	_ "github.com/jackc/pgx"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/models"
)

func (db *DB) CreateDisease(userID, projectID, sceneID int32, obj *models.Disease) error {
	obj.UserID = userID
	obj.SceneID = sceneID
	return db.workInScene(userID, projectID,
		func(tx *sqlx.Tx) error {
			return db.createDisease(tx, sceneID, obj)
		})
}

func (db *DB) UpdateDisease(userID, projectID int32, obj models.Disease) error {
	return db.workInScene(userID, projectID,
		func(tx *sqlx.Tx) error {
			return db.updateDisease(tx, &obj)
		})
}

func (db *DB) DeleteDisease(userID, projectID, objID int32) error {
	return db.workInScene(userID, projectID,
		func(tx *sqlx.Tx) error {
			return db.deleteDisease(tx, objID)
		})
}
