// internal/user/repository.go
package user

import (
	"gorm.io/gorm"
)

// UserRepository define a interface para operações de persistência de usuário
type UserRepository interface {
	CreateUser(user *User) error
	GetUserByUsername(username string) (*User, error)
	GetUserByID(id uint) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id uint) error
	// --- ESTES DOIS MÉTODOS ESTAVAM FALTANDO NA INTERFACE! ---
	GetUserByUsernameOrEmail(identifier string) (*User, error)
	GetUserByEmail(email string) (*User, error) // <--- MÉTODO ADICIONADO À INTERFACE
}

// userRepositoryImpl é a implementação concreta do UserRepository
type userRepositoryImpl struct {
	db *gorm.DB // O campo do GORM DB é 'db' minúsculo
}

// NewUserRepository cria uma nova instância de UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db} // Inicializa o campo 'db' corretamente
}

// CreateUser cria um novo usuário no banco de dados
func (r *userRepositoryImpl) CreateUser(user *User) error {
	return r.db.Create(user).Error
}

// GetUserByUsername busca um usuário pelo nome de usuário
func (r *userRepositoryImpl) GetUserByUsername(username string) (*User, error) {
	var user User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsernameOrEmail busca um usuário pelo nome de usuário ou email
func (r *userRepositoryImpl) GetUserByUsernameOrEmail(identifier string) (*User, error) {
	var user User
	// Esta é a query chave: busca onde o username OU o email correspondem ao identificador
	if err := r.db.Where("username = ? OR email = ?", identifier, identifier).First(&user).Error; err != nil {
		return nil, err // Retornará gorm.ErrRecordNotFound se não encontrar
	}
	return &user, nil
}

// GetUserByEmail busca um usuário pelo email
// --- IMPLEMENTAÇÃO ADICIONADA/CONFIRMADA AQUI ---
func (r *userRepositoryImpl) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID busca um usuário pelo ID
func (r *userRepositoryImpl) GetUserByID(id uint) (*User, error) {
	var user User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser atualiza um usuário existente no banco de dados
func (r *userRepositoryImpl) UpdateUser(user *User) error {
	return r.db.Save(user).Error
}

// DeleteUser deleta um usuário pelo ID
func (r *userRepositoryImpl) DeleteUser(id uint) error {
	return r.db.Delete(&User{}, id).Error
}
