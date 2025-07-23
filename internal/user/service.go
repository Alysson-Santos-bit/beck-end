package user

import (
	"net/http"
	"strconv" // Para converter string para uint

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10" // Para validação de requisições
	"gorm.io/gorm"                           // Importe gorm para verificar "record not found"

	"api_authentication/internal/auth" // Para hashing de senha e JWT
)

// UserService define a interface para operações de serviço de usuário
type UserService interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetUserByID(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	GetCurrentUser(c *gin.Context) // <--- Esta linha está correta aqui na interface
}

// userServiceImpl é a implementação concreta do UserService
type userServiceImpl struct {
	repo     UserRepository
	validate *validator.Validate // Validador para structs
}

// NewUserService cria uma nova instância de UserService
func NewUserService(repo UserRepository) UserService {
	return &userServiceImpl{
		repo:     repo,
		validate: validator.New(),
	}
}

// Register lida com o registro de um novo usuário
func (s *userServiceImpl) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados de requisição inválidos: " + err.Error()})
		return
	}

	// Validação dos campos da requisição
	if err := s.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Campos obrigatórios ausentes ou inválidos: " + err.Error()})
		return
	}

	// --- VERIFICAÇÃO DE UNICIDADE DO USERNAME ---
	_, err := s.repo.GetUserByUsername(req.Username)
	if err == nil { // Se não houve erro, o username já existe
		c.JSON(http.StatusConflict, gin.H{"message": "Nome de usuário já existe."})
		return
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao verificar nome de usuário existente."})
		return
	}

	// --- VERIFICAÇÃO DE UNICIDADE DO EMAIL ---
	_, err = s.repo.GetUserByEmail(req.Email)
	if err == nil { // Se não houve erro, o email já existe
		c.JSON(http.StatusConflict, gin.H{"message": "Email já cadastrado."})
		return
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao verificar email existente."})
		return
	}

	// Hash da senha
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao hashear senha."})
		return
	}

	// Criar novo usuário
	newUser := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.repo.CreateUser(newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao registrar usuário."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuário registrado com sucesso!"})
}

// Login lida com a autenticação do usuário
func (s *userServiceImpl) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados de requisição inválidos: " + err.Error()})
		return
	}

	// Validação dos campos da requisição
	if err := s.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Campos obrigatórios ausentes ou inválidos: " + err.Error()})
		return
	}

	// --- ALTERAÇÃO CRÍTICA AQUI: USANDO GetUserByUsernameOrEmail ---
	user, err := s.repo.GetUserByUsernameOrEmail(req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Credenciais inválidas."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao buscar usuário."})
		return
	}

	// Verificar a senha
	if !auth.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Credenciais inválidas."})
		return
	}

	// Gerar JWT
	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao gerar token JWT."})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: token})
}

// GetCurrentUser busca o perfil do usuário logado usando o ID do JWT
// --- ESTA FUNÇÃO FOI MOVIDA PARA FORA DO MÉTODO Login ---
func (s *userServiceImpl) GetCurrentUser(c *gin.Context) {
	// O userID é definido pelo AuthMiddleware após a validação do token
	userID, exists := c.Get("userID") // Obtenha o userID do contexto
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ID do usuário não encontrado no contexto."})
		return
	}

	// Converta para uint, pois o c.Get retorna interface{}
	id := userID.(uint)

	user, err := s.repo.GetUserByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Usuário não encontrado."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao buscar usuário logado."})
		return
	}

	// Não retornar a senha hasheada
	user.Password = ""
	c.JSON(http.StatusOK, user) // Retorna os dados do usuário
}

// GetUserByID (Rota protegida para obter um usuário por ID)
func (s *userServiceImpl) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 32) // Converte string para uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID de usuário inválido."})
		return
	}

	user, err := s.repo.GetUserByID(uint(userID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Usuário não encontrado."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao buscar usuário."})
		return
	}

	// Não retornar a senha hasheada
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// UpdateUser (Rota protegida para atualizar um usuário)
func (s *userServiceImpl) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID de usuário inválido."})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados de requisição inválidos: " + err.Error()})
		return
	}

	// Validação dos campos da requisição de atualização
	if err := s.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Campos obrigatórios ausentes ou inválidos: " + err.Error()})
		return
	}

	// Buscar o usuário existente
	user, err := s.repo.GetUserByID(uint(userID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Usuário não encontrado."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao buscar usuário para atualização."})
		return
	}

	// Aplicar as atualizações apenas se os campos forem fornecidos
	if req.Username != nil {
		// Verificar se o novo nome de usuário já existe, se for diferente do atual
		if *req.Username != user.Username {
			existingUser, err := s.repo.GetUserByUsername(*req.Username)
			if err == nil && existingUser.ID != user.ID { // Se encontrou outro usuário com o mesmo nome
				c.JSON(http.StatusConflict, gin.H{"message": "Nome de usuário já em uso."})
				return
			}
			if err != nil && err != gorm.ErrRecordNotFound { // Outro erro ao verificar
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao verificar nome de usuário."})
				return
			}
		}
		user.Username = *req.Username
	}
	if req.Email != nil {
		// Verificar se o novo email já existe, se for diferente do atual
		if *req.Email != user.Email {
			existingUser, err := s.repo.GetUserByEmail(*req.Email)
			if err == nil && existingUser.ID != user.ID {
				c.JSON(http.StatusConflict, gin.H{"message": "Email já em uso."})
				return
			}
			if err != nil && err != gorm.ErrRecordNotFound {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao verificar email."})
				return
			}
		}
		user.Email = *req.Email
	}
	if req.Password != nil {
		hashedPassword, err := auth.HashPassword(*req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao hashear nova senha."})
			return
		}
		user.Password = hashedPassword
	}

	// Salvar as alterações no banco de dados
	if err := s.repo.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao atualizar usuário."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuário atualizado com sucesso!"})
}

// DeleteUser (Rota protegida para deletar um usuário)
func (s *userServiceImpl) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID de usuário inválido."})
		return
	}

	// Verificar se o usuário existe antes de tentar deletar
	_, err = s.repo.GetUserByID(uint(userID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Usuário não encontrado."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao verificar usuário para exclusão."})
		return
	}

	if err := s.repo.DeleteUser(uint(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao deletar usuário."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuário deletado com sucesso!"})
}
