package models

import "database/sql"

// UserPublicInfo information about person
// available for unauthorized users
type UserPublicInfo struct {
	ID        int            `json:"-"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	PhotoURL  string         `json:"photo,omitempty"`
	FileKey   string         `json:"-"`
	BestScore sql.NullString `json:"bestScore"`
	BestTime  sql.NullString `json:"bestTime"`
	Difficult int            `json:"difficult"`
}
