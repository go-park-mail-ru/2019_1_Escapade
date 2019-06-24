package api

import (
	"context"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	ran "math/rand"
)

func (h *Handler) register(ctx context.Context,
	user models.UserPrivateInfo) (userID int, sessionID string, err error) {

	sessionID = utils.RandomString(16)
	if userID, err = h.DB.Register(&user, sessionID); err != nil {
		return
	}

	/*sessID, err := h.Clients.Session.Create(ctx,
		&session.Session{
			UserID: int32(userID),
			Login:  user.Name,
		})

	if err != nil {
		return
	}

	sessionID = sessID.ID*/
	return
}

func (h *Handler) deleteAccount(ctx context.Context,
	user *models.UserPrivateInfo, sessionID string) (err error) {

	if err = h.DB.DeleteAccount(user); err != nil {
		return
	}

	// _, err = h.Clients.Session.Delete(ctx,
	// 	&session.SessionID{
	// 		ID: sessionID,
	// 	})

	return
}

// RandomUsers create n random users
func (h *Handler) RandomUsers(limit int) {

	n := 16
	for i := 0; i < limit; i++ {
		ran.Seed(time.Now().UnixNano())
		user := models.UserPrivateInfo{
			Name:     utils.RandomString(n),
			Password: utils.RandomString(n)}
		userID, _, err := h.register(context.Background(), user)
		if err != nil {
			fmt.Println("some error is here", err)
			return
		}

		for j := 0; j < 4; j++ {
			record := &models.Record{
				Score:       ran.Intn(1000000),
				Time:        float64(ran.Intn(10000)),
				Difficult:   j,
				SingleTotal: ran.Intn(2),
				OnlineTotal: ran.Intn(2),
				SingleWin:   ran.Intn(2),
				OnlineWin:   ran.Intn(2)}
			h.DB.UpdateRecords(userID, record)
		}

	}
}
