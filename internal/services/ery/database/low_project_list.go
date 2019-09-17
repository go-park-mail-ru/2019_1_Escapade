package database

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	//
	_ "github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
)

func (db *DB) createProject(tx *sqlx.Tx, project *models.Project) error {
	sqlInsert := `
		INSERT INTO Projects(name, public_access, company_access,
			 public_edit, company_edit, about, editor_id) VALUES
			(:name, :public_access, :company_access, :public_edit, :company_edit,
				:about, :editor_id) returning *;
			`
	return createAndReturnStruct(tx, sqlInsert, project)
}

func (db *DB) createProjectToken(tx *sqlx.Tx, token models.ProjectToken) (int32, error) {

	sqlInsert := `
	INSERT INTO ProjectTokens(owner, edit_name, edit_info, edit_access, edit_scene, edit_members_list, edit_members_token) VALUES
		(:owner, :edit_name, :edit_info, :edit_access, :edit_scene, :edit_members_list, :edit_members_token) returning id;
		`
	id, err := createAndReturnID(tx, sqlInsert, token)
	return id, err
}

func (db *DB) createUserInProject(tx *sqlx.Tx, userInProject models.UserInProject) (int32, error) {

	sqlInsert := `
	INSERT INTO UsersInProjects(position, user_id, token_id, project_id, "from", "to", user_confirmed, project_confirmed) VALUES
		(:position, :user_id, :token_id, :project_id, :from, :to, :user_confirmed, :project_confirmed) returning id;
		`
	id, err := createAndReturnID(tx, sqlInsert, userInProject)
	return id, err
}

func (db *DB) getProjectListByUserID(tx *sqlx.Tx, userID int32) ([]models.ProjectWithMembers, error) {

	statement := `
	select p.id, p.name, p.about, p.add, up.user_confirmed, up.project_confirmed from Projects as p
	join UsersInProjects as up on p.id=up.project_id
	WHERE up.user_id = $1 
	ORDER BY up."from"
		`
	rows, err := tx.Queryx(statement, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]models.ProjectWithMembers, 0)

	for rows.Next() {
		var project models.ProjectWithMembers
		if err = rows.StructScan(&project.Project); err != nil {
			break
		}
		projects = append(projects, project)
	}
	if err != nil {
		return nil, err
	}
	return projects, err
}

// checkCanAccess возвращает структуру участника проекта
func (db *DB) getProjectMember(userID, projectID int32) (models.Projectmember, error) {
	var member models.Projectmember
	statement := `
	select u.id, u.name, u.photo_title, up.position, up.from, up.to,
	up.user_confirmed, up.project_confirmed, pt.owner, pt.edit_name,
	pt.edit_info, pt.edit_access, pt.edit_scene, pt.edit_members_list,
	pt.edit_members_token from Users as u
join UsersInProjects as up on u.id=up.user_id
join ProjectTokens as pt on pt.id=up.token_id
WHERE up.project_id = $1 and user_id = $2
`
	rows, err := db.db.Query(statement, projectID, userID)
	if err != nil {
		return member, err
	}
	defer rows.Close()

	var length int
	for rows.Next() {
		length++
		var user models.UserInProject
		var token models.ProjectToken
		err = rows.Scan(&member.ID, &member.Name, &member.PhotoTitle,
			&user.Position, &user.From, &user.To, &user.UserConfirmed,
			&user.ProjectConfirmed, &token.Owner, &token.EditName,
			&token.EditInfo, &token.EditAccess, &token.EditScene,
			&token.EditMembersList, &token.EditMembersToken)
		if err != nil {
			break
		}

		member.User = user
		member.Token = token
	}
	if err != nil {
		return member, err
	}
	if length != 1 {
		return member, re.UserInProjectNotFound()
	}
	return member, err
}

func (db *DB) GetProjectToken(userID, projectID int32) (models.ProjectToken, error) {
	var token models.ProjectToken
	utils.Debug(false, "we have these id:", userID, projectID)
	statement := `
	select pt.id, pt.owner, pt.edit_name, pt.edit_info, pt.edit_access,
		pt.edit_scene, pt.edit_members_list, pt.edit_members_token 
			from ProjectTokens as pt
		join UsersInProjects as up on pt.id=up.token_id
		WHERE up.project_id = $1 and up.user_id = $2
`
	row := db.db.QueryRowx(statement, projectID, userID)
	err := row.StructScan(&token)
	if err != nil {
		utils.Debug(false, "we have error", err.Error())
	}

	return token, err
}

