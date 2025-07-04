package vote

import (
	"desafio-tecnico-fullstack/backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func VoteHandler(voteService services.VoteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		topicID, err := strconv.Atoi(c.Param("topic_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Erro": "topic_id inválido"})
			return
		}
		var req struct {
			Choice string `json:"choice"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Erro": "Voto deve ser 'Sim' ou 'Não'"})
			return
		}
		cpf := c.GetString("cpf")
		err = voteService.Vote(topicID, cpf, req.Choice)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Erro": err.Error()})
			return
		}
		c.Status(http.StatusCreated)
	}
}

func ResultHandler(voteService services.VoteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		topicID, err := strconv.Atoi(c.Param("topic_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Erro": "topic_id inválido"})
			return
		}
		yes, no, err := voteService.GetResult(topicID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Erro": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Sim": yes, "Não": no})
	}
}
