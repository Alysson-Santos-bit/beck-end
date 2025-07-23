// internal/middlewares/auth.go
package middlewares

import (
	"log" // Importe o pacote log
	"net/http"
	"strings"

	"api_authentication/internal/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de autenticação ausente"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		log.Printf("AuthMiddleware received token string: '%s'", tokenString) // Add this log

		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			log.Printf("Erro de validação JWT no middleware: %v", err) // Adicione este log
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		log.Printf("Token validado com sucesso para userID: %d", claims.UserID) // Adicione este log
		c.Next()
	}
}
