package chat

import (
	"context"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/models"
)

const (
	ErrNoChatUseCase    = "No chat usecase given"
	ErrNoUserUseCase    = "No user usecase given"
	ErrNoMessageUseCase = "No message usecase given"
)

type ChatUseCase interface {
	Create(
		ctx context.Context,
		chat *models.ChatWithUsers,
	) (int32, error)

	GetOne(
		ctx context.Context,
		chat *models.Chat,
	) (int32, error)
}

type UserUseCase interface {
	InviteToChat(
		ctx context.Context,
		userInChat *models.UserInGroup,
	) error

	LeaveChat(
		ctx context.Context,
		userInChat *models.UserInGroup,
	) error
}

type MessageUseCase interface {
	AppendOne(
		ctx context.Context,
		message *models.Message,
	) (int32, error)

	AppendMany(
		ctx context.Context,
		messages *models.Messages,
	) ([]int32, error)

	Update(
		ctx context.Context,
		message *models.Message,
	) error

	Delete(
		ctx context.Context,
		message *models.Message,
	) error

	GetAll(
		ctx context.Context,
		chatID int32,
	) ([]*models.Message, error)
}
