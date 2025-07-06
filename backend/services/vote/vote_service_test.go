package vote

import (
	"desafio-tecnico-fullstack/backend/models"
	"errors"
	"testing"
	"time"
)

type mockVoteRepo struct {
	votes       []models.Vote
	hasVoted    bool
	hasVotedErr error
	registerErr error
	resultYes   int
	resultNo    int
	resultErr   error
}

func (m *mockVoteRepo) RegisterVote(vote models.Vote) error {
	if m.registerErr != nil {
		return m.registerErr
	}
	m.votes = append(m.votes, vote)
	return nil
}

func (m *mockVoteRepo) HasUserVoted(topicID int, userID int) (bool, error) {
	if m.hasVotedErr != nil {
		return false, m.hasVotedErr
	}
	return m.hasVoted, nil
}

func (m *mockVoteRepo) GetResult(topicID int) (yes int, no int, err error) {
	return m.resultYes, m.resultNo, m.resultErr
}

type mockSessionRepo struct {
	session    *models.Session
	sessionErr error
}

func (m *mockSessionRepo) OpenSession(topicID int, openAt, closeAt int64) error {
	return nil
}

func (m *mockSessionRepo) GetSessionByTopic(topicID int) (*models.Session, error) {
	if m.sessionErr != nil {
		return nil, m.sessionErr
	}
	return m.session, nil
}

func (m *mockSessionRepo) UpdateExpiredSessions() error {
	return nil
}

func TestVoteService_Vote_Success(t *testing.T) {
	now := time.Now().Unix()
	voteRepo := &mockVoteRepo{hasVoted: false}
	sessionRepo := &mockSessionRepo{
		session: &models.Session{
			ID:      1,
			TopicID: 1,
			OpenAt:  now - 100,
			CloseAt: now + 100,
		},
	}

	service := NewVoteService(voteRepo, sessionRepo)

	err := service.Vote(1, 123, "Sim")
	if err != nil {
		t.Errorf("esperava sucesso, obteve erro: %v", err)
	}

	if len(voteRepo.votes) != 1 {
		t.Errorf("esperava 1 voto registrado, obteve %d", len(voteRepo.votes))
	}

	vote := voteRepo.votes[0]
	if vote.TopicID != 1 || vote.UserID != 123 || vote.Choice != "Sim" {
		t.Errorf("voto registrado incorretamente: %+v", vote)
	}
}

func TestVoteService_Vote_InvalidChoice(t *testing.T) {
	voteRepo := &mockVoteRepo{}
	sessionRepo := &mockSessionRepo{}

	service := NewVoteService(voteRepo, sessionRepo)

	err := service.Vote(1, 123, "Talvez")
	if err == nil || err.Error() != "voto deve ser 'Sim' ou 'Não'" {
		t.Errorf("esperava erro de escolha inválida, obteve: %v", err)
	}
}

func TestVoteService_Vote_SessionNotFound(t *testing.T) {
	voteRepo := &mockVoteRepo{}
	sessionRepo := &mockSessionRepo{
		sessionErr: errors.New("session not found"),
	}

	service := NewVoteService(voteRepo, sessionRepo)

	err := service.Vote(1, 123, "Sim")
	if err == nil || err.Error() != "sessão não encontrada para a pauta" {
		t.Errorf("esperava erro de sessão não encontrada, obteve: %v", err)
	}
}

func TestVoteService_Vote_SessionClosed(t *testing.T) {
	now := time.Now().Unix()
	voteRepo := &mockVoteRepo{}
	sessionRepo := &mockSessionRepo{
		session: &models.Session{
			ID:      1,
			TopicID: 1,
			OpenAt:  now - 200,
			CloseAt: now - 100, // Sessão já fechada
		},
	}

	service := NewVoteService(voteRepo, sessionRepo)

	err := service.Vote(1, 123, "Sim")
	if err == nil || err.Error() != "sessão de votação não está aberta" {
		t.Errorf("esperava erro de sessão fechada, obteve: %v", err)
	}
}

func TestVoteService_Vote_SessionNotOpen(t *testing.T) {
	now := time.Now().Unix()
	voteRepo := &mockVoteRepo{}
	sessionRepo := &mockSessionRepo{
		session: &models.Session{
			ID:      1,
			TopicID: 1,
			OpenAt:  now + 100, // Sessão ainda não abriu
			CloseAt: now + 200,
		},
	}

	service := NewVoteService(voteRepo, sessionRepo)

	err := service.Vote(1, 123, "Sim")
	if err == nil || err.Error() != "sessão de votação não está aberta" {
		t.Errorf("esperava erro de sessão não aberta, obteve: %v", err)
	}
}

func TestVoteService_Vote_AlreadyVoted(t *testing.T) {
	now := time.Now().Unix()
	voteRepo := &mockVoteRepo{hasVoted: true}
	sessionRepo := &mockSessionRepo{
		session: &models.Session{
			ID:      1,
			TopicID: 1,
			OpenAt:  now - 100,
			CloseAt: now + 100,
		},
	}

	service := NewVoteService(voteRepo, sessionRepo)

	err := service.Vote(1, 123, "Sim")
	if err == nil || err.Error() != "voto já registrado" {
		t.Errorf("esperava erro de voto já registrado, obteve: %v", err)
	}
}

func TestVoteService_Vote_HasVotedError(t *testing.T) {
	now := time.Now().Unix()
	voteRepo := &mockVoteRepo{
		hasVotedErr: errors.New("database error"),
	}
	sessionRepo := &mockSessionRepo{
		session: &models.Session{
			ID:      1,
			TopicID: 1,
			OpenAt:  now - 100,
			CloseAt: now + 100,
		},
	}

	service := NewVoteService(voteRepo, sessionRepo)

	err := service.Vote(1, 123, "Sim")
	if err == nil || err.Error() != "database error" {
		t.Errorf("esperava erro de banco de dados, obteve: %v", err)
	}
}

func TestVoteService_GetResult_Success(t *testing.T) {
	voteRepo := &mockVoteRepo{
		resultYes: 10,
		resultNo:  5,
	}
	sessionRepo := &mockSessionRepo{}

	service := NewVoteService(voteRepo, sessionRepo)

	yes, no, err := service.GetResult(1)
	if err != nil {
		t.Errorf("esperava sucesso, obteve erro: %v", err)
	}

	if yes != 10 || no != 5 {
		t.Errorf("esperava yes=10, no=5, obteve yes=%d, no=%d", yes, no)
	}
}

func TestVoteService_GetResult_Error(t *testing.T) {
	voteRepo := &mockVoteRepo{
		resultErr: errors.New("database error"),
	}
	sessionRepo := &mockSessionRepo{}

	service := NewVoteService(voteRepo, sessionRepo)

	_, _, err := service.GetResult(1)
	if err == nil || err.Error() != "database error" {
		t.Errorf("esperava erro de banco de dados, obteve: %v", err)
	}
}
