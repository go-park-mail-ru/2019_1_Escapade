package message

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	chat "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/proto"
)

// MessagesI control access to messages
// Proxy Pattern
type MessagesI interface {
	synced.PublisherI

	Fix(message *models.Message, user *models.UserPublicInfo)
	Proto(message *models.Message) (*chat.Message, error)

	Write(message *models.Message, send *chat.Message) error
	Update(message *models.Message, send *chat.Message) error
	Delete(message *models.Message, send *chat.Message) error

	Send(message *models.Message)
	HandleError(message *models.Message, send *chat.Message)

	setMessages(messages []*models.Message)
	Messages() []*models.Message

	ChatID() int32
	setChatID(id int32)
}
