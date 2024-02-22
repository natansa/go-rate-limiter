package main

import "github.com/gin-gonic/gin"

type PostgresPersistence struct{}

func (p *PostgresPersistence) Initialize() {
	// Lógica de inicialização específica do PostgreSQL
}

func (p *PostgresPersistence) SetTimeToExpireKey(context *gin.Context, expire int, key string) error {
	// Lógica para definir o tempo de expiração de uma chave no PostgreSQL
	return nil
}

func (p *PostgresPersistence) GetRequestCount(context *gin.Context, key string) (int64, error) {
	// incrementa o contador de requisições por IP ou API_KEY no PostgreSQL
	return 0, nil
}
