package repository

import (
	"encoding/json"
	"log"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/ulwisibaq/kumparan-test/internal/article"
	"github.com/ulwisibaq/kumparan-test/internal/models"
	"github.com/ulwisibaq/kumparan-test/pkg/redis"
)

type redisRepository struct {
}

func NewRedisRepository() article.RedisRepositoryInterface {
	return &redisRepository{}
}

func (r *redisRepository) GetArticlesFromRedis(redisKey string) (resp []models.Articles, err error) {

	conn := redis.Pool.Get()
	defer conn.Close()

	isKeyExists, err := redigo.Bool(conn.Do("EXISTS", redisKey))
	if err != nil {
		return
	}

	if isKeyExists {
		articlesBytes, errGet := redigo.Bytes(conn.Do("GET", redisKey))
		if errGet != nil {
			return resp, errGet
		}

		errGet = json.Unmarshal(articlesBytes, &resp)
		if errGet != nil {
			return resp, errGet
		}
	}

	return
}

func (r *redisRepository) CacheArticlesData(redisKey string, articles []models.Articles) {

	conn := redis.Pool.Get()
	defer conn.Close()

	value, err := json.Marshal(articles)
	if err != nil {
		log.Println("error marshalling to json.")
	}

	_, err = conn.Do("SET", redisKey, value)
	if err != nil {
		log.Println("error when SET to redis.")
	}

	// set to expire in 1 hour
	expInSec := int64(3600)
	_, err = conn.Do("EXPIRE", redisKey, expInSec)
	if err != nil {
		log.Println("error when set EXPIRE to redis.")
	}
}

func (r *redisRepository) DeleteCache(redisKey string) {
	conn := redis.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", redisKey)
	if err != nil {
		log.Println("error when DELETE key to redis")
	}
}
