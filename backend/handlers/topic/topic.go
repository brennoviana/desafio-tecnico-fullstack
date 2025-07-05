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
			Name   string `json:"name"`
			Status string `json:"status,omitempty"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			utils.RespondError(c, http.StatusBadRequest, "Requisição inválida")
			return
		}

		if req.Name == "" {
			utils.RespondError(c, http.StatusBadRequest, "Nome da pauta é obrigatório")
			return
		}

		// Validate status if provided
		if req.Status != "" {
			validStatuses := []string{"PENDING", "OPEN", "CLOSED", "FINISHED"}
			valid := false
			for _, status := range validStatuses {
				if req.Status == status {
					valid = true
					break
				}
			}
			if !valid {
				utils.RespondError(c, http.StatusBadRequest, "Status deve ser um dos seguintes: PENDING, OPEN, CLOSED, FINISHED")
				return
			}
		}

		err := topicService.CreateTopic(req.Name, req.Status)
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
