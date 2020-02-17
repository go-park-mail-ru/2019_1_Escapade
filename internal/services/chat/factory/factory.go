package factory

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/delivery/grpc"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/usecase/database"
)

func NewService(
	db infrastructure.Database,
	logger infrastructure.Logger,
	trace infrastructure.ErrorTrace,
	photo infrastructure.PhotoService,
	timeout time.Duration,
) (*grpc.ChatServiceServer, error) {
	chatUC, err := database.NewChat(db, logger, trace, timeout)
	if err != nil {
		return nil, err
	}
	userUC, err := database.NewUser(db, logger, trace, timeout)
	if err != nil {
		return nil, err
	}
	messageUC, err := database.NewMessage(
		db,
		logger,
		trace,
		photo,
		timeout,
	)
	if err != nil {
		return nil, err
	}
	return grpc.NewChatServiceServer(
		chatUC,
		userUC,
		messageUC,
	)
}
