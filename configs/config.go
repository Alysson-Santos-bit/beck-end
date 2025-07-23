// configs/configs.go
package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv" // Biblioteca para carregar variáveis de um arquivo .env
)

// LoadEnv carrega as variáveis de ambiente de um arquivo .env
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Atenção: Não foi possível carregar o arquivo .env. As variáveis serão lidas do ambiente.")
	}
}

// GetEnv recupera o valor de uma variável de ambiente, com um fallback opcional
func GetEnv(key string, fallback ...string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	log.Fatalf("Variável de ambiente %s não definida e sem valor padrão.", key)
	return "" // Nunca será alcançado devido ao Fatalf
}
