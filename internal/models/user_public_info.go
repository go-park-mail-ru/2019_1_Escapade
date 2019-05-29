package models

import (
	"database/sql"
	"encoding/json"
)

// UserPublicInfo information about person
// available for unauthorized users
type UserPublicInfo struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	PhotoURL  string         `json:"photo,omitempty"`
	FileKey   string         `json:"-"`
	BestScore sql.NullString `json:"bestScore"`
	BestTime  sql.NullString `json:"bestTime"`
	Difficult int            `json:"difficult"`
}

type UserPublicInfoSQL struct {
	ID        sql.NullInt64  `json:"id"`
	Name      sql.NullString `json:"name"`
	PhotoURL  sql.NullString `json:"photo,omitempty"`
	FileKey   sql.NullString `json:"-"`
	BestScore sql.NullString `json:"bestScore"`
	BestTime  sql.NullString `json:"bestTime"`
	Difficult sql.NullInt64  `json:"difficult"`
}

// Compare compare two users
func (a UserPublicInfo) Compare(b UserPublicInfo) bool {
	return a.ID == b.ID && a.Name == b.Name
}

// ComparePublicUsers unmarshal strings to struct UserPublicInfo and compare them
func ComparePublicUsers(a, b string) bool {
	var (
		err   error
		userA UserPublicInfo
		userB UserPublicInfo
	)
	if err = json.Unmarshal([]byte(a), &userA); err != nil {
		return false
	}
	if err = json.Unmarshal([]byte(b), &userB); err != nil {
		return false
	}

	return userA.Compare(userB)
}
