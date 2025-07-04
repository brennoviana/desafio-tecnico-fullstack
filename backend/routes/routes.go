package routes

import (
	"desafio-tecnico-fullstack/backend/handlers/auth"
	sessionhandler "desafio-tecnico-fullstack/backend/handlers/session"
	topichandler "desafio-tecnico-fullstack/backend/handlers/topic"
	votehandler "desafio-tecnico-fullstack/backend/handlers/vote"
	"desafio-tecnico-fullstack/backend/middleware"
	sessionrepo "desafio-tecnico-fullstack/backend/storage/repository/session"
	topicrepo "desafio-tecnico-fullstack/backend/storage/repository/topic"
	userrepo "desafio-tecnico-fullstack/backend/storage/repository/user"
	voterepo "desafio-tecnico-fullstack/backend/storage/repository/vote"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	UserRepo    userrepo.UserRepository
	TopicRepo   topicrepo.TopicRepository
	SessionRepo sessionrepo.SessionRepository
	VoteRepo    voterepo.VoteRepository
}

func RegisterRoutes(router *gin.Engine, deps *Dependencies) {
	router.POST("/register", auth.RegisterHandler(deps.UserRepo))
	router.POST("/login", auth.LoginHandler(deps.UserRepo))

	router.POST("/topics", middleware.AuthMiddleware(), topichandler.CreateTopicHandler(deps.TopicRepo))
	router.GET("/topics", topichandler.ListTopicsHandler(deps.TopicRepo))
	router.POST("/topics/:topic_id/session", middleware.AuthMiddleware(), sessionhandler.OpenSessionHandler(deps.SessionRepo))
	router.POST("/topics/:topic_id/vote", middleware.AuthMiddleware(), votehandler.VoteHandler(deps.VoteRepo, deps.SessionRepo))
	router.GET("/topics/:topic_id/result", votehandler.ResultHandler(deps.VoteRepo, deps.SessionRepo))
}
