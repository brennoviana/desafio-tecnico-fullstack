package topic

import (
	"desafio-tecnico-fullstack/backend/models"
	topicrepo "desafio-tecnico-fullstack/backend/storage/repository/topic"
	"errors"
)

type TopicService interface {
	CreateTopic(name string, status ...string) error
	ListTopics() ([]models.Topic, error)
}

type topicService struct {
	repo topicrepo.TopicRepository
}

func NewTopicService(repo topicrepo.TopicRepository) TopicService {
	return &topicService{repo: repo}
}

func (s *topicService) CreateTopic(name string, status ...string) error {
	if len(name) < 3 {
		return errors.New("nome da pauta muito curto")
	}

	topicStatus := "PENDING"
	if len(status) > 0 && status[0] != "" {
		topicStatus = status[0]
	}

	topic := models.Topic{Name: name, Status: topicStatus}
	return s.repo.CreateTopic(topic)
}

func (s *topicService) ListTopics() ([]models.Topic, error) {
	return s.repo.ListTopics()
}
