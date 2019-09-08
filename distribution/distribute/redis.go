package distribute

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	IDLE_COUNT   = 20  //连接池空闲个数
	ACTIVE_COUNT = 20  //连接池活动个数
	IDLE_TIMEOUT = 180 //空闲超时时间
	REQKEY       = "REQUEST_KEY"
	RESPKEY      = "RESPONSE_KEY"
)

var (
	REDIS_DB     int
	RedisClients *redis.Pool
)

func GetRedisConn() (conn redis.Conn) {
	return RedisClients.Get()
}

func init() {
	REDIS_DB = 0
	RedisClients = createPool(IDLE_COUNT, ACTIVE_COUNT, IDLE_TIMEOUT, "127.0.0.1:6379")
}

func createPool(maxIdle, maxActive, idleTimeout int, address string) (obj *redis.Pool) {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	// obj = new(redis.Pool)
	// obj.MaxIdle = maxIdle
	// obj.MaxActive = maxActive
	// obj.Wait = true
	// obj.IdleTimeout = (time.Duration)(idleTimeout) * time.Second
	// obj.Dial = func() (redis.Conn, error) {
	// 	c, err := redis.Dial("tcp", address)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	c.Do("SELECT", REDIS_DB)
	// 	return c, err
	// }
	// obj.TestOnBorrow = func(c redis.Conn, t time.Time) error {
	// 	if time.Since(t) < time.Minute {
	// 		return nil
	// 	}
	// 	_, err := c.Do("PING")
	// 	return err
	// }
	// return
}
