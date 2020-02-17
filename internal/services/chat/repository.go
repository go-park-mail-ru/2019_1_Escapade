package chat

import (
	"context"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/models"
)

type ChatRepository interface {
	Create(
		ctx context.Context,
		chatType,
		typeID int32,
	) (int32, error)

	Get(
		ctx context.Context,
		chatModel *models.Chat,
	) (int32, error)
}

type UserRepository interface {
	Create(
		ctx context.Context,
		chatID int32,
		users ...*models.User,
	) error

	Delete(
		ctx context.Context,
		userInGroup *models.UserInGroup,
	) error
}

type MessageRepository interface {
	CreateOne(
		ctx context.Context,
		message *models.Message,
	) (int32, error)

	CreateMany(
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
