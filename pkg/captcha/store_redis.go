package captcha

import (
	"fmt"
	"go-oj/pkg/config"
	"go-oj/pkg/redis"
	"time"
)

type StoreRedis struct {
	Client    *redis.RedisClient
	KeyPrefix string
}

func (s *StoreRedis) Set(key, value string) error {
	expireTime := time.Duration(config.GetInt64("captcha.expire_time")) * time.Minute
	if ok := s.Client.Set(s.KeyPrefix+key, value, expireTime); !ok {
		return fmt.Errorf("存储captcha失败")
	}
	return nil
}

func (s *StoreRedis) Get(key string, clear bool) string {
	v := s.Client.Get(s.KeyPrefix + key)
	if clear {
		s.Client.Del(key)
	}
	return v
}

func (s *StoreRedis) Verify(id string, answer string, clear bool) bool {
	return s.Get(id, clear) == answer
}
