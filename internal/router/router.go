package router

import (
	"api_authentication/internal/middlewares"
	"api_authentication/internal/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Inicialize o repositório e serviço de usuário
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)

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