func (db *DB) getProjectMembers(tx *sqlx.Tx, projectID int32, needTokens bool) ([]models.Projectmember, error) {

	statement := `
	select u.id, u.name, u.photo_title, up.position, up.from, up.to,
		up.user_confirmed, up.project_confirmed, pt.owner, pt.edit_name,
		pt.edit_info, pt.edit_access, pt.edit_scene, pt.edit_members_list,
		pt.edit_members_token from Users as u
	join UsersInProjects as up on u.id=up.user_id
	join ProjectTokens as pt on pt.id=up.token_id
	WHERE up.project_id = $1
		`

	rows, err := tx.Query(statement, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]models.Projectmember, 0)

	for rows.Next() {
		var member models.Projectmember
		var user models.UserInProject
		var token models.ProjectToken
		err = rows.Scan(&member.ID, &member.Name, &member.PhotoTitle,
			&user.Position, &user.From, &user.To, &user.UserConfirmed,
			&user.ProjectConfirmed, &token.Owner, &token.EditName,
			&token.EditInfo, &token.EditAccess, &token.EditScene,
			&token.EditMembersList, &token.EditMembersToken)
		if err != nil {
			break
		}

		member.User = user
		if needTokens {
			member.Token = token
		}
		members = append(members, member)
	}
	if err != nil {
		return nil, err
	}
	return members, err
}

func (db *DB) addUserToProject(tx *sqlx.Tx, projectID, userID int32) error {

	sqlInsert := `
	UPDATE UsersInProjects 
	SET project_confirmed = true, user_confirmed = true, position = 'не указана'
	WHERE project_id = $1 and user_id = $2
		`
	_, err := tx.Exec(sqlInsert, projectID, userID)
	if err != nil {
		utils.Debug(false, "addUserToProject err", err.Error())
	}
	return err
}

func (db *DB) deleteUserFromProject(tx *sqlx.Tx, projectID, userID int32) error {

	sqlInsert := `
DELETE FROM UsersInProjects WHERE project_id = $1 and user_id = $2
		`
	_, err := tx.Exec(sqlInsert, projectID, userID)

	_, owners, _, _, err := db.getMembersInfo(tx, projectID, userID)

	if err != nil {
		return err
	}
	if owners == 0 {
		err = db.deleteProject(tx, projectID)
	}

	return err
}

func (db *DB) deleteProject(tx *sqlx.Tx, projectID int32) error {
	sqlInsert := `DELETE FROM Projects WHERE id = $1`

	_, err := tx.Exec(sqlInsert, projectID)
	return err
}

func (db *DB) updateProject(tx *sqlx.Tx, project *models.Project) error {
	statement := `
	UPDATE Projects 
	SET name = :name, public_access = :public_access, company_access = :company_access,
		public_edit = :public_edit, company_edit = :company_edit, about = :about
	WHERE id = :id
	`
	_, err := tx.NamedExec(statement, project)
	if err != nil {
		utils.Debug(false, "updateProject error", err.Error())
	}
	return err
}

func (db *DB) updateProjectToken(tx *sqlx.Tx, tokenID int32, token *models.ProjectToken) error {
	token.ID = tokenID

	statement := `
	UPDATE ProjectTokens 
	SET owner = :owner, edit_name = :edit_name, edit_info = :edit_info,
		edit_access = :edit_access, edit_scene = :edit_scene,
		edit_members_list = :edit_members_list, edit_members_token = :edit_members_token
	WHERE id = :id
	`
	_, err := tx.NamedExec(statement, token)
	if err != nil {
		utils.Debug(false, "done with err", err.Error())
	}
	return err
}

func (db *DB) updateUserInProject(tx *sqlx.Tx, userID, projectID int32, user *models.UserInProject) error {
	user.UserID = userID
	user.ProjectID = projectID

	statement := `
	UPDATE UsersInProjects
	SET position = :position, "from" = :from, "to" = :to
	WHERE user_id = :user_id and project_id = :project_id
	`
	_, err := tx.NamedExec(statement, user)
	return err
}

// scene

