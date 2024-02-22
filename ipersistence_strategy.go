package main

import "github.com/gin-gonic/gin"

// IPersistenceStrategy define as operações de persistência para o rate limiter.
type IPersistenceStrategy interface {
	Initialize()
	SetTimeToExpireKey(context *gin.Context, expire int, key string) error
	GetRequestCount(context *gin.Context, key string) (int64, error)
}
