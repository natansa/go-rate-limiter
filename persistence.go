package main

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
)

func initialize() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}

func setTimeToExpireKey(context *gin.Context, expire int, key string) error {
	ctx := context.Request.Context()
	// define a expiração do contador de requisições por IP ou API_KEY no redis
	expireInternal := time.Second * time.Duration(expire)
	if err := redisClient.Expire(ctx, key, expireInternal).Err(); err != nil {
		return err
	}
	return nil
}

func getRequestCount(context *gin.Context, key string) (int64, error) {
	ctx := context.Request.Context()
	// incrementa o contador de requisições por IP ou API_KEY no redis
	requests_count, err := redisClient.Incr(ctx, key).Result()
	return requests_count, err
}
