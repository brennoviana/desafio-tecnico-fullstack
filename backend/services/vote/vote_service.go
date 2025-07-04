package vote

import (
	"desafio-tecnico-fullstack/backend/models"
	sessionRepoPkg "desafio-tecnico-fullstack/backend/storage/repository/session"
	voteRepoPkg "desafio-tecnico-fullstack/backend/storage/repository/vote"
	"errors"
	"time"
)

type VoteService interface {
	Vote(topicID int, userCPF, choice string) error
	GetResult(topicID int) (yes int, no int, err error)
}

type voteService struct {
	voteRepo    voteRepoPkg.VoteRepository
	sessionRepo sessionRepoPkg.SessionRepository
}

func NewVoteService(voteRepo voteRepoPkg.VoteRepository, sessionRepo sessionRepoPkg.SessionRepository) VoteService {
	return &voteService{voteRepo: voteRepo, sessionRepo: sessionRepo}
}

func (s *voteService) Vote(topicID int, userCPF, choice string) error {
	if choice != "Sim" && choice != "Não" {
		return errors.New("voto deve ser 'Sim' ou 'Não'")
	}
	session, err := s.sessionRepo.GetSessionByTopic(topicID)
	if err != nil {
		return errors.New("sessão não encontrada para a pauta")
	}
	now := time.Now().Unix()
	if now < session.OpenAt || now > session.CloseAt {
		return errors.New("sessão de votação não está aberta")
	}
	voted, err := s.voteRepo.HasUserVoted(topicID, userCPF)
	if err != nil {
		return err
	}
	if voted {
		return errors.New("usuário já votou nesta pauta")
	}
	vote := models.Vote{TopicID: topicID, UserCPF: userCPF, Choice: choice}
	return s.voteRepo.RegisterVote(vote)
}

func (s *voteService) GetResult(topicID int) (yes int, no int, err error) {
	return s.voteRepo.GetResult(topicID)
}
