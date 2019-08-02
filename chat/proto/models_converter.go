package chat

import (

	//
	"time"

	_ "github.com/lib/pq"

	"github.com/golang/protobuf/ptypes"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

func UserFromNullUser(nullUser *models.MessageUserSQL) *User {
	if val, _ := nullUser.ID.Value(); val != nil {
		return &User{
			Id:     int32(nullUser.ID.Int64),
			Name:   nullUser.Name.String,
			Photo:  nullUser.Photo.String,
			Status: Status(nullUser.Status.Int64),
		}
	}
	return nil
}

func MessageFromNullMessage(nullMessage *models.MessageSQL) *Message {
	if val, _ := nullMessage.ID.Value(); val != nil {
		pMessage := &Message{
			Id:     int32(nullMessage.ID.Int64),
			Text:   nullMessage.Text.String,
			ChatId: int32(nullMessage.ChatID.Int64),
			From:   UserFromNullUser(nullMessage.From),
			To:     UserFromNullUser(nullMessage.To),
		}
		pMessage.Time, _ = ptypes.TimestampProto(nullMessage.Time)
		return pMessage
	}
	return nil
}

func UserFromProto(pUser *User) *models.UserPublicInfo {
	return &models.UserPublicInfo{
		ID:       pUser.Id,
		Name:     pUser.Name,
		PhotoURL: pUser.Photo,
	}
}

func UserToProto(user *models.UserPublicInfo) *User {
	return &User{
		Id:    user.ID,
		Name:  user.Name,
		Photo: user.PhotoURL,
	}
}

func MessagesFromProto(loc *time.Location, pMessages ...*Message) []*models.Message {
	var (
		mMessages = make([]*models.Message, 0)
		err       error
	)
	for _, message := range pMessages {

		mMessage := &models.Message{
			ID:     message.Id,
			User:   UserFromProto(message.From),
			Text:   message.Text,
			Status: int32(message.From.Status),
			Edited: message.Edited,
		}
		mMessage.Time, err = ptypes.Timestamp(message.Time)
		if err != nil {
			panic("MessagesFromProto panic")
		}

		mMessage.Time = mMessage.Time.In(loc)
		mMessages = append(mMessages, mMessage)
	}
	return mMessages
}

func MessagesToProto(mMessages ...*models.Message) *Messages {
	var (
		pMessages = make([]*Message, 0)
	)
	for _, message := range mMessages {
		pMessage := MessageToProto(message)
		pMessages = append(pMessages, pMessage)
	}
	return &Messages{
		Messages: pMessages,
	}
}

func MessageToProto(message *models.Message) *Message {
	var err error
	pMessage := &Message{
		Id:     message.ID,
		From:   UserToProto(message.User),
		Text:   message.Text,
		Edited: message.Edited,
	}
	pMessage.Time, err = ptypes.TimestampProto(message.Time)
	if err != nil {
		panic("MessagesToProto panic")
	}
	pMessage.From.Status = Status(message.Status)
	return pMessage
}
