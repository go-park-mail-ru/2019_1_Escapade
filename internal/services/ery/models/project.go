package models

import (
	"time"

	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
)

//easyjson:json
type Project struct {
	ID               int32     `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	PublicAccess     bool      `json:"public_access" db:"public_access"`
	CompanyAccess    bool      `json:"company_access" db:"company_access"`
	PublicEdit       bool      `json:"public_edit" db:"public_edit"`
	CompanyEdit      bool      `json:"company_edit" db:"company_edit"`
	About            string    `json:"about" db:"about"`
	Add              time.Time `json:"add" db:"add"`
	UserConfirmed    bool      `json:"user_confirmed" db:"user_confirmed"`
	ProjectConfirmed bool      `json:"project_confirmed" db:"project_confirmed"`
	Edit             time.Time `json:"edit" db:"edit"`
	EditorID         int32     `json:"editor_id" db:"editor_id"`
	MembersAmount    int32     `json:"members_amount"`
	ScenesAmount     int32     `json:"scenes_amount"`
	OwnersAmount     int32     `json:"owners_Amount"`
	YouOwner         bool      `json:"you_owner"`
}

//easyjson:json
type ProjectUpdate struct {
	Name          *string `json:"name,omitempty"`
	PublicAccess  *bool   `json:"public_access,omitempty"`
	CompanyAccess *bool   `json:"company_access,omitempty"`
	PublicEdit    *bool   `json:"public_edit,omitempty"`
	CompanyEdit   *bool   `json:"company_edit,omitempty"`
	About         *string `json:"about,omitempty"`
}

func (updated *ProjectUpdate) Update(projectI api.JSONtype) bool {
	var needUpdate bool
	switch project := projectI.(type) {
	case *Project:
		updateString(&project.Name, updated.Name, &needUpdate)
		updateString(&project.About, updated.About, &needUpdate)

		updateBool(&project.PublicAccess, updated.PublicAccess, &needUpdate)
		updateBool(&project.CompanyAccess, updated.CompanyAccess, &needUpdate)
		updateBool(&project.PublicEdit, updated.PublicEdit, &needUpdate)
		updateBool(&project.CompanyEdit, updated.CompanyEdit, &needUpdate)
	}
	return needUpdate
}

// Для корректной работы newProject
func (project *Project) Update(newProject *Project, token ProjectToken) *Project {
	owner := token.Owner

	if owner || token.EditName {
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
	Projects []ProjectWithMembers `json:"projects"`
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

//easyjson:json
type ProjectTokenUpdate struct {
	Owner            *bool `json:"owner,omitempty"`
	EditName         *bool `json:"edit_name,omitempty"`
	EditInfo         *bool `json:"edit_info,omitempty"`
	EditAccess       *bool `json:"edit_access,omitempty"`
	EditScene        *bool `json:"edit_scene,omitempty"`
	EditMembersList  *bool `json:"edit_members_list,omitempty"`
	EditMembersToken *bool `json:"edit_members_token,omitempty"`
}

func (updated *ProjectTokenUpdate) Update(tokenI api.JSONtype) bool {
	var needUpdate bool

	switch token := tokenI.(type) {
	case *ProjectToken:
		updateBool(&token.Owner, updated.Owner, &needUpdate)
		updateBool(&token.EditName, updated.EditName, &needUpdate)
		updateBool(&token.EditInfo, updated.EditInfo, &needUpdate)
		updateBool(&token.EditAccess, updated.EditAccess, &needUpdate)
		updateBool(&token.EditScene, updated.EditScene, &needUpdate)
		updateBool(&token.EditMembersList, updated.EditMembersList, &needUpdate)
		updateBool(&token.EditMembersToken, updated.EditMembersToken, &needUpdate)
	}

	return needUpdate
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

//easyjson:json
type UserInProjectUpdate struct {
	Position *string    `json:"position,omitempty"`
	From     *time.Time `json:"from,omitempty"`
	To       *time.Time `json:"to,omitempty"`
}

func (updated *UserInProjectUpdate) Update(userInprojectI api.JSONtype) bool {
	var needUpdate bool

	switch up := userInprojectI.(type) {
	case *UserInProject:
		updateString(&up.Position, updated.Position, &needUpdate)
		updateTime(&up.From, updated.From, &needUpdate)
		updateTime(&up.To, updated.To, &needUpdate)
	}

	return needUpdate
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
	Token      ProjectToken  `json:"token,omitempty"`
}

//easyjson:json
type ProjectWithMembers struct {
	Project Project            `json:"project"`
	Members []Projectmember    `json:"members"`
	Scenes  []SceneWithObjects `json:"scenes"`
	You     Projectmember      `json:"you"`
}
