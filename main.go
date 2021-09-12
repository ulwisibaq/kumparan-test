package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/ulwisibaq/kumparan-test/internal/article/repository"
	"github.com/ulwisibaq/kumparan-test/internal/article/service"
	"github.com/ulwisibaq/kumparan-test/internal/config"
	"github.com/ulwisibaq/kumparan-test/internal/handler"
	"github.com/ulwisibaq/kumparan-test/pkg/redis"
)

func main() {

	cfg := &config.MainConfig{}
	config.ReadConfig(cfg)

	//init redis connection
	redisConfig := redis.RedisConfig{
		Connection: cfg.Redis.Connection,
		Password:   cfg.Redis.Password,
		Timeout:    cfg.Redis.Timeout,
		MaxIdle:    cfg.Redis.MaxIdle,
	}
	redis.InitRedis(redisConfig)

	// setup mysql db connection
	db, err := sqlx.Open(`mysql`, cfg.Database.MysqlDSN)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	articleRepo := repository.NewArticleRepository(db)
	redisRepo := repository.NewRedisRepository()
	articleService := service.NewArticleService(articleRepo, redisRepo)
	articleHandler := handler.NewArticleHandler(articleService)

	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/articles", articleHandler.GetArticles).Methods("GET")
	myRouter.HandleFunc("/articles", articleHandler.CreateArticle).Methods("POST")

	srv := http.Server{
		Addr:    ":8080",
		Handler: myRouter,
	}
	go func() {
		log.Println("http server starting in port :8080")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}
