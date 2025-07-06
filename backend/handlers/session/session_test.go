package session

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"desafio-tecnico-fullstack/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockSessionService struct {
	openErr    error
	getSession *models.Session
	getErr     error
	updateErr  error
}

func (m *mockSessionService) OpenSession(topicID int, durationMinutes int) error {
	return m.openErr
}

func (m *mockSessionService) GetSessionByTopic(topicID int) (*models.Session, error) {
	return m.getSession, m.getErr
}

func (m *mockSessionService) UpdateExpiredSessions() error {
	return m.updateErr
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestOpenSessionHandler_Success(t *testing.T) {
	service := &mockSessionService{}
	router := setupTestRouter()

	router.POST("/api/topics/:topic_id/session", OpenSessionHandler(service))

	requestBody := map[string]int{"duration_minutes": 5}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/topics/1/session", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, "success", response["status"])
}

func TestOpenSessionHandler_InvalidTopicID(t *testing.T) {
	service := &mockSessionService{}
	router := setupTestRouter()

	router.POST("/api/topics/:topic_id/session", OpenSessionHandler(service))

	requestBody := map[string]int{"duration_minutes": 5}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/topics/invalid/session", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "topic_id inválido", response["error"])
}

func TestOpenSessionHandler_InvalidJSON(t *testing.T) {
	service := &mockSessionService{}
	router := setupTestRouter()

	router.POST("/api/topics/:topic_id/session", OpenSessionHandler(service))

	req, _ := http.NewRequest("POST", "/api/topics/1/session", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Should still succeed because the handler sets default duration when JSON bind fails
	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, "success", response["status"])
}

func TestOpenSessionHandler_EmptyBody(t *testing.T) {
	service := &mockSessionService{}
	router := setupTestRouter()

	router.POST("/api/topics/:topic_id/session", OpenSessionHandler(service))

	req, _ := http.NewRequest("POST", "/api/topics/1/session", nil)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Should succeed with default duration of 1 minute
	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, "success", response["status"])
}

func TestOpenSessionHandler_ServiceError(t *testing.T) {
	service := &mockSessionService{
		openErr: errors.New("sessão já existe para esta pauta"),
	}
	router := setupTestRouter()

	router.POST("/api/topics/:topic_id/session", OpenSessionHandler(service))

	requestBody := map[string]int{"duration_minutes": 5}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/topics/1/session", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "sessão já existe para esta pauta", response["error"])
}

func TestOpenSessionHandler_DifferentDurations(t *testing.T) {
	tests := []struct {
		name     string
		duration int
	}{
		{"1 minute", 1},
		{"5 minutes", 5},
		{"10 minutes", 10},
		{"30 minutes", 30},
		{"60 minutes", 60},
		{"zero duration", 0},
		{"negative duration", -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &mockSessionService{}
			router := setupTestRouter()

			router.POST("/api/topics/:topic_id/session", OpenSessionHandler(service))

			requestBody := map[string]int{"duration_minutes": tt.duration}
			jsonBody, _ := json.Marshal(requestBody)

			req, _ := http.NewRequest("POST", "/api/topics/1/session", bytes.NewBuffer(jsonBody))
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

func TestOpenSessionHandler_MultipleTopicIDs(t *testing.T) {
	topicIDs := []string{"1", "100", "999"}

	for _, topicID := range topicIDs {
		t.Run("TopicID_"+topicID, func(t *testing.T) {
			service := &mockSessionService{}
			router := setupTestRouter()

			router.POST("/api/topics/:topic_id/session", OpenSessionHandler(service))

			requestBody := map[string]int{"duration_minutes": 5}
			jsonBody, _ := json.Marshal(requestBody)

			req, _ := http.NewRequest("POST", "/api/topics/"+topicID+"/session", bytes.NewBuffer(jsonBody))
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

func TestOpenSessionHandler_MissingDurationField(t *testing.T) {
	service := &mockSessionService{}
	router := setupTestRouter()

	router.POST("/api/topics/:topic_id/session", OpenSessionHandler(service))

	// Send JSON without duration_minutes field
	requestBody := map[string]string{"other_field": "value"}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/topics/1/session", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Should still succeed because missing field defaults to 0, which is handled by service
	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, "success", response["status"])
}

func TestOpenSessionHandler_ValidJSON_ValidDuration(t *testing.T) {
	service := &mockSessionService{}
	router := setupTestRouter()

	router.POST("/api/topics/:topic_id/session", OpenSessionHandler(service))

	requestBody := map[string]int{"duration_minutes": 15}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/topics/42/session", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, "success", response["status"])
	assert.Nil(t, response["data"])
}

func TestOpenSessionHandler_HTTPMethods(t *testing.T) {
	service := &mockSessionService{}
	router := setupTestRouter()

	router.POST("/api/topics/:topic_id/session", OpenSessionHandler(service))

	// Test with GET method (should not match the route)
	req, _ := http.NewRequest("GET", "/api/topics/1/session", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)

	// Test with PUT method (should not match the route)
	req, _ = http.NewRequest("PUT", "/api/topics/1/session", nil)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestOpenSessionHandler_DatabaseErrors(t *testing.T) {
	dbErrors := []string{
		"connection timeout",
		"duplicate key violation",
		"foreign key constraint",
		"table not found",
	}

	for _, errorMsg := range dbErrors {
		t.Run("Error_"+errorMsg, func(t *testing.T) {
			service := &mockSessionService{
				openErr: errors.New(errorMsg),
			}
			router := setupTestRouter()

			router.POST("/api/topics/:topic_id/session", OpenSessionHandler(service))

			requestBody := map[string]int{"duration_minutes": 5}
			jsonBody, _ := json.Marshal(requestBody)

			req, _ := http.NewRequest("POST", "/api/topics/1/session", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)

			var response map[string]interface{}
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Equal(t, "error", response["status"])
			assert.Equal(t, errorMsg, response["error"])
		})
	}
}
