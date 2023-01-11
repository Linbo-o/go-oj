package verifycode

import (
	"go-oj/pkg/config"
	"go-oj/pkg/redis"
	"time"
)

type StoreRedis struct {
	Client    *redis.RedisClient
	KeyPrefix string
}

func (s *StoreRedis) Set(id string, value string) bool {
	expireTime := time.Duration(int64(config.GetInt64("verifycode.expire_time"))) * time.Minute
	return s.Client.Set(s.KeyPrefix+id, value, expireTime)
}

func (s *StoreRedis) Get(id string, clear bool) string {
	id = s.KeyPrefix + id
	v := s.Client.Get(id)
	if clear {
		s.Client.Del(id)
	}
	return v
}

func (s *StoreRedis) Verify(id, answer string, clear bool) bool {
	return s.Get(id, clear) == answer
}
