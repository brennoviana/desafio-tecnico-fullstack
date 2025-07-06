package vote

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockVoteService struct {
	voteErr   error
	resultYes int
	resultNo  int
	resultErr error
}

func (m *mockVoteService) Vote(topicID int, userID int, choice string) error {
	return m.voteErr
}

func (m *mockVoteService) GetResult(topicID int) (yes int, no int, err error) {
	return m.resultYes, m.resultNo, m.resultErr
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestVoteHandler_Success(t *testing.T) {
	service := &mockVoteService{}
	router := setupTestRouter()

	router.POST("/api/topics/:topic_id/vote", func(c *gin.Context) {
		c.Set("user_id", 123)
		VoteHandler(service)(c)
	})

	requestBody := map[string]string{"choice": "Sim"}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/topics/1/vote", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, "success", response["status"])
}

func TestVoteHandler_InvalidTopicID(t *testing.T) {
	service := &mockVoteService{}
	router := setupTestRouter()

	router.POST("/api/topics/:topic_id/vote", VoteHandler(service))

	requestBody := map[string]string{"choice": "Sim"}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/topics/invalid/vote", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "topic_id inválido", response["error"])
}

func TestVoteHandler_InvalidJSON(t *testing.T) {
	service := &mockVoteService{}
	router := setupTestRouter()

	router.POST("/api/topics/:topic_id/vote", func(c *gin.Context) {
		c.Set("user_id", 123)
		VoteHandler(service)(c)
	})

	req, _ := http.NewRequest("POST", "/api/topics/1/vote", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "voto deve ser 'Sim' ou 'Não'", response["error"])
}

func TestVoteHandler_NoUserID(t *testing.T) {
	service := &mockVoteService{}
	router := setupTestRouter()

	router.POST("/api/topics/:topic_id/vote", VoteHandler(service))

	requestBody := map[string]string{"choice": "Sim"}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/topics/1/vote", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "usuário não autenticado", response["error"])
}

func TestVoteHandler_ServiceError(t *testing.T) {
	service := &mockVoteService{
		voteErr: errors.New("voto já registrado"),
	}
	router := setupTestRouter()

	router.POST("/api/topics/:topic_id/vote", func(c *gin.Context) {
		c.Set("user_id", 123)
		VoteHandler(service)(c)
	})

	requestBody := map[string]string{"choice": "Sim"}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/topics/1/vote", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "voto já registrado", response["error"])
}

func TestVoteHandler_DifferentChoices(t *testing.T) {
	tests := []struct {
		name   string
		choice string
	}{
		{"Sim vote", "Sim"},
		{"Não vote", "Não"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &mockVoteService{}
			router := setupTestRouter()

			router.POST("/api/topics/:topic_id/vote", func(c *gin.Context) {
				c.Set("user_id", 123)
				VoteHandler(service)(c)
			})

			requestBody := map[string]string{"choice": tt.choice}
			jsonBody, _ := json.Marshal(requestBody)

			req, _ := http.NewRequest("POST", "/api/topics/1/vote", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)

			var response map[string]interface{}
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Equal(t, "success", response["status"])
		})
	}
}

func TestResultHandler_Success(t *testing.T) {
	service := &mockVoteService{
		resultYes: 10,
		resultNo:  5,
	}
	router := setupTestRouter()

	router.GET("/api/topics/:topic_id/result", ResultHandler(service))

	req, _ := http.NewRequest("GET", "/api/topics/1/result", nil)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, "success", response["status"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(10), data["Sim"])
	assert.Equal(t, float64(5), data["Não"])
}

func TestResultHandler_InvalidTopicID(t *testing.T) {
	service := &mockVoteService{}
	router := setupTestRouter()

	router.GET("/api/topics/:topic_id/result", ResultHandler(service))

	req, _ := http.NewRequest("GET", "/api/topics/invalid/result", nil)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "topic_id inválido", response["error"])
}

func TestResultHandler_ServiceError(t *testing.T) {
	service := &mockVoteService{
		resultErr: errors.New("database error"),
	}
	router := setupTestRouter()

	router.GET("/api/topics/:topic_id/result", ResultHandler(service))

	req, _ := http.NewRequest("GET", "/api/topics/1/result", nil)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "database error", response["error"])
}

func TestResultHandler_ZeroResults(t *testing.T) {
	service := &mockVoteService{
		resultYes: 0,
		resultNo:  0,
	}
	router := setupTestRouter()

	router.GET("/api/topics/:topic_id/result", ResultHandler(service))

	req, _ := http.NewRequest("GET", "/api/topics/1/result", nil)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, "success", response["status"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(0), data["Sim"])
	assert.Equal(t, float64(0), data["Não"])
}

func TestResultHandler_MultipleTopicIDs(t *testing.T) {
	topicIDs := []string{"1", "100", "999"}

	for _, topicID := range topicIDs {
		t.Run("TopicID_"+topicID, func(t *testing.T) {
			service := &mockVoteService{
				resultYes: 3,
				resultNo:  7,
			}
			router := setupTestRouter()

			router.GET("/api/topics/:topic_id/result", ResultHandler(service))

			req, _ := http.NewRequest("GET", "/api/topics/"+topicID+"/result", nil)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)

			var response map[string]interface{}
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Equal(t, "success", response["status"])
		})
	}
}
