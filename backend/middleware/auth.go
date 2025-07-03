package middleware

import (
    "desafio-tecnico-fullstack/backend/utils"
    "github.com/gin-gonic/gin"
    "net/http"
    "strings"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if !strings.HasPrefix(authHeader, "Bearer ") {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }
        token := strings.TrimPrefix(authHeader, "Bearer ")
        cpf, err := utils.ValidateJWT(token)
        if err != nil {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }
        c.Set("cpf", cpf)
        c.Next()
    }
} 