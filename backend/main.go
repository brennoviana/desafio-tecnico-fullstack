package main

import (
	"desafio-tecnico-fullstack/backend/handlers"
	"desafio-tecnico-fullstack/backend/middleware"
	"desafio-tecnico-fullstack/backend/storage/connection"
	"desafio-tecnico-fullstack/backend/storage/repository"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := connection.NewDB()
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)

	router := gin.Default()

	router.POST("/register", handlers.RegisterHandler(repo))
	router.POST("/login", handlers.LoginHandler(repo))

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/protected", func(c *gin.Context) {
		cpf := c.GetString("cpf")
		c.JSON(http.StatusOK, gin.H{"message": "Acesso autorizado", "cpf": cpf})
	})

	router.Run(":8080")
}
