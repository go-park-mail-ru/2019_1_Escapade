package database

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	//
	_ "github.com/jackc/pgx"
)

/*
ProjectListCreate - создать проект
Создает новый проект и добавляет текущего пользователя в него с максимальными
правами доступа(поле owner=true). Возвращает информацию о проекте и его членах

project - модель проекта
userID - идентификатор пользователя, создающего проект
*/
func (db *DB) ProjectListCreate(project *models.Project, userID int32) (models.ProjectWithMembers, error) {
	var info models.ProjectWithMembers
	tx, err := db.db.Beginx()
	if err != nil {
		return info, err
	}
	defer tx.Rollback()

	project.EditorID = userID
	if err = db.createProject(tx, project); err != nil {
		return info, err
	}

	token := models.ProjectToken{
		Owner:            true,
		EditName:         true,
		EditInfo:         true,
		EditAccess:       true,
		EditScene:        true,
		EditMembersList:  true,
		EditMembersToken: true,
	}

	token.ID, err = db.createProjectToken(tx, token)
	if err != nil {
		return info, err
	}

	user := models.UserInProject{
		Position:         "Создатель",
		UserID:           userID,
		TokenID:          token.ID,
		ProjectID:        project.ID,
		From:             time.Now(),
		To:               time.Unix(0, 0),
		UserConfirmed:    true,
		ProjectConfirmed: true,
	}

	user.ID, err = db.createUserInProject(tx, user)
	if err != nil {
		return info, err
	}

	info.Project = *project
	project.Edit = time.Now()
	member := models.Projectmember{
		ID:    userID,
		User:  user,
		Token: token,
	}
	info.Members = make([]models.Projectmember, 1)
	info.Members[0] = member
	info.You = member

	err = tx.Commit()
	return info, err
}

/*
ProjectListGet - получить список проектов, принадлежащих пользователю с
идентификатором userID
*/

func (db *DB) ProjectListGet(userID int32) (models.Projects, error) {
	var list models.Projects
	tx, err := db.db.Beginx()
	if err != nil {
		return list, err
	}
	defer tx.Rollback()

	list.Projects, err = db.getProjectListByUserID(tx, userID)
	if err != nil {
		return list, err
	}
	for i, project := range list.Projects {
		list.Projects[i], err = db.ProjectGet(userID, project.Project.ID)
		if err != nil {
			break
		}
		utils.Debug(false, "list.Projects[i]", i, list.Projects[i].Project.ID, list.Projects[i].Project.Name)
	}
	if err != nil {
		return list, err
	}

	err = tx.Commit()
	return list, err
}

/*
ProjectGet получить информацию о проекте, его сценах, списке участников проекта и
о вашем токене(предоставленным вам правам доступа). В случае, если у
вас есть право управлять правами доступами других участников, в поле
members у каждого участника будет проинициализировано поле Token.
эритроциты, пользовательские настройки)
Если вы не состоите в данном проекте - не подавали заявку на вступление,
либо вашу еще не одобрили, либо вы еще не приняли приглашение - то
информация о проекте может быть скрыта(зависит от того, является ли
проект публичным - поле PublicAccess у проекта). Если проект не является
публичным и вы в нем не состоите функция вернет ошибку ProjectNotAllowed()
*/
func (db *DB) ProjectGet(userID, projectID int32) (models.ProjectWithMembers, error) {
	var projectWithMembers models.ProjectWithMembers

	tx, err := db.db.Beginx()
	if err != nil {
		return projectWithMembers, err
	}
	defer tx.Rollback()

	project, err := db.GetProject(projectID)
	if err != nil {
		return projectWithMembers, err
	}

	var (
		members, owners int32
		isOwner         bool
	)
	members, owners, isOwner, projectWithMembers.You, err = db.getMembersInfo(tx, projectID, userID)

	if err != nil {
		return projectWithMembers, err
	}

	project.MembersAmount = members
	project.OwnersAmount = owners
	project.YouOwner = isOwner
	projectWithMembers.Project = *project

	project.ScenesAmount, err = db.GetScenesInProjectAmount(projectID)
	if err != nil {
		return projectWithMembers, err
	}

	if project.PublicAccess || (err == nil && projectWithMembers.You.User.Confirmed()) {

		needTokens := projectWithMembers.You.Token.HasAccessToTokens()
		projectWithMembers.Members, err = db.getProjectMembers(tx, projectID, needTokens)
		if err != nil {
			return projectWithMembers, err
		}
		scenes, err := db.getProjectScenes(tx, projectID)
		if err != nil {
			return projectWithMembers, err
		}

		projectWithMembers.Scenes = make([]models.SceneWithObjects, len(scenes))

		for i, scene := range scenes {
			projectWithMembers.Scenes[i].Diseases, err = db.getSceneDiseases(tx, scene.ID)
			if err != nil {
				break
			}

			projectWithMembers.Scenes[i].Files, err = db.getSceneEryObjects(tx, scene.ID)
			if err != nil {
				break
			}

			projectWithMembers.Scenes[i].Erythrocytes, err = db.getSceneErythrocytes(tx, scene.ID)
			if err != nil {
				break
			}
			projectWithMembers.Scenes[i].Scene = scene
		}
		if err != nil {
			return projectWithMembers, err
		}

	}
	err = tx.Commit()
	return projectWithMembers, err
}

