package main

import (
	"desafio-tecnico-fullstack/backend/config"
	"desafio-tecnico-fullstack/backend/routes"
	sessionService "desafio-tecnico-fullstack/backend/services/session"
	topicService "desafio-tecnico-fullstack/backend/services/topic"
	userService "desafio-tecnico-fullstack/backend/services/user"
	voteService "desafio-tecnico-fullstack/backend/services/vote"
	"desafio-tecnico-fullstack/backend/storage/connection"
	sessionRepo "desafio-tecnico-fullstack/backend/storage/repository/session"
	topicRepo "desafio-tecnico-fullstack/backend/storage/repository/topic"
	userRepo "desafio-tecnico-fullstack/backend/storage/repository/user"
	voteRepo "desafio-tecnico-fullstack/backend/storage/repository/vote"

	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	db, err := connection.NewDB()
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}
	defer db.Close()

	userRepository := userRepo.NewUserRepository(db)
	topicRepository := topicRepo.NewTopicRepository(db)
	sessionRepository := sessionRepo.NewSessionRepository(db)
	voteRepository := voteRepo.NewVoteRepository(db)

	userService := userService.NewUserService(userRepository)
	topicService := topicService.NewTopicService(topicRepository)
	sessionService := sessionService.NewSessionService(sessionRepository)
	voteService := voteService.NewVoteService(voteRepository, sessionRepository)

	deps := &routes.Services{
		UserService:    userService,
		TopicService:   topicService,
		SessionService: sessionService,
		VoteService:    voteService,
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	routes.RegisterRoutes(router, deps)
	router.Run(":8080")
}
