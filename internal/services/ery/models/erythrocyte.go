package models

import (
	"net/http"
	"time"

	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/return_errors"
)

//easyjson:json
type Disease struct {
	ID       int32     `json:"id" db:"id"`
	UserID   int32     `json:"user_id" db:"user_id"`
	SceneID  int32     `json:"scene_id" db:"scene_id"`
	Form     float32   `json:"form" db:"form"`
	Oxygen   float32   `json:"oxygen" db:"oxygen"`
	Gemoglob float32   `json:"gemoglob" db:"gemoglob"`
	Add      time.Time `json:"add" db:"add"`
}

//easyjson:json
type EryObject struct {
	ID        int32     `json:"id" db:"id"`
	UserID    int32     `json:"user_id" db:"user_id"`
	SceneID   int32     `json:"scene_id" db:"scene_id"`
	Name      string    `json:"name" db:"name"`
	Path      string    `json:"path" db:"path"`
	About     string    `json:"about" db:"about"`
	Source    string    `json:"source" db:"source"`
	IsForm    bool      `json:"is_form" db:"is_form"`
	IsTexture bool      `json:"is_texture" db:"is_texture"`
	IsImage   bool      `json:"is_image" db:"is_image"`
	Public    bool      `json:"public" db:"public"`
	Add       time.Time `json:"add" db:"add"`
}

func (obj *EryObject) Set(r *http.Request, name, size, ftype, path string) error {
	obj.Name = name
	obj.About = "Размер файла:" + size + ", Расширение:" + ftype
	obj.Path = path

	IsFormString := r.FormValue("is_form")
	IsTextureString := r.FormValue("is_texture")
	IsImageString := r.FormValue("is_image")

	found := 0

	if len(IsFormString) > 0 {
		found++
		obj.IsForm = true
	}
	if len(IsTextureString) > 0 {
		found++
		obj.IsTexture = true
	}
	if len(IsImageString) > 0 {
		found++
		obj.IsImage = true
	}
	if found != 1 {
		return re.InvalidObjectType()
	}
	return nil
}

//easyjson:json
type Erythrocyte struct {
	ID        int32 `json:"id" db:"id"`
	UserID    int32 `json:"user_id" db:"user_id"`
	TextureID int32 `json:"texture_id" db:"texture_id"`
	FormID    int32 `json:"form_id" db:"form_id"`
	SceneID   int32 `json:"scene_id" db:"scene_id"`
	DiseaseID int32 `json:"disease_id" db:"disease_id"`

	SizeX float32 `json:"size_x" db:"size_x"`
	SizeY float32 `json:"size_y" db:"size_y"`
	SizeZ float32 `json:"size_z" db:"size_z"`

	AngleX float32 `json:"angle_x" db:"angle_x"`
	AngleY float32 `json:"angle_y" db:"angle_y"`
	AngleZ float32 `json:"angle_z" db:"angle_z"`

	ScaleX float32 `json:"scale_x" db:"scale_x"`
	ScaleY float32 `json:"scale_y" db:"scale_y"`
	ScaleZ float32 `json:"scale_z" db:"scale_z"`

	PositionX float32 `json:"position_x" db:"position_x"`
	PositionY float32 `json:"position_y" db:"position_y"`
	PositionZ float32 `json:"position_z" db:"position_z"`

	Form     float32 `json:"form" db:"form"`
	Oxygen   float32 `json:"oxygen" db:"oxygen"`
	Gemoglob float32 `json:"gemoglob" db:"gemoglob"`

	Add time.Time `json:"add" db:"add"`
}
