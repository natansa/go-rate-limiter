package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// middleware que limita o número de requisições por IP ou API_KEY
func rateLimiterMiddleware(persistence IPersistenceStrategy) gin.HandlerFunc {
	return func(context *gin.Context) {
		// checa se o número de requisições por IP ou API_KEY excedeu o limite
		ratelimit, requests := checkRateLimit(context, persistence)
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
