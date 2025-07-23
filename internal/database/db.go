package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"api_authentication/internal/user"
)

func ConnectDB() (*gorm.DB, error) {
	var dsn string

	// Tenta obter a DATABASE_URL (usada pelo Heroku)
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		dsn = databaseURL
		log.Println("Usando DATABASE_URL para conexão com o banco de dados.")
	} else {
		// Se DATABASE_URL não estiver definida, usa as variáveis de ambiente locais
		log.Println("DATABASE_URL não definida, usando variáveis de ambiente locais para conexão.")
		dbHost := os.Getenv("DB_HOST")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		dbPort := os.Getenv("DB_PORT")
		dbSSLMode := os.Getenv("DB_SSLMODE")

		if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" {
			return nil, fmt.Errorf("variáveis de ambiente DB_HOST, DB_USER, DB_PASSWORD, DB_NAME ou DB_PORT não definidas para conexão local")
		}

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao banco de dados: %w", err)
	}

	err = db.AutoMigrate(&user.User{})
	if err != nil {
		log.Fatalf("Falha ao migrar o banco de dados: %v", err)
		return nil, fmt.Errorf("falha ao migrar o banco de dados: %w", err)
	}

	return db, nil
}
