package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var (
	redisClient *redis.Client
)

func init() {
	// carrega o arquivo .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Erro ao carregar o arquivo .env")
		os.Exit(1)
	}

	// inicializa o redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}

func main() {
	// inicializa o servidor
	router := gin.Default()
	router.Use(rateLimiterMiddleware())
	router.GET("/")
	router.Run(":8080")
}

// middleware que limita o número de requisições por IP ou API_KEY
func rateLimiterMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		// checa se o número de requisições por IP ou API_KEY excedeu o limite
		ratelimit, requests := checkRateLimit(context)
		if ratelimit {
			// retorna erro 429 caso o número de requisições por IP ou API_KEY tenha excedido o limite
			context.String(http.StatusTooManyRequests, "you have reached the maximum number of requests or actions allowed within a certain time frame")
			context.Abort()
			return
		}

		// retorna ok caso o número de requisições por IP ou API_KEY não tenha excedido o limite
		context.String(http.StatusOK, "REQUEST "+strconv.FormatInt(requests, 10)+" ALLOWED")
	}
}

// verifica se o número de requisições por IP ou API_KEY excedeu o limite
func checkRateLimit(context *gin.Context) (bool, int64) {
	limit_type := os.Getenv("LIMIT_TYPE")
	var expire int
	var limit int
	var key string
	var err error

	if limit_type == "IP" {
		// busca o tempo de bloqueio por IP no arquivo .env
		expire, err = strconv.Atoi(os.Getenv("BLOCKING_TIME_OUT_IN_SECONDS_PER_IP"))
		if err != nil {
			// retorna erro 500 em caso de erro ao buscar o tempo de bloqueio por IP no arquivo .env
			context.String(http.StatusInternalServerError, err.Error())
			context.Abort()
		}
		limit, err = getLimitFromEnv("IP")
		if err != nil {
			// retorna erro 500 em caso de erro ao buscar o limite de requisições por IP no arquivo .env
			context.String(http.StatusInternalServerError, err.Error())
			context.Abort()
		}
		key = context.ClientIP()
	} else if limit_type == "API_KEY" {
		// busca o tempo de bloqueio por API_KEY no arquivo .env
		expire, err = strconv.Atoi(os.Getenv("BLOCKING_TIME_OUT_IN_SECONDS_PER_API_KEY"))
		if err != nil {
			// retorna erro 500 em caso de erro ao buscar o tempo de bloqueio por API_KEY no arquivo .env
			context.String(http.StatusInternalServerError, err.Error())
			context.Abort()
		}
		key = context.GetHeader("API_KEY")
		limit, err = getLimitFromEnv(key)
		if err != nil {
			// retorna erro 500 em caso de erro ao buscar o limite de requisições por API_KEY no arquivo .env
			context.String(http.StatusInternalServerError, err.Error())
			context.Abort()
		}
	} else {
		// retorna erro 500 caso não tenha sido definido o tipo de limite de requisições no arquivo .env
		context.String(http.StatusInternalServerError, "Blocking time in seconds not configured")
		context.Abort()
	}

	ctx := context.Request.Context()
	// incrementa o contador de requisições por IP ou API_KEY no redis
	requests_count, err := redisClient.Incr(ctx, key).Result()

	if err != nil {
		// retorna erro 500 em caso de erro ao incrementar o contador de requisições por IP ou API_KEY no redis
		context.String(http.StatusInternalServerError, err.Error())
		context.Abort()
	}

	if requests_count == 1 {
		// define a expiração do contador de requisições por IP ou API_KEY no redis
		expire := time.Second * time.Duration(expire)
		if err := redisClient.Expire(ctx, key, expire).Err(); err != nil {
			// retorna erro 500 em caso de erro ao definir a expiração do contador de requisições por IP ou API_KEY no redis
			context.String(http.StatusInternalServerError, err.Error())
			context.Abort()
		}
	}

	if int(requests_count) > limit {
		// retorna true caso o número de requisições por IP ou API_KEY tenha excedido o limite
		return true, requests_count
	}

	// retorna false caso o número de requisições por IP ou API_KEY não tenha excedido o limite
	return false, requests_count
}

// retorna o limite de requisições por IP ou API_KEY do arquivo .env
func getLimitFromEnv(key string) (int, error) {
	if key == "IP" {
		// retorna o limite de requisições por IP no arquivo .env
		limit_ip, err := strconv.Atoi(os.Getenv("LIMIT_IP"))
		return limit_ip, err
	}

	// retorna o limite de requisições por API_KEY no arquivo .env
	limit_key_str := os.Getenv(key)
	limit_key, err := strconv.Atoi(limit_key_str)
	return limit_key, err
}
