package main

import (
	"desafio-tecnico-fullstack/backend/config"
	"desafio-tecnico-fullstack/backend/routes"
	"desafio-tecnico-fullstack/backend/storage/connection"
	sessionrepo "desafio-tecnico-fullstack/backend/storage/repository/session"
	topicrepo "desafio-tecnico-fullstack/backend/storage/repository/topic"
	userrepo "desafio-tecnico-fullstack/backend/storage/repository/user"
	voterepo "desafio-tecnico-fullstack/backend/storage/repository/vote"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	db, err := connection.NewDB()
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}
	defer db.Close()

	deps := &routes.Dependencies{
		UserRepo:    userrepo.NewUserRepository(db),
		TopicRepo:   topicrepo.NewTopicRepository(db),
		SessionRepo: sessionrepo.NewSessionRepository(db),
		VoteRepo:    voterepo.NewVoteRepository(db),
	}

	router := gin.Default()
	routes.RegisterRoutes(router, deps)
	router.Run(":8080")
}
