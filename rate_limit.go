package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// verifica se o número de requisições por IP ou API_KEY excedeu o limite
func checkRateLimit(context *gin.Context) (bool, int64) {
	var expire int // tempo de bloqueio/expiracao por IP ou API_KEY
	var limit int  // limite de requisições por IP ou API_KEY
	var key string // IP ou API_KEY
	var err error

	// busca API_KEY no header e escolhe o tipo de limite de requisições por IP ou API_KEY
	key = context.GetHeader("API_KEY")
	if key != "" {
		// busca o tempo de bloqueio por API_KEY no arquivo .env
		expire, err = strconv.Atoi(os.Getenv("BLOCKING_TIME_OUT_IN_SECONDS_PER_API_KEY"))
		if err != nil {
			// retorna erro 500 em caso de erro ao buscar o tempo de bloqueio por API_KEY no arquivo .env
			throwError(context, err)
		}

		// busca o limite de requisições por API_KEY no arquivo .env
		limit, err = getLimitFromEnv(key, "")
		if err != nil {
			// retorna erro 500 em caso de erro ao buscar o limite de requisições por API_KEY no arquivo .env
			throwError(context, err)
		}
	} else {
		// define o IP como key
		key = context.ClientIP()
		// substitui o IP ::1 (caso seja localhost) por 127.0.0.1
		if key == "::1" {
			key = "127.0.0.1"
		}
		// busca o tempo de bloqueio por IP no arquivo .env
		expire, err = strconv.Atoi(os.Getenv("BLOCKING_TIME_OUT_IN_SECONDS_PER_IP"))
		if err != nil {
			// retorna erro 500 em caso de erro ao buscar o tempo de bloqueio por IP no arquivo .env
			throwError(context, err)
		}

		// busca o limite de requisições por IP no arquivo .env
		limit, err = getLimitFromEnv("IP", key)
		if err != nil {
			// retorna erro 500 em caso de erro ao buscar o limite de requisições por IP no arquivo .env
			throwError(context, err)
		}
	}

	requests_count, err := getRequestCount(context, key)
	if err != nil {
		// retorna erro 500 em caso de erro ao incrementar o contador de requisições por IP ou API_KEY
		throwError(context, err)
	}

	if requests_count == 1 {
		err = setTimeToExpireKey(context, expire, key)
		if err != nil {
			// retorna erro 500 em caso de erro ao definir a expiração do contador de requisições por IP ou API_KEY
			throwError(context, err)
		}
	}

	if int(requests_count) > limit {
		// retorna true caso o número de requisições por IP ou API_KEY tenha excedido o limite
		return true, requests_count
	}

	// retorna false caso o número de requisições por IP ou API_KEY não tenha excedido o limite
	return false, requests_count
}

func throwError(context *gin.Context, err error) {
	context.String(http.StatusInternalServerError, err.Error())
	context.Abort()
}

// retorna o limite de requisições por IP ou API_KEY do arquivo .env
func getLimitFromEnv(key string, ip string) (int, error) {
	if key == "IP" {
		// retorna o limite de requisições por IP no arquivo .env
		limit_ip, err := strconv.Atoi(os.Getenv("LIMIT_IP_" + ip))
		return limit_ip, err
	}

	// retorna o limite de requisições por API_KEY no arquivo .env
	limit_key, err := strconv.Atoi(os.Getenv("LIMIT_" + key))
	return limit_key, err
}
