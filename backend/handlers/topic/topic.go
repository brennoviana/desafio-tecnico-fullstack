package topic

import (
	"desafio-tecnico-fullstack/backend/services/topic"
	"desafio-tecnico-fullstack/backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateTopicHandler(topicService topic.TopicService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name string `json:"name"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.RespondError(c, http.StatusBadRequest, "Requisição inválida")
			return
		}
		if req.Name == "" {
			utils.RespondError(c, http.StatusBadRequest, "Nome da pauta é obrigatório")
			return
		}
		err := topicService.CreateTopic(req.Name)
		if err != nil {
			utils.RespondError(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondSuccess(c, nil)
	}
}

func ListTopicsHandler(topicService topic.TopicService) gin.HandlerFunc {
	return func(c *gin.Context) {
		topics, err := topicService.ListTopics()
		if err != nil {
			utils.RespondError(c, http.StatusInternalServerError, err.Error())
			return
		}
		utils.RespondSuccess(c, topics)
	}
}
