package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ulwisibaq/kumparan-test/internal/article"
	"github.com/ulwisibaq/kumparan-test/internal/models"
)

type articleHandler struct {
	articleService article.Service
}

func NewArticleHandler(
	articleService article.Service,
) *articleHandler {
	return &articleHandler{
		articleService: articleService,
	}
}

func (ah *articleHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	author := r.FormValue("author")
	keyword := r.FormValue("keyword")

	resp, err := ah.articleService.GetArticles(author, keyword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := models.Response{
		Data: resp,
	}

	respondWithJSON(w, http.StatusOK, response)

}

func (ah *articleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	article := models.Articles{}

	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
	}

	resp, err := ah.articleService.CreateArticle(article)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := models.Response{
		Data: resp,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"Error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
