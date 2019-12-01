package database

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	//
	_ "github.com/lib/pq"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/proto"
)

// UserFromNullUser converts the structure to retrieve user from the database
// into the structure for transmission over grpc
func UserFromNullUser(nullUser *models.MessageUserSQL) *proto.User {
	if nullUser == nil {
		return nil
	}
	if val, _ := nullUser.ID.Value(); val != nil {
		return &proto.User{
			Id:     int32(nullUser.ID.Int64),
			Name:   nullUser.Name.String,
			Photo:  nullUser.Photo.String,
			Status: proto.Status(nullUser.Status.Int64),
		}
	}
	return nil
}

// MessageFromNullMessage converts the structure to retrieve message from the
// database into the structure for transmission over grpc
func MessageFromNullMessage(nullMessage *models.MessageSQL) (*proto.Message, error) {
	if nullMessage == nil {
		return nil, nil
	}
	if val, _ := nullMessage.ID.Value(); val != nil {
		var (
			pMessage = &proto.Message{
				Id:     int32(nullMessage.ID.Int64),
				Text:   nullMessage.Text.String,
				ChatId: int32(nullMessage.ChatID.Int64),
				From:   UserFromNullUser(nullMessage.From),
				To:     UserFromNullUser(nullMessage.To),
				Edited: nullMessage.Edited.Bool,
			}
			err error
		)
		if nullMessage.Answer != nil {
			pMessage.Answer, err = MessageFromNullMessage(nullMessage)
		}
		pMessage.Time, err = ptypes.TimestampProto(nullMessage.Time)
		return pMessage, err
	}
	return nil, nil
}

// UserFromProto converts the structure for transmission over grpc into the
// structure to retrieve user from the database
func UserFromProto(pUser *proto.User) *models.UserPublicInfo {
	if pUser == nil {
		return nil
	}
	utils.Debug(false, "pUser.Photo", pUser.Photo)
	return &models.UserPublicInfo{
		ID:      pUser.Id,
		Name:    pUser.Name,
		FileKey: pUser.Photo,
	}
}

// UserToProto converts the structure to retrieve user from the database
// into the structure for transmission over grpc
func UserToProto(user *models.UserPublicInfo) *proto.User {
	if user == nil {
		return nil
	}
	return &proto.User{
		Id:    user.ID,
		Name:  user.Name,
		Photo: user.FileKey,
	}
}

// MessagesFromProto converts the structure for transmission over grpc into the
// structure to retrieve messages from the database
func MessagesFromProto(loc *time.Location, pMessages ...*proto.Message) ([]*models.Message, error) {
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
			return mMessages, err
		}

		mMessage.Time = mMessage.Time.In(loc)
		mMessages = append(mMessages, mMessage)
	}
	return mMessages, err
}

// MessagesToProto converts the structure to retrieve messages from the database
// into the structure for transmission over grpc
func MessagesToProto(chatID int32, mMessages ...*models.Message) (*proto.Messages, error) {
	var (
		pMessages = make([]*proto.Message, 0)
	)
	for _, message := range mMessages {
		pMessage, err := MessageToProto(message, chatID)
		if err != nil {
			return &proto.Messages{}, err
		}
		pMessages = append(pMessages, pMessage)
	}
	return &proto.Messages{
		Messages: pMessages,
	}, nil
}

// MessageToProto converts the structure to retrieve message from the database
// into the structure for transmission over grpc
func MessageToProto(message *models.Message, chatID int32) (*proto.Message, error) {
	var (
		err      error
		pMessage = &proto.Message{
			Id:     message.ID,
			From:   UserToProto(message.User),
			Text:   message.Text,
			Edited: message.Edited,
			ChatId: chatID,
		}
	)
	pMessage.Time, err = ptypes.TimestampProto(message.Time)
	if err != nil {
		return pMessage, err
	}
	pMessage.From.Status = proto.Status(message.Status)
	return pMessage, err
}
