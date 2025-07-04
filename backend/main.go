package main

import (
	"desafio-tecnico-fullstack/backend/config"
	"desafio-tecnico-fullstack/backend/routes"
	"desafio-tecnico-fullstack/backend/services"
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

	userRepository := userrepo.NewUserRepository(db)
	topicRepository := topicrepo.NewTopicRepository(db)
	sessionRepository := sessionrepo.NewSessionRepository(db)
	voteRepository := voterepo.NewVoteRepository(db)

	userService := services.NewUserService(userRepository)
	topicService := services.NewTopicService(topicRepository)
	sessionService := services.NewSessionService(sessionRepository)
	voteService := services.NewVoteService(voteRepository, sessionRepository)

	deps := &routes.Dependencies{
		UserService:    userService,
		TopicService:   topicService,
		SessionService: sessionService,
		VoteService:    voteService,
	}

	router := gin.Default()
	routes.RegisterRoutes(router, deps)
	router.Run(":8080")
}
