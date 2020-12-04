package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	pool *redis.Pool
	redisHost = "192.168.0.17:6379"
)

func init() {
	pool = NewRedisPool()
}

func getConn()(redis.Conn, error){
	c, err := redis.Dial("tcp", redisHost)
	if err != nil{
		fmt.Println(err)
		return nil, err
	}
	return c, nil
}

func NewRedisPool() *redis.Pool{
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: getConn,
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},

	}
}

func RedisPool() *redis.Pool{
	return pool
}

