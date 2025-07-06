package topic

import (
	"desafio-tecnico-fullstack/backend/models"
	"errors"
	"testing"
)

type mockTopicRepo struct {
	topics    []models.Topic
	createErr error
	listErr   error
}

func (m *mockTopicRepo) CreateTopic(topic models.Topic) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.topics = append(m.topics, topic)
	return nil
}

func (m *mockTopicRepo) ListTopics() ([]models.Topic, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.topics, nil
}

type mockSessionService struct {
	updateErr error
}

func (m *mockSessionService) OpenSession(topicID int, durationMinutes int) error {
	return nil
}

func (m *mockSessionService) GetSessionByTopic(topicID int) (*models.Session, error) {
	return nil, nil
}

func (m *mockSessionService) UpdateExpiredSessions() error {
	return m.updateErr
}

func TestTopicService_CreateTopic_Success(t *testing.T) {
	repo := &mockTopicRepo{}
	sessionService := &mockSessionService{}

	service := NewTopicService(repo, sessionService)

	err := service.CreateTopic("Nova Pauta", "Ativa")
	if err != nil {
		t.Errorf("esperava sucesso, obteve erro: %v", err)
	}

	if len(repo.topics) != 1 {
		t.Errorf("esperava 1 tópico criado, obteve %d", len(repo.topics))
	}

	topic := repo.topics[0]
	if topic.Name != "Nova Pauta" || topic.Status != "Ativa" {
		t.Errorf("tópico criado incorretamente: %+v", topic)
	}
}

func TestTopicService_CreateTopic_EmptyStatus(t *testing.T) {
	repo := &mockTopicRepo{}
	sessionService := &mockSessionService{}

	service := NewTopicService(repo, sessionService)

	err := service.CreateTopic("Nova Pauta", "")
	if err != nil {
		t.Errorf("esperava sucesso, obteve erro: %v", err)
	}

	if len(repo.topics) != 1 {
		t.Errorf("esperava 1 tópico criado, obteve %d", len(repo.topics))
	}

	topic := repo.topics[0]
	if topic.Status != "Aguardando Abertura" {
		t.Errorf("esperava status padrão 'Aguardando Abertura', obteve '%s'", topic.Status)
	}
}

func TestTopicService_CreateTopic_RepoError(t *testing.T) {
	repo := &mockTopicRepo{
		createErr: errors.New("database error"),
	}
	sessionService := &mockSessionService{}

	service := NewTopicService(repo, sessionService)

	err := service.CreateTopic("Nova Pauta", "Ativa")
	if err == nil || err.Error() != "database error" {
		t.Errorf("esperava erro de banco de dados, obteve: %v", err)
	}
}

func TestTopicService_ListTopics_Success(t *testing.T) {
	expectedTopics := []models.Topic{
		{ID: 1, Name: "Pauta 1", Status: "Ativa"},
		{ID: 2, Name: "Pauta 2", Status: "Inativa"},
	}

	repo := &mockTopicRepo{
		topics: expectedTopics,
	}
	sessionService := &mockSessionService{}

	service := NewTopicService(repo, sessionService)

	topics, err := service.ListTopics()
	if err != nil {
		t.Errorf("esperava sucesso, obteve erro: %v", err)
	}

	if len(topics) != 2 {
		t.Errorf("esperava 2 tópicos, obteve %d", len(topics))
	}

	for i, topic := range topics {
		if topic.ID != expectedTopics[i].ID ||
			topic.Name != expectedTopics[i].Name ||
			topic.Status != expectedTopics[i].Status {
			t.Errorf("tópico %d incorreto: esperava %+v, obteve %+v", i, expectedTopics[i], topic)
		}
	}
}

func TestTopicService_ListTopics_UpdateSessionError(t *testing.T) {
	repo := &mockTopicRepo{
		topics: []models.Topic{{ID: 1, Name: "Pauta 1", Status: "Ativa"}},
	}
	sessionService := &mockSessionService{
		updateErr: errors.New("session update error"),
	}

	service := NewTopicService(repo, sessionService)

	topics, err := service.ListTopics()
	if err == nil || err.Error() != "session update error" {
		t.Errorf("esperava erro de atualização de sessão, obteve: %v", err)
	}

	if topics != nil {
		t.Errorf("esperava topics nil quando há erro, obteve: %+v", topics)
	}
}

func TestTopicService_ListTopics_RepoError(t *testing.T) {
	repo := &mockTopicRepo{
		listErr: errors.New("database error"),
	}
	sessionService := &mockSessionService{}

	service := NewTopicService(repo, sessionService)

	topics, err := service.ListTopics()
	if err == nil || err.Error() != "database error" {
		t.Errorf("esperava erro de banco de dados, obteve: %v", err)
	}

	if topics != nil {
		t.Errorf("esperava topics nil quando há erro, obteve: %+v", topics)
	}
}

func TestTopicService_ListTopics_EmptyList(t *testing.T) {
	repo := &mockTopicRepo{
		topics: []models.Topic{},
	}
	sessionService := &mockSessionService{}

	service := NewTopicService(repo, sessionService)

	topics, err := service.ListTopics()
	if err != nil {
		t.Errorf("esperava sucesso, obteve erro: %v", err)
	}

	if len(topics) != 0 {
		t.Errorf("esperava lista vazia, obteve %d tópicos", len(topics))
	}
}
