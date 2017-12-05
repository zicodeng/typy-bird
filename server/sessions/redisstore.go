package sessions

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/go-redis/redis"
)

//RedisStore represents a session.Store backed by redis.
type RedisStore struct {
	//Redis client used to talk to redis server.
	Client *redis.Client
	//Used for key expiry time on redis.
	SessionDuration time.Duration
}

//NewRedisStore constructs a new RedisStore
func NewRedisStore(client *redis.Client, sessionDuration time.Duration) *RedisStore {
	//initialize and return a new RedisStore struct
	return &RedisStore{
		Client:          client,
		SessionDuration: sessionDuration,
	}
}

//Store implementation

//Save saves the provided `sessionState` and associated SessionID to the store.
//The `sessionState` parameter is typically a pointer to a struct containing
//all the data you want to associated with the given SessionID.
func (rs *RedisStore) Save(sid SessionID, sessionState interface{}) error {
	key := sid.getRedisKey()
	session, err := json.Marshal(sessionState)
	if err != nil {
		return fmt.Errorf("could not marshal to json: %v", err)
	}
	rs.Client.Set(string(key), session, cache.DefaultExpiration)
	return nil
}

//Get populates `sessionState` with the data previously saved
//for the given SessionID
func (rs *RedisStore) Get(sid SessionID, sessionState interface{}) error {
	key := sid.getRedisKey()
	sessionData := rs.Client.Get(string(key))
	sessionDataBytes, err := sessionData.Bytes()
	if err != nil {
		return ErrStateNotFound
	}
	err = json.Unmarshal(sessionDataBytes, sessionState)
	if err != nil {
		return fmt.Errorf("could not unmarshal data: %v", err)
	}
	rs.Client.Expire(string(key), 0)

	return nil
}

//GetAll retrieves all of the currently running sessions
func (rs *RedisStore) GetAll(sessionState interface{}) ([]string, error) {
	var cursor uint64
	keys, cursor, err := rs.Client.Scan(cursor, "", 10).Result()
	if err != nil {
		return nil, fmt.Errorf("could not get session ids: %v", err)
	}
	return keys, nil
}

//Delete deletes all state data associated with the SessionID from the store.
func (rs *RedisStore) Delete(sid SessionID) error {
	key := sid.getRedisKey()
	rs.Client.Del(string(key))
	return nil
}

//getRedisKey() returns the redis key to use for the SessionID
func (sid SessionID) getRedisKey() string {
	//convert the SessionID to a string and add the prefix "sid:" to keep
	//SessionID keys separate from other keys that might end up in this
	//redis instance
	return "sid:" + sid.String()
}
