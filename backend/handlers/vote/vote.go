package vote

import (
	"desafio-tecnico-fullstack/backend/services/vote"
	"desafio-tecnico-fullstack/backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func VoteHandler(voteService vote.VoteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		topicID, err := strconv.Atoi(c.Param("topic_id"))
		if err != nil {
			utils.RespondError(c, http.StatusBadRequest, "topic_id inválido")
			return
		}

		var req struct {
			Choice string `json:"choice"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			utils.RespondError(c, http.StatusBadRequest, "voto deve ser 'Sim' ou 'Não'")
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			utils.RespondError(c, http.StatusUnauthorized, "usuário não autenticado")
			return
		}

		err = voteService.Vote(topicID, userID.(int), req.Choice)
		if err != nil {
			utils.RespondError(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondSuccess(c, nil)
	}
}

func ResultHandler(voteService vote.VoteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		topicID, err := strconv.Atoi(c.Param("topic_id"))
		if err != nil {
			utils.RespondError(c, http.StatusBadRequest, "topic_id inválido")
			return
		}
		yes, no, err := voteService.GetResult(topicID)
		if err != nil {
			utils.RespondError(c, http.StatusInternalServerError, err.Error())
			return
		}
		utils.RespondSuccess(c, gin.H{"Sim": yes, "Não": no})
	}
}
