package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/models"

	_ "github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
)

// createPlayer создать пользователя
func (db *DB) createUser(tx *sqlx.Tx, user *models.User) error {
	sqlInsert := `
	INSERT INTO Users(name, password, photo_title, website,
		 about, email, phone) VALUES
		(:name, :password, :photo_title, :website, :about,
			:email, :phone)
		RETURNING *;
		`
	return createAndReturnStruct(tx, sqlInsert, &user)
}

// updateUser обновить публичную информацию о пользователе
func (db *DB) updateUser(tx *sqlx.Tx, user *models.User) error {
	statement := `
	UPDATE Users 
	SET photo_title = :photo_title , website = :website, about = :about, email = :email,
	phone = :phone, birthday = :birthday
	WHERE id = :id
		`
	_, err := tx.NamedExec(statement, user)
	return err
}

// updateUser обновить публичную информацию о пользователе
func (db *DB) SetNewImage(filekey string, userID int32) error {
	statement := `
	UPDATE Users 
	SET photo_title = $1
	WHERE id = $2
		`
	_, err := db.db.Exec(statement, filekey, userID)
	return err
}

// updateUserPassword обновить имя или пароль пользователя
func (db *DB) updateUserPrivate(tx *sqlx.Tx, oldUser *models.User, newUser *models.User) error {
	statement := `
	UPDATE Users 
	SET name = $1, password = $2
	WHERE name like $3 and password like $4
		`
	_, err := tx.Exec(statement, newUser.Name, newUser.Password, oldUser.Name, oldUser.Password)
	return err
}

func (db *DB) getOneUser(tx *sqlx.Tx, userID int32) (models.User, error) {
	statement := `
	select * from Users 
	WHERE id=$1
		`

	var (
		row  = tx.QueryRowx(statement, userID)
		user models.User
	)

	err := row.StructScan(&user)
	user.Password = ""
	return user, err
}

// searchUsersWithName - найти людей с именем name
// Поиск людей в одном проекте реализовать на стороне клиента при получении массива участников
func (db *DB) searchUsersWithName(tx *sqlx.Tx, name string) ([]models.User, error) {
	statement := `
	select * from Users where POSITION (lower($1) IN lower(name)) > 0;
	`
	//name ~* $1
	rows, err := tx.Queryx(statement, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)

	for rows.Next() {
		var user models.User
		if err = rows.StructScan(&user); err != nil {
			break
		}
		user.Password = ""
		users = append(users, user)
	}
	if err != nil {
		return users, err
	}
	return users, err
}

func (db *DB) getAllUsers(tx *sqlx.Tx) ([]models.User, error) {
	statement := `
	select  * from Users
	`
	rows, err := tx.Queryx(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)

	for rows.Next() {
		var user models.User
		if err = rows.StructScan(&user); err != nil {
			break
		}
		user.Password = ""
		users = append(users, user)
	}
	if err != nil {
		return users, err
	}
	return users, err
}
