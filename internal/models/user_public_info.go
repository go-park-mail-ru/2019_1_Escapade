package models

import (
	"database/sql"
	"encoding/json"
)

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

// Compare comapre two users
func (a UserPublicInfo) Compare(b UserPublicInfo) bool {
	return a.ID == b.ID && a.Name == b.Name && a.Email == b.Email
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
