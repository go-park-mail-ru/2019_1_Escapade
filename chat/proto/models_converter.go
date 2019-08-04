package chat

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"time"

	//
	_ "github.com/lib/pq"

	"github.com/golang/protobuf/ptypes"
)

// UserFromNullUser converts the structure to retrieve user from the database
// into the structure for transmission over grpc
func UserFromNullUser(nullUser *models.MessageUserSQL) *User {
	if nullUser == nil {
		return nil
	}
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

// MessageFromNullMessage converts the structure to retrieve message from the
// database into the structure for transmission over grpc
func MessageFromNullMessage(nullMessage *models.MessageSQL) (*Message, error) {
	if nullMessage == nil {
		return nil, nil
	}
	if val, _ := nullMessage.ID.Value(); val != nil {
		var (
			pMessage = &Message{
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
func UserFromProto(pUser *User) *models.UserPublicInfo {
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
func UserToProto(user *models.UserPublicInfo) *User {
	if user == nil {
		return nil
	}
	return &User{
		Id:    user.ID,
		Name:  user.Name,
		Photo: user.FileKey,
	}
}

// MessagesFromProto converts the structure for transmission over grpc into the
// structure to retrieve messages from the database
func MessagesFromProto(loc *time.Location, pMessages ...*Message) ([]*models.Message, error) {
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
func MessagesToProto(mMessages ...*models.Message) (*Messages, error) {
	var (
		pMessages = make([]*Message, 0)
	)
	for _, message := range mMessages {
		pMessage, err := MessageToProto(message)
		if err != nil {
			return &Messages{}, err
		}
		pMessages = append(pMessages, pMessage)
	}
	return &Messages{
		Messages: pMessages,
	}, nil
}

// MessageToProto converts the structure to retrieve message from the database
// into the structure for transmission over grpc
func MessageToProto(message *models.Message) (*Message, error) {
	var (
		err      error
		pMessage = &Message{
			Id:     message.ID,
			From:   UserToProto(message.User),
			Text:   message.Text,
			Edited: message.Edited,
		}
	)
	pMessage.Time, err = ptypes.TimestampProto(message.Time)
	if err != nil {
		return pMessage, err
	}
	pMessage.From.Status = Status(message.Status)
	return pMessage, err
}
