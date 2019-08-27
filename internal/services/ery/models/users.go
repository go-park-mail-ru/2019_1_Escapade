package models

import (
	"time"

	"database/sql"
)

// User структура пользователя
//easyjson:json
type User struct {
	ID         int32          `json:"id" db:"id"`
	Name       string         `json:"name" db:"name" maxLength:"30" example:"John"`
	Password   string         `json:"password" db:"password" minLength:"6" maxLength:"30" example:"easyPassword" `
	PhotoTitle string         `json:"photo_title" db:"photo_title" maxLength:"40" example:"image12.jpg" `
	Website    sql.NullString `json:"website" db:"website" maxLength:"80" example:"https://github.com/SmartPhoneJava" `
	About      sql.NullString `json:"about" db:"about" maxLength:"400" example:"Student of BMSTU" `
	Email      sql.NullString `json:"email" db:"email" maxLength:"50" example:"artyom2019@gmail.com" `
	Phone      sql.NullString `json:"phone" db:"phone" minLength:"8" maxLength:"18" example:"81234567" `
	Birthday   time.Time      `json:"birthday" db:"birthday" example:"1992-09-25 00:00:00" `
	Add        time.Time      `json:"add" db:"add" example:"2006-01-02 15:04:05" `
	LastSeen   time.Time      `json:"last_seen" db:"last_seen" example:"2006-01-02 15:04:05" `
}

// Users структура массива пользователей
//easyjson:json
type Users struct {
	Users []User `json:"users"`
}

// UpdatePrivateUser структура обновления имени или пароля пользователя. В
// полях Old и New должны быть заполнены поля имя, пароль для подтвеждения
// личности и для указания нового имени/пароля соответственно
//easyjson:json
type UpdatePrivateUser struct {
	Old User `json:"old"`
	New User `json:"new"`
}

// GetName получить имя
// функция интерфейска UserI(пакет api)
func (user *User) GetName() string {
	return user.Name
}

// GetPassword получить пароль
// функция интерфейска UserI(пакет api)
func (user *User) GetPassword() string {
	return user.Password
}

// SetName установить имя
// функция интерфейска UserI(пакет api)
func (user *User) SetName(name string) {
	user.Name = name
}

// SetPassword уставновить пароль
// функция интерфейска UserI(пакет api)
func (user *User) SetPassword(password string) {
	user.Password = password
}
