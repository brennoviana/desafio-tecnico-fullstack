package utils

import (
	"time"

	"desafio-tecnico-fullstack/backend/config"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})
	return token.SignedString([]byte(config.AppConfig.JWT.Secret))
}

func ValidateJWT(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.Secret), nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return 0, jwt.ErrTokenMalformed
		}
		return int(userID), nil
	}
	return 0, err
}
