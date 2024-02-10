package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// carrega o arquivo .env
	err := godotenv.Load()
	if err != nil {
		panic("Erro ao carregar o arquivo .env")
	}

	// inicializa o redis
	initialize()
}

func main() {
	// inicializa o servidor
	router := gin.Default()
	router.Use(rateLimiterMiddleware())
	router.GET("/")
	router.Run(":8080")
}
