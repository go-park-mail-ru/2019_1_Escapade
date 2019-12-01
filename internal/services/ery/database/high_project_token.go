package database

import (
	//
	_ "github.com/jackc/pgx"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/return_errors"
)

// CheckPermissions функция проверки прав доступа
type CheckPermissions func(models.ProjectToken) bool

/*
CheckToken - проверить наличие пользователя с идентификатором userID в проекте
с идентификатором projectID. В случае существования соответствующего токена
проводится проверка прав с помощью переданной функции check

Если пользователь не связан с проектом, функция вернёт ошибку UserInProjectNotFoundWrapper
Если токен не прошел переданную проверку(check), функцяя вернёт ошибку ProjectNotAllowed
*/
func (db *DB) CheckToken(userID, projectID int32, check CheckPermissions) (models.ProjectToken, error) {

	token, err := db.GetProjectToken(userID, projectID)
	if err != nil {
		return token, re.UserInProjectNotFoundWrapper(err)
	}
	if !check(token) {
		return token, re.ProjectNotAllowed()
	}

	return token, nil
}
