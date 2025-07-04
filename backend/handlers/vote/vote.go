package vote

import (
	"desafio-tecnico-fullstack/backend/models"
	sessionRepoPkg "desafio-tecnico-fullstack/backend/storage/repository/session"
	voteRepoPkg "desafio-tecnico-fullstack/backend/storage/repository/vote"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func VoteHandler(voteRepo voteRepoPkg.VoteRepository, sessionRepo sessionRepoPkg.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		topicID, err := strconv.Atoi(c.Param("topic_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Erro": "topic_id inválido"})
			return
		}
		var req struct {
			Choice string `json:"choice"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || (req.Choice != "Sim" && req.Choice != "Não") {
			c.JSON(http.StatusBadRequest, gin.H{"Erro": "Voto deve ser 'Sim' ou 'Não'"})
			return
		}
		cpf := c.GetString("cpf")

		session, err := sessionRepo.GetSessionByTopic(topicID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Erro": "Sessão não encontrada para a pauta"})
			return
		}
		now := time.Now().Unix()
		if now < session.OpenAt || now > session.CloseAt {
			c.JSON(http.StatusBadRequest, gin.H{"Erro": "Sessão de votação não está aberta"})
			return
		}

		voted, err := voteRepo.HasUserVoted(topicID, cpf)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Erro": err.Error()})
			return
		}
		if voted {
			c.JSON(http.StatusConflict, gin.H{"Erro": "Usuário já votou nesta pauta"})
			return
		}
		vote := models.Vote{TopicID: topicID, UserCPF: cpf, Choice: req.Choice}
		err = voteRepo.RegisterVote(vote)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Erro": err.Error()})
			return
		}
		c.Status(http.StatusCreated)
	}
}

func ResultHandler(voteRepo voteRepoPkg.VoteRepository, sessionRepo sessionRepoPkg.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		topicID, err := strconv.Atoi(c.Param("topic_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Erro": "topic_id inválido"})
			return
		}
		session, err := sessionRepo.GetSessionByTopic(topicID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Erro": "Sessão não encontrada para a pauta"})
			return
		}
		now := time.Now().Unix()
		if now < session.CloseAt {
			c.JSON(http.StatusBadRequest, gin.H{"Erro": "Sessão de votação ainda está aberta"})
			return
		}
		yes, no, err := voteRepo.GetResult(topicID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Erro": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Sim": yes, "Não": no})
	}
}
