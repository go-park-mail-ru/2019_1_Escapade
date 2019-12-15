package constants

// FieldConfiguration - the limits of the characteristics of the field
//easyjson:json
type FieldConfiguration struct {
	Set       bool
	WidthMin  int32 `json:"widthMin"`
	WidthMax  int32 `json:"widthMax"`
	HeightMin int32 `json:"heightMin"`
	HeightMax int32 `json:"heightMax"`
}

// FIELD - singleton of field constants
var FIELD = FieldConfiguration{}

// InitField initializes FIELD
func InitField(rep RepositoryI, path string) error {
	field, err := rep.GetField(path)
	if err != nil {
		return err
	}
	FIELD = field
	FIELD.Set = true

	return nil
}
