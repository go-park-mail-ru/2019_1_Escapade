package models

import (
	"time"

	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
)

//easyjson:json
type Scene struct {
	ID        int32     `json:"id" db:"id"`
	UserID    int32     `json:"user_id" db:"user_id"`
	UserName  string    `json:"user_name"`
	UserPhoto string    `json:"user_photo"`
	Name      string    `json:"name" db:"name"`
	About     string    `json:"about" db:"about"`
	ProjectID int32     `json:"project_id" db:"project_id"`
	Edit      time.Time `json:"edit" db:"edit"`
	EditorID  int32     `json:"editor_id" db:"editor_id"`
	Add       time.Time `json:"add" db:"add"`
}

//easyjson:json
type SceneUpdate struct {
	Name  *string `json:"name,omitempty"`
	About *string `json:"about,omitempty"`
}

func (updated *SceneUpdate) Update(sceneI api.JSONtype) bool {
	var needUpdate bool
	switch scene := sceneI.(type) {
	case *Scene:
		updateString(&scene.Name, updated.Name, &needUpdate)
		updateString(&scene.About, updated.About, &needUpdate)
	}
	return needUpdate
}

//easyjson:json
type SceneWithObjects struct {
	Scene        Scene         `json:"scene"`
	Erythrocytes []Erythrocyte `json:"erythrocytes"`
	Files        []EryObject   `json:"files"`
	Diseases     []Disease     `json:"diseases"`
}
