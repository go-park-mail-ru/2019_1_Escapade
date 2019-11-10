package chat

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
)

type UserUseCase struct {
	database.UseCaseBase
	user UserRepositoryI
}

func (db *UserUseCase) Init(user UserRepositoryI) {
	db.user = user
}

// InviteToChat invite user to the chat
// to work correctly, specify user and id of the chat
func (db *UserUseCase) InviteToChat(userInChat *UserInGroup) (*Result, error) {
	if err := db.check(userInChat); err != nil {
		return &Result{}, err
	}

	var (
		tx  database.TransactionI
		err error
	)

	if tx, err = db.Db.Begin(); err != nil {
		return &Result{Done: false}, err
	}
	defer tx.Rollback()

	if err = db.user.create(tx, userInChat.Chat.Id, userInChat.User); err != nil {
		return &Result{Done: false}, err
	}

	err = tx.Commit()

	return &Result{Done: true}, err
}

// LeaveChat leave user from the chat
// to work correctly, specify user and id of the chat
func (db *UserUseCase) LeaveChat(userInChat *UserInGroup) (*Result, error) {

	if err := db.check(userInChat); err != nil {
		return &Result{}, err
	}

	res, err := db.user.delete(db.Db, userInChat)
	return res, err
}

func (db *UserUseCase) check(userInChat *UserInGroup) error {
	if userInChat == nil {
		return re.InvalidUser()
	}

	if userInChat.User == nil || userInChat.User.Id <= 0 {
		return re.InvalidUser()
	}

	if userInChat.Chat == nil || userInChat.Chat.Id <= 0 {
		return re.InvalidChatID()
	}
	return nil
}
