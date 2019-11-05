package chat

import (
	context "context"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"

	//
	_ "github.com/lib/pq"
)

// repositories stores all implementations of operations in the database
type repositories struct {
	user    UserRepositoryI
	chat    ChatRepositoryI
	message MessageRepositoryI
}

type Handler struct {
	user    UserUseCaseI
	chat    ChatUseCaseI
	message MessageUseCaseI
}

// InitWithPostgreSQL apply postgreSQL as database
func (h *Handler) InitWithPostgreSQL(c *config.Configuration) error {
	var (
		reps = repositories{
			user:    &UserRepositoryPQ{},
			message: &MessageRepositoryPQ{},
			chat:    &ChatRepositoryPQ{},
		}
		database = &database.PostgresSQL{}
	)
	return h.Init(c, database, reps)
}

func (h *Handler) Init(c *config.Configuration, db database.DatabaseI, reps repositories) error {
	err := db.Open(c.DataBase)
	if err != nil {
		return err
	}

	var user = &UserUseCase{}
	user.Init(reps.user)
	h.user = user
	err = h.user.Use(db)
	if err != nil {
		return err
	}

	var message = &MessageUseCase{}
	message.Init(reps.message)
	h.message = message
	err = h.message.Use(db)
	if err != nil {
		return err
	}

	var chat = &ChatUseCase{}
	chat.Init(reps.chat, reps.user)
	h.chat = chat
	err = h.chat.Use(db)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) Check() (bool, error) {
	err := h.user.Get().Ping()
	if err != nil {
		return false, err
	}

	return false, nil
}

func (h *Handler) Close() {
	h.user.Close()
	h.message.Close()
	h.chat.Close()
	return
}

func (h *Handler) CreateChat(ctx context.Context, in *ChatWithUsers) (*ChatID, error) {
	return h.chat.Create(ctx, in)
}
func (h *Handler) GetChat(ctx context.Context, in *Chat) (*ChatID, error) {
	return h.chat.GetOne(ctx, in)
}
func (h *Handler) InviteToChat(ctx context.Context, in *UserInGroup) (*Result, error) {
	return h.user.InviteToChat(ctx, in)
}
func (h *Handler) LeaveChat(ctx context.Context, in *UserInGroup) (*Result, error) {
	return h.user.LeaveChat(ctx, in)
}
func (h *Handler) AppendMessage(ctx context.Context, in *Message) (*MessageID, error) {
	return h.message.AppendOne(ctx, in)
}
func (h *Handler) AppendMessages(ctx context.Context, in *Messages) (*MessagesID, error) {
	return h.message.AppendMany(ctx, in)
}
func (h *Handler) UpdateMessage(ctx context.Context, in *Message) (*Result, error) {
	return h.message.Update(ctx, in)
}
func (h *Handler) DeleteMessage(ctx context.Context, in *Message) (*Result, error) {
	return h.message.Delete(ctx, in)
}
func (h *Handler) GetChatMessages(ctx context.Context, in *ChatID) (*Messages, error) {
	return h.message.GetAll(ctx, in)
}
