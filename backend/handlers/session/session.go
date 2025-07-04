package session

import (
	sessionrepo "desafio-tecnico-fullstack/backend/storage/repository/session"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func OpenSessionHandler(repo sessionrepo.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		topicID, err := strconv.Atoi(c.Param("topic_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Erro": "topic_id inválido"})
			return
		}
		var req struct {
			DurationMinutes int `json:"duration_minutes"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			req.DurationMinutes = 1 // padrão
		}
		if req.DurationMinutes <= 0 {
			req.DurationMinutes = 1
		}
		now := time.Now().Unix()
		closeAt := now + int64(req.DurationMinutes*60)
		err = repo.OpenSession(topicID, now, closeAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Erro": err.Error()})
			return
		}
		c.Status(http.StatusCreated)
	}
}
