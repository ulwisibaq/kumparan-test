package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/ulwisibaq/kumparan-test/internal/article"
	"github.com/ulwisibaq/kumparan-test/internal/models"
)

type ArticleRepo struct {
	db *sqlx.DB
}

func NewArticleRepository(db *sqlx.DB) article.MysqlRepository {
	return &ArticleRepo{db: db}
}

func (ar ArticleRepo) GetArticles(author string, keyword string) (resp []models.Articles, err error) {

	err = ar.db.Select(
		&resp,
		GetArticlesQuery,
		"%"+keyword+"%",
		"%"+keyword+"%",
		author,
		author,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return resp, nil
		}
		return
	}

	return
}

func (ar ArticleRepo) CreateArticle(article models.Articles) (resp models.Articles, err error) {

	res, err := ar.db.Exec(
		CreateArticleQuery,
		article.Author,
		article.Title,
		article.Body,
		article.Created,
	)
	if err != nil {
		return
	}
	id, _ := res.LastInsertId()

	resp = article
	resp.ID = id

	return
}