func (db *DB) createScene(tx *sqlx.Tx, userID, projectID int32, scene *models.Scene) error {
	scene.ProjectID = projectID
	scene.UserID = userID
	scene.EditorID = userID
	scene.Edit = time.Now()
	statement := `
	INSERT INTO Scene(name, about, project_id, user_id, edit, editor_id) VALUES
	   (:name, :about, :project_id, :user_id, :edit, :editor_id) returning *;
	`
	err := createAndReturnStruct(tx, statement, scene)
	if err != nil {
		utils.Debug(false, "createScene err", err.Error())
	}
	return err
}

func (db *DB) updateScene(tx *sqlx.Tx, scene *models.Scene) error {

	statement := `
	UPDATE ProjectTokens 
	SET name = :name, about = :about, project_id = :project_id.
		user_id = :user_id, edit = :edit, editor_id = :editor_id
	WHERE id = :id
	`
	_, err := tx.NamedExec(statement, scene)
	return err
}

func (db *DB) deleteScene(tx *sqlx.Tx, sceneID int32) error {
	statement := `DELETE FROM Scene WHERE id = $1`
	_, err := tx.Exec(statement, sceneID)
	return err
}

func (db *DB) getScene(tx *sqlx.Tx, sceneID int32) (models.Scene, error) {
	statement := `
	select s.id, s.user_id, u.name, u.photo_title,
		s.name, s.about, s.project_id, s.edit, s.editor_id, s.add
			from Scene as s join Users as u on s.user_id=u.id
	WHERE s.id = $1 
	`
	row := tx.QueryRow(statement, sceneID)
	var scene models.Scene
	err := row.Scan(&scene.ID, &scene.UserID, &scene.UserName, &scene.UserPhoto,
		&scene.Name, &scene.About, &scene.ProjectID, &scene.Edit, &scene.EditorID,
		&scene.Add)
	return scene, err
}

