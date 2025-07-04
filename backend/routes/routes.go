package routes

import (
	"desafio-tecnico-fullstack/backend/handlers/auth"
	sessionhandler "desafio-tecnico-fullstack/backend/handlers/session"
	topichandler "desafio-tecnico-fullstack/backend/handlers/topic"
	votehandler "desafio-tecnico-fullstack/backend/handlers/vote"
	"desafio-tecnico-fullstack/backend/middleware"
	"desafio-tecnico-fullstack/backend/services/session"
	"desafio-tecnico-fullstack/backend/services/topic"
	"desafio-tecnico-fullstack/backend/services/user"
	"desafio-tecnico-fullstack/backend/services/vote"

	"github.com/gin-gonic/gin"
)

type Services struct {
	UserService    user.UserService
	TopicService   topic.TopicService
	SessionService session.SessionService
	VoteService    vote.VoteService
}

func RegisterRoutes(router *gin.Engine, deps *Services) {
	router.POST("/register", auth.RegisterHandler(deps.UserService))
	router.POST("/login", auth.LoginHandler(deps.UserService))

	router.POST("/topics", middleware.AuthMiddleware(), topichandler.CreateTopicHandler(deps.TopicService))
	router.GET("/topics", topichandler.ListTopicsHandler(deps.TopicService))
	router.POST("/topics/:topic_id/session", middleware.AuthMiddleware(), sessionhandler.OpenSessionHandler(deps.SessionService))
	router.POST("/topics/:topic_id/vote", middleware.AuthMiddleware(), votehandler.VoteHandler(deps.VoteService))
	router.GET("/topics/:topic_id/result", votehandler.ResultHandler(deps.VoteService))
}
