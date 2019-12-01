package constants

import (
	"io/ioutil"
)

// FieldConfiguration - the limits of the characteristics of the field
//easyjson:json
type fieldConfiguration struct {
	Set       bool
	WidthMin  int32 `json:"widthMin"`
	WidthMax  int32 `json:"widthMax"`
	HeightMin int32 `json:"heightMin"`
	HeightMax int32 `json:"heightMax"`
}

// FIELD - singleton of field constants
var FIELD = fieldConfiguration{}

// InitField initializes FIELD
func InitField(path string) error {
	var (
		data []byte
		err  error
	)

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
