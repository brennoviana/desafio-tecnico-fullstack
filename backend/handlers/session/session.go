package session

import (
	"desafio-tecnico-fullstack/backend/services/session"
	"desafio-tecnico-fullstack/backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func OpenSessionHandler(sessionService session.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		topicID, err := strconv.Atoi(c.Param("topic_id"))
		if err != nil {
			utils.RespondError(c, http.StatusBadRequest, "topic_id inválido")
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
			utils.RespondError(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondSuccess(c, nil)
	}
}
