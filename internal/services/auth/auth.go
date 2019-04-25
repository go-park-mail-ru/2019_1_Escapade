package session

import (
	"context"
	session "escapade/internal/services/auth/proto"
	"escapade/internal/utils"
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
)

type SessionManager struct {
	redisConn redis.Conn
}

func NewSessionManager(conn redis.Conn) *SessionManager {
	return &SessionManager{
		redisConn: conn,
	}
}

func (sm *SessionManager) Create(ctx context.Context, sess *session.Session) (sid *session.SessionID, err error) {
	fmt.Println("Creating sess for: ", sess.Login)
	sid = &session.SessionID{ID: utils.RandomString(10)}
	result, err := redis.String(sm.redisConn.Do("SET", sid.ID, sess.Login, "EX", 86400))
	if err != nil {
		return nil, err
	}
	if result != "OK" {
		return nil, fmt.Errorf("result not OK")
	}
	fmt.Println("OK")
	return
}

func (sm *SessionManager) Delete(ctx context.Context, cookie *session.SessionID) (i *session.Nothing, r error) {
	return
}

func (sm *SessionManager) Check(ctx context.Context, cookie *session.SessionID) (sess *session.Session, err error) {
	login, err := redis.String(sm.redisConn.Do("GET", cookie.ID))
	if err != nil {
		log.Println("cant get data:", err)
		return nil, nil
	}
	sess = &session.Session{Login: login}
	return
}
