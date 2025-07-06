package topic

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

type mockTopicService struct {
	createErr error
	topics    []models.Topic
	listErr   error
}

func (m *mockTopicService) CreateTopic(name, status string) error {
	return m.createErr
}

func (m *mockTopicService) ListTopics() ([]models.Topic, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.topics, nil
}

func setupTopicRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestCreateTopicHandler_Success(t *testing.T) {
	service := &mockTopicService{}
	router := setupTopicRouter()

	// Simulate authentication middleware
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 123)
		c.Next()
	})

	router.POST("/topics", CreateTopicHandler(service))

	reqBody := `{"name":"Nova Pauta de Votação","status":"Ativa"}`
	req, _ := http.NewRequest("POST", "/topics", strings.NewReader(reqBody))
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
}

func TestCreateTopicHandler_SuccessWithoutStatus(t *testing.T) {
	service := &mockTopicService{}
	router := setupTopicRouter()

	router.Use(func(c *gin.Context) {
		c.Set("user_id", 123)
		c.Next()
	})

	router.POST("/topics", CreateTopicHandler(service))

	reqBody := `{"name":"Nova Pauta de Votação"}`
	req, _ := http.NewRequest("POST", "/topics", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperava status 200, obteve %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "success" {
		t.Errorf("esperava status 'success', obteve '%v'", response["status"])
	}
}

func TestCreateTopicHandler_InvalidJSON(t *testing.T) {
	service := &mockTopicService{}
	router := setupTopicRouter()

	router.Use(func(c *gin.Context) {
		c.Set("user_id", 123)
		c.Next()
	})

	router.POST("/topics", CreateTopicHandler(service))

	req, _ := http.NewRequest("POST", "/topics", strings.NewReader(`{invalid json}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperava status 400, obteve %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["error"] != "Requisição inválida" {
		t.Errorf("esperava erro 'Requisição inválida', obteve '%v'", response["error"])
	}
}

func TestCreateTopicHandler_MissingName(t *testing.T) {
	testCases := []struct {
		name    string
		reqBody string
	}{
		{"empty name", `{"name":"","status":"Ativa"}`},
		{"missing name", `{"status":"Ativa"}`},
		{"null name", `{"name":null,"status":"Ativa"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := &mockTopicService{}
			router := setupTopicRouter()

			router.Use(func(c *gin.Context) {
				c.Set("user_id", 123)
				c.Next()
			})

			router.POST("/topics", CreateTopicHandler(service))

			req, _ := http.NewRequest("POST", "/topics", strings.NewReader(tc.reqBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("esperava status 400, obteve %d", w.Code)
			}

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if response["error"] != "Nome da pauta é obrigatório" {
				t.Errorf("esperava erro sobre nome obrigatório, obteve '%v'", response["error"])
			}
		})
	}
}

func TestCreateTopicHandler_ServiceError(t *testing.T) {
	service := &mockTopicService{
		createErr: errors.New("database connection failed"),
	}
	router := setupTopicRouter()

	router.Use(func(c *gin.Context) {
		c.Set("user_id", 123)
		c.Next()
	})

	router.POST("/topics", CreateTopicHandler(service))

	reqBody := `{"name":"Nova Pauta","status":"Ativa"}`
	req, _ := http.NewRequest("POST", "/topics", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperava status 400, obteve %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["error"] != "database connection failed" {
		t.Errorf("esperava erro de banco de dados, obteve '%v'", response["error"])
	}
}

func TestCreateTopicHandler_ValidNames(t *testing.T) {
	validNames := []string{
		"Pauta Simples",
		"Pauta com números 123",
		"Pauta-com-hífens",
		"Pauta_com_underscores",
		"Pauta com caracteres especiais: @#$%",
		"Nome muito longo que deveria ser aceito se não houver validação de tamanho específica no handler",
	}

	for _, name := range validNames {
		t.Run("name_"+name, func(t *testing.T) {
			service := &mockTopicService{}
			router := setupTopicRouter()

			router.Use(func(c *gin.Context) {
				c.Set("user_id", 123)
				c.Next()
			})

			router.POST("/topics", CreateTopicHandler(service))

			reqBody := `{"name":"` + name + `","status":"Ativa"}`
			req, _ := http.NewRequest("POST", "/topics", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("esperava status 200 para nome '%s', obteve %d", name, w.Code)
			}
		})
	}
}

func TestListTopicsHandler_Success(t *testing.T) {
	expectedTopics := []models.Topic{
		{ID: 1, Name: "Primeira Pauta", Status: "Ativa"},
		{ID: 2, Name: "Segunda Pauta", Status: "Inativa"},
		{ID: 3, Name: "Terceira Pauta", Status: "Aguardando Abertura"},
	}

	service := &mockTopicService{
		topics: expectedTopics,
	}
	router := setupTopicRouter()
	router.GET("/topics", ListTopicsHandler(service))

	req, _ := http.NewRequest("GET", "/topics", nil)

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

	data := response["data"].([]interface{})
	if len(data) != len(expectedTopics) {
		t.Errorf("esperava %d tópicos, obteve %d", len(expectedTopics), len(data))
	}

	// Verificar primeiro tópico
	if len(data) > 0 {
		firstTopic := data[0].(map[string]interface{})
		if int(firstTopic["id"].(float64)) != expectedTopics[0].ID {
			t.Errorf("esperava ID %d no primeiro tópico, obteve %v", expectedTopics[0].ID, firstTopic["id"])
		}
		if firstTopic["name"] != expectedTopics[0].Name {
			t.Errorf("esperava nome '%s' no primeiro tópico, obteve '%v'", expectedTopics[0].Name, firstTopic["name"])
		}
		if firstTopic["status"] != expectedTopics[0].Status {
			t.Errorf("esperava status '%s' no primeiro tópico, obteve '%v'", expectedTopics[0].Status, firstTopic["status"])
		}
	}
}

func TestListTopicsHandler_EmptyList(t *testing.T) {
	service := &mockTopicService{
		topics: []models.Topic{},
	}
	router := setupTopicRouter()
	router.GET("/topics", ListTopicsHandler(service))

	req, _ := http.NewRequest("GET", "/topics", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperava status 200, obteve %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "success" {
		t.Errorf("esperava status 'success', obteve '%v'", response["status"])
	}

	data := response["data"].([]interface{})
	if len(data) != 0 {
		t.Errorf("esperava lista vazia, obteve %d tópicos", len(data))
	}
}

func TestListTopicsHandler_ServiceError(t *testing.T) {
	service := &mockTopicService{
		listErr: errors.New("database connection failed"),
	}
	router := setupTopicRouter()
	router.GET("/topics", ListTopicsHandler(service))

	req, _ := http.NewRequest("GET", "/topics", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperava status 500, obteve %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["error"] != "database connection failed" {
		t.Errorf("esperava erro de banco de dados, obteve '%v'", response["error"])
	}
}

func TestListTopicsHandler_NilTopics(t *testing.T) {
	service := &mockTopicService{
		topics: nil,
	}
	router := setupTopicRouter()
	router.GET("/topics", ListTopicsHandler(service))

	req, _ := http.NewRequest("GET", "/topics", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperava status 200, obteve %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "success" {
		t.Errorf("esperava status 'success', obteve '%v'", response["status"])
	}

	// data should be nil for nil slice
	if response["data"] != nil {
		t.Errorf("esperava data nil, obteve %v", response["data"])
	}
}
