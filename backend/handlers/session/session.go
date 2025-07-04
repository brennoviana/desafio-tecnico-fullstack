package session

import (
	"desafio-tecnico-fullstack/backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func OpenSessionHandler(sessionService services.SessionService) gin.HandlerFunc {
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
			req.DurationMinutes = 1
		}
		err = sessionService.OpenSession(topicID, req.DurationMinutes)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Erro": err.Error()})
			return
		}
		c.Status(http.StatusCreated)
	}
}
