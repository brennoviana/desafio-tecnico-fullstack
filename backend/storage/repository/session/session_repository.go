package session

import (
	"database/sql"
	"desafio-tecnico-fullstack/backend/models"
)

type SessionRepository interface {
	OpenSession(topicID int, openAt, closeAt int64) error
	GetSessionByTopic(topicID int) (*models.Session, error)
}

type sessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) OpenSession(topicID int, openAt, closeAt int64) error {
	_, err := r.db.Exec("INSERT INTO sessions (topic_id, open_at, close_at) VALUES ($1, $2, $3)", topicID, openAt, closeAt)
	return err
}

func (r *sessionRepository) GetSessionByTopic(topicID int) (*models.Session, error) {
	var s models.Session
	err := r.db.QueryRow("SELECT id, topic_id, open_at, close_at FROM sessions WHERE topic_id = $1", topicID).Scan(&s.ID, &s.TopicID, &s.OpenAt, &s.CloseAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
