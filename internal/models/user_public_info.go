package models

import (
	"database/sql"
)

type UserPublicInfo struct {
	Name      string         `json:"username"`
	Email     string         `json:"email"`
	BestScore sql.NullString `json:"bestScore"`
	BestTime  sql.NullString `json:"bestTime"`
}
