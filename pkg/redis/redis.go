package redis

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	// redis "github.com/go-redis/redis/v8"
	"github.com/gomodule/redigo/redis"
)

var (
	Pool *redis.Pool
)

type RedisConfig struct {
	Connection string
	Password   string
	Timeout    time.Duration
	MaxIdle    int
}

func InitRedis(cfg RedisConfig) {
	if cfg.Connection == "" {
		cfg.Connection = ":6379"
	}
	Pool = NewPool(&cfg)
	cleanupHook()

	err := Ping()
	if err != nil {
		log.Println("failed to init redis")
	}
}

func NewPool(config *RedisConfig) *redis.Pool {

	return &redis.Pool{
		MaxIdle:     config.MaxIdle,
		IdleTimeout: config.Timeout * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Connection)
			if err != nil {
				return nil, err
			}

			// _, err = c.Do("AUTH", config.Password)
			// if err != nil {
			// 	return nil, err
			// }
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func cleanupHook() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		Pool.Close()
	}()
}

func Ping() error {

	conn := Pool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("PING"))
	if err != nil {
		fmt.Printf("cannot 'PING' db: %v", err)
		return fmt.Errorf("cannot 'PING' db: %v", err)
	}
	return nil
}
