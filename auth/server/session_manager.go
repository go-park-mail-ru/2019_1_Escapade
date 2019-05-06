package session

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"context"
	"fmt"
	"strconv"

	"github.com/gomodule/redigo/redis"
)

// SessionManager session struct
type SessionManager struct {
	redisConn redis.Conn
	config    config.SessionConfig
}

func NewSessionManager(redis redis.Conn, c config.SessionConfig) *SessionManager {
	return &SessionManager{
		redisConn: redis,
		config:    c,
	}
}

func (sm *SessionManager) Create(ctx context.Context, sess *Session) (sid *SessionID, err error) {
	fmt.Println("Creating sess for: ", sess.UserID)
	sid = &SessionID{ID: utils.RandomString(sm.config.Length)}
	// UserID - в конфиг, name - не хранить
	result, err := redis.String(sm.redisConn.Do("HMSET", sid.ID,
		"UserID", sess.UserID,
		"Name", sess.Login,
		"EX", sm.config.LifetimeSeconds))

	if err != nil {
		return &SessionID{ID: ""}, err
	}
	if result != "OK" {
		fmt.Println("NOT OK")
		return &SessionID{ID: ""}, re.ErrorSessionQueryNotOK(result)
	}
	fmt.Println("ALL OK", sid.ID, "!")

	return
}

func (sm *SessionManager) Delete(ctx context.Context, sess *SessionID) (i *Nothing, err error) {
	_, err = redis.Int(sm.redisConn.Do("DEL", sess.ID))
	if err != nil {
		fmt.Println("redis error:", err)
	}
	fmt.Println("Deleted session: ", sess.ID)
	i = &Nothing{}
	return
}

func (sm *SessionManager) Check(ctx context.Context, cookie *SessionID) (sess *Session, err error) {
	var (
		userID, login string
		id            int
	)
	if userID, err = redis.String(sm.redisConn.Do("HGET", cookie.ID, "UserID")); err != nil {
		fmt.Println("cant get userID:", err)
		return &Session{UserID: -1}, err
	}

	if login, err = redis.String(sm.redisConn.Do("HGET", cookie.ID, "Name")); err != nil {
		fmt.Println("cant get login:", err)
		return &Session{UserID: -1}, err
	}

	if id, err = strconv.Atoi(userID); err != nil {
		fmt.Println("cant convert:", userID)
		return &Session{UserID: -1}, err
	}

	fmt.Println("Got session for: ", login, id)
	sess = &Session{
		UserID: int32(id),
		Login:  login,
	}
	return
}
