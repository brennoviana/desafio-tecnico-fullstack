package handlers

import (
	"database/sql"
	"desafio-tecnico-fullstack/backend/models"
	"desafio-tecnico-fullstack/backend/storage"
	"desafio-tecnico-fullstack/backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	var req struct {
		Name     string `json:"name"`
		CPF      string `json:"cpf"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Erro": "Requisição inválida"})
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := models.User{Name: req.Name, CPF: req.CPF, Password: string(hash)}
	err := storage.AddUser(db, user)
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
	db := c.MustGet("db").(*sql.DB)
	var req struct {
		CPF      string `json:"cpf"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Erro": "Requisição inválida"})
		return
	}
	user := storage.GetUserByCPF(db, req.CPF)
	if user == nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Erro": "Requisição inválida"})
		return
	}
	token, _ := utils.GenerateJWT(user.CPF)
	c.JSON(http.StatusOK, gin.H{"token": token})
}
