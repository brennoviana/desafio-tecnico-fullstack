package routes

import (
	"desafio-tecnico-fullstack/backend/handlers/auth"
	sessionhandler "desafio-tecnico-fullstack/backend/handlers/session"
	topichandler "desafio-tecnico-fullstack/backend/handlers/topic"
	votehandler "desafio-tecnico-fullstack/backend/handlers/vote"
	"desafio-tecnico-fullstack/backend/middleware"
	"desafio-tecnico-fullstack/backend/services"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	UserService    services.UserService
	TopicService   services.TopicService
	SessionService services.SessionService
	VoteService    services.VoteService
}

func RegisterRoutes(router *gin.Engine, deps *Dependencies) {
	router.POST("/register", auth.RegisterHandler(deps.UserService))
	router.POST("/login", auth.LoginHandler(deps.UserService))

	router.POST("/topics", middleware.AuthMiddleware(), topichandler.CreateTopicHandler(deps.TopicService))
	router.GET("/topics", topichandler.ListTopicsHandler(deps.TopicService))
	router.POST("/topics/:topic_id/session", middleware.AuthMiddleware(), sessionhandler.OpenSessionHandler(deps.SessionService))
	router.POST("/topics/:topic_id/vote", middleware.AuthMiddleware(), votehandler.VoteHandler(deps.VoteService))
	router.GET("/topics/:topic_id/result", votehandler.ResultHandler(deps.VoteService))
}
