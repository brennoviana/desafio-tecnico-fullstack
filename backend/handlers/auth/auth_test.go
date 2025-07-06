package auth

import (
	"desafio-tecnico-fullstack/backend/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockUserService struct {
	registerErr       error
	authenticateErr   error
	authenticateUser  *models.User
	authenticateToken string
}

func (m *mockUserService) RegisterUser(name, cpf, password string) error {
	return m.registerErr
}

func (m *mockUserService) AuthenticateUser(cpf, password string) (string, *models.User, error) {
	if m.authenticateErr != nil {
		return "", nil, m.authenticateErr
	}
	return m.authenticateToken, m.authenticateUser, nil
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestRegisterHandler_Success(t *testing.T) {
	service := &mockUserService{
		authenticateToken: "test-token",
		authenticateUser: &models.User{
			ID:   1,
			Name: "João Silva",
			CPF:  "12345678901",
		},
	}

	router := setupRouter()
	router.POST("/register", RegisterHandler(service))

	reqBody := `{"name":"João Silva","cpf":"12345678901","password":"senha123"}`
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperava status 200, obteve %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}

	if response["status"] != "success" {
		t.Errorf("esperava status 'success', obteve '%v'", response["status"])
	}

	data := response["data"].(map[string]interface{})
	if data["token"] != "test-token" {
		t.Errorf("esperava token 'test-token', obteve '%v'", data["token"])
	}

	if data["name"] != "João Silva" {
		t.Errorf("esperava nome 'João Silva', obteve '%v'", data["name"])
	}

	if data["cpf"] != "12345678901" {
		t.Errorf("esperava CPF '12345678901', obteve '%v'", data["cpf"])
	}
}

func TestRegisterHandler_InvalidJSON(t *testing.T) {
	service := &mockUserService{}
	router := setupRouter()
	router.POST("/register", RegisterHandler(service))

	req, _ := http.NewRequest("POST", "/register", strings.NewReader(`{invalid json}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperava status 400, obteve %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "error" {
		t.Errorf("esperava status 'error', obteve '%v'", response["status"])
	}

	if response["error"] != "requisição inválida" {
		t.Errorf("esperava erro 'requisição inválida', obteve '%v'", response["error"])
	}
}

func TestRegisterHandler_MissingFields(t *testing.T) {
	testCases := []struct {
		name    string
		reqBody string
	}{
		{"missing name", `{"cpf":"12345678901","password":"senha123"}`},
		{"missing cpf", `{"name":"João Silva","password":"senha123"}`},
		{"missing password", `{"name":"João Silva","cpf":"12345678901"}`},
		{"empty name", `{"name":"","cpf":"12345678901","password":"senha123"}`},
		{"empty cpf", `{"name":"João Silva","cpf":"","password":"senha123"}`},
		{"empty password", `{"name":"João Silva","cpf":"12345678901","password":""}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := &mockUserService{}
			router := setupRouter()
			router.POST("/register", RegisterHandler(service))

			req, _ := http.NewRequest("POST", "/register", strings.NewReader(tc.reqBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("esperava status 400, obteve %d", w.Code)
			}

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if response["error"] != "campos obrigatórios não preenchidos" {
				t.Errorf("esperava erro sobre campos obrigatórios, obteve '%v'", response["error"])
			}
		})
	}
}

func TestRegisterHandler_UserAlreadyExists(t *testing.T) {
	service := &mockUserService{
		registerErr: errors.New("usuário já existe"),
	}

	router := setupRouter()
	router.POST("/register", RegisterHandler(service))

	reqBody := `{"name":"João Silva","cpf":"12345678901","password":"senha123"}`
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("esperava status 409, obteve %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["error"] != "usuário já existe" {
		t.Errorf("esperava erro 'usuário já existe', obteve '%v'", response["error"])
	}
}

func TestRegisterHandler_ValidationErrors(t *testing.T) {
	testCases := []struct {
		name         string
		serviceErr   string
		expectedCode int
	}{
		{"invalid CPF", "cpf inválido", http.StatusBadRequest},
		{"short password", "senha muito curta", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := &mockUserService{
				registerErr: errors.New(tc.serviceErr),
			}

			router := setupRouter()
			router.POST("/register", RegisterHandler(service))

			reqBody := `{"name":"João Silva","cpf":"12345678901","password":"senha123"}`
			req, _ := http.NewRequest("POST", "/register", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Errorf("esperava status %d, obteve %d", tc.expectedCode, w.Code)
			}

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if response["error"] != tc.serviceErr {
				t.Errorf("esperava erro '%s', obteve '%v'", tc.serviceErr, response["error"])
			}
		})
	}
}

func TestRegisterHandler_DatabaseError(t *testing.T) {
	service := &mockUserService{
		registerErr: errors.New("database connection failed"),
	}

	router := setupRouter()
	router.POST("/register", RegisterHandler(service))

	reqBody := `{"name":"João Silva","cpf":"12345678901","password":"senha123"}`
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperava status 500, obteve %d", w.Code)
	}
}

func TestRegisterHandler_AuthenticationFailsAfterRegistration(t *testing.T) {
	service := &mockUserService{
		authenticateErr: errors.New("jwt generation failed"),
	}

	router := setupRouter()
	router.POST("/register", RegisterHandler(service))

	reqBody := `{"name":"João Silva","cpf":"12345678901","password":"senha123"}`
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperava status 500, obteve %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["error"] != "erro ao gerar token" {
		t.Errorf("esperava erro 'erro ao gerar token', obteve '%v'", response["error"])
	}
}

func TestLoginHandler_Success(t *testing.T) {
	service := &mockUserService{
		authenticateToken: "login-token",
		authenticateUser: &models.User{
			ID:   1,
			Name: "Maria Santos",
			CPF:  "98765432109",
		},
	}

	router := setupRouter()
	router.POST("/login", LoginHandler(service))

	reqBody := `{"cpf":"98765432109","password":"senha123"}`
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperava status 200, obteve %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}

	if response["status"] != "success" {
		t.Errorf("esperava status 'success', obteve '%v'", response["status"])
	}

	data := response["data"].(map[string]interface{})
	if data["token"] != "login-token" {
		t.Errorf("esperava token 'login-token', obteve '%v'", data["token"])
	}

	if data["name"] != "Maria Santos" {
		t.Errorf("esperava nome 'Maria Santos', obteve '%v'", data["name"])
	}
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	service := &mockUserService{}
	router := setupRouter()
	router.POST("/login", LoginHandler(service))

	req, _ := http.NewRequest("POST", "/login", strings.NewReader(`{invalid json}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperava status 400, obteve %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["error"] != "requisição inválida" {
		t.Errorf("esperava erro 'requisição inválida', obteve '%v'", response["error"])
	}
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	service := &mockUserService{
		authenticateErr: errors.New("usuário ou senha inválidos"),
	}

	router := setupRouter()
	router.POST("/login", LoginHandler(service))

	reqBody := `{"cpf":"98765432109","password":"wrongpassword"}`
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperava status 401, obteve %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["error"] != "usuário ou senha inválidos" {
		t.Errorf("esperava erro de credenciais inválidas, obteve '%v'", response["error"])
	}
}

func TestLoginHandler_InternalServerError(t *testing.T) {
	service := &mockUserService{
		authenticateErr: errors.New("database connection failed"),
	}

	router := setupRouter()
	router.POST("/login", LoginHandler(service))

	reqBody := `{"cpf":"98765432109","password":"senha123"}`
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperava status 500, obteve %d", w.Code)
	}
}
