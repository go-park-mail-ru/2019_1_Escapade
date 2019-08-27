package database

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/return_errors"

	_ "github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
)

func (db *DB) createProject(tx *sqlx.Tx, project models.Project) (int32, error) {
	sqlInsert := `
		INSERT INTO Projects(public_access, company_access,
			 public_edit, company_edit, about) VALUES
			(:public_access, :company_access, :public_edit, :company_edit,
				:about) returning id;
			`
	id, err := createAndReturnID(tx, sqlInsert, project)
	/*

		rows, err := tx.NamedQuery(sqlInsert, project)
		if err != nil {
			return 0, err
		}
		defer rows.Close()

		var id int32
		for rows.Next() {
			err = rows.Scan(&id)
		}
	*/

	return id, err
}

func (db *DB) createProjectToken(tx *sqlx.Tx, token models.ProjectToken) (int32, error) {

	sqlInsert := `
	INSERT INTO ProjectTokens(owner, access, edit_name, edit_info, edit_access, edit_scene, edit_members_list, edit_members_token) VALUES
		(:owner, :access, :edit_name, :edit_info, :edit_access, :edit_scene, :edit_members_list, :edit_members_token) returning id;
		`
	id, err := createAndReturnID(tx, sqlInsert, token)
	/*
		rows, err := tx.NamedQuery(sqlInsert, token)
		if err != nil {
			return 0, err
		}
		defer rows.Close()

		var id int32
		for rows.Next() {
			err = rows.Scan(&id)
		}*/
	return id, err
}

func (db *DB) createUserInProject(tx *sqlx.Tx, userInProject models.UserInProject) (int32, error) {

	sqlInsert := `
	INSERT INTO UsersInProjects(position, user_id, token_id, project_id, from, to, user_Confirmed, project_confirmed) VALUES
		(:position, :user_id, :token_id, :project_id, :from, :to, :user_Confirmed, :project_confirmed) returning id;
		`
	id, err := createAndReturnID(tx, sqlInsert, userInProject)
	/*
		rows, err := tx.NamedQuery(sqlInsert, userInProject)
		if err != nil {
			return 0, err
		}
		defer rows.Close()

		var id int32
		for rows.Next() {
			err = rows.Scan(&id)
		}*/
	return id, err
}

