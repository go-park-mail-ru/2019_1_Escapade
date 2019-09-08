package models

import (
	"fmt"
	"net/http"
	"time"

	"strconv"

	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

//easyjson:json
type Disease struct {
	ID       int32     `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	About    string    `json:"about" db:"about"`
	UserID   int32     `json:"user_id" db:"user_id"`
	SceneID  int32     `json:"scene_id" db:"scene_id"`
	Form     float32   `json:"form" db:"form"`
	Oxygen   float32   `json:"oxygen" db:"oxygen"`
	Gemoglob float32   `json:"gemoglob" db:"gemoglob"`
	Add      time.Time `json:"add" db:"add"`
}

//easyjson:json
type DiseaseUpdate struct {
	Name     *string  `json:"name,omitempty"`
	About    *string  `json:"about,omitempty"`
	Form     *float32 `json:"form,omitempty"`
	Oxygen   *float32 `json:"oxygen,omitempty"`
	Gemoglob *float32 `json:"gemoglob,omitempty"`
}

func (updated *DiseaseUpdate) Update(diseaseI api.JSONtype) bool {
	var needUpdate bool

	switch disease := diseaseI.(type) {
	case *Disease:
		updateString(&disease.Name, updated.Name, &needUpdate)
		updateString(&disease.About, updated.About, &needUpdate)
		updateFloat32(&disease.Form, updated.Form, &needUpdate)
		updateFloat32(&disease.Oxygen, updated.Oxygen, &needUpdate)
		updateFloat32(&disease.Gemoglob, updated.Gemoglob, &needUpdate)
	}
	return needUpdate
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

//easyjson:json
type EryObjectUpdate struct {
	Name   *string `json:"name,omitempty"`
	About  *string `json:"about,omitempty"`
	Source *string `json:"source,omitempty"`
	Public *bool   `json:"public,omitempty"`
}

func (updated *EryObjectUpdate) Update(eryObjectI api.JSONtype) bool {
	var needUpdate bool

	switch eryObject := eryObjectI.(type) {
	case *EryObject:
		updateString(&eryObject.Name, updated.Name, &needUpdate)
		updateString(&eryObject.About, updated.About, &needUpdate)
		updateString(&eryObject.Source, updated.Source, &needUpdate)
		updateBool(&eryObject.Public, updated.Public, &needUpdate)
	}

	return needUpdate
}

func (obj *EryObject) Set(r *http.Request, name, size, path string) error {
	bString := "Байт"
	sizeInt, _ := strconv.Atoi(size)
	if sizeInt > 1024 {
		sizeInt /= 1024
		bString = "КБайт"
	}
	if sizeInt > 1024 {
		sizeInt /= 1024
		bString = "МБайт"
	}
	if sizeInt > 1024 {
		sizeInt /= 1024
		bString = "ГБайт"
	}
	size = strconv.Itoa(sizeInt)
	obj.Name = name
	obj.About = "Размер файла:" + size + " " + bString
	obj.Path = path

	IsFormString := r.FormValue("is_form")
	IsTextureString := r.FormValue("is_texture")
	IsImageString := r.FormValue("is_image")

	found := 0

	if len(IsFormString) > 0 {
		utils.Debug(false, "IsFormString", IsFormString)
		found++
		obj.IsForm = true
	}
	if len(IsTextureString) > 0 {
		utils.Debug(false, "IsTextureString", IsTextureString)
		found++
		obj.IsTexture = true
	}
	if len(IsImageString) > 0 {
		utils.Debug(false, "IsImageString", IsImageString)
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
	ImageID   int32 `json:"image_id" db:"image_id"`
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

//easyjson:json
type ErythrocyteUpdate struct {
	TextureID *int32 `json:"texture_id,omitempty"`
	FormID    *int32 `json:"form_id,omitempty"`
	ImageID   *int32 `json:"image_id,omitempty"`
	DiseaseID *int32 `json:"disease_id,omitempty"`

	SizeX *float32 `json:"size_x,omitempty"`
	SizeY *float32 `json:"size_y,omitempty"`
	SizeZ *float32 `json:"size_z,omitempty"`

	AngleX *float32 `json:"angle_x,omitempty"`
	AngleY *float32 `json:"angle_y,omitempty"`
	AngleZ *float32 `json:"angle_z,omitempty"`

	ScaleX *float32 `json:"scale_x,omitempty"`
	ScaleY *float32 `json:"scale_y,omitempty"`
	ScaleZ *float32 `json:"scale_z,omitempty"`

	PositionX *float32 `json:"position_x,omitempty"`
	PositionY *float32 `json:"position_y,omitempty"`
	PositionZ *float32 `json:"position_z,omitempty"`

	Form     *float32 `json:"form,omitempty"`
	Oxygen   *float32 `json:"oxygen,omitempty"`
	Gemoglob *float32 `json:"gemoglob,omitempty"`
}

func (updated *ErythrocyteUpdate) Update(erythrocyteI api.JSONtype) bool {
	var needUpdate bool

	switch erythrocyte := erythrocyteI.(type) {
	case *Erythrocyte:
		updateInt32(&erythrocyte.TextureID, updated.TextureID, &needUpdate)
		updateInt32(&erythrocyte.FormID, updated.FormID, &needUpdate)
		updateInt32(&erythrocyte.DiseaseID, updated.DiseaseID, &needUpdate)

		updateFloat32(&erythrocyte.SizeX, updated.SizeX, &needUpdate)
		updateFloat32(&erythrocyte.SizeY, updated.SizeY, &needUpdate)
		updateFloat32(&erythrocyte.SizeZ, updated.SizeZ, &needUpdate)

		updateFloat32(&erythrocyte.AngleX, updated.AngleX, &needUpdate)
		updateFloat32(&erythrocyte.AngleY, updated.AngleY, &needUpdate)
		updateFloat32(&erythrocyte.AngleZ, updated.AngleZ, &needUpdate)

		updateFloat32(&erythrocyte.ScaleX, updated.ScaleX, &needUpdate)
		updateFloat32(&erythrocyte.ScaleY, updated.ScaleY, &needUpdate)
		updateFloat32(&erythrocyte.ScaleZ, updated.ScaleZ, &needUpdate)

		updateFloat32(&erythrocyte.PositionX, updated.PositionX, &needUpdate)
		updateFloat32(&erythrocyte.PositionY, updated.PositionY, &needUpdate)
		updateFloat32(&erythrocyte.PositionZ, updated.PositionZ, &needUpdate)

		updateFloat32(&erythrocyte.Form, updated.Form, &needUpdate)
		updateFloat32(&erythrocyte.Oxygen, updated.Oxygen, &needUpdate)
		updateFloat32(&erythrocyte.Gemoglob, updated.Gemoglob, &needUpdate)
	}

	return needUpdate
}

func updateInt32(oldValue, newValue *int32, needUpdate *bool) {
	if newValue != nil {
		*needUpdate = true
		*oldValue = *newValue
	}
}

func updateFloat32(oldValue, newValue *float32, needUpdate *bool) {
	if newValue != nil {
		*needUpdate = true
		*oldValue = *newValue
	}
}

func updateBool(oldValue, newValue *bool, needUpdate *bool) {
	if newValue != nil {
		*needUpdate = true
		*oldValue = *newValue
	}
}

func updateString(oldValue, newValue *string, needUpdate *bool) {
	if newValue != nil {
		fmt.Println("newValue", *newValue)
		*needUpdate = true
		*oldValue = *newValue
	}
}

func updateTime(oldValue, newValue *time.Time, needUpdate *bool) {
	if newValue != nil {
		*needUpdate = true
		*oldValue = *newValue
	}
}

// 266 -> 238
