package api

import (

	//session "github.com/go-park-mail-ru/2019_1_Escapade/auth/proto"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	//re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	// "github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
	// "context"
	// "fmt"
	// "strconv"
	// "github.com/gomodule/redigo/redis"
)

// SessionManager session struct
type APIManager struct {
	config config.Configuration
	api    *Handler
}

/*
func NewAPIManager(redis redis.Conn, c config.SessionConfig) *APIManager {
	return &SessionManager{
		redisConn: redis,
		config:    c,
	}
}

func (sm *APIManager) Create(ctx context.Context, sess *session.Session) (sid *session.SessionID, err error) {
	fmt.Println("Creating sess for: ", sess.UserID)
	sid = &session.SessionID{ID: utils.RandomString(sm.config.Length)}
	// UserID - в конфиг, name - не хранить
	result, err := redis.String(sm.redisConn.Do("HMSET", sid.ID,
		"UserID", sess.UserID,
		"Name", sess.Login,
		"EX", sm.config.LifetimeSeconds))

	if err != nil {
		return &session.SessionID{ID: ""}, err
	}
	if result != "OK" {
		return &session.SessionID{ID: ""}, re.ErrorSessionQueryNotOK(result)
	}
	fmt.Println("OK")
	return
}

func (sm *APIManager) Delete(ctx context.Context, sess *session.SessionID) (i *session.Nothing, err error) {
	_, err = redis.Int(sm.redisConn.Do("DEL", sess.ID))
	if err != nil {
		fmt.Println("redis error:", err)
	}
	fmt.Println("Deleted session: ", sess.ID)
	i = &session.Nothing{}
	return
}

func (sm *APIManager) Check(ctx context.Context, cookie *session.SessionID) (sess *session.Session, err error) {
	var (
		userID, login string
		id            int
	)
	if userID, err = redis.String(sm.redisConn.Do("HGET", cookie.ID, "UserID")); err != nil {
		fmt.Println("cant get userID:", err)
		return &session.Session{UserID: -1}, err
	}

	if login, err = redis.String(sm.redisConn.Do("HGET", cookie.ID, "Name")); err != nil {
		fmt.Println("cant get login:", err)
		return &session.Session{UserID: -1}, err
	}

	if id, err = strconv.Atoi(userID); err != nil {
		fmt.Println("cant convert:", userID)
		return &session.Session{UserID: -1}, err
	}

	fmt.Println("Got session for: ", login, id)
	sess = &session.Session{
		UserID: int32(id),
		Login:  login,
	}
	return
}
*/
