package database

import (
	//"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/return_errors"

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

/*
func (db *DB) CreateScene(userID, projectID int32, scene models.Scene) (models.Scene, error) {
	tx, err := db.db.Beginx()
	if err != nil {
		return scene, err
	}
	defer tx.Rollback()

	// токен пользователя, который инициализировал действие
	mainToken, err := db.getProjectToken(tx, userID, projectID)
	if err != nil {
		return scene, re.UserInProjectNotFoundWrapper(err)
	}

	// имеет ли данный пользователь право управлять сценами
	if !mainToken.CanEditScene() {
		return scene, re.ProjectNotAllowed()
	}

	scene.ID, err = db.createScene(tx, userID, projectID, scene)
	if err != nil {
		return scene, err
	}

	err = tx.Commit()
	return scene, err
}

func (db *DB) UpdateScene(userID, projectID int32, scene models.Scene) error {
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

	err = db.updateScene(tx, &scene)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func (db *DB) DeleteScene(userID, projectID, sceneID int32) error {
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

	err = db.deleteScene(tx, sceneID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}
*/

func (db *DB) GetScene(userID, projectID, sceneID int32) (models.Scene, error) {
	var scene models.Scene
	tx, err := db.db.Beginx()
	if err != nil {
		return scene, err
	}
	defer tx.Rollback()

	// токен пользователя, который инициализировал действие
	_, err = db.getProjectToken(tx, userID, projectID)
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

func (db *DB) GetSceneObjects(userID, projectID, sceneID int32) (models.SceneObjects, error) {
	var scene models.SceneObjects
	tx, err := db.db.Beginx()
	if err != nil {
		return scene, err
	}
	defer tx.Rollback()

	// токен пользователя, который инициализировал действие
	_, err = db.getProjectToken(tx, userID, projectID)
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

	scene.Erythrocytes, err = db.getSceneErythrocytes(tx, sceneID)
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

	if err = call(tx); err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

//				deprecated							deprecated
// deprecated							deprecated
//							deprecated							deprecated

/* inSceneUpdate - функция обновления объекта в БД. Принимает на вход
транзакцию и возвращает ошибку.

Не принимает "аргументом" объект, который необходимо создать, так как
этот объект передается в контексте вышестоящей функции(т.е.
inSceneUpdate должна быть вложенной функцией)
*/
type inSceneUpdate func(tx *sqlx.Tx) error

/* updateInScene - создать объект в сцене
Сначала функция проверяет, что пользователь является участником проекта.
В случае, если не является функция вернёт ошибку UserInProjectNotFoundWrapper
Затем функция проверяет, что у пользователя есть право редактировать сцену
В случае, если у пользователя нет соответствующих прав, функция вернёт
ошибку ProjectNotAllowed

В случае успешного прохождения обеих проверок будет выполнена функция
inSceneCreate, создающая объект в БД
*/
func (db *DB) updateInScene(userID, projectID int32, update inSceneUpdate) error {
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

	if err = update(tx); err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

/* inSceneDelete - функция удаления объекта в БД. Принимает на вход
транзакцию и идентификатор объекта. Возвращает ошибку.

Не принимает "аргументом" объект, который необходимо создать, так как
этот объект передается в контексте вышестоящей функции(т.е.
inSceneDelete должна быть вложенной функцией)
*/
type inSceneDelete func(tx *sqlx.Tx) error

func (db *DB) deleteInScene(userID, projectID int32, delete inSceneDelete) error {
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

	if err = delete(tx); err != nil {
		return err
	}

	err = tx.Commit()
	return err
}
