package auth

import (
	"desafio-tecnico-fullstack/backend/services"
	"desafio-tecnico-fullstack/backend/utils"
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
			utils.RespondError(c, http.StatusBadRequest, "requisição inválida")
			return
		}

		if req.Name == "" || req.CPF == "" || req.Password == "" {
			utils.RespondError(c, http.StatusBadRequest, "campos obrigatórios não preenchidos")
			return
		}

		err := userService.RegisterUser(req.Name, req.CPF, req.Password)
		if err != nil {
			if err.Error() == "usuário já existe" {
				utils.RespondError(c, http.StatusConflict, err.Error())
			} else if err.Error() == "cpf inválido" || err.Error() == "senha muito curta" {
				utils.RespondError(c, http.StatusBadRequest, err.Error())
			} else {
				utils.RespondError(c, http.StatusInternalServerError, err.Error())
			}
			return
		}
		utils.RespondSuccess(c, nil)
	}
}

func LoginHandler(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			CPF      string `json:"cpf"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.RespondError(c, http.StatusBadRequest, "requisição inválida")
			return
		}
		token, err := userService.AuthenticateUser(req.CPF, req.Password)
		if err != nil {
			if err.Error() == "usuário ou senha inválidos" {
				utils.RespondError(c, http.StatusUnauthorized, err.Error())
			} else {
				utils.RespondError(c, http.StatusInternalServerError, err.Error())
			}
			return
		}
		utils.RespondSuccess(c, gin.H{"token": token})
	}
}
