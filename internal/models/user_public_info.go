package models

import (
	"database/sql"
)

type UserPublicInfo struct {
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	BestScore sql.NullString `json:"bestScore"`
	BestTime  sql.NullString `json:"bestTime"`
}
