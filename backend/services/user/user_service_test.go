package user

import (
	"testing"

	"desafio-tecnico-fullstack/backend/models"

	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	user *models.User
}

func (m *mockUserRepo) GetUserByCPF(cpf string) *models.User {
	return m.user
}

func (m *mockUserRepo) AddUser(u models.User) error {
	return nil
}

func TestAuthenticateUser_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)
	user := &models.User{ID: 1, CPF: "12345678901", Password: string(hash)}
	repo := &mockUserRepo{user: user}
	service := &userService{
		repo:        repo,
		generateJWT: func(userID int) (string, error) { return "token123", nil },
	}

	token, _, err := service.AuthenticateUser("12345678901", "senha123")
	if err != nil {
		t.Fatalf("esperava sucesso, obteve erro: %v", err)
	}
	if token != "token123" {
		t.Errorf("esperava token 'token123', obteve '%s'", token)
	}
}

func TestAuthenticateUser_UserNotFound(t *testing.T) {
	repo := &mockUserRepo{user: nil}
	service := &userService{
		repo:        repo,
		generateJWT: func(userID int) (string, error) { return "token123", nil },
	}

	_, _, err := service.AuthenticateUser("00000000000", "senha123")
	if err == nil || err.Error() != "usuário ou senha inválidos" {
		t.Errorf("esperava erro de usuário ou senha inválidos, obteve: %v", err)
	}
}

func TestAuthenticateUser_WrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)
	user := &models.User{ID: 1, CPF: "12345678901", Password: string(hash)}
	repo := &mockUserRepo{user: user}
	service := &userService{
		repo:        repo,
		generateJWT: func(userID int) (string, error) { return "token123", nil },
	}

	_, _, err := service.AuthenticateUser("12345678901", "errada")
	if err == nil || err.Error() != "usuário ou senha inválidos" {
		t.Errorf("esperava erro de usuário ou senha inválidos, obteve: %v", err)
	}
}
