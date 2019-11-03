package redis_utl

import (
	"MyProject/Short_Url/pkg/setting"
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"
)

var RedisCon *redis.Pool

func Start() error {
	RedisCon = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,
		MaxActive:   setting.RedisSetting.MaxActive,
		IdleTimeout: setting.RedisSetting.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				return nil, err
			}

			if setting.RedisSetting.Password != "" {
				if _, err := conn.Do("AUTH", setting.RedisSetting.Password); err != nil {
					conn.Close()
					return nil, err
				}
			}
			return conn, err
		},

		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}

	return nil
}

func Set(key string, data interface{}, time int) error {
	conn := RedisCon.Get()
	defer conn.Close()

	// value, err :=

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, value)
	if err != nil {
		return err
	}

	return nil
}

func Get(key string) ([]byte, error) {
	conn := RedisCon.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

func Delete(key string) (bool, error) {
	conn := RedisCon.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}
