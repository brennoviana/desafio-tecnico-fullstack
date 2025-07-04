package topic

import (
	"desafio-tecnico-fullstack/backend/models"
	topicrepo "desafio-tecnico-fullstack/backend/storage/repository/topic"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateTopicHandler(repo topicrepo.TopicRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name string `json:"name"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Name) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Erro": "Nome da pauta é obrigatório"})
			return
		}
		topic := models.Topic{Name: req.Name}
		err := repo.CreateTopic(topic)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Erro": err.Error()})
			return
		}
		c.Status(http.StatusCreated)
	}
}

func ListTopicsHandler(repo topicrepo.TopicRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		topics, err := repo.ListTopics()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Erro": err.Error()})
			return
		}
		c.JSON(http.StatusOK, topics)
	}
}