func (db *DB) getProjectScenes(tx *sqlx.Tx, projectID int32) ([]models.Scene, error) {
	statement := `
	select s.id, s.user_id, u.name, u.photo_title,
		s.name, s.about, s.project_id, s.edit, s.editor_id, s.add
			from Scene as s join Users as u on s.user_id=u.id
	WHERE s.project_id = $1 
	`
	rows, err := tx.Query(statement, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scenes := make([]models.Scene, 0)

	for rows.Next() {
		var scene models.Scene
		err := rows.Scan(&scene.ID, &scene.UserID, &scene.UserName,
			&scene.UserPhoto, &scene.Name, &scene.About, &scene.ProjectID,
			&scene.Edit, &scene.EditorID, &scene.Add)
		if err != nil {
			break
		}
		scenes = append(scenes, scene)
	}
	if err != nil {
		return nil, err
	}
	return scenes, err
}

func (db *DB) createEryObject(tx *sqlx.Tx, sceneID int32, obj *models.EryObject) error {
	obj.SceneID = sceneID
	statement := `
	INSERT INTO EryObject(user_id, scene_id, path, name, about,
		source, public, is_form, is_texture, is_image) VALUES
		(:user_id, :scene_id, :path, :name, :about, :source, :public,
		:is_form, :is_texture, :is_image) returning id;
		`
	id, err := createAndReturnID(tx, statement, obj)
	obj.ID = id
	return err
}

func (db *DB) updateEryObject(tx *sqlx.Tx, obj *models.EryObject) error {

	statement := `
	UPDATE EryObject 
	SET path = :path, name = :name, about = :about, 
	source = :source, public = :public, is_form = :is_form,
	is_texture = :is_texture, is_image = :is_image
	WHERE id = :id
	`
	_, err := tx.NamedExec(statement, obj)
	return err
}

func (db *DB) deleteEryObject(tx *sqlx.Tx, objID int32) error {
	statement := `DELETE FROM EryObject WHERE id = $1`
	_, err := tx.Exec(statement, objID)
	return err
}

func (db *DB) GetEryObject(objectID int32) (*models.EryObject, error) {
	var eryObject models.EryObject
	statement := `select * from EryObject where id=$1`
	row := db.db.QueryRowx(statement, objectID)
	err := row.StructScan(&eryObject)

	return &eryObject, err
}

func (db *DB) getSceneEryObjects(tx *sqlx.Tx, sceneID int32) ([]models.EryObject, error) {
	statement := `select * from EryObject where scene_id=$1`
	rows, err := tx.Queryx(statement, sceneID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	objs := make([]models.EryObject, 0)

	for rows.Next() {
		var obj models.EryObject
		if err = rows.StructScan(&obj); err != nil {
			break
		}
		objs = append(objs, obj)
	}
	if err != nil {
		return nil, err
	}
	return objs, err
}

func (db *DB) createDisease(tx *sqlx.Tx, sceneID int32, obj *models.Disease) error {
	obj.SceneID = sceneID

	statement := `
	INSERT INTO Diseases(user_id, scene_id, form, oxygen, gemoglob) VALUES
		(:user_id, :scene_id, :form, :oxygen, :gemoglob) returning id;
		`
	id, err := createAndReturnID(tx, statement, obj)
	obj.ID = id
	return err
}

func (db *DB) updateDisease(tx *sqlx.Tx, obj *models.Disease) error {

	statement := `
	UPDATE Disease 
	SET user_id = :user_id, scene_id = :scene_id, form = :form.
		oxygen = :oxygen, gemoglob = :gemoglob
	WHERE id = :id
	`
	_, err := tx.NamedExec(statement, obj)
	return err
}

func (db *DB) deleteDisease(tx *sqlx.Tx, objID int32) error {
	statement := `DELETE FROM Diseases WHERE id = $1`
	_, err := tx.Exec(statement, objID)
	return err
}

func (db *DB) GetDisease(objectID int32) (*models.Disease, error) {
	var disease models.Disease
	statement := `select * from Diseases where id=$1`
	row := db.db.QueryRowx(statement, objectID)
	err := row.StructScan(&disease)

	return &disease, err
}

func (db *DB) getSceneDiseases(tx *sqlx.Tx, sceneID int32) ([]models.Disease, error) {
	statement := `select * from Diseases where scene_id=$1`
	rows, err := tx.Queryx(statement, sceneID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	objs := make([]models.Disease, 0)

	for rows.Next() {
		var obj models.Disease
		if err = rows.StructScan(&obj); err != nil {
			break
		}
		objs = append(objs, obj)
	}
	if err != nil {
		return nil, err
	}
	return objs, err
}

func (db *DB) createErythrocyte(tx *sqlx.Tx, sceneID int32, obj *models.Erythrocyte) error {
	obj.SceneID = sceneID

	statement := `
	INSERT INTO Erythrocytes(user_id, texture_id, form_id, image_id, scene_id,
		disease_id, size_x, size_y, size_z, angle_x, angle_y,
		angle_z, scale_x, scale_y, scale_z, position_x, position_y,
		position_z, form, oxygen, gemoglob) VALUES
			(:user_id, :texture_id, :form_id, :image_id, :scene_id, :disease_id,
			 :size_x, :size_y, :size_z, :angle_x, :angle_y, :angle_z,
			 :scale_x, :scale_y, :scale_z, :position_x, :position_y,
			 :position_z, :form, :oxygen, :gemoglob) returning id;
		`

	id, err := createAndReturnID(tx, statement, obj)
	obj.ID = id
	return err
}

func (db *DB) updateErythrocyte(tx *sqlx.Tx, obj *models.Erythrocyte) error {

	statement := `
	UPDATE Erythrocytes 
	SET user_id = :user_id, texture_id = :texture_id, image_id = :image_id, form_id = :form_id,
	scene_id = :scene_id, disease_id = :disease_id, size_x = :size_x,
	size_y = :size_y, size_z = :size_z, angle_x = :angle_x, angle_y = :angle_y,
	angle_z = :angle_z, scale_x = :scale_x, scale_y = :scale_y, scale_z = :scale_z,
	position_x = :position_x, position_y = :position_y, position_z = :position_z,
	form = :form, oxygen = :oxygen, gemoglob = :gemoglob
	WHERE id = :id
	`
	_, err := tx.NamedExec(statement, obj)
	return err
}

func (db *DB) deleteErythrocyte(tx *sqlx.Tx, objID int32) error {
	statement := `DELETE FROM Erythrocytes WHERE id = $1`
	_, err := tx.Exec(statement, objID)
	return err
}

func (db *DB) GetErythrocyte(objectID int32) (*models.Erythrocyte, error) {
	var erythrocyte models.Erythrocyte
	statement := `select * from Erythrocytes where id=$1`
	row := db.db.QueryRowx(statement, objectID)
	err := row.StructScan(&erythrocyte)
	if err != nil {
		utils.Debug(false, "GetErythrocyte cant get", err.Error())
	}

	return &erythrocyte, err
}

func (db *DB) getSceneErythrocytes(tx *sqlx.Tx, sceneID int32) ([]models.Erythrocyte, error) {
	statement := `select * from Erythrocytes where scene_id=$1`
	rows, err := tx.Queryx(statement, sceneID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	objs := make([]models.Erythrocyte, 0)

	for rows.Next() {
		var obj models.Erythrocyte
		if err = rows.StructScan(&obj); err != nil {
			break
		}
		objs = append(objs, obj)
	}
	if err != nil {
		return nil, err
	}
	return objs, err
}

// createAndReturn создает объект obj
func createAndReturnID(tx *sqlx.Tx, statement string, obj interface{}) (int32, error) {
	rows, err := tx.NamedQuery(statement, obj)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var id int32
	for rows.Next() {
		err = rows.Scan(&id)
	}

	return id, err
}

func createAndReturnStruct(tx *sqlx.Tx, statement string, obj interface{}) error {
	rows, err := tx.NamedQuery(statement, obj)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.StructScan(obj)
	}
	return err
}

func (db *DB) searchProjectsWithName(tx *sqlx.Tx, name string) ([]models.ProjectWithMembers, error) {
	statement := `
	select  id, name, public_access, company_access, about, edit, editor_id, add from Projects where 
	POSITION (lower($1) IN lower(name)) > 0 order by edit DESC;
	` // name ~* 
	rows, err := tx.Queryx(statement, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]models.ProjectWithMembers, 0)

	for rows.Next() {
		var project models.ProjectWithMembers
		if err = rows.StructScan(&project.Project); err != nil {
			break
		}
		projects = append(projects, project)
	}
	if err != nil {
		return projects, err
	}
	return projects, err
}

func (db *DB) getAllProjects(tx *sqlx.Tx) ([]models.ProjectWithMembers, error) {
	statement := `
	select id, name, public_access, company_access, about, edit, editor_id, add from Projects order by edit DESC
	`
	rows, err := tx.Queryx(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]models.ProjectWithMembers, 0)
	for rows.Next() {
		var project models.ProjectWithMembers
		if err = rows.StructScan(&project.Project); err != nil {
			break
		}
		projects = append(projects, project)
	}
	if err != nil {
		utils.Debug(false, "getAllProjects error", err.Error())
		return projects, err
	}
	return projects, err
}

func (db *DB) GetProject(projectID int32) (*models.Project, error) {
	var project models.Project
	statement := `select * from Projects where id = $1`
	row := db.db.QueryRowx(statement, projectID)
	err := row.StructScan(&project)
	return &project, err
}

func (db *DB) GetUserInProject(projectID, userID int32) (*models.UserInProject, error) {
	var userInProject models.UserInProject
	statement := `select * from UsersInProjects where project_id = $1 and user_id = $2`
	row := db.db.QueryRowx(statement, projectID, userID)
	err := row.StructScan(&userInProject)
	return &userInProject, err
}

func (db *DB) GetScenesInProjectAmount(projectID int32) (int32, error) {
	var amount int32
	statement := `select count(id) from Scene where project_id = $1`
	row := db.db.QueryRowx(statement, projectID)
	err := row.Scan(&amount)
	return amount, err
}

func (db *DB) getMembersInfo(tx *sqlx.Tx, projectID, userID int32) (int32, int32, bool, models.Projectmember, error) {
	members, err := db.getProjectMembers(tx, projectID, true)
	var you models.Projectmember
	if err != nil {
		return 0, 0, false, you, err
	}
	var (
		owners        int32
		membersAmount int32
		isOwner       bool
	)
	for _, member := range members {
		if member.Token.Owner {
			owners++
			if member.ID == userID {
				isOwner = true
			}
		}

		if member.ID == userID {
			you = member
		}
		if member.User.UserConfirmed && member.User.ProjectConfirmed {
			membersAmount++
		}
	}
	return membersAmount, owners, isOwner, you, nil
}

// 654
