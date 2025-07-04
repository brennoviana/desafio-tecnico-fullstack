package auth

import (
	"desafio-tecnico-fullstack/backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name     string `json:"name"`
			CPF      string `json:"cpf"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Erro": "Requisição inválida"})
			return
		}

		if req.Name == "" || req.CPF == "" || req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Erro": "Campos obrigatórios não preenchidos"})
			return
		}

		err := userService.RegisterUser(req.Name, req.CPF, req.Password)
		if err != nil {
			if err.Error() == "usuário já existe" {
				c.JSON(http.StatusConflict, gin.H{"Erro": err.Error()})
			} else if err.Error() == "cpf inválido" || err.Error() == "senha muito curta" {
				c.JSON(http.StatusBadRequest, gin.H{"Erro": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"Erro": err.Error()})
			}
			return
		}
		c.Status(http.StatusCreated)
	}
}

func LoginHandler(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			CPF      string `json:"cpf"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Erro": "Requisição inválida"})
			return
		}
		token, err := userService.AuthenticateUser(req.CPF, req.Password)
		if err != nil {
			if err.Error() == "usuário ou senha inválidos" {
				c.JSON(http.StatusUnauthorized, gin.H{"Erro": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"Erro": err.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
