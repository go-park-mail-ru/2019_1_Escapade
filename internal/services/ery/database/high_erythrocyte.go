package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/models"
	"github.com/jmoiron/sqlx"

	//
	_ "github.com/jackc/pgx"
)

func (db *DB) CreateErythrocyte(userID, projectID, sceneID int32, obj *models.Erythrocyte) error {
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
