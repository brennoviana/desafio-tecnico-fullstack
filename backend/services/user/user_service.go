package user

import (
	"desafio-tecnico-fullstack/backend/models"
	"desafio-tecnico-fullstack/backend/storage/repository/user"
	"desafio-tecnico-fullstack/backend/utils"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(name, cpf, password string) error
	AuthenticateUser(cpf, password string) (string, *models.User, error)
}

type userService struct {
	repo        user.UserRepository
	generateJWT func(userID int) (string, error)
}

func NewUserService(repo user.UserRepository) UserService {
	return &userService{
		repo:        repo,
		generateJWT: utils.GenerateJWT,
	}
}

func (s *userService) RegisterUser(name, cpf, password string) error {
	if !isValidCPF(cpf) {
		return errors.New("cpf inválido")
	}
	if len(password) < 6 {
		return errors.New("senha muito curta")
	}
	if existing := s.repo.GetUserByCPF(cpf); existing != nil {
		return errors.New("usuário já existe")
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{Name: name, CPF: cpf, Password: string(hash)}
	err := s.repo.AddUser(user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return errors.New("usuário já existe")
		}
		return err
	}
	return nil
}

func (s *userService) AuthenticateUser(cpf, password string) (string, *models.User, error) {
	user := s.repo.GetUserByCPF(cpf)
	if user == nil {
		return "", nil, errors.New("usuário ou senha inválidos")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", nil, errors.New("usuário ou senha inválidos")
	}
	token, err := s.generateJWT(user.ID)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func isValidCPF(cpf string) bool {
	return len(cpf) == 11
}
