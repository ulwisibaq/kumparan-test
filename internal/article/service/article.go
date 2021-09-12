package service

import (
	"time"

	"github.com/ulwisibaq/kumparan-test/internal/article"
	"github.com/ulwisibaq/kumparan-test/internal/models"
)

type ArticleService struct {
	articleRepository article.MysqlRepository
	redisRepository   article.RedisRepositoryInterface
}

func NewArticleService(articleRepo article.MysqlRepository, redisRepo article.RedisRepositoryInterface) article.Service {
	return &ArticleService{articleRepository: articleRepo, redisRepository: redisRepo}
}

func (as ArticleService) GetArticles(author string, keyword string) (resp []models.Articles, err error) {

	redisKey := getArticleRedisKey(author, keyword)

	articles, err := as.redisRepository.GetArticlesFromRedis(redisKey)
	if err != nil {
		return
	}

	if articles != nil {
		return articles, nil
	}

	resp, err = as.articleRepository.GetArticles(author, keyword)
	if err != nil {
		return
	}

	if len(resp) > 0 {
		as.redisRepository.CacheArticlesData(redisKey, resp)
	}

	return
}

func (as ArticleService) CreateArticle(article models.Articles) (resp models.Articles, err error) {

	article.Created = time.Now()

	resp, err = as.articleRepository.CreateArticle(article)
	if err != nil {
		return
	}

	// delete cache for all data
	key := getArticleRedisKey("", "")
	as.redisRepository.DeleteCache(key)

	return
}