/*
ProjectDelete - удалить проект и все связанные с ним объекты(токены,
эритроциты, пользовательские настройки)
Выполнить данное действие может только владелец проекта. В противном
случае вернётся ошибка NotProjectOwner
*/
func (db *DB) ProjectDelete(userID, projectID int32) error {
	_, err := db.CheckToken(userID, projectID, func(token models.ProjectToken) bool {
		return token.Owner
	})
	if err != nil {
		return err
	}

	tx, err := db.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = db.deleteProject(tx, projectID); err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

/*
ProjectUpdate - обновить информацию о проекте и доступ к нему

	- Если вы владелец проекта, то принимаются все предложенные изменения.
	- Если вы не являетесь участником проекта, функция вернет ошибку
		UserInProjectNotFoundWrapper
	- Если вы являетесь участником проекта, но у вас нет ни одного
		необходимого разрешения на изменение настроек проекта, то функция
		вернет ошибку ProjectNotAllowed
	- Если у вас есть права на изменение информации о проекте, будут
		утверждены те изменения, на внесение которых вам предоставлено
		право(см. ProjectToken)

Для корректной работы все поля структуры project должны быть проинициализированы
(кроме идентификатора проекта - его можно не задавать)

userID - идентификатор пользователя, совершающего действие
projectID - идентификатор проекта, над котором совершается действие
project - структура проекта с изменениями, которые необходимо занести в БД
*/
func (db *DB) ProjectUpdate(userID, projectID int32, project *models.Project) error {

	token, err := db.CheckToken(userID, projectID, func(token models.ProjectToken) bool {
		return token.CanUpdateProjectInfo()
	})
	if err != nil {
		return err
	}

	tx, err := db.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	realProject, err := db.GetProject(projectID)
	if err != nil {
		return err
	}

	project = realProject.Update(project, token)

	project.ID = projectID
	err = db.updateProject(tx, project)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

/*
ProjectTokenUpdate - обновить токен участника проекта(права доступа)
Внести изменения в токен пользователя может только владелец проекта или
пользователь с соответствующим правом(поле EditMembersToken). Иначе функция
вернёт ошибку ProjectNotAllowed.

Если изменения будут вноситься пользователем, который не является участником
проекта или над пользователем, который не является участником проекта*, то
предложенные изменения будут отклонены с ошибкой UserInProjectNotFoundWrapper.
* - в данном контексте участником проекта может являться человек, который только
	подал заявку, но ее еще не одобрили, или человек, который был приглашен,
	но еще не согласился.

Назначить пользователя создателем или снять с должности создателя
может только другой создатель. В случае попытки проведения данных операций
не создателем, функция вернет ошибку NotProjectOwner

userID - идентификатор пользователя, совершающего действие
goalID - идентификатор пользователя, чей токен будет изменён
projectID - идентификатор проекта, над котором совершается действие
token - структура токена с изменениями, которые необходимо занести в БД
*/
func (db *DB) ProjectTokenUpdate(userID, goalID, projectID int32, newToken *models.ProjectToken) error {
	tx, err := db.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// имеет ли данный пользователь право менять токены других участников
	mainToken, err := db.CheckToken(userID, projectID, func(token models.ProjectToken) bool {
		return token.CanUpdateToken()
	})
	if err != nil {
		return err
	}

	/* Проверка, не пытаестя не "создатель" назначить кого то на должность
	создателя или снять создателя с должности */
	goalToken, err := db.CheckToken(goalID, projectID, func(token models.ProjectToken) bool {
		return token.Owner == newToken.Owner || (token.Owner != newToken.Owner && mainToken.Owner)
	})
	if err != nil {
		return err
	}

	err = db.updateProjectToken(tx, goalToken.ID, newToken)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

/*
ProjectUserUpdate - обновить информацию о пользователя(положение, сроки работы)
Внести изменения в токен пользователя может только владелец проекта или
пользователь с соответствующим правом(поле EditMembersList). Иначе функция
вернёт ошибку ProjectNotAllowed.

userID - идентификатор пользователя, совершающего действие
goalID - идентификатор пользователя, чей токен будет изменён
projectID - идентификатор проекта, над котором совершается действие
user - структура пользователя с изменениями, которые необходимо занести в БД
*/
func (db *DB) ProjectUserUpdate(userID, goalID, projectID int32, user *models.UserInProject) error {
	tx, err := db.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// имеет ли данный пользователь право менять информацию других участников
	_, err = db.CheckToken(userID, projectID, func(token models.ProjectToken) bool {
		return token.CanUpdateUser()
	})
	if err != nil {
		return err
	}

	err = db.updateUserInProject(tx, goalID, projectID, user)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

/*
MembersWork отвечает за работу с участниками проекта

add - флаг действия. Если истина, значит требуется добавить пользователя.
	Если ложь, значит требуется удалить
goalID - идентификатор пользователя над которым выполняется действие
memberID - идентификатор пользователя, который выполняет действие.
	Если goalID = memberID, значит пользователь не приглашен, а подает заявку
	(либо не отменяет приглашение, а отменяет заявку)
projectID - идентификатор проекта, над котором совершается действие

Возможны следующие ситуации:
add = true, goalID = memberID
	- пользователь goalID подал заявку на вступление
	- пользователь goalID принял приглашение
add = true, goalID != memberID
	- memberID пригласил пользователя goalID
	- memberID одобрил заявку пользователя goalID
	* Если goalID дважды подал заявку или был приглашен, функция вернет ошибку
		AlreadyApplied или AlreadyInvited соответственно
	* Если goalID итак уже является участником проекта, функция вернет ошибку
		AlreadyInProject
	* Если memberID не является участником чата или не имеет права управлять
		участниками чата, функция вернёт UserInProjectNotFoundWrapper или
		ProjectNotAllowed соответственно
add = false, goalID = memberID
	- пользователь отменил заявку на вступление
	- пользователь отклонил приглашение
	- пользователь вышел из проекта
add = false, goalID != memberID
	- участник memberID отклонил заявку пользователя goalID
	- участник memberID отменил приглашение пользователю goalID
	- участник memberID выгнал пользователя goalID
	* При попытке удалить пользователя goalID, которого нет в проекте функция
		возвращает ошибку NoSuchUserInProject
	* Если memberID не является участником чата или не имеет права управлять
		участниками чата, функция вернёт UserInProjectNotFoundWrapper или
		ProjectNotAllowed соответственно
*/
func (db *DB) MembersWork(projectID, goalID, memberID int32, add bool) error {
	tx, err := db.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	invite := goalID != memberID
	// Проверка, что memberID имеет право взаимодействовать с участниками проектаmembersWork
	if invite {
		// имеет ли данный пользователь право менять информацию других участников
		_, err = db.CheckToken(memberID, projectID, func(token models.ProjectToken) bool {
			return token.CanUpdateUser()
		})
		if err != nil {
			return err
		}
	}

	//Проверка не связан ли пользователь с данным проектом
	member, err := db.getProjectMember(goalID, projectID)
	if err != nil {
		utils.Debug(false, "совсем не связан!", goalID, projectID, err.Error())
	}
	// Связан(подал заявку/был приглашен/является участников)
	if err == nil && (member.User.UserConfirmed || member.User.ProjectConfirmed) {
		utils.Debug(false, "связан!", goalID)
		utils.Debug(false, "Юзер конфермед!", member.User.UserConfirmed)
		utils.Debug(false, "Проект конфермед!", member.User.ProjectConfirmed)
		// является участником
		if member.User.Confirmed() {
			if add {
				return re.AlreadyInProject()
			}
			err = db.deleteUserFromProject(tx, projectID, goalID)
			// Пользователь goalID подал заявку
		} else if member.User.UserConfirmed {
			utils.Debug(false, "подал заявку!")
			// добавление в проект
			if add {
				// пользователь memberID пригласил goalID
				if invite {
					// одобрить заявку
					err = db.addUserToProject(tx, projectID, goalID)
				} else {
					utils.Debug(false, "заявка уже подана!")
					// заявка уже подана, ошибка
					return re.AlreadyApplied()
				}
			} else {
				utils.Debug(false, "отклонить заявку!")
				// отклонить заявку или отменить - без разницы(разве что для логов)
				err = db.deleteUserFromProject(tx, projectID, goalID)
			}
			// Пользователь goalID приглашен
		} else {
			utils.Debug(false, "приглашен")
			// добавление в проект
			if add {
				// пользователь goalID подал заявку
				if !invite {
					// принять приглашение
					err = db.addUserToProject(tx, projectID, goalID)
				} else {
					// пользоавтель уже приглашен
					return re.AlreadyInvited()
				}
			} else {
				// отклонить приглашение или отменить - без разницы(разве что для логов)
				err = db.deleteUserFromProject(tx, projectID, goalID)
			}
		}
		if err == nil {
			err = tx.Commit()
		}
		return err
	}
	// пользователь не связан с проектом
	utils.Debug(false, "не связан!", err.Error())
	if !add {
		// удалить пользователя, которого нет в проекте нельзя
		return re.NoSuchUserInProject()
	}

	token := models.ProjectToken{}

	token.ID, err = db.createProjectToken(tx, token)
	if err != nil {
		return err
	}

	userInRole := models.UserInProject{
		Position:         "Ожидает принятие заявки",
		UserID:           goalID,
		TokenID:          token.ID,
		ProjectID:        projectID,
		From:             time.Now(),
		To:               time.Unix(0, 0),
		UserConfirmed:    true,
		ProjectConfirmed: false,
	}

	// если участник проекта memberID проводит действие над пользователем goalID
	if invite {
		userInRole.Position = "Приглашен"
		userInRole.UserConfirmed = false
		userInRole.ProjectConfirmed = true
	}

	userInRole.ID, err = db.createUserInProject(tx, userInRole)
	if err == nil {
		err = tx.Commit()
	}

	return err
}

/*
GetProjects получить список всех проектов. Если name пустое, будут получены
 все существующие проекты. В ином случае в массиве проектов будут
 присутсвовать только те проекты, в имена которых присутсвует подстрока name
*/
func (db *DB) GetProjects(userID int32, name string) (models.Projects, error) {

	var projects models.Projects
	tx, err := db.db.Beginx()
	if err != nil {
		return projects, err
	}
	defer tx.Rollback()

	if name == "" {
		projects.Projects, err = db.getAllProjects(tx)
	} else {
		projects.Projects, err = db.searchProjectsWithName(tx, name)
	}
	if err != nil {
		return projects, err
	}

	for i, project := range projects.Projects {
		var (
			members, owners, projectID int32
			isOwner                    bool
			you                        models.Projectmember
		)
		projectID = project.Project.ID
		members, owners, isOwner, you, err = db.getMembersInfo(tx, projectID, userID)

		if err != nil {
			break
		}

		projects.Projects[i].Project.MembersAmount = members
		projects.Projects[i].Project.OwnersAmount = owners
		projects.Projects[i].Project.YouOwner = isOwner
		projects.Projects[i].You = you

		projects.Projects[i].Project.ScenesAmount, err = db.GetScenesInProjectAmount(projectID)
		if err != nil {
			break
		}
		utils.Debug(false, "projects.Projects[i].You ", you)
		utils.Debug(false, "members amount", project.Project.MembersAmount)
		utils.Debug(false, "scenes amount", project.Project.ScenesAmount)
	}
	if err != nil {
		return projects, err
	}

	err = tx.Commit()
	return projects, err
}

// 666 -> 519 -> 503
