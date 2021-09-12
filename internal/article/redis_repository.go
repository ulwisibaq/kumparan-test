package article

import "github.com/ulwisibaq/kumparan-test/internal/models"

type RedisRepositoryInterface interface {
	GetArticlesFromRedis(redisKey string) (resp []models.Articles, err error)
	CacheArticlesData(redisKey string, articles []models.Articles)
	DeleteCache(redisKey string)
}
