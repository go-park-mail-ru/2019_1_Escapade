package models

import (
	"database/sql"
)

// UserPublicInfo information about person
// available for unauthorized users
//easyjson:json
type UserPublicInfo struct {
	ID        int            `json:"id"`
	Name      string         `json:"name" minLength:"3" maxLength:"30"`
	PhotoURL  string         `json:"photo,omitempty"  maxLength:"50"`
	FileKey   string         `json:"-"`
	BestScore sql.NullString `json:"bestScore"`
	BestTime  sql.NullString `json:"bestTime"`
	Difficult int            `json:"difficult"`
}

// UsersPublicInfo is the slice of UserPublicInfo
//easyjson:json
type UsersPublicInfo struct {
	Users []*UserPublicInfo `json:"users"`
}

// UserPublicInfoSQL wrapper of UserPublicInfo
// Required to obtain data from the database
type UserPublicInfoSQL struct {
	ID        sql.NullInt64
	Name      sql.NullString
	PhotoURL  sql.NullString
	FileKey   sql.NullString
	BestScore sql.NullString
	BestTime  sql.NullString
	Difficult sql.NullInt64
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
	if err = userA.UnmarshalJSON([]byte(a)); err != nil {
		return false
	}
	if err = userB.UnmarshalJSON([]byte(b)); err != nil {
		return false
	}

	return userA.Compare(userB)
}
