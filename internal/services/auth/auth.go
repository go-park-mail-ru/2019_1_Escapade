package session

import (
	session "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/proto"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/gomodule/redigo/redis"
)

// SessionManager session struct
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
	result, err := redis.String(sm.redisConn.Do("HMSET", sid.ID, "id", sess.UserID, "login", sess.Login, "EX", 86400))
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
	userID, err := redis.String(sm.redisConn.Do("HGET", cookie.ID, "id"))
	if err != nil {
		log.Println("cant get userID:", err)
		return &session.Session{UserID: -1}, err
	}
	login, err := redis.String(sm.redisConn.Do("HGET", cookie.ID, "login"))
	id, _ := strconv.Atoi(userID)
	if err != nil {
		log.Println("cant get login:", err)
		return &session.Session{UserID: -1}, err
	}
	log.Println("Got session for: ", login, id)
	sess = &session.Session{UserID: int32(id), Login: login}
	return
}
