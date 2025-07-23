// internal/user/models.go
package user

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"` // `json:"-"` para não serializar a senha
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Para payload de registro
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// Para payload de login
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Para payload de atualização de usuário (campos opcionais)
type UpdateUserRequest struct {
	Username *string `json:"username" validate:"omitempty,min=3,max=30"` // Ponteiro para indicar que é opcional
	Email    *string `json:"email" validate:"omitempty,email"`
	Password *string `json:"password" validate:"omitempty,min=6"`
}

// Para payload de resposta de login
type LoginResponse struct {
	Token string `json:"token"`
}
