package router

import (
	"api_authentication/internal/middlewares"
	"api_authentication/internal/user"

	"time" // Para configurar o MaxAge, se desejar

	"github.com/gin-contrib/cors" // Importe o middleware CORS para Gin
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// --- Adição do Middleware CORS para Gin ---
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5500",                           // Para desenvolvimento local
			"http://127.0.0.1:5500",                           // Para desenvolvimento local
			"https://Alysson-Santos-bit.github.io/front-end",  // Seu front-end (sem barra final)
			"https://Alysson-Santos-bit.github.io/front-end/", // Seu front-end (com barra final)
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			// Se você tiver muitas origens, pode usar esta função
			// para verificar dinamicamente. Por enquanto, a lista acima é suficiente.
			return true // Ou uma lógica mais complexa baseada na lista AllowedOrigins
		},
		MaxAge: 12 * time.Hour, // Tempo em que as informações de preflight podem ser cacheadas
	}))
	// --- Fim da adição do Middleware CORS para Gin ---

	// Inicialize o repositório e serviço de usuário
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)

	// Inicialize o repositório e serviço de usuário

	// Rotas de autenticação (públicas)
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", userService.Register)
		authRoutes.POST("/login", userService.Login)
	}

	// Rotas protegidas (exigem JWT)
	authMiddleware := middlewares.AuthMiddleware() // Instancie o middleware
	privateRoutes := r.Group("/api", authMiddleware)
	{
		// ... (outras rotas existentes)
		privateRoutes.GET("/users/:id", userService.GetUserByID)
		privateRoutes.PUT("/users/:id", userService.UpdateUser)
		privateRoutes.DELETE("/users/:id", userService.DeleteUser)

		// --- NOVA ROTA PROTEGIDA PARA BUSCAR O USUÁRIO LOGADO ---
		privateRoutes.GET("/perfil", userService.GetCurrentUser) // <--- ADICIONE ESTA LINHA
		// Agora o frontend pode chamar /api/perfil
	}

	return r
}
