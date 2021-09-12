## Kumparan technical test: REST API for get and create articles

### How to start
#### 1. Start the dev environment(mysql and redis) using this command
```
make start-dev
```
or use docker-compose command
```
docker-compose up -d
```



#### 2. Create table articles on database
```
docker exec -i kumparan-test_db_1 mysql -uroot -proot test_db < article.sql
```
or you can usq the DDL in the file.

#### 3. run the app
```
go run main.go
```

