// internal/auth/jwt.go
package auth

import (
	"errors" // Import errors package
	"log"    // Import log
	"time"

	"api_authentication/configs" // Importe seu pacote de configs

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func init() {
	configs.LoadEnv()
	secret := configs.GetEnv("JWT_SECRET", "")
	if secret == "" {
		panic("JWT_SECRET não está definido nas variáveis de ambiente.")
	}
	jwtSecret = []byte(secret)
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uint) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Token válido por 24 horas
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (*Claims, error) {
	log.Printf("Attempting to validate tokenString: '%s'", tokenString) // Log the token string

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de assinatura inesperado")
		}
		return jwtSecret, nil
	})

	// CRITICAL FIX: Check for error immediately after parsing the token
	if err != nil {
		log.Printf("Error parsing JWT: %v", err) // Log the error from parsing
		return nil, err                          // Return the error if parsing failed (e.g., malformed token, invalid signature, expired)
	}

	// Ensure the token is not nil and is valid
	if !token.Valid {
		log.Printf("JWT token is not valid after parsing. Token Valid: %t", token.Valid) // Log invalid token
		return nil, errors.New("token JWT inválido")
	}

	// Type assert the claims
	claims, ok := token.Claims.(*Claims)
	if !ok {
		log.Printf("Claims are not of expected type. Claims Type: %T", token.Claims) // Log claims type mismatch
		return nil, errors.New("claims do token JWT inválidos")
	}

	log.Printf("JWT token valid for userID: %d", claims.UserID) // Log successful validation
	return claims, nil
}
