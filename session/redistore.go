package session

import (
	"errors"
	"net/http"
	"sync"

	"github.com/boj/redistore"
	ginsessions "github.com/gin-gonic/contrib/sessions"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/sessions"

	"github.com/liumingmin/goutils/log4go"
)

type RedisStoreOptions struct {
	SessionName     string
	PoolSize        int
	Network         string
	Address         string
	Password        string
	Db              string
	SecureCookieKey string
	BizKey          string
}

type RedisStore struct {
	*redistore.RediStore

	sessionName string
	maxLength   int
	keyPrefix   string
	serializer  redistore.SessionSerializer
	bizIdMap    sync.Map
	bizKey      string
}

func NewRediStoreWithDB(opt RedisStoreOptions) (*RedisStore, error) {
	if opt.SessionName == "" {
		opt.SessionName = "websession"
	}

	if opt.Network == "" {
		opt.Network = "tcp"
	}

	if opt.Db == "" {
		opt.Db = "0"
	}

	if opt.SecureCookieKey == "" {
		opt.SecureCookieKey = "secret"
	}

	store, err := redistore.NewRediStoreWithDB(opt.PoolSize, opt.Network, opt.Address, opt.Password,
		opt.Db, []byte(opt.SecureCookieKey))
	if err != nil {
		return nil, err
	}

	rdsStore := &RedisStore{
		RediStore:   store,
		sessionName: opt.SessionName,
		serializer:  redistore.GobSerializer{},
		maxLength:   4096,
		keyPrefix:   "session_",
		bizKey:      opt.BizKey,
	}

	store.SetSerializer(rdsStore.serializer)
	store.SetMaxLength(rdsStore.maxLength)
	store.SetKeyPrefix(rdsStore.keyPrefix)

	return rdsStore, nil
}

func (c *RedisStore) Options(options ginsessions.Options) {
	c.RediStore.Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}

func (c *RedisStore) SessionName() string {
	return c.sessionName
}

func (s *RedisStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	session, err := s.RediStore.Get(r, name)
	if err == nil {
		if value, ok := session.Values[s.bizKey]; ok {
			if bizId, ok2 := value.(string); ok2 {
				s.bizIdMap.Store(bizId, session.ID)
			}
		}
	}
	return session, err
}

func (s *RedisStore) SetKV(bizId, key string, value interface{}) error {
	id, ok := s.bizIdMap.Load(bizId)
	if !ok {
		log4go.Error("sessionId is empty: %v", bizId)
		return nil
	}

	session, err := s.getSession(id.(string))
	if err != nil {
		return err
	}

	if session.IsNew {
		log4go.Info("SetKV no session found: %v", id)
		return nil
	}

	session.Values[key] = value
	return s.saveSession(session)
}

func (s *RedisStore) SetKVs(bizId string, kvs map[string]interface{}) error {
	id, ok := s.bizIdMap.Load(bizId)
	if !ok {
		log4go.Error("sessionId is empty: %v", bizId)
		return nil
	}

	session, err := s.getSession(id.(string))
	if err != nil {
		return err
	}

	if session.IsNew {
		log4go.Info("SetKV no session found: %v", id)
		return nil
	}

	for key, value := range kvs {
		session.Values[key] = value
	}

	return s.saveSession(session)
}

func (s *RedisStore) getSession(id string) (*sessions.Session, error) {
	var (
		err error
		ok  bool
	)
	session := sessions.NewSession(s, s.sessionName)
	session.ID = id
	// make a copy
	options := *s.RediStore.Options
	session.Options = &options
	session.IsNew = true

	ok, err = s.load(session)
	session.IsNew = !(err == nil && ok)

	return session, err
}

func (s *RedisStore) saveSession(session *sessions.Session) error {
	if err := s.save(session); err != nil {
		return err
	}

	return nil
}

func (s *RedisStore) save(session *sessions.Session) error {
	b, err := s.serializer.Serialize(session)
	if err != nil {
		return err
	}
	if s.maxLength != 0 && len(b) > s.maxLength {
		return errors.New("SessionStore: the value to store is too big")
	}
	conn := s.Pool.Get()
	defer conn.Close()
	if err = conn.Err(); err != nil {
		return err
	}
	age := session.Options.MaxAge
	if age == 0 {
		age = s.DefaultMaxAge
	}
	_, err = conn.Do("SETEX", s.keyPrefix+session.ID, age, b)
	return err
}

// load reads the session from redis.
// returns true if there is a sessoin data in DB
func (s *RedisStore) load(session *sessions.Session) (bool, error) {
	conn := s.Pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return false, err
	}
	data, err := conn.Do("GET", s.keyPrefix+session.ID)
	if err != nil {
		return false, err
	}
	if data == nil {
		return false, nil // no data was associated with this key
	}
	b, err := redis.Bytes(data, err)
	if err != nil {
		return false, err
	}
	return true, s.serializer.Deserialize(b, session)
}
