package session

import (
	"desafio-tecnico-fullstack/backend/models"
	sessionrepo "desafio-tecnico-fullstack/backend/storage/repository/session"
	"time"
)

type SessionService interface {
	OpenSession(topicID int, durationMinutes int) error
	GetSessionByTopic(topicID int) (*models.Session, error)
	UpdateExpiredSessions() error
}

type sessionService struct {
	repo sessionrepo.SessionRepository
}

func NewSessionService(repo sessionrepo.SessionRepository) SessionService {
	return &sessionService{repo: repo}
}

func (s *sessionService) OpenSession(topicID int, durationMinutes int) error {
	if durationMinutes <= 0 {
		durationMinutes = 1
	}

	now := time.Now().Unix()
	closeAt := now + int64(durationMinutes*60)
	return s.repo.OpenSession(topicID, now, closeAt)
}

func (s *sessionService) GetSessionByTopic(topicID int) (*models.Session, error) {
	return s.repo.GetSessionByTopic(topicID)
}

func (s *sessionService) UpdateExpiredSessions() error {
	return s.repo.UpdateExpiredSessions()
}
