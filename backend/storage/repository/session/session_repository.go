package session

import (
	"database/sql"
	"desafio-tecnico-fullstack/backend/models"
	"time"
)

type SessionRepository interface {
	OpenSession(topicID int, openAt, closeAt int64) error
	GetSessionByTopic(topicID int) (*models.Session, error)
	UpdateExpiredSessions() error
}

type sessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) OpenSession(topicID int, openAt, closeAt int64) error {
	_, err := r.db.Exec("INSERT INTO sessions (topic_id, open_at, close_at) VALUES ($1, $2, $3)", topicID, openAt, closeAt)
	if err != nil {
		return err
	}

	_, err = r.db.Exec("UPDATE topics SET status = 'Sessão Aberta' WHERE id = $1", topicID)
	if err != nil {
		return err
	}

	return nil
}

func (r *sessionRepository) GetSessionByTopic(topicID int) (*models.Session, error) {
	var s models.Session
	err := r.db.QueryRow("SELECT id, topic_id, open_at, close_at FROM sessions WHERE topic_id = $1", topicID).Scan(&s.ID, &s.TopicID, &s.OpenAt, &s.CloseAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *sessionRepository) UpdateExpiredSessions() error {
	now := time.Now().Unix()

	_, err := r.db.Exec(`
		UPDATE topics 
		SET status = 'Votação Encerrada' 
		WHERE id IN (
			SELECT topic_id 
			FROM sessions 
			WHERE close_at < $1
		) AND status = 'Sessão Aberta'
	`, now)

	return err
}