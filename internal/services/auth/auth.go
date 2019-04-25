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
	fmt.Println("Creating sess for: ", sess.UserID)
	sid = &session.SessionID{ID: utils.RandomString(10)}
	result, err := redis.String(sm.redisConn.Do("SET", sid.ID, sess.UserID, "EX", 86400))
	if err != nil {
		return &session.SessionID{ID: ""}, err
	}
	if result != "OK" {
		return &session.SessionID{ID: ""}, fmt.Errorf("result not OK")
	}
	fmt.Println("OK")
	return
}

func (sm *SessionManager) Delete(ctx context.Context, cookie *session.SessionID) (i *session.Nothing, err error) {
	_, err = redis.Int(sm.redisConn.Do("DEL", cookie.ID))
	if err != nil {
		log.Println("redis error:", err)
	}
	fmt.Println("Deleted session: ", cookie.ID)
	i = &session.Nothing{}
	return
}

func (sm *SessionManager) Check(ctx context.Context, cookie *session.SessionID) (sess *session.Session, err error) {
	userID, err := redis.Int(sm.redisConn.Do("GET", cookie.ID))
	if err != nil {
		log.Println("cant get data:", err)
		return &session.Session{UserID: -1}, err
	}
	log.Println("Got session for: ", userID)
	sess = &session.Session{UserID: int32(userID)}
	return
}
