package handlers

import (
	"desafio-tecnico-fullstack/backend/models"
	"desafio-tecnico-fullstack/backend/storage/repository"
	"desafio-tecnico-fullstack/backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	repo := c.MustGet("userRepository").(*repository.UserRepository)
	var req struct {
		Name     string `json:"name"`
		CPF      string `json:"cpf"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Erro": "Requisição inválida"})
		return
	}

	missingFields := []string{}
	if req.Name == "" {
		missingFields = append(missingFields, "Nome")
	}
	if req.CPF == "" {
		missingFields = append(missingFields, "CPF")
	}
	if req.Password == "" {
		missingFields = append(missingFields, "Senha")
	}
	if len(missingFields) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Erro": "Campos obrigatórios não preenchidos", "campos": missingFields})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := models.User{Name: req.Name, CPF: req.CPF, Password: string(hash)}
	err := repo.AddUser(user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{"Erro": "Usuário já existe"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"Erro": err.Error()})
		}
		return
	}
	c.Status(http.StatusCreated)
}

func Login(c *gin.Context) {
	repo := c.MustGet("userRepository").(*repository.UserRepository)
	var req struct {
		CPF      string `json:"cpf"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Erro": "Requisição inválida"})
		return
	}
	user := repo.GetUserByCPF(req.CPF)
	if user == nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Erro": "Requisição inválida"})
		return
	}
	token, _ := utils.GenerateJWT(user.CPF)
	c.JSON(http.StatusOK, gin.H{"token": token})
}
