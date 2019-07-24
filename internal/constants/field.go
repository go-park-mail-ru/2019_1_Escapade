package constants

import (
	"io/ioutil"
)

// FieldConfiguration - the limits of the characteristics of the field
//easyjson:json
type fieldConfiguration struct {
	Set       bool
	WidthMin  int `json:"widthMin"`
	WidthMax  int `json:"widthMax"`
	HeightMin int `json:"heightMin"`
	HeightMax int `json:"heightMax"`
}

// FIELD - singleton of field constants
var FIELD = fieldConfiguration{}

// InitField initializes FIELD
func InitField() error {
	var (
		data []byte
		err  error
	)

	path := "field.json"

	if data, err = ioutil.ReadFile(path); err != nil {
		return err
	}

	var tmp fieldConfiguration
	if err = tmp.UnmarshalJSON(data); err != nil {
		return err
	}
	tmp.Set = true
	FIELD = tmp

	return nil
}
