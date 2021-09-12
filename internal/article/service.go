package article

import "github.com/ulwisibaq/kumparan-test/internal/models"

type Service interface {
	GetArticles(author string, keyword string) (resp []models.Articles, err error)
	CreateArticle(article models.Articles) (resp models.Articles, err error)
}
