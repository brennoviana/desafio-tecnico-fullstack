package session

import (
	"desafio-tecnico-fullstack/backend/models"
	"errors"
	"testing"
	"time"
)

type mockSessionRepo struct {
	openErr     error
	getSession  *models.Session
	getErr      error
	updateErr   error
	openedCalls []openSessionCall
}

type openSessionCall struct {
	topicID int
	openAt  int64
	closeAt int64
}

func (m *mockSessionRepo) OpenSession(topicID int, openAt, closeAt int64) error {
	if m.openErr != nil {
		return m.openErr
	}
	m.openedCalls = append(m.openedCalls, openSessionCall{
		topicID: topicID,
		openAt:  openAt,
		closeAt: closeAt,
	})
	return nil
}

func (m *mockSessionRepo) GetSessionByTopic(topicID int) (*models.Session, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.getSession, nil
}

func (m *mockSessionRepo) UpdateExpiredSessions() error {
	return m.updateErr
}

func TestSessionService_OpenSession_Success(t *testing.T) {
	repo := &mockSessionRepo{}
	service := NewSessionService(repo)

	now := time.Now().Unix()

	err := service.OpenSession(1, 5)
	if err != nil {
		t.Errorf("esperava sucesso, obteve erro: %v", err)
	}

	if len(repo.openedCalls) != 1 {
		t.Errorf("esperava 1 chamada para OpenSession, obteve %d", len(repo.openedCalls))
	}

	call := repo.openedCalls[0]
	if call.topicID != 1 {
		t.Errorf("esperava topicID 1, obteve %d", call.topicID)
	}

	expectedCloseAt := call.openAt + int64(5*60)
	if call.closeAt != expectedCloseAt {
		t.Errorf("esperava closeAt %d, obteve %d", expectedCloseAt, call.closeAt)
	}

	// Verifica se openAt está próximo do tempo atual (dentro de 5 segundos)
	if abs(call.openAt-now) > 5 {
		t.Errorf("openAt muito diferente do tempo atual: esperava próximo de %d, obteve %d", now, call.openAt)
	}
}

func TestSessionService_OpenSession_ZeroDuration(t *testing.T) {
	repo := &mockSessionRepo{}
	service := NewSessionService(repo)

	err := service.OpenSession(1, 0)
	if err != nil {
		t.Errorf("esperava sucesso, obteve erro: %v", err)
	}

	if len(repo.openedCalls) != 1 {
		t.Errorf("esperava 1 chamada para OpenSession, obteve %d", len(repo.openedCalls))
	}

	call := repo.openedCalls[0]
	expectedCloseAt := call.openAt + int64(1*60) // Deve usar 1 minuto como padrão
	if call.closeAt != expectedCloseAt {
		t.Errorf("esperava closeAt %d (1 minuto), obteve %d", expectedCloseAt, call.closeAt)
	}
}

func TestSessionService_OpenSession_NegativeDuration(t *testing.T) {
	repo := &mockSessionRepo{}
	service := NewSessionService(repo)

	err := service.OpenSession(1, -10)
	if err != nil {
		t.Errorf("esperava sucesso, obteve erro: %v", err)
	}

	if len(repo.openedCalls) != 1 {
		t.Errorf("esperava 1 chamada para OpenSession, obteve %d", len(repo.openedCalls))
	}

	call := repo.openedCalls[0]
	expectedCloseAt := call.openAt + int64(1*60) // Deve usar 1 minuto como padrão
	if call.closeAt != expectedCloseAt {
		t.Errorf("esperava closeAt %d (1 minuto), obteve %d", expectedCloseAt, call.closeAt)
	}
}

func TestSessionService_OpenSession_RepoError(t *testing.T) {
	repo := &mockSessionRepo{
		openErr: errors.New("database error"),
	}
	service := NewSessionService(repo)

	err := service.OpenSession(1, 5)
	if err == nil || err.Error() != "database error" {
		t.Errorf("esperava erro de banco de dados, obteve: %v", err)
	}
}

func TestSessionService_GetSessionByTopic_Success(t *testing.T) {
	expectedSession := &models.Session{
		ID:      1,
		TopicID: 123,
		OpenAt:  1000,
		CloseAt: 2000,
	}

	repo := &mockSessionRepo{
		getSession: expectedSession,
	}
	service := NewSessionService(repo)

	session, err := service.GetSessionByTopic(123)
	if err != nil {
		t.Errorf("esperava sucesso, obteve erro: %v", err)
	}

	if session == nil {
		t.Fatalf("esperava sessão, obteve nil")
	}

	if session.ID != expectedSession.ID ||
		session.TopicID != expectedSession.TopicID ||
		session.OpenAt != expectedSession.OpenAt ||
		session.CloseAt != expectedSession.CloseAt {
		t.Errorf("sessão incorreta: esperava %+v, obteve %+v", expectedSession, session)
	}
}

func TestSessionService_GetSessionByTopic_NotFound(t *testing.T) {
	repo := &mockSessionRepo{
		getSession: nil,
		getErr:     nil,
	}
	service := NewSessionService(repo)

	session, err := service.GetSessionByTopic(123)
	if err != nil {
		t.Errorf("esperava sucesso, obteve erro: %v", err)
	}

	if session != nil {
		t.Errorf("esperava nil, obteve sessão: %+v", session)
	}
}

func TestSessionService_GetSessionByTopic_RepoError(t *testing.T) {
	repo := &mockSessionRepo{
		getErr: errors.New("database error"),
	}
	service := NewSessionService(repo)

	session, err := service.GetSessionByTopic(123)
	if err == nil || err.Error() != "database error" {
		t.Errorf("esperava erro de banco de dados, obteve: %v", err)
	}

	if session != nil {
		t.Errorf("esperava nil quando há erro, obteve: %+v", session)
	}
}

func TestSessionService_UpdateExpiredSessions_Success(t *testing.T) {
	repo := &mockSessionRepo{}
	service := NewSessionService(repo)

	err := service.UpdateExpiredSessions()
	if err != nil {
		t.Errorf("esperava sucesso, obteve erro: %v", err)
	}
}

func TestSessionService_UpdateExpiredSessions_RepoError(t *testing.T) {
	repo := &mockSessionRepo{
		updateErr: errors.New("database error"),
	}
	service := NewSessionService(repo)

	err := service.UpdateExpiredSessions()
	if err == nil || err.Error() != "database error" {
		t.Errorf("esperava erro de banco de dados, obteve: %v", err)
	}
}

func TestSessionService_OpenSession_DurationCalculation(t *testing.T) {
	tests := []struct {
		name             string
		inputDuration    int
		expectedDuration int64
	}{
		{"1 minute", 1, 60},
		{"5 minutes", 5, 300},
		{"10 minutes", 10, 600},
		{"60 minutes", 60, 3600},
		{"120 minutes", 120, 7200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockSessionRepo{}
			service := NewSessionService(repo)

			err := service.OpenSession(1, tt.inputDuration)
			if err != nil {
				t.Errorf("esperava sucesso, obteve erro: %v", err)
			}

			if len(repo.openedCalls) != 1 {
				t.Errorf("esperava 1 chamada para OpenSession, obteve %d", len(repo.openedCalls))
			}

			call := repo.openedCalls[0]
			actualDuration := call.closeAt - call.openAt
			if actualDuration != tt.expectedDuration {
				t.Errorf("duração incorreta: esperava %d segundos, obteve %d", tt.expectedDuration, actualDuration)
			}
		})
	}
}

// Helper function to calculate absolute value
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
