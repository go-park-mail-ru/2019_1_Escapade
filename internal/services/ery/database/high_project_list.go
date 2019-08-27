package database

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	//
	_ "github.com/jackc/pgx"
)

func (db *DB) ProjectListCreate(project models.Project, userID int32) (models.ProjectWithMembers, error) {
	var info models.ProjectWithMembers
	tx, err := db.db.Beginx()
	if err != nil {
		return info, err
	}
	defer tx.Rollback()

	project.ID, err = db.createProject(tx, project)
	if err != nil {
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

	info.Project = project
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

/////////// Вступление в другой проект

/*
ProjectListApply - подать заявку о вступлении в существующий проект
При попытке повторной подачи заявки в проект(с учетом, что предыдущая
заявка еще не была рассмотрена) функция вернет ошибку

При попытке подать заявку на вступление в проект, в который данный
пользователь приглашен, функция вернет ошибку. Чтобы принять приглашение
используйте функцию ProjectListAcceptInvitation

userID - идентификатор пользователя, совершающего действие
projectID - идентификатор проекта, над котором совершается действие
*/
/*
func (db *DB) ProjectListApply(projectID, userID int32) error {
	tx, err := db.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//Проверка не состоит ли пользователь в данном проекте
	member, err := db.getProjectMember(tx, userID, projectID)
	// Пользователь состоит
	if err == nil {

	}

	token := models.ProjectToken{}

	token.ID, err = db.createProjectToken(tx, token)
	if err != nil {
		return err
	}

	userInRole := models.UserInProject{
		Position:         "Ожидает принятие заявки",
		UserID:           userID,
		TokenID:          token.ID,
		ProjectID:        projectID,
		From:             time.Now(),
		To:               time.Unix(0, 0),
		UserConfirmed:    true,
		ProjectConfirmed: false,
	}

	userInRole.ID, err = db.createUserInProject(tx, userInRole)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}
*/
/*
ProjectListAcceptInvitation - принять приглашение вступить в
существующий проект

userID - идентификатор пользователя, совершающего действие
projectID - идентификатор проекта, над котором совершается действие
*/
/*
func (db *DB) ProjectListAcceptInvitation(projectID, userID int32) error {
	tx, err := db.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = db.acceptInvitationToProject(tx, projectID, userID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}
*/
/*
ProjectListAcceptInvitation - отказаться от приглашения в
существующий проект(или отменить поданную заявку)

userID - идентификатор пользователя, совершающего действие
projectID - идентификатор проекта, над котором совершается действие
*/
/*
func (db *DB) ProjectListRefuseInvitation(projectID, userID int32) error {
	tx, err := db.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = db.deleteUserFromProject(tx, projectID, userID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}
*/
/////////// Добавление других людей в своей проект

/*
ProjectListApply - подать заявку о вступлении в существующий проект
При попытке повторной повторного приглашения в проект(с учетом, что
предыдущее приглашение не было рассмотрена) функция вернет ошибку.

При попытке пригласить пользователя, который сам подал заявку, функция
вернет ошибку. Чтобы принять заявку Используйте функцию ProjectListAcceptApply

userID - идентификатор пользователя, совершающего действие
goalID - идентификатор пользователя, который будет приглашен
projectID - идентификатор проекта, над котором совершается действие
*/
/*
func (db *DB) ProjectListInvite(projectID, userID, goalID int32) error {
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

	// имеет ли данный пользователь право приглашать пользователей
	if !mainToken.CanUpdateUser() {
		return re.ProjectNotAllowed()
	}

	token := models.ProjectToken{}

	token.ID, err = db.createProjectToken(tx, token)
	if err != nil {
		return err
	}

	userInRole := models.UserInProject{
		Position:         "Приглашен",
		UserID:           goalID,
		TokenID:          token.ID,
		ProjectID:        projectID,
		From:             time.Now(),
		To:               time.Unix(0, 0),
		UserConfirmed:    false,
		ProjectConfirmed: true,
	}

	userInRole.ID, err = db.createUserInProject(tx, userInRole)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}*/

func (db *DB) ProjectListGet(userID int32) (models.ProjectsList, error) {
	var list models.ProjectsList
	tx, err := db.db.Beginx()
	if err != nil {
		return list, err
	}
	defer tx.Rollback()

	list.Projects, err = db.getProjectListByUserID(tx, userID)
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

	project, err := db.getProject(tx, projectID)
	if err != nil {
		return projectWithMembers, err
	}

	member, err := db.getProjectMember(tx, userID, projectID)
	if err != nil {
		utils.Debug(false, "getProjectMember err:", err.Error())
		member = models.Projectmember{}
	}
	projectWithMembers.You = member

	if project.PublicAccess || (err == nil && member.User.Confirmed()) {
		projectWithMembers.Project = project

		needTokens := member.Token.HasAccessToTokens()
		projectWithMembers.Members, err = db.getProjectMembers(tx, projectID, needTokens)
		if err != nil {
			return projectWithMembers, err
		}

		projectWithMembers.Scenes, err = db.getProjectScenes(tx, projectID)
		if err != nil {
			return projectWithMembers, err
		}

	} else {
		return projectWithMembers, re.ProjectNotAllowed()
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
	tx, err := db.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	member, err := db.getProjectMember(tx, userID, projectID)
	if err != nil {
		return err
	}

	if !member.Token.Owner {
		return re.NotProjectOwner()
	}

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

userID - идентификатор пользователя, совершающего действие
projectID - идентификатор проекта, над котором совершается действие
project - структура проекта с изменениями, которые необходимо занести в БД
*/
func (db *DB) ProjectUpdate(userID, projectID int32, project *models.Project) error {
	tx, err := db.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	token, err := db.getProjectToken(tx, userID, projectID)
	if err != nil {
		return re.UserInProjectNotFoundWrapper(err)
	}

	if !token.CanUpdateProjectInfo() {
		return re.ProjectNotAllowed()
	}

	if !token.Owner {
		realProject, err := db.getProject(tx, project.ID)
		if err != nil {
			return err
		}

		project = realProject.Update(project, token)
	}
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

	// токен пользователя, который инициализировал действие
	mainToken, err := db.getProjectToken(tx, userID, projectID)
	if err != nil {
		return re.UserInProjectNotFoundWrapper(err)
	}

	// имеет ли данный пользователь право менять токены других участников
	if mainToken.CanUpdateToken() {
		return re.ProjectNotAllowed()
	}

	// токен пользователя, над которым проводится действие
	goalToken, err := db.getProjectToken(tx, goalID, projectID)
	if err != nil {
		return re.UserInProjectNotFoundWrapper(err)
	}

	/* Проверка, не пытаестя не "создатель" назначить кого то на должность
	создателя или снять создателя с должности */
	if !mainToken.Owner && (goalToken.Owner != newToken.Owner) {
		return re.NotProjectOwner()
	}

	newToken.ID = goalToken.ID
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

	// токен пользователя, который инициализировал действие
	mainToken, err := db.getProjectToken(tx, userID, projectID)
	if err != nil {
		return re.UserInProjectNotFoundWrapper(err)
	}

	// имеет ли данный пользователь право менять информацию других участников
	if !mainToken.CanUpdateUser() {
		return re.ProjectNotAllowed()
	}

	err = db.updateUserInProject(tx, goalID, projectID, user)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

/*
membersWork отвечает за работу с участниками проекта

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
		// токен пользователя, который инициализировал действие
		mainToken, err := db.getProjectToken(tx, memberID, projectID)
		if err != nil {
			return re.UserInProjectNotFoundWrapper(err)
		}

		// имеет ли данный пользователь право менять информацию других участников
		if !mainToken.CanUpdateUser() {
			return re.ProjectNotAllowed()
		}
	}

	//Проверка не связан ли пользователь с данным проектом
	member, err := db.getProjectMember(tx, goalID, projectID)
	// Связан(подал заявку/был приглашен/является участников)
	if err == nil {
		// является участником
		if member.User.Confirmed() {
			if add {
				return re.AlreadyInProject()
			}
			err = db.deleteUserFromProject(tx, projectID, goalID)
			return err
		}
		// Пользователь goalID подал заявку
		if member.User.UserConfirmed {
			// добавление в проект
			if add {
				// пользователь memberID пригласил goalID
				if invite {
					// одобрить заявку
					err = db.addUserToProject(tx, projectID, goalID)
				} else {
					// заявка уже подана, ошибка
					return re.AlreadyApplied()
				}
			} else {
				// отклонить заявку или отменить - без разницы(разве что для логов)
				err = db.deleteUserFromProject(tx, projectID, goalID)
			}
			return err
		}
		// Пользователь goalID приглашен
		if member.User.ProjectConfirmed {
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
		return err
	}
	// пользователь не связан с проектом

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
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func (db *DB) GetProjects(name string) (models.Projects, error) {

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

	err = tx.Commit()
	return projects, err
}
