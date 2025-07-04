package topic

import (
	"desafio-tecnico-fullstack/backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateTopicHandler(topicService services.TopicService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name string `json:"name"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Erro": "Requisição inválida"})
			return
		}
		if req.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Erro": "Nome da pauta é obrigatório"})
			return
		}
		err := topicService.CreateTopic(req.Name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Erro": err.Error()})
			return
		}
		c.Status(http.StatusCreated)
	}
}

func ListTopicsHandler(topicService services.TopicService) gin.HandlerFunc {
	return func(c *gin.Context) {
		topics, err := topicService.ListTopics()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Erro": err.Error()})
			return
		}
		c.JSON(http.StatusOK, topics)
	}
}
