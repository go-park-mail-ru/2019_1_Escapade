package models

import (
	"database/sql"
	"time"
)

//easyjson:json
type ProjectsList struct {
	Projects []Project `json:"projects"`
}

//easyjson:json
type Project struct {
	ID               int32          `json:"id" db:"id"`
	Name             sql.NullString `json:"name" db:"name"`
	PublicAccess     bool           `json:"public_access" db:"public_access"`
	CompanyAccess    bool           `json:"company_access" db:"company_access"`
	PublicEdit       bool           `json:"public_edit" db:"public_edit"`
	CompanyEdit      bool           `json:"company_edit" db:"company_edit"`
	About            sql.NullString `json:"about" db:"about"`
	Add              time.Time      `json:"add" db:"add"`
	UserConfirmed    bool           `json:"user_confirmed" db:"user_confirmed"`
	ProjectConfirmed bool           `json:"project_confirmed" db:"project_confirmed"`
}

func (project *Project) Update(newProject *Project, token ProjectToken) *Project {
	if token.Owner {
		return newProject
	}

	if token.EditName {
		project.Name = newProject.Name
	}
	if token.EditInfo {
		project.About = newProject.About
	}
	if token.EditAccess {
		project.PublicAccess = newProject.PublicAccess
		project.CompanyAccess = newProject.CompanyAccess
		project.PublicEdit = newProject.PublicEdit
		project.CompanyEdit = newProject.CompanyEdit
	}

	return project

}

//easyjson:json
type Projects struct {
	Projects []Project `json:"projects"`
}

//easyjson:json
type ProjectToken struct {
	ID               int32 `json:"id" db:"id"`
	Owner            bool  `json:"owner" db:"owner"`
	EditName         bool  `json:"edit_name" db:"edit_name"`
	EditInfo         bool  `json:"edit_info" db:"edit_info"`
	EditAccess       bool  `json:"edit_access" db:"edit_access"`
	EditScene        bool  `json:"edit_scene" db:"edit_scene"`
	EditMembersList  bool  `json:"edit_members_list" db:"edit_members_list"`
	EditMembersToken bool  `json:"edit_members_token" db:"edit_members_token"`
}

// CanUpdateProjectInfo проверяет позволяет ли токен обновлять настройки проекта
func (token *ProjectToken) CanUpdateProjectInfo() bool {
	return token.Owner || token.EditName || token.EditInfo || token.EditAccess
}

// CanUpdateToken проверяет, позволяет ли токен обновлять токены других участников
func (token *ProjectToken) CanUpdateToken() bool {
	return token.Owner || token.EditMembersToken
}

// CanUpdateUser проверяет, позволяет ли токен обновлять информацию других участников
func (token *ProjectToken) CanUpdateUser() bool {
	return token.Owner || token.EditMembersList
}

// CanEditScene проверяет, позволяет ли токен управлять сценами
func (token *ProjectToken) CanEditScene() bool {
	return token.Owner || token.EditScene
}

// HasAccessToTokens проверяет, позволяет ли токен обновлять другие токены
func (token *ProjectToken) HasAccessToTokens() bool {
	return token.Owner || token.EditMembersToken
}

//easyjson:json
type UserInProject struct {
	ID               int32     `json:"id" db:"id"`
	Position         string    `json:"position" db:"position"`
	UserID           int32     `json:"user_id" db:"user_id"`
	TokenID          int32     `json:"token_id" db:"token_id"`
	ProjectID        int32     `json:"project_id" db:"project_id"`
	From             time.Time `json:"from" db:"from"`
	To               time.Time `json:"to" db:"to"`
	UserConfirmed    bool      `json:"user_confirmed" db:"user_confirmed"`
	ProjectConfirmed bool      `json:"project_confirmed" db:"project_confirmed"`
}

func (user *UserInProject) Confirmed() bool {
	return user.UserConfirmed && user.ProjectConfirmed
}

//easyjson:json
type Projectmember struct {
	ID         int32         `json:"id" db:"id"`
	Name       string        `json:"name" db:"name" maxLength:"30" example:"John"`
	PhotoTitle string        `json:"photo_title" db:"photo_title" maxLength:"40" example:"image12.jpg" `
	User       UserInProject `json:"user"`
	Token      ProjectToken  `json:"token"`
}

//easyjson:json
type ProjectWithMembers struct {
	Project Project         `json:"project"`
	Members []Projectmember `json:"members"`
	Scenes  []Scene         `json:"scene"`
	You     Projectmember   `json:"you"`
}
