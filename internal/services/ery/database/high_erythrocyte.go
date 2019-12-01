package database

import (
	"github.com/jmoiron/sqlx"
	//
	_ "github.com/jackc/pgx"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/models"
)

func (db *DB) CreateErythrocyte(userID, projectID, sceneID int32, obj *models.Erythrocyte) error {
	obj.UserID = userID
	obj.SceneID = sceneID
	return db.workInScene(userID, projectID,
		func(tx *sqlx.Tx) error {
			return db.createErythrocyte(tx, sceneID, obj)
		})
}

func (db *DB) UpdateErythrocyte(userID, projectID int32, obj models.Erythrocyte) error {
	return db.workInScene(userID, projectID,
		func(tx *sqlx.Tx) error {
			return db.updateErythrocyte(tx, &obj)
		})
}

func (db *DB) DeleteErythrocyte(userID, projectID, objID int32) error {
	return db.workInScene(userID, projectID,
		func(tx *sqlx.Tx) error {
			return db.deleteErythrocyte(tx, objID)
		})
}
