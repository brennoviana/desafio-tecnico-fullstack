package topic

import (
	"desafio-tecnico-fullstack/backend/models"
	topicrepo "desafio-tecnico-fullstack/backend/storage/repository/topic"
)

type TopicService interface {
	CreateTopic(name string, status string) error
	ListTopics() ([]models.Topic, error)
}

type topicService struct {
	repo topicrepo.TopicRepository
}

func NewTopicService(repo topicrepo.TopicRepository) TopicService {
	return &topicService{repo: repo}
}

func (s *topicService) CreateTopic(name string, status string) error {

	if status == "" {
		status = "Aguardando Abertura"
	}

	topic := models.Topic{Name: name, Status: status}
	return s.repo.CreateTopic(topic)
}

func (s *topicService) ListTopics() ([]models.Topic, error) {
	return s.repo.ListTopics()
}