func (db *DB) getProjectListByUserID(tx *sqlx.Tx, userID int32) ([]models.Project, error) {

	statement := `
	select p.id, p.name, p.about, p.add, up.user_confirmed, up.project_confirmed from Projects as p
	join UsersInProjects as up on p.id=up.project_id
	WHERE up.user_id = $1 
	ORDER BY up.from, up.to
		`

	rows, err := tx.Queryx(statement, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]models.Project, 0)

	for rows.Next() {
		var project models.Project
		if err = rows.StructScan(&project); err != nil {
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
func (db *DB) getProjectMember(tx *sqlx.Tx, userID, projectID int32) (models.Projectmember, error) {
	var member models.Projectmember
	statement := `
	select u.id, u.name, u.photo_title, up.position, up.from, up.to,
	up.user_confirmed, up.project_confirmed, pt.owner, pt.edit_name,
	pt.edit_info, pt.edit_access, pt.edit_scene, pt.edit_members_list,
	pt.edit_members_token from Users as u
join UsersInProjects as up on u.id=up.user_id
join ProjectToken as pt on pt.id=up.token_id
WHERE up.project_id = $1 and 
`
	rows, err := tx.Query(statement, userID)
	if err != nil {
		return member, err
	}
	defer rows.Close()

	var length int
	for rows.Next() {
		length++
		var user models.UserInProject
		var token models.ProjectToken
		err = rows.Scan(member.ID, member.Name, member.PhotoTitle,
			user.Position, user.From, user.To, user.UserConfirmed,
			user.ProjectConfirmed, token.Owner, token.EditName,
			token.EditInfo, token.EditAccess, token.EditScene,
			token.EditMembersList, token.EditMembersToken)
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

func (db *DB) getProjectToken(tx *sqlx.Tx, userID, projectID int32) (models.ProjectToken, error) {
	var token models.ProjectToken
	statement := `
	select pt.owner, pt.edit_name, pt.edit_info, pt.edit_access,
		pt.edit_scene, pt.edit_members_list, pt.edit_members_token 
			from ProjectToken as pt
		join UsersInProjects as up on pt.id=up.token_id
		WHERE up.project_id = $1 and up.user_id = $2
`
	rows, err := tx.Query(statement, projectID, userID)
	if err != nil {
		return token, err
	}
	defer rows.Close()

	var length int
	for rows.Next() {
		length++
		err = rows.Scan(token, token.EditName,
			token.EditInfo, token.EditAccess, token.EditScene,
			token.EditMembersList, token.EditMembersToken)
		if err != nil {
			break
		}
	}
	if err != nil {
		return token, err
	}
	if length != 1 {
		return token, re.UserInProjectNotFound()
	}
	return token, err
}

func (db *DB) getProject(tx *sqlx.Tx, projectID int32) (models.Project, error) {
	statement := `
	select * from Projects as p
	WHERE up.project_id = $1`
	var (
		row     = tx.QueryRowx(statement, projectID)
		project models.Project
	)

	err := row.StructScan(&project)
	return project, err
}

func (db *DB) getProjectMembers(tx *sqlx.Tx, projectID int32, needTokens bool) ([]models.Projectmember, error) {

	statement := `
	select u.id, u.name, u.photo_title, up.position, up.from, up.to,
		up.user_confirmed, up.project_confirmed, pt.owner, pt.edit_name,
		pt.edit_info, pt.edit_access, pt.edit_scene, pt.edit_members_list,
		pt.edit_members_token from Users as u
	join UsersInProjects as up on u.id=up.user_id
	join ProjectToken as pt on pt.id=up.token_id
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
		err = rows.Scan(member.ID, member.Name, member.PhotoTitle,
			user.Position, user.From, user.To, user.UserConfirmed,
			user.ProjectConfirmed, token.Owner, token.EditName,
			token.EditInfo, token.EditAccess, token.EditScene,
			token.EditMembersList, token.EditMembersToken)
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
	SET project_confirmed = true, user_confirmed = true, position = "не указана"
	WHERE project_id = $1 and user_id = $2
		`

	_, err := tx.Exec(sqlInsert, projectID, userID)
	return err
}

func (db *DB) deleteUserFromProject(tx *sqlx.Tx, projectID, userID int32) error {

	sqlInsert := `
	DELETE FROM UsersInProjects
	WHERE project_id = $1 and user_id = $2
		`

	_, err := tx.Exec(sqlInsert, projectID, userID)
	return err
}

func (db *DB) deleteProject(tx *sqlx.Tx, projectID int32) error {

	sqlInsert := `
	DELETE FROM Projects
	WHERE project_id = $1
		`

	_, err := tx.Exec(sqlInsert, projectID)
	return err
}

func (db *DB) updateProject(tx *sqlx.Tx, project *models.Project) error {
	statement := `
	UPDATE Projects 
	SET name = :name, public_access = :public_access, company_access = :company_access.
		public_edit = :public_edit, company_edit = :company_edit, about = :about
	WHERE id = :id
	`
	_, err := tx.NamedExec(statement, project)
	return err
}

func (db *DB) updateProjectToken(tx *sqlx.Tx, tokenID int32, token *models.ProjectToken) error {
	token.ID = tokenID

	statement := `
	UPDATE ProjectTokens 
	SET owner = :owner, edit_name = :edit_name, edit_info = :edit_info.
		edit_access = :edit_access, edit_scene = :edit_scene,
		edit_members_list = :edit_members_list, edit_members_token = :edit_members_token
	WHERE id = :id
	`
	_, err := tx.NamedExec(statement, token)
	return err
}

func (db *DB) updateUserInProject(tx *sqlx.Tx, userID, projectID int32, user *models.UserInProject) error {
	user.UserID = userID
	user.ProjectID = projectID

	statement := `
	UPDATE UsersInCompany 
	SET position = :position, from = :from, to = :to.
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
	scene.Edit = time.Unix(0, 0)
	statement := `
	INSERT INTO Scene(name, about, project_id, user_id, edit, editor_id) VALUES
	   (:name, :about, :project_id. :user_id, :edit, :editor_id) returning id;
	`
	id, err := createAndReturnID(tx, statement, scene)
	/*
		rows, err := tx.NamedQuery(statement, scene)
		if err != nil {
			return 0, err
		}
		defer rows.Close()

		var sceneID int32
		for rows.Next() {
			err = rows.Scan(&sceneID)
		}*/
	scene.ID = id
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
	statement := `
	DELETE FROM Scene
	WHERE id = $1
	`
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
	err := row.Scan(scene.ID, scene.UserID, scene.UserName, scene.UserPhoto,
		scene.Name, scene.About, scene.ProjectID, scene.Edit, scene.EditorID,
		scene.Add)
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
		err := rows.Scan(scene.ID, scene.UserID, scene.UserName, scene.UserPhoto,
			scene.Name, scene.About, scene.ProjectID, scene.Edit, scene.EditorID,
			scene.Add)
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
	/*
		rows, err := tx.NamedQuery(statement, obj)
		if err != nil {
			return 0, err
		}
		defer rows.Close()

		var id int32
		for rows.Next() {
			err = rows.Scan(&id)
		}
	*/
	obj.ID = id
	return err
}

func (db *DB) updateEryObject(tx *sqlx.Tx, obj *models.EryObject) error {

	statement := `
	UPDATE EryObject 
	SET user_id = :user_id, scene_id = :scene_id, path = :path.
		name = :name, about = :about, source = :source. public = :public,
		is_form = :is_form, is_texture = :is_texture, is_image = :is_image
	WHERE id = :id
	`
	_, err := tx.NamedExec(statement, obj)
	return err
}

func (db *DB) deleteEryObject(tx *sqlx.Tx, objID int32) error {
	statement := `
	DELETE FROM EryObject
	WHERE id = $1
	`
	_, err := tx.Exec(statement, objID)
	return err
}

func (db *DB) getSceneEryObjects(tx *sqlx.Tx, sceneID int32) ([]models.EryObject, error) {
	statement := `
	select * from EryObject where scene_id=$1
	`
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
	/*
		rows, err := tx.NamedQuery(statement, obj)
		if err != nil {
			return 0, err
		}
		defer rows.Close()

		var id int32
		for rows.Next() {
			err = rows.Scan(&id)
		}*/

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
	statement := `
	DELETE FROM Diseases
	WHERE id = $1
	`
	_, err := tx.Exec(statement, objID)
	return err
}

func (db *DB) getSceneDiseases(tx *sqlx.Tx, sceneID int32) ([]models.Disease, error) {
	statement := `
	select * from Diseases where scene_id=$1
	`
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
	INSERT INTO Erythrocytes(user_id, texture_id, form_id, scene_id,
		disease_id, size_x, size_y, size_z, angle_x, angle_y,
		angle_z, scale_x, scale_y, scale_z, position_x, position_y,
		position_z, form, oxygen, gemoglob) VALUES
			(:user_id, :texture_id, :form_id, :scene_id, :disease_id,
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
	SET user_id = :user_id, texture_id = :texture_id, form_id = :form_id.
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
	statement := `
	DELETE FROM Erythrocytes
	WHERE id = $1
	`
	_, err := tx.Exec(statement, objID)
	return err
}

func (db *DB) getSceneErythrocytes(tx *sqlx.Tx, sceneID int32) ([]models.Erythrocyte, error) {
	statement := `
	select * from Erythrocytes where scene_id=$1
	`
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

func (db *DB) searchProjectsWithName(tx *sqlx.Tx, name string) ([]models.Project, error) {
	statement := `
	select  id, name, public_access, company_access, about from Projects where name ~*  '.$1.'
	`
	rows, err := tx.Queryx(statement, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]models.Project, 0)

	for rows.Next() {
		var project models.Project
		if err = rows.StructScan(&project); err != nil {
			break
		}
		projects = append(projects, project)
	}
	if err != nil {
		return projects, err
	}
	return projects, err
}

func (db *DB) getAllProjects(tx *sqlx.Tx) ([]models.Project, error) {
	statement := `
	select  id, name, public_access, company_access, about from Projects
	`
	rows, err := tx.Queryx(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]models.Project, 0)

	for rows.Next() {
		var user models.Project
		if err = rows.StructScan(&user); err != nil {
			break
		}
		projects = append(projects, user)
	}
	if err != nil {
		return projects, err
	}
	return projects, err
}
