package service

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func getArticleRedisKey(author string, keyword string) string {
	paramKey := strings.TrimSpace(fmt.Sprintf("%s-%s", author, keyword))
	encodedParamKey := base64.StdEncoding.EncodeToString([]byte(paramKey))
	redisKey := fmt.Sprintf("articles:%s", encodedParamKey)
	return redisKey
}
