package models

import "database/sql"

type UserPublicInfo struct {
	ID        int            `json:"-"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Photo     []byte         `json:"photo"`
	FileName  string         `json:"-"`
	BestScore sql.NullString `json:"bestScore"`
	BestTime  sql.NullString `json:"bestTime"`
	Difficult int            `json:"difficult"`
}
