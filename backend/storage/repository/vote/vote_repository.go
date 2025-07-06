package vote

import (
	"database/sql"
	"desafio-tecnico-fullstack/backend/models"
)

type VoteRepository interface {
	RegisterVote(vote models.Vote) error
	HasUserVoted(topicID int, userID int) (bool, error)
	GetResult(topicID int) (yes int, no int, err error)
}

type voteRepository struct {
	db *sql.DB
}

func NewVoteRepository(db *sql.DB) VoteRepository {
	return &voteRepository{db: db}
}

func (r *voteRepository) RegisterVote(vote models.Vote) error {
	_, err := r.db.Exec("INSERT INTO votes (topic_id, user_id, choice) VALUES ($1, $2, $3)", vote.TopicID, vote.UserID, vote.Choice)
	return err
}

func (r *voteRepository) HasUserVoted(topicID int, userID int) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM votes WHERE topic_id = $1 AND user_id = $2", topicID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *voteRepository) GetResult(topicID int) (yes int, no int, err error) {
	err = r.db.QueryRow("SELECT COUNT(*) FROM votes WHERE topic_id = $1 AND choice = 'Sim'", topicID).Scan(&yes)
	if err != nil {
		return
	}

	err = r.db.QueryRow("SELECT COUNT(*) FROM votes WHERE topic_id = $1 AND choice = 'Não'", topicID).Scan(&no)
	return
}
