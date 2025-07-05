package topic

import (
	"desafio-tecnico-fullstack/backend/models"
	"desafio-tecnico-fullstack/backend/services/session"
	topicrepo "desafio-tecnico-fullstack/backend/storage/repository/topic"
)

type TopicService interface {
	CreateTopic(name string, status string) error
	ListTopics() ([]models.Topic, error)
}

type topicService struct {
	repo           topicrepo.TopicRepository
	sessionService session.SessionService
}

func NewTopicService(repo topicrepo.TopicRepository, sessionService session.SessionService) TopicService {
	return &topicService{
		repo:           repo,
		sessionService: sessionService,
	}
}

func (s *topicService) CreateTopic(name string, status string) error {

	if status == "" {
		status = "Aguardando Abertura"
	}

	topic := models.Topic{Name: name, Status: status}
	return s.repo.CreateTopic(topic)
}

func (s *topicService) ListTopics() ([]models.Topic, error) {
	if err := s.sessionService.UpdateExpiredSessions(); err != nil {
		return nil, err
	}

	return s.repo.ListTopics()
}
