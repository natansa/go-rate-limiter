package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var persistence IPersistenceStrategy

func init() {
	// carrega o arquivo .env
	err := godotenv.Load()
	if err != nil {
		panic("Erro ao carregar o arquivo .env")
	}

	// inicializa o redis

	switch os.Getenv("PERSISTENCE_TYPE") {
	case "REDIS":
		persistence = &RedisPersistence{}
	case "POSTGRES":
		persistence = &PostgresPersistence{}
	default:
		panic("Tipo de persistência não suportado")
	}

	persistence.Initialize()
}

func main() {
	// inicializa o servidor
	router := gin.Default()
	router.Use(rateLimiterMiddleware(persistence))
	router.GET("/")
	router.Run(":8080")
}
