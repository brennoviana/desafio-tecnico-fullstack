package utils

import (
	"time"

	"desafio-tecnico-fullstack/backend/config"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(cpf string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"cpf": cpf,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

func ValidateJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWTSecret), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		cpf, ok := claims["cpf"].(string)
		if !ok {
			return "", jwt.ErrTokenMalformed
		}
		return cpf, nil
	}
	return "", err
}
