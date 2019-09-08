package database

import (
	//"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	//"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	//
	_ "github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
)

func (db *DB) CreateScene(userID, projectID int32, obj *models.Scene) error {
	return db.workInScene(userID, projectID,
		func(tx *sqlx.Tx) error {
			return db.createScene(tx, userID, projectID, obj)
		})
}

func (db *DB) UpdateScene(userID, projectID int32, obj models.Scene) error {
	return db.workInScene(userID, projectID,
		func(tx *sqlx.Tx) error {
			return db.updateScene(tx, &obj)
		})
}

func (db *DB) DeleteScene(userID, projectID, objID int32) error {
	return db.workInScene(userID, projectID,
		func(tx *sqlx.Tx) error {
			return db.deleteScene(tx, objID)
		})
}

func (db *DB) GetScene(userID, projectID, sceneID int32) (models.Scene, error) {
	var scene models.Scene
	tx, err := db.db.Beginx()
	if err != nil {
		return scene, err
	}
	defer tx.Rollback()

	// токен пользователя, который инициализировал действие
	_, err = db.GetProjectToken(userID, projectID)
	if err != nil {
		return scene, re.UserInProjectNotFoundWrapper(err)
	}

	scene, err = db.getScene(tx, sceneID)
	if err != nil {
		return scene, err
	}

	err = tx.Commit()
	return scene, err
}

func (db *DB) GetSceneWithObjects(userID, projectID, sceneID int32) (models.SceneWithObjects, error) {
	var scene models.SceneWithObjects
	tx, err := db.db.Beginx()
	if err != nil {
		return scene, err
	}
	defer tx.Rollback()

	// токен пользователя, который инициализировал действие
	_, err = db.GetProjectToken(userID, projectID)
	if err != nil {
		return scene, re.UserInProjectNotFoundWrapper(err)
	}

	scene.Diseases, err = db.getSceneDiseases(tx, sceneID)
	if err != nil {
		return scene, err
	}

	scene.Files, err = db.getSceneEryObjects(tx, sceneID)
	if err != nil {
		return scene, err
	}
	scene.Files = GetImages(scene.Files)

	scene.Erythrocytes, err = db.getSceneErythrocytes(tx, sceneID)
	if err != nil {
		return scene, err
	}

	scene.Scene, err = db.getScene(tx, sceneID)
	if err != nil {
		return scene, err
	}

	err = tx.Commit()
	return scene, err
}

/* inSceneFunc - функция работы с объектом в БД. Принимает на вход
транзакцию. Возвращает ошибку.

Не принимает "аргументом" объект, с которым связана работа, или
дополнительные параметры, так как этот объект(и все необходимые
дополнительные параметры) передаются в контексте вышестоящей функции
(т.е. inSceneFunc должна быть вложенной функцией)
*/
type inSceneFunc func(tx *sqlx.Tx) error

/* createInScene - создать объект в сцене
Сначала функция проверяет, что пользователь является участником проекта.
В случае, если не является функция вернёт ошибку UserInProjectNotFoundWrapper
Затем функция проверяет, что у пользователя есть право редактировать сцену
В случае, если у пользователя нет соответствующих прав, функция вернёт
ошибку ProjectNotAllowed

В случае успешного прохождения обеих проверок будет выполнена функция
inSceneCreate, создающая объект в БД
*/
func (db *DB) workInScene(userID, projectID int32, call inSceneFunc) error {
	// имеет ли данный пользователь право управлять сценами
	_, err := db.CheckToken(userID, projectID, func(token models.ProjectToken) bool {
		return token.CanEditScene()
	})
	if err != nil {
		return err
	}

	tx, err := db.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = call(tx); err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

//GetImages get image from image storage and set it to every user
func GetImages(objs []models.EryObject) []models.EryObject {
	if len(objs) == 0 {
		return objs
	}
	var err error
	for i, object := range objs {
		objs[i].Path, err = photo.GetImageFromS3(object.Path)
		if err != nil {
			utils.Debug(false, "catched error", err.Error())
		}
	}
	return objs
}

// 306 -> 132
