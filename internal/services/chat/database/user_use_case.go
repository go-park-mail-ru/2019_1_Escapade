package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/proto"
)

type UserUseCase struct {
	database.UseCaseBase
	user UserRepositoryI
}

func (db *UserUseCase) Init(user UserRepositoryI) UserUseCaseI {
	db.user = user
	return db
}

// InviteToChat invite user to the chat
// to work correctly, specify user and id of the chat
func (db *UserUseCase) InviteToChat(userInChat *proto.UserInGroup) (*proto.Result, error) {
	if err := db.check(userInChat); err != nil {
		return &proto.Result{}, err
	}

	var (
		tx  database.TransactionI
		err error
	)

	if tx, err = db.Db.Begin(); err != nil {
		return &proto.Result{Done: false}, err
	}
	defer tx.Rollback()

	if err = db.user.create(tx, userInChat.Chat.Id, userInChat.User); err != nil {
		return &proto.Result{Done: false}, err
	}

	err = tx.Commit()

	return &proto.Result{Done: true}, err
}

// LeaveChat leave user from the chat
// to work correctly, specify user and id of the chat
func (db *UserUseCase) LeaveChat(userInChat *proto.UserInGroup) (*proto.Result, error) {

	if err := db.check(userInChat); err != nil {
		return &proto.Result{}, err
	}

	res, err := db.user.delete(db.Db, userInChat)
	return res, err
}

func (db *UserUseCase) check(userInChat *proto.UserInGroup) error {
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
