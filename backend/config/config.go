package config

import (
	"log"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
}

var AppConfig *Config

func LoadConfig() {
	AppConfig = &Config{
		DBHost:     getEnv("POSTGRES_HOST", ""),
		DBPort:     getEnv("POSTGRES_PORT", ""),
		DBUser:     getEnv("POSTGRES_USER", ""),
		DBPassword: getEnv("POSTGRES_PASSWORD", ""),
		DBName:     getEnv("POSTGRES_DB", ""),
		JWTSecret:  getEnv("JWT_SECRET", ""),
	}
	log.Printf("Config loaded: host=%s port=%s user=%s db=%s", AppConfig.DBHost, AppConfig.DBPort, AppConfig.DBUser, AppConfig.DBName)
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
