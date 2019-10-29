package game

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// FieldJSON is a wrapper for sending Field by JSON
//easyjson:json
type FieldJSON struct {
	History   []*Cell `json:"history"`
	CellsLeft int32   `json:"cellsLeft"`

	Width     int32   `json:"width"`
	Height    int32   `json:"height"`
	Mines     int32   `json:"mines"`
	Difficult float64 `json:"difficult"`
}

// JSON convert Field to FieldJSON
func (field *Field) JSON() FieldJSON {
	utils.Debug(false, "FieldJSON")
	return FieldJSON{
		History:   field.History(),
		CellsLeft: field.cellsLeft(),
		Width:     field.Width,
		Height:    field.Height,
		Mines:     field.Mines,
		Difficult: field.Difficult,
	}
}

// MarshalJSON - overriding the standard method json.Marshal
func (field *Field) MarshalJSON() ([]byte, error) {
	return field.JSON().MarshalJSON()
}

// UnmarshalJSON - overriding the standard method json.Unmarshal
func (field *Field) UnmarshalJSON(b []byte) error {
	temp := &FieldJSON{}

	if err := temp.UnmarshalJSON(b); err != nil {
		return err
	}
	field.setHistory(temp.History)
	field.setCellsLeft(temp.CellsLeft)

	return nil
}
