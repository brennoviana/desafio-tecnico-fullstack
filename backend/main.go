package main

import (
	authhandler "desafio-tecnico-fullstack/backend/handlers/auth"
	sessionhandler "desafio-tecnico-fullstack/backend/handlers/session"
	topichandler "desafio-tecnico-fullstack/backend/handlers/topic"
	votehandler "desafio-tecnico-fullstack/backend/handlers/vote"
	"desafio-tecnico-fullstack/backend/middleware"
	"desafio-tecnico-fullstack/backend/storage/connection"
	sessionrepo "desafio-tecnico-fullstack/backend/storage/repository/session"
	topicrepo "desafio-tecnico-fullstack/backend/storage/repository/topic"
	userrepo "desafio-tecnico-fullstack/backend/storage/repository/user"
	voterepo "desafio-tecnico-fullstack/backend/storage/repository/vote"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := connection.NewDB()
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}
	defer db.Close()

	repo := userrepo.NewUserRepository(db)
	topicRepo := topicrepo.NewTopicRepository(db)
	sessionRepo := sessionrepo.NewSessionRepository(db)
	voteRepo := voterepo.NewVoteRepository(db)

	router := gin.Default()

	router.POST("/register", authhandler.RegisterHandler(repo))
	router.POST("/login", authhandler.LoginHandler(repo))
	router.POST("/topics", middleware.AuthMiddleware(), topichandler.CreateTopicHandler(topicRepo))
	router.GET("/topics", topichandler.ListTopicsHandler(topicRepo))
	router.POST("/topics/:topic_id/session", middleware.AuthMiddleware(), sessionhandler.OpenSessionHandler(sessionRepo))
	router.POST("/topics/:topic_id/vote", middleware.AuthMiddleware(), votehandler.VoteHandler(voteRepo, sessionRepo))
	router.GET("/topics/:topic_id/result", votehandler.ResultHandler(voteRepo, sessionRepo))

	router.Run(":8080")
}
