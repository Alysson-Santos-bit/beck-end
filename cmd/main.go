package main

import (
	"api_authentication/configs"
	"api_authentication/internal/database"
	"api_authentication/internal/router"
	"log"
	"net/http"

	"github.com/rs/cors" // 1. Importe a biblioteca CORS
)

func main() {
	// 1. Carregar variáveis de ambiente
	configs.LoadEnv()

	// 2. Conectar ao banco de dados
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Não foi possível conectar ao banco de dados: %v", err)
	}
	log.Println("Conexão com o banco de dados estabelecida com sucesso!")

	// 3. Configurar e iniciar o roteador (gorilla/mux)
	r := router.SetupRouter(db)

	// --- Começo da adição para CORS ---

	// 4. Configuração do CORS
	// É CRÍTICO definir AllowedOrigins para o domínio onde seu FRONT-END está rodando.
	// Se você estiver usando o Live Server do VS Code, ele geralmente roda em http://127.0.0.1:5500 ou http://localhost:5500.
	// Verifique a URL exata do seu front-end no navegador.
	// Na sua main.go da API Go
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			// Origens para desenvolvimento local (se você testar localmente no futuro)
			"http://localhost:5500",
			"http://127.0.0.1:5500",
			// IMPORTANTE: Suas URLs do GitHub Pages - COM E SEM A BARRA FINAL
			"https://alysson-santos-bit.github.io/front-end",  // URL do seu front-end (sem barra final)
			"https://alysson-santos-bit.github.io/front-end/", // URL do seu front-end (com barra final)
			// URL da sua própria API no Render (útil para testes diretos ou se a API tiver front-end próprio)
			"https://beck-end-oafv.onrender.com",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	// 5. Envolver o roteador com o manipulador CORS
	// Agora, 'handler' é o roteador 'r' com as configurações de CORS aplicadas.
	handler := c.Handler(r)

	// --- Fim da adição para CORS ---

	// 6. Iniciar o servidor HTTP
	port := configs.GetEnv("PORT", configs.GetEnv("API_PORT", "8080")) // Preferir "PORT" do Heroku, senão "API_PORT"
	log.Printf("Servidor iniciado na porta :%s", port)
	// Passe o 'handler' (que inclui CORS) para ListenAndServe
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
