package topic

import (
	"database/sql"
	"desafio-tecnico-fullstack/backend/models"
)

type TopicRepository interface {
	CreateTopic(topic models.Topic) error
	ListTopics() ([]models.Topic, error)
}

type topicRepository struct {
	db *sql.DB
}

func NewTopicRepository(db *sql.DB) TopicRepository {
	return &topicRepository{db: db}
}

func (r *topicRepository) CreateTopic(topic models.Topic) error {
	_, err := r.db.Exec("INSERT INTO topics (name) VALUES ($1)", topic.Name)
	return err
}

func (r *topicRepository) ListTopics() ([]models.Topic, error) {
	rows, err := r.db.Query("SELECT id, name FROM topics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	topics := []models.Topic{}
	for rows.Next() {
		var t models.Topic
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		topics = append(topics, t)
	}
	return topics, nil
}
