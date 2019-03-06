package models

type Game struct {
	FieldWidth  uint `json:"fieldWidth"`
	FieldHeight uint `json:"fieldHeight"`
	MinsTotal   uint `json:"minsTotal"`
	MinsFound   uint `json:"minsFound"`
	Finished    bool `json:"finihsed"`
	Exploded    bool `json:"exploded"`
}
