package database

import (
	"github.com/jmoiron/sqlx"
	//
	_ "github.com/jackc/pgx"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/models"
)

func (db *DB) CreateEryObject(userID, projectID, sceneID int32, obj *models.EryObject) error {
	obj.UserID = userID
	obj.SceneID = sceneID
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
			return db.deleteEryObject(tx, objID)
		})
}

// 119 -> 36
