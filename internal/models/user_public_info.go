package models

import (
	"database/sql"
)

type UserPublicInfo struct {
	Name      string         `json:"name"`
	Photo     sql.NullString `json:"photo"`
	BestScore sql.NullString `json:"bestScore"`
	BestTime  sql.NullString `json:"bestTime"`
}
